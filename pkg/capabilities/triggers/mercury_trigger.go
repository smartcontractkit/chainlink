package triggers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var mercuryInfo = capabilities.MustNewCapabilityInfo(
	"mercury-trigger",
	capabilities.CapabilityTypeTrigger,
	"An example mercury trigger.",
	"v1.0.0",
)

// This Trigger Service allows for the registration and deregistration of triggers. You can also send reports to the service.
type MercuryTriggerService struct {
	capabilities.CapabilityInfo
	chans               map[string]chan<- capabilities.CapabilityResponse
	feedIDsForTriggerID map[string][]mercury.FeedID
	triggerIDsForFeedID map[mercury.FeedID]map[string]bool
	mu                  sync.Mutex
}

var _ capabilities.TriggerCapability = (*MercuryTriggerService)(nil)

func NewMercuryTriggerService() *MercuryTriggerService {
	return &MercuryTriggerService{
		CapabilityInfo:      mercuryInfo,
		chans:               map[string]chan<- capabilities.CapabilityResponse{},
		feedIDsForTriggerID: make(map[string][]mercury.FeedID),
		triggerIDsForFeedID: make(map[mercury.FeedID]map[string]bool),
	}
}

type FeedReport struct {
	FeedID               [mercury.FeedIDBytesLen]byte `json:"feedId"`
	FullReport           []byte                       `json:"fullReport"`
	BenchmarkPrice       int64                        `json:"benchmarkPrice"`
	ObservationTimestamp int64                        `json:"observationTimestamp"`
}

func (o *MercuryTriggerService) ProcessReport(reports []FeedReport) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	currentTime := time.Now()
	unixTimestampMillis := currentTime.UnixNano() / int64(time.Millisecond)
	triggerIDsToReports := make(map[string][]int)

	for reportIndex, report := range reports {
		for triggerID := range o.triggerIDsForFeedID[mercury.FeedIDFromBytes(report.FeedID)] {
			// if its not initialized, initialize it
			if _, ok := triggerIDsToReports[triggerID]; !ok {
				triggerIDsToReports[triggerID] = make([]int, 0)
			}
			triggerIDsToReports[triggerID] = append(triggerIDsToReports[triggerID], reportIndex)
		}
	}

	// Then for each trigger id, find which reports correspond to that trigger and create an event bundling the reports
	// and send it to the channel associated with the trigger id.
	for triggerID, reportIDs := range triggerIDsToReports {
		reportList := make([]mercury.FeedReport, 0)
		reportMap := make(map[string]any)
		for _, reportID := range reportIDs {
			rep := reports[reportID]
			feedID := mercury.FeedIDFromBytes(rep.FeedID)
			mercRep := mercury.FeedReport{
				FeedID:               feedID.String(),
				FullReport:           rep.FullReport,
				BenchmarkPrice:       rep.BenchmarkPrice,
				ObservationTimestamp: rep.ObservationTimestamp,
			}
			reportList = append(reportList, mercRep)
			reportMap[feedID.String()] = mercRep
		}

		triggerEvent := capabilities.TriggerEvent{
			TriggerType:    "mercury",
			ID:             GenerateTriggerEventID(reportList),
			Timestamp:      strconv.FormatInt(unixTimestampMillis, 10),
			BatchedPayload: reportMap,
		}

		val, err := mercury.Codec{}.WrapMercuryTriggerEvent(triggerEvent)
		if err != nil {
			return err
		}

		// Create a new CapabilityResponse with the MercuryTriggerEvent
		capabilityResponse := capabilities.CapabilityResponse{
			Value: val,
		}

		ch, ok := o.chans[triggerID]
		if !ok {
			return fmt.Errorf("no registration for %s", triggerID)
		}
		ch <- capabilityResponse
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
	if _, ok := o.chans[triggerID]; ok {
		return fmt.Errorf("triggerId %s already registered", triggerID)
	}

	feedIDs, err := o.GetFeedIDs(req)
	if err != nil {
		return err
	}

	if len(feedIDs) == 0 {
		return errors.New("no feedIDs to register")
	}

	o.chans[triggerID] = callback
	o.feedIDsForTriggerID[triggerID] = feedIDs
	for _, feedID := range feedIDs {
		// check if we need to initialize the map first
		if _, ok := o.triggerIDsForFeedID[feedID]; !ok {
			o.triggerIDsForFeedID[feedID] = make(map[string]bool)
		}
		o.triggerIDsForFeedID[feedID][triggerID] = true
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

	if _, ok := o.chans[triggerID]; !ok {
		return fmt.Errorf("triggerId %s not registered", triggerID)
	}

	ch, ok := o.chans[triggerID]
	if ok {
		close(ch)
	}

	for _, feedID := range o.feedIDsForTriggerID[triggerID] {
		delete(o.triggerIDsForFeedID[feedID], triggerID)
	}

	delete(o.chans, triggerID)
	delete(o.feedIDsForTriggerID, triggerID)

	return nil
}

// Get array of feedIds from CapabilityRequest req
func (o *MercuryTriggerService) GetFeedIDs(req capabilities.CapabilityRequest) ([]mercury.FeedID, error) {
	feedIDs := make([]mercury.FeedID, 0)
	// Unwrap the inputs which should return pair (map, nil) and then get the feedIds from the map
	if config, err := req.Config.Unwrap(); err == nil {
		if feeds, ok := config.(map[string]interface{})["feedIds"].([]any); ok {
			// Copy to feedIds
			for _, feed := range feeds {
				if id, ok := feed.(string); ok {
					mfid, err := mercury.NewFeedID(id)
					if err != nil {
						return nil, err
					}
					feedIDs = append(feedIDs, mfid)
				}
			}
		}
	}

	return feedIDs, nil
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

func GenerateTriggerEventID(reports []mercury.FeedReport) string {
	// Let's hash all the feedIds and timestamps together
	sort.Slice(reports, func(i, j int) bool {
		if reports[i].FeedID == reports[j].FeedID {
			return reports[i].ObservationTimestamp < reports[j].ObservationTimestamp
		}
		return reports[i].FeedID < reports[j].FeedID
	})
	s := ""
	for _, report := range reports {
		s += report.FeedID + strconv.FormatInt(report.ObservationTimestamp, 10) + ","
	}
	return sha256Hash(s)
}

func ValidateInput(mercuryTriggerEvent values.Value) error {
	// TODO: Fill this in
	return nil
}

func ExampleOutput() (values.Value, error) {
	feedOne := "0x111111111111111111110000000000000000000000000000000000000000"
	feedTwo := "0x222222222222222222220000000000000000000000000000000000000000"

	feeds := map[string]any{
		feedOne: mercury.FeedReport{
			FeedID:               feedOne,
			FullReport:           []byte("hello"),
			BenchmarkPrice:       100,
			ObservationTimestamp: 123,
		},
		feedTwo: mercury.FeedReport{
			FeedID:               feedTwo,
			FullReport:           []byte("world"),
			BenchmarkPrice:       100,
			ObservationTimestamp: 123,
		},
	}
	event := capabilities.TriggerEvent{
		TriggerType:    "mercury",
		ID:             "123",
		Timestamp:      "2024-01-17T04:00:10Z",
		BatchedPayload: feeds,
	}
	return mercury.Codec{}.WrapMercuryTriggerEvent(event)
}

func ValidateConfig(config values.Value) error {
	// TODO: Fill this in
	return nil
}

func sha256Hash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
