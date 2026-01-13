package cre

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
)

const (
	nodagCapabilityID          = "streams-trigger@2.0.0"
	nodagTickerResolutionMs    = 1000
	nodagSendChannelBufferSize = 1000
)

// nodagTransmitter is a standalone NoDAG implementation of the Streams LLO Trigger.
// It receives reports from the LLO plugin via the Transmitter interface and emits them
// to workflow subscribers via the StreamsCapability interface.
//
// Architecture:
//
//	LLO Plugin → nodagTransmitter.Transmit() → processReport() → subscribers
//
// This is a clean NoDAG implementation that uses proto-based types throughout.
type nodagTransmitter struct {
	services.Service
	eng *services.Engine
	capabilities.CapabilityInfo

	lggr        logger.Logger
	donID       uint32
	fromAccount ocr2types.Account
	registry    coretypes.CapabilitiesRegistry

	// Subscriber management
	subscribers  map[string]*nodagSubscriber
	lastReportMs uint64
	mu           sync.Mutex
}

// nodagSubscriber holds the state for a single workflow subscription
type nodagSubscriber struct {
	ch         chan<- capabilities.TriggerAndId[*streams.Report]
	workflowID string
	config     *streams.Config
}

// NewNodagTransmitter creates a new standalone NoDAG transmitter
func NewNodagTransmitter(lggr logger.Logger, donID uint32, registry coretypes.CapabilitiesRegistry) (*nodagTransmitter, error) {
	t := &nodagTransmitter{
		lggr:        lggr,
		donID:       donID,
		fromAccount: ocr2types.Account(lggr.Name() + strconv.FormatUint(uint64(donID), 10)),
		registry:    registry,
		subscribers: make(map[string]*nodagSubscriber),
	}

	capInfo, err := capabilities.NewCapabilityInfo(
		nodagCapabilityID,
		capabilities.CapabilityTypeTrigger,
		"Streams LLO NoDAG Trigger",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create capability info: %w", err)
	}
	t.CapabilityInfo = capInfo

	t.Service, t.eng = services.Config{
		Name:  "NodagTransmitter",
		Start: t.start,
		Close: t.close,
	}.NewServiceEngine(lggr)

	return t, nil
}

func (t *nodagTransmitter) start(ctx context.Context) error {
	return t.registry.Add(ctx, t)
}

func (t *nodagTransmitter) close() error {
	return t.registry.Remove(context.Background(), t.ID)
}

// FromAccount implements llotypes.Transmitter
func (t *nodagTransmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return t.fromAccount, nil
}

// Transmit receives reports from the LLO plugin and distributes them to subscribers.
// This implements the llotypes.Transmitter interface.
func (t *nodagTransmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) error {
	// Only process capability trigger reports
	if report.Info.ReportFormat != llotypes.ReportFormatCapabilityTrigger {
		return nil
	}

	// Only process production reports
	if report.Info.LifeCycleStage != datastreamsllo.LifeCycleStageProduction {
		return nil
	}

	// Convert OCR signatures to proto format
	protoSigs := make([]*streams.OCRSignature, len(sigs))
	for i, sig := range sigs {
		protoSigs[i] = &streams.OCRSignature{
			Signer:    uint32(sig.Signer),
			Signature: sig.Signature,
		}
	}

	// Create proto Report
	protoReport := &streams.Report{
		ConfigDigest: cd[:],
		SeqNr:        seqNr,
		Report:       report.Report,
		Sigs:         protoSigs,
	}

	return t.processReport(ctx, protoReport)
}

