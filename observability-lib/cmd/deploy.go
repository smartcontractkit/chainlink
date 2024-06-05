package cmd

import (
	"errors"

	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/spf13/cobra"

	atlasdon "github.com/smartcontractkit/chainlink-common/observability-lib/atlas-don"
	corenode "github.com/smartcontractkit/chainlink-common/observability-lib/core-node"
	corenodecomponents "github.com/smartcontractkit/chainlink-common/observability-lib/core-node-components"
	k8sresources "github.com/smartcontractkit/chainlink-common/observability-lib/k8s-resources"
	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Grafana Dashboards, Prometheus Alerts",
	RunE: func(cmd *cobra.Command, args []string) error {
		dataSources, errDataSources := cmd.Flags().GetStringArray("grafana-data-sources")
		if errDataSources != nil || len(dataSources) < 1 {
			return errDataSources
		}

		dataSourcesType := SetDataSources(dataSources)
		name := cmd.Flag("dashboard-name").Value.String()
		platform := cmd.Flag("platform").Value.String()
		url := cmd.Flag("grafana-url").Value.String()
		folder := cmd.Flag("dashboard-folder").Value.String()
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
				Str("URL", url).
				Str("Folder", folder).
				Str("Type", typeDashboard).
				Msg("Could not build dashboard")
			return err
		}

		dashboard := NewDashboard(
			name,
			cmd.Flag("grafana-token").Value.String(),
			url,
			folder,
			dataSourcesType,
			platform,
			builder,
		)
		errDeploy := dashboard.Deploy()
		if errDeploy != nil {
			return errDeploy
		}

		return nil
	},
}

func init() {
	DeployCmd.Flags().String("dashboard-name", "", "Name of the dashboard to deploy")
	errName := DeployCmd.MarkFlagRequired("dashboard-name")
	if errName != nil {
		panic(errName)
	}
	DeployCmd.Flags().String("dashboard-folder", "", "Dashboard folder")
	errFolder := DeployCmd.MarkFlagRequired("dashboard-folder")
	if errFolder != nil {
		panic(errFolder)
	}
	DeployCmd.Flags().String("grafana-url", "", "Grafana URL")
	errURL := DeployCmd.MarkFlagRequired("grafana-url")
	if errURL != nil {
		panic(errURL)
	}
	DeployCmd.Flags().String("grafana-token", "", "Grafana API token")
	errToken := DeployCmd.MarkFlagRequired("grafana-token")
	if errToken != nil {
		panic(errToken)
	}
	DeployCmd.Flags().StringArray("grafana-data-sources", []string{"Prometheus"}, "Data sources to add to the dashboard, at least one required")
	DeployCmd.Flags().String("platform", "docker", "Platform where the dashboard is deployed (docker or kubernetes)")
	DeployCmd.Flags().String("type", "core-node", "Dashboard type can be either core-node | core-node-components")
}
