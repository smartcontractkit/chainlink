package main

import (
	"github.com/K-Phoen/grabana/dashboard"
	lib "github.com/smartcontractkit/chainlink/dashboard-lib"
	ccipLoadTestView "github.com/smartcontractkit/chainlink/dashboard-lib/ccip-load-test-view"
	coreDon "github.com/smartcontractkit/chainlink/dashboard-lib/core-don"
	coreOCRv2ccip "github.com/smartcontractkit/chainlink/dashboard-lib/core-ocrv2-ccip"
	k8spods "github.com/smartcontractkit/chainlink/dashboard-lib/k8s-pods"
	waspdb "github.com/smartcontractkit/wasp/dashboard"
)

func main() {
	cfg := lib.ReadEnvDeployOpts()
	db := lib.NewDashboard(cfg.DashboardName, cfg,
		[]dashboard.Option{
			dashboard.AutoRefresh("10s"),
			dashboard.Tags([]string{"generated"}),
		},
	)
	db.Add(
		ccipLoadTestView.New(
			ccipLoadTestView.Props{
				LokiDataSource: cfg.DataSources.Loki,
			},
		),
	)
	db.Add(
		coreOCRv2ccip.New(
			coreOCRv2ccip.Props{
				PrometheusDataSource: cfg.DataSources.Prometheus,
				PluginName:           "CCIPCommit",
			},
		),
	)
	db.Add(
		coreOCRv2ccip.New(
			coreOCRv2ccip.Props{
				PrometheusDataSource: cfg.DataSources.Prometheus,
				PluginName:           "CCIPExecution",
			},
		),
	)
	db.Add(
		coreDon.New(
			coreDon.Props{
				PrometheusDataSource: cfg.DataSources.Prometheus,
				PlatformOpts:         coreDon.PlatformPanelOpts(cfg.Platform),
			},
		),
	)
	if cfg.Platform == "kubernetes" {
		db.Add(
			k8spods.New(
				k8spods.Props{
					PrometheusDataSource: cfg.DataSources.Prometheus,
					LokiDataSource:       cfg.DataSources.Loki,
				},
			),
		)
	}
	db.Add(waspdb.AddVariables(cfg.DataSources.Loki))
	if err := db.Deploy(); err != nil {
		lib.L.Fatal().Err(err).Msg("failed to deploy the dashboard")
	}
	lib.L.Info().
		Str("Name", cfg.DashboardName).
		Str("GrafanaURL", db.DeployOpts.GrafanaURL).
		Str("GrafanaFolder", db.DeployOpts.GrafanaFolder).
		Msg("Dashboard deployed")
}
