package securemint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func TestTransmitter_NewTransmitter(t *testing.T) {
	lggr := logger.Test(t)

	// Create a mock capabilities registry
	mockRegistry := &mockCapabilitiesRegistry{}

	config := TransmitterConfig{
		Logger:                       lggr,
		CapabilitiesRegistry:         mockRegistry,
		DonID:                        1,
		TriggerCapabilityName:        "test-trigger",
		TriggerCapabilityVersion:     "1.0.0",
		TriggerTickerMinResolutionMs: 1000,
		TriggerSendChannelBufferSize: 1000,
	}

	transmitter, err := config.NewTransmitter("test-transmitter")
	require.NoError(t, err)
	require.NotNil(t, transmitter)

	// Verify it implements the required interfaces
	assert.Implements(t, (*Transmitter)(nil), transmitter)
	assert.Implements(t, (*capabilities.TriggerCapability)(nil), transmitter)

	// Test service lifecycle
	err = transmitter.Start(context.Background())
	require.NoError(t, err)

	err = transmitter.Close()
	require.NoError(t, err)
}

func TestTransmitter_RegisterTrigger(t *testing.T) {
	lggr := logger.Test(t)
	mockRegistry := &mockCapabilitiesRegistry{}

	config := TransmitterConfig{
		Logger:                       lggr,
		CapabilitiesRegistry:         mockRegistry,
		DonID:                        1,
		TriggerCapabilityName:        "test-trigger",
		TriggerCapabilityVersion:     "1.0.0",
		TriggerTickerMinResolutionMs: 1000,
		TriggerSendChannelBufferSize: 1000,
	}

	transmitter, err := config.NewTransmitter("test-transmitter")
	require.NoError(t, err)

	err = transmitter.Start(context.Background())
	require.NoError(t, err)
	defer transmitter.Close()

	// Create trigger config as values.Map
	triggerConfig, err := values.NewMap(map[string]any{
		"maxFrequencyMs": uint64(2000),
	})
	require.NoError(t, err)

	// Test trigger registration
	req := capabilities.TriggerRegistrationRequest{
		TriggerID: "test-trigger-1",
		Config:    triggerConfig,
		Metadata: capabilities.RequestMetadata{
			WorkflowID: "test-workflow",
		},
	}

	ch, err := transmitter.RegisterTrigger(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, ch)

	// Test duplicate registration
	_, err = transmitter.RegisterTrigger(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")

	// Test unregister
	err = transmitter.UnregisterTrigger(context.Background(), req)
	require.NoError(t, err)

	// Test unregister non-existent
	err = transmitter.UnregisterTrigger(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not registered")
}

func TestTransmitter_FromAccount(t *testing.T) {
	lggr := logger.Test(t)
	mockRegistry := &mockCapabilitiesRegistry{}

	config := TransmitterConfig{
		Logger:                       lggr,
		CapabilitiesRegistry:         mockRegistry,
		DonID:                        1,
		TriggerCapabilityName:        "test-trigger",
		TriggerCapabilityVersion:     "1.0.0",
		TriggerTickerMinResolutionMs: 1000,
		TriggerSendChannelBufferSize: 1000,
	}

	transmitter, err := config.NewTransmitter("0x11234")
	require.NoError(t, err)

	account, err := transmitter.FromAccount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// Verify account format includes logger name and don ID
	assert.Contains(t, string(account), lggr.Name())
	assert.Equal(t, "0x11234", string(account))
}

// Mock capabilities registry for testing
type mockCapabilitiesRegistry struct{}

func (m *mockCapabilitiesRegistry) Add(ctx context.Context, c capabilities.BaseCapability) error {
	return nil
}

func (m *mockCapabilitiesRegistry) Remove(ctx context.Context, ID string) error {
	return nil
}

func (m *mockCapabilitiesRegistry) Get(ctx context.Context, ID string) (capabilities.BaseCapability, error) {
	return nil, nil
}

func (m *mockCapabilitiesRegistry) List(ctx context.Context) ([]capabilities.BaseCapability, error) {
	return nil, nil
}

func (m *mockCapabilitiesRegistry) GetExecutable(ctx context.Context, ID string) (capabilities.ExecutableCapability, error) {
	return nil, nil
}

func (m *mockCapabilitiesRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	return capabilities.CapabilityConfiguration{}, nil
}

func (m *mockCapabilitiesRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	return capabilities.Node{}, nil
}

func (m *mockCapabilitiesRegistry) GetTrigger(ctx context.Context, ID string) (capabilities.TriggerCapability, error) {
	return nil, nil
}

func (m *mockCapabilitiesRegistry) NodeByPeerID(ctx context.Context, peerID p2ptypes.PeerID) (capabilities.Node, error) {
	return capabilities.Node{}, nil
}
