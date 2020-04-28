FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=ethereum/client-go:stable /usr/local/bin/geth /usr/local/bin/

WORKDIR /gethnet
COPY datadir datadir/
COPY config.toml .
COPY apocalypse.json .
RUN geth init --datadir /gethnet/datadir /gethnet/apocalypse.json

EXPOSE 8545 8546 8547 30303 30303/udp
CMD geth --networkid 1337 \
         --datadir /gethnet/datadir \
         --mine \
         --miner.threads 1 \
         --miner.etherbase $ETHERBASE \
         --nodekeyhex $NODE_KEY \
         --ipcdisable \
         --dev.period 2 \
         --config /gethnet/config.toml \
         --rpc --rpcaddr 0.0.0.0 \
         --rpcport 8545 \
         --rpcapi eth,web3,net,admin,db,txpool \
         --ws \
         --wsaddr 0.0.0.0 \
         --wsport 8546 \
         --wsapi eth,web3,net,admin,db,txpool \
         --bootnodes enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.3:30303
