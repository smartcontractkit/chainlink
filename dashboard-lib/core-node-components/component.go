package core_node_components

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
)

type Props struct {
	PrometheusDataSource string
	PlatformOpts         PlatformOpts
}

func vars(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.VariableAsInterval(
			"interval",
			interval.Values([]string{"30s", "1m", "5m", "15m", "30m", "1h", "6h", "12h"}),
		),
		dashboard.VariableAsQuery(
			"env",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up, env)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"cluster",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\"}, cluster)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"blockchain",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\", cluster=\"$cluster\"}, blockchain)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"product",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\", cluster=\"$cluster\", blockchain=\"$blockchain\"}, product)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"network_type",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\", cluster=\"$cluster\", blockchain=\"$blockchain\", product=\"$product\"}, network_type)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"component",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\", cluster=\"$cluster\", blockchain=\"$blockchain\", network_type=\"$network_type\"}, component)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"service",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(up{env=\"$env\", cluster=\"$cluster\", blockchain=\"$blockchain\", network_type=\"$network_type\", component=\"$component\"}, service)"),
			query.Sort(query.NumericalAsc),
		),
		dashboard.VariableAsQuery(
			"service_id",
			query.DataSource(p.PrometheusDataSource),
			query.Multiple(),
			query.IncludeAll(),
			query.Request("label_values(health{env=\"$env\", cluster=\"$cluster\", blockchain=\"$blockchain\", network_type=\"$network_type\", component=\"$component\", service=\"$service\"}, service_id)"),
			query.Sort(query.NumericalAsc),
		),
	}
}

func generalInfoRow(p Props) []dashboard.Option {
	return []dashboard.Option{
		dashboard.Row(
			"General CL Cluster Info",
			// row.Collapse(),
			row.WithTable(
				"List Nodes",
				table.Span(12),
				table.HideColumn("Time"),
				table.HideColumn("Value"),
				table.DataSource(p.PrometheusDataSource),
				table.WithPrometheusTarget(
					`max(up{`+p.PlatformOpts.LabelQuery+`}) by (env, cluster, blockchain, product, network_type, network, version, team, component, service)`,
					prometheus.Legend(""),
					prometheus.Format("table"),
					prometheus.Instant(),
				),
			),
			row.WithTimeSeries(
				"Uptime",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Axis(
					axis.Max(1),
					axis.Max(0),
					axis.Unit("bool"),
					axis.Label("Alive"),
				),
				timeseries.WithPrometheusTarget(
					`up{`+p.PlatformOpts.LabelQuery+`}`,
					prometheus.Legend(""),
				),
			),
			row.WithTimeSeries(
				"Service Components Health by Service",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Legend(timeseries.ToTheRight),
				timeseries.WithPrometheusTarget(
					`health{`+p.PlatformOpts.LabelQuery+`service_id=~"${service_id}"}`,
					prometheus.Legend("{{service_id}}"),
				),
			),
			row.WithTimeSeries(
				"Service Components Health Avg by Service",
				timeseries.Span(12),
				timeseries.Height("200px"),
				timeseries.DataSource(p.PrometheusDataSource),
				timeseries.Legend(timeseries.ToTheRight),
				timeseries.WithPrometheusTarget(
					`avg(avg_over_time(health{`+p.PlatformOpts.LabelQuery+`service_id=~"${service_id}"}[$interval])) by (service_id, version, service, cluster, env)`,
					prometheus.Legend("{{service_id}}"),
				),
			),
			row.WithStat(
				"Service Components Health Avg by Service",
				stat.Span(12),
				stat.Height("200px"),
				stat.DataSource(p.PrometheusDataSource),
				stat.Text(stat.TextValueAndName),
				stat.Orientation(stat.OrientationVertical),
				stat.SparkLine(),
				stat.TitleFontSize(4),
				stat.ValueFontSize(12),
				stat.WithPrometheusTarget(
					`avg(avg_over_time(health{`+p.PlatformOpts.LabelQuery+`service_id=~"${service_id}"}[$interval])) by (service_id, version, service, cluster, env)`,
					prometheus.Legend("{{service_id}}"),
				),
			),
		),
	}
}

func New(p Props) []dashboard.Option {
	opts := vars(p)
	opts = append(opts, generalInfoRow(p)...)
	return opts
}
