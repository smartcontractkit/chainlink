# Performance tests for CL jobs

This folder container performance e2e tests for different job types, currently implemented:

- VRFv2
- CCIP (wip)

All the tests have 4 groups:

- one product soak
- one product load
- multiple product instances soak
- multiple product instances load

When you develop an e2e performance suite for a new product you can implement the tests one by one to answer the questions:

What are performance characteristics of a one instance of a product (jobs + contracts):

- is our product stable at all, no memory leaks, no flaking performance under some RPS? (test #1)
- what are the limits for one product instance, figuring out the max/optimal performance params by increasing RPS and varying configuration (test #2)
- update test #1 with optimal params and workload to constantly run in CI

What are performance and capacity characteristics of Chainlink node(s) that run multiple products of the same type simultaneously:

- how many products of the same type we can run at once at a stable load with optimal configuration? (test #3)
- what are the limits if we add more and more products of the same type, each product have a stable RPS, we vary only amount of products
- update test #3 with optimal params and workload to constantly run in CI

Implementing test #1 is **mandatory** for each product.
Tests #2,#3,#4 are optional if you need to figure out your product scaling or node/cluster capacity.

## Usage

```sh
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestVRFV2Load/vrfv2_soak_test
```

### Dashboards

Each product has its own generated dashboard in `cmd/dashboard.go`

Deploying dashboard:

```sh
export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export DATA_SOURCE_NAME=grafanacloud-logs
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=${JobTypeName, for example WaspVRFv2}

go run dashboard.go
```

### Assertions

You can implement your product requirements assertions in `onchain_monitoring.go`. For Chainlink products (VRF/OCR/Keepers/Automation) it's easier to assert on-chain part

If you need to assert some metrics in `Prometheus/Loki`, here is an [example](https://github.com/smartcontractkit/wasp/blob/master/examples/alerts/main_test.go#L88)

Do not mix workload logic with assertions, separate them.

### Implementation

To implement a standard e2e performance suite for a new product please look at `gun.go` and `vu.go`.

Gun should be working with one instance of your product.

VU(Virtual user) creates a new instance of your product and works with it in `Call()`

### Cluster mode (k8s)

Add configuration to `overrides.toml`

```toml
[WaspAutoBuild]
namespace = "wasp"
update_image = true
repo_image_version_uri = "${staging_ecr_registry}/wasp-tests:wb-core"
test_binary_name = "ocr.test"
test_name = "TestOCRLoad"
test_timeout = "24h"
wasp_log_level = "debug"
wasp_jobs = "1"
keep_jobs = true
```

And run your tests using `go test -v -run TestClusterEntrypoint`
