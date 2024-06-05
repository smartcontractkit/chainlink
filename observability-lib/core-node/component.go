package corenode

import (
	"fmt"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"
)

func BuildDashboard(name string, dataSourceMetric string, platform string) (dashboard.Dashboard, error) {
	props := Props{
		MetricsDataSource: dataSourceMetric,
		PlatformOpts:      PlatformPanelOpts(platform),
	}

	builder := dashboard.NewDashboardBuilder(name).
		Tags([]string{"Core", "Node"}).
		Refresh("30s").
		Time("now-30m", "now")

	utils.AddVars(builder, vars(props))

	builder.WithRow(dashboard.NewRowBuilder("General CL Cluster Info"))
	utils.AddPanels(builder, panelsGeneralClusterInfo(props))

	builder.WithRow(dashboard.NewRowBuilder("LogPoller"))
	utils.AddPanels(builder, logPoller(props))

	builder.WithRow(dashboard.NewRowBuilder("Feeds Jobs"))
	utils.AddPanels(builder, feedsJobs(props))

	builder.WithRow(dashboard.NewRowBuilder("Mailbox"))
	utils.AddPanels(builder, mailbox(props))

	builder.WithRow(dashboard.NewRowBuilder("PromReporter"))
	utils.AddPanels(builder, promReporter(props))

	builder.WithRow(dashboard.NewRowBuilder("TxManager"))
	utils.AddPanels(builder, txManager(props))

	builder.WithRow(dashboard.NewRowBuilder("HeadTracker"))
	utils.AddPanels(builder, headTracker(props))

	builder.WithRow(dashboard.NewRowBuilder("AppDBConnections"))
	utils.AddPanels(builder, appDBConnections(props))

	builder.WithRow(dashboard.NewRowBuilder("SQLQueries"))
	utils.AddPanels(builder, sqlQueries(props))

	builder.WithRow(dashboard.NewRowBuilder("LogsCounters"))
	utils.AddPanels(builder, logsCounters(props))

	builder.WithRow(dashboard.NewRowBuilder("EvmPoolLifecycle"))
	utils.AddPanels(builder, evmPoolLifecycle(props))

	builder.WithRow(dashboard.NewRowBuilder("Node RPC State"))
	utils.AddPanels(builder, nodesRPC(props))

	builder.WithRow(dashboard.NewRowBuilder("EVM Pool RPC Node Metrics (App)"))
	utils.AddPanels(builder, evmNodeRPC(props))

	builder.WithRow(dashboard.NewRowBuilder("EVM Pool RPC Node Latencies (App)"))
	utils.AddPanels(builder, evmRPCNodeLatencies(props))

	builder.WithRow(dashboard.NewRowBuilder("Block History Estimator"))
	utils.AddPanels(builder, evmBlockHistoryEstimator(props))

	builder.WithRow(dashboard.NewRowBuilder("Pipeline Metrics (Runner)"))
	utils.AddPanels(builder, pipelines(props))

	builder.WithRow(dashboard.NewRowBuilder("HTTP API Metrics"))
	utils.AddPanels(builder, httpAPI(props))

	builder.WithRow(dashboard.NewRowBuilder("PromHTTP Metrics"))
	utils.AddPanels(builder, promHTTP(props))

	builder.WithRow(dashboard.NewRowBuilder("Go Metrics"))
	utils.AddPanels(builder, goMetrics(props))

	return builder.Build()
}

func vars(p Props) []cog.Builder[dashboard.VariableModel] {
	var variables []cog.Builder[dashboard.VariableModel]
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "instance", "Instance", fmt.Sprintf("label_values(%s)", p.PlatformOpts.LabelFilter), true))
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "evmChainID", "EvmChainID", fmt.Sprintf("label_values(%s)", "evmChainID"), true))

	return variables
}

func panelsGeneralClusterInfo(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"App Version",
		"app version",
		4,
		4,
		1,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `version{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{version}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Go Version",
		"golang version",
		4,
		4,
		1,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `go_info{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{version}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Uptime in days",
		"instance uptime",
		4,
		16,
		1,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `uptime_seconds{` + p.PlatformOpts.LabelQuery + `} / 86400`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"ETH Balance",
		"ETH balance",
		4,
		12,
		2,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `eth_balance{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{account}}`,
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "red"},
				{Value: utils.Float64Ptr(0.99), Color: "green"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Solana Balance",
		"Solana balance",
		4,
		12,
		2,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `solana_balance{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{account}}`,
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "red"},
				{Value: utils.Float64Ptr(0.99), Color: "green"},
			})),
	)

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Service Components Health",
		"service components health",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `health{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{service_id}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"ETH Balance",
		"eth balance graph",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `eth_balance{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{account}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"SOL Balance",
		"sol balance graph",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `solana_balance{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{account}}`,
		},
	))

	return panelsArray
}

