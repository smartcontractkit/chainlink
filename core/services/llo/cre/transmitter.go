package cre

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	streams "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/streams"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	streamstypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
)

const (
	defaultCapabilityName        = "streams-trigger"
	defaultCapabilityVersion     = "2.0.0" // v2 = LLO
	defaultTickerResolutionMs    = 1000
	defaultSendChannelBufferSize = 1000
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type TransmitterConfig struct {
	Logger               logger.Logger                  `json:"-"`
	CapabilitiesRegistry coretypes.CapabilitiesRegistry `json:"-"`
	DonID                uint32                         `json:"-"`

	TriggerCapabilityName        string `json:"triggerCapabilityName"`
	TriggerCapabilityVersion     string `json:"triggerCapabilityVersion"`
	TriggerTickerMinResolutionMs int    `json:"triggerTickerMinResolutionMs"`
	TriggerSendChannelBufferSize int    `json:"triggerSendChannelBufferSize"`
}

var _ Transmitter = &transmitter{}
var _ capabilities.TriggerCapability = &transmitter{}

type transmitter struct {
	services.Service
	eng *services.Engine
	capabilities.CapabilityInfo

	config      TransmitterConfig
	fromAccount ocr2types.Account
	registry    coretypes.CapabilitiesRegistry

	subscribers  map[string]*subscriber
	lastReportMs uint64
	mu           sync.Mutex
}

type subscriber struct {
	ch         chan<- capabilities.TriggerResponse
	workflowID string
	config     streamstypes.LLOTriggerConfig
}

func (c TransmitterConfig) NewTransmitter() (*transmitter, error) {
	return c.newTransmitter(c.Logger)
}

func (c TransmitterConfig) newTransmitter(lggr logger.Logger) (*transmitter, error) {
	t := &transmitter{
		config:      c,
		fromAccount: ocr2types.Account(lggr.Name() + strconv.FormatUint(uint64(c.DonID), 10)),
		registry:    c.CapabilitiesRegistry,
		subscribers: make(map[string]*subscriber),
	}
	if t.config.TriggerCapabilityName == "" {
		t.config.TriggerCapabilityName = defaultCapabilityName
	}
	if t.config.TriggerCapabilityVersion == "" {
		t.config.TriggerCapabilityVersion = defaultCapabilityVersion
	}
	if t.config.TriggerTickerMinResolutionMs == 0 {
		t.config.TriggerTickerMinResolutionMs = defaultTickerResolutionMs
	}
	if t.config.TriggerSendChannelBufferSize == 0 {
		t.config.TriggerSendChannelBufferSize = defaultSendChannelBufferSize
	}

	capInfo, err := capabilities.NewCapabilityInfo(
		// TODO(CAPPL-645): add labels
		t.config.TriggerCapabilityName+"@"+t.config.TriggerCapabilityVersion,
		capabilities.CapabilityTypeTrigger,
		"Streams LLO Trigger",
	)
	if err != nil {
		return nil, err
	}
	t.CapabilityInfo = capInfo

	t.Service, t.eng = services.Config{
		Name:  "CRETransmitter",
		Start: t.start,
		Close: t.close,
	}.NewServiceEngine(lggr)

	return t, nil
}

func (t *transmitter) start(ctx context.Context) error {
	t.eng.Infow("CRETransmitter starting, registering capability",
		"capabilityID", t.ID,
		"donID", t.config.DonID)
	return t.registry.Add(ctx, t)
}

func (t *transmitter) close() error {
	return t.registry.Remove(context.Background(), t.ID)
}

func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return t.fromAccount, nil
}

func (t *transmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) error {
	switch report.Info.ReportFormat {
	case llotypes.ReportFormatCapabilityTrigger:
		// Format 5: Native protobuf format designed for CRE
	case llotypes.ReportFormatEVMABIEncodeUnpackedExpr:
		// Format 7: ABI-encoded format (can also be used for CRE with ABI decoding)
	default:
		// NOTE: Silently ignore non-capability format reports here. All
		// channels are broadcast to all transmitters but this transmitter only
		// cares about channels of type ReportFormatCapabilityTrigger (5)
		// or ReportFormatEVMABIEncodeUnpackedExpr (7)
		return nil
	}
	switch report.Info.LifeCycleStage {
	case datastreamsllo.LifeCycleStageProduction:
	default:
		// NOTE: Ignore retirement and staging reports; for now we assume that
		// we only care about sending production reports.
		//
		// Support could be added in future e.g. for verifying blue-green
		// deploys etc.
		return nil
	}

	capSigs := make([]capabilities.OCRAttributedOnchainSignature, len(sigs))
	for i, sig := range sigs {
		capSigs[i] = capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(sig.Signer),
			Signature: sig.Signature,
		}
	}
	ev := &capabilities.OCRTriggerEvent{
		ConfigDigest: cd[:],
		SeqNr:        seqNr,
		Report:       report.Report,
		Sigs:         capSigs,
	}
	return t.processNewEvent(ctx, ev, report.Info.ReportFormat)
}

