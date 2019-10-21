FROM ethereum/client-go

RUN apk add curl bash

COPY root /root

# some handy commands for messing with geth and chain data
RUN chmod +x /root/bin/*
ENV PATH="$PATH:/root/bin"

# create genesis block
RUN geth --nousb --config /root/geth-config.toml init /root/genesis.json

EXPOSE 30303
EXPOSE 30303/udp