func logPoller(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Goroutines",
		"goroutines",
		6,
		12,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `count(count by (evmChainID) (log_poller_query_duration_sum{job=~"$instance"}))`,
			Legend: "Goroutines",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"RPS",
		"requests per second",
		6,
		12,
		2,
		"reqps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (query) (sum by (query, job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
			Legend: "{{query}} - {{job}}",
		},
		utils.PrometheusQuery{
			Query:  `avg (sum by(job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
			Legend: "Total",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"RPS by Type",
		"",
		6,
		12,
		2,
		"reqps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (type) (sum by (type, job) (rate(log_poller_query_duration_count{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])))`,
			Legend: "{{query}} - {{job}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Avg number of logs returned",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (query) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
			Legend: "{{query}} - {{job}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Max number of logs returned",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `max by (query) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
			Legend: "{{query}} - {{job}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Logs returned by chain",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `max by (evmChainID) (log_poller_query_dataset_size{job=~"$instance", evmChainID=~"$evmChainID"})`,
			Legend: "{{evmChainID}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Queries duration by type (0.5 perc)",
		"",
		6,
		12,
		2,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.5, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
			Legend: "{{query}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Queries duration by type (0.9 perc)",
		"",
		6,
		12,
		2,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.9, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
			Legend: "{{query}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Queries duration by type (0.99 perc)",
		"",
		6,
		12,
		2,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.99, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, query)) / 1e6`,
			Legend: "{{query}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Queries duration by chain (0.99 perc)",
		"",
		6,
		12,
		2,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.99, sum(rate(log_poller_query_duration_bucket{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval])) by (le, evmChainID)) / 1e6`,
			Legend: "{{query}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Number of logs inserted",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (evmChainID) (log_poller_logs_inserted{job=~"$instance", evmChainID=~"$evmChainID"})`,
			Legend: "{{evmChainID}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Logs insertion rate",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (evmChainID) (rate(log_poller_logs_inserted{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval]))`,
			Legend: "{{evmChainID}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Number of blocks inserted",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (evmChainID) (log_poller_blocks_inserted{job=~"$instance", evmChainID=~"$evmChainID"})`,
			Legend: "{{evmChainID}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Blocks insertion rate",
		"",
		6,
		12,
		2,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg by (evmChainID) (rate(log_poller_blocks_inserted{job=~"$instance", evmChainID=~"$evmChainID"}[$__rate_interval]))`,
			Legend: "{{evmChainID}}",
		},
	))

	return panelsArray
}

func feedsJobs(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Feeds Job Proposal Requests",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `feeds_job_proposal_requests{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Feeds Job Proposal Count",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `feeds_job_proposal_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func mailbox(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Mailbox Load Percent",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `mailbox_load_percent{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{name}}`,
		},
	))

	return panelsArray
}

