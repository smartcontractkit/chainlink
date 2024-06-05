package cmd_test

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/smartcontractkit/chainlink-common/observability-lib/cmd"

	"github.com/stretchr/testify/require"
)

func TestNewDashboard(t *testing.T) {
	t.Run("NewDashboard returns a new Dashboard struct", func(t *testing.T) {
		builder, err := dashboard.NewDashboardBuilder("test").Build()
		if err != nil {
			t.Errorf("Error building dashboard: %v", err)
		}

		testDashboard := cmd.NewDashboard(
			"test",
			"",
			"",
			"",
			cmd.DataSources{
				Metrics: "test",
			},
			"kubernetes",
			builder,
		)

		require.IsType(t, &cmd.Dashboard{}, testDashboard)
	})
}

func TestSetDataSources(t *testing.T) {
	t.Run("SetDataSources returns DataSources struct with metrics Prometheus", func(t *testing.T) {
		dataSources := []string{"Prometheus"}
		dataSourcesType := cmd.SetDataSources(dataSources)
		require.Equal(t, "Prometheus", dataSourcesType.Metrics)
		require.Equal(t, "", dataSourcesType.Logs)
	})

	t.Run("SetDataSources returns DataSources struct with metrics Thanos", func(t *testing.T) {
		dataSources := []string{"Thanos"}
		dataSourcesType := cmd.SetDataSources(dataSources)
		require.Equal(t, "Thanos", dataSourcesType.Metrics)
		require.Equal(t, "", dataSourcesType.Logs)
	})

	t.Run("SetDataSources returns DataSources struct with logs Loki", func(t *testing.T) {
		dataSources := []string{"Loki"}
		dataSourcesType := cmd.SetDataSources(dataSources)
		require.Equal(t, "", dataSourcesType.Metrics)
		require.Equal(t, "Loki", dataSourcesType.Logs)
	})
}

func TestGetJSON(t *testing.T) {
	builder, err := dashboard.NewDashboardBuilder("test").Build()
	if err != nil {
		t.Errorf("Error building dashboard: %v", err)
	}

	testDashboard := cmd.NewDashboard(
		"test",
		"",
		"",
		"",
		cmd.DataSources{
			Metrics: "test",
		},
		"kubernetes",
		builder,
	)
	jsonDashboard, errGenerate := testDashboard.GetJSON()
	if errGenerate != nil {
		t.Errorf("Error generating JSON: %v", errGenerate)
	}

	t.Run("GetJSON return dashboard in JSON", func(t *testing.T) {
		expected := `{
          "title": "test",
          "timezone": "browser",
          "graphTooltip": 0,
          "fiscalYearStartMonth": 0,
          "schemaVersion": 0,
          "templating": {},
          "annotations": {}
        }`
		require.JSONEq(t, expected, string(jsonDashboard))
	})
}
