package securemint

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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

func TestTransmitter_Transmit(t *testing.T) {
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

	// Register a trigger to receive events
	triggerConfig, err := values.NewMap(map[string]any{
		"maxFrequencyMs": uint64(2000),
	})
	require.NoError(t, err)

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

	t.Run("successful transmission", func(t *testing.T) {
		// Create test data
		cd := ocr2types.ConfigDigest{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		seqNr := uint64(123)
		report := ocr3types.ReportWithInfo[securemint.ChainSelector]{
			Report: []byte("test report data"),
			Info:   securemint.ChainSelector(1), // Use ChainSelector as the Info type
		}
		sigs := []types.AttributedOnchainSignature{
			{
				Signature: []byte("signature1"),
				Signer:    1,
			},
			{
				Signature: []byte("signature2"),
				Signer:    2,
			},
		}

		// Transmit the report
		err := transmitter.Transmit(context.Background(), cd, seqNr, report, sigs)
		require.NoError(t, err)

		// Wait for the event to be processed and sent to the channel
		select {
		case response := <-ch:
			// Verify the trigger response
			assert.Equal(t, "securemint_123", response.Event.ID)
			assert.Equal(t, transmitter.CapabilityInfo.ID, response.Event.TriggerType)

			// Verify outputs
			outputs := response.Event.Outputs
			assert.NotNil(t, outputs)

			// Check seqNr
			var seqNrVal uint64
			err = outputs.Underlying["seqNr"].UnwrapTo(&seqNrVal)
			require.NoError(t, err)
			assert.Equal(t, uint64(123), seqNrVal)

			// Check configDigest
			var cdVal ocr2types.ConfigDigest
			err = outputs.Underlying["configDigest"].UnwrapTo(&cdVal)
			require.NoError(t, err)
			assert.Equal(t, cd, cdVal)

			// Check signatures
			var capSigs []capabilities.OCRAttributedOnchainSignature
			err = outputs.Underlying["sigs"].UnwrapTo(&capSigs)
			require.NoError(t, err)
			assert.Len(t, capSigs, 2)
			assert.Equal(t, uint32(1), capSigs[0].Signer)
			assert.Equal(t, []byte("signature1"), capSigs[0].Signature)
			assert.Equal(t, uint32(2), capSigs[1].Signer)
			assert.Equal(t, []byte("signature2"), capSigs[1].Signature)

			// Check report
			var reportBytes []byte
			err = outputs.Underlying["report"].UnwrapTo(&reportBytes)
			require.NoError(t, err)

			// json umarshal bytes to string and check if it contains "test report data"
			var report ocr3types.ReportWithInfo[securemint.ChainSelector]
			err = json.Unmarshal(reportBytes, &report)
			require.NoError(t, err)
			assert.Equal(t, "test report data", string(report.Report))

		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for trigger response")
		}
	})
}

func TestTransmitter_Transmit_NoSubscribers(t *testing.T) {
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

	// Test transmission without any registered triggers
	cd := ocr2types.ConfigDigest{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	seqNr := uint64(127)
	report := ocr3types.ReportWithInfo[securemint.ChainSelector]{
		Report: []byte("test report data no subscribers"),
		Info:   securemint.ChainSelector(1), // Use ChainSelector as the Info type
	}
	sigs := []types.AttributedOnchainSignature{}

	// This should succeed even without subscribers
	err = transmitter.Transmit(context.Background(), cd, seqNr, report, sigs)
	require.NoError(t, err)
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

func (m *mockCapabilitiesRegistry) DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error) {
	return nil, nil
}
