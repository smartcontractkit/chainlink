package atlasdon

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"

	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"
)

func BuildDashboard(name string, dataSourceMetric string, platform string, ocrVersion string) (dashboard.Dashboard, error) {
	props := Props{
		MetricsDataSource: dataSourceMetric,
		PlatformOpts:      PlatformPanelOpts(platform, ocrVersion),
		OcrVersion:        ocrVersion,
	}

	builder := dashboard.NewDashboardBuilder(name).
		Tags([]string{"DON", ocrVersion}).
		Refresh("30s").
		Time("now-30m", "now")

	utils.AddVars(builder, vars(props))

	builder.WithRow(dashboard.NewRowBuilder("Summary"))
	utils.AddPanels(builder, summary(props))

	builder.WithRow(dashboard.NewRowBuilder("OCR Contract Oracle"))
	utils.AddPanels(builder, ocrContractConfigOracle(props))

	builder.WithRow(dashboard.NewRowBuilder("DON Nodes"))
	utils.AddPanels(builder, ocrContractConfigNodes(props))

	builder.WithRow(dashboard.NewRowBuilder("Price Reporting"))
	utils.AddPanels(builder, priceReporting(props))

	builder.WithRow(dashboard.NewRowBuilder("Round / Epoch Progression"))
	utils.AddPanels(builder, roundEpochProgression(props))

	builder.WithRow(dashboard.NewRowBuilder("OCR Contract Config Delta"))
	utils.AddPanels(builder, ocrContractConfigDelta(props))

	return builder.Build()
}

func vars(p Props) []cog.Builder[dashboard.VariableModel] {
	var variables []cog.Builder[dashboard.VariableModel]

	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "namespace", "Namespace",
			`label_values(namespace)`, true).Regex("otpe[1-3]?$"))

	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "job", "Job",
			`label_values(up{namespace="$namespace"}, job)`, true))

	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "pod", "Pod",
			`label_values(up{namespace="$namespace", job="$job"}, pod)`, true))

	variableFeedID := "feed_id"
	if p.OcrVersion == "ocr3" {
		variableFeedID = "feed_id_name"
	}

	variableQueryContract := utils.QueryVariable(p.MetricsDataSource, "contract", "Contract",
		`label_values(`+p.OcrVersion+`_contract_config_f{job="$job"}, contract)`, true)

	variableQueryFeedID := utils.QueryVariable(p.MetricsDataSource, variableFeedID, "Feed ID",
		`label_values(`+p.OcrVersion+`_contract_config_f{job="$job", contract="$contract"}, `+variableFeedID+`)`, true)

	variables = append(variables, variableQueryContract)

	switch p.OcrVersion {
	case "ocr2":
		variables = append(variables, variableQueryFeedID)
	case "ocr3":
		variables = append(variables, variableQueryFeedID)
	}

	return variables
}

func summary(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Telemetry Down",
		"Which jobs are not receiving any telemetry?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_telemetry_down{` + p.PlatformOpts.LabelQuery + `} == 1`,
			Legend: "{{job}} | {{report_type}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Oracle Down",
		"Which NOPs are not providing any telemetry?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_oracle_telemetry_down_except_telemetry_down{job=~"${job}", oracle!="csa_unknown"} == 1`,
			Legend: "{{oracle}} | {{report_type}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Feeds reporting failure",
		"Which feeds are failing to report?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_feed_reporting_failure_except_feed_telemetry_down{job=~"${job}", oracle!="csa_unknown"} == 1`,
			Legend: "{{feed_id_name}} on {{job}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Feed telemetry Down",
		"Which feeds are not receiving any telemetry?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_feed_telemetry_down_except_telemetry_down{job=~"${job}"} == 1`,
			Legend: "{{feed_id_name}} on {{job}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Oracles no observations",
		"Which NOPs are not providing observations?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_oracle_blind_except_telemetry_down{job=~"${job}"} == 1`,
			Legend: "{{oracle}} | {{report_type}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Oracles not contributing observations to feeds",
		"Which oracles are failing to make observations on feeds they should be participating in?",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `bool:` + p.OcrVersion + `_oracle_feed_no_observations_except_oracle_blind_except_feed_reporting_failure_except_feed_telemetry_down{job=~"${job}"} == 1`,
			Legend: "{{oracle}} | {{report_type}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "green"},
				{Value: utils.Float64Ptr(0.99), Color: "red"},
			})),
	)

	return panelsArray
}

func ocrContractConfigOracle(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"OCR Contract Oracle Active",
		"set to one as long as an oracle is on a feed",
		8,
		24,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(` + p.OcrVersion + `_contract_oracle_active{` + p.PlatformOpts.LabelQuery + `}) by (contract, oracle)`,
			Legend: "{{contract}} - {{oracle}}",
		},
	).Thresholds(
		dashboard.NewThresholdsConfigBuilder().
			Mode(dashboard.ThresholdsModeAbsolute).
			Steps([]dashboard.Threshold{
				{Value: utils.Float64Ptr(0), Color: "red"},
				{Value: utils.Float64Ptr(0.99), Color: "green"},
			})),
	)

	return panelsArray
}

func ocrContractConfigNodes(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	var variableFeedID string
	switch p.OcrVersion {
	case "ocr":
		variableFeedID = "contract"
	case "ocr2":
		variableFeedID = "feed_id"
	case "ocr3":
		variableFeedID = "feed_id_name"
	}

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Number of NOPs",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_n{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{` + variableFeedID + `}}`,
		},
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_r_max{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `Max nodes`,
		},
		utils.PrometheusQuery{
			Query:  `avg(2 * ` + p.OcrVersion + `_contract_config_f{` + p.PlatformOpts.LabelQuery + `} + 1)`,
			Legend: `Min nodes`,
		},
	).Min(0))

	return panelsArray
}

