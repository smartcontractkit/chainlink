package core_ocrv2_ccip

import (
	"fmt"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
)

type Props struct {
	PrometheusDataSource string
	PluginName           string
}

func quantileRowOpts(ds string, pluginName string, perc string) row.Option {
	return row.WithTimeSeries(
		fmt.Sprintf("(%s) OCR2 duration (%s)", pluginName, perc),
		timeseries.Span(6),
		timeseries.Height("200px"),
		timeseries.DataSource(ds),
		timeseries.WithPrometheusTarget(
			fmt.Sprintf(`histogram_quantile(%s, sum(rate(ocr2_reporting_plugin_observation_time_bucket{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__rate_interval])) by (le)) / 1e9`, perc, pluginName),
			prometheus.Legend("Observation"),
		),
		timeseries.WithPrometheusTarget(
			fmt.Sprintf(`histogram_quantile(%s, sum(rate(ocr2_reporting_plugin_report_time_bucket{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__rate_interval])) by (le)) / 1e9`, perc, pluginName),
			prometheus.Legend("Report"),
		),
		timeseries.WithPrometheusTarget(
			fmt.Sprintf(`histogram_quantile(%s, sum(rate(ocr2_reporting_plugin_should_accept_finalized_report_time_bucket{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__rate_interval])) by (le)) / 1e9`, perc, pluginName),
			prometheus.Legend("ShouldAcceptFinalizedReport"),
		),
		timeseries.WithPrometheusTarget(
			fmt.Sprintf(`histogram_quantile(%s, sum(rate(ocr2_reporting_plugin_should_transmit_accepted_report_time_bucket{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__rate_interval])) by (le)) / 1e9`, perc, pluginName),
			prometheus.Legend("ShouldTransmitAcceptedReport"),
		),
	)
}

func ocrv2PluginObservationStageQuantiles(p Props) []dashboard.Option {
	opts := make([]row.Option, 0)
	opts = append(opts,
		row.Collapse(),
		row.WithTimeSeries(
			fmt.Sprintf("(%s) OCR2 RPS by phase", p.PluginName),
			timeseries.Span(6),
			timeseries.Height("200px"),
			timeseries.DataSource(p.PrometheusDataSource),
			timeseries.WithPrometheusTarget(
				fmt.Sprintf(`sum(rate(ocr2_reporting_plugin_observation_time_count{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__range]))`, p.PluginName),
				prometheus.Legend("Observation"),
			),
			timeseries.WithPrometheusTarget(
				fmt.Sprintf(`sum(rate(ocr2_reporting_plugin_report_time_count{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__range]))`, p.PluginName),
				prometheus.Legend("Report"),
			),
			timeseries.WithPrometheusTarget(
				fmt.Sprintf(`sum(rate(ocr2_reporting_plugin_should_accept_finalized_report_time_count{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__range]))`, p.PluginName),
				prometheus.Legend("ShouldAcceptFinalizedReport"),
			),
			timeseries.WithPrometheusTarget(
				fmt.Sprintf(`sum(rate(ocr2_reporting_plugin_should_transmit_accepted_report_time_count{plugin="%s", job=~"$instance", chainID=~"$evmChainID"}[$__range]))`, p.PluginName),
				prometheus.Legend("ShouldTransmitAcceptedReport"),
			),
		),
		quantileRowOpts(p.PrometheusDataSource, p.PluginName, "0.5"),
		quantileRowOpts(p.PrometheusDataSource, p.PluginName, "0.9"),
		quantileRowOpts(p.PrometheusDataSource, p.PluginName, "0.99"),
	)
	return []dashboard.Option{
		dashboard.Row(
			fmt.Sprintf("OCRv2 Metrics - Plugin: %s", p.PluginName),
			opts...,
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := make([]dashboard.Option, 0)
	opts = append(opts, ocrv2PluginObservationStageQuantiles(p)...)
	return opts
}
