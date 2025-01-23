# CCIP Load Tests 1.6+

Resources
- [CRIB setup](https://smartcontract-it.atlassian.net/wiki/spaces/CRIB/pages/1024622593/CCIP+v2+CRIB+-+Deploy+Access+Instructions+WIP)
- [Load Test](https://smartcontract-it.atlassian.net/wiki/spaces/CCIP/pages/edit-v2/1247707289?draftShareId=5d462e12-a0bb-4e05-99d4-f9687d40b627)
- [Demo](https://drive.google.com/file/d/1U3WAiyuoCVYGWMRWnG8ESRR3ie7lp7er/view?usp=drive_web)


## Usage

### Local
Follow the instructions in CRIB Setup to get CRIB running. If you want to execute against a specific Chainlink image,
you can set the `DEVSPACE_IMAGE_TAG` environment variable in the .env you create in the crib repo. This image tag can be found in ECR. 

Make sure you're using `kind` as the provider to run these changes locally. 
```sh
DEVSPACE_IMAGE= <found in ECR>
DEVSPACE_IMAGE_TAG= <found in ECR>

DEVSPACE_NAMESPACE=crib-local
PROVIDER=kind
```

Modify the `ccip.toml` or `overrides.toml` file to your desired configuration.
```sh
# replace this with the loki endpoint of the crib stack or from running `devspace ingress-hosts`
LokiEndpoint=https://find.this.in.crib.output
# MessageTypeWeights corresponds with [data only, token only, message with data]
MessageTypeWeights=[100,0,0]
# all enabled destination chains will receive 1 incoming request per RequestFrequency for the duration of LoadDuration
RequestFrequency="10s"
LoadDuration="1m"
# destination chain selectors to send messages to ie [3379446385462418246,909606746561742123, etc.]
EnabledDestionationChains=[3379446385462418246]
# Directory where we receive environment configuration from crib
CribEnvDirectory="directory/to/crib/output"

```

Execute the test using 
```sh
export TIMEOUT=6h
go test -run ^TestCCIPLoad_RPS$ ./integration-tests/load/ccip -v -timeout $TIMEOUT`
```

### Remote
Update the `PROVIDER=aws` and `DEVSPACE_NAMESPACE` in crib environment and deploy. Everything else should be the same. 