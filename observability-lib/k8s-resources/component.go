package k8sresources

import (
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"

	"github.com/smartcontractkit/chainlink-common/observability-lib/utils"
)

func BuildDashboard(name string, dataSourceMetric string, dataSourceLog string) (dashboard.Dashboard, error) {
	props := Props{
		MetricsDataSource: dataSourceMetric,
		LogsDataSource:    dataSourceLog,
	}

	builder := dashboard.NewDashboardBuilder(name).
		Tags([]string{"Core", "Node", "Kubernetes", "Resources"}).
		Refresh("30s").
		Time("now-30m", "now")

	utils.AddVars(builder, vars(props))

	builder.WithRow(dashboard.NewRowBuilder("Headlines"))
	utils.AddPanels(builder, headlines(props))

	builder.WithRow(dashboard.NewRowBuilder("Pod Status"))
	utils.AddPanels(builder, podStatus(props))

	builder.WithRow(dashboard.NewRowBuilder("Resources Usage"))
	utils.AddPanels(builder, resourcesUsage(props))

	builder.WithRow(dashboard.NewRowBuilder("Network Usage"))
	utils.AddPanels(builder, networkUsage(props))

	builder.WithRow(dashboard.NewRowBuilder("Disk Usage"))
	utils.AddPanels(builder, diskUsage(props))

	return builder.Build()
}

func vars(p Props) []cog.Builder[dashboard.VariableModel] {
	var variables []cog.Builder[dashboard.VariableModel]

	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "env", "Environment", `label_values(up, env)`, false))
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "cluster", "Cluster", `label_values(up{env="$env"}, cluster)`, false))
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "namespace", "Namespace", `label_values(up{env="$env", cluster="$cluster"}, namespace)`, false))
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "job", "Job", `label_values(up{env="$env", cluster="$cluster", namespace="$namespace"}, job)`, false))
	variables = append(variables,
		utils.QueryVariable(p.MetricsDataSource, "pod", "Pod", `label_values(up{env="$env", cluster="$cluster", namespace="$namespace", job="$job"}, pod)`, true))

	return variables
}

func headlines(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"CPU Utilisation (from requests)",
		"",
		4,
		6,
		1,
		"percent",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValue,
		common.VizOrientationHorizontal,
	).WithTarget(
		prometheus.NewDataqueryBuilder().
			Expr(`100 * sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container)`).
			LegendFormat("{{pod}}").
			Instant(true),
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"CPU Utilisation (from limits)",
		"",
		4,
		6,
		1,
		"percent",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValue,
		common.VizOrientationHorizontal,
	).WithTarget(
		prometheus.NewDataqueryBuilder().
			Expr(`100 * sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container)`).
			LegendFormat("{{pod}}").
			Instant(true),
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Memory Utilisation (from requests)",
		"",
		4,
		6,
		2,
		"percent",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValue,
		common.VizOrientationHorizontal,
	).WithTarget(
		prometheus.NewDataqueryBuilder().
			Expr(`100 * sum(container_memory_working_set_bytes{job="kubelet", metrics_path="/metrics/cadvisor", cluster="$cluster", namespace="$namespace", pod="$pod", image!=""}) by (container) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container)`).
			LegendFormat("{{pod}}").
			Instant(true),
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Memory Utilisation (from limits)",
		"",
		4,
		6,
		1,
		"percent",
		common.BigValueColorModeValue,
		common.BigValueGraphModeNone,
		common.BigValueTextModeValue,
		common.VizOrientationHorizontal,
	).WithTarget(
		prometheus.NewDataqueryBuilder().
			Expr(`100 * sum(container_memory_working_set_bytes{job="kubelet", metrics_path="/metrics/cadvisor", cluster="$cluster", namespace="$namespace", pod="$pod", container!="", image!=""}) by (container) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_limits{cluster="$cluster", namespace="$namespace", pod="$pod"}) by (container)`).
			LegendFormat("{{pod}}").
			Instant(true),
	))

	return panelsArray
}

