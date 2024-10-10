package main

import (
	"context"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/keystore"
)

const (
	loggerName = "PluginKeystore"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	stop := make(chan struct{})
	defer close(stop)

	kss := NewPlugin(s.Logger)

	s.MustRegister(kss)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginKeystoreHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginKeystoreName: &loop.GRPCPluginKeystore{
				BrokerConfig: loop.BrokerConfig{
					Logger:   s.Logger,
					StopCh:   stop,
					GRPCOpts: s.GRPCOpts,
				},
				PluginServer: kss,
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

func NewPlugin(lggr logger.Logger) *Plugin {
	return &Plugin{
		Plugin: loop.Plugin{Logger: lggr},
		stop:   make(services.StopChan),
	}
}

var _ keystore.Keystore = (*Plugin)(nil)
var _ keystore.Management = (*Plugin)(nil)

type Plugin struct {
	loop.Plugin
	stop services.StopChan
}

func (p *Plugin) AddPolicy(ctx context.Context, policy []byte) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) RemovePolicy(ctx context.Context, policyID string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) ListPolicy(ctx context.Context) []byte {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Import(ctx context.Context, keyType string, data []byte, tags []string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Export(ctx context.Context, keyID []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Create(ctx context.Context, keyType string, tags []string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Delete(ctx context.Context, keyID []byte) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) AddTag(ctx context.Context, keyID []byte, tag string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) RemoveTag(ctx context.Context, keyID []byte, tag string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) ListTags(ctx context.Context, keyID []byte) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Sign(ctx context.Context, keyID []byte, data []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) SignBatch(ctx context.Context, keyID []byte, data [][]byte) ([][]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Verify(ctx context.Context, keyID []byte, data []byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) VerifyBatch(ctx context.Context, keyID []byte, data [][]byte) ([]bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) List(ctx context.Context, tags []string) ([][]byte, error) {
	return [][]byte{[]byte("test 123")}, nil
}

func (p *Plugin) RunUDF(ctx context.Context, udfName string, keyID []byte, data []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
