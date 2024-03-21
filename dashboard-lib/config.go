package dashboard_lib

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type EnvConfig struct {
	Platform                 string
	GrafanaURL               string
	GrafanaToken             string
	GrafanaBasicAuthUser     string
	GrafanaBasicAuthPassword string
	GrafanaFolder            string
	DataSources              DataSources
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
	ba := os.Getenv("GRAFANA_BASIC_AUTH")
	if ba == "" {
		L.Fatal().Msg("GRAFANA_BASIC_AUTH is empty")
	}
	user, password, err := decodeBasicAuth(ba)
	if err != nil {
		L.Fatal().Err(err).Msg("failed to decode basic auth")
	}
	return EnvConfig{
		GrafanaURL:               grafanaURL,
		GrafanaToken:             grafanaToken,
		GrafanaBasicAuthUser:     user,
		GrafanaBasicAuthPassword: password,
		GrafanaFolder:            grafanaFolder,
		Platform:                 platform,
		DataSources: DataSources{
			Loki:       loki,
			Prometheus: prom,
		},
	}
}

func decodeBasicAuth(encodedAuth string) (string, string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedAuth)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(string(decodedBytes), ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid basic authentication format")
	}
	return parts[0], parts[1], nil
}
