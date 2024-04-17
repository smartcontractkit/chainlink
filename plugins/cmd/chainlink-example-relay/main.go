package main

import (
	"context"
	"errors"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const (
	loggerName = "PluginExample"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := &pluginRelayer{lggr: s.Logger, ds: s.DataSource}
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stopCh,
					Logger:   s.Logger,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

type pluginRelayer struct {
	lggr logger.Logger
	ds   sqlutil.DataSource
}

func (p *pluginRelayer) Ready() error { return nil }

func (p *pluginRelayer) HealthReport() map[string]error { return map[string]error{p.Name(): nil} }

func (p *pluginRelayer) Name() string { return p.lggr.Name() }

func (p *pluginRelayer) NewRelayer(ctx context.Context, config string, keystore types.Keystore) (loop.Relayer, error) {
	var names []string
	err := p.ds.SelectContext(ctx, names, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	p.lggr.Info("Queried table names", "names", names)
	return nil, errors.New("example relayer unimplemented") //TODO sentinel error/message
}

func (p *pluginRelayer) Close() error { return nil }
