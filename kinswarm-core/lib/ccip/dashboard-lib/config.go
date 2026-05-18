package dashboard_lib

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type EnvConfig struct {
	Name                     string
	Platform                 string
	GrafanaURL               string
	GrafanaToken             string
	GrafanaBasicAuthUser     string
	GrafanaBasicAuthPassword string
	GrafanaFolder            string
	DataSources              DataSources
	PanelsIncluded           map[string]bool
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
	panelsIncludedString := os.Getenv("PANELS_INCLUDED")
	panelsIncludedArray := strings.Split(panelsIncludedString, ",")
	panelsIncluded := make(map[string]bool)

	if panelsIncludedString != "" {
		for _, panelName := range panelsIncludedArray {
			panelsIncluded[panelName] = true
		}
	}

	ba := os.Getenv("GRAFANA_BASIC_AUTH")
	grafanaToken := os.Getenv("GRAFANA_TOKEN")
	if grafanaToken == "" && ba == "" {
		L.Fatal().Msg("GRAFANA_TOKEN or GRAFANA_BASIC_AUTH must be provided")
	}
	var user, password string
	var err error
	if ba != "" {
		user, password, err = DecodeBasicAuth(ba)
		if err != nil {
			L.Fatal().Err(err).Msg("failed to decode basic auth")
		}
	}

	return EnvConfig{
		Name:                     name,
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
		PanelsIncluded: panelsIncluded,
	}
}

func DecodeBasicAuth(authString string) (string, string, error) {
	var data string
	decodedBytes, err := base64.StdEncoding.DecodeString(authString)
	if err != nil {
		L.Warn().Err(err).Msg("failed to decode basic auth, plain text? reading auth data")
		data = authString
	} else {
		data = string(decodedBytes[1 : len(decodedBytes)-1])
	}
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return "", "", errors.New("invalid basic authentication format")
	}
	return parts[0], parts[1], nil
}