func priceReporting(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	telemetryP2PReceivedTotal := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"P2P messages received",
		"From an individual node's perspective, how many messages are they receiving from other nodes? Uses ocr_telemetry_p2p_received_total",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum by (sender, receiver) (increase(` + p.OcrVersion + `_telemetry_p2p_received_total{job=~"${job}"}[5m]))`,
			Legend: `{{sender}} > {{receiver}}`,
		},
	)

	telemetryP2PReceivedTotalRate := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"P2P messages received Rate",
		"From an individual node's perspective, how many messages are they receiving from other nodes? Uses ocr_telemetry_p2p_received_total",
		6,
		24,
		1,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum by (sender, receiver) (rate(` + p.OcrVersion + `_telemetry_p2p_received_total{job=~"${job}"}[5m]))`,
			Legend: `{{sender}} > {{receiver}}`,
		},
	)

	telemetryObservationAsk := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Ask observation in MessageObserve sent",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_observation_ask{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryObservation := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Price observation in MessageObserve sent",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_observation{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryObservationBid := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Bid observation in MessageObserve sent",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_observation_bid{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryMessageProposeObservationAsk := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Ask MessagePropose observations",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_message_propose_observation_ask{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryMessageProposeObservation := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Price MessagePropose observations",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_message_propose_observation{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryMessageProposeObservationBid := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Bid MessagePropose observations",
		"",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_message_propose_observation_bid{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryMessageProposeObservationTotal := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Total number of observations included in MessagePropose",
		"How often is a node's observation included in the report?",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_message_propose_observation_total{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: `{{oracle}}`,
		},
	)

	telemetryMessageObserveTotal := utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Total MessageObserve sent",
		"From an individual node's perspective, how often are they sending an observation?",
		6,
		24,
		1,
		"",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `rate(` + p.OcrVersion + `_telemetry_message_observe_total{` + p.PlatformOpts.LabelQuery + `}[5m])`,
			Legend: `{{oracle}}`,
		},
	)

	switch p.OcrVersion {
	case "ocr":
		panelsArray = append(panelsArray, telemetryP2PReceivedTotal)
		panelsArray = append(panelsArray, telemetryP2PReceivedTotalRate)
		panelsArray = append(panelsArray, telemetryObservation)
		panelsArray = append(panelsArray, telemetryMessageObserveTotal)
	case "ocr2":
		panelsArray = append(panelsArray, telemetryP2PReceivedTotal)
		panelsArray = append(panelsArray, telemetryP2PReceivedTotalRate)
		panelsArray = append(panelsArray, telemetryObservation)
		panelsArray = append(panelsArray, telemetryMessageObserveTotal)
	case "ocr3":
		panelsArray = append(panelsArray, telemetryP2PReceivedTotal)
		panelsArray = append(panelsArray, telemetryP2PReceivedTotalRate)
		panelsArray = append(panelsArray, telemetryObservationAsk)
		panelsArray = append(panelsArray, telemetryObservation)
		panelsArray = append(panelsArray, telemetryObservationBid)
		panelsArray = append(panelsArray, telemetryMessageProposeObservationAsk)
		panelsArray = append(panelsArray, telemetryMessageProposeObservation)
		panelsArray = append(panelsArray, telemetryMessageProposeObservationBid)
		panelsArray = append(panelsArray, telemetryMessageProposeObservationTotal)
		panelsArray = append(panelsArray, telemetryMessageObserveTotal)
	}

	return panelsArray
}

func roundEpochProgression(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	variableFeedID := "feed_id"
	if p.OcrVersion == "ocr3" {
		variableFeedID = "feed_id_name"
	}

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Agreed Epoch Progression",
		"",
		6,
		8,
		1,
		"short",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_feed_agreed_epoch{` + variableFeedID + `=~"${` + variableFeedID + `}"}`,
			Legend: `{{` + variableFeedID + `}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Round Epoch Progression",
		"",
		6,
		8,
		1,
		"short",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_telemetry_epoch_round{` + variableFeedID + `=~"${` + variableFeedID + `}"}`,
			Legend: `{{oracle}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Rounds Started",
		`Tracks individual nodes firing "new round" message via telemetry (not part of P2P messages)`,
		6,
		8,
		1,
		"short",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `rate(` + p.OcrVersion + `_telemetry_round_started_total{` + variableFeedID + `=~"${` + variableFeedID + `}"}[1m])`,
			Legend: `{{oracle}}`,
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Telemetry Ingested",
		"",
		6,
		8,
		1,
		"short",
		common.LegendPlacementRight,
		utils.PrometheusQuery{
			Query:  `rate(` + p.OcrVersion + `_telemetry_ingested_total{` + variableFeedID + `=~"${` + variableFeedID + `}"}[1m])`,
			Legend: `{{oracle}}`,
		},
	))

	return panelsArray
}

func ocrContractConfigDelta(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Relative Deviation Threshold",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_alpha{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Max Contract Value Age Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_c_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Observation Grace Period Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_grace_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Bad Epoch Timeout Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_progress_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Resend Interval Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_resend_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Round Interval Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_round_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Transmission Stage Timeout Seconds",
		"",
		4,
		8,
		1,
		"",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `` + p.OcrVersion + `_contract_config_delta_stage_seconds{` + p.PlatformOpts.LabelQuery + `}`,
			Legend: "{{contract}}",
		},
	))

	return panelsArray
}
