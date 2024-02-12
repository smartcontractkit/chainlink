package main

import (
	"fmt"
	"github.com/smartcontractkit/chainlink/charts/chainlink-cluster/dashboard/dashboard"
	"os"
)

func main() {
	name := os.Getenv("DASHBOARD_NAME")
	if name == "" {
		panic("DASHBOARD_NAME must be provided")
	}
	// Can be empty
	lokiDataSourceName := os.Getenv("LOKI_DATA_SOURCE_NAME")

	prometheusDataSourceName := os.Getenv("PROMETHEUS_DATA_SOURCE_NAME")
	if prometheusDataSourceName == "" {
		panic("PROMETHEUS_DATA_SOURCE_NAME must be provided")
	}

	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		panic("GRAFANA_URL must be provided")
	}

	grafanaToken := os.Getenv("GRAFANA_TOKEN")
	if grafanaToken == "" {
		panic("GRAFANA_TOKEN must be provided")
	}

	grafanaFolder := os.Getenv("GRAFANA_FOLDER")
	if grafanaFolder == "" {
		panic("GRAFANA_FOLDER must be provided")
	}

	infraPlatform := os.Getenv("INFRA_PLATFORM")
	if infraPlatform == "" {
		panic("INFRA_PLATFORM must be provided, can be either docker|kubernetes")
	}

	err := dashboard.NewDashboard(
		name,
		grafanaURL,
		grafanaToken,
		grafanaFolder,
		[]string{"generated"},
		lokiDataSourceName,
		prometheusDataSourceName,
		infraPlatform,
		[]string{"core"},
		nil,
	)
	if err != nil {
		fmt.Printf("Could not create dashbard: %s\n", name)
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully deployed %s dashboard on grafana instance %s in folder %s\n",
		name,
		grafanaURL,
		grafanaFolder,
	)
}