// processReport distributes a report to all subscribers that should receive it
func (t *nodagTransmitter) processReport(ctx context.Context, report *streams.Report) error {
	// Extract timestamp from the report for frequency throttling
	p := &capabilitiespb.OCRTriggerReport{}
	if err := proto.Unmarshal(report.Report, p); err != nil {
		return fmt.Errorf("failed to unmarshal OCRTriggerReport: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Global frequency throttling - ignore reports that are too frequent
	tsMs := p.Timestamp / 1000000 // nanoseconds -> milliseconds
	if tsMs/uint64(nodagTickerResolutionMs) == t.lastReportMs/uint64(nodagTickerResolutionMs) {
		return nil
	}
	t.lastReportMs = tsMs

	// Align timestamp to ticker resolution for subscriber filtering
	alignedTsMs := tsMs - tsMs%uint64(nodagTickerResolutionMs)

	t.eng.Debugw("ProcessReport distributing", "eventID", p.EventID, "seqNr", report.SeqNr, "tsMs", tsMs, "alignedTsMs", alignedTsMs, "nSubscribers", len(t.subscribers))

	// Distribute to subscribers based on their frequency configuration
	nIncluded := 0
	for triggerID, sub := range t.subscribers {
		// Check if this subscriber should receive this report based on frequency
		if alignedTsMs%sub.config.MaxFrequencyMs == 0 {
			// Check if report contains any of the subscriber's requested stream IDs
			if shouldIncludeSubscriber(report, sub.config) {
				triggerEvent := capabilities.TriggerAndId[*streams.Report]{
					Id:      triggerID,
					Trigger: report,
				}

				select {
				case sub.ch <- triggerEvent:
					t.eng.Debugw("Sent report to subscriber", "triggerID", triggerID, "seqNr", report.SeqNr)
					nIncluded++
				case <-ctx.Done():
					t.eng.Error("Context done, stopping report distribution")
					return ctx.Err()
				default:
					// Non-blocking send - drop if channel is full
					t.eng.Warnw("Subscriber channel full, dropping report", "triggerID", triggerID, "workflowID", sub.workflowID, "seqNr", report.SeqNr)
				}
			}
		}
	}

	t.eng.Debugw("ProcessReport done", "eventID", p.EventID, "nIncluded", nIncluded, "nTotal", len(t.subscribers))
	return nil
}

// shouldIncludeSubscriber checks if a report should be sent to a subscriber
// based on the stream IDs they've configured
func shouldIncludeSubscriber(report *streams.Report, config *streams.Config) bool {
	// If no stream IDs configured, send all reports
	if len(config.StreamIds) == 0 {
		return true
	}

	// Parse the report to extract which streams it contains
	// For now, we send all reports to all subscribers
	// TODO: Implement stream ID filtering once we have the report format defined
	return true
}

// RegisterTrigger implements the NoDAG API for trigger registration
func (t *nodagTransmitter) RegisterTrigger(
	ctx context.Context,
	triggerID string,
	metadata capabilities.RequestMetadata,
	input *streams.Config,
) (<-chan capabilities.TriggerAndId[*streams.Report], caperrors.Error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Validate config
	if input == nil {
		return nil, caperrors.NewError(
			fmt.Errorf("config is nil"),
			caperrors.VisibilityPublic,
			caperrors.OriginSystem,
			caperrors.InvalidArgument,
		)
	}

	// Validate MaxFrequencyMs is a multiple of ticker resolution
	if int64(input.MaxFrequencyMs)%int64(nodagTickerResolutionMs) != 0 {
		return nil, caperrors.NewError(
			fmt.Errorf("MaxFrequencyMs must be a multiple of %d", nodagTickerResolutionMs),
			caperrors.VisibilityPublic,
			caperrors.OriginSystem,
			caperrors.InvalidArgument,
		)
	}

	// Check if already registered
	if _, exists := t.subscribers[triggerID]; exists {
		return nil, caperrors.NewError(
			fmt.Errorf("trigger %s already registered", triggerID),
			caperrors.VisibilityPublic,
			caperrors.OriginSystem,
			caperrors.InvalidArgument,
		)
	}

	// Create channel for this subscriber
	ch := make(chan capabilities.TriggerAndId[*streams.Report], nodagSendChannelBufferSize)

	// Register subscriber
	t.subscribers[triggerID] = &nodagSubscriber{
		ch:         ch,
		workflowID: metadata.WorkflowID,
		config:     input,
	}

	t.eng.Debugw("Registered trigger", "triggerID", triggerID, "workflowID", metadata.WorkflowID, "streamIds", input.StreamIds, "maxFrequencyMs", input.MaxFrequencyMs)

	return ch, nil
}

// UnregisterTrigger implements the NoDAG API for trigger unregistration
func (t *nodagTransmitter) UnregisterTrigger(
	ctx context.Context,
	triggerID string,
	metadata capabilities.RequestMetadata,
	input *streams.Config,
) caperrors.Error {
	t.mu.Lock()
	defer t.mu.Unlock()

	sub, exists := t.subscribers[triggerID]
	if !exists {
		return caperrors.NewError(
			fmt.Errorf("trigger %s not registered", triggerID),
			caperrors.VisibilityPublic,
			caperrors.OriginSystem,
			caperrors.InvalidArgument,
		)
	}

	// Close channel and remove subscriber
	close(sub.ch)
	delete(t.subscribers, triggerID)

	t.eng.Debugw("Unregistered trigger", "triggerID", triggerID, "workflowID", metadata.WorkflowID)

	return nil
}

// Ensure nodagTransmitter implements required interfaces
var _ llotypes.Transmitter = &nodagTransmitter{}
var _ services.Service = &nodagTransmitter{}

// Note: The NoDAG trigger registration methods (RegisterTrigger/UnregisterTrigger) provide
// the StreamsCapability interface, which can be wrapped by the generated server
