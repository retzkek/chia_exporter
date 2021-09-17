// Copyright 2021 Kevin Retzke
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	addr      = flag.String("listen", ":9133", "The address to listen on for HTTP requests.")
	cert      = flag.String("cert", "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.crt", "The full node SSL certificate.")
	key       = flag.String("key", "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.key", "The full node SSL key.")
	url       = flag.String("url", "https://localhost:8555", "The base URL for the full node RPC endpoint.")
	wallet    = flag.String("wallet", "https://localhost:9256", "The base URL for the wallet RPC endpoint.")
	farmer    = flag.String("farmer", "https://localhost:8559", "The base URL for the farmer RPC endpoint.")
	harvester = flag.String("harvester", "https://localhost:8560", "The base URL for the harvester RPC endpoint.")
	timeout   = flag.String("timeout", "5s", "HTTP client timeout per request, as duration string.")
)

var (
	Version = "0.5.2"
)

func main() {
	log.Printf("chia_exporter version %s", Version)
	flag.Parse()

	client, err := newClient(os.ExpandEnv(*cert), os.ExpandEnv(*key))
	if err != nil {
		log.Fatal(err)
	}
	var info NetworkInfo
	if err := queryAPI(client, *url, "get_network_info", "", &info); err != nil {
		log.Print(err)
	} else {
		log.Printf("Connected to node at %s on %s", *url, info.NetworkName)
	}

	cc := ChiaCollector{
		client:       client,
		baseURL:      *url,
		walletURL:    *wallet,
		farmerURL:    *farmer,
		harvesterURL: *harvester,
	}
	prometheus.MustRegister(cc)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "chia_exporter version %s\n", Version)
		fmt.Fprintf(w, "metrics are published on /metrics\n\n")
		fmt.Fprintf(w, "This program is free software released under the GNU AGPL.\n")
		fmt.Fprintf(w, "The source code is availabe at https://github.com/retzkek/chia_exporter\n")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %s. Serving metrics on /metrics.", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func newClient(cert, key string) (*http.Client, error) {
	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	to, err := time.ParseDuration(*timeout)
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
		Timeout: to,
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
	//t := io.TeeReader(r.Body, os.Stdout)
	t := io.TeeReader(r.Body, ioutil.Discard)
	if err := json.NewDecoder(t).Decode(result); err != nil {
		if err != nil {
			return fmt.Errorf("error decoding %s response: %w", endpoint, err)
		}
	}
	return nil
}

type ChiaCollector struct {
	client       *http.Client
	baseURL      string
	walletURL    string
	farmerURL    string
	harvesterURL string
}

// Describe is implemented with DescribeByCollect.
func (cc ChiaCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect queries Chia and returns metrics on ch.
func (cc ChiaCollector) Collect(ch chan<- prometheus.Metric) {
	cc.collectConnections(ch)
	cc.collectBlockchainState(ch)
	cc.collectWallets(ch)
	cc.collectPoolState(ch)
	cc.collectPlots(ch)
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

func (cc ChiaCollector) collectBlockchainState(ch chan<- prometheus.Metric) {
	var bs BlockchainState
	if err := queryAPI(cc.client, cc.baseURL, "get_blockchain_state", "", &bs); err != nil {
		log.Print(err)
		return
	}
	sync := 0.0
	if bs.BlockchainState.Sync.SyncMode {
		sync = 1.0
	} else if bs.BlockchainState.Sync.Synced {
		sync = 2.0
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_blockchain_sync_status",
			"Sync status, 0=not synced, 1=syncing, 2=synced",
			nil, nil,
		),
		prometheus.GaugeValue,
		sync,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_blockchain_height",
			"Current height",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(bs.BlockchainState.Peak.Height),
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_blockchain_difficulty",
			"Current difficulty",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(bs.BlockchainState.Difficulty),
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_blockchain_space_bytes",
			"Estimated current netspace",
			nil, nil,
		),
		prometheus.GaugeValue,
		bs.BlockchainState.Space,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_blockchain_total_iters",
			"Current total iterations",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(bs.BlockchainState.Peak.TotalIters),
	)
}

