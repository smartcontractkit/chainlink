### VRFv2 Load tests

## Usage
```
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestVRFV2Load/vrfv2_soak_test
```

### Dashboards

Deploying dashboard:
```
export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export DATA_SOURCE_NAME=grafanacloud-logs
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=${JobTypeName, for example WaspVRFv2}

go run dashboard.go
```