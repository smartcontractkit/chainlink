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
	"github.com/K-Phoen/grabana/variable/query"
	"github.com/pkg/errors"
)

/*
Use ripgrep to get the full list
rg -oU ".*promauto.*\n.*Name: \"(.*)\"" -r '$1' > metrics.txt

duplicates?

common/client/node.go:pool_rpc_node_verifies
common/client/node.go:pool_rpc_node_verifies_failed
common/client/node.go:pool_rpc_node_verifies_success
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_alive
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_in_sync
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_out_of_sync
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_unreachable
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_invalid_chain_id
common/client/node_fsm.go:pool_rpc_node_num_transitions_to_unusable
common/client/node_lifecycle.go:pool_rpc_node_highest_seen_block
common/client/node_lifecycle.go:pool_rpc_node_num_seen_blocks
common/client/node_lifecycle.go:pool_rpc_node_polls_total
common/client/node_lifecycle.go:pool_rpc_node_polls_failed
common/client/node_lifecycle.go:pool_rpc_node_polls_success

covered

core/logger/prometheus.go:log_warn_count
core/logger/prometheus.go:log_error_count
core/logger/prometheus.go:log_critical_count
core/logger/prometheus.go:log_panic_count
core/logger/prometheus.go:log_fatal_count
common/client/multi_node.go:multi_node_states
common/txmgr/broadcaster.go:tx_manager_time_until_tx_broadcast
common/txmgr/confirmer.go:tx_manager_num_gas_bumps
common/txmgr/confirmer.go:tx_manager_gas_bump_exceeds_limit
common/txmgr/confirmer.go:tx_manager_num_confirmed_transactions
common/txmgr/confirmer.go:tx_manager_num_successful_transactions
common/txmgr/confirmer.go:tx_manager_num_tx_reverted
common/txmgr/confirmer.go:tx_manager_fwd_tx_count
common/txmgr/confirmer.go:tx_manager_tx_attempt_count
common/txmgr/confirmer.go:tx_manager_time_until_tx_confirmed
common/txmgr/confirmer.go:tx_manager_blocks_until_tx_confirmed
common/headtracker/head_tracker.go:head_tracker_current_head
common/headtracker/head_tracker.go:head_tracker_very_old_head
common/headtracker/head_listener.go:head_tracker_heads_received
common/headtracker/head_listener.go:head_tracker_connection_errors
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_alive
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_in_sync
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_out_of_sync
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_unreachable
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_invalid_chain_id
core/chains/evm/client/node_fsm.go:evm_pool_rpc_node_num_transitions_to_unusable
core/services/promreporter/prom_reporter.go:unconfirmed_transactions
core/services/promreporter/prom_reporter.go:max_unconfirmed_tx_age
core/services/promreporter/prom_reporter.go:max_unconfirmed_blocks
core/services/promreporter/prom_reporter.go:pipeline_runs_queued
core/services/promreporter/prom_reporter.go:pipeline_task_runs_queued
core/services/pipeline/task.bridge.go:bridge_latency_seconds
core/services/pipeline/task.bridge.go:bridge_errors_total
core/services/pipeline/task.bridge.go:bridge_cache_hits_total
core/services/pipeline/task.bridge.go:bridge_cache_errors_total
core/services/pipeline/task.http.go:pipeline_task_http_fetch_time
core/services/pipeline/task.http.go:pipeline_task_http_response_body_size
core/services/pipeline/task.eth_call.go:pipeline_task_eth_call_execution_time
core/services/pipeline/runner.go:pipeline_task_execution_time
core/services/pipeline/runner.go:pipeline_run_errors
core/services/pipeline/runner.go:pipeline_run_total_time_to_completion
core/services/pipeline/runner.go:pipeline_tasks_total_finished
core/chains/evm/client/node.go:evm_pool_rpc_node_dials_total
core/chains/evm/client/node.go:evm_pool_rpc_node_dials_failed
core/chains/evm/client/node.go:evm_pool_rpc_node_dials_success
core/chains/evm/client/node.go:evm_pool_rpc_node_verifies
core/chains/evm/client/node.go:evm_pool_rpc_node_verifies_failed
core/chains/evm/client/node.go:evm_pool_rpc_node_verifies_success
core/chains/evm/client/node.go:evm_pool_rpc_node_calls_total
core/chains/evm/client/node.go:evm_pool_rpc_node_calls_failed
core/chains/evm/client/node.go:evm_pool_rpc_node_calls_success
core/chains/evm/client/node.go:evm_pool_rpc_node_rpc_call_time
core/chains/evm/client/pool.go:evm_pool_rpc_node_states
core/utils/mailbox_prom.go:mailbox_load_percent
core/services/pg/stats.go:db_conns_max
core/services/pg/stats.go:db_conns_open
core/services/pg/stats.go:db_conns_used
core/services/pg/stats.go:db_wait_count
core/services/pg/stats.go:db_wait_time_seconds
core/chains/evm/client/node_lifecycle.go:evm_pool_rpc_node_highest_seen_block
core/chains/evm/client/node_lifecycle.go:evm_pool_rpc_node_num_seen_blocks
core/chains/evm/client/node_lifecycle.go:evm_pool_rpc_node_polls_total
core/chains/evm/client/node_lifecycle.go:evm_pool_rpc_node_polls_failed
core/chains/evm/client/node_lifecycle.go:evm_pool_rpc_node_polls_success
core/services/relay/evm/config_poller.go:ocr2_failed_rpc_contract_calls
core/services/feeds/service.go:feeds_job_proposal_requests
core/services/feeds/service.go:feeds_job_proposal_count
core/services/ocrcommon/prom.go:bridge_json_parse_values
core/services/ocrcommon/prom.go:ocr_median_values
core/chains/evm/logpoller/observability.go:log_poller_query_dataset_size

not-covered and product specific (definitions/usage should be moved to plugins)

mercury

core/services/relay/evm/mercury/types/types.go:mercury_price_feed_missing
core/services/relay/evm/mercury/types/types.go:mercury_price_feed_errors
core/services/relay/evm/mercury/queue.go:mercury_transmit_queue_load
core/services/relay/evm/mercury/v1/data_source.go:mercury_insufficient_blocks_count
core/services/relay/evm/mercury/v1/data_source.go:mercury_zero_blocks_count
core/services/relay/evm/mercury/wsrpc/client.go:mercury_transmit_timeout_count
core/services/relay/evm/mercury/wsrpc/client.go:mercury_dial_count
core/services/relay/evm/mercury/wsrpc/client.go:mercury_dial_success_count
core/services/relay/evm/mercury/wsrpc/client.go:mercury_dial_error_count
core/services/relay/evm/mercury/wsrpc/client.go:mercury_connection_reset_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_success_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_duplicate_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_connection_error_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_queue_delete_error_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_queue_insert_error_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_queue_push_error_count
core/services/relay/evm/mercury/transmitter.go:mercury_transmit_server_error_count

functions

core/services/gateway/connectionmanager.go:gateway_heartbeats_sent
core/services/gateway/gateway.go:gateway_request
core/services/gateway/handlers/functions/handler.functions.go:gateway_functions_handler_error
core/services/gateway/handlers/functions/handler.functions.go:gateway_functions_secrets_set_success
core/services/gateway/handlers/functions/handler.functions.go:gateway_functions_secrets_set_failure
core/services/gateway/handlers/functions/handler.functions.go:gateway_functions_secrets_list_success
core/services/gateway/handlers/functions/handler.functions.go:gateway_functions_secrets_list_failure
core/services/functions/external_adapter_client.go:functions_external_adapter_client_latency
core/services/functions/external_adapter_client.go:functions_external_adapter_client_errors_total
core/services/functions/listener.go:functions_request_received
core/services/functions/listener.go:functions_request_internal_error
core/services/functions/listener.go:functions_request_computation_error
core/services/functions/listener.go:functions_request_computation_success
core/services/functions/listener.go:functions_request_timeout
core/services/functions/listener.go:functions_request_confirmed
core/services/functions/listener.go:functions_request_computation_result_size
core/services/functions/listener.go:functions_request_computation_error_size
core/services/functions/listener.go:functions_request_computation_duration
core/services/functions/listener.go:functions_request_pruned
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_restarts
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_query
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_observation
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_report
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_report_num_observations
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_accept
core/services/ocr2/plugins/functions/reporting.go:functions_reporting_plugin_transmit
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_query
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_observation
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_report
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_accept
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_query_byte_size
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_query_rows_count
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_observation_rows_count
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_report_rows_count
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_wrong_sig_count
core/services/ocr2/plugins/s4/prometheus.go:s4_reporting_plugin_expired_rows

vrf

core/services/vrf/vrfcommon/metrics.go:vrf_request_queue_size
core/services/vrf/vrfcommon/metrics.go:vrf_processed_request_count
core/services/vrf/vrfcommon/metrics.go:vrf_dropped_request_count
core/services/vrf/vrfcommon/metrics.go:vrf_duplicate_requests
core/services/vrf/vrfcommon/metrics.go:vrf_request_time_between_sims
core/services/vrf/vrfcommon/metrics.go:vrf_request_time_until_initial_sim

keeper
core/services/keeper/upkeep_executer.go:keeper_check_upkeep_execution_time
*/

