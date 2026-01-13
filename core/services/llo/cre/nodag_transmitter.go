package cre

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	streamstypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
)

// nodagTransmitter wraps the existing legacy transmitter to provide the NoDAG capability API.
// This is a thin adapter layer that converts between the proto-based NoDAG API and the
// existing legacy RegisterTrigger/UnregisterTrigger API.
//
// Architecture:
//
//	LLO Plugin → LLO Transmitter → CRE Sub-Transmitter (this) → NoDAG API → Workflows
//
// The NoDAG API uses:
//   - Proto-based Config message (type-safe configuration)
//   - Proto-based Report message (type-safe output)
//   - Streaming RPC pattern (instead of RegisterTrigger/UnregisterTrigger)
type nodagTransmitter struct {
	*transmitter // Embed the existing legacy transmitter

	// NoDAG-specific state
	lggr logger.Logger
}

// NewNodagTransmitter creates a new NoDAG-compatible transmitter.
// This wraps the existing legacy transmitter and provides the NoDAG API.
func NewNodagTransmitter(cfg TransmitterConfig) (*nodagTransmitter, error) {
	// Create the legacy transmitter
	legacy, err := cfg.newTransmitter(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create legacy transmitter: %w", err)
	}

	return &nodagTransmitter{
		transmitter: legacy,
		lggr:        cfg.Logger,
	}, nil
}

// RegisterTrigger implements the NoDAG capability API.
// This is the main entry point for workflows to subscribe to LLO reports.
//
// The NoDAG API differs from the legacy API in several ways:
//  1. Takes typed proto Config instead of generic values.Map
//  2. Returns typed channel of proto Report instead of generic TriggerResponse
//  3. Returns caperrors.Error instead of standard error
//  4. Separates triggerID and metadata as explicit parameters
//
// This method bridges between the two APIs by:
//  1. Converting proto Config → legacy LLOTriggerConfig
//  2. Calling legacy RegisterTrigger
//  3. Converting legacy channel → typed proto channel
func (t *nodagTransmitter) RegisterTrigger(
	ctx context.Context,
	triggerID string,
	metadata capabilities.RequestMetadata,
	input *streams.Config,
) (<-chan capabilities.TriggerAndId[*streams.Report], caperrors.Error) {
	// 1. Convert proto Config to legacy LLOTriggerConfig
	legacyConfig, err := convertProtoConfigToLegacy(input)
	if err != nil {
		return nil, caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.InvalidArgument)
	}

	// 2. Convert legacy config to values.Map for legacy API
	configMap, err := values.WrapMap(legacyConfig)
	if err != nil {
		return nil, caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.InvalidArgument)
	}

	// 3. Call legacy RegisterTrigger
	legacyReq := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID,
		Metadata:  metadata,
		Config:    configMap,
	}

	legacyCh, err := t.transmitter.RegisterTrigger(ctx, legacyReq)
	if err != nil {
		return nil, caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.Internal)
	}

	// 4. Create typed proto channel
	protoCh := make(chan capabilities.TriggerAndId[*streams.Report], t.config.TriggerSendChannelBufferSize)

	// 5. Bridge between legacy channel and proto channel
	go t.bridgeChannels(ctx, triggerID, legacyCh, protoCh)

	return protoCh, nil
}

// UnregisterTrigger implements the NoDAG capability API for unsubscribing.
func (t *nodagTransmitter) UnregisterTrigger(
	ctx context.Context,
	triggerID string,
	metadata capabilities.RequestMetadata,
	input *streams.Config,
) caperrors.Error {
	// Convert proto Config to legacy format
	legacyConfig, err := convertProtoConfigToLegacy(input)
	if err != nil {
		return caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.InvalidArgument)
	}

	configMap, err := values.WrapMap(legacyConfig)
	if err != nil {
		return caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.InvalidArgument)
	}

	// Call legacy UnregisterTrigger
	legacyReq := capabilities.TriggerRegistrationRequest{
		TriggerID: triggerID,
		Metadata:  metadata,
		Config:    configMap,
	}

	if err := t.transmitter.UnregisterTrigger(ctx, legacyReq); err != nil {
		return caperrors.NewError(err, caperrors.VisibilityPublic, caperrors.OriginSystem, caperrors.Internal)
	}

	return nil
}

