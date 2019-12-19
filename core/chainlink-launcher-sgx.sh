#!/bin/bash -e

ETHEREUM_CRT="${ETHEREUM_CRT:-/secrets/ethereum/tls.crt}"
if [ -f $ETHEREUM_CRT ]; then
  cp $ETHEREUM_CRT /usr/local/share/ca-certificates/eth.crt
  update-ca-certificates
fi

command=`echo $1 | tr A-Z a-z`
if [ "$command" != "n" ] && [ "$command" != "node" ]; then
  chainlink "$@"
  exit
fi

if [ "$SGX_SIMULATION" != "yes" ]; then
  LD_LIBRARY_PATH=/opt/intel/libsgx-enclave-common/aesm /opt/intel/sgxpsw/aesm/aesm_service &
  aesm_pid=$!
  trap "kill $aesm_pid 2>/dev/null || true" SIGINT SIGTERM EXIT
fi

# XXX: Since chainlink has to run in the background in SGX mode, prevent it
# from detecting a TTY so that it does not prompt for a password
chainlink "$@" | cat &
chainlink_pid=$!

if [ "$SGX_SIMULATION" != "yes" ]; then
  while sleep 10; do
    kill -0 $aesm_pid 2>/dev/null
    kill -0 $chainlink_pid 2>/dev/null
  done
fi
wait $chainlink_pid
exit $?
