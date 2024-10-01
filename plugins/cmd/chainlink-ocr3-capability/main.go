package main

import (
	"context"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins"
	ocr3rp "github.com/smartcontractkit/chainlink-common/pkg/loop/reportingplugins/ocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
)

const (
	loggerName = "PluginOCR3Capability"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	c := ocr3.Config{
		Logger:            s.Logger,
		EncoderFactory:    capabilities.NewEncoder,
		AggregatorFactory: capabilities.NewAggregator,
	}
	p := ocr3.NewOCR3(c)
	if err := p.Start(context.Background()); err != nil {
		s.Logger.Fatal("Failed to start OCR3 capability", err)
	}

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
