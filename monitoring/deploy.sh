export GRAFANA_URL=http://localhost:3000
export GRAFANA_TOKEN=secret
export LOKI_DATA_SOURCE_NAME=Prometheus
export PROMETHEUS_DATA_SOURCE_NAME=Prometheus
export DASHBOARD_FOLDER=Node
export DASHBOARD_NAME=ChainlinkCluster

go run cmd/dashboard_deploy.go