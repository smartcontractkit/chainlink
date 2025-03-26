package framework

import (
	"context"
	"sync"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// TriggerSink is a TriggerFactory implementation that sends output to all triggers that are created by it.
type TriggerSink struct {
	services.StateMachine
	triggerID   string
	triggerName string
	version     string

	triggers []fakeTrigger

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func NewTriggerSink(t *testing.T, triggerName string, version string) *TriggerSink {
	triggersFactory := &TriggerSink{
		triggerID:   triggerName + "@" + version,
		triggerName: triggerName,
		version:     version,
		stopCh:      make(services.StopChan),
	}
	servicetest.Run(t, triggersFactory)
	return triggersFactory
}

func (r *TriggerSink) GetTriggerVersion() string {
	return r.version
}

func (r *TriggerSink) GetTriggerName() string {
	return r.triggerName
}

func (r *TriggerSink) GetTriggerID() string {
	return r.triggerID
}

func (r *TriggerSink) Start(ctx context.Context) error {
	return r.StartOnce("TriggerSink", func() error {
		return nil
	})
}

func (r *TriggerSink) Close() error {
	return r.StopOnce("TriggerSink", func() error {
		close(r.stopCh)
		r.wg.Wait()
		return nil
	})
}

// SendOutput wraps the given output in a TriggerEvent and sends it to all triggers created by this factory
func (r *TriggerSink) SendOutput(outputs *values.Map, eventID string) {
	triggerEvent := capabilities.TriggerEvent{
		TriggerType: r.triggerID,
		ID:          eventID,
		Outputs:     outputs,
	}

	resp := capabilities.TriggerResponse{
		Event: triggerEvent,
	}

	for _, trigger := range r.triggers {
		trigger.sendResponse(resp)
	}
}

func (r *TriggerSink) CreateNewTrigger(t *testing.T) capabilities.TriggerCapability {
	trigger := newFakeTrigger(t, r.triggerID, &r.wg, r.stopCh)
	r.triggers = append(r.triggers, trigger)
	return &trigger
}

type fakeTrigger struct {
	t         *testing.T
	triggerID string
	cancel    context.CancelFunc
	toSend    chan capabilities.TriggerResponse

	wg     *sync.WaitGroup
	stopCh services.StopChan
}

func newFakeTrigger(t *testing.T, triggerID string, wg *sync.WaitGroup, stopCh services.StopChan) fakeTrigger {
	return fakeTrigger{
		t:         t,
		triggerID: triggerID,
		toSend:    make(chan capabilities.TriggerResponse, 1000),
		wg:        wg,
		stopCh:    stopCh,
	}
}

func (s *fakeTrigger) sendResponse(resp capabilities.TriggerResponse) {
	s.toSend <- resp
}

func (s *fakeTrigger) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		s.triggerID,
		capabilities.CapabilityTypeTrigger,
		"fake trigger for trigger id "+s.triggerID,
	), nil
}

func (s *fakeTrigger) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
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

func (s *fakeTrigger) UnregisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) error {
	if s.cancel == nil {
		s.t.Fatal("trigger not registered")
	}

	s.cancel()
	s.cancel = nil
	return nil
}