const (
	ErrFailedToCreateDashboard = "failed to create dashboard"
	ErrFailedToCreateFolder    = "failed to create folder"
)

// CLClusterDashboard is a dashboard for a Chainlink cluster
type CLClusterDashboard struct {
	Nodes                    int
	Name                     string
	LokiDataSourceName       string
	PrometheusDataSourceName string
	Folder                   string
	GrafanaURL               string
	GrafanaToken             string
	opts                     []dashboard.Option
	extendedOpts             []dashboard.Option
	builder                  dashboard.Builder
}

// NewCLClusterDashboard returns a new dashboard for a Chainlink cluster, can be used as a base for more complex plugin based dashboards
func NewCLClusterDashboard(nodes int, name, ldsn, pdsn, dbf, grafanaURL, grafanaToken string, opts []dashboard.Option) (*CLClusterDashboard, error) {
	db := &CLClusterDashboard{
		Nodes:                    nodes,
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

func (m *CLClusterDashboard) Opts() []dashboard.Option {
	return m.opts
}

// logsRowOption returns a row option for a node's logs with name and instance selector
func (m *CLClusterDashboard) logsRowOption(name, q string) row.Option {
	return row.WithLogs(
		name,
		logs.DataSource(m.LokiDataSourceName),
		logs.Span(12),
		logs.Height("300px"),
		logs.Transparent(),
		logs.WithLokiTarget(q),
	)
}

func (m *CLClusterDashboard) logsRowOptionsForNodes(nodes int) []row.Option {
	opts := make([]row.Option, 0)
	for i := 1; i <= nodes; i++ {
		opts = append(opts, row.WithLogs(
			fmt.Sprintf("Node %d", i),
			logs.DataSource(m.LokiDataSourceName),
			logs.Span(12),
			logs.Height("300px"),
			logs.Transparent(),
			logs.WithLokiTarget(fmt.Sprintf(`{namespace="${namespace}", app="app", instance="node-%d", container="node"}`, i)),
		))
	}
	return opts
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
		dashboard.Row(
			"Cluster health",
			row.Collapse(),
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
		// HeadTracker
		dashboard.Row("Head tracker",
			row.Collapse(),
			m.timeseriesRowOption(
				"Head tracker current head",
				"Block",
				`head_tracker_current_head{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Head tracker very old head",
				"Block",
				`head_tracker_very_old_head{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Head tracker heads received",
				"Block",
				`head_tracker_heads_received{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Head tracker connection errors",
				"Errors",
				`head_tracker_connection_errors{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("LogPoller",
			row.Collapse(),
			m.timeseriesRowOption(
				"LogPoller Query Dataset Size",
				"",
				`log_poller_query_dataset_size{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("OCRCommon",
			row.Collapse(),
			m.timeseriesRowOption(
				"Bridge JSON Parse Values",
				"",
				`bridge_json_parse_values{namespace="${namespace}"}`,
				"{{ pod }} JobID: {{ job_id }}",
			),
			m.timeseriesRowOption(
				"OCR Median Values",
				"",
				`ocr_median_values{namespace="${namespace}"}`,
				"{{pod}} JobID: {{ job_id }}",
			),
		),
		dashboard.Row("Relay Config Poller",
			row.Collapse(),
			m.timeseriesRowOption(
				"Relay Config Poller RPC Contract Calls",
				"",
				`ocr2_failed_rpc_contract_calls{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("Feeds Jobs",
			row.Collapse(),
			m.timeseriesRowOption(
				"Feeds Job Proposal Requests",
				"",
				`feeds_job_proposal_requests{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Feeds Job Proposal Count",
				"",
				`feeds_job_proposal_count{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("Mailbox",
			row.Collapse(),
			m.timeseriesRowOption(
				"Mailbox Load Percent",
				"",
				`mailbox_load_percent{namespace="${namespace}"}`,
				"{{ pod }} {{ name }}",
			),
		),
		dashboard.Row("Multi Node States",
			row.Collapse(),
			m.timeseriesRowOption(
				"Multi Node States",
				"",
				`multi_node_states{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("Block History Estimator",
			row.Collapse(),
			m.timeseriesRowOption(
				"Gas Updater All Gas Price Percentiles",
				"",
				`gas_updater_all_gas_price_percentiles{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Gas Updater All Tip Cap Percentiles",
				"",
				`gas_updater_all_tip_cap_percentiles{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Gas Updater Set Gas Price",
				"",
				`gas_updater_set_gas_price{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Gas Updater Set Tip Cap",
				"",
				`gas_updater_set_tip_cap{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Gas Updater Current Base Fee",
				"",
				`gas_updater_current_base_fee{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Block History Estimator Connectivity Failure Count",
				"",
				`block_history_estimator_connectivity_failure_count{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		// PromReporter
		dashboard.Row("Prom Reporter",
			row.Collapse(),
			m.timeseriesRowOption(
				"Unconfirmed Transactions",
				"Tx",
				`unconfirmed_transactions{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Unconfirmed TX Age",
				"Sec",
				`max_unconfirmed_tx_age{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"Unconfirmed TX Blocks",
				"Blocks",
				`max_unconfirmed_blocks{namespace="${namespace}"}`,
				"{{ pod }}",
			),
		),
		dashboard.Row("TX Manager",
			row.Collapse(),
			m.timeseriesRowOption(
				"TX Manager Time Until TX Broadcast",
				"",
				`tx_manager_time_until_tx_broadcast{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Gas Bumps",
				"",
				`tx_manager_num_gas_bumps{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Gas Bumps Exceeds Limit",
				"",
				`tx_manager_gas_bump_exceeds_limit{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Confirmed Transactions",
				"",
				`tx_manager_num_confirmed_transactions{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Successful Transactions",
				"",
				`tx_manager_num_successful_transactions{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Reverted Transactions",
				"",
				`tx_manager_num_tx_reverted{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Fwd Transactions",
				"",
				`tx_manager_fwd_tx_count{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Num Transactions Attempts",
				"",
				`tx_manager_tx_attempt_count{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Time Until TX Confirmed",
				"",
				`tx_manager_time_until_tx_confirmed{namespace="${namespace}"}`,
				"{{ pod }}",
			),
			m.timeseriesRowOption(
				"TX Manager Block Until TX Confirmed",
				"",
				`tx_manager_blocks_until_tx_confirmed{namespace="${namespace}"}`,
				"{{ pod }}",
			),
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
			"EVM Pool Lifecycle",
			row.Collapse(),
			m.timeseriesRowOption(
				"EVM Pool Highest Seen Block",
				"Block",
				`evm_pool_rpc_node_highest_seen_block{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool Num Seen Blocks",
				"Block",
				`evm_pool_rpc_node_num_seen_blocks{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool Node Polls Total",
				"Block",
				`evm_pool_rpc_node_polls_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool Node Polls Failed",
				"Block",
				`evm_pool_rpc_node_polls_failed{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool Node Polls Success",
				"Block",
				`evm_pool_rpc_node_polls_success{namespace="${namespace}"}`,
				"{{pod}}",
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
				"DB Wait Count",
				"",
				`db_wait_count{namespace="${namespace}"}`,
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
				"EVM Pool RPC Node Dials Failed",
				"",
				`evm_pool_rpc_node_dials_failed{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Dials Total",
				"",
				`evm_pool_rpc_node_dials_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Node Dials Failed",
				"",
				`evm_pool_rpc_node_dials_failed{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to Alive",
				"",
				`evm_pool_rpc_node_num_transitions_to_alive{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to In Sync",
				"",
				`evm_pool_rpc_node_num_transitions_to_in_sync{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to Out of Sync",
				"",
				`evm_pool_rpc_node_num_transitions_to_out_of_sync{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to Unreachable",
				"",
				`evm_pool_rpc_node_num_transitions_to_unreachable{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to invalid Chain ID",
				"",
				`evm_pool_rpc_node_num_transitions_to_invalid_chain_id{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"EVM Pool RPC Total Transitions to unusable",
				"",
				`evm_pool_rpc_node_num_transitions_to_unusable{namespace="${namespace}"}`,
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
			m.timeseriesRowOption(
				"EVM Pool RPC Node Verifies Failed",
				"",
				`evm_pool_rpc_node_verifies_failed{namespace="${namespace}"}`,
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
			"Pipeline Metrics (Runner)",
			row.Collapse(),
			m.timeseriesRowOption(
				"Pipeline Task Execution Time",
				"Sec",
				`pipeline_task_execution_time{namespace="${namespace}"} / 1e6`,
				"{{ pod }} JobID: {{ job_id }}",
			),
			m.timeseriesRowOption(
				"Pipeline Run Errors",
				"",
				`pipeline_run_errors{namespace="${namespace}"}`,
				"{{ pod }} JobID: {{ job_id }}",
			),
			m.timeseriesRowOption(
				"Pipeline Run Total Time to Completion",
				"Sec",
				`pipeline_run_total_time_to_completion{namespace="${namespace}"} / 1e6`,
				"{{ pod }} JobID: {{ job_id }}",
			),
			m.timeseriesRowOption(
				"Pipeline Tasks Total Finished",
				"",
				`pipeline_tasks_total_finished{namespace="${namespace}"}`,
				"{{ pod }} JobID: {{ job_id }}",
			),
		),
		dashboard.Row(
			"Pipeline Metrics (ETHCall)",
			row.Collapse(),
			m.timeseriesRowOption(
				"Pipeline Task ETH Call Execution Time",
				"Sec",
				`pipeline_task_eth_call_execution_time{namespace="${namespace}"}`,
				"{{pod}}",
			),
		),
		dashboard.Row(
			"Pipeline Metrics (HTTP)",
			row.Collapse(),
			m.timeseriesRowOption(
				"Pipeline Task HTTP Fetch Time",
				"Sec",
				`pipeline_task_http_fetch_time{namespace="${namespace}"} / 1e6`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"Pipeline Task HTTP Response Body Size",
				"Bytes",
				`pipeline_task_http_response_body_size{namespace="${namespace}"}`,
				"{{pod}}",
			),
		),
		dashboard.Row(
			"Pipeline Metrics (Bridge)",
			row.Collapse(),
			m.timeseriesRowOption(
				"Pipeline Bridge Latency",
				"Sec",
				`bridge_latency_seconds{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"Pipeline Bridge Errors Total",
				"",
				`bridge_errors_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"Pipeline Bridge Cache Hits Total",
				"",
				`bridge_cache_hits_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
			m.timeseriesRowOption(
				"Pipeline Bridge Cache Errors Total",
				"",
				`bridge_cache_errors_total{namespace="${namespace}"}`,
				"{{pod}}",
			),
		),
		dashboard.Row(
			"Pipeline Metrics",
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
	logOptsFinal := make([]row.Option, 0)
	logOptsFinal = append(
		logOptsFinal,
		row.Collapse(),
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
				`log_warn_count{namespace="${namespace}"}`,
				prometheus.Legend("{{pod}} - warn"),
			),
			timeseries.WithPrometheusTarget(
				`log_error_count{namespace="${namespace}"}`,
				prometheus.Legend("{{pod}} - error"),
			),
		),
		m.logsRowOption("All errors", `
			{namespace="${namespace}", app="app", container="node"} 
			| json 
			| level="error" 
			|  line_format "{{ .instance }} {{ .level }} {{ .ts }} {{ .logger }} {{ .caller }} {{ .msg }} {{ .version }} {{ .nodeTier }} {{ .nodeName }} {{ .node }} {{ .evmChainID }} {{ .nodeOrder }} {{ .mode }} {{ .nodeState }} {{ .sentryEventID }} {{ .stacktrace }}"`),
	)
	logOptsFinal = append(logOptsFinal, m.logsRowOptionsForNodes(m.Nodes)...)
	logRowOpts := dashboard.Row(
		"Logs",
		logOptsFinal...,
	)
	opts = append(opts, logRowOpts)
	opts = append(opts, m.extendedOpts...)
	builder, err := dashboard.New(
		"Chainlink Cluster Dashboard",
		opts...,
	)
	m.opts = opts
	m.builder = builder
	return err
}

// Deploy deploys the dashboard to Grafana
func (m *CLClusterDashboard) Deploy(ctx context.Context) error {
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
