package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen", ":9133", "The address to listen on for HTTP requests.")
	cert = flag.String("cert", "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.crt", "The full node SSL certificate.")
	key  = flag.String("key", "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.key", "The full node SSL key.")
	url  = flag.String("url", "https://localhost:8555", "The base URL for the full node RPC endpoint.")
)

func main() {
	flag.Parse()

	client, err := newClient(os.ExpandEnv(*cert), os.ExpandEnv(*key))
	if err != nil {
		log.Fatal(err)
	}
	var info NetworkInfo
	if err := queryAPI(client, *url, "get_network_info", "", &info); err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to node at %s on %s", *url, info.NetworkName)

	// TODO: add labels for node type (can't use NewGuageFunc, need a collector)
	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "chia",
			Name:      "peers_count",
			Help:      "Number of peers currently connected.",
		},
		peerCounter(client, *url),
	)); err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s. Serving metrics on /metrics.", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func newClient(cert, key string) (*http.Client, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{c},
				InsecureSkipVerify: true,
			},
		},
		Timeout: 5 * time.Second,
	}, nil
}

func queryAPI(client *http.Client, base, endpoint, query string, result interface{}) error {
	if query == "" {
		query = `{"":""}`
	}
	b := strings.NewReader(query)
	r, err := client.Post(base+"/"+endpoint, "application/json", b)
	if err != nil {
		return fmt.Errorf("error calling %s: %w", endpoint, err)
	}
	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		if err != nil {
			return fmt.Errorf("error decoding %s response: %w", endpoint, err)
		}
	}
	return nil
}

func peerCounter(client *http.Client, base string) func() float64 {
	return func() float64 {
		var conns Connections
		if err := queryAPI(client, base, "get_connections", "", &conns); err != nil {
			log.Print(err)
			return -1.0
		}
		peers := 0
		for _, p := range conns.Connections {
			if p.Type == NodeTypeFullNode {
				peers++
			}
		}
		return float64(peers)
	}
}
