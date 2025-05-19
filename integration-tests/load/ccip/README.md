# CCIP Load and Chaos Tests 1.6+

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

## Remote
Update the `PROVIDER=aws` and `DEVSPACE_NAMESPACE` in crib environment and deploy. Everything else should be the same. 

## Running Chaos Tests (Remote only)

### Realistic RPC Latency

Go to `integration-tests/testconfig/ccip/ccip.toml` and change params as required, select chaos mode first
```
[Load.CCIP.Load]
# 0 - no chaos, 1 - rpc latency, 2 - full chaos suite
ChaosMode = 1

# works only with Load.CCIP.ChaosMode = 1
RPCLatency = "400ms"
RPCJitter = "50ms"
```
Prefer using `ChaosMode = 1` with `400ms` RPC latency by default

Run the load test
```
go test -run ^TestCCIPLoad_RPS$ -v -timeout 12h
```

### Full Chaos Suite

Go to `integration-tests/testconfig/ccip/ccip.toml` and change params as required, select chaos mode first

```
[Load.CCIP.Load]
# 0 - no chaos, 1 - rpc latency, 2 - full chaos suite
ChaosMode = 2
```

Then check chaos settings, change `Namespace`, `ExperimentFullInterval`, `ExperimentInjectionInterval` accordingly, add any dashboard UIDs you need
```
[CCIP.Chaos]
Namespace = "crib-ccip-chaos-tests"
# RPC, commit, exec dashboards, can be found here: https://grafana.ops.prod.cldev.sh/d/dde396ff-5d22-42fb-9e92-00845c17688c/load-ccipv2-exec-plugin-v2?orgId=1&editview=dashboard_json
DashboardUIDs = ["e08d9f98-a39a-4603-8b44-e9a2958330e4", "ed3d5742-57cb-440f-b432-65f229c124ec", "dde396ff-5d22-42fb-9e92-00845c17688c"]
WaitBeforeStart = "30s"
# works only with Load.CCIP.ChaosMode = 2
# Chaos experiment total duration (chaos + recovery)
ExperimentFullInterval = "10m"
# Chaos time
ExperimentInjectionInterval = "9m"
```
By default, we keep all chaos CRDs after each run so it's easier to debug, but when you suite is stable you can set this flag to `true` to remove CRDs [automatically](https://github.com/smartcontractkit/chainlink/blob/aw/LoadTest-aws-chaos/integration-tests/load/ccip/ccip_chaos_test.go#L46)

Run the load test
```
go test -run ^TestCCIPLoad_RPS$ -v -timeout 12h
```