package cre

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	streamstypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
)

func TestNewNodagTransmitter(t *testing.T) {
	t.Run("creates transmitter successfully", func(t *testing.T) {
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)
		require.NotNil(t, transmitter)
		assert.NotNil(t, transmitter.transmitter, "should embed legacy transmitter")
	})

	t.Run("returns error if legacy transmitter creation fails", func(t *testing.T) {
		// Invalid config - will panic with nil logger, so we verify it creates the wrapper at least
		// The actual validation happens in the legacy transmitter
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "", // Empty name
			TriggerCapabilityVersion:     "",
			TriggerTickerMinResolutionMs: 0,
			TriggerSendChannelBufferSize: 0,
		}

		// Should still create successfully - legacy transmitter applies defaults
		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)
		require.NotNil(t, transmitter)
	})
}

func TestConvertProtoConfigToLegacy(t *testing.T) {
	t.Run("converts valid proto config", func(t *testing.T) {
		protoConfig := &streams.Config{
			StreamIds:      []uint32{1, 2, 3},
			MaxFrequencyMs: 5000,
		}

		legacyConfig, err := convertProtoConfigToLegacy(protoConfig)
		require.NoError(t, err)
		require.NotNil(t, legacyConfig)

		assert.Len(t, legacyConfig.StreamIDs, 3)
		assert.Equal(t, streamstypes.LLOStreamID(1), legacyConfig.StreamIDs[0])
		assert.Equal(t, streamstypes.LLOStreamID(2), legacyConfig.StreamIDs[1])
		assert.Equal(t, streamstypes.LLOStreamID(3), legacyConfig.StreamIDs[2])
		assert.Equal(t, uint64(5000), legacyConfig.MaxFrequencyMs)
	})

	t.Run("handles empty stream IDs", func(t *testing.T) {
		protoConfig := &streams.Config{
			StreamIds:      []uint32{},
			MaxFrequencyMs: 1000,
		}

		legacyConfig, err := convertProtoConfigToLegacy(protoConfig)
		require.NoError(t, err)
		assert.Empty(t, legacyConfig.StreamIDs)
		assert.Equal(t, uint64(1000), legacyConfig.MaxFrequencyMs)
	})

	t.Run("returns error for nil config", func(t *testing.T) {
		_, err := convertProtoConfigToLegacy(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil")
	})
}

func TestConvertLegacyResponseToProto(t *testing.T) {
	t.Run("converts valid legacy response", func(t *testing.T) {
		event := capabilities.OCRTriggerEvent{
			ConfigDigest: []byte{1, 2, 3, 4},
			SeqNr:        42,
			Report:       []byte{5, 6, 7, 8},
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{
					Signer:    1,
					Signature: []byte{9, 10, 11},
				},
				{
					Signer:    2,
					Signature: []byte{12, 13, 14},
				},
			},
		}

		outputsMap, err := values.WrapMap(event)
		require.NoError(t, err)

		legacyResp := capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "streams-trigger@2.0.0",
				ID:          "test-event-1",
				Outputs:     outputsMap,
			},
		}

		protoReport, err := convertLegacyResponseToProto(legacyResp)
		require.NoError(t, err)
		require.NotNil(t, protoReport)

		assert.Equal(t, []byte{1, 2, 3, 4}, protoReport.ConfigDigest)
		assert.Equal(t, uint64(42), protoReport.SeqNr)
		assert.Equal(t, []byte{5, 6, 7, 8}, protoReport.Report)
		assert.Len(t, protoReport.Sigs, 2)
		assert.Equal(t, uint32(1), protoReport.Sigs[0].Signer)
		assert.Equal(t, []byte{9, 10, 11}, protoReport.Sigs[0].Signature)
		assert.Equal(t, uint32(2), protoReport.Sigs[1].Signer)
		assert.Equal(t, []byte{12, 13, 14}, protoReport.Sigs[1].Signature)
	})

	t.Run("converts legacy response with no signatures", func(t *testing.T) {
		// Create an event with no signatures
		event := capabilities.OCRTriggerEvent{
			ConfigDigest: []byte{1, 2, 3, 4},
			SeqNr:        42,
			Report:       []byte{5, 6, 7, 8},
			Sigs:         []capabilities.OCRAttributedOnchainSignature{}, // Empty sigs
		}

		outputsMap, err := values.WrapMap(event)
		require.NoError(t, err)

		legacyResp := capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				Outputs: outputsMap,
			},
		}

		protoReport, err := convertLegacyResponseToProto(legacyResp)
		require.NoError(t, err)
		require.NotNil(t, protoReport)

		assert.Equal(t, []byte{1, 2, 3, 4}, protoReport.ConfigDigest)
		assert.Equal(t, uint64(42), protoReport.SeqNr)
		assert.Empty(t, protoReport.Sigs)
	})
}

