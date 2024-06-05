package atlasdon_test

import (
	"testing"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"

	"github.com/stretchr/testify/require"
)

func TestBuildDashboard(t *testing.T) {
	t.Run("BuildDashboard creates a dashboard", func(t *testing.T) {
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
