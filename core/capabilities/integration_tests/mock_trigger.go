package integration_tests

import (
	"context"
	"sync"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const triggerID = "streams-trigger@1.0.0"

type reportsSink struct {
	services.StateMachine
	triggers []streamsTrigger

	stopCh services.StopChan
	wg     sync.WaitGroup
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
	for _, trigger := range r.triggers {
		resp, err := wrapReports(reportList, "1", 12, datastreams.Metadata{})
		if err != nil {
			panic(err)
		}
		trigger.sendResponse(resp)
	}
}

func (r *reportsSink) getNewTrigger(t *testing.T) *streamsTrigger {
	trigger := streamsTrigger{t: t, toSend: make(chan capabilities.TriggerResponse, 1000),
		wg: &r.wg, stopCh: r.stopCh}
	r.triggers = append(r.triggers, trigger)
	return &trigger
}

type streamsTrigger struct {
	t      *testing.T
	cancel context.CancelFunc
	toSend chan capabilities.TriggerResponse

	wg     *sync.WaitGroup
	stopCh services.StopChan
}

func (s *streamsTrigger) sendResponse(resp capabilities.TriggerResponse) {
	s.toSend <- resp
}

func (s *streamsTrigger) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		triggerID,
		capabilities.CapabilityTypeTrigger,
		"issues a trigger when a report is received.",
	), nil
}

func (s *streamsTrigger) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if s.cancel != nil {
		s.t.Fatal("trigger already registered")
	}

	responseCh := make(chan capabilities.TriggerResponse)

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.stopCh:
				return
			case <-ctxWithCancel.Done():
				return
			case resp := <-s.toSend:
				responseCh <- resp
			}
		}
	}()

	return responseCh, nil
}

func (s *streamsTrigger) UnregisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) error {
	if s.cancel == nil {
		s.t.Fatal("trigger not registered")
	}

	s.cancel()
	s.cancel = nil
	return nil
}

func wrapReports(reportList []*datastreams.FeedReport, eventID string, timestamp int64, meta datastreams.Metadata) (capabilities.TriggerResponse, error) {
	rl := []datastreams.FeedReport{}
	for _, r := range reportList {
		rl = append(rl, *r)
	}
	outputs, err := values.WrapMap(datastreams.StreamsTriggerEvent{
		Payload:   rl,
		Metadata:  meta,
		Timestamp: timestamp,
	})
	if err != nil {
		return capabilities.TriggerResponse{}, err
	}

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: triggerID,
		ID:          eventID,
		Outputs:     outputs,
	}

	// Create a new TriggerResponse with the MercuryTriggerEvent
	return capabilities.TriggerResponse{
		Event: triggerEvent,
	}, nil
}
