package integration_tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const triggerID = "streams-trigger@1.0.0"

type reportsSink struct {
	triggers []streamsTrigger
}

func (r *reportsSink) sendReports(reportList []*datastreams.FeedReport) {
	for _, trigger := range r.triggers {
		resp, err := wrapReports(reportList, "1", 12, datastreams.SignersMetadata{})
		if err != nil {
			panic(err)
		}
		trigger.sendResponse(resp)
	}
}

func (r *reportsSink) getNewTrigger(t *testing.T) *streamsTrigger {
	trigger := streamsTrigger{t: t, toSend: make(chan capabilities.CapabilityResponse, 1000)}
	r.triggers = append(r.triggers, trigger)
	return &trigger
}

type streamsTrigger struct {
	t      *testing.T
	cancel context.CancelFunc
	toSend chan capabilities.CapabilityResponse
}

func (s *streamsTrigger) sendResponse(resp capabilities.CapabilityResponse) {
	s.toSend <- resp
}

func (s *streamsTrigger) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		triggerID,
		capabilities.CapabilityTypeTrigger,
		"issues a trigger when a report is received.",
	), nil
}

func (s *streamsTrigger) RegisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if s.cancel != nil {
		s.t.Fatal("trigger already registered")
	}

	responseCh := make(chan capabilities.CapabilityResponse)

	ctxWithCancel, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go func() {
		select {
		case <-ctxWithCancel.Done():
			return
		case resp := <-s.toSend:
			responseCh <- resp
		}
	}()

	return responseCh, nil
}

func (s *streamsTrigger) UnregisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) error {
	if s.cancel == nil {
		s.t.Fatal("trigger not registered")
	}

	s.cancel()
	s.cancel = nil
	return nil
}

func wrapReports(reportList []*datastreams.FeedReport, eventID string, timestamp int64, meta datastreams.SignersMetadata) (capabilities.CapabilityResponse, error) {
	val, err := values.Wrap(reportList)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	metaVal, err := values.Wrap(meta)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: triggerID,
		ID:          eventID,
		Timestamp:   strconv.FormatInt(timestamp, 10),
		Metadata:    metaVal,
		Payload:     val,
	}

	triggerEventMapValue, err := values.WrapMap(triggerEvent)
	if err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("failed to wrap trigger event: %w", err)
	}

	// Create a new CapabilityResponse with the MercuryTriggerEvent
	return capabilities.CapabilityResponse{
		Value: triggerEventMapValue,
	}, nil
}
