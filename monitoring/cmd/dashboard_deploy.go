package main

import (
	"context"
	"fmt"
	"github.com/K-Phoen/grabana"
	"github.com/smartcontractkit/chainlink/v2/core/dashboard"
	"net/http"
	"os"
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
	os.Setenv("DATA_SOURCE_NAME", ldsn)
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
	// if you'll use this dashboard base in other projects, you can add your own opts here to extend it
	db, err := dashboard.NewDashboard(name, ldsn, pdsn, "kubernetes", nil)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, grafanaURL, grabana.WithAPIToken(grafanaToken))
	fo, err := client.FindOrCreateFolder(ctx, dbf)
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}
	if _, err := client.UpsertDashboard(ctx, fo, db.Builder); err != nil {
		panic(err)
	}
}