func promReporter(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Unconfirmed Transactions",
		"",
		6,
		8,
		1,
		"Tx",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `unconfirmed_transactions{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Unconfirmed TX Age",
		"",
		6,
		8,
		1,
		"s",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `max_unconfirmed_tx_age{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Unconfirmed TX Blocks",
		"",
		6,
		8,
		1,
		"Blocks",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `max_unconfirmed_blocks{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func txManager(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Time Until TX Broadcast",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_time_until_tx_broadcast{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Gas Bumps",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_num_gas_bumps{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Gas Bumps Exceeds Limit",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_gas_bump_exceeds_limit{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Confirmed Transactions",
		"",
		6,
		6,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_num_confirmed_transactions{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Successful Transactions",
		"",
		6,
		6,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_num_successful_transactions{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Reverted Transactions",
		"",
		6,
		6,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_num_tx_reverted{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Fwd Transactions",
		"",
		6,
		6,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_fwd_tx_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Num Transactions Attempts",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_tx_attempt_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Time Until TX Confirmed",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_time_until_tx_confirmed{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"TX Manager Block Until TX Confirmed",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `tx_manager_blocks_until_tx_confirmed{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func headTracker(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Head tracker current head",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `head_tracker_current_head{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Head tracker very old head",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `head_tracker_very_old_head{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Head tracker heads received",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `head_tracker_heads_received{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Head tracker connection errors",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `head_tracker_connection_errors{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func appDBConnections(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"DB Connections",
		"",
		6,
		24,
		1,
		"Conn",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `db_conns_max{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Max`,
		},
		utils.PrometheusQuery{
			Query:  `db_conns_open{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Open`,
		},
		utils.PrometheusQuery{
			Query:  `db_conns_used{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Used`,
		},
		utils.PrometheusQuery{
			Query:  `db_conns_wait{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Wait`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"DB Wait Count",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `db_wait_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"DB Wait Time",
		"",
		6,
		12,
		1,
		"Sec",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `db_wait_time_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func sqlQueries(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"SQL Query Timeout Percent",
		"",
		6,
		24,
		1,
		"percent",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.9, sum(rate(sql_query_timeout_percent_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (le))`,
			Legend: "p90",
		},
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.95, sum(rate(sql_query_timeout_percent_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (le))`,
			Legend: "p95",
		},
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.99, sum(rate(sql_query_timeout_percent_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (le))`,
			Legend: "p99",
		},
	))

	return panelsArray
}

func logsCounters(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Logs Counters",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `log_panic_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - panic`,
		},
		utils.PrometheusQuery{
			Query:  `log_fatal_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - fatal`,
		},
		utils.PrometheusQuery{
			Query:  `log_critical_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - critical`,
		},
		utils.PrometheusQuery{
			Query:  `log_warn_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - warn`,
		},
		utils.PrometheusQuery{
			Query:  `log_error_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - error`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Logs Rate",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(rate(log_panic_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - panic`,
		},
		utils.PrometheusQuery{
			Query:  `sum(rate(log_fatal_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - fatal`,
		},
		utils.PrometheusQuery{
			Query:  `sum(rate(log_critical_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - critical`,
		},
		utils.PrometheusQuery{
			Query:  `sum(rate(log_warn_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - warn`,
		},
		utils.PrometheusQuery{
			Query:  `sum(rate(log_error_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - error`,
		},
	))

	return panelsArray
}

func evmPoolLifecycle(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool Highest Seen Block",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_highest_seen_block{` + p.PlatformOpts.LabelQuery + `evmChainID="${evmChainID}"}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool Num Seen Blocks",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_seen_blocks{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool Node Polls Total",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_polls_total{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool Node Polls Failed",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_polls_failed{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool Node Polls Success",
		"",
		6,
		12,
		1,
		"Block",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_polls_success{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func nodesRPC(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC Alive",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Alive"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC Closed",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Closed"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC Dialed",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Dialed"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC InvalidChainID",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="InvalidChainID"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC OutOfSync",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="OutOfSync"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC UnDialed",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Undialed"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC Unreachable",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Unreachable"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Node RPC Unusable",
		"",
		6,
		6,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(multi_node_states{` + p.PlatformOpts.LabelQuery + `chainId=~"$evmChainID", state="Unusable"}) by (` + p.PlatformOpts.LegendString + `, chainId)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{chainId}}`,
		},
	))

	return panelsArray
}

func evmNodeRPC(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Calls Success Rate",
		"",
		6,
		24,
		1,
		"percentunit",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(increase(evm_pool_rpc_node_calls_success{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_calls_total{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{evmChainID}} - {{nodeName}}`,
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "red"},
				{Value: utils.Float64Ptr(0.8), Color: "orange"},
				{Value: utils.Float64Ptr(0.99), Color: "green"},
			})),
	)

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Dials Failure Rate",
		"",
		6,
		24,
		1,
		"percentunit",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(increase(evm_pool_rpc_node_dials_failed{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_calls_total{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{evmChainID}} - {{nodeName}}`,
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.3), Color: "orange"},
				{Value: utils.Float64Ptr(0.7), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Transitions",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_alive{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "Alive",
		},
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_in_sync{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "InSync",
		},
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_out_of_sync{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "OutOfSync",
		},
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_unreachable{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "UnReachable",
		},
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_invalid_chain_id{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "InvalidChainID",
		},
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_num_transitions_to_unusable{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "TransitionToUnusable",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node States",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `evm_pool_rpc_node_states{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{evmChainID}} - {{state}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Verifies Success Rate",
		"",
		6,
		12,
		1,
		"percentunit",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(increase(evm_pool_rpc_node_verifies_success{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_verifies{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) * 100`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{evmChainID}} - {{nodeName}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Verifies Failure Rate",
		"",
		6,
		12,
		1,
		"percentunit",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(increase(evm_pool_rpc_node_verifies_failed{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) / sum(increase(evm_pool_rpc_node_verifies{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, evmChainID, nodeName) * 100`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{evmChainID}} - {{nodeName}}`,
		},
	))

	return panelsArray
}

func evmRPCNodeLatencies(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Calls Latency 0.90 quantile",
		"",
		6,
		24,
		1,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.90, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, le, rpcCallName)) / 1e6`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{rpcCallName}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Calls Latency 0.95 quantile",
		"",
		6,
		24,
		1,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.95, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, le, rpcCallName)) / 1e6`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{rpcCallName}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"EVM Pool RPC Node Calls Latency 0.99 quantile",
		"",
		6,
		24,
		1,
		"ms",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.99, sum(rate(evm_pool_rpc_node_rpc_call_time_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, le, rpcCallName)) / 1e6`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{rpcCallName}}`,
		},
	))

	return panelsArray
}

func evmBlockHistoryEstimator(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Gas Updater All Gas Price Percentiles",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `gas_updater_all_gas_price_percentiles{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{percentile}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Gas Updater All Tip Cap Percentiles",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `gas_updater_all_tip_cap_percentiles{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{percentile}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Gas Updater Set Gas Price",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `gas_updater_set_gas_price{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Gas Updater Set Tip Cap",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `gas_updater_set_tip_cap{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Gas Updater Current Base Fee",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `gas_updater_current_base_fee{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Block History Estimator Connectivity Failure Count",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `block_history_estimator_connectivity_failure_count{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func pipelines(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Pipeline Task Execution Time",
		"",
		6,
		24,
		1,
		"s",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `pipeline_task_execution_time{` + p.PlatformOpts.LabelQuery + `} / 1e6`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} JobID: {{job_id}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Pipeline Run Errors",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `pipeline_run_errors{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} JobID: {{job_id}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Pipeline Run Total Time to Completion",
		"",
		6,
		24,
		1,
		"s",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `pipeline_run_total_time_to_completion{` + p.PlatformOpts.LabelQuery + `} / 1e6`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} JobID: {{job_id}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Pipeline Tasks Total Finished",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `pipeline_tasks_total_finished{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} JobID: {{job_id}}`,
		},
	))

	return panelsArray
}

func httpAPI(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Request Duration p95",
		"",
		6,
		24,
		1,
		"s",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `histogram_quantile(0.95, sum(rate(service_gonic_request_duration_bucket{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, le, path, method))`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{method}} - {{path}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Request Total Rate over interval",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(rate(service_gonic_requests_total{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, path, method, code)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{method}} - {{path}} - {{code}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Average Request Size",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg(rate(service_gonic_request_size_bytes_sum{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)/avg(rate(service_gonic_request_size_bytes_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Response Size",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg(rate(service_gonic_response_size_bytes_sum{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)/avg(rate(service_gonic_response_size_bytes_count{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}

func promHTTP(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.GaugePanel(
		p.MetricsDataSource,
		"HTTP Request in flight",
		"",
		6,
		24,
		1,
		"",
		utils.PrometheusQuery{
			Query:  `promhttp_metric_handler_requests_in_flight{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"HTTP rate by return code",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(rate(promhttp_metric_handler_requests_total{` + p.PlatformOpts.LabelQuery + `}[$__rate_interval])) by (` + p.PlatformOpts.LegendString + `, code)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - {{code}}`,
		},
	))

	return panelsArray
}

func goMetrics(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Threads",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(go_threads{` + p.PlatformOpts.LabelQuery + `}) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Heap Allocations",
		"",
		4,
		24,
		1,
		"bytes",
		common.BigValueColorModeNone,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValue,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(go_memstats_heap_alloc_bytes{` + p.PlatformOpts.LabelQuery + `}) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: "",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Heap allocations",
		"",
		6,
		24,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(go_memstats_heap_alloc_bytes{` + p.PlatformOpts.LabelQuery + `}) by (` + p.PlatformOpts.LegendString + `)`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Heap allocations",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_memstats_heap_alloc_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Alloc`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_heap_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Sys`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_heap_idle_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Idle`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_heap_inuse_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - InUse`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_heap_released_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Released`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Memory in Off-Heap",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_memstats_mspan_inuse_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Total InUse`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_mspan_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Total Sys`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_mcache_inuse_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Cache InUse`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_mcache_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Cache Sys`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_buck_hash_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Hash Sys`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_gc_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - GC Sys`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_other_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - bytes of memory are used for other runtime allocations`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_next_gc_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Next GC`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Memory in Stack",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_memstats_stack_inuse_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - InUse`,
		},
		utils.PrometheusQuery{
			Query:  `go_memstats_stack_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}} - Sys`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Total Used Memory",
		"",
		6,
		12,
		1,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_memstats_sys_bytes{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Number of Live Objects",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_memstats_mallocs_total{` + p.PlatformOpts.LabelQuery + `} - go_memstats_frees_total{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Rate of Objects Allocated",
		"",
		6,
		12,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `rate(go_memstats_mallocs_total{` + p.PlatformOpts.LabelQuery + `}[1m])`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Rate of a Pointer Dereferences",
		"",
		6,
		12,
		1,
		"ops",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `rate(go_memstats_lookups_total{` + p.PlatformOpts.LabelQuery + `}[1m])`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Goroutines",
		"",
		6,
		12,
		1,
		"ops",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `go_goroutines{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + p.PlatformOpts.LegendString + `}}`,
		},
	))

	return panelsArray
}
