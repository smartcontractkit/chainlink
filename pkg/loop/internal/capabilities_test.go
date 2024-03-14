package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type mockTrigger struct {
	capabilities.BaseCapability
	callback chan<- capabilities.CapabilityResponse
}

func (m *mockTrigger) RegisterTrigger(ctx context.Context, callback chan<- capabilities.CapabilityResponse, request capabilities.CapabilityRequest) error {
	m.callback = callback
	return nil
}

func (m *mockTrigger) UnregisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) error {
	m.callback = nil
	return nil
}

func mustMockTrigger(t *testing.T) *mockTrigger {
	return &mockTrigger{
		BaseCapability: capabilities.MustNewCapabilityInfo("trigger", capabilities.CapabilityTypeTrigger, "a mock trigger", "v0.0.1"),
	}
}

type mockCallback struct {
	capabilities.BaseCapability
	callback     chan<- capabilities.CapabilityResponse
	regRequest   capabilities.RegisterToWorkflowRequest
	unregRequest capabilities.UnregisterFromWorkflowRequest
}

func (m *mockCallback) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	m.regRequest = request
	return nil
}

func (m *mockCallback) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	m.unregRequest = request
	return nil
}

func (m *mockCallback) Execute(ctx context.Context, callback chan<- capabilities.CapabilityResponse, request capabilities.CapabilityRequest) error {
	m.callback = callback
	return nil
}

func mustMockCallback(t *testing.T, _type capabilities.CapabilityType) *mockCallback {
	return &mockCallback{
		BaseCapability: capabilities.MustNewCapabilityInfo(fmt.Sprintf("callback %s", _type), _type, fmt.Sprintf("a mock %s", _type), "v0.0.1"),
	}
}

type capabilityPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	brokerCfg  BrokerConfig
	capability capabilities.BaseCapability
}

func (c *capabilityPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, client *grpc.ClientConn) (any, error) {
	bext := &BrokerExt{
		BrokerConfig: c.brokerCfg,
		Broker:       broker,
	}
	switch c.capability.(type) {
	case capabilities.TriggerExecutable:
		return NewTriggerCapabilityClient(bext, client), nil
	case capabilities.CallbackExecutable:
		return NewCallbackCapabilityClient(bext, client), nil
	}

	panic("unreachable")
}

func (c *capabilityPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	switch tc := c.capability.(type) {
	case capabilities.TriggerCapability:
		return RegisterTriggerCapabilityServer(server, broker, c.brokerCfg, tc)
	case CallbackCapability:
		return RegisterCallbackCapabilityServer(server, broker, c.brokerCfg, tc)
	}

	return nil
}

func newCapabilityPlugin(t *testing.T, capability capabilities.BaseCapability) (capabilities.BaseCapability, error) {
	stopCh := make(chan struct{})
	logger := logger.Test(t)
	pluginName := "registry"

	client, _ := plugin.TestPluginGRPCConn(
		t,
		false,
		map[string]plugin.Plugin{
			pluginName: &capabilityPlugin{
				brokerCfg: BrokerConfig{
					StopCh: stopCh,
					Logger: logger,
				},
				capability: capability,
			},
		},
	)

	regClient, err := client.Dispense(pluginName)
	require.NoError(t, err)

	return regClient.(capabilities.BaseCapability), nil
}

