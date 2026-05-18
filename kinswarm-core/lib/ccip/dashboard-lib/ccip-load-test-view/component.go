package ccip_load_test_view

import (
	"encoding/json"
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/loki"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/query"
	cLoki "github.com/grafana/grafana-foundation-sdk/go/loki"
	cXYChart "github.com/grafana/grafana-foundation-sdk/go/xychart"
)

type Props struct {
	LokiDataSource string
}

func vars(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.VariableAsQuery(
			"Test Run Name",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(namespace)"),
		),
		dashboard.VariableAsQuery(
			"cluster",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(cluster)"),
		),
		dashboard.VariableAsQuery(
			"test_group",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(test_group)"),
		),
		dashboard.VariableAsQuery(
			"test_id",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(test_id)"),
		),
		dashboard.VariableAsQuery(
			"source_chain",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(source_chain)"),
		),
		dashboard.VariableAsQuery(
			"dest_chain",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(dest_chain)"),
		),
		dashboard.VariableAsQuery(
			"geth_node",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(geth_node)"),
		),
		dashboard.VariableAsQuery(
			"remote_runner",
			query.DataSource(p.LokiDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("namespace"),
		),
	}
}

func XYChartSeqNum() map[string]interface{} {
	// TODO: https://github.com/grafana/grafana-foundation-sdk/tree/v10.4.x%2Bcog-v0.0.x/go has a lot of useful components
	// TODO: need to change upload API and use combined upload in lib/dashboard.go
	xAxisName := "seq_num"
	builder := cXYChart.NewPanelBuilder().
		Title("XYChart").
		Dims(cXYChart.XYDimensionConfig{
			X: &xAxisName,
		}).
		WithTarget(
			cLoki.NewDataqueryBuilder().
				Expr(`{namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_duration!= "" | data_Commit_ReportAccepted_success="✅"`).
				LegendFormat("Commit Report Accepted"),
		)
	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(dashboardJson, &data); err != nil {
		panic(err)
	}
	fmt.Println(string(dashboardJson))
	return data
}

func statsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"CCIP Duration Stats",
			row.Collapse(),
			row.WithTimeSeries(
				"Sequence numbers",
				timeseries.Transparent(),
				timeseries.Description("Sequence Numbers triggered by Test"),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", source_chain="${source_chain}", dest_chain="${dest_chain}"} | json | data_CCIPSendRequested_success="✅" | unwrap  data_CCIPSendRequested_seq_num [$__range]) by (test_id)`,
					loki.Legend("Starts"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"} | json | data_CCIPSendRequested_success="✅" | unwrap  data_CCIPSendRequested_seq_num [$__range]) by (test_id)`,
					loki.Legend("Ends"),
				),
			),
			row.WithTimeSeries(
				"Source Router Fees ( /1e18)",
				timeseries.Transparent(),
				timeseries.Description("Router.GetFee"),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_fee [$__range]) by (test_id) /1e18`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_fee [$__range]) by (test_id) /1e18`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_fee [$__range]) by (test_id) /1e18 `,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Commit Duration Summary",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("seconds"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_duration [$__range]) by (data_Commit_ReportAccepted_seqNum)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json  | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_duration [$__range]) by (data_Commit_ReportAccepted_seqNum)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_duration [$__range]) by (data_Commit_ReportAccepted_seqNum)`,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Report Blessing Summary",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("seconds"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_duration [$__range]) by (data_ReportBlessedByARM_seqNum)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({ namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json |  data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_duration [$__range]) by (data_ReportBlessedByARM_seqNum)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_duration [$__range]) by (data_ReportBlessedByARM_seqNum)`,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Execution Duration Summary",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("seconds"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json  | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_duration [$__range]) by (data_ExecutionStateChanged_seqNum)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_duration [$__range]) by (data_ExecutionStateChanged_seqNum)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json  | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_duration [$__range]) by (data_ExecutionStateChanged_seqNum)`,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"E2E (Commit, ARM, Execution)",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("seconds"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json |  data_CommitAndExecute_success="✅"| unwrap  data_CommitAndExecute_duration [$__range]) by (data_CommitAndExecute_seqNum)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CommitAndExecute_success="✅"| unwrap  data_CommitAndExecute_duration [$__range]) by (data_CommitAndExecute_seqNum)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CommitAndExecute_success="✅"| unwrap  data_CommitAndExecute_duration [$__range]) by (data_CommitAndExecute_seqNum)`,
					loki.Legend("Max"),
				),
			),
		),
	}
}

func failedMessagesRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"Failed Messages",
			row.Collapse(),
			row.WithTimeSeries(
				"Failed Commit",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`count_over_time({ namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="❌" [$__range])`,
					loki.Legend("{{error}}"),
				),
			),
			row.WithTimeSeries(
				"Failed Bless",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`count_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="❌" [$__range])`,
					loki.Legend("{{error}}"),
				),
			),
			row.WithTimeSeries(
				"Failed Execution",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`count_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ExecutionStateChanged_success="❌" [$__range])`,
					loki.Legend("{{error}}"),
				),
			),
			row.WithTimeSeries(
				"Failed Commit and Execution",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`count_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CommitAndExecute_success="❌" [$__range])`,
					loki.Legend("{{error}}"),
				),
			),
		),
	}
}

func reqRespRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"Requests/Responses",
			row.Collapse(),
			row.WithStat(
				"Stats",
				stat.DataSource(p.LokiDataSource),
				stat.Transparent(),
				stat.Text(stat.TextValueAndName),
				stat.Height("100px"),
				stat.TitleFontSize(20),
				stat.ValueFontSize(20),
				stat.Span(12),
				stat.WithPrometheusTarget(
					`max_over_time({namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="stats", test_group="$test_group", test_id=~"${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"}
| json
| unwrap current_time_unit [$__range]) by (test_id)`,
					prometheus.Legend("Time Unit"),
				),
				stat.WithPrometheusTarget(
					`max_over_time({namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="stats", test_group="$test_group", test_id=~"${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"}
| json
| unwrap load_duration [$__range]) by (test_id)/ 1e9 `,
					prometheus.Legend("Total Duration"),
				),
				stat.WithPrometheusTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_message_bytes_length [$__range]) by (test_id)`,
					prometheus.Legend("Max Byte Len Sent"),
				),
				stat.WithPrometheusTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_no_of_tokens_sent [$__range]) by (test_id)`,
					prometheus.Legend("Max No Of Tokens Sent"),
				),
			),
			row.WithTimeSeries(
				"Request Rate",
				timeseries.Transparent(),
				timeseries.Description("Requests triggered over test duration"),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`last_over_time({namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="stats", test_group="$test_group", test_id="${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"}| json | unwrap current_rps [$__interval]) by (test_id,gen_name)`,
					loki.Legend("Request Triggered/TimeUnit"),
				),
			),
			row.WithTimeSeries(
				"Trigger Summary",
				timeseries.Transparent(),
				timeseries.Points(),
				timeseries.Description("Latest Stage Stats"),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="stats", test_group="$test_group", test_id=~"${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"}
| json
| unwrap success [$__range]) by (test_id)`,
					loki.Legend("Successful Requests"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="stats", test_group="$test_group", test_id=~"${test_id:pipe}", source_chain="${source_chain}", dest_chain="${dest_chain}"}
| json
| unwrap failed [$__range]) by (test_id)`,
					loki.Legend("Failed Requests"),
				),
			),
			row.WithLogs(
				"All CCIP Phases Stats",
				logs.DataSource(p.LokiDataSource),
				logs.Span(12),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(
					`{namespace="${namespace}", go_test_name="${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json `,
				),
			),
		),
	}
}

func gasStatsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"CCIP Gas Stats",
			row.Collapse(),
			row.WithTimeSeries(
				"Gas Used in CCIP-Send⛽️",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("wei"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({ namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_CCIP_Send_Transaction_success="✅"| unwrap  data_CCIP_Send_Transaction_ccip_send_data_gas_used [$__range]) by (test_id)  `,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Gas Used in Commit⛽️",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("wei"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_Commit_ReportAccepted_success="✅"| unwrap  data_Commit_ReportAccepted_ccip_send_data_gas_used [$__range]) by (test_id)  `,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Gas Used in ARM Blessing⛽️",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("wei"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ReportBlessedByARM_success="✅"| unwrap  data_ReportBlessedByARM_ccip_send_data_gas_used [$__range]) by (test_id)  `,
					loki.Legend("Max"),
				),
			),
			row.WithTimeSeries(
				"Gas Used in Execution⛽️",
				timeseries.Transparent(),
				timeseries.Description(""),
				timeseries.Span(6),
				timeseries.Height("200px"),
				timeseries.DataSource(p.LokiDataSource),
				timeseries.Axis(
					axis.Unit("wei"),
				),
				timeseries.WithLokiTarget(
					`avg_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Avg"),
				),
				timeseries.WithLokiTarget(
					`min_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_ccip_send_data_gas_used [$__range]) by (test_id)`,
					loki.Legend("Min"),
				),
				timeseries.WithLokiTarget(
					`max_over_time({namespace="${namespace}", go_test_name=~"${go_test_name:pipe}", test_data_type="responses", test_group="${test_group}", test_id=~"${test_id:pipe}",source_chain="${source_chain}",dest_chain="${dest_chain}"} | json | data_ExecutionStateChanged_success="✅"| unwrap  data_ExecutionStateChanged_ccip_send_data_gas_used [$__range]) by (test_id)  `,
					loki.Legend("Max"),
				),
			),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, statsRow(p)...)
	opts = append(opts, gasStatsRow(p)...)
	opts = append(opts, failedMessagesRow(p)...)
	opts = append(opts, reqRespRow(p)...)
	return opts
}