func (cc ChiaCollector) collectWallets(ch chan<- prometheus.Metric) {
	var ws Wallets
	if err := queryAPI(cc.client, cc.walletURL, "get_wallets", "", &ws); err != nil {
		log.Print(err)
		return
	}
	for _, w := range ws.Wallets {
		w.StringID = strconv.Itoa(w.ID)
		w.PublicKey = cc.getWalletPublicKey(w)
		cc.collectWalletBalance(ch, w)
		cc.collectWalletSync(ch, w)
		cc.collectFarmedAmount(ch, w)
	}
}

// getWalletPublicKey returns the fingerprint of first public key associated
// with the wallet.
func (cc ChiaCollector) getWalletPublicKey(w Wallet) string {
	var wpks WalletPublicKeys
	q := fmt.Sprintf(`{"wallet_id":%d}`, w.ID)
	if err := queryAPI(cc.client, cc.walletURL, "get_public_keys", q, &wpks); err != nil {
		log.Print(err)
		return ""
	}
	if len(wpks.PublicKeyFingerprints) < 1 {
		log.Print("no public key")
		return ""
	}
	if len(wpks.PublicKeyFingerprints) > 1 {
		log.Print("more than one public key; returning first")
	}
	return strconv.Itoa(wpks.PublicKeyFingerprints[0])
}

var (
	confirmedBalanceDesc = prometheus.NewDesc(
		"chia_wallet_confirmed_balance_mojo",
		"Confirmed wallet balance.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
	unconfirmedBalanceDesc = prometheus.NewDesc(
		"chia_wallet_unconfirmed_balance_mojo",
		"Unconfirmed wallet balance.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
	spendableBalanceDesc = prometheus.NewDesc(
		"chia_wallet_spendable_balance_mojo",
		"Spendable wallet balance.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
	maxSendDesc = prometheus.NewDesc(
		"chia_wallet_max_send_mojo",
		"Maximum sendable amount.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
	pendingChangeDesc = prometheus.NewDesc(
		"chia_wallet_pending_change_mojo",
		"Pending change amount.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
)

func (cc ChiaCollector) collectWalletBalance(ch chan<- prometheus.Metric, w Wallet) {
	var wb WalletBalance
	q := fmt.Sprintf(`{"wallet_id":%d}`, w.ID)
	if err := queryAPI(cc.client, cc.walletURL, "get_wallet_balance", q, &wb); err != nil {
		log.Print(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		confirmedBalanceDesc,
		prometheus.GaugeValue,
		float64(wb.WalletBalance.ConfirmedBalance),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		unconfirmedBalanceDesc,
		prometheus.GaugeValue,
		float64(wb.WalletBalance.UnconfirmedBalance),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		spendableBalanceDesc,
		prometheus.GaugeValue,
		float64(wb.WalletBalance.SpendableBalance),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		maxSendDesc,
		prometheus.GaugeValue,
		float64(wb.WalletBalance.MaxSendAmount),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		pendingChangeDesc,
		prometheus.GaugeValue,
		float64(wb.WalletBalance.PendingChange),
		w.StringID, w.PublicKey,
	)
}

var (
	walletSyncStatusDesc = prometheus.NewDesc(
		"chia_wallet_sync_status",
		"Sync status, 0=not synced, 1=syncing, 2=synced",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
	walletHeightDesc = prometheus.NewDesc(
		"chia_wallet_height",
		"Wallet synced height.",
		[]string{"wallet_id", "wallet_fingerprint"}, nil,
	)
)

func (cc ChiaCollector) collectWalletSync(ch chan<- prometheus.Metric, w Wallet) {
	var wss WalletSyncStatus
	q := fmt.Sprintf(`{"wallet_id":%d}`, w.ID)
	if err := queryAPI(cc.client, cc.walletURL, "get_sync_status", q, &wss); err != nil {
		log.Print(err)
		return
	}
	sync := 0.0
	if wss.Syncing {
		sync = 1.0
	} else if wss.Synced {
		sync = 2.0
	}
	ch <- prometheus.MustNewConstMetric(
		walletSyncStatusDesc,
		prometheus.GaugeValue,
		sync,
		w.StringID, w.PublicKey,
	)

	var whi WalletHeightInfo
	if err := queryAPI(cc.client, cc.walletURL, "get_height_info", q, &whi); err != nil {
		log.Print(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		walletHeightDesc,
		prometheus.GaugeValue,
		float64(whi.Height),
		w.StringID, w.PublicKey,
	)
}

func (cc ChiaCollector) collectPoolState(ch chan<- prometheus.Metric) {
	var pools PoolState
	if err := queryAPI(cc.client, cc.farmerURL, "get_pool_state", "", &pools); err != nil {
		log.Print(err)
		return
	}
	for _, p := range pools.PoolState {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"chia_pool_current_difficulty",
				"Current difficulty on pool.",
				[]string{"launcher_id", "pool_url"}, nil,
			),
			prometheus.GaugeValue,
			float64(p.CurrentDificulty),
			p.PoolConfig.LauncherId,
			p.PoolConfig.PoolURL,
		)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"chia_pool_current_points",
				"Current points on pool.",
				[]string{"launcher_id", "pool_url"}, nil,
			),
			prometheus.GaugeValue,
			float64(p.CurrentPoints),
			p.PoolConfig.LauncherId,
			p.PoolConfig.PoolURL,
		)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"chia_pool_points_acknowledged_24h",
				"Points acknowledged last 24h on pool.",
				[]string{"launcher_id", "pool_url"}, nil,
			),
			prometheus.GaugeValue,
			float64(len(p.PointsAcknowledged24h)),
			p.PoolConfig.LauncherId,
			p.PoolConfig.PoolURL,
		)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"chia_pool_points_found_24h",
				"Points found last 24h on pool.",
				[]string{"launcher_id", "pool_url"}, nil,
			),
			prometheus.GaugeValue,
			float64(len(p.PointsFound24h)),
			p.PoolConfig.LauncherId,
			p.PoolConfig.PoolURL,
		)
	}
}

