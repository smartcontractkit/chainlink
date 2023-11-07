package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/dashboard/dashboard"
)

func main() {
	name := os.Getenv("DASHBOARD_NAME")
	if name == "" {
		panic("DASHBOARD_NAME must be provided")
	}
	ldsn := os.Getenv("LOKI_DATA_SOURCE_NAME")
	if ldsn == "" {
		panic("DATA_SOURCE_NAME must be provided")
	}
	pdsn := os.Getenv("PROMETHEUS_DATA_SOURCE_NAME")
	if ldsn == "" {
		panic("DATA_SOURCE_NAME must be provided")
	}
	dbf := os.Getenv("DASHBOARD_FOLDER")
	if dbf == "" {
		panic("DASHBOARD_FOLDER must be provided")
	}
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		panic("GRAFANA_URL must be provided")
	}
	grafanaToken := os.Getenv("GRAFANA_TOKEN")
	if grafanaToken == "" {
		panic("GRAFANA_TOKEN must be provided")
	}
	db, err := dashboard.NewCLClusterDashboard(name, ldsn, pdsn, dbf, grafanaURL, grafanaToken, nil)
	if err != nil {
		panic(err)
	}
	if err := db.Deploy(); err != nil {
		panic(err)
	}
}
