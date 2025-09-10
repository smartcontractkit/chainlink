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

// These constants are used to identify the secure mint trigger capability.
const (
	defaultCapabilityName        = "securemint-trigger"
	defaultCapabilityVersion     = "1.0.0"
	defaultTickerResolutionMs    = 1000
	defaultSendChannelBufferSize = 1000
)

// Transmitter is a wrapper for ocr3types.ContractTransmitter[securemint.ChainSelector] to add the Service interface to it.
type Transmitter interface {
	ocr3types.ContractTransmitter[securemint.ChainSelector]
	services.Service
}

// TransmitterConfig is the configuration for the secure mint transmitter capability.
type TransmitterConfig struct {
	Logger               logger.Logger                  `json:"-"`
	CapabilitiesRegistry coretypes.CapabilitiesRegistry `json:"-"`
	DonID                uint32                         `json:"-"`

	TriggerCapabilityName        string `json:"triggerCapabilityName"`
	TriggerCapabilityVersion     string `json:"triggerCapabilityVersion"`
	TriggerTickerMinResolutionMs int    `json:"triggerTickerMinResolutionMs"`
	TriggerSendChannelBufferSize int    `json:"triggerSendChannelBufferSize"`
}

// NewTransmitter creates a new secure mint transmitter.
func (c TransmitterConfig) NewTransmitter(transmitterID string) (*transmitter, error) {
	c.Logger.Infow("Initializing SecureMintTransmitter", "triggerCapabilityName", c.TriggerCapabilityName, "triggerCapabilityVersion", c.TriggerCapabilityVersion)

	t := &transmitter{
		config:      c,
		fromAccount: ocr2types.Account(transmitterID),
		registry:    c.CapabilitiesRegistry,
		subscribers: make(map[string]*subscriber),
	}

	// set default values if not provided
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
	}.NewServiceEngine(c.Logger)

	t.eng.Infow("SecureMintTransmitter initialized", "triggerCapabilityName", c.TriggerCapabilityName, "triggerCapabilityVersion", c.TriggerCapabilityVersion)
	return t, nil
}

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

var _ Transmitter = &transmitter{}
var _ capabilities.TriggerCapability = &transmitter{}

func (t *transmitter) start(ctx context.Context) error {
	t.eng.Infow("Starting SecureMintTransmitter", "triggerCapabilityName", t.config.TriggerCapabilityName, "triggerCapabilityVersion", t.config.TriggerCapabilityVersion)
	err := t.registry.Add(ctx, t)
	if err != nil {
		return fmt.Errorf("failed to add transmitter to registry: %w", err)
	}
	return nil
}

func (t *transmitter) close() error {
	t.eng.Infow("Closing SecureMintTransmitter", "triggerCapabilityName", t.config.TriggerCapabilityName, "triggerCapabilityVersion", t.config.TriggerCapabilityVersion)
	return t.registry.Remove(context.Background(), t.CapabilityInfo.ID)
}

// FromAccount returns the CSA public key of this node.
func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	t.eng.Debugw("FromAccount", "fromAccount", t.fromAccount)
	return t.fromAccount, nil
}

// Transmit processes the secure mint report and transmits it as a trigger event to any subscribed workflows.
func (t *transmitter) Transmit(
	ctx context.Context,
	cd ocr2types.ConfigDigest,
	seqNr uint64,
	ocr3Report ocr3types.ReportWithInfo[securemint.ChainSelector],
	sigs []types.AttributedOnchainSignature,
) error {
	t.eng.Debugw("Transmit called", "cd", cd, "seqNr", seqNr, "report", ocr3Report, "sigs", sigs)

	// convert the secure mint report to a trigger event
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

	// use the seqNr as eventID to make sure we have unique event ids per report
	// and that nodes sending the same report use the same event id (to enable consensus in the Workflow DON to work properly).
	eventID := fmt.Sprintf("securemint_%d", seqNr)

	ev := &capabilities.TriggerEvent{
		TriggerType: t.CapabilityInfo.ID,
		ID:          eventID,
		Outputs:     outputs,
	}
	return t.processNewEvent(ctx, ev)
}

// processNewEvent sends the trigger event to any subscribed workflows.
func (t *transmitter) processNewEvent(ctx context.Context, event *capabilities.TriggerEvent) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.eng.Debugw("processNewEvent pushing event", "eventID", event.ID)

	capResponse := capabilities.TriggerResponse{
		Event: *event,
	}

	numIncludedSubscribers := 0
	for _, sub := range t.subscribers {
		// include all subscribers (no frequency limiting for now)
		select {
		case sub.ch <- capResponse:
		case <-ctx.Done():
			t.eng.Errorw("context done, dropping event", "eventID", event.ID)
			return ctx.Err()
		default:
			// drop event if channel is full - processNewEvent() should be non-blocking
			t.eng.Errorw("subscriber channel full, dropping event", "eventID", event.ID, "workflowID", sub.workflowID)
		}
		numIncludedSubscribers++
	}

	t.eng.Debugw("ProcessReport done", "eventID", event.ID, "numIncludedSubscribers", numIncludedSubscribers)
	return nil
}

// RegisterTrigger registers a new subscription to the secure mint trigger capability.
// This means that the workflow will receive a trigger event for each secure mint report.
func (t *transmitter) RegisterTrigger(ctx context.Context, req capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	t.eng.Debugw("RegisterTrigger", "triggerID", req.TriggerID, "metadata", req.Metadata)

	config, err := validateConfig(req.Config, &t.config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscribers[req.TriggerID]; ok {
		return nil, fmt.Errorf("triggerId %s already registered", t.ID)
	}

	ch := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	t.subscribers[req.TriggerID] = &subscriber{
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

// UnregisterTrigger unregisters a subscription to the secure mint trigger capability.
// This means that the workflow will no longer receive a trigger event for each secure mint report.
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

// subscriber contains the channel to send a trigger response to (normally a CRE workflow).
type subscriber struct {
	ch         chan<- capabilities.TriggerResponse
	workflowID string
	config     config.SecureMintTriggerConfig
}
