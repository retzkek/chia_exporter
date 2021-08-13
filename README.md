# chia_exporter

[Prometheus](https://prometheus.io) metric collector for
[Chia](https://chia.net) nodes, using the local [RPC
API](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces)

## Quick Install

1. Get the latest tarball for your platform from the [release
   page](https://github.com/retzkek/chia_exporter/releases)

2. Extract the tarball, move the `chia_exporter` binary to `/usr/bin` and move
   the `chia-exporter@.service` unit file to `/etc/systemd/system`.

``` sh
tar xzvf chia_exporter*.tar.gz
sudo mv chia_exporter /usr/bin
sudo mv chia-exporter@.service /etc/systemd/system/
```

3. Start the exporter, substituting your username (or the name of the user that
   chia runs under) where it says MYUSERNAME below. Optionally enable it to
   start at boot.

``` sh
sudo systemctl start chia-exporter@MYUSERNAME.service
sudo systemctl enable chia-exporter@MYUSERNAME.service
```

4. Check that it's running.

``` sh
sudo systemctl status chia-exporter@MYUSERNAME.service
```

### Troubleshooting

If it fails to start and says it can't find the certificate or key: check if the
chia SSL certificate and key (needed to talk to the APIs) are in
`/home/MYUSERNAME/.chia/mainnet/config/ssl/full_node`, which is the default
location. If they're somewhere else you'll need to modify the service file. If
you're running the exporter as a different user than chia then you need to
modify the service file and make sure the user can access the key and cert.

If it says it can't reach one or more chia daemons: if you're not running all
the daemons, you can safely ignore these warnings. Otherwise you may need to
update the daemon URLs, see configuration options below.

## Building and Running

With the [Go](http://golang.org) compiler tools installed:

    go build

Run `./chia_exporter -h` to see the command configuration options:

    -cert string
          The full node SSL certificate. (default "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.crt")
    -key string
          The full node SSL key. (default "$HOME/.chia/mainnet/config/ssl/full_node/private_full_node.key")
    -listen string
          The address to listen on for HTTP requests. (default ":9133")
    -url string
          The base URL for the full node RPC endpoint. (default "https://localhost:8555")
    -wallet string
          The base URL for the wallet RPC endpoint. (default "https://localhost:9256")
    -harvester string
          The base URL for the harvester RPC endpoint. (default "https://localhost:8560")
    -farmer string
          The base URL for the farmer RPC endpoint. (default "https://localhost:8559")

## Metrics

Example of all metrics currently exposed:

``` sh
# HELP chia_blockchain_difficulty Current difficulty
# TYPE chia_blockchain_difficulty gauge
chia_blockchain_difficulty 112
# HELP chia_blockchain_height Current height
# TYPE chia_blockchain_height gauge
chia_blockchain_height 221609
# HELP chia_blockchain_space_bytes Estimated current netspace
# TYPE chia_blockchain_space_bytes gauge
chia_blockchain_space_bytes 1.8771214186533368e+18
# HELP chia_blockchain_sync_status Sync status, 0=not synced, 1=syncing, 2=synced
# TYPE chia_blockchain_sync_status gauge
chia_blockchain_sync_status 2
# HELP chia_blockchain_total_iters Current total iterations
# TYPE chia_blockchain_total_iters gauge
chia_blockchain_total_iters 7.20695891692e+11
# HELP chia_peers_count Number of peers currently connected.
# TYPE chia_peers_count gauge
chia_peers_count{type="1"} 52
chia_peers_count{type="2"} 0
chia_peers_count{type="3"} 1
chia_peers_count{type="4"} 0
chia_peers_count{type="5"} 0
chia_peers_count{type="6"} 1
# HELP chia_wallet_confirmed_balance_mojo Confirmed wallet balance.
# TYPE chia_wallet_confirmed_balance_mojo gauge
chia_wallet_confirmed_balance_mojo{wallet_id="1",wallet_fingerprint="103402894"} 100
# HELP chia_wallet_height Wallet synced height.
# TYPE chia_wallet_height gauge
chia_wallet_height{wallet_id="1",wallet_fingerprint="103402894"} 30756
# HELP chia_wallet_max_send_mojo Maximum sendable amount.
# TYPE chia_wallet_max_send_mojo gauge
chia_wallet_max_send_mojo{wallet_id="1",wallet_fingerprint="103402894"} 100
# HELP chia_wallet_pending_change_mojo Pending change amount.
# TYPE chia_wallet_pending_change_mojo gauge
chia_wallet_pending_change_mojo{wallet_id="1",wallet_fingerprint="103402894"} 0
# HELP chia_wallet_spendable_balance_mojo Spendable wallet balance.
# TYPE chia_wallet_spendable_balance_mojo gauge
chia_wallet_spendable_balance_mojo{wallet_id="1",wallet_fingerprint="103402894"} 100
# HELP chia_wallet_sync_status Sync status, 0=not synced, 1=syncing, 2=synced
# TYPE chia_wallet_sync_status gauge
chia_wallet_sync_status{wallet_id="1",wallet_fingerprint="103402894"} 0
# HELP chia_wallet_unconfirmed_balance_mojo Unconfirmed wallet balance.
# TYPE chia_wallet_unconfirmed_balance_mojo gauge
chia_wallet_unconfirmed_balance_mojo{wallet_id="1",wallet_fingerprint="103402894"} 100
# HELP chia_wallet_farmed_amount Farmed amount
# TYPE chia_wallet_farmed_amount gauge
chia_wallet_farmed_amount{wallet_fingerprint="103402894",wallet_id="1"} 0
# HELP chia_wallet_fee_amount Fee amount amount
# TYPE chia_wallet_fee_amount gauge
chia_wallet_fee_amount{wallet_fingerprint="103402894",wallet_id="1"} 0
# HELP chia_wallet_last_height_farmed Last height farmed
# TYPE chia_wallet_last_height_farmed gauge
chia_wallet_last_height_farmed{wallet_fingerprint="103402894",wallet_id="1"} 0
# HELP chia_wallet_pool_reward_amount Pool Reward amount
# TYPE chia_wallet_pool_reward_amount gauge
chia_wallet_pool_reward_amount{wallet_fingerprint="103402894",wallet_id="1"} 0
# HELP chia_wallet_reward_amount Reward amount
# TYPE chia_wallet_reward_amount gauge
chia_wallet_reward_amount{wallet_fingerprint="103402894",wallet_id="1"} 0
# HELP chia_pool_current_difficulty Current difficulty on pool.
# TYPE chia_pool_current_difficulty gauge
chia_pool_current_difficulty{launcher_id="0x...",pool_url="https://pool.yyy.y"} 1
# HELP chia_pool_current_points Current points on pool.
# TYPE chia_pool_current_points gauge
chia_pool_current_points{launcher_id="0x...",pool_url="https://pool.yyy.y"} 12
# HELP chia_pool_points_acknowledged_24h Points acknowledged last 24h on pool.
# TYPE chia_pool_points_acknowledged_24h gauge
chia_pool_points_acknowledged_24h{launcher_id="0x...",pool_url="https://pool.yyy.y"} 5
# HELP chia_pool_points_found_24h Points found last 24h on pool.
# TYPE chia_pool_points_found_24h gauge
chia_pool_points_found_24h{launcher_id="0x...",pool_url="https://pool.xchpool.org"} 5
# HELP chia_plots Number of plots currently using.
# TYPE chia_plots gauge
chia_plots 54
# HELP chia_plots_failed_to_open Number of plots files failed to open.
# TYPE chia_plots_failed_to_open gauge
chia_plots_failed_to_open 0
# HELP chia_plots_not_found Number of plots files not found.
# TYPE chia_plots_not_found gauge
chia_plots_not_found 0
```

### Blockchain and Connections (full node)

Various node and blockchain metrics are collected from the
[get_blockchain_state](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_blockchain_state)
endpoint.

* The number of connections are collected for each node type from the
  [get_connections](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_connections)
  endpoint.

Node types (from
[chia/server/outbound_message.py](https://github.com/Chia-Network/chia-blockchain/blob/main/chia/server/outbound_message.py#L10)):

    FULL_NODE = 1
    HARVESTER = 2
    FARMER = 3
    TIMELORD = 4
    INTRODUCER = 5
    WALLET = 6

### Wallet

The list of wallets is obtained from the
[get_wallets](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_wallets)
endpoint. The wallet metrics are collected for each wallet, and include
`wallet_id` and `wallet_fingerprint` labels.

* Balances are collected from the
  [get_wallet_balance](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_wallet_balance)
  endpoint.

* Sync status is collected from the
  [get_sync_status](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_sync_status)
  endpoint.

* Height is collected from the
  [get_height_info](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_height_info)
  endpoint.

* Farmed ammount and reward are collected from the
  [get_farmed_amount](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_farmed_amount)

### Pool (farmer)

* Pool state is collected from the
  [get_pool_state](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_pool_state)
  endpoint (not yet documented). Need chia client version 1.2.0 or later

### Plots (harvester)

* Plots data are collected from the
  [get_plots](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_plots)
  endpoint.

