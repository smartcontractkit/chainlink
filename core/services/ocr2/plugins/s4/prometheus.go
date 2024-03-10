package s4

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promReportingPluginQuery = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "s4_reporting_plugin_query",
		Help: "Metric to track number of ReportingPlugin.Query() calls",
	}, []string{"product"})

	promReportingPluginObservation = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "s4_reporting_plugin_observation",
		Help: "Metric to track number of ReportingPlugin.Observation() calls",
	}, []string{"product"})

	promReportingPluginReport = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "s4_reporting_plugin_report",
		Help: "Metric to track number of ReportingPlugin.Report() calls",
	}, []string{"product"})

	promReportingPluginShouldAccept = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "s4_reporting_plugin_accept",
		Help: "Metric to track number of ReportingPlugin.ShouldAcceptFinalizedReport() calls",
	}, []string{"product"})

	promReportingPluginsQueryByteSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "s4_reporting_plugin_query_byte_size",
		Help: "Metric to track query byte size returned by ReportingPlugin.Query()",
	}, []string{"product"})

	promReportingPluginsQueryRowsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "s4_reporting_plugin_query_rows_count",
		Help: "Metric to track rows count returned by ReportingPlugin.Query()",
	}, []string{"product"})

	promReportingPluginsObservationRowsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "s4_reporting_plugin_observation_rows_count",
		Help: "Metric to track rows count returned by ReportingPlugin.Observation()",
	}, []string{"product"})

	promReportingPluginsReportRowsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "s4_reporting_plugin_report_rows_count",
		Help: "Metric to track rows count returned by ReportingPlugin.Report()",
	}, []string{"product"})

	promReportingPluginWrongSigCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "s4_reporting_plugin_wrong_sig_count",
		Help: "Metric to track number of rows having wrong signature",
	}, []string{"product"})

	promReportingPluginsExpiredRows = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "s4_reporting_plugin_expired_rows",
		Help: "Metric to track number of expired rows",
	}, []string{"product"})
)
