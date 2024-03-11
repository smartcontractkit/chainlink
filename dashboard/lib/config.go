package dashboardlib

import "os"

type EnvConfig struct {
	Platform      string
	GrafanaURL    string
	GrafanaToken  string
	GrafanaFolder string
	DataSources   DataSources
}

type DataSources struct {
	Loki       string
	Prometheus string
}

type DashboardOpts struct {
	Tags        []string
	AutoRefresh string
}

func ReadEnvDeployOpts() EnvConfig {
	name := os.Getenv("DASHBOARD_NAME")
	if name == "" {
		L.Fatal().Msg("DASHBOARD_NAME must be provided")
	}
	lokiDataSourceName := os.Getenv("LOKI_DATA_SOURCE_NAME")
	if lokiDataSourceName == "" {
		L.Fatal().Msg("LOKI_DATA_SOURCE_NAME is empty, panels with logs will be disabled")
	}
	prometheusDataSourceName := os.Getenv("PROMETHEUS_DATA_SOURCE_NAME")
	if prometheusDataSourceName == "" {
		L.Fatal().Msg("PROMETHEUS_DATA_SOURCE_NAME must be provided")
	}
	grafanaURL := os.Getenv("GRAFANA_URL")
	if grafanaURL == "" {
		L.Fatal().Msg("GRAFANA_URL must be provided")
	}
	grafanaToken := os.Getenv("GRAFANA_TOKEN")
	if grafanaToken == "" {
		L.Fatal().Msg("GRAFANA_TOKEN must be provided")
	}
	grafanaFolder := os.Getenv("GRAFANA_FOLDER")
	if grafanaFolder == "" {
		L.Fatal().Msg("GRAFANA_FOLDER must be provided")
	}
	platform := os.Getenv("INFRA_PLATFORM")
	if platform == "" {
		L.Fatal().Msg("INFRA_PLATFORM must be provided, can be either docker|kubernetes")
	}
	loki := os.Getenv("LOKI_DATA_SOURCE_NAME")
	if lokiDataSourceName == "" {
		L.Fatal().Msg("LOKI_DATA_SOURCE_NAME is empty, panels with logs will be disabled")
	}
	prom := os.Getenv("PROMETHEUS_DATA_SOURCE_NAME")
	if prometheusDataSourceName == "" {
		L.Fatal().Msg("PROMETHEUS_DATA_SOURCE_NAME must be provided")
	}
	return EnvConfig{
		GrafanaURL:    grafanaURL,
		GrafanaToken:  grafanaToken,
		GrafanaFolder: grafanaFolder,
		Platform:      platform,
		DataSources: DataSources{
			Loki:       loki,
			Prometheus: prom,
		},
	}
}
