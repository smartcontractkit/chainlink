package utils_test

import (
	"fmt"
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"

	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"

	"github.com/stretchr/testify/require"
)

func TestQueryVariable(t *testing.T) {
	t.Run("QueryVariable creates a query variable", func(t *testing.T) {
		queryVariableTest := utils.QueryVariable("Prometheus", "test", "Test", fmt.Sprintf("label_values(%s)", "test"), false)

		require.IsType(t, &dashboard.QueryVariableBuilder{}, queryVariableTest)
	})
}

func TestIntervalVariable(t *testing.T) {
	t.Run("IntervalVariable creates an interval variable", func(t *testing.T) {
		intervalVariableTest := utils.IntervalVariable("test", "Test", "1m")

		require.IsType(t, &dashboard.IntervalVariableBuilder{}, intervalVariableTest)
	})
}

func TestAddVars(t *testing.T) {
	t.Run("AddVars adds variables to the dashboard", func(t *testing.T) {
		builder := dashboard.NewDashboardBuilder("test")

		utils.AddVars(builder, []cog.Builder[dashboard.VariableModel]{
			utils.QueryVariable("Prometheus", "test", "Test", fmt.Sprintf("label_values(%s)", "test"), false),
		})

		testBuild, err := builder.Build()
		if err != nil {
			t.Errorf("Error building dashboard: %v", err)
		}

		require.IsType(t, dashboard.Dashboard{}, testBuild)
	})
}

func TestStatPanel(t *testing.T) {
	t.Run("StatPanel creates a stat panel", func(t *testing.T) {
		statPanelTest := utils.StatPanel(
			"Prometheus",
			"Test",
			"Test",
			1,
			1,
			1,
			"",
			common.BigValueColorModeNone,
			common.BigValueGraphModeNone,
			common.BigValueTextModeName,
			common.VizOrientationHorizontal,
			utils.PrometheusQuery{
				Query:  `test`,
				Legend: "{{test}}",
			})

		require.IsType(t, &stat.PanelBuilder{}, statPanelTest)
	})
}

func TestAddPanels(t *testing.T) {
	t.Run("AddPanels adds panels to the dashboard", func(t *testing.T) {
		builder := dashboard.NewDashboardBuilder("test")
		utils.AddPanels(builder, []cog.Builder[dashboard.Panel]{
			utils.StatPanel(
				"Prometheus",
				"Test",
				"Test",
				1,
				1,
				1,
				"",
				common.BigValueColorModeNone,
				common.BigValueGraphModeNone,
				common.BigValueTextModeName,
				common.VizOrientationHorizontal,
				utils.PrometheusQuery{
					Query:  `test`,
					Legend: "{{test}}",
				}),
		})

		testBuild, err := builder.Build()
		if err != nil {
			t.Errorf("Error building dashboard: %v", err)
		}

		require.IsType(t, dashboard.Dashboard{}, testBuild)
	})
}
