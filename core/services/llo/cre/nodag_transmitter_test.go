package cre

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	rageptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// mockCapabilityRegistry for testing
type mockCapabilityRegistry struct {
	added   []capabilities.BaseCapability
	removed []string
}

func (m *mockCapabilityRegistry) Add(ctx context.Context, capability capabilities.BaseCapability) error {
	m.added = append(m.added, capability)
	return nil
}

func (m *mockCapabilityRegistry) Remove(ctx context.Context, id string) error {
	m.removed = append(m.removed, id)
	return nil
}

func (m *mockCapabilityRegistry) Get(ctx context.Context, id string) (capabilities.BaseCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) GetTrigger(ctx context.Context, id string) (capabilities.TriggerCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) GetAction(ctx context.Context, id string) (capabilities.ActionCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) GetExecutable(ctx context.Context, id string) (capabilities.ExecutableCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) GetConsensus(ctx context.Context, id string) (capabilities.ConsensusCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) GetTarget(ctx context.Context, id string) (capabilities.TargetCapability, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) List(ctx context.Context) ([]capabilities.BaseCapability, error) {
	return m.added, nil
}

func (m *mockCapabilityRegistry) ConfigForCapability(ctx context.Context, capabilityID string, capabilityDonID uint32) (capabilities.CapabilityConfiguration, error) {
	return capabilities.CapabilityConfiguration{}, nil
}

func (m *mockCapabilityRegistry) DONsForCapability(ctx context.Context, id string) ([]capabilities.DONWithNodes, error) {
	return nil, nil
}

func (m *mockCapabilityRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	return capabilities.Node{}, nil
}

func (m *mockCapabilityRegistry) NodeByPeerID(ctx context.Context, peerID rageptypes.PeerID) (capabilities.Node, error) {
	return capabilities.Node{}, nil
}

func TestNewNodagTransmitter(t *testing.T) {
	t.Run("creates transmitter successfully", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		lggr := logger.Test(t)

		transmitter, err := NewNodagTransmitter(lggr, 1, registry)
		require.NoError(t, err)
		require.NotNil(t, transmitter)
		assert.Equal(t, nodagCapabilityID, transmitter.ID)
		assert.Empty(t, transmitter.subscribers)
	})
}

func TestNodagTransmitter_RegisterTrigger(t *testing.T) {
	t.Run("registers trigger successfully", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{
			WorkflowID:   "test-workflow",
			WorkflowName: "Test Workflow",
		}

		config := &streams.Config{
			StreamIds:      []uint32{1, 2, 3},
			MaxFrequencyMs: 5000,
		}

		// Register trigger
		ch, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		require.Nil(t, capErr)
		require.NotNil(t, ch)

		// Verify subscriber was added
		assert.Len(t, transmitter.subscribers, 1)
		assert.Contains(t, transmitter.subscribers, "trigger-1")

		// Cleanup
		capErr = transmitter.UnregisterTrigger(ctx, "trigger-1", metadata, config)
		assert.Nil(t, capErr)
	})

	t.Run("returns error for nil config", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}

		// Try to register with nil config
		_, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, nil)
		assert.NotNil(t, capErr)
	})

	t.Run("returns error for invalid frequency", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}

		// Invalid frequency (not a multiple of ticker resolution)
		config := &streams.Config{
			StreamIds:      []uint32{1},
			MaxFrequencyMs: 1500, // Not a multiple of 1000
		}

		_, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		assert.NotNil(t, capErr)
	})

	t.Run("returns error for duplicate registration", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}
		config := &streams.Config{
			StreamIds:      []uint32{1},
			MaxFrequencyMs: 1000,
		}

		// First registration should succeed
		_, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		require.Nil(t, capErr)

		// Second registration with same ID should fail
		_, capErr = transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		assert.NotNil(t, capErr)
	})
}

func TestNodagTransmitter_UnregisterTrigger(t *testing.T) {
	t.Run("unregisters trigger successfully", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}
		config := &streams.Config{
			StreamIds:      []uint32{1},
			MaxFrequencyMs: 1000,
		}

		// Register
		ch, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		require.Nil(t, capErr)

		// Unregister
		capErr = transmitter.UnregisterTrigger(ctx, "trigger-1", metadata, config)
		assert.Nil(t, capErr)

		// Channel should be closed
		_, ok := <-ch
		assert.False(t, ok, "channel should be closed")

		// Subscriber should be removed
		assert.Empty(t, transmitter.subscribers)
	})

	t.Run("returns error for non-existent trigger", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}
		config := &streams.Config{
			StreamIds:      []uint32{1},
			MaxFrequencyMs: 1000,
		}

		// Try to unregister non-existent trigger
		capErr := transmitter.UnregisterTrigger(ctx, "trigger-1", metadata, config)
		assert.NotNil(t, capErr)
	})
}

