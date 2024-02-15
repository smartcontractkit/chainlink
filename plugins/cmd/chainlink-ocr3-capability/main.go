package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	ocr3rp "github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const (
	loggerName = "PluginOCR3Capability"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := ocr3.NewOCR3(s.Logger)
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stop := make(chan struct{})
	defer close(stop)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: reportingplugins.ReportingPluginHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			ocr3rp.PluginServiceName: &ocr3rp.GRPCService[types.PluginProvider]{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					Logger:   s.Logger,
					StopCh:   stop,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}
