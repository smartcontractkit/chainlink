package gateway

import (
	"encoding/json"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/monitoring"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type capturingHandlerFactory struct {
	t        *testing.T
	handlers map[HandlerType]handlers.Handler
	calls    map[HandlerType][]config.ShardedDONConfig
}

func (f *capturingHandlerFactory) NewHandler(handlerType HandlerType, _ json.RawMessage, shardedDONs []config.ShardedDONConfig, _ [][]handlers.DON) (handlers.Handler, error) {
	f.calls[handlerType] = shardedDONs
	h, ok := f.handlers[handlerType]
	require.True(f.t, ok, "missing test handler for type %q", handlerType)
	return h, nil
}

func TestSetupFromNewConfig_SharedDONAndLegacyMethodRouting(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)

	cfg := &config.GatewayConfig{
		NodeServerConfig: gw_net.WebSocketServerConfig{
			HTTPServerConfig: gw_net.HTTPServerConfig{Path: "/node"},
		},
		ShardedDONs: []config.ShardedDONConfig{
			{
				DonName: "shared-don",
				F:       0,
				Shards: []config.Shard{
					{
						Nodes: []config.NodeConfig{
							{
								Name:    "node-1",
								Address: "0x68902D681C28119F9B2531473A417088BF008E59",
							},
						},
					},
				},
			},
		},
		Services: []config.ServiceConfig{
			{
				ServiceName: "workflows",
				DONs:        []string{"shared-don"},
				Handlers: []config.Handler{
					{Name: "workflows-handler"},
				},
			},
			{
				ServiceName: "vault",
				DONs:        []string{"shared-don"},
				Handlers: []config.Handler{
					{Name: "vault-handler"},
				},
			},
		},
	}

	connMgr, err := NewConnectionManager(cfg, clockwork.NewFakeClock(), gMetrics, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	workflowsHandler := handlermocks.NewHandler(t)
	workflowsHandler.On("Methods").Return([]string{"workflows.execute"})
	vaultHandler := handlermocks.NewHandler(t)
	vaultHandler.On("Methods").Return([]string{"vault.store"})

	factory := &capturingHandlerFactory{
		t: t,
		handlers: map[HandlerType]handlers.Handler{
			"workflows-handler": workflowsHandler,
			"vault-handler":     vaultHandler,
		},
		calls: make(map[HandlerType][]config.ShardedDONConfig),
	}

	serviceToHandler, err := setupFromNewConfig(cfg, factory, connMgr, lggr)
	require.NoError(t, err)
	require.Len(t, serviceToHandler, 2, "shared DON should be attachable to multiple services")
	require.NotNil(t, serviceToHandler["workflows"])
	require.NotNil(t, serviceToHandler["vault"])

	require.Equal(t,
		"0x68902d681c28119f9b2531473a417088bf008e59",
		factory.calls["workflows-handler"][0].Shards[0].Nodes[0].Address,
		"setupFromNewConfig should normalize node addresses before passing DON configs to handlers",
	)
	require.Equal(t,
		"0x68902d681c28119f9b2531473a417088bf008e59",
		factory.calls["vault-handler"][0].Shards[0].Nodes[0].Address,
		"setupFromNewConfig should normalize node addresses before passing DON configs to handlers",
	)

	donConnMgr := connMgr.DONConnectionManager(config.ShardDONID("shared-don", 0))
	require.NotNil(t, donConnMgr)

	legacyHandler, err := donConnMgr.getHandler("legacyMethodWithoutServicePrefix")
	require.NoError(t, err)
	require.Same(t, serviceToHandler["workflows"], legacyHandler, "legacy methods should route to workflows service")

	vaultMethodHandler, err := donConnMgr.getHandler("vault.store")
	require.NoError(t, err)
	require.Same(t, serviceToHandler["vault"], vaultMethodHandler, "service-prefixed methods should route by prefix")
}
