package dashboard

import (
	"context"
	"fmt"
	"net/http"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
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

// CLClusterDashboard is a dashboard for a Chainlink cluster
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

// NewCLClusterDashboard returns a new dashboard for a Chainlink cluster, can be used as a base for more complex plugin based dashboards
func NewCLClusterDashboard(name, ldsn, pdsn, dbf, grafanaURL, grafanaToken string, opts []dashboard.Option) (*CLClusterDashboard, error) {
	db := &CLClusterDashboard{
		Name:                     name,
		Folder:                   dbf,
		LokiDataSourceName:       ldsn,
		PrometheusDataSourceName: pdsn,
		GrafanaURL:               grafanaURL,
		GrafanaToken:             grafanaToken,
		extendedOpts:             opts,
	}
	if err := db.generate(); err != nil {
		return db, err
	}
	return db, nil
}

// logsRowOption returns a row option for a node's logs with name and instance selector
func (m *CLClusterDashboard) logsRowOption(name, instanceSelector string) row.Option {
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

// timeseriesRowOption returns a row option for a timeseries with name, axis unit, query and legend template
func (m *CLClusterDashboard) timeseriesRowOption(name, axisUnit, query, legendTemplate string) row.Option {
	var tsq timeseries.Option
	if legendTemplate != "" {
		tsq = timeseries.WithPrometheusTarget(
			query,
			prometheus.Legend(legendTemplate),
		)
	} else {
		tsq = timeseries.WithPrometheusTarget(query)
	}
	var au timeseries.Option
	if axisUnit != "" {
		au = timeseries.Axis(
			axis.Unit(axisUnit),
		)
	} else {
		au = timeseries.Axis()
	}
	return row.WithTimeSeries(
		name,
		timeseries.Span(6),
		timeseries.Height("300px"),
		timeseries.DataSource(m.PrometheusDataSourceName),
		au,
		tsq,
	)
}

// statRowOption returns a row option for a stat with name, prometheus target and legend template
func (m *CLClusterDashboard) statRowOption(name, target, legend string) row.Option {
	return row.WithStat(
		name,
		stat.Transparent(),
		stat.DataSource(m.PrometheusDataSourceName),
		stat.Text(stat.TextValueAndName),
		stat.Orientation(stat.OrientationVertical),
		stat.TitleFontSize(12),
		stat.ValueFontSize(20),
		stat.Span(12),
		stat.Height("100px"),
		stat.WithPrometheusTarget(target, prometheus.Legend(legend)),
	)
}

// generate generates the dashboard, adding extendedOpts to the default options
func (m *CLClusterDashboard) generate() error {
	opts := []dashboard.Option{
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
		dashboard.Row(
			"Cluster health",
			m.statRowOption(
				"App Version",
				`version{namespace="${namespace}"}`,
				"{{pod}} - {{version}}",
			),
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
			row.WithTimeSeries(
				"Service Components Health",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`health{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - {{service_id}}"),
				),
			),
			row.WithTimeSeries(
				"Log Counters",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`log_panic_count{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - panic"),
				),
				timeseries.WithPrometheusTarget(
					`log_fatal_count{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - fatal"),
				),
				timeseries.WithPrometheusTarget(
					`log_critical_count{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - critical"),
				),
				timeseries.WithPrometheusTarget(
					`log_error_count{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - error"),
				),
			),
			row.WithTimeSeries(
				"ETH Balance",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`eth_balance{namespace="${namespace}"}`,
					prometheus.Legend("{{pod}} - {{account}}"),
				),
			),
		),
		// logs
		dashboard.Row(
			"Logs",
			row.Collapse(),
			m.logsRowOption("Node 1", "node-1"),
			m.logsRowOption("Node 2", "node-2"),
			m.logsRowOption("Node 3", "node-3"),
			m.logsRowOption("Node 4", "node-4"),
		),
		// DON report metrics
		dashboard.Row("DON Report metrics",
			row.Collapse(),
			m.timeseriesRowOption(
				"Plugin Query() count",
				"Count",
				`sum(rate(ocr2_reporting_plugin_query_count{namespace="${namespace}", app="app"}[$__rate_interval])) by (service)`,
				"",
			),
			m.timeseriesRowOption(
				"Plugin Observation() time (95th)",
				"Sec",
				`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_observation_time_bucket{namespace="${namespace}", app="app"}[$__rate_interval])) by (le, service)) / 1e9`,
				"",
			),
			m.timeseriesRowOption(
				"Plugin ShouldAcceptReport() time (95th)",
				"Sec",
				`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_should_accept_report_time_bucket{namespace="${namespace}", app="app"}[$__rate_interval])) by (le, service)) / 1e9`,
				"",
			),
			m.timeseriesRowOption(
				"Plugin Report() time (95th)",
				"Sec",
				`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_report_time_bucket{namespace="${namespace}", app="app"}[$__rate_interval])) by (le, service)) / 1e9`,
				"",
			),
			m.timeseriesRowOption(
				"Plugin ShouldTransmitReport() time (95th)",
				"Sec",
				`histogram_quantile(0.95, sum(rate(ocr2_reporting_plugin_should_transmit_report_time_bucket{namespace="${namespace}", app="app"}[$__rate_interval])) by (le, service)) / 1e9`,
				"",
			),
		),
		dashboard.Row(
			"DB Connection Metrics (App)",
			row.Collapse(),
			m.timeseriesRowOption(
				"DB Connections MAX",
				"Conn",
				`db_conns_max{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"DB Connections Open",
				"Conn",
				`db_conns_open{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"DB Connections Used",
				"Conn",
				`db_conns_used{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"DB Connections Wait",
				"Conn",
				`db_conns_wait{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"DB Wait time",
				"Sec",
				`db_wait_time_seconds{namespace="${namespace}"}`,
				"{{pod}}",
			),
		),
		dashboard.Row(
			"EVM Pool RPC Node Metrics (App)",
			row.Collapse(),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Calls Success",
				"",
				`evm_pool_rpc_node_calls_success{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Calls Total",
				"",
				`evm_pool_rpc_node_calls_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Dials Success",
				"",
				`evm_pool_rpc_node_dials_success{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Dials Total",
				"",
				`evm_pool_rpc_node_dials_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Highest Seen Block",
				"",
				`evm_pool_rpc_node_highest_seen_block{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to Alive",
				"",
				`evm_pool_rpc_node_num_transitions_to_alive{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Polls Success",
				"",
				`evm_pool_rpc_node_polls_success{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Polls Total",
				"",
				`evm_pool_rpc_node_polls_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node States",
				"",
				`evm_pool_rpc_node_states{namespace="${namespace}"}`,
				"{{pod}} - {{evmChainID}} - {{state}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Verifies Total",
				"",
				`evm_pool_rpc_node_verifies{namespace="${namespace}"}`,
				"{{pod}} - {{evmChainID}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Verifies Success",
				"",
				`evm_pool_rpc_node_verifies_success{namespace="${namespace}"}`,
				"{{pod}} - {{evmChainID}}",
			),
		),
		dashboard.Row(
			"EVM Pool RPC Node Latencies (App)",
			row.Collapse(),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Calls Latency 0.95 quantile",
				"ms",
				`histogram_quantile(0.95, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{namespace="${namespace}"}[$__rate_interval])) by (le, rpcCallName)) / 1e6`,
				"{{pod}}",
			),
		),
		dashboard.Row(
			"Pipeline Tasks Metrics (App)",
			row.Collapse(),
			m.timeseriesRowOption(
				"Pipeline Runs Queued",
				"",
				`pipeline_runs_queued{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"Pipeline Runs Tasks Queued",
				"",
				`pipeline_task_runs_queued{namespace="${namespace}"}`,
				"{{pod}}",
			),
		),
	}
	opts = append(opts, m.extendedOpts...)
	builder, err := dashboard.New(
		"Chainlink Cluster Dashboard",
		opts...,
	)
	m.builder = builder
	return err
}

// Deploy deploys the dashboard to Grafana
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
