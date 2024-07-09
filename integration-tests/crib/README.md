### CRIB Health Check Test

## Setup CRIB
This is a simple smoke + chaos test for CRIB deployment.
It runs OCRv1 and reboots the environment confirming integration with environment is working and data is properly saved even after reboots.

```shell
CRIB_NAMESPACE=any-crib-namespace
CRIB_NETWORK=geth
CRIB_NODES=5
go test -v -run TestCRIB
```