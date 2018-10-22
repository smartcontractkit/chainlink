#!/bin/bash -e

ETHEREUM_CRT="${ETHEREUM_CRT:-/secrets/ethereum/tls.crt}"
if [ -f $ETHEREUM_CRT ]; then
  cp $ETHEREUM_CRT /usr/local/share/ca-certificates/eth.crt
  update-ca-certificates
fi

chainlink $@
