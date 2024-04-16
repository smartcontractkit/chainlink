package triggers

import (
	"context"
	"errors"
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

var mercuryInfo = capabilities.MustNewCapabilityInfo(
	"mercury-trigger",
	capabilities.CapabilityTypeTrigger,
	"An example mercury trigger.",
	"v1.0.0",
)

const defaultTickerResolutionMs = 1000

// This Trigger Service allows for the registration and deregistration of triggers. You can also send reports to the service.
type MercuryTriggerService struct {
	capabilities.CapabilityInfo
	tickerResolutionMs int64
	subscribers        map[string]*subscriber
	latestReports      map[mercury.FeedID]mercury.FeedReport
	mu                 sync.Mutex
	stopCh             services.StopChan
	wg                 sync.WaitGroup
	lggr               logger.Logger
}

var _ services.Service = &MercuryTriggerService{}

type subscriberConfig struct {
	FeedIds        []string
	MaxFrequencyMs int
}

type subscriber struct {
	ch         chan<- capabilities.CapabilityResponse
	workflowID string
	config     subscriberConfig
}

var _ capabilities.TriggerCapability = (*MercuryTriggerService)(nil)

// Mercury Trigger will send events to each subscriber every MaxFrequencyMs (configurable per subscriber).
// Event generation happens whenever local unix time is a multiple of tickerResolutionMs. Therefore,
// all subscribers' MaxFrequencyMs values need to be a multiple of tickerResolutionMs.
func NewMercuryTriggerService(tickerResolutionMs int64, lggr logger.Logger) *MercuryTriggerService {
	if tickerResolutionMs == 0 {
		tickerResolutionMs = defaultTickerResolutionMs
	}
	return &MercuryTriggerService{
		CapabilityInfo:     mercuryInfo,
		tickerResolutionMs: tickerResolutionMs,
		subscribers:        make(map[string]*subscriber),
		latestReports:      make(map[mercury.FeedID]mercury.FeedReport),
		stopCh:             make(services.StopChan),
		lggr:               logger.Named(lggr, "MercuryTriggerService"),
	}
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

func (o *MercuryTriggerService) RegisterTrigger(ctx context.Context, callback chan<- capabilities.CapabilityResponse, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	triggerID, err := o.GetTriggerID(req, wid)
	if err != nil {
		return err
	}

	// If triggerId is already registered, return an error
	if _, ok := o.subscribers[triggerID]; ok {
		return fmt.Errorf("triggerId %s already registered", triggerID)
	}

	cfg := subscriberConfig{}
	err = req.Config.UnwrapTo(&cfg)
	if err != nil {
		return err
	}

	feedIDs := []mercury.FeedID{}
	for _, feedID := range cfg.FeedIds {
		mfid, err := mercury.NewFeedID(feedID)
		if err != nil {
			return err
		}
		feedIDs = append(feedIDs, mfid)
	}

	if len(feedIDs) == 0 {
		return errors.New("no feedIDs to register")
	}

	if int64(cfg.MaxFrequencyMs)%o.tickerResolutionMs != 0 {
		return fmt.Errorf("MaxFrequencyMs must be a multiple of %d", o.tickerResolutionMs)
	}

	o.subscribers[triggerID] =
		&subscriber{
			ch:         callback,
			workflowID: wid,
			config:     cfg,
		}
	return nil
}

func (o *MercuryTriggerService) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	wid := req.Metadata.WorkflowID

	o.mu.Lock()
	defer o.mu.Unlock()

	triggerID, err := o.GetTriggerID(req, wid)
	if err != nil {
		return err
	}

	subscriber, ok := o.subscribers[triggerID]
	if !ok {
		return fmt.Errorf("triggerId %s not registered", triggerID)
	}
	close(subscriber.ch)
	delete(o.subscribers, triggerID)
	return nil
}

// Get the triggerId from the CapabilityRequest req map
func (o *MercuryTriggerService) GetTriggerID(req capabilities.CapabilityRequest, workflowID string) (string, error) {
	// Unwrap the inputs which should return pair (map, nil) and then get the triggerId from the map
	inputs, err := req.Inputs.Unwrap()
	if err != nil {
		return "", err
	}
	if id, ok := inputs.(map[string]interface{})["triggerId"].(string); ok {
		// TriggerIDs should be namespaced to the workflowID
		return workflowID + "|" + id, nil
	}

	return "", fmt.Errorf("triggerId not found in inputs")
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
			for _, feedID := range sub.config.FeedIds {
				if latest, ok := o.latestReports[mercury.FeedID(feedID)]; ok {
					reportList = append(reportList, latest)
				}
			}

			// use 32-byte-padded timestamp as EventID (human-readable)
			eventID := fmt.Sprintf("mercury_%024s", strconv.FormatInt(timestamp, 10))
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
	val, err := mercury.Codec{}.Wrap(reportList)
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

func ValidateInput(mercuryTriggerEvent values.Value) error {
	// TODO: Fill this in
	return nil
}

func ExampleOutput() (values.Value, error) {
	feedOne := "0x111111111111111111110000000000000000000000000000000000000000"
	feedTwo := "0x222222222222222222220000000000000000000000000000000000000000"

	reportSet := []mercury.FeedReport{
		{
			FeedID:               feedOne,
			FullReport:           []byte("hello"),
			BenchmarkPrice:       100,
			ObservationTimestamp: 123,
		},
		{
			FeedID:               feedTwo,
			FullReport:           []byte("world"),
			BenchmarkPrice:       100,
			ObservationTimestamp: 123,
		},
	}

	val, err := mercury.Codec{}.Wrap(reportSet)
	if err != nil {
		return val, err
	}

	event := capabilities.TriggerEvent{
		TriggerType: "mercury",
		ID:          "1712963290",
		Timestamp:   "1712963290",
		Payload:     val,
	}

	return values.Wrap(event)
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

func ValidateConfig(config values.Value) error {
	// TODO: Fill this in
	return nil
}
