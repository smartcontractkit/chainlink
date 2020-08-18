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
         --miner.gastarget 8000000 \
         --miner.gaslimit 8000000 \
         --miner.recommit 1s \
         --miner.noverify \
         --fakepow \
         --gcmode archive \
         --nodekeyhex $NODE_KEY \
         --ipcdisable \
         # --dev.period 2 \
         --config /gethnet/config.toml \
         --http \
         --http.addr 0.0.0.0 \
         --http.port 8545 \
         --http.api eth,web3,net,admin,debug,txpool \
         --ws \
         --ws.addr 0.0.0.0 \
         --ws.port 8546 \
         --ws.api eth,web3,net,admin,debug,txpool \
         --bootnodes enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.4:30303,enode://9f87a48e17c61f00d1dbbfe3692d890359ccb811d10ef21517913aeb05f47947ce3c9ac6bf96ae75be527f1b1c301a7dc17c680f0bea90f33e28dd8d76d26dbe@172.17.0.7:30303,enode://c1cad3139b0ab583de214e3d64f7fb7793995023559f7fa1e6b01e87603145ca8e60d5d9f8e23d08df3d1c0c82294bd9515b729efec210f060b2fe3a193f9ae0@172.17.0.6:30303
