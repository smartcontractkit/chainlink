package cre

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	streamstypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	donID = 4
)

func Test_Transmitter(t *testing.T) {
	digest := types.ConfigDigest{1, 2, 3}
	sigs := []types.AttributedOnchainSignature{
		{
			Signer:    6,
			Signature: []byte{4, 5, 6},
		},
	}

	cfg := TransmitterConfig{
		Logger:               logger.TestLogger(t),
		CapabilitiesRegistry: nil,
		DonID:                donID,
	}
	tr, err := cfg.NewTransmitter()
	require.NoError(t, err)

	t.Run("invalid config", func(t *testing.T) {
		req := buildRegistrationRequest(t, "myID123", []streamstypes.LLOStreamID{12345, 67890}, 2300)
		_, err = tr.RegisterTrigger(t.Context(), req)
		require.Error(t, err)
	})

	t.Run("two registrations", func(t *testing.T) {
		req1 := buildRegistrationRequest(t, "wf1_trigger1", []streamstypes.LLOStreamID{12345, 67890}, 1000)
		req2 := buildRegistrationRequest(t, "wf2_trigger1", []streamstypes.LLOStreamID{67890}, 3000)
		respCh1, err := tr.RegisterTrigger(t.Context(), req1)
		require.NoError(t, err)
		respCh2, err := tr.RegisterTrigger(t.Context(), req2)
		require.NoError(t, err)

		require.NoError(t, tr.Transmit(t.Context(), digest, 1, encodeReport(t, 1023000000), sigs))
		require.NoError(t, tr.Transmit(t.Context(), digest, 2, encodeReport(t, 1803000000), sigs))
		require.NoError(t, tr.Transmit(t.Context(), digest, 3, encodeReport(t, 2101000000), sigs))
		require.NoError(t, tr.Transmit(t.Context(), digest, 4, encodeReport(t, 3456000000), sigs))
		require.NoError(t, tr.Transmit(t.Context(), digest, 5, encodeReport(t, 4502000000), sigs))
		require.NoError(t, tr.Transmit(t.Context(), digest, 6, encodeReport(t, 4777000000), sigs))
		require.Len(t, respCh1, 4) // every second
		require.Len(t, respCh2, 1) // every 3 seconds
	})
}

// mockRegistry is a CapabilitiesRegistry that no-ops Add/Remove so Start/Close don't fail.
type mockRegistry struct {
	coretypes.UnimplementedCapabilitiesRegistry
}

func (*mockRegistry) Add(ctx context.Context, c capabilities.BaseCapability) error { return nil }
func (*mockRegistry) Remove(ctx context.Context, ID string) error                  { return nil }

// Test_nextTransmissionBoundary verifies that the next boundary is aligned to the window.
func Test_nextTransmissionBoundary(t *testing.T) {
	// Window 100ms: boundaries at 0, 100, 200, ... ms from epoch
	now := time.UnixMilli(12345)
	next := nextTransmissionBoundary(now, 100)
	require.Equal(t, time.UnixMilli(12400), next)
	require.True(t, next.UnixMilli()%100 == 0)
}

// Test_Transmitter_delayedSend uses 1/8 of minimum window so the test completes quickly
// while still exercising the "wait until top of window" path.
func Test_Transmitter_delayedSend(t *testing.T) {
	windowMs := MinTransmissionWindowMs / 8 // 12ms for breathing room
	if windowMs < 1 {
		windowMs = 1
	}
	digest := types.ConfigDigest{1, 2, 3}
	sigs := []types.AttributedOnchainSignature{
		{Signer: 6, Signature: []byte{4, 5, 6}},
	}

	cfg := TransmitterConfig{
		Logger:               logger.TestLogger(t),
		CapabilitiesRegistry: &mockRegistry{},
		DonID:                donID,
		TransmissionWindowMs: windowMs,
	}
	tr, err := cfg.NewTransmitter()
	require.NoError(t, err)
	require.NoError(t, tr.Start(context.Background()))

	req := buildRegistrationRequest(t, "wf1", []streamstypes.LLOStreamID{12345}, 0)
	respCh, err := tr.RegisterTrigger(context.Background(), req)
	require.NoError(t, err)

	require.NoError(t, tr.Transmit(context.Background(), digest, 1, encodeReport(t, 1023000000), sigs))
	// Event is delayed until next boundary; wait for it (at most 2*windowMs)
	select {
	case r := <-respCh:
		require.Equal(t, "streams_4_1023000000", r.Event.ID)
	case <-time.After(time.Duration(2*windowMs) * time.Millisecond):
		t.Fatal("timed out waiting for delayed trigger event")
	}

	require.NoError(t, tr.Close())
}

func buildRegistrationRequest(t *testing.T, triggerID string, streamIDs []streamstypes.LLOStreamID, maxFrequencyMs uint64) capabilities.TriggerRegistrationRequest {
	cfg := &streamstypes.LLOTriggerConfig{
		StreamIDs:      streamIDs,
		MaxFrequencyMs: maxFrequencyMs,
	}
	wrappedCfg, err := values.WrapMap(cfg)
	require.NoError(t, err)

	return capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID,
		Config:    wrappedCfg,
	}
}

func encodeReport(t *testing.T, timestamp uint64) ocr3types.ReportWithInfo[llotypes.ReportInfo] {
	codec := NewReportCodecCapabilityTrigger(logger.TestLogger(t), donID)
	rep := llo.Report{
		ConfigDigest:                    types.ConfigDigest{1, 2, 3},
		SeqNr:                           32,
		ChannelID:                       llotypes.ChannelID(31),
		ValidAfterNanoseconds:           28,
		ObservationTimestampNanoseconds: timestamp,
		Values:                          []llo.StreamValue{llo.ToDecimal(decimal.NewFromInt(35)), llo.ToDecimal(decimal.NewFromInt(36))},
		Specimen:                        false,
	}
	cd := llotypes.ChannelDefinition{
		ReportFormat: llotypes.ReportFormatCapabilityTrigger,
		Streams: []llotypes.Stream{
			{StreamID: 1},
			{StreamID: 2},
		},
	}
	rawReport, err := codec.Encode(rep, cd, nil)
	require.NoError(t, err)

	return ocr3types.ReportWithInfo[llotypes.ReportInfo]{
		Report: rawReport,
		Info: llotypes.ReportInfo{
			LifeCycleStage: datastreamsllo.LifeCycleStageProduction,
			ReportFormat:   llotypes.ReportFormatCapabilityTrigger,
		},
	}
}
