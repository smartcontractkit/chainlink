# Chainlink Grafana Dashboards Library

This library offers dashboard components and tools for constructing and testing Grafana dashboards at Chainlink.

Components structure is as follows:
```
dashboard
  |- lib
     |- component_1
        |- component.go
     |- component_2
        |- component.go
|- go.mod
|- go.sum
```

Each component should contain rows, logic and unique variables in `component.go`

Components should be imported from this module, see [example](../crib/dashboard/cmd/deploy.go)

## How to convert from JSON using Grabana codegen utility
1. Download Grabana binary [here](https://github.com/K-Phoen/grabana/releases)
2. ./bin/grabana convert-go -i dashboard.json > lib/my_new_component/rows.go
3. Create a [component](k8s-pods/component.go)