func TestNodagTransmitter_Transmit(t *testing.T) {
	t.Run("processes capability trigger reports", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		// Start the service
		ctx := context.Background()
		err = transmitter.start(ctx)
		require.NoError(t, err)
		defer transmitter.close()

		// Register a subscriber
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}
		config := &streams.Config{
			StreamIds:      []uint32{1, 2, 3},
			MaxFrequencyMs: 1000,
		}

		ch, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		require.Nil(t, capErr)

		// Create a test report
		reportPayload := &capabilitiespb.OCRTriggerReport{
			EventID:   "test-event-1",
			Timestamp: 1000000000, // 1ms in nanoseconds
		}
		reportBytes, err := proto.Marshal(reportPayload)
		require.NoError(t, err)

		configDigest := ocr2types.ConfigDigest{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		seqNr := uint64(42)
		sigs := []types.AttributedOnchainSignature{
			{Signer: 1, Signature: []byte("sig1")},
			{Signer: 2, Signature: []byte("sig2")},
		}

		// Transmit the report
		err = transmitter.Transmit(
			ctx,
			configDigest,
			seqNr,
			ocr3types.ReportWithInfo[llotypes.ReportInfo]{
				Report: reportBytes,
				Info: llotypes.ReportInfo{
					ReportFormat:   llotypes.ReportFormatCapabilityTrigger,
					LifeCycleStage: datastreamsllo.LifeCycleStageProduction,
				},
			},
			sigs,
		)
		require.NoError(t, err)

		// Verify subscriber received the report
		select {
		case triggerEvent := <-ch:
			assert.Equal(t, "trigger-1", triggerEvent.Id)
			assert.Equal(t, seqNr, triggerEvent.Trigger.SeqNr)
			assert.Equal(t, configDigest[:], triggerEvent.Trigger.ConfigDigest)
			assert.Equal(t, reportBytes, triggerEvent.Trigger.Report)
			assert.Len(t, triggerEvent.Trigger.Sigs, 2)
		default:
			t.Fatal("Expected to receive trigger event")
		}
	})

	t.Run("ignores non-capability reports", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()

		// Transmit a non-capability report
		err = transmitter.Transmit(
			ctx,
			ocr2types.ConfigDigest{},
			1,
			ocr3types.ReportWithInfo[llotypes.ReportInfo]{
				Report: []byte("test"),
				Info: llotypes.ReportInfo{
					ReportFormat:   llotypes.ReportFormatJSON, // Not capability trigger
					LifeCycleStage: datastreamsllo.LifeCycleStageProduction,
				},
			},
			nil,
		)

		// Should not error, just silently ignore
		assert.NoError(t, err)
	})

	t.Run("ignores non-production reports", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()

		// Transmit a staging report
		err = transmitter.Transmit(
			ctx,
			ocr2types.ConfigDigest{},
			1,
			ocr3types.ReportWithInfo[llotypes.ReportInfo]{
				Report: []byte("test"),
				Info: llotypes.ReportInfo{
					ReportFormat:   llotypes.ReportFormatCapabilityTrigger,
					LifeCycleStage: datastreamsllo.LifeCycleStageStaging, // Not production
				},
			},
			nil,
		)

		// Should not error, just silently ignore
		assert.NoError(t, err)
	})
}

func TestNodagTransmitter_FrequencyThrottling(t *testing.T) {
	t.Run("respects per-subscriber frequency", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()

		// Register subscriber with 2000ms frequency
		metadata := capabilities.RequestMetadata{WorkflowID: "test-workflow"}
		config := &streams.Config{
			StreamIds:      []uint32{1},
			MaxFrequencyMs: 2000,
		}

		ch, capErr := transmitter.RegisterTrigger(ctx, "trigger-1", metadata, config)
		require.Nil(t, capErr)

		// Send reports at different timestamps
		// Note: First report at t=0 will be filtered by global throttling since lastReportMs starts at 0
		// So we start from t=1000ms
		for i, tsNs := range []uint64{1000000000, 2000000000, 3000000000, 4000000000} {
			reportPayload := &capabilitiespb.OCRTriggerReport{
				EventID:   fmt.Sprintf("event-%d", i),
				Timestamp: tsNs,
			}
			reportBytes, err := proto.Marshal(reportPayload)
			require.NoError(t, err)

			report := &streams.Report{
				ConfigDigest: []byte{1, 2, 3, 4},
				SeqNr:        uint64(i),
				Report:       reportBytes,
			}

			err = transmitter.processReport(ctx, report)
			require.NoError(t, err)
		}

		// Should receive reports at t=2000ms and t=4000ms (every 2000ms)
		// t=1000ms and t=3000ms should be skipped

		received := 0
		for {
			select {
			case event := <-ch:
				received++
				// Should only receive events at aligned timestamps (2000, 4000)
				// Which correspond to seqNr 1 and 3
				assert.True(t, event.Trigger.SeqNr == 1 || event.Trigger.SeqNr == 3)
			default:
				goto done
			}
		}
	done:
		assert.Equal(t, 2, received, "should receive exactly 2 reports")
	})
}

func TestNodagTransmitter_Lifecycle(t *testing.T) {
	t.Run("registers and unregisters with capability registry", func(t *testing.T) {
		registry := &mockCapabilityRegistry{}
		transmitter, err := NewNodagTransmitter(logger.Test(t), 1, registry)
		require.NoError(t, err)

		ctx := context.Background()

		// Start should register
		err = transmitter.start(ctx)
		require.NoError(t, err)
		assert.Len(t, registry.added, 1)

		// Close should unregister
		err = transmitter.close()
		require.NoError(t, err)
		assert.Len(t, registry.removed, 1)
		assert.Equal(t, nodagCapabilityID, registry.removed[0])
	})
}
