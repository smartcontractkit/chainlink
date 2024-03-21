package atlas_don

import (
	"fmt"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/gauge"
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
	return []dashboard.Option{
		dashboard.VariableAsQuery(
			"contract",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "contract")),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"feed_id_name",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(fmt.Sprintf("label_values(%s)", "feed_id_name")),
			query.Sort(query.NumericalAsc),
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
				stat.Text(stat.TextValueAndName),
				stat.Description("set to one as long as an oracle is on a feed"),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(12),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_oracle_active{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }} - {{oracle}}"),
				),
			),
			row.WithGauge(
				"OCR Contract Oracle Active",
				gauge.Span(12),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.Description("set to one as long as an oracle is on a feed"),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_oracle_active{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }} - {{oracle}}"),
				),
			),
		),
	}
}

func ocrContractConfigNodes(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Contract Config Nodes",
			row.Collapse(),
			row.WithStat(
				"Node Count",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"Max Faulty Node Count",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_f{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Max Faulty Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_f{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"Max Round Node Count",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_r_max{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Max Round Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_r_max{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
		),
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
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_alpha{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"relativeDeviationThreshold",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_alpha{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"maxContractValueAgeSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_c_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"maxContractValueAgeSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_c_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"observationGracePeriodSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_grace_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"observationGracePeriodSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_grace_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"badEpochTimeoutSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_progress_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"badEpochTimeoutSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_progress_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"resendIntervalSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_resend_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"resendIntervalSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_resend_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"roundIntervalSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_round_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"roundIntervalSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_round_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithStat(
				"transmissionStageTimeoutSeconds",
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationHorizontal),
				stat.TitleFontSize(12),
				stat.ValueFontSize(20),
				stat.Span(6),
				stat.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_stage_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"transmissionStageTimeoutSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_stage_seconds{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
		),
	}
}

func roundEpochProgression(p Props) []dashboard.Option {
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
					``+p.OcrVersion+`_telemetry_feed_agreed_epoch{`+p.PlatformOpts.LabelQuery+`feed_id_name=~"${feed_id_name}"}`,
					prometheus.Legend("{{feed_id_name}}"),
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
					``+p.OcrVersion+`_telemetry_epoch_round{`+p.PlatformOpts.LabelQuery+`feed_id_name=~"${feed_id_name}"}`,
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
					`rate(`+p.OcrVersion+`_telemetry_round_started_total{`+p.PlatformOpts.LabelQuery+`feed_id_name=~"${feed_id_name}"}[1m])`,
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
				timeseries.WithPrometheusTarget(
					`rate(`+p.OcrVersion+`_telemetry_ingested_total{`+p.PlatformOpts.LabelQuery+`feed_id_name=~"${feed_id_name}"}[1m])`,
					prometheus.Legend("{{oracle}}"),
				),
			),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, ocrContractConfigOracle(p)...)
	opts = append(opts, ocrContractConfigNodes(p)...)
	opts = append(opts, ocrContractConfigDelta(p)...)
	opts = append(opts, roundEpochProgression(p)...)
	return opts
}
