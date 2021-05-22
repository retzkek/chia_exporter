FROM golang:alpine AS builder

WORKDIR /build
COPY . /build/chia_exporter
RUN apk add --update --no-cache --virtual build-dependencies \
 && cd chia_exporter \
 && go build -tags netgo

FROM alpine
COPY --from=builder /build/chia_exporter/chia_exporter /usr/bin/chia_exporter

EXPOSE 9133

ENV FULL_NODE_CERT=/chia_exporter/private_full_node.crt
ENV FULL_NODE_KEY=/chia_exporter/private_full_node.key
ENV FULL_NODE_RPC_ENDPOINT=https://localhost:8555
ENV WALLET_RPC_ENDPOINT=https://localhost:9256

CMD /usr/bin/chia_exporter -cert $FULL_NODE_CERT -key $FULL_NODE_KEY -url $FULL_NODE_RPC_ENDPOINT -wallet $WALLET_RPC_ENDPOINT
