package relayerset

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/relayerset"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_RelayerSet(t *testing.T) {
	ctx := tests.Context(t)
	stopCh := make(chan struct{})
	log := logger.Test(t)

	relayer1 := mocks.NewRelayer(t)
	relayer2 := mocks.NewRelayer(t)

	relayers := map[types.RelayID]core.Relayer{
		{
			Network: "N1",
			ChainID: "C1",
		}: relayer1,
		{
			Network: "N2",
			ChainID: "C2",
		}: relayer2,
	}

	pluginName := "relayerset-test"
	client, server := plugin.TestPluginGRPCConn(
		t,
		true,
		map[string]plugin.Plugin{
			pluginName: &testRelaySetPlugin{
				log:  log,
				impl: &TestRelayerSet{relayers: relayers},
				brokerExt: &net.BrokerExt{
					BrokerConfig: net.BrokerConfig{
						StopCh: stopCh,
						Logger: log,
					},
				},
			},
		},
	)

	defer client.Close()
	defer server.Stop()

	relayerSetClient, err := client.Dispense(pluginName)
	require.NoError(t, err)

	rc, ok := relayerSetClient.(*Client)
	require.True(t, ok)

	relayerClient, err := rc.Get(ctx, types.RelayID{
		Network: "N1",
		ChainID: "C1",
	})

	require.NoError(t, err)

	relayer1.On("Start", mock.Anything).Return(nil)
	err = relayerClient.Start(ctx)
	require.NoError(t, err)
	relayer1.AssertCalled(t, "Start", mock.Anything)

	relayer1.On("Close").Return(nil)
	err = relayerClient.Close()
	require.NoError(t, err)
	relayer1.AssertCalled(t, "Close")

	relayer1.On("Ready").Return(nil)
	err = relayerClient.Ready()
	require.NoError(t, err)
	relayer1.AssertCalled(t, "Ready")

	relayer1.On("HealthReport").Return(map[string]error{"stat1": errors.New("error1")})
	healthReport := relayerClient.HealthReport()
	require.Len(t, healthReport, 1)
	require.Equal(t, "error1", healthReport["stat1"].Error())
	relayer1.AssertCalled(t, "HealthReport")

	relayer1.On("Name").Return("test-relayer")
	name := relayerClient.Name()
	require.Equal(t, "test-relayer", name)
	relayer1.AssertCalled(t, "Name")
}

type TestRelayerSet struct {
	relayers map[types.RelayID]core.Relayer
}

func (t *TestRelayerSet) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	if relayer, ok := t.relayers[relayID]; ok {
		return relayer, nil
	}

	return nil, fmt.Errorf("relayer with id %s not found", relayID)
}

func (t *TestRelayerSet) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	return t.relayers, nil
}

type testRelaySetPlugin struct {
	log logger.Logger
	plugin.NetRPCUnsupportedPlugin
	brokerExt *net.BrokerExt
	impl      core.RelayerSet
}

func (r *testRelaySetPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, client *grpc.ClientConn) (any, error) {
	r.brokerExt.Broker = broker

	return NewRelayerSetClient(logger.Nop(), r.brokerExt, client), nil
}

func (r *testRelaySetPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	r.brokerExt.Broker = broker

	rs, _ := NewRelayerSetServer(r.log, r.impl, r.brokerExt)
	relayerset.RegisterRelayerSetServer(server, rs)
	return nil
}
