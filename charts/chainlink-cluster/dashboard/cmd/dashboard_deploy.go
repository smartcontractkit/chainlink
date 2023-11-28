package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/dashboard/dashboard"
	"github.com/smartcontractkit/wasp"
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
	waspDsn := os.Getenv("DATA_SOURCE_NAME")
	if waspDsn == "" {
		panic("DATA_SOURCE_NAME must be provided, should be the same as LOKI_DATA_SOURCE_NAME")
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
	// if you'll use this dashboard base in other projects, you can add your own opts here to extend it
	db, err := dashboard.NewCLClusterDashboard(name, ldsn, pdsn, dbf, grafanaURL, grafanaToken, nil)
	if err != nil {
		panic(err)
	}
	// here we are extending load testing dashboard with core metrics, for example
	wdb, err := wasp.NewDashboard(nil, db.Opts())
	if err != nil {
		panic(err)
	}
	if _, err := wdb.Deploy(); err != nil {
		panic(err)
	}
}
