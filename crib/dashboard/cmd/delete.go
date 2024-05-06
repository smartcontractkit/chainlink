package main

import (
	lib "github.com/smartcontractkit/chainlink/dashboard-lib"
)

func main() {
	cfg := lib.ReadEnvDeployOpts()
	db := lib.NewDashboard(cfg.Name, cfg, nil)
	err := db.Delete()
	if err != nil {
		lib.L.Fatal().Err(err).Msg("failed to delete the dashboard")
	}
	lib.L.Info().
		Str("Name", db.Name).
		Str("GrafanaURL", db.DeployOpts.GrafanaURL).
		Str("GrafanaFolder", db.DeployOpts.GrafanaFolder).
		Msg("Dashboard deleted")
}
