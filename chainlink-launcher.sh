#!/bin/bash -e

command=`echo $1 | tr A-Z a-z`
if [ "$command" != "n" ] && [ "$command" != "node" ]; then
  ./chainlink "$@"
  exit 0
fi

trap "kill -- -$$ 2>/dev/null || true" SIGINT SIGTERM EXIT

if [ "$SGX_SIMULATION" != yes ]; then
  /opt/intel/sgxpsw/aesm/aesm_service &
  aesm_pid=$!
fi

./chainlink "$@" | cat &
chainlink_pid=$!

if [ "$SGX_SIMULATION" != yes ]; then
  while sleep 10; do
    kill -0 $aesm_pid 2>/dev/null
    kill -0 $chainlink_pid 2>/dev/null
  done
fi
wait $chainlink_pid
exit $?
