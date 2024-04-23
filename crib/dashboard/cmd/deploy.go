package main

import (
	"github.com/K-Phoen/grabana/dashboard"
	lib "github.com/smartcontractkit/chainlink/dashboard-lib"
	atlas_don "github.com/smartcontractkit/chainlink/dashboard-lib/atlas-don"
	core_don "github.com/smartcontractkit/chainlink/dashboard-lib/core-don"
	core_node_components "github.com/smartcontractkit/chainlink/dashboard-lib/core-node-components"
	k8spods "github.com/smartcontractkit/chainlink/dashboard-lib/k8s-pods"
	waspdb "github.com/smartcontractkit/wasp/dashboard"
	"strings"
)

func main() {
	cfg := lib.ReadEnvDeployOpts()
	db := lib.NewDashboard(cfg.Name, cfg,
		[]dashboard.Option{
			dashboard.AutoRefresh("10s"),
			dashboard.Tags([]string{"generated"}),
		},
	)
	if len(cfg.PanelsIncluded) == 0 || cfg.PanelsIncluded["core"] {
		db.Add(
			core_don.New(
				core_don.Props{
					PrometheusDataSource: cfg.DataSources.Prometheus,
					PlatformOpts:         core_don.PlatformPanelOpts(cfg.Platform),
				},
			),
		)
		// TODO: refactor as a component later
		addWASPRows(db, cfg)
	}
	if cfg.PanelsIncluded["core_components"] {
		db.Add(
			core_node_components.New(
				core_node_components.Props{
					PrometheusDataSource: cfg.DataSources.Prometheus,
					PlatformOpts:         core_node_components.PlatformPanelOpts(),
				},
			),
		)
	}
	if cfg.PanelsIncluded["ocr"] || cfg.PanelsIncluded["ocr2"] || cfg.PanelsIncluded["ocr3"] {
		for key := range cfg.PanelsIncluded {
			if strings.Contains(key, "ocr") {
				db.Add(
					atlas_don.New(
						atlas_don.Props{
							PrometheusDataSource: cfg.DataSources.Prometheus,
							PlatformOpts:         atlas_don.PlatformPanelOpts(cfg.Platform, key),
							OcrVersion:           key,
						},
					),
				)
			}
		}
	}
	if !cfg.PanelsIncluded["core_components"] && cfg.Platform == "kubernetes" {
		db.Add(
			k8spods.New(
				k8spods.Props{
					PrometheusDataSource: cfg.DataSources.Prometheus,
					LokiDataSource:       cfg.DataSources.Loki,
				},
			),
		)
	}
	if err := db.Deploy(); err != nil {
		lib.L.Fatal().Err(err).Msg("failed to deploy the dashboard")
	}
	lib.L.Info().
		Str("Name", db.Name).
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
