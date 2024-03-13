package atlas_don

import (
	"fmt"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/gauge"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/prometheus"
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
					``+p.OcrVersion+`_contract_oracle_active{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_oracle_active{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_n{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_f{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Max Faulty Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_f{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_r_max{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"Max Round Node Count",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_r_max{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_alpha{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"relativeDeviationThreshold",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_alpha{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_c_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"maxContractValueAgeSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_c_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_grace_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"observationGracePeriodSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_grace_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_progress_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"badEpochTimeoutSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_progress_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_resend_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"resendIntervalSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_resend_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_round_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"roundIntervalSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_round_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
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
					``+p.OcrVersion+`_contract_config_delta_stage_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
			row.WithGauge(
				"transmissionStageTimeoutSeconds",
				gauge.Span(6),
				gauge.Orientation(gauge.OrientationVertical),
				gauge.DataSource(p.PrometheusDataSource),
				gauge.WithPrometheusTarget(
					``+p.OcrVersion+`_contract_config_delta_stage_seconds{`+p.PlatformOpts.LabelQuery+`contract=~"${contract}"}`,
					prometheus.Legend("{{"+p.PlatformOpts.LegendString+"}} - {{ contract }}"),
				),
			),
		),
	}
}

func ocrTelemetryFeed(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Telemetry Feed",
			row.Collapse(),
		),
	}
}

func ocrTelemetryP2P(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Telemetry P2P",
			row.Collapse(),
		),
	}
}

func ocrTelemetryOthers(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row("OCR Telemetry Others",
			row.Collapse(),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, ocrContractConfigOracle(p)...)
	opts = append(opts, ocrContractConfigNodes(p)...)
	opts = append(opts, ocrContractConfigDelta(p)...)
	return opts
}
