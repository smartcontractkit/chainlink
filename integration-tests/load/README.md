## Performance tests for CL jobs

This folder container performance e2e tests for different job types, currently implemented:
- VRFv2

All the tests have 4 groups:
- one product soak
- one product load
- multiple product instances soak
- multiple product instances load

## Usage
```
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestVRFV2Load/vrfv2_soak_test
```

### Dashboards
Each job type has its own generated dashboard in `cmd/dashboard.go`

Deploying dashboard:
```
export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export DATA_SOURCE_NAME=grafanacloud-logs
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=${JobTypeName, for example WaspVRFv2}

go run dashboard.go
```