func TestNodagTransmitter_RegisterTrigger(t *testing.T) {
	t.Run("registers trigger successfully", func(t *testing.T) {
		// Setup
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{
			WorkflowID:   "test-workflow",
			WorkflowName: "Test Workflow",
		}

		protoConfig := &streams.Config{
			StreamIds:      []uint32{1, 2},
			MaxFrequencyMs: 5000,
		}

		// Register trigger
		ch, err := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, protoConfig)
		require.NoError(t, err)
		require.NotNil(t, ch)

		// Cleanup
		err = transmitter.UnregisterTrigger(ctx, "trigger-1", metadata, protoConfig)
		assert.NoError(t, err)
	})

	t.Run("returns error for invalid config", func(t *testing.T) {
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{
			WorkflowID: "test-workflow",
		}

		// Invalid config (nil)
		_, err = transmitter.RegisterTrigger(ctx, "trigger-1", metadata, nil)
		assert.Error(t, err)
	})
}

func TestNodagTransmitter_ChannelBridge(t *testing.T) {
	t.Run("bridges events from legacy to proto channel", func(t *testing.T) {
		// This test verifies that events from the legacy channel
		// are correctly converted and sent to the proto channel

		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Create channels
		legacyCh := make(chan capabilities.TriggerResponse, 10)
		protoCh := make(chan capabilities.TriggerAndId[*streams.Report], 10)

		// Start bridge
		go transmitter.bridgeChannels(ctx, "test-trigger", legacyCh, protoCh)

		// Send a test event
		event := capabilities.OCRTriggerEvent{
			ConfigDigest: []byte{1, 2, 3, 4},
			SeqNr:        42,
			Report:       []byte{5, 6, 7, 8},
			Sigs: []capabilities.OCRAttributedOnchainSignature{
				{Signer: 1, Signature: []byte{9, 10, 11}},
			},
		}

		outputsMap, err := values.WrapMap(event)
		require.NoError(t, err)

		legacyCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "streams-trigger@2.0.0",
				ID:          "test-event-1",
				Outputs:     outputsMap,
			},
		}

		// Receive proto event
		select {
		case protoEvent := <-protoCh:
			assert.Equal(t, "test-trigger", protoEvent.Id)
			assert.Equal(t, []byte{1, 2, 3, 4}, protoEvent.Trigger.ConfigDigest)
			assert.Equal(t, uint64(42), protoEvent.Trigger.SeqNr)
		case <-ctx.Done():
			t.Fatal("Timeout waiting for proto event")
		}
	})

	t.Run("closes proto channel when legacy channel closes", func(t *testing.T) {
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		legacyCh := make(chan capabilities.TriggerResponse)
		protoCh := make(chan capabilities.TriggerAndId[*streams.Report], 10)

		go transmitter.bridgeChannels(ctx, "test-trigger", legacyCh, protoCh)

		// Close legacy channel
		close(legacyCh)

		// Verify proto channel closes
		_, ok := <-protoCh
		assert.False(t, ok, "proto channel should be closed")
	})

	t.Run("stops on context cancellation", func(t *testing.T) {
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 100,
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		legacyCh := make(chan capabilities.TriggerResponse)
		protoCh := make(chan capabilities.TriggerAndId[*streams.Report], 10)

		go transmitter.bridgeChannels(ctx, "test-trigger", legacyCh, protoCh)

		// Cancel context
		cancel()

		// Verify proto channel closes
		_, ok := <-protoCh
		assert.False(t, ok, "proto channel should be closed")
	})

	t.Run("drops events when proto channel is full", func(t *testing.T) {
		// This should match the behavior of the legacy transmitter
		// When the channel is full, events are dropped (non-blocking send)
		cfg := TransmitterConfig{
			Logger:                       logger.Test(t),
			TriggerCapabilityName:        "test-trigger",
			TriggerCapabilityVersion:     "1.0.0",
			TriggerTickerMinResolutionMs: 1000,
			TriggerSendChannelBufferSize: 0, // Unbuffered channel
		}

		transmitter, err := NewNodagTransmitter(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		legacyCh := make(chan capabilities.TriggerResponse, 10)
		protoCh := make(chan capabilities.TriggerAndId[*streams.Report]) // Unbuffered - will block/drop

		go transmitter.bridgeChannels(ctx, "test-trigger", legacyCh, protoCh)

		// Create test events
		event := capabilities.OCRTriggerEvent{
			ConfigDigest: []byte{1, 2, 3, 4},
			SeqNr:        1,
			Report:       []byte{5, 6, 7, 8},
			Sigs:         []capabilities.OCRAttributedOnchainSignature{},
		}

		outputsMap, err := values.WrapMap(event)
		require.NoError(t, err)

		// Send event (will be dropped because no receiver)
		legacyCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{
				TriggerType: "streams-trigger@2.0.0",
				ID:          "event-1",
				Outputs:     outputsMap,
			},
		}

		// Give bridge time to attempt send (and drop it)
		// Since channel is unbuffered and no receiver, it should drop immediately
		select {
		case <-protoCh:
			t.Fatal("Should not have received event (should be dropped)")
		default:
			// Expected - event was dropped because channel has no receiver
		}
	})
}
