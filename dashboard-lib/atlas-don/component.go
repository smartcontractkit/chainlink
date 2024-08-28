package atlas_don

import (
	"fmt"

	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/query"
)

type Props struct {
	PrometheusDataSource string
	PlatformOpts         PlatformOpts
	OcrVersion           string
}

func vars(p Props) []dashboard.Option {
	variableFeedId := "feed_id"
	if p.OcrVersion == "ocr3" {
		variableFeedId = "feed_id_name"
	}

	variableQueryContract := dashboard.VariableAsQuery(
		"contract",
		query.DataSource(p.PrometheusDataSource),
		query.Multiple(),
		query.IncludeAll(),
		query.Request(fmt.Sprintf(`label_values(`+p.OcrVersion+`_contract_config_f{job="$job"}, %s)`, "contract")),
		query.Sort(query.NumericalAsc),
	)

	variableQueryFeedId := dashboard.VariableAsQuery(
		variableFeedId,
		query.DataSource(p.PrometheusDataSource),
		query.Multiple(),
		query.IncludeAll(),
		query.Request(fmt.Sprintf(`label_values(`+p.OcrVersion+`_contract_config_f{job="$job", contract="$contract"}, %s)`, variableFeedId)),
		query.Sort(query.NumericalAsc),
	)

	variables := []dashboard.Option{
		variableQueryContract,
	}

	switch p.OcrVersion {
	case "ocr":
		break
	case "ocr2":
		variables = append(variables, variableQueryFeedId)
		break
	case "ocr3":
		variables = append(variables, variableQueryFeedId)
		break
	}

	return variables
}

