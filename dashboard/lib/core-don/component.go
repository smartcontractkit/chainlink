package core_don

import (
	"fmt"
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

type Props struct {
	PrometheusDataSource string
	PlatformOpts         PlatformOpts
}

func vars(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.VariableAsQuery(
			"instance",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", p.PlatformOpts.LabelFilter)),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"evmChainID",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "evmChainID")),
			query.Sort(query.NumericalAsc),
		),
	}
}

func generalInfoRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"General CL Cluster Info",
			row.Collapse(),
			row.WithStat(
				"App Version",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationAuto),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(2),
				stat.Text("name"),
				stat.WithPrometheusTarget(
					`version{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{version}}"),
				),
			),
			row.WithStat(
				"Go Version",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationAuto),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(2),
				stat.Text("name"),
				stat.WithPrometheusTarget(
					`go_info{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{version}}"),
				),
			),
			row.WithStat(
				"Uptime in days",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(8),
				stat.WithPrometheusTarget(
					`uptime_seconds{`+p.PlatformOpts.LabelQuery+`} / 86400`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithStat(
				"ETH Balance",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.Decimals(2),
				stat.WithPrometheusTarget(
					`eth_balance{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{account}}"),
				),
			),
			row.WithStat(
				"Solana Balance",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.Decimals(2),
				stat.WithPrometheusTarget(
					`solana_balance{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LabelFilter+"}} - {{account}}"),
				),
			),
			row.WithTimeSeries(
				"Service Components Health",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`health{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{service_id}}"),
				),
			),
			row.WithTimeSeries(
				"ETH Balance",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Decimals(2),
				),
				timeseries.WithPrometheusTarget(
					`eth_balance{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{account}}"),
				),
			),
			row.WithTimeSeries(
				"SOL Balance",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Decimals(2),
				),
				timeseries.WithPrometheusTarget(
					`solana_balance{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{account}}"),
				),
			),
		),
	}
}

func logPollerRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("LogPoller",
			row.Collapse(),
			row.WithStat(
				"Goroutines",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationAuto),
				stat.Height("200px"),
				stat.TitleFontSize(30),
				stat.ValueFontSize(30),
				stat.Span(6),
				stat.Text("Goroutines"),
				stat.WithPrometheusTarget(
					`count(count by (evmChainID) (log_poller_query_duration_sum{job=~"$instance"}))`,
					prometheus.Legend("Goroutines"),
				),
			),
			row.WithTimeSeries(
				"RPS",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("requests"),
				),
				timeseries.WithPrometheusTarget(
					`avg by (query) (sum by (query, job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
					prometheus.Legend("{{query}} - {{job}}"),
				),
				timeseries.WithPrometheusTarget(
					`avg (sum by(job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
					prometheus.Legend("Total"),
				),
			),
			row.WithTimeSeries(
				"RPS by type",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("requests"),
				),
				timeseries.WithPrometheusTarget(
					`avg by (type) (sum by (type, job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
				),
			),
			row.WithTimeSeries(
				"Avg number of logs returned",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("logs"),
				),
				timeseries.WithPrometheusTarget(
					`avg by (query) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
					prometheus.Legend("{{query}} - {{job}}"),
				),
			),
			row.WithTimeSeries(
				"Max number of logs returned",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("logs"),
				),
				timeseries.WithPrometheusTarget(
					`max by (query) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
					prometheus.Legend("{{query}} - {{job}}"),
				),
			),
			row.WithTimeSeries(
				"Logs returned by chain",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("logs"),
				),
				timeseries.WithPrometheusTarget(
					`max by (evmChainID) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"Queries duration by type (0.5 perc)",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.5, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"queries duration by type (0.9 perc)",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.9, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"Queries duration by type (0.99 perc)",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"Queries duration by chain (0.99 perc)",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, evmChainID)) / 1e6`,
					prometheus.Legend("{{query}}"),
				),
			),
			row.WithTimeSeries(
				"Number of logs inserted",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("logs"),
				),
				timeseries.WithPrometheusTarget(
					`avg by (evmChainID) (log_poller_logs_inserted{job=~"$instance", evmChainID=~"$evmChainID"})`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"Logs insertion rate",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`avg by (evmChainID) (rate(log_poller_logs_inserted{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval]))`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"Number of blocks inserted",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("blocks"),
				),
				timeseries.WithPrometheusTarget(
					`avg by (evmChainID) (log_poller_blocks_inserted{job=~"$instance", evmChainID=~"$evmChainID"})`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
			row.WithTimeSeries(
				"Blocks insertion rate",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`avg by (evmChainID) (rate(log_poller_blocks_inserted{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval]))`,
					prometheus.Legend("{{evmChainID}}"),
				),
			),
		),
	}
}

func feedJobsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Feeds Jobs",
			row.Collapse(),
			row.WithTimeSeries(
				"Feeds Job Proposal Requests",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`feeds_job_proposal_requests{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Feeds Job Proposal Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`feeds_job_proposal_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func mailBoxRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Mailbox",
			row.Collapse(),
			row.WithTimeSeries(
				"Mailbox Load Percent",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`mailbox_load_percent{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ name }}"),
				),
			),
		),
	}
}

func promReporterRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Prom Reporter",
			row.Collapse(),
			row.WithTimeSeries(
				"Unconfirmed Transactions",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Tx"),
				),
				timeseries.WithPrometheusTarget(
					`unconfirmed_transactions{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Unconfirmed TX Age",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`max_unconfirmed_tx_age{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Unconfirmed TX Blocks",
				timeseries.Span(4),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Blocks"),
				),
				timeseries.WithPrometheusTarget(
					`max_unconfirmed_blocks{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func txManagerRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("TX Manager",
			row.Collapse(),
			row.WithTimeSeries(
				"TX Manager Time Until TX Broadcast",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_time_until_tx_broadcast{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Gas Bumps",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_gas_bumps{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Gas Bumps Exceeds Limit",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_gas_bump_exceeds_limit{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Confirmed Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_confirmed_transactions{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Successful Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_successful_transactions{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Reverted Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_num_tx_reverted{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Fwd Transactions",
				timeseries.Span(3),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_fwd_tx_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Num Transactions Attempts",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_tx_attempt_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Time Until TX Confirmed",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_time_until_tx_confirmed{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"TX Manager Block Until TX Confirmed",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`tx_manager_blocks_until_tx_confirmed{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func headTrackerRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Head tracker",
			row.Collapse(),
			row.WithTimeSeries(
				"Head tracker current head",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_current_head{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker very old head",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_very_old_head{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker heads received",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_heads_received{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Head tracker connection errors",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`head_tracker_connection_errors{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func appDBConnectionsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("DB Connection Metrics (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"DB Connections",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Conn"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_max{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Max"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_open{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Open"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_used{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Used"),
				),
				timeseries.WithPrometheusTarget(
					`db_conns_wait{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Wait"),
				),
			),
			row.WithTimeSeries(
				"DB Wait Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`db_wait_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"DB Wait Time",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`db_wait_time_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func sqlQueriesRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"SQL Query",
			row.Collapse(),
			row.WithTimeSeries(
				"SQL Query Timeout Percent",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("percent"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.9, sum(rate(sql_query_timeout_percent_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (le))`,
					prometheus.Legend("p90"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(sql_query_timeout_percent_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (le))`,
					prometheus.Legend("p95"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(sql_query_timeout_percent_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (le))`,
					prometheus.Legend("p99"),
				),
			),
		),
	}
}

func logsCountersRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Logs Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Logs Counters",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`log_panic_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - panic"),
				),
				timeseries.WithPrometheusTarget(
					`log_fatal_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - fatal"),
				),
				timeseries.WithPrometheusTarget(
					`log_critical_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - critical"),
				),
				timeseries.WithPrometheusTarget(
					`log_warn_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - warn"),
				),
				timeseries.WithPrometheusTarget(
					`log_error_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - error"),
				),
			),
			row.WithTimeSeries(
				"Logs Rate",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`sum(rate(log_panic_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - panic"),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(log_fatal_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - fatal"),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(log_critical_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - critical"),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(log_warn_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - warn"),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(log_error_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - error"),
				),
			),
		),
	}
}

// TODO: fix, no data points for OCRv1
func evmPoolLifecycleRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"EVM Pool Lifecycle",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool Highest Seen Block",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_highest_seen_block{`+p.PlatformOpts.LabelQuery+`evmChainID="${evmChainID}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Num Seen Blocks",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_seen_blocks{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_total{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Failed",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_failed{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool Node Polls Success",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Block"),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_polls_success{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func nodesRPCRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"Node RPC State",
			row.Collapse(),
			row.WithStat(
				"Node RPC Alive",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Alive"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC Closed",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Closed"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC Dialed",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Dialed"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC InvalidChainID",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="InvalidChainID"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC OutOfSync",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="OutOfSync"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC UnDialed",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Undialed"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC Unreachable",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Unreachable"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
			row.WithStat(
				"Node RPC Unusable",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(3),
				stat.WithPrometheusTarget(
					`sum(multi_node_states{`+p.PlatformOpts.LabelQuery+`chainId=~"$evmChainID", state="Unusable"}) by (`+p.PlatformOpts.LegendString+`, chainId)`,
					prometheus.Legend("{{pod}} - {{chainId}}"),
				),
			),
		),
	}
}

func evmNodeRPCRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"EVM Pool RPC Node Metrics (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Success Rate",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Label("%"),
					axis.SoftMin(0),
					axis.SoftMax(100),
				),
				timeseries.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_calls_success{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_calls_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) * 100`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
			),
			row.WithGauge(
				"EVM Pool RPC Node Calls Success Rate",
				gauge.Span(12),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.Unit("percentunit"),
				gauge.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_calls_success{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_calls_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
				gauge.AbsoluteThresholds([]gauge.ThresholdStep{
					{Color: "#ff0000"},
					{Color: "#ffa500", Value: float64Ptr(0.8)},
					{Color: "#00ff00", Value: float64Ptr(0.9)},
				}),
			),
			// issue when value is 0
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Success Rate",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Label("%"),
					axis.SoftMin(0),
					axis.SoftMax(100),
				),
				timeseries.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_dials_success{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_dials_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) * 100`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
			),
			// issue when value is 0
			row.WithTimeSeries(
				"EVM Pool RPC Node Dials Failure Rate",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Label("%"),
					axis.SoftMin(0),
					axis.SoftMax(100),
				),
				timeseries.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_dials_failed{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_dials_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) * 100`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Transitions",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_alive{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_in_sync{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_out_of_sync{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_unreachable{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_invalid_chain_id{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_num_transitions_to_unusable{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node States",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`evm_pool_rpc_node_states{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{state}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Verifies Success Rate",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Label("%"),
					axis.SoftMin(0),
					axis.SoftMax(100),
				),
				timeseries.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_verifies_success{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_verifies{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) * 100`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Verifies Failure Rate",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
					axis.Label("%"),
					axis.SoftMin(0),
					axis.SoftMax(100),
				),
				timeseries.WithPrometheusTarget(
					`sum(increase(evm_pool_rpc_node_verifies_failed{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_verifies{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, evmChainID, nodeName) * 100`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{evmChainID}} - {{nodeName}}"),
				),
			),
		),
	}
}

func evmRPCNodeLatenciesRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"EVM Pool RPC Node Latencies (App)",
			row.Collapse(),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Latency 0.90 quantile",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.90, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, le, rpcCallName)) / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{rpcCallName}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Latency 0.95 quantile",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, le, rpcCallName)) / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{rpcCallName}}"),
				),
			),
			row.WithTimeSeries(
				"EVM Pool RPC Node Calls Latency 0.99 quantile",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("ms"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.99, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, le, rpcCallName)) / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{rpcCallName}}"),
				),
			),
		),
	}
}

func evmBlockHistoryEstimatorRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Block History Estimator",
			row.Collapse(),
			row.WithTimeSeries(
				"Gas Updater All Gas Price Percentiles",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_all_gas_price_percentiles{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ percentile }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater All Tip Cap Percentiles",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_all_tip_cap_percentiles{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ percentile }}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Set Gas Price",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_set_gas_price{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Set Tip Cap",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_set_tip_cap{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Gas Updater Current Base Fee",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`gas_updater_current_base_fee{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Block History Estimator Connectivity Failure Count",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`block_history_estimator_connectivity_failure_count{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func pipelinesRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Pipeline Metrics (Runner)",
			row.Collapse(),
			row.WithTimeSeries(
				"Pipeline Task Execution Time",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_execution_time{`+p.PlatformOpts.LabelQuery+`} / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Run Errors",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_run_errors{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Run Total Time to Completion",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_run_total_time_to_completion{`+p.PlatformOpts.LabelQuery+`} / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} JobID: {{ job_id }}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Tasks Total Finished",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_tasks_total_finished{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} JobID: {{ job_id }}"),
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
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_eth_call_execution_time{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
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
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_http_fetch_time{`+p.PlatformOpts.LabelQuery+`} / 1e6`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Task HTTP Response Body Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_http_response_body_size{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
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
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`bridge_latency_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Errors Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_errors_total{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Cache Hits Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_cache_hits_total{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Bridge Cache Errors Total",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`bridge_cache_errors_total{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
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
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_runs_queued{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Pipeline Runs Tasks Queued",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`pipeline_task_runs_queued{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func httpAPIRow(p Props) []dashboard.Option {
	return []dashboard.Option{

		dashboard.Row(
			"HTTP API Metrics",
			row.Collapse(),
			row.WithTimeSeries(
				"Request Duration p95",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Sec"),
				),
				timeseries.WithPrometheusTarget(
					`histogram_quantile(0.95, sum(rate(service_gonic_request_duration_bucket{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, le, path, method))`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ method }} - {{ path }}"),
				),
			),
			row.WithTimeSeries(
				"Request Total Rate over interval",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(service_gonic_requests_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, path, method, code)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ method }} - {{ path }} - {{ code }}"),
				),
			),
			row.WithTimeSeries(
				"Average Request Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`avg(rate(service_gonic_request_size_bytes_sum{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)/avg(rate(service_gonic_request_size_bytes_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"Response Size",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("Bytes"),
				),
				timeseries.WithPrometheusTarget(
					`avg(rate(service_gonic_response_size_bytes_sum{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)/avg(rate(service_gonic_response_size_bytes_count{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func promHTTPRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"PromHTTP Metrics",
			row.Collapse(),
			row.WithGauge("HTTP Request in flight",
				gauge.Span(12),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					`promhttp_metric_handler_requests_in_flight{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithTimeSeries(
				"HTTP rate",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(rate(promhttp_metric_handler_requests_total{`+p.PlatformOpts.LabelQuery+`}[$__rate_interval])) by (`+p.PlatformOpts.LegendString+`, code)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
		),
	}
}

func goMetricsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"Go Metrics",
			row.Collapse(),
			row.WithTable(
				"Threads",
				table.Span(3),
				table.Height("200px"),
				table.DataSource(p.PrometheusDataSource),
				table.WithPrometheusTarget(
					`sum(go_threads{`+p.PlatformOpts.LabelQuery+`}) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}")),
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
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit(""),
				),
				timeseries.WithPrometheusTarget(
					`sum(go_threads{`+p.PlatformOpts.LabelQuery+`}) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
			),
			row.WithStat(
				"Heap Allocations",
				stat.Span(12),
				stat.Orientation(stat.OrientationVertical),
				stat.DataSource(p.PrometheusDataSource),
				stat.Unit("bytes"),
				stat.ColorValue(),
				stat.WithPrometheusTarget(
					`sum(go_memstats_heap_alloc_bytes{`+p.PlatformOpts.LabelQuery+`}) by (`+p.PlatformOpts.LegendString+`)`,
				),
			),
			row.WithTimeSeries(
				"Heap allocations",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`sum(go_memstats_heap_alloc_bytes{`+p.PlatformOpts.LabelQuery+`}) by (`+p.PlatformOpts.LegendString+`)`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Memory in Heap",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_heap_alloc_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Alloc"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_heap_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Sys"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_heap_idle_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Idle"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_heap_inuse_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - InUse"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_heap_released_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Released"),
				),
			),
			row.WithTimeSeries(
				"Memory in Off-Heap",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_mspan_inuse_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Total InUse"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_mspan_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Total Sys"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_mcache_inuse_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Cache InUse"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_mcache_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Cache Sys"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_buck_hash_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Hash Sys"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_gc_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - GC Sys"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_other_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - bytes of memory are used for other runtime allocations"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_next_gc_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Next GC"),
				),
			),
			row.WithTimeSeries(
				"Memory in Stack",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`go_memstats_stack_inuse_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - InUse"),
				),
				timeseries.WithPrometheusTarget(
					`go_memstats_stack_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - Sys"),
				),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Total Used Memory",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`go_memstats_sys_bytes{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Number of Live Objects",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`go_memstats_mallocs_total{`+p.PlatformOpts.LabelQuery+`} - go_memstats_frees_total{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Rate of Objects Allocated",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`rate(go_memstats_mallocs_total{`+p.PlatformOpts.LabelQuery+`}[1m])`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Rate of a Pointer Dereferences",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`rate(go_memstats_lookups_total{`+p.PlatformOpts.LabelQuery+`}[1m])`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.Unit("ops"),
					axis.SoftMin(0),
				),
			),
			row.WithTimeSeries(
				"Goroutines",
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.WithPrometheusTarget(
					`go_goroutines{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}}"),
				),
				timeseries.Axis(
					axis.SoftMin(0),
				),
			),
		),
	}
}

func float64Ptr(input float64) *float64 {
	return &input
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, generalInfoRow(p)...)
	opts = append(opts, logPollerRow(p)...)
	opts = append(opts, feedJobsRow(p)...)
	opts = append(opts, mailBoxRow(p)...)
	opts = append(opts, promReporterRow(p)...)
	opts = append(opts, txManagerRow(p)...)
	opts = append(opts, headTrackerRow(p)...)
	opts = append(opts, appDBConnectionsRow(p)...)
	opts = append(opts, sqlQueriesRow(p)...)
	opts = append(opts, logsCountersRow(p)...)
	opts = append(opts, evmPoolLifecycleRow(p)...)
	opts = append(opts, nodesRPCRow(p)...)
	opts = append(opts, evmNodeRPCRow(p)...)
	opts = append(opts, evmRPCNodeLatenciesRow(p)...)
	opts = append(opts, evmBlockHistoryEstimatorRow(p)...)
	opts = append(opts, pipelinesRow(p)...)
	opts = append(opts, httpAPIRow(p)...)
	opts = append(opts, promHTTPRow(p)...)
	opts = append(opts, goMetricsRow(p)...)
	return opts
}