func (t *transmitter) processNewEvent(ctx context.Context, event *capabilities.OCRTriggerEvent, reportFormat llotypes.ReportFormat) error {
	var tsMs uint64
	var eventID string

	// Extract timestamp and eventID based on report format
	switch reportFormat {
	case llotypes.ReportFormatCapabilityTrigger:
		// Format 5: Protobuf OCRTriggerReport
		p := &capabilitiespb.OCRTriggerReport{}
		err := proto.Unmarshal(event.Report, p)
		if err != nil {
			return fmt.Errorf("failed to unmarshal OCRTriggerReport (Format 5): %w", err)
		}
		tsMs = p.Timestamp / 1000000 // nanoseconds -> milliseconds
		eventID = p.EventID
	case llotypes.ReportFormatEVMABIEncodeUnpackedExpr:
		// Format 7: ABI-encoded report
		// The timestamp is embedded in the ABI header at offset 64 (32+32) as uint32
		// For simplicity, use current time and generate unique event ID
		tsMs = uint64(time.Now().UnixMilli())
		eventID = fmt.Sprintf("streams_%d_%d_f7", t.config.DonID, tsMs)
	default:
		return fmt.Errorf("unsupported report format: %d", reportFormat)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if tsMs/uint64(t.config.TriggerTickerMinResolutionMs) == t.lastReportMs/uint64(t.config.TriggerTickerMinResolutionMs) { //nolint:gosec // disable G115
		// ignore reports that are too frequent
		return nil
	}
	t.lastReportMs = tsMs
	alignedTsMs := tsMs - tsMs%uint64(t.config.TriggerTickerMinResolutionMs) //nolint:gosec // disable G115

	// Convert OCRTriggerEvent to streams.Report proto (V2 format)
	// This allows the workflow SDK to properly deserialize the trigger event
	streamsReport := &streams.Report{
		ConfigDigest: event.ConfigDigest,
		SeqNr:        event.SeqNr,
		Report:       event.Report,
		Sigs:         make([]*streams.OCRSignature, len(event.Sigs)),
	}
	for i, sig := range event.Sigs {
		streamsReport.Sigs[i] = &streams.OCRSignature{
			Signer:    sig.Signer,
			Signature: sig.Signature,
		}
	}

	// Wrap in anypb.Any for V2 proto format - this goes in Event.Payload
	payload, err := anypb.New(streamsReport)
	if err != nil {
		return fmt.Errorf("failed to wrap streams.Report in anypb.Any: %w", err)
	}

	// Also keep backward compatibility with V1 format using Outputs
	// Note: For Format 7 (ABI), ToMap will likely fail but that's ok since V1 is deprecated
	o, mapErr := event.ToMap()
	if mapErr != nil && reportFormat == llotypes.ReportFormatCapabilityTrigger {
		t.eng.Warnw("failed to convert OCRTriggerEvent to map (V1 compat)", "error", mapErr)
	}

	capResponse := capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: t.ID,
			ID:          eventID,
			Payload:     payload, // V2 format: proto wrapped in anypb.Any
			Outputs:     o,       // V1 format: values.Map (for backward compat)
		},
	}

	t.eng.Debugw("ProcessReport pushing event", "eventID", eventID, "tsMs", tsMs, "alignedTsMs", alignedTsMs, "nSubscribers", len(t.subscribers))
	nIncludedSubscribers := 0
	for _, sub := range t.subscribers {
		// Handle case where MaxFrequencyMs is 0 (default/unset) - treat as "include every report"
		includeByFrequency := sub.config.MaxFrequencyMs == 0 || alignedTsMs%sub.config.MaxFrequencyMs == 0
		if includeByFrequency {
			// include this subscriber
			select {
			case sub.ch <- capResponse:
			case <-ctx.Done():
				t.eng.Error("context done, dropping event")
				return ctx.Err()
			default:
				// drop event if channel is full - processNewEvent() should be non-blocking
				t.eng.Errorw("subscriber channel full, dropping event", "eventID", eventID, "workflowID", sub.workflowID)
			}
			nIncludedSubscribers++
		}
	}
	t.eng.Debugw("ProcessReport done", "eventID", eventID, "nIncludedSubscribers", nIncludedSubscribers)
	return nil
}

func (t *transmitter) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	config, err := validateConfig(req.Config, &t.config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if _, ok := t.subscribers[req.TriggerID]; ok {
		return nil, fmt.Errorf("triggerId %s already registered", t.ID)
	}

	ch := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	t.subscribers[req.TriggerID] =
		&subscriber{
			ch:         ch,
			workflowID: req.Metadata.WorkflowID,
			config:     *config,
		}

	t.eng.Infow("RegisterTrigger: subscriber added",
		"triggerID", req.TriggerID,
		"workflowID", req.Metadata.WorkflowID,
		"streamIDs", config.StreamIDs,
		"maxFrequencyMs", config.MaxFrequencyMs,
		"totalSubscribers", len(t.subscribers))

	return ch, nil
}

func validateConfig(registerConfig *values.Map, capabilityConfig *TransmitterConfig) (*streamstypes.LLOTriggerConfig, error) {
	cfg := &streamstypes.LLOTriggerConfig{}
	if err := registerConfig.UnwrapTo(cfg); err != nil {
		return nil, err
	}
	if int64(cfg.MaxFrequencyMs)%int64(capabilityConfig.TriggerTickerMinResolutionMs) != 0 { //nolint:gosec // disable G115
		return nil, fmt.Errorf("MaxFrequencyMs must be a multiple of %d", capabilityConfig.TriggerTickerMinResolutionMs)
	}
	return cfg, nil
}

func (t *transmitter) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	subscriber, ok := t.subscribers[req.TriggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", t.ID)
	}
	close(subscriber.ch)
	delete(t.subscribers, req.TriggerID)

	t.eng.Infow("UnregisterTrigger: subscriber removed",
		"triggerID", req.TriggerID,
		"workflowID", subscriber.workflowID,
		"remainingSubscribers", len(t.subscribers))

	return nil
}
