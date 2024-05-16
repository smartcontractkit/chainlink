package triggers

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var capInfo = capabilities.MustNewCapabilityInfo(
	"streams-trigger",
	capabilities.CapabilityTypeTrigger,
	"Streams Trigger",
	"v1.0.0",
	nil,
)

const defaultTickerResolutionMs = 1000

// TODO pending capabilities configuration implementation - this should be configurable with a sensible default
const defaultSendChannelBufferSize = 1000

type config struct {
	// strings should be hex-encoded 32-byte values, prefixed with "0x", all lowercase, minimum 1 item
	FeedIDs []string `json:"feedIds" jsonschema:"pattern=^0x[0-9a-f]{64}$,minItems=1"`
	// must be greater than 0
	MaxFrequencyMs int `json:"maxFrequencyMs" jsonschema:"minimum=1"`
}

type inputs struct {
	TriggerID string `json:"triggerId"`
}

var mercuryTriggerValidator = capabilities.NewValidator[config, inputs, capabilities.TriggerEvent](capabilities.ValidatorArgs{Info: capInfo})

// This Trigger Service allows for the registration and deregistration of triggers. You can also send reports to the service.
type MercuryTriggerService struct {
	capabilities.Validator[config, inputs, capabilities.TriggerEvent]
	capabilities.CapabilityInfo
	tickerResolutionMs int64
	subscribers        map[string]*subscriber
	latestReports      map[mercury.FeedID]mercury.FeedReport
	mu                 sync.Mutex
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	lggr               logger.Logger
}

var _ capabilities.TriggerCapability = (*MercuryTriggerService)(nil)
var _ services.Service = &MercuryTriggerService{}

type subscriber struct {
	ch         chan<- capabilities.CapabilityResponse
	workflowID string
	config     config
}

// Mercury Trigger will send events to each subscriber every MaxFrequencyMs (configurable per subscriber).
// Event generation happens whenever local unix time is a multiple of tickerResolutionMs. Therefore,
// all subscribers' MaxFrequencyMs values need to be a multiple of tickerResolutionMs.
func NewMercuryTriggerService(tickerResolutionMs int64, lggr logger.Logger) *MercuryTriggerService {
	if tickerResolutionMs == 0 {
		tickerResolutionMs = defaultTickerResolutionMs
	}
	return &MercuryTriggerService{
		Validator:          mercuryTriggerValidator,
		CapabilityInfo:     capInfo,
		tickerResolutionMs: tickerResolutionMs,
		subscribers:        make(map[string]*subscriber),
		latestReports:      make(map[mercury.FeedID]mercury.FeedReport),
		stopCh:             make(services.StopChan),
		lggr:               logger.Named(lggr, "MercuryTriggerService")}
}

func (o *MercuryTriggerService) ProcessReport(reports []mercury.FeedReport) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lggr.Debugw("ProcessReport", "nReports", len(reports))
	for _, report := range reports {
		feedID := mercury.FeedID(report.FeedID)
		o.latestReports[feedID] = report
	}
	return nil
}

func (o *MercuryTriggerService) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	config, err := o.ValidateConfig(req.Config)
	if err != nil {
		return nil, err
	}

	inputs, err := o.ValidateInputs(req.Inputs)
	if err != nil {
		return nil, err
	}

	triggerID := o.getTriggerID(inputs.TriggerID, wid)
	// If triggerId is already registered, return an error
	if _, ok := o.subscribers[triggerID]; ok {
		return nil, fmt.Errorf("triggerId %s already registered", triggerID)
	}

	if int64(config.MaxFrequencyMs)%o.tickerResolutionMs != 0 {
		return nil, fmt.Errorf("MaxFrequencyMs must be a multiple of %d", o.tickerResolutionMs)
	}

	ch := make(chan capabilities.CapabilityResponse, defaultSendChannelBufferSize)
	o.subscribers[triggerID] =
		&subscriber{
			ch:         ch,
			workflowID: wid,
			config:     *config,
		}
	return ch, nil
}