func Test_Capabilities(t *testing.T) {
	mtr := mustMockTrigger(t)
	ma := mustMockCallback(t, capabilities.CapabilityTypeAction)
	ctx := tests.Context(t)

	t.Run("fetching a trigger capability, and executing it", func(t *testing.T) {
		tr, err := newCapabilityPlugin(t, mtr)
		require.NoError(t, err)

		ctr := tr.(capabilities.TriggerCapability)

		ch := make(chan capabilities.CapabilityResponse)
		err = ctr.RegisterTrigger(
			ctx,
			ch,
			capabilities.CapabilityRequest{})
		require.NoError(t, err)

		vs := values.NewString("hello")
		require.NoError(t, err)
		cr := capabilities.CapabilityResponse{
			Value: vs,
		}
		mtr.callback <- cr
		assert.Equal(t, cr, <-ch)
	})

	t.Run("fetching a trigger capability, and closing the channel", func(t *testing.T) {
		tr, err := newCapabilityPlugin(t, mtr)
		require.NoError(t, err)

		ctr := tr.(capabilities.TriggerCapability)

		ch := make(chan capabilities.CapabilityResponse)
		err = ctr.RegisterTrigger(
			ctx,
			ch,
			capabilities.CapabilityRequest{})
		require.NoError(t, err)

		// Close the channel from the server, to signal no further results.
		close(mtr.callback)

		// This should propagate to the client.
		_, isOpen := <-ch
		assert.False(t, isOpen)
	})

	t.Run("fetching a trigger capability, and unregistering", func(t *testing.T) {
		tr, err := newCapabilityPlugin(t, mtr)
		require.NoError(t, err)

		ctr := tr.(capabilities.TriggerCapability)

		ch := make(chan capabilities.CapabilityResponse)
		err = ctr.RegisterTrigger(
			ctx,
			ch,
			capabilities.CapabilityRequest{})
		require.NoError(t, err)
		assert.NotNil(t, mtr.callback)

		err = ctr.UnregisterTrigger(
			ctx,
			capabilities.CapabilityRequest{})
		require.NoError(t, err)

		assert.Nil(t, mtr.callback)
	})

	t.Run("fetching a trigger capability and calling Info", func(t *testing.T) {
		tr, err := newCapabilityPlugin(t, mtr)
		require.NoError(t, err)

		gotInfo, err := tr.Info(ctx)
		require.NoError(t, err)

		expectedInfo, err := mtr.Info(ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedInfo, gotInfo)
	})

	t.Run("fetching an action capability, and (un)registering it", func(t *testing.T) {
		c, err := newCapabilityPlugin(t, ma)
		require.NoError(t, err)

		act := c.(capabilities.ActionCapability)

		vmap, err := values.NewMap(map[string]any{"foo": "bar"})
		require.NoError(t, err)
		expectedRequest := capabilities.RegisterToWorkflowRequest{
			Config: vmap,
		}
		err = act.RegisterToWorkflow(
			ctx,
			expectedRequest)
		require.NoError(t, err)

		assert.Equal(t, expectedRequest, ma.regRequest)

		expectedUnrRequest := capabilities.UnregisterFromWorkflowRequest{
			Config: vmap,
		}
		err = act.UnregisterFromWorkflow(
			ctx,
			expectedUnrRequest)
		require.NoError(t, err)

		assert.Equal(t, expectedUnrRequest, ma.unregRequest)
	})

	t.Run("fetching an action capability, and executing it", func(t *testing.T) {
		c, err := newCapabilityPlugin(t, ma)
		require.NoError(t, err)

		cmap, err := values.NewMap(map[string]any{"foo": "bar"})
		require.NoError(t, err)

		imap, err := values.NewMap(map[string]any{"bar": "baz"})
		require.NoError(t, err)
		expectedRequest := capabilities.CapabilityRequest{
			Config: cmap,
			Inputs: imap,
		}
		ch := make(chan capabilities.CapabilityResponse)
		err = c.(capabilities.ActionCapability).Execute(
			ctx,
			ch,
			expectedRequest)
		require.NoError(t, err)

		expectedErr := errors.New("an error")
		expectedResp := capabilities.CapabilityResponse{
			Err: expectedErr,
		}

		ma.callback <- expectedResp
		assert.Equal(t, expectedResp, <-ch)
	})

	t.Run("fetching an action capability, and closing it", func(t *testing.T) {
		c, err := newCapabilityPlugin(t, ma)
		require.NoError(t, err)

		cmap, err := values.NewMap(map[string]any{"foo": "bar"})
		require.NoError(t, err)

		imap, err := values.NewMap(map[string]any{"bar": "baz"})
		require.NoError(t, err)
		expectedRequest := capabilities.CapabilityRequest{
			Config: cmap,
			Inputs: imap,
		}
		ch := make(chan capabilities.CapabilityResponse)
		err = c.(capabilities.ActionCapability).Execute(
			ctx,
			ch,
			expectedRequest)
		require.NoError(t, err)

		close(ma.callback)
		_, isOpen := <-ch
		assert.False(t, isOpen)
	})
}
