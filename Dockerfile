FROM golang:alpine AS builder

WORKDIR /build
COPY . /build/chia_exporter
RUN apk add --update --no-cache --virtual build-dependencies \
 && cd chia_exporter \
 && go build -tags netgo

FROM ghcr.io/chia-network/chia:1.2.11
COPY --from=builder /build/chia_exporter/chia_exporter /usr/bin/chia_exporter
COPY docker/exporter_init.sh /usr/local/bin/exporter_init.sh
EXPOSE 9133

ENV CERT=/root/.chia/mainnet/config/ssl/full_node/private_full_node.crt
ENV KEY=/root/.chia/mainnet/config/ssl/full_node/private_full_node.key
ENV FULL_NODE_RPC_ENDPOINT=https://localhost:8555
ENV WALLET_RPC_ENDPOINT=https://localhost:9256
ENV FARMER_RPC_ENDPOINT=https://localhost:8559
ENV HARVESTER_RPC_ENDPOINT=https://localhost:8560

CMD /usr/local/bin/exporter_init.sh "$@"
