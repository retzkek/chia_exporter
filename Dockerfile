FROM golang:latest

EXPOSE 9133

ENV FULL_NODE_CERT=/chia_exporter/private_full_node.crt
ENV FULL_NODE_KEY=/chia_exporter/private_full_node.key
ENV FULL_NODE_RPC_ENDPOINT=https://localhost:8555
ENV WALLET_RPC_ENDPOINT=https://localhost:9256

RUN mkdir /chia_exporter

WORKDIR /chia_exporter

COPY . ./

RUN go build

CMD ./chia_exporter -cert $FULL_NODE_CERT -key $FULL_NODE_KEY -url $FULL_NODE_RPC_ENDPOINT -wallet $WALLET_RPC_ENDPOINT
