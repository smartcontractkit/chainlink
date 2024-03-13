package main

import (
	"github.com/K-Phoen/grabana/dashboard"
	lib "github.com/smartcontractkit/chainlink/dashboard-lib/lib"
	core_don "github.com/smartcontractkit/chainlink/dashboard-lib/lib/core-don"
	k8spods "github.com/smartcontractkit/chainlink/dashboard-lib/lib/k8s-pods"
	waspdb "github.com/smartcontractkit/wasp/dashboard"
)

const (
	DashboardName = "Chainlink Cluster (DON)"
)

func main() {
	cfg := lib.ReadEnvDeployOpts()
	db := lib.NewDashboard(DashboardName, cfg,
		[]dashboard.Option{
			dashboard.AutoRefresh("10s"),
			dashboard.Tags([]string{"experimental", "generated"}),
		},
	)
	db.Add(
		core_don.New(
			core_don.Props{
				PrometheusDataSource: cfg.DataSources.Prometheus,
				PlatformOpts:         core_don.PlatformPanelOpts(cfg.Platform),
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
	// TODO: refactor as a component later
	addWASPRows(db, cfg)
	if err := db.Deploy(); err != nil {
		lib.L.Fatal().Err(err).Msg("failed to deploy the dashboard")
	}
	lib.L.Info().
		Str("Name", DashboardName).
		Str("GrafanaURL", db.DeployOpts.GrafanaURL).
		Str("GrafanaFolder", db.DeployOpts.GrafanaFolder).
		Msg("Dashboard deployed")
}

func addWASPRows(db *lib.Dashboard, cfg lib.EnvConfig) {
	if cfg.Platform == "docker" {
		return
	}
	selectors := map[string]string{
		"branch": `=~"${branch:pipe}"`,
		"commit": `=~"${commit:pipe}"`,
	}
	db.Add(waspdb.AddVariables(cfg.DataSources.Loki))
	db.Add(
		[]dashboard.Option{
			waspdb.WASPLoadStatsRow(
				cfg.DataSources.Loki,
				selectors,
			),
		},
	)
	db.Add(
		[]dashboard.Option{
			waspdb.WASPDebugDataRow(
				cfg.DataSources.Loki,
				selectors,
				true,
			),
		},
	)
}
