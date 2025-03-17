package k8spods

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/query"
)

type Props struct {
	LokiDataSource       string
	PrometheusDataSource string
}

func vars(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.VariableAsQuery(
			"namespace",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(namespace)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"job",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(`label_values(up{namespace="$namespace"}, job)`),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"pod",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request(`label_values(up{namespace="$namespace", job="$job"}, pod)`),
			query.Sort(query.NumericalAsc),
		),
	}
}

func logsRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"K8s Logs",
			row.Collapse(),
			row.WithLogs(
				"All Logs",
				logs.DataSource(p.LokiDataSource),
				logs.Span(12),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(`{namespace="$namespace", pod=~"${pod:pipe}"}`),
			),
			row.WithLogs(
				"All Errors",
				logs.DataSource(p.LokiDataSource),
				logs.Span(12),
				logs.Height("300px"),
				logs.Transparent(),
				logs.WithLokiTarget(`{namespace="$namespace", pod=~"${pod:pipe}"} | json | level=~"error|warn|fatal|panic"`),
			),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts,
		[]dashboard.Option{
			dashboard.Row(
				"K8s Pods",
				row.Collapse(),
				row.WithStat(
					"Pod Restarts",
					stat.Span(4),
					stat.Text(stat.TextValueAndName),
					stat.Orientation(stat.OrientationHorizontal),
					stat.DataSource(p.PrometheusDataSource),
					stat.SparkLine(),
					stat.SparkLineYMin(0),
					stat.WithPrometheusTarget(
						`sum(increase(kube_pod_container_status_restarts_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithStat(
					"OOM Events",
					stat.Span(4),
					stat.Text(stat.TextValueAndName),
					stat.Orientation(stat.OrientationHorizontal),
					stat.DataSource(p.PrometheusDataSource),
					stat.SparkLine(),
					stat.SparkLineYMin(0),
					stat.WithPrometheusTarget(
						`sum(container_oom_events_total{pod=~"$pod", namespace=~"${namespace}"}) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithStat(
					"OOM Killed",
					stat.Span(4),
					stat.Text(stat.TextValueAndName),
					stat.Orientation(stat.OrientationHorizontal),
					stat.DataSource(p.PrometheusDataSource),
					stat.SparkLine(),
					stat.SparkLineYMin(0),
					stat.WithPrometheusTarget(
						`kube_pod_container_status_last_terminated_reason{reason="OOMKilled", pod=~"$pod", namespace=~"${namespace}"}`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"CPU Usage",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.WithPrometheusTarget(
						`sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{pod=~"$pod", namespace=~"${namespace}"}) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"Memory Usage",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.Axis(
						axis.Unit("bytes"),
						axis.Label("Memory"),
						axis.SoftMin(0),
					),
					timeseries.WithPrometheusTarget(
						`sum(container_memory_rss{pod=~"$pod", namespace=~"${namespace}", container!=""}) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"Receive Bandwidth",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.Axis(
						axis.Unit("Bps"),
						axis.SoftMin(0),
					),
					timeseries.WithPrometheusTarget(
						`sum(irate(container_network_receive_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"Transmit Bandwidth",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.Axis(
						axis.Unit("Bps"),
						axis.SoftMin(0),
					),
					timeseries.WithPrometheusTarget(
						`sum(irate(container_network_transmit_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"Average Container Bandwidth by Namespace: Received",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.Axis(
						axis.Unit("Bps"),
						axis.SoftMin(0),
					),
					timeseries.WithPrometheusTarget(
						`avg(irate(container_network_receive_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
				row.WithTimeSeries(
					"Average Container Bandwidth by Namespace: Transmitted",
					timeseries.Span(6),
					timeseries.Height("200px"),
					timeseries.DataSource(p.PrometheusDataSource),
					timeseries.Axis(
						axis.Unit("Bps"),
						axis.SoftMin(0),
					),
					timeseries.WithPrometheusTarget(
						`avg(irate(container_network_transmit_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
						prometheus.Legend("{{pod}}"),
					),
				),
			),
		}...,
	)
	opts = append(opts, logsRow(p)...)
	return opts
}