func (cc ChiaCollector) collectPlots(ch chan<- prometheus.Metric) {
	var plots PlotFiles
	if err := queryAPI(cc.client, cc.harvesterURL, "get_plots", "", &plots); err != nil {
		log.Print(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_plots_failed_to_open",
			"Number of plots files failed to open.",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(len(plots.FailedToOpen)),
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_plots_not_found",
			"Number of plots files not found.",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(len(plots.NotFound)),
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_plots",
			"Number of plots currently using.",
			nil, nil,
		),
		prometheus.GaugeValue,
		float64(len(plots.Plots)),
	)
}

func (cc ChiaCollector) collectFarmedAmount(ch chan<- prometheus.Metric, w Wallet) {
	var farmed FarmedAmount
	q := fmt.Sprintf(`{"wallet_id":%d}`, w.ID)
	if err := queryAPI(cc.client, cc.walletURL, "get_farmed_amount", q, &farmed); err != nil {
		log.Print(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_wallet_farmed_amount",
			"Farmed amount",
			[]string{"wallet_id", "wallet_fingerprint"}, nil,
		),
		prometheus.GaugeValue,
		float64(farmed.FarmedAmount),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_wallet_reward_amount",
			"Reward amount",
			[]string{"wallet_id", "wallet_fingerprint"}, nil,
		),
		prometheus.GaugeValue,
		float64(farmed.RewardAmount),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_wallet_fee_amount",
			"Fee amount amount",
			[]string{"wallet_id", "wallet_fingerprint"}, nil,
		),
		prometheus.GaugeValue,
		float64(farmed.FeeAmount),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_wallet_last_height_farmed",
			"Last height farmed",
			[]string{"wallet_id", "wallet_fingerprint"}, nil,
		),
		prometheus.GaugeValue,
		float64(farmed.LastHeightFarmed),
		w.StringID, w.PublicKey,
	)
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"chia_wallet_pool_reward_amount",
			"Pool Reward amount",
			[]string{"wallet_id", "wallet_fingerprint"}, nil,
		),
		prometheus.GaugeValue,
		float64(farmed.PoolRewardAmount),
		w.StringID, w.PublicKey,
	)
}
