FROM ethereum/client-go

RUN apk add curl bash

COPY ./genesis.json /root
COPY ./geth-config.toml /root

# create genesis block
RUN geth --nousb --config /root/geth-config.toml init /root/genesis.json

EXPOSE 30303
EXPOSE 30303/udp

