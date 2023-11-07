package dashboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
	"github.com/pkg/errors"
)

const (
	ErrFailedToCreateDashboard = "failed to create dashboard"
	ErrFailedToCreateFolder    = "failed to create folder"
)

type CLClusterDashboard struct {
	Name                     string
	LokiDataSourceName       string
	PrometheusDataSourceName string
	Folder                   string
	GrafanaURL               string
	GrafanaToken             string
	extendedOpts             []dashboard.Option
	builder                  dashboard.Builder
}

func NewCLClusterDashboard(name, ldsn, pdsn, dbf, grafanaURL, grafanaToken string, opts []dashboard.Option) (*CLClusterDashboard, error) {
	db := &CLClusterDashboard{
		Name:                     name,
		Folder:                   dbf,
		LokiDataSourceName:       ldsn,
		PrometheusDataSourceName: pdsn,
		GrafanaURL:               grafanaURL,
		GrafanaToken:             grafanaToken,
	}
	if err := db.generate(); err != nil {
		return db, err
	}
	return db, nil
}

func (m *CLClusterDashboard) nodeLogsRowOption(name, instanceSelector string) row.Option {
	return row.WithLogs(
		name,
		logs.DataSource(m.LokiDataSourceName),
		logs.Span(12),
		logs.Height("300px"),
		logs.Transparent(),
		logs.WithLokiTarget(fmt.Sprintf(`
			{namespace="${namespace}", app="app", instance="%s", container="node"}
		`, instanceSelector)),
	)
}

func (m *CLClusterDashboard) generate() error {
	builder, err := dashboard.New(
		"Chainlink Cluster Dashboard",
		dashboard.AutoRefresh("10s"),
		dashboard.Tags([]string{"generated"}),
		dashboard.VariableAsQuery(
			"namespace",
			query.DataSource(m.LokiDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "namespace")),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsInterval(
			"interval",
			interval.Values([]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}),
		),
		// logs
		dashboard.Row(
			"Logs",
			row.Collapse(),
			m.nodeLogsRowOption("Node 1", "node-1"),
			m.nodeLogsRowOption("Node 2", "node-2"),
			m.nodeLogsRowOption("Node 3", "node-3"),
			m.nodeLogsRowOption("Node 4", "node-4"),
		),
		dashboard.Row(
			"Cluster health",
			row.WithTimeSeries(
				"Restarts",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`sum(increase(kube_pod_container_status_restarts_total{namespace=~"${namespace}"}[5m])) by (pod)`,
					prometheus.Legend("{{pod}}"),
				),
			),
		),
		// FIXME: these metrics are not exposed by the node for some reason
		// DON report metrics
		dashboard.Row("DON Report metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Plugin Query() time (95th)",
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_query_time_bucket{job=~".*"}[$__rate_interval])) by (le, service)) / 1e9`,
				),
			),
			row.WithTimeSeries(
				"Plugin Observation() time (95th)",
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_observation_time_bucket{job=~".*"}[$__rate_interval])) by (le, service)) / 1e9`,
				),
			),
			row.WithTimeSeries(
				"Plugin ShouldAcceptFinalizedReport() time (95th)",
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_should_accept_finalized_report_time_bucket{job=~".*"}[$__rate_interval])) by (le, service)) / 1e9`,
				),
			),
			row.WithTimeSeries(
				"Plugin Report() time (95th)",
				timeseries.Span(6),
				timeseries.Height("300px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`	histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_report_time_bucket{job=~".*"}[$__rate_interval])) by (le, service)) / 1e9`,
				),
			),
			row.WithTimeSeries(
				"Plugin ShouldTransmitAcceptedReport() time (95th)",
				timeseries.Span(6),
				timeseries.Height("300px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_should_transmit_accepted_report_time_bucket{job=~".*"}[$__rate_interval])) by (le, service)) / 1e9`,
				),
			),
		),
	)
	m.builder = builder
	return err
}

func (m *CLClusterDashboard) Deploy() error {
	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, m.GrafanaURL, grabana.WithAPIToken(m.GrafanaToken))
	folder, err := client.FindOrCreateFolder(ctx, m.Folder)
	if err != nil {
		return errors.Wrap(err, ErrFailedToCreateFolder)
	}
	if _, err := client.UpsertDashboard(ctx, folder, m.builder); err != nil {
		return errors.Wrap(err, ErrFailedToCreateDashboard)
	}
	return nil
}
