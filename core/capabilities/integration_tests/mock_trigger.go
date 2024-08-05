package integration_tests

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const triggerID = "streams-trigger@1.0.0"

type reportsSink struct {
	services.StateMachine
	triggers []streamsTrigger

	nextEventID int

	sentResponses []capabilities.CapabilityResponse

	stopCh services.StopChan
	wg     sync.WaitGroup

	mux sync.Mutex
}

func newReportsSink() *reportsSink {
	return &reportsSink{
		stopCh: make(services.StopChan),
	}
}

func (r *reportsSink) Start(ctx context.Context) error {
	return r.StartOnce("reportsSink", func() error {
		return nil
	})
}

func (r *reportsSink) Close() error {
	return r.StopOnce("reportsSink", func() error {
		close(r.stopCh)
		r.wg.Wait()
		return nil
	})
}

func (r *reportsSink) sendReports(reportList []*datastreams.FeedReport) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.nextEventID++

	resp, err := wrapReports(reportList, strconv.Itoa(r.nextEventID), time.Now().UnixMilli(), datastreams.SignersMetadata{})
	if err != nil {
		panic(fmt.Errorf("failed to wrap reports: %w", err))
	}
	r.sentResponses = append(r.sentResponses, resp)

	for _, trigger := range r.triggers {
		if err := sendResponse(trigger, resp); err != nil {
			panic(fmt.Errorf("failed to send response: %w", err))
		}
	}
}

func (r *reportsSink) getNewTrigger(t *testing.T) *streamsTrigger {
	r.mux.Lock()
	defer r.mux.Unlock()
	trigger := streamsTrigger{t: t, toSend: make(chan capabilities.CapabilityResponse, 1000),
		wg: &r.wg, stopCh: r.stopCh}
	r.triggers = append(r.triggers, trigger)

	for _, resp := range r.sentResponses {
		if err := sendResponse(trigger, resp); err != nil {
			panic(fmt.Errorf("failed to send response: %w", err))
		}
	}

	return &trigger
}

func sendResponse(trigger streamsTrigger, response capabilities.CapabilityResponse) error {
	// clone response before sending so each trigger has its own instance to avoid any cross trigger mutation
	marshalledResponse, err := pb.MarshalCapabilityResponse(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	responseCopy, err := pb.UnmarshalCapabilityResponse(marshalledResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	trigger.sendResponse(responseCopy)

	return nil
}

type streamsTrigger struct {
	t      *testing.T
	cancel context.CancelFunc
	toSend chan capabilities.CapabilityResponse

	wg     *sync.WaitGroup
	stopCh services.StopChan
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
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		select {
		case <-s.stopCh:
			return
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
