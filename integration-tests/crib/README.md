### CRIB Health Check Test

## Setup CRIB
This is a simple smoke + chaos test for CRIB deployment.
It runs OCRv1 and reboots the environment confirming integration with environment is working and data is properly saved even after reboots.
Go to the [CRIB](https://github.com/smartcontractkit/crib) repository and spin up a cluster.

```shell
./scripts/cribbit.sh crib-oh-my-crib
devspace deploy --debug --profile local-dev-simulated-core-ocr1
```

## Run the tests
```shell
CRIB_NAMESPACE=crib-oh-my-crib
CRIB_NETWORK=geth # only "geth" is supported for now
CRIB_NODES=5 # min 5 nodes
go test -v -run TestCRIB
```