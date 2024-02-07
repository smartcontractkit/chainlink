package dashboard

import (
	"fmt"
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/query"
)

func (m *Dashboard) init() {
	opts := []dashboard.Option{
		dashboard.AutoRefresh("10s"),
		dashboard.Tags([]string{"Chainlink Dashboard"}),
	}

	m.opts = opts
}

func (m *Dashboard) addVariables() {
	opts := []dashboard.Option{
		dashboard.VariableAsQuery(
			"instance",
			query.DataSource(m.LokiDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "instance")),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"evmChainID",
			query.DataSource(m.LokiDataSourceName),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "evmChainID")),
			query.Sort(query.NumericalAsc),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addMainPanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"Global health",
			row.WithStat(
				"App Version",
				stat.DataSource(m.PrometheusDataSourceName),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationVertical),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(2),
				stat.Height("100px"),
				stat.WithPrometheusTarget(
					`version{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - {{version}}"),
				),
			),
			row.WithStat(
				"Go Version",
				stat.DataSource(m.PrometheusDataSourceName),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationVertical),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(2),
				stat.Height("100px"),
				stat.WithPrometheusTarget(
					`go_info{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - {{version}}"),
				),
			),
			row.WithStat(
				"Uptime in seconds",
				stat.DataSource(m.PrometheusDataSourceName),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationVertical),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(2),
				stat.Height("100px"),
				stat.WithPrometheusTarget(
					`uptime_seconds{instance=~"$instance"}`,
					prometheus.Legend("{{instance}}"),
				),
			),
			row.WithStat(
				"ETH Balance",
				stat.DataSource(m.PrometheusDataSourceName),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationVertical),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.Height("100px"),
				stat.WithPrometheusTarget(
					`eth_balance{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - {{account}}"),
				),
			),
			row.WithTimeSeries(
				"Service Components Health",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`health{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - {{service_id}}"),
				),
			),
			row.WithTimeSeries(
				"ETH Balance",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`eth_balance{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - {{account}}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addLogPollerPanels() {
	opts := []dashboard.Option{
		dashboard.Row("LogPoller",
			row.Collapse(),
			row.WithTimeSeries(
				"LogPoller RPS",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`avg(sum(rate(log_poller_query_duration_count{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (query, instance)) by (query)`,
					prometheus.Legend("{{query}}"),
				),
				timeseries.WithPrometheusTarget(
					`avg(sum(rate(log_poller_query_duration_count{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval]))) by (instance)`,
					prometheus.Legend("Total"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Logs Number Returned",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`log_poller_query_dataset_size{instance=~"$instance", evmChainID=~"$evmChainID"}`,
					prometheus.Legend("{{query}} : {{type}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Average Logs Number Returned",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`avg(log_poller_query_dataset_size{instance=~"$instance", evmChainID=~"$evmChainID"}) by (query)`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Max Logs Number Returned",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`max(log_poller_query_dataset_size{instance=~"$instance", evmChainID=~"$evmChainID"}) by (query)`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Logs Number Returned by Chain",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`max(log_poller_query_dataset_size{instance=~"$instance"}) by (evmChainID)`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Queries Duration Avg",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`(sum(rate(log_poller_query_duration_sum{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (query) / sum(rate(log_poller_query_duration_count{instance=~"$instance"}[$__rate_interval])) by (query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Queries Duration p99",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(log_poller_query_duration_bucket{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Queries Duration p95",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(log_poller_query_duration_bucket{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Queries Duration p90",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(log_poller_query_duration_bucket{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"LogPoller Queries Duration Median",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.5, sum(rate(log_poller_query_duration_bucket{instance=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addFeedsJobsPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Feeds Jobs",
			row.Collapse(),
			row.WithTimeSeries(
				"Feeds Job Proposal Requests",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`feeds_job_proposal_requests{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Feeds Job Proposal Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`feeds_job_proposal_count{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addMailboxPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Mailbox",
			row.Collapse(),
			row.WithTimeSeries(
				"Mailbox Load Percent",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`mailbox_load_percent{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{ name }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addPromReporterPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Prom Reporter",
			row.Collapse(),
			row.WithTimeSeries(
				"Unconfirmed Transactions",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Tx"),
				),
				timeseries.WithPrometheusTarget(
					`unconfirmed_transactions{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Unconfirmed TX Age",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`max_unconfirmed_tx_age{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Unconfirmed TX Blocks",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Blocks"),
				),
				timeseries.WithPrometheusTarget(
					`max_unconfirmed_blocks{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addTxManagerPanels() {
	opts := []dashboard.Option{
		dashboard.Row("TX Manager",
			row.Collapse(),
			row.WithTimeSeries(
				"TX Manager Time Until TX Broadcast",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_time_until_tx_broadcast{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Gas Bumps",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_gas_bumps{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Gas Bumps Exceeds Limit",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_gas_bump_exceeds_limit{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Confirmed Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_confirmed_transactions{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Successful Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_successful_transactions{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Reverted Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_tx_reverted{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Fwd Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_fwd_tx_count{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Transactions Attempts",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_tx_attempt_count{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Time Until TX Confirmed",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_time_until_tx_confirmed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Block Until TX Confirmed",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_blocks_until_tx_confirmed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addHeadTrackerPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Head tracker",
			row.Collapse(),
			row.WithTimeSeries(
				"Head tracker current head",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_current_head{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker very old head",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_very_old_head{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker heads received",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_heads_received{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker connection errors",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_connection_errors{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addDatabasePanels() {
	opts := []dashboard.Option{
		// DB Metrics
		dashboard.Row("DB Connection Metrics (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"DB Connections",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Conn"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_max{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - Max}"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_open{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - Open"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_used{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - Used"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_wait{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - Wait"),
				),
			),
			row.WithTimeSeries(
				"DB Wait Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`db_wait_count{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"DB Wait Time",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`db_wait_time_seconds{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addSQLQueryPanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"SQL Query",
			row.Collapse(),
			row.WithTimeSeries(
				"SQL Query Timeout Percent",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("percent"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.9, sum(rate(sql_query_timeout_percent_bucket{instance=~"$instance"}[$__rate_interval])) by (le))`,
					prometheus.Legend("p90"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(sql_query_timeout_percent_bucket{instance=~"$instance"}[$__rate_interval])) by (le))`,
					prometheus.Legend("p95"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(sql_query_timeout_percent_bucket{instance=~"$instance"}[$__rate_interval])) by (le))`,
					prometheus.Legend("p99"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addLogsPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Logs Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Log Counters",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`log_panic_count{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - panic"),
				),
				timeseries.WithPrometheusTarget(
					`log_fatal_count{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - fatal"),
				),
				timeseries.WithPrometheusTarget(
					`log_critical_count{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - critical"),
				),
				timeseries.WithPrometheusTarget(
					`log_warn_count{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - warn"),
				),
				timeseries.WithPrometheusTarget(
					`log_error_count{instance=~"$instance"}`,
					prometheus.Legend("{{instance}} - error"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addEVMPoolLifecyclePanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"EVM Pool Lifecycle",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool Highest Seen Block",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_highest_seen_block{instance=~"$instance", evmChainID="${evmChainID}"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Num Seen Blocks",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_seen_blocks{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Failed",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_failed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_success{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addEVMPoolRPCNodePanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"EVM Pool RPC Node Metrics (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_calls_success{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_calls_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_dials_success{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Failed",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_dials_failed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_dials_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Failed",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_dials_failed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to Alive",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_alive{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to In Sync",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_in_sync{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to Out of Sync",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_out_of_sync{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to Unreachable",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_unreachable{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to invalid Chain ID",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_invalid_chain_id{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Total Transitions to unusable",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_unusable{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Polls Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_success{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Polls Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node States",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_states{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{evmChainID}} - {{state}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Verifies Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_verifies{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Verifies Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_verifies_success{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Verifies Failed",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_verifies_failed{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{evmChainID}}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addEVMRPCNodeLatenciesPanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"EVM Pool RPC Node Latencies (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Latency 0.95 quantile",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{instance=~"$instance"}[$__rate_interval])) by (le, rpcCallName)) / 1e6`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addBlockHistoryEstimatorPanels() {
	opts := []dashboard.Option{
		dashboard.Row("Block History Estimator",
			row.Collapse(),
			row.WithTimeSeries(
				"Gas Updater All Gas Price Percentiles",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_all_gas_price_percentiles{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{ percentile }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater All Tip Cap Percentiles",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_all_tip_cap_percentiles{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} - {{ percentile }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Set Gas Price",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_set_gas_price{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Set Tip Cap",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_set_tip_cap{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Current Base Fee",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_current_base_fee{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Block History Estimator Connectivity Failure Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`block_history_estimator_connectivity_failure_count{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addPipelinePanels() {
	opts := []dashboard.Option{
		dashboard.Row("Pipeline Metrics (Runner)",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Task Execution Time",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_execution_time{instance=~"$instance"} / 1e6`,
					prometheus.Legend("{{ instance }} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Run Errors",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_run_errors{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Run Total Time to Completion",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_run_total_time_to_completion{instance=~"$instance"} / 1e6`,
					prometheus.Legend("{{ instance }} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Tasks Total Finished",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_tasks_total_finished{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }} JobID: {{ job_id }}"),
				),
			),
		),
		dashboard.Row(
			"Pipeline Metrics (ETHCall)",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Task ETH Call Execution Time",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_eth_call_execution_time{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
		dashboard.Row(
			"Pipeline Metrics (HTTP)",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Task HTTP Fetch Time",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_http_fetch_time{instance=~"$instance"} / 1e6`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Task HTTP Response Body Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_http_response_body_size{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
		dashboard.Row(
			"Pipeline Metrics (Bridge)",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Bridge Latency",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`bridge_latency_seconds{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Errors Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_errors_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Cache Hits Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_cache_hits_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Cache Errors Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_cache_errors_total{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
		dashboard.Row(
			"Pipeline Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Runs Queued",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_runs_queued{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Runs Tasks Queued",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_runs_queued{instance=~"$instance"}`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addHTTPAPIPanels() {
	opts := []dashboard.Option{
		// HTTP API Metrics
		dashboard.Row(
			"HTTP API Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Request Duration p95",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(service_gonic_request_duration_bucket{instance=~"$instance"}[$__rate_interval])) by (le, path, method))`,
					prometheus.Legend("{{ method }} {{ path }}"),
				),
			),
			row.WithTimeSeries(
				"Request Total Rate over interval",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(service_gonic_requests_total{instance=~"$instance"}[$__rate_interval])) by (path, method, code)`,
					prometheus.Legend("{{ method }} {{ path }} {{ code }}"),
				),
			),
			row.WithTimeSeries(
				"Request Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`avg(rate(service_gonic_request_size_bytes_sum{instance=~"$instance"}[$__rate_interval]))/avg(rate(service_gonic_request_size_bytes_count{instance=~"$instance"}[$__rate_interval]))`,
					prometheus.Legend("Average"),
				),
			),
			row.WithTimeSeries(
				"Response Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`avg(rate(service_gonic_response_size_bytes_sum{instance=~"$instance"}[$__rate_interval]))/avg(rate(service_gonic_response_size_bytes_count{instance=~"$instance"}[$__rate_interval]))`,
					prometheus.Legend("Average"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addPromHTTPPanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"PromHTTP Metrics",
			row.Collapse(),
			row.WithGauge("HTTP Request in flight",
				gauge.Span(2),
				gauge.Height("200px"),
				gauge.DataSource(m.PrometheusDataSourceName),
				gauge.WithPrometheusTarget(
					`promhttp_metric_handler_requests_in_flight`,
					prometheus.Legend(""),
				),
			),
			row.WithTimeSeries(
				"HTTP rate",
				timeseries.Span(10),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(promhttp_metric_handler_requests_total{instance=~"$instance"}[$__rate_interval])) by (code)`,
					prometheus.Legend("{{ code }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) addGoMetricsPanels() {
	opts := []dashboard.Option{
		dashboard.Row(
			"Go Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Heap allocations",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.WithPrometheusTarget(
					`sum(go_memstats_heap_alloc_bytes{instance=~"$instance"}) by (instance)`,
					prometheus.Legend("{{ instance }}"),
				),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
				timeseries.Legend(timeseries.Last, timeseries.AsTable),
				timeseries.Alert(
					"Too many heap allocations",
					alert.Description("Above 128Mib on {{ instance }}"),
					alert.Runbook("https://google.com"),
					alert.Tags(map[string]string{
						"service": "chainlink-node",
						"owner":   "team-core",
					}),
					alert.WithPrometheusQuery(
						"A",
						`sum(go_memstats_heap_alloc_bytes{instance=~"$instance"}) by (instance)`,
					),
					alert.If(alert.Avg, "A", alert.IsAbove(1.342e+8)), // 128Mib
				),
			),
			row.WithStat(
				"Heap Allocations",
				stat.Span(6),
				stat.Height("200px"),
				stat.DataSource(m.PrometheusDataSourceName),
				stat.Unit("bytes"),
				stat.ColorValue(),
				stat.WithPrometheusTarget("sum(go_memstats_heap_alloc_bytes)"),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{
						Color: "green",
						Value: nil,
					},
					{
						Color: "orange",
						Value: float64Ptr(6.711e+7),
					},
					{
						Color: "red",
						Value: float64Ptr(1.342e+8),
					},
				}),
			),
			row.WithTable(
				"Threads",
				table.Span(3),
				table.Height("200px"),
				table.DataSource(m.PrometheusDataSourceName),
				table.WithPrometheusTarget(
					`sum(go_threads{instance=~"$instance"}) by (instance)`,
					prometheus.Legend("{{ instance }}")),
				table.HideColumn("Time"),
				table.AsTimeSeriesAggregations([]table.Aggregation{
					{Label: "AVG", Type: table.AVG},
					{Label: "Current", Type: table.Current},
				}),
			),
			row.WithTimeSeries(
				"Threads",
				timeseries.Span(9),
				timeseries.Height("200px"),
				timeseries.DataSource(m.PrometheusDataSourceName),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(go_threads{instance=~"$instance"}) by (instance)`,
					prometheus.Legend("{{ instance }}"),
				),
			),
		),
	}

	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) AddPanels(opts []dashboard.Option) {
	m.opts = append(m.opts, opts...)
}

func (m *Dashboard) generate() error {
	m.init()
	m.addVariables()
	m.addMainPanels()
	m.addLogPollerPanels()
	m.addFeedsJobsPanels()
	m.addMailboxPanels()
	m.addPromReporterPanels()
	m.addTxManagerPanels()
	m.addHeadTrackerPanels()
	m.addDatabasePanels()
	m.addSQLQueryPanels()
	m.addLogsPanels()
	m.addEVMPoolLifecyclePanels()
	m.addEVMPoolRPCNodePanels()
	m.addEVMRPCNodeLatenciesPanels()
	m.addBlockHistoryEstimatorPanels()
	m.addPipelinePanels()
	m.addHTTPAPIPanels()
	m.addPromHTTPPanels()
	m.addGoMetricsPanels()

	opts := append(m.opts, m.extendedOpts...)

	builder, err := dashboard.New(
		m.Name,
		opts...,
	)
	m.Builder = builder
	return err
}

func float64Ptr(input float64) *float64 {
	return &input
}
