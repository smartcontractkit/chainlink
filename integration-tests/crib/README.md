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
export CRIB_NAMESPACE=crib-oh-my-crib
export CRIB_NETWORK=geth # only "geth" is supported for now
export CRIB_NODES=5 # min 5 nodes
#export SETH_LOG_LEVEL=debug # these two can be enabled to debug connection issues
#export RESTY_DEBUG=true
export GAP_URL=https://localhost:8080/primary # only applicable in CI, unset the var to connect locally
go test -v -run TestCRIB
```