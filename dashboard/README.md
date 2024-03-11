# Chainlink Grafana Dashboards Library

This library offers dashboard components and tools for constructing and testing Grafana dashboards at Chainlink.

Components structure is as follows:
```
dashboard
  |- lib
     |- component_1
        |- component.go
        |- component.spec.ts
     |- component_2
        |- component.go
        |- component.spec.ts
|- go.mod
|- go.sum
```

Each component should contain rows, logic and unique variables in `component.go`

Components should be imported from this module, see [example](../charts/chainlink-cluster/dashboard/cmd/deploy.go)

`component.spec.ts` is a Playwright test step that can be used when testing project [dashboards](../charts/chainlink-cluster/dashboard/tests/specs/core-don.spec.ts)

## How to convert from JSON using Grabana codegen utility
1. Download Grabana binary [here](https://github.com/K-Phoen/grabana/releases)
2. ./bin/grabana convert-go -i dashboard.json > lib/my_new_component/rows.go
3. Create a [component](lib/k8s-pods/component.go)