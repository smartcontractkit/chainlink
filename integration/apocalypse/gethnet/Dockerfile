FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=ethereum/client-go:stable /usr/local/bin/geth /usr/local/bin/

WORKDIR /gethnet
COPY datadir datadir/
COPY config.toml .
COPY apocalypse.json .
RUN geth init --datadir /gethnet/datadir /gethnet/apocalypse.json

EXPOSE 8545 8546 8547 30303 30303/udp
CMD geth --networkid 1337 --datadir /gethnet/datadir --mine --miner.threads=1 --nodekeyhex $NODE_KEY --ipcdisable --dev.period 2 --config /gethnet/config.toml

