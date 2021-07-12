# chia_exporter

[Prometheus](https://prometheus.io) metric collector for
[Chia](https://chia.net) nodes, using the local [RPC
API](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces)

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
```

### Blockchain

Various node and blockchain metrics are collected from the
[get_blockchain_state](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_blockchain_state)
endpoint.

### Connections

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

* Pool state is collected from the
  [get_pool_state](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_pool_state)
  endpoint (not yet documented). Need chia client version 1.2.0 or later

