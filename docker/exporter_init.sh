#!/usr/bin/env bash

# build param stack
args=()

if [ -z "${FULL_NODE_RPC_ENDPOINT}" ]; then
  >&2 echo "FULL_NODE_RPC_ENDPOINT not set, feature disabled."
else
  args+=("-url=$FULL_NODE_RPC_ENDPOINT")
fi

if [ -z "${WALLET_RPC_ENDPOINT}" ]; then
  >&2 echo "WALLET_RPC_ENDPOINT not set, feature disabled."
else
  args+=("-wallet=$WALLET_RPC_ENDPOINT")
fi

if [ -z "${FARMER_RPC_ENDPOINT}" ]; then
  >&2 echo "FARMER_RPC_ENDPOINT not set, feature disabled."
else
  args+=("-farmer=$FARMER_RPC_ENDPOINT")
fi

if [ -z "${HARVESTER_RPC_ENDPOINT}" ]; then
  >&2 echo "HARVESTER_RPC_ENDPOINT not set, feature disabled."
else
  args+=("-harvester=$HARVESTER_RPC_ENDPOINT")
fi

/usr/bin/chia_exporter \
  -cert "$CERT" -key "$KEY" \
  "${args[@]}" &

# Finally, exec the normal Chia docker-start
exec /usr/local/bin/docker-start.sh