// bridgeChannels converts events from the legacy channel to the typed proto channel.
// This runs in a goroutine and handles the conversion of each event.
func (t *nodagTransmitter) bridgeChannels(
	ctx context.Context,
	triggerID string,
	legacyCh <-chan capabilities.TriggerResponse,
	protoCh chan<- capabilities.TriggerAndId[*streams.Report],
) {
	defer close(protoCh)

	for {
		select {
		case <-ctx.Done():
			t.lggr.Debugw("Context cancelled, stopping channel bridge", "triggerID", triggerID)
			return

		case legacyResp, ok := <-legacyCh:
			if !ok {
				t.lggr.Debugw("Legacy channel closed, stopping bridge", "triggerID", triggerID)
				return
			}

			// Convert legacy TriggerResponse to proto Report
			protoReport, err := convertLegacyResponseToProto(legacyResp)
			if err != nil {
				t.lggr.Errorw("Failed to convert legacy response to proto",
					"triggerID", triggerID,
					"error", err,
				)
				continue
			}

			// Send to proto channel
			select {
			case protoCh <- capabilities.TriggerAndId[*streams.Report]{
				Id:      triggerID,
				Trigger: protoReport,
			}:
				t.lggr.Debugw("Sent proto report",
					"triggerID", triggerID,
					"seqNr", protoReport.SeqNr,
				)

			case <-ctx.Done():
				t.lggr.Debugw("Context cancelled while sending", "triggerID", triggerID)
				return

			default:
				// Channel full, drop event (same behavior as legacy transmitter)
				t.lggr.Warnw("Proto channel full, dropping event",
					"triggerID", triggerID,
					"seqNr", protoReport.SeqNr,
				)
			}
		}
	}
}

// convertProtoConfigToLegacy converts the proto Config message to the legacy LLOTriggerConfig.
//
// Proto Config:
//   - stream_ids: repeated uint32
//   - max_frequency_ms: uint64
//
// Legacy LLOTriggerConfig:
//   - StreamIDs: []LLOStreamID (which is uint32)
//   - MaxFrequencyMs: uint64
func convertProtoConfigToLegacy(proto *streams.Config) (*streamstypes.LLOTriggerConfig, error) {
	if proto == nil {
		return nil, fmt.Errorf("proto config is nil")
	}

	// Convert stream IDs
	streamIDs := make([]streamstypes.LLOStreamID, len(proto.StreamIds))
	for i, id := range proto.StreamIds {
		streamIDs[i] = streamstypes.LLOStreamID(id)
	}

	return &streamstypes.LLOTriggerConfig{
		StreamIDs:      streamIDs,
		MaxFrequencyMs: proto.MaxFrequencyMs,
	}, nil
}

// convertLegacyResponseToProto converts a legacy TriggerResponse to a proto Report.
//
// Legacy TriggerResponse contains:
//   - Event.Outputs["event"] = OCRTriggerEvent with:
//   - ConfigDigest: []byte
//   - SeqNr: uint64
//   - Report: []byte
//   - Sigs: []OCRAttributedOnchainSignature
//
// Proto Report contains:
//   - config_digest: bytes
//   - seq_nr: uint64
//   - report: bytes
//   - sigs: []OCRSignature
func convertLegacyResponseToProto(legacy capabilities.TriggerResponse) (*streams.Report, error) {
	// Extract OCRTriggerEvent from legacy response
	// legacy.Event.Outputs is a *values.Map, need to get the "event" field
	var event capabilities.OCRTriggerEvent
	if err := legacy.Event.Outputs.UnwrapTo(&event); err != nil {
		return nil, fmt.Errorf("failed to unwrap event from outputs: %w", err)
	}

	// Convert signatures
	sigs := make([]*streams.OCRSignature, len(event.Sigs))
	for i, sig := range event.Sigs {
		sigs[i] = &streams.OCRSignature{
			Signer:    sig.Signer,
			Signature: sig.Signature,
		}
	}

	// Create proto Report
	return &streams.Report{
		ConfigDigest: event.ConfigDigest,
		SeqNr:        event.SeqNr,
		Report:       event.Report,
		Sigs:         sigs,
	}, nil
}

// Ensure nodagTransmitter implements the required interfaces
var _ Transmitter = (*nodagTransmitter)(nil)
var _ services.Service = (*nodagTransmitter)(nil)