func podStatus(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"Pod Restarts",
		"Number of pod restarts",
		4,
		8,
		0,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(increase(kube_pod_container_status_restarts_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"OOM Events",
		"Out-of-memory number of events",
		4,
		8,
		0,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `sum(container_oom_events_total{pod=~"$pod", namespace=~"${namespace}"}) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.StatPanel(
		p.MetricsDataSource,
		"OOM Killed",
		"",
		4,
		8,
		0,
		"",
		common.BigValueColorModeNone,
		common.BigValueGraphModeLine,
		common.BigValueTextModeValueAndName,
		common.VizOrientationHorizontal,
		utils.PrometheusQuery{
			Query:  `kube_pod_container_status_last_terminated_reason{reason="OOMKilled", pod=~"$pod", namespace=~"${namespace}"}`,
			Legend: "{{pod}}",
		},
	))

	return panelsArray
}

func resourcesUsage(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"CPU Usage",
		"",
		6,
		12,
		3,
		"",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{pod=~"$pod", namespace=~"${namespace}"}) by (pod)`,
			Legend: "{{pod}}",
		},
		utils.PrometheusQuery{
			Query:  `sum(kube_pod_container_resource_requests{job="kube-state-metrics", cluster="$cluster", namespace="$namespace", pod="$pod", resource="cpu"})`,
			Legend: "Requests",
		},
		utils.PrometheusQuery{
			Query:  `sum(kube_pod_container_resource_limits{job="kube-state-metrics", cluster="$cluster", namespace="$namespace", pod="$pod", resource="cpu"})`,
			Legend: "Limits",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Memory Usage",
		"",
		6,
		12,
		0,
		"bytes",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(container_memory_rss{pod=~"$pod", namespace=~"${namespace}", container!=""}) by (pod)`,
			Legend: "{{pod}}",
		},
		utils.PrometheusQuery{
			Query:  `sum(kube_pod_container_resource_requests{job="kube-state-metrics", cluster="$cluster", namespace="$namespace", pod="$pod", resource="memory"})`,
			Legend: "Requests",
		},
		utils.PrometheusQuery{
			Query:  `sum(kube_pod_container_resource_limits{job="kube-state-metrics", cluster="$cluster", namespace="$namespace", pod="$pod", resource="memory"})`,
			Legend: "Limits",
		},
	))

	return panelsArray
}

func networkUsage(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Receive Bandwidth",
		"",
		6,
		8,
		0,
		"bps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(irate(container_network_receive_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Transmit Bandwidth",
		"",
		6,
		8,
		0,
		"bps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum(irate(container_network_transmit_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Average Container Bandwidth by Namespace: Received",
		"",
		6,
		8,
		0,
		"bps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg(irate(container_network_receive_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"Average Container Bandwidth by Namespace: Transmitted",
		"",
		6,
		8,
		0,
		"bps",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `avg(irate(container_network_transmit_bytes_total{pod=~"$pod", namespace=~"${namespace}"}[$__rate_interval])) by (pod)`,
			Legend: "{{pod}}",
		},
	))

	return panelsArray
}

func diskUsage(p Props) []cog.Builder[dashboard.Panel] {
	var panelsArray []cog.Builder[dashboard.Panel]

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"IOPS(Read+Write)",
		"",
		6,
		12,
		2,
		"short",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `ceil(sum by(container, pod) (rate(container_fs_reads_total{job="kubelet", metrics_path="/metrics/cadvisor", container!="", cluster="$cluster", namespace="$namespace", pod="$pod"}[$__rate_interval]) + rate(container_fs_writes_total{job="kubelet", metrics_path="/metrics/cadvisor", container!="", cluster="$cluster", namespace="$namespace", pod="$pod"}[$__rate_interval])))`,
			Legend: "{{pod}}",
		},
	))

	panelsArray = append(panelsArray, utils.TimeSeriesPanel(
		p.MetricsDataSource,
		"ThroughPut(Read+Write)",
		"",
		6,
		12,
		2,
		"short",
		common.LegendPlacementBottom,
		utils.PrometheusQuery{
			Query:  `sum by(container, pod) (rate(container_fs_reads_bytes_total{job="kubelet", metrics_path="/metrics/cadvisor", container!="", cluster="$cluster", namespace="$namespace", pod="$pod"}[$__rate_interval]) + rate(container_fs_writes_bytes_total{job="kubelet", metrics_path="/metrics/cadvisor", container!="", cluster="$cluster", namespace="$namespace", pod="$pod"}[$__rate_interval]))`,
			Legend: "{{pod}}",
		},
	))

	return panelsArray
}
