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
	"strconv"
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

	cc := ChiaCollector{
		client:  client,
		baseURL: *url,
	}
	prometheus.MustRegister(cc)

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

type ChiaCollector struct {
	client  *http.Client
	baseURL string
}

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc ChiaCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect queries Chia and returns metrics on ch.
func (cc ChiaCollector) Collect(ch chan<- prometheus.Metric) {
	cc.collectConnections(ch)
}

func (cc ChiaCollector) collectConnections(ch chan<- prometheus.Metric) {
	var conns Connections
	if err := queryAPI(cc.client, cc.baseURL, "get_connections", "", &conns); err != nil {
		log.Print(err)
		return
	}
	peers := make([]int, NumNodeTypes)
	for _, p := range conns.Connections {
		peers[p.Type-1]++
	}
	desc := prometheus.NewDesc(
		"chia_peers_count",
		"Number of peers currently connected.",
		[]string{"type"}, nil,
	)
	for nt, cnt := range peers {
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(cnt),
			strconv.Itoa(nt+1),
		)
	}
}
