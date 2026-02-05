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

	// MinTransmissionWindowMs is the minimum transmission window for production.
	// When TransmissionWindowMs > 0, sends are delayed until the next wall-clock boundary
	// so Streams DON nodes tend to send at the same time. Tests may use 1/4 or 1/8 of this
	// for faster runs (e.g. MinTransmissionWindowMs/8).
	MinTransmissionWindowMs = 100
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

	// TransmissionWindowMs delays pushing to subscribers until the next wall-clock boundary
	// (top of window), so Streams DON nodes are more likely to send at the same time.
	// 0 = no delay (immediate send). When > 0, use at least MinTransmissionWindowMs in production;
	TransmissionWindowMs int `json:"transmissionWindowMs"`
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

	// Delayed send until top of transmission window (when TransmissionWindowMs > 0)
	pendingMu      sync.Mutex
	pendingQueue   []pendingSend
	pendingWake    chan struct{}
	sendWorkerDone chan struct{}
	stopCtx        context.Context
	stopCancel     context.CancelFunc
}

type subscriber struct {
	ch         chan<- capabilities.TriggerResponse
	workflowID string
	config     streamstypes.LLOTriggerConfig
}

// pendingSend is a trigger response scheduled for the top of a transmission window (per-subscriber).
type pendingSend struct {
	response     capabilities.TriggerResponse
	targetTime   time.Time
	alignedTsMs  uint64
	eventID      string
	reportFormat uint32 // for subscriber filtering by accepted_report_formats
	triggerID    string // which subscriber to send to (workflow defines its transmission_window_ms)
}

func (c TransmitterConfig) NewTransmitter() (*transmitter, error) {
	return c.newTransmitter(c.Logger)
}

func (c TransmitterConfig) newTransmitter(lggr logger.Logger) (*transmitter, error) {
	stopCtx, stopCancel := context.WithCancel(context.Background())
	t := &transmitter{
		config:         c,
		fromAccount:    ocr2types.Account(lggr.Name() + strconv.FormatUint(uint64(c.DonID), 10)),
		registry:       c.CapabilitiesRegistry,
		subscribers:    make(map[string]*subscriber),
		pendingWake:    make(chan struct{}, 1),
		sendWorkerDone: make(chan struct{}),
		stopCtx:        stopCtx,
		stopCancel:     stopCancel,
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
	// Always run send worker so workflows that set transmission_window_ms in Config can use delayed send
	go t.sendWorker()
	t.eng.Infow("CRETransmitter starting, registering capability",
		"capabilityID", t.ID,
		"donID", t.config.DonID,
		"defaultTransmissionWindowMs", t.config.TransmissionWindowMs)
	return t.registry.Add(ctx, t)
}

func (t *transmitter) close() error {
	t.stopCancel()
	select {
	case t.pendingWake <- struct{}{}:
	default:
	}
	<-t.sendWorkerDone
	// Drain any remaining pending sends immediately
	t.pendingMu.Lock()
	for _, p := range t.pendingQueue {
		t.doPushToSubscriber(p.response, p.alignedTsMs, p.eventID, p.reportFormat, p.triggerID)
	}
	t.pendingQueue = nil
	t.pendingMu.Unlock()
	return t.registry.Remove(context.Background(), t.ID)
}

func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return t.fromAccount, nil
}

// nextTransmissionBoundary returns the next wall-clock time aligned to the given window (ms).
// Boundaries are at 0, windowMs, 2*windowMs, ... ms from Unix epoch.
func nextTransmissionBoundary(now time.Time, windowMs int) time.Time {
	if windowMs <= 0 {
		return now
	}
	unixMs := now.UnixMilli()
	boundaryMs := (unixMs/int64(windowMs) + 1) * int64(windowMs)
	return time.UnixMilli(boundaryMs)
}

// doPushToSubscribers sends the response to all subscribers that pass the frequency and accepted_report_formats filters.
// Caller must not hold t.mu; doPushToSubscribers locks as needed.
func (t *transmitter) doPushToSubscribers(response capabilities.TriggerResponse, alignedTsMs uint64, eventID string, reportFormat uint32) {
	hasPayload := response.Event.Payload != nil
	payloadType := "nil"
	if hasPayload {
		payloadType = response.Event.Payload.TypeUrl
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	nSent := 0
	nDropped := 0
	nFiltered := 0
	for _, sub := range t.subscribers {
		includeByFrequency := sub.config.MaxFrequencyMs == 0 || alignedTsMs%sub.config.MaxFrequencyMs == 0
		includeByReportFormat := len(sub.config.AcceptedReportFormats) == 0 || containsUint32(sub.config.AcceptedReportFormats, reportFormat)
		if includeByFrequency && includeByReportFormat {
			select {
			case sub.ch <- response:
				nSent++
				t.eng.Infow("CRETransmitter: Sent TriggerResponse to subscriber channel", "eventID", eventID, "workflowID", sub.workflowID, "hasPayload", hasPayload, "payloadType", payloadType)
			default:
				nDropped++
				t.eng.Errorw("CRETransmitter: subscriber channel full, dropping event", "eventID", eventID, "workflowID", sub.workflowID, "channelBufferSize", defaultSendChannelBufferSize)
			}
		} else {
			nFiltered++
		}
	}
	t.eng.Infow("ProcessReport done", "eventID", eventID, "nSubscribers", len(t.subscribers), "nSent", nSent, "nDropped", nDropped, "nFiltered", nFiltered)
}

// doPushToSubscriber sends the response to a single subscriber if it passes frequency and accepted_report_formats filters.
// Caller must not hold t.mu; doPushToSubscriber locks as needed.
func (t *transmitter) doPushToSubscriber(response capabilities.TriggerResponse, alignedTsMs uint64, eventID string, reportFormat uint32, triggerID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	sub, ok := t.subscribers[triggerID]
	if !ok {
		return
	}
	includeByFrequency := sub.config.MaxFrequencyMs == 0 || alignedTsMs%sub.config.MaxFrequencyMs == 0
	includeByReportFormat := len(sub.config.AcceptedReportFormats) == 0 || containsUint32(sub.config.AcceptedReportFormats, reportFormat)
	if !includeByFrequency || !includeByReportFormat {
		return
	}
	hasPayload := response.Event.Payload != nil
	payloadType := "nil"
	if hasPayload {
		payloadType = response.Event.Payload.TypeUrl
	}
	select {
	case sub.ch <- response:
		t.eng.Infow("CRETransmitter: Sent TriggerResponse to subscriber channel", "eventID", eventID, "workflowID", sub.workflowID, "triggerID", triggerID, "hasPayload", hasPayload, "payloadType", payloadType)
	default:
		t.eng.Errorw("CRETransmitter: subscriber channel full, dropping event", "eventID", eventID, "workflowID", sub.workflowID, "triggerID", triggerID, "channelBufferSize", defaultSendChannelBufferSize)
	}
}

func containsUint32(slice []uint32, v uint32) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}

func (t *transmitter) sendWorker() {
	defer close(t.sendWorkerDone)
	for {
		t.pendingMu.Lock()
		for len(t.pendingQueue) == 0 {
			t.pendingMu.Unlock()
			select {
			case <-t.pendingWake:
				t.pendingMu.Lock()
				continue
			case <-t.stopCtx.Done():
				return
			}
		}
		first := t.pendingQueue[0]
		t.pendingQueue = t.pendingQueue[1:]
		targetTime := first.targetTime
		t.pendingMu.Unlock()

		waitDur := time.Until(targetTime)
		if waitDur > 0 {
			select {
			case <-time.After(waitDur):
			case <-t.pendingWake:
				t.pendingMu.Lock()
				t.pendingQueue = append([]pendingSend{first}, t.pendingQueue...)
				t.pendingMu.Unlock()
				continue
			case <-t.stopCtx.Done():
				return
			}
		}

		t.doPushToSubscriber(first.response, first.alignedTsMs, first.eventID, first.reportFormat, first.triggerID)
	}
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
		// Protobuf-encoded reports (Streams DON format type)
	case llotypes.ReportFormatEVMABIEncodeUnpackedExpr:
		// ABI-encoded reports (Streams DON format type)
	default:
		// NOTE: Silently ignore non-capability format reports here. All
		// channels are broadcast to all transmitters but this transmitter only
		// cares about channels of type ReportFormatCapabilityTrigger (protobuf)
		// or ReportFormatEVMABIEncodeUnpackedExpr (ABI-encoded)
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
		// Protobuf-encoded OCRTriggerReport (Streams DON format type)
		p := &capabilitiespb.OCRTriggerReport{}
		err := proto.Unmarshal(event.Report, p)
		if err != nil {
			return fmt.Errorf("failed to unmarshal protobuf-encoded OCRTriggerReport: %w", err)
		}
		tsMs = p.Timestamp / 1000000 // nanoseconds -> milliseconds
		eventID = p.EventID
	case llotypes.ReportFormatEVMABIEncodeUnpackedExpr:
		// ABI-encoded report (Streams DON format type)
		// The timestamp is embedded in the ABI header at offset 64 (32+32) as uint32
		// For simplicity, use current time and generate unique event ID
		tsMs = uint64(time.Now().UnixMilli())
		eventID = fmt.Sprintf("streams_%d_%d_f7", t.config.DonID, tsMs)
	default:
		return fmt.Errorf("unsupported report format: %d", reportFormat)
	}

	t.mu.Lock()
	if tsMs/uint64(t.config.TriggerTickerMinResolutionMs) == t.lastReportMs/uint64(t.config.TriggerTickerMinResolutionMs) { //nolint:gosec // disable G115
		// ignore reports that are too frequent
		t.mu.Unlock()
		return nil
	}
	t.lastReportMs = tsMs
	alignedTsMs := tsMs - tsMs%uint64(t.config.TriggerTickerMinResolutionMs) //nolint:gosec // disable G115

	// Convert OCRTriggerEvent to streams.Report proto (V2 format)
	// This allows the workflow SDK to properly deserialize the trigger event.
	// ReportFormat is set so the workflow can decode without "try both" (5 = CapabilityTrigger, 7 = EVMABIEncodeUnpackedExpr).
	streamsReport := &streams.Report{
		ConfigDigest: event.ConfigDigest,
		SeqNr:        event.SeqNr,
		Report:       event.Report,
		Sigs:         make([]*streams.OCRSignature, len(event.Sigs)),
	}
	if reportFormat == llotypes.ReportFormatCapabilityTrigger || reportFormat == llotypes.ReportFormatEVMABIEncodeUnpackedExpr {
		streamsReport.ReportFormat = uint32(reportFormat)
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
	// Note: For ABI-encoded reports, ToMap will likely fail but that's ok since V1 is deprecated
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

	hasPayload := capResponse.Event.Payload != nil
	payloadType := "nil"
	if hasPayload {
		payloadType = capResponse.Event.Payload.TypeUrl
	}
	t.eng.Infow("ProcessReport pushing event", "eventID", eventID, "tsMs", tsMs, "alignedTsMs", alignedTsMs, "nSubscribers", len(t.subscribers), "hasPayload", hasPayload, "payloadType", payloadType)

	reportFormatU32 := uint32(reportFormat)
	now := time.Now()
	var sendNow []string
	var toQueue []pendingSend
	for triggerID, sub := range t.subscribers {
		includeByFrequency := sub.config.MaxFrequencyMs == 0 || alignedTsMs%sub.config.MaxFrequencyMs == 0
		includeByReportFormat := len(sub.config.AcceptedReportFormats) == 0 || containsUint32(sub.config.AcceptedReportFormats, reportFormatU32)
		if !includeByFrequency || !includeByReportFormat {
			continue
		}
		windowMs := sub.config.TransmissionWindowMs
		if windowMs == 0 {
			windowMs = t.config.TransmissionWindowMs
		}
		if windowMs > 0 {
			toQueue = append(toQueue, pendingSend{
				response:     capResponse,
				targetTime:   nextTransmissionBoundary(now, windowMs),
				alignedTsMs:  alignedTsMs,
				eventID:      eventID,
				reportFormat: reportFormatU32,
				triggerID:    triggerID,
			})
		} else {
			sendNow = append(sendNow, triggerID)
		}
	}
	t.mu.Unlock()

	for _, triggerID := range sendNow {
		t.doPushToSubscriber(capResponse, alignedTsMs, eventID, reportFormatU32, triggerID)
	}
	if len(toQueue) > 0 {
		t.pendingMu.Lock()
		t.pendingQueue = append(t.pendingQueue, toQueue...)
		t.pendingMu.Unlock()
		select {
		case t.pendingWake <- struct{}{}:
		default:
		}
	}
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
		"acceptedReportFormats", config.AcceptedReportFormats,
		"transmissionWindowMs", config.TransmissionWindowMs,
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
