package cmd

import (
	"errors"
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/spf13/cobra"

	atlasdon "github.com/smartcontractkit/chainlink-common/observability-lib/atlas-don"
	corenode "github.com/smartcontractkit/chainlink-common/observability-lib/core-node"
	corenodecomponents "github.com/smartcontractkit/chainlink-common/observability-lib/core-node-components"
	k8sresources "github.com/smartcontractkit/chainlink-common/observability-lib/k8s-resources"
	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Grafana Dashboard JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		dataSources, errDataSources := cmd.Flags().GetStringArray("grafana-data-sources")
		if errDataSources != nil || len(dataSources) < 1 {
			return errDataSources
		}

		dataSourcesType := SetDataSources(dataSources)
		name := cmd.Flag("dashboard-name").Value.String()
		platform := cmd.Flag("platform").Value.String()
		typeDashboard := cmd.Flag("type").Value.String()

		var builder dashboard.Dashboard
		var err error

		switch typeDashboard {
		case "core-node":
			builder, err = corenode.BuildDashboard(name, dataSourcesType.Metrics, platform)
		case "core-node-components":
			builder, err = corenodecomponents.BuildDashboard(name, dataSourcesType.Metrics)
		case "core-node-resources":
			builder, err = k8sresources.BuildDashboard(name, dataSourcesType.Metrics, dataSourcesType.Logs)
		case "ocr":
			builder, err = atlasdon.BuildDashboard(name, dataSourcesType.Metrics, platform, typeDashboard)
		case "ocr2":
			builder, err = atlasdon.BuildDashboard(name, dataSourcesType.Metrics, platform, typeDashboard)
		case "ocr3":
			builder, err = atlasdon.BuildDashboard(name, dataSourcesType.Metrics, platform, typeDashboard)
		default:
			return errors.New("invalid dashboard type")
		}
		if err != nil {
			utils.Logger.Error().
				Str("Name", name).
				Str("Type", typeDashboard).
				Msg("Could not build dashboard")
			return err
		}

		dashboard := NewDashboard(
			name,
			"",
			"",
			"",
			dataSourcesType,
			platform,
			builder,
		)
		jsonDashboard, errGenerate := dashboard.GetJSON()
		if errGenerate != nil {
			return errGenerate
		}

		fmt.Print(string(jsonDashboard))

		return nil
	},
}

func init() {
	GenerateCmd.Flags().String("dashboard-name", "", "Name of the dashboard to deploy")
	errName := GenerateCmd.MarkFlagRequired("dashboard-name")
	if errName != nil {
		panic(errName)
	}
	GenerateCmd.Flags().StringArray("grafana-data-sources", []string{"Prometheus"}, "Data sources to add to the dashboard, at least one required")
	GenerateCmd.Flags().String("platform", "docker", "Platform where the dashboard is deployed (docker or kubernetes)")
	GenerateCmd.Flags().String("type", "core-node", "Dashboard type can be either core-node | core-node-components")
}
