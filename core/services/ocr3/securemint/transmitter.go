package securemint

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	coretypes "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core/securemint"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const (
	defaultCapabilityName        = "securemint-trigger"
	defaultCapabilityVersion     = "1.0.0"
	defaultTickerResolutionMs    = 1000
	defaultSendChannelBufferSize = 1000
)

type Transmitter interface {
	ocr3types.ContractTransmitter[securemint.ChainSelector]
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

	subscribers map[string]*subscriber
	mu          sync.Mutex
}

type subscriber struct {
	ch         chan<- capabilities.TriggerResponse
	workflowID string
	config     config.SecureMintTriggerConfig
}

func (c TransmitterConfig) NewTransmitter(transmitterID string) (*transmitter, error) {
	return c.newTransmitter(c.Logger, transmitterID)
}

func (c TransmitterConfig) newTransmitter(lggr logger.Logger, transmitterID string) (*transmitter, error) {
	lggr.Infow("Initializing SecureMintTransmitter", "triggerCapabilityName", c.TriggerCapabilityName, "triggerCapabilityVersion", c.TriggerCapabilityVersion)
	t := &transmitter{
		config:      c,
		fromAccount: ocr2types.Account(transmitterID),
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
		"Secure Mint Trigger",
	)
	if err != nil {
		return nil, err
	}
	t.CapabilityInfo = capInfo

	t.Service, t.eng = services.Config{
		Name:  "SecureMintTransmitter",
		Start: t.start,
		Close: t.close,
	}.NewServiceEngine(lggr)

	t.eng.Infow("SecureMintTransmitter initialized", "triggerCapabilityName", t.config.TriggerCapabilityName, "triggerCapabilityVersion", t.config.TriggerCapabilityVersion)
	return t, nil
}

func (t *transmitter) start(ctx context.Context) error {
	t.eng.Infow("Starting SecureMintTransmitter", "triggerCapabilityName", t.config.TriggerCapabilityName, "triggerCapabilityVersion", t.config.TriggerCapabilityVersion)
	err := t.registry.Add(ctx, t)
	if err != nil {
		return fmt.Errorf("failed to add transmitter to registry: %w", err)
	}
	t.eng.Infow("SecureMintTransmitter registered", "triggerCapabilityInfo", t.CapabilityInfo)
	return nil
}

func (t *transmitter) close() error {
	t.eng.Infow("Closing SecureMintTransmitter", "triggerCapabilityName", t.config.TriggerCapabilityName, "triggerCapabilityVersion", t.config.TriggerCapabilityVersion)
	return t.registry.Remove(context.Background(), t.CapabilityInfo.ID)
}

func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	t.eng.Debugw("FromAccount", "fromAccount", t.fromAccount)
	return t.fromAccount, nil
}

func (t *transmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	ocr3Report ocr3types.ReportWithInfo[securemint.ChainSelector],
	sigs []types.AttributedOnchainSignature,
) error {
	t.eng.Debugw("Transmit called", "cd", cd, "seqNr", seqNr, "report", ocr3Report, "sigs", sigs)
	// Process the secure mint report and convert it to a trigger event
	capSigs := make([]capabilities.OCRAttributedOnchainSignature, len(sigs))
	for i, sig := range sigs {
		capSigs[i] = capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(sig.Signer),
			Signature: sig.Signature,
		}
	}

	jsonOcr3Report, err := json.Marshal(ocr3Report)
	if err != nil {
		return fmt.Errorf("failed to marshal ocr3 report: %w", err)
	}

	outputs, err := values.NewMap(map[string]any{
		"report":       jsonOcr3Report,
		"sigs":         capSigs,
		"seqNr":        seqNr,
		"configDigest": cd,
	})
	if err != nil {
		return fmt.Errorf("failed to create outputs map: %w", err)
	}

	// use the seqNr as eventID
	eventID := fmt.Sprintf("securemint_%d", seqNr)

	ev := &capabilities.TriggerEvent{
		TriggerType: t.CapabilityInfo.ID,
		ID:          eventID,
		Outputs:     outputs,
	}
	return t.processNewEvent(ctx, ev)
}

func (t *transmitter) processNewEvent(ctx context.Context, event *capabilities.TriggerEvent) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	capResponse := capabilities.TriggerResponse{
		Event: *event,
	}

	t.eng.Debugw("ProcessReport pushing event", "eventID", event.ID)
	nIncludedSubscribers := 0
	for _, sub := range t.subscribers {
		// include this subscriber (no frequency limiting as requested)
		select {
		case sub.ch <- capResponse:
		case <-ctx.Done():
			t.eng.Error("context done, dropping event")
			return ctx.Err()
		default:
			// drop event if channel is full - processNewEvent() should be non-blocking
			t.eng.Errorw("subscriber channel full, dropping event", "eventID", event.ID, "workflowID", sub.workflowID)
		}
		nIncludedSubscribers++
	}
	t.eng.Debugw("ProcessReport done", "eventID", event.ID, "nIncludedSubscribers", nIncludedSubscribers)
	return nil
}

func (t *transmitter) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	t.eng.Debugw("RegisterTrigger", "triggerID", req.TriggerID, "metadata", req.Metadata)
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
	return ch, nil
}

func validateConfig(registerConfig *values.Map, capabilityConfig *TransmitterConfig) (*config.SecureMintTriggerConfig, error) {
	cfg := &config.SecureMintTriggerConfig{}
	if err := registerConfig.UnwrapTo(cfg); err != nil {
		return nil, err
	}
	if int64(cfg.MaxFrequencyMs)%int64(capabilityConfig.TriggerTickerMinResolutionMs) != 0 { //nolint:gosec // disable G115
		return nil, fmt.Errorf("MaxFrequencyMs must be a multiple of %d", capabilityConfig.TriggerTickerMinResolutionMs)
	}
	return cfg, nil
}

func (t *transmitter) UnregisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) error {
	t.eng.Debugw("UnregisterTrigger", "triggerID", req.TriggerID)
	t.mu.Lock()
	defer t.mu.Unlock()

	subscriber, ok := t.subscribers[req.TriggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", t.ID)
	}
	close(subscriber.ch)
	delete(t.subscribers, req.TriggerID)
	return nil
}
