#!/bin/sh

# Check if networkid is provided as an argument
if [ -z "$1" ]; then
  echo "Network ID is not provided. Usage: ./entrypoint.sh <networkid>"
  exit 1
fi

NETWORKID=$1

# Initialize geth
geth --datadir "./data-temp" init ./data-temp/genesis.json

# Run geth
geth --dev \
  --datadir "./data-temp" \
  --networkid "$NETWORKID" \
  --ipcdisable \
  --unlock 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 \
  --mine \
  --miner.etherbase 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 \
  --dev.period 1 \
  --http \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.api "admin,eth,net,web3,personal,txpool,miner,debug" \
  --http.corsdomain "*" \
  --http.vhosts "*" \
  --password /data-temp/password.txt \
  --ws \
  --ws.addr 0.0.0.0 \
  --ws.port 8546 \
  --ws.api "admin,eth,net,web3,personal,txpool,miner,debug" \
  --ws.origins "*" \
  --allow-insecure-unlock \
  --authrpc.addr 0.0.0.0 \
  --authrpc.port 8551 \
  --authrpc.vhosts "*" \
  --rpc.allow-unprotected-txs \
  --rpc.txfeecap 0 \
  --nodiscover \
  --verbosity 3 \
  --log.vmodule "rpc=5"