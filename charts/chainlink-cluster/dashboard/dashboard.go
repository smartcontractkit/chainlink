package dashboard

import (
	"context"
	"fmt"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/variable/query"
	"net/http"
	"os"
)

type PanelOption struct {
	labelFilter string
}

type Dashboard struct {
	Name                     string
	grafanaURL               string
	grafanaToken             string
	grafanaFolder            string
	grafanaTags              []string
	LokiDataSourceName       string
	PrometheusDataSourceName string
	platform                 string
	panels                   []string
	panelOption              PanelOption
	Builder                  dashboard.Builder
	opts                     []dashboard.Option
	extendedOpts             []dashboard.Option
}

// NewDashboard returns a new Grafana dashboard, it comes empty and depending on user inputs panels are added to it
func NewDashboard(
	name string,
	grafanaURL string,
	grafanaToken string,
	grafanaFolder string,
	grafanaTags []string,
	lokiDataSourceName string,
	prometheusDataSourceName string,
	platform string,
	panels []string,
	extendedOpts []dashboard.Option,
) error {
	db := &Dashboard{
		Name:                     name,
		grafanaURL:               grafanaURL,
		grafanaToken:             grafanaToken,
		grafanaFolder:            grafanaFolder,
		grafanaTags:              grafanaTags,
		LokiDataSourceName:       lokiDataSourceName,
		PrometheusDataSourceName: prometheusDataSourceName,
		platform:                 platform,
		panels:                   panels,
		extendedOpts:             extendedOpts,
	}
	db.init()
	db.addCoreVariables()
	if Contains(db.panels, "core") {
		db.addCorePanels()
	}

	switch db.platform {
	case "kubernetes":
		db.addKubernetesVariables()
		db.addKubernetesPanels()
		break
	}

	db.opts = append(db.opts, db.extendedOpts...)
	err := db.deploy()
	if err != nil {
		os.Exit(1)
		return err
	}
	return nil
}

func (m *Dashboard) deploy() error {
	ctx := context.Background()

	builder, builderErr := dashboard.New(
		m.Name,
		m.opts...,
	)
	if builderErr != nil {
		fmt.Printf("Could not build dashboard: %s\n", builderErr)
		return builderErr
	}

	client := grabana.NewClient(&http.Client{}, m.grafanaURL, grabana.WithAPIToken(m.grafanaToken))
	fo, folderErr := client.FindOrCreateFolder(ctx, m.grafanaFolder)
	if folderErr != nil {
		fmt.Printf("Could not find or create folder: %s\n", folderErr)
		return folderErr
	}
	if _, err := client.UpsertDashboard(ctx, fo, builder); err != nil {
		fmt.Printf("Could not upsert dashboard: %s\n", err)
		return err
	}

	return nil
}

func (m *Dashboard) init() {
	opts := []dashboard.Option{
		dashboard.AutoRefresh("10s"),
		dashboard.Tags(m.grafanaTags),
	}

	switch m.platform {
	case "kubernetes":
		m.panelOption.labelFilter = "job"
		break
	case "docker":
		m.panelOption.labelFilter = "instance"
		break
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addCoreVariables() {
	opts := []dashboard.Option{
		dashboard.VariableAsQuery(
			"instance",
			query.DataSource(m.PrometheusDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", m.panelOption.labelFilter)),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"evmChainID",
			query.DataSource(m.PrometheusDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "evmChainID")),
			query.Sort(query.NumericalAsc),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addKubernetesVariables() {
	opts := []dashboard.Option{
		dashboard.VariableAsQuery(
			"namespace",
			query.DataSource(m.LokiDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "namespace")),
			query.Sort(query.NumericalAsc),
		),
	}

	m.opts = append(m.opts, opts...)
}