func summary(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("Summary",
			row.WithStat(
				"Telemetry Down",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which jobs are not receiving any telemetry?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_telemetry_down{`+p.PlatformOpts.LabelQuery+`} == 1`,
					prometheus.Legend("{{job}} | {{report_type}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
			row.WithStat(
				"Oracles Down",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which NOPs are not providing any telemetry?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.ValueType(stat.Last),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_oracle_telemetry_down_except_telemetry_down{job=~"${job}", oracle!="csa_unknown"} == 1`,
					prometheus.Legend("{{oracle}} | {{report_type}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
			row.WithStat(
				"Feeds reporting failure",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which feeds are failing to report?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.ValueType(stat.Last),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_feed_reporting_failure_except_feed_telemetry_down{job=~"${job}", oracle!="csa_unknown"} == 1`,
					prometheus.Legend("{{feed_id_name}} on {{job}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
			row.WithStat(
				"Feed telemetry Down",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which feeds are not receiving any telemetry?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.ValueType(stat.Last),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_feed_telemetry_down_except_telemetry_down{job=~"${job}"} == 1`,
					prometheus.Legend("{{feed_id_name}} on {{job}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
			row.WithStat(
				"Oracles no observations",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which NOPs are not providing observations?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.ValueType(stat.Last),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_oracle_blind_except_telemetry_down{job=~"${job}"} == 1`,
					prometheus.Legend("{{oracle}} | {{report_type}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
			row.WithStat(
				"Oracles not contributing observations to feeds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("Which oracles are failing to make observations on feeds they should be participating in?"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(12),
				stat.Span(4),
				stat.ValueType(stat.Last),
				stat.WithPrometheusTarget(
					`bool:`+p.OcrVersion+`_oracle_feed_no_observations_except_oracle_blind_except_feed_reporting_failure_except_feed_telemetry_down{job=~"${job}"} == 1`,
					prometheus.Legend("{{oracle}} | {{report_type}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#008000", Value: float64Ptr(0.0)},
					{Color: "#FF0000", Value: float64Ptr(1.0)},
				}),
			),
		),
	}
}

func ocrContractConfigOracle(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Contract Oracle",
			row.Collapse(),
			row.WithStat(
				"OCR Contract Oracle Active",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextName),
				stat.Description("set to one as long as an oracle is on a feed"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.ValueFontSize(12),
				stat.Span(12),
				stat.WithPrometheusTarget(
					`sum(`+p.OcrVersion+`_contract_oracle_active{`+p.PlatformOpts.LabelQuery+`}) by (contract, oracle)`,
					prometheus.Legend("{{ contract }} - {{oracle}}"),
				),
				stat.AbsoluteThresholds([]stat.ThresholdStep{
					{Color: "#FF0000", Value: float64Ptr(0.0)},
					{Color: "#008000", Value: float64Ptr(1.0)},
				}),
			),
		),
	}
}

func ocrContractConfigNodes(p Props) []dashboard.Option {
	variableFeedId := "feed_id"
	if p.OcrVersion == "ocr3" {
		variableFeedId = "feed_id_name"
	}

	var options []timeseries.Option

	options = append(options, timeseries.Span(12),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.Axis(
			axis.Min(0),
		),
	)

	switch p.OcrVersion {
	case "ocr":
		options = append(options, timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{contract}}"),
		))
		break
	case "ocr2":
		options = append(options, timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{"+variableFeedId+"}}"),
		))
		break
	case "ocr3":
		options = append(options, timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_message_observe_total_nop_count{contract=~"${contract}", `+variableFeedId+`=~"${`+variableFeedId+`}", job=~"${job}"}`,
			prometheus.Legend("{{"+variableFeedId+"}}"),
		))
		break
	}

	options = append(options,
		timeseries.WithPrometheusTarget(
			`avg(2 * `+p.OcrVersion+`_contract_config_r_max{`+p.PlatformOpts.LabelQuery+`} + 4)`,
			prometheus.Legend("Max nodes"),
		),
		timeseries.WithPrometheusTarget(
			`avg(2 * `+p.OcrVersion+`_contract_config_f{`+p.PlatformOpts.LabelQuery+`} + 1)`,
			prometheus.Legend("Min nodes"),
		),
	)

	return []dashboard.Option{
		dashboard.Row("DON Nodes",
			row.Collapse(),
			row.WithTimeSeries(
				"Number of NOPs",
				options...,
			),
		),
	}
}

func priceReporting(p Props) []dashboard.Option {
	telemetryP2PReceivedTotal := row.WithTimeSeries(
		"P2P messages received",
		timeseries.Span(12),
		timeseries.Height("600px"),
		timeseries.Description("From an individual node's perspective, how many messages are they receiving from other nodes? Uses ocr_telemetry_p2p_received_total"),
		timeseries.Axis(
			axis.Min(0),
		),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			`sum by (sender, receiver) (increase(`+p.OcrVersion+`_telemetry_p2p_received_total{job=~"${job}"}[5m]))`,
			prometheus.Legend("{{sender}} > {{receiver}}"),
		),
	)

	telemetryP2PReceivedTotalRate := row.WithTimeSeries(
		"P2P messages received Rate",
		timeseries.Span(12),
		timeseries.Height("600px"),
		timeseries.Description("From an individual node's perspective, how many messages are they receiving from other nodes? Uses ocr_telemetry_p2p_received_total"),
		timeseries.Axis(
			axis.Min(0),
		),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			`sum by (sender, receiver) (rate(`+p.OcrVersion+`_telemetry_p2p_received_total{job=~"${job}"}[5m]))`,
			prometheus.Legend("{{sender}} > {{receiver}}"),
		),
	)

	telemetryObservationAsk := row.WithTimeSeries(
		"Ask observation in MessageObserve sent",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_observation_ask{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryObservation := row.WithTimeSeries(
		"Price observation in MessageObserve sent",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_observation{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryObservationBid := row.WithTimeSeries(
		"Bid observation in MessageObserve sent",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_observation_bid{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryMessageProposeObservationAsk := row.WithTimeSeries(
		"Ask MessagePropose observations",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_message_propose_observation_ask{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryMessageProposeObservation := row.WithTimeSeries(
		"Price MessagePropose observations",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_message_propose_observation{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryMessageProposeObservationBid := row.WithTimeSeries(
		"Bid MessagePropose observations",
		timeseries.Span(12),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			``+p.OcrVersion+`_telemetry_message_propose_observation_bid{`+p.PlatformOpts.LabelQuery+`}`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryMessageProposeObservationTotal := row.WithTimeSeries(
		"Total number of observations included in MessagePropose",
		timeseries.Span(12),
		timeseries.Description("How often is a node's observation included in the report?"),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.Axis(
			axis.Min(0),
		),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			`rate(`+p.OcrVersion+`_telemetry_message_propose_observation_total{`+p.PlatformOpts.LabelQuery+`}[5m])`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	telemetryMessageObserveTotal := row.WithTimeSeries(
		"Total MessageObserve sent",
		timeseries.Span(12),
		timeseries.Description("From an individual node's perspective, how often are they sending an observation?"),
		timeseries.Legend(timeseries.ToTheRight),
		timeseries.Axis(
			axis.Min(0),
		),
		timeseries.DataSource(p.PrometheusDataSource),
		timeseries.WithPrometheusTarget(
			`rate(`+p.OcrVersion+`_telemetry_message_observe_total{`+p.PlatformOpts.LabelQuery+`}[5m])`,
			prometheus.Legend("{{oracle}}"),
		),
	)

	panels := []row.Option{
		row.Collapse(),
	}

	switch p.OcrVersion {
	case "ocr":
		panels = append(panels, telemetryP2PReceivedTotal)
		panels = append(panels, telemetryP2PReceivedTotalRate)
		panels = append(panels, telemetryObservation)
		panels = append(panels, telemetryMessageObserveTotal)
		break
	case "ocr2":
		panels = append(panels, telemetryP2PReceivedTotal)
		panels = append(panels, telemetryP2PReceivedTotalRate)
		panels = append(panels, telemetryObservation)
		panels = append(panels, telemetryMessageObserveTotal)
		break
	case "ocr3":
		panels = append(panels, telemetryP2PReceivedTotal)
		panels = append(panels, telemetryP2PReceivedTotalRate)
		panels = append(panels, telemetryObservationAsk)
		panels = append(panels, telemetryObservation)
		panels = append(panels, telemetryObservationBid)
		panels = append(panels, telemetryMessageProposeObservationAsk)
		panels = append(panels, telemetryMessageProposeObservation)
		panels = append(panels, telemetryMessageProposeObservationBid)
		panels = append(panels, telemetryMessageProposeObservationTotal)
		panels = append(panels, telemetryMessageObserveTotal)
		break
	}

	return []dashboard.Option{
		dashboard.Row("Price Reporting", panels...),
	}
}

func ocrContractConfigDelta(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Contract Config Delta",
			row.Collapse(),
			row.WithStat(
				"relativeDeviationThreshold",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_alpha{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"maxContractValueAgeSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_c_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"observationGracePeriodSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_grace_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"badEpochTimeoutSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_progress_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"resendIntervalSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_resend_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"roundIntervalSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_round_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
			row.WithStat(
				"transmissionStageTimeoutSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(28),
				stat.Span(4),
				stat.SparkLine(),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_stage_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{ contract }}"),
				),
			),
		),
	}
}

func roundEpochProgression(p Props) []dashboard.Option {
	variableFeedId := "feed_id"
	if p.OcrVersion == "ocr3" {
		variableFeedId = "feed_id_name"
	}

	return []dashboard.Option{
		dashboard.Row("Round / Epoch Progression",
			row.Collapse(),
			row.WithTimeSeries(
				"Agreed Epoch Progression",
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("short"),
				),
				timeseries.WithPrometheusTarget(
					``+p.OcrVersion+`_telemetry_feed_agreed_epoch{`+variableFeedId+`=~"${`+variableFeedId+`}"}`,
					prometheus.Legend("{{"+variableFeedId+"}}"),
				),
			),
			row.WithTimeSeries(
				"Round Epoch Progression",
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("short"),
				),
				timeseries.WithPrometheusTarget(
					``+p.OcrVersion+`_telemetry_epoch_round{`+variableFeedId+`=~"${`+variableFeedId+`}"}`,
					prometheus.Legend("{{oracle}}"),
				),
			),
			row.WithTimeSeries(
				"Rounds Started",
				timeseries.Description("Tracks individual nodes firing \"new round\" message via telemetry (not part of P2P messages)"),
				timeseries.Span(4),
				timeseries.Height("300px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("short"),
				),
				timeseries.WithPrometheusTarget(
					`rate(`+p.OcrVersion+`_telemetry_round_started_total{`+variableFeedId+`=~"${`+variableFeedId+`}"}[1m])`,
					prometheus.Legend("{{oracle}}"),
				),
			),
			row.WithTimeSeries(
				"Telemetry Ingested",
				timeseries.Span(12),
				timeseries.Height("300px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Unit("short"),
				),
				timeseries.Legend(timeseries.ToTheRight),
				timeseries.WithPrometheusTarget(
					`rate(`+p.OcrVersion+`_telemetry_ingested_total{`+variableFeedId+`=~"${`+variableFeedId+`}"}[1m])`,
					prometheus.Legend("{{oracle}}"),
				),
			),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, summary(p)...)
	opts = append(opts, ocrContractConfigOracle(p)...)
	opts = append(opts, ocrContractConfigNodes(p)...)
	opts = append(opts, priceReporting(p)...)
	opts = append(opts, roundEpochProgression(p)...)
	opts = append(opts, ocrContractConfigDelta(p)...)
	return opts
}
