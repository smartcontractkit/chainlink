#!/bin/bash -e

ETHEREUM_CRT="${ETHEREUM_CRT:-/secrets/ethereum/tls.crt}"
if [ -f $ETHEREUM_CRT ]; then
  cp $ETHEREUM_CRT /usr/local/share/ca-certificates/eth.crt
  update-ca-certificates
fi

command=`echo $1 | tr A-Z a-z`
if [ "$command" != "n" ] && [ "$command" != "node" ]; then
  chainlink "$@"
  exit 0
fi

trap "kill -- -$$ 2>/dev/null || true" SIGINT SIGTERM EXIT

if [ "$SGXENABLED" = "yes" ] && [ "$SGX_SIMULATION" != "yes" ]; then
  /opt/intel/sgxpsw/aesm/aesm_service &
  aesm_pid=$!
fi

chainlink "$@" | cat &
chainlink_pid=$!

if [ "$SGXENABLED" = "yes" ] && [ "$SGX_SIMULATION" != "yes" ]; then
  while sleep 10; do
    kill -0 $aesm_pid 2>/dev/null
    kill -0 $chainlink_pid 2>/dev/null
  done
fi
wait $chainlink_pid
exit $?
