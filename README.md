# chia_exporter

[Prometheus](https://prometheus.io) metric collector for
[Chia](https://chia.net) nodes, using the local [RPC
API](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces)

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
```

### Blockchain

Various node and blockchain metrics are collected from the
[get_blockchain_state](https://github.com/Chia-Network/chia-blockchain/wiki/RPC-Interfaces#get_blockchain_state)
endpoint.

### Connections

The number of connections are collected for each node type from the
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