func (o *MercuryTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	inputs, err := o.ValidateInputs(req.Inputs)
	if err != nil {
		return err
	}
	triggerID := o.getTriggerID(inputs.TriggerID, wid)

	subscriber, ok := o.subscribers[triggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", triggerID)
	}
	close(subscriber.ch)
	delete(o.subscribers, triggerID)
	return nil
}

func (o *MercuryTriggerService) getTriggerID(triggerID string, wid string) string {
	tid := wid + "|" + triggerID
	return tid
}

func (o *MercuryTriggerService) loop() {
	defer o.wg.Done()
	now := time.Now().UnixMilli()
	nextWait := o.tickerResolutionMs - now%o.tickerResolutionMs

	for {
		select {
		case <-o.stopCh:
			return
		case <-time.After(time.Duration(nextWait) * time.Millisecond):
			startTs := time.Now().UnixMilli()
			// find closest timestamp that is a multiple of o.tickerResolutionMs
			aligned := (startTs + o.tickerResolutionMs/2) / o.tickerResolutionMs * o.tickerResolutionMs
			o.process(aligned)
			endTs := time.Now().UnixMilli()
			if endTs-startTs > o.tickerResolutionMs {
				o.lggr.Errorw("processing took longer than ticker resolution", "duration", endTs-startTs, "tickerResolutionMs", o.tickerResolutionMs)
			}
			nextWait = getNextWaitIntervalMs(aligned, o.tickerResolutionMs, endTs)
		}
	}
}

func getNextWaitIntervalMs(lastTs, tickerResolutionMs, currentTs int64) int64 {
	desiredNext := lastTs + tickerResolutionMs
	nextWait := desiredNext - currentTs
	if nextWait <= 0 {
		nextWait = 0
	}
	return nextWait
}

func (o *MercuryTriggerService) process(timestamp int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, sub := range o.subscribers {
		if timestamp%int64(sub.config.MaxFrequencyMs) == 0 {
			reportList := make([]mercury.FeedReport, 0)
			for _, feedID := range sub.config.FeedIDs {
				if latest, ok := o.latestReports[mercury.FeedID(feedID)]; ok {
					reportList = append(reportList, latest)
				}
			}

			// use 32-byte-padded timestamp as EventID (human-readable)
			eventID := fmt.Sprintf("streams_%024s", strconv.FormatInt(timestamp, 10))
			capabilityResponse, err := wrapReports(reportList, eventID, timestamp)
			if err != nil {
				o.lggr.Errorw("error wrapping reports", "err", err)
				continue
			}

			o.lggr.Debugw("ProcessReport pushing event", "nReports", len(reportList), "eventID", eventID)
			select {
			case sub.ch <- capabilityResponse:
			default:
				o.lggr.Errorw("subscriber channel full, dropping event", "eventID", eventID, "workflowID", sub.workflowID)
			}
		}
	}
}

func wrapReports(reportList []mercury.FeedReport, eventID string, timestamp int64) (capabilities.CapabilityResponse, error) {
	val, err := values.Wrap(reportList)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: "mercury",
		ID:          eventID,
		Timestamp:   strconv.FormatInt(timestamp, 10),
		Payload:     val,
	}

	eventVal, err := values.Wrap(triggerEvent)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	// Create a new CapabilityResponse with the MercuryTriggerEvent
	return capabilities.CapabilityResponse{
		Value: eventVal,
	}, nil
}

func (o *MercuryTriggerService) Start(ctx context.Context) error {
	o.wg.Add(1)
	go o.loop()
	o.lggr.Info("MercuryTriggerService started")
	return nil
}

func (o *MercuryTriggerService) Close() error {
	close(o.stopCh)
	o.wg.Wait()
	o.lggr.Info("MercuryTriggerService closed")
	return nil
}

func (o *MercuryTriggerService) Ready() error {
	return nil
}

func (o *MercuryTriggerService) HealthReport() map[string]error {
	return nil
}

func (o *MercuryTriggerService) Name() string {
	return "MercuryTriggerService"
}
