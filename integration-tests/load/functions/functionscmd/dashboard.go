package main

import (
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/grabana/row"
	"github.com/smartcontractkit/wasp"
)

func main() {
	lokiDS := "grafanacloud-logs"
	d, err := wasp.NewDashboard(nil,
		[]dashboard.Option{
			dashboard.Row("DON logs (errors)",
				row.Collapse(),
				row.WithLogs(
					"DON logs",
					logs.DataSource(lokiDS),
					logs.Span(12),
					logs.Height("300px"),
					logs.Transparent(),
					logs.WithLokiTarget(`
					{ cluster="staging-us-west-2-main", app=~"clc-ocr2-dr-matic-testnet" } | json | level="error"
				`),
				)),
		},
	)
	if err != nil {
		panic(err)
	}
	if _, err := d.Deploy(); err != nil {
		panic(err)
	}
}
