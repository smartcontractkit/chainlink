FROM ethereum/client-go:v1.10.8
# docker build . -t smartcontract/gethnet:london
WORKDIR /gethnet
COPY node_password .
COPY genesis.json .
# Initializes genesis file with london forks enabled
RUN geth --datadir /gethnet/datadir init genesis.json
# Copy a prefunded devnet key into the keystore
COPY keys/* /gethnet/datadir/keystore/
EXPOSE 8545 8546 8547 30303 30303/udp
ENTRYPOINT [ "geth" ]
CMD [ \
"--networkid=34055", \
"--mine", \
"--miner.threads=1", \
"--miner.noverify", \
"--miner.recommit=1s", \
"--datadir=/gethnet/datadir", \
"--fakepow", \
"--nodiscover", \
"--http", \
"--http.addr=0.0.0.0", \
"--http.port=8545", \
"--port=30303", \
"--http.corsdomain", "*", \
"--http.api", "eth,web3,personal,net", \
"--password=node_password", \
"--ipcdisable", \
"--unlock", "0", \
"--allow-insecure-unlock", \
"--ws", \
"--ws.addr=0.0.0.0", \
"--ws.port=8546", \
"--ws.api","eth,web3,net,admin,debug,txpool", \
"--txpool.accountslots=1024", \
"--txpool.accountqueue=1024" \
]