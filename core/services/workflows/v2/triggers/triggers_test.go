package triggers

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// validWorkflowID is a syntactically valid workflow ID: a 64-character hex string (32 bytes).
var validWorkflowID = strings.Repeat("ab", 32)

func TestRegistrationID(t *testing.T) {
	tests := []struct {
		name         string
		workflowID   string
		triggerIndex int
		want         string
	}{
		{
			name:         "simple workflow ID",
			workflowID:   "abc123",
			triggerIndex: 0,
			want:         "trigger_reg_abc123_0",
		},
		{
			name:         "workflow ID containing underscores",
			workflowID:   "wf_with_underscores",
			triggerIndex: 5,
			want:         "trigger_reg_wf_with_underscores_5",
		},
		{
			name:         "empty workflow ID",
			workflowID:   "",
			triggerIndex: 2,
			want:         "trigger_reg__2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RegistrationID(tt.workflowID, tt.triggerIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseWorkflowID(t *testing.T) {
	tests := []struct {
		name             string
		registrationID   string
		wantWorkflowID   string
		wantErrSubstring string
	}{
		{
			name:           "valid workflow ID",
			registrationID: "trigger_reg_" + validWorkflowID + "_0",
			wantWorkflowID: validWorkflowID,
		},
		{
			name:           "valid workflow ID with large trigger index",
			registrationID: "trigger_reg_" + validWorkflowID + "_12345",
			wantWorkflowID: validWorkflowID,
		},
		{
			name:             "missing prefix",
			registrationID:   "not_a_registration_id_0",
			wantErrSubstring: "missing prefix",
		},
		{
			name:             "missing trigger index",
			registrationID:   "trigger_reg_" + validWorkflowID,
			wantErrSubstring: "missing trigger index",
		},
		{
			name:             "non-numeric trigger index",
			registrationID:   "trigger_reg_" + validWorkflowID + "_notanumber",
			wantErrSubstring: "invalid trigger index",
		},
		{
			name:             "empty workflow ID",
			registrationID:   "trigger_reg__2",
			wantErrSubstring: "invalid workflow ID",
		},
		{
			name:             "workflow ID too short",
			registrationID:   "trigger_reg_abc123_0",
			wantErrSubstring: "invalid workflow ID",
		},
		{
			name:             "workflow ID not hex-encoded",
			registrationID:   "trigger_reg_" + strings.Repeat("z", 64) + "_0",
			wantErrSubstring: "invalid workflow ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWorkflowID(tt.registrationID)
			if tt.wantErrSubstring != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrSubstring)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantWorkflowID, got)
		})
	}
}

func newTestMetrics(t *testing.T) *monitoring.WorkflowsMetricLabeler {
	t.Helper()
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	return monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
}

func TestAck(t *testing.T) {
	const (
		triggerCapID          = "trigger-cap-id"
		triggerRegistrationID = "trigger-reg-id"
		eventID               = "event-id"
		method                = "test-method"
	)

	t.Run("returns an error when the handle is not found", func(t *testing.T) {
		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), triggerRegistrationID)
	})

	t.Run("delegates to the handle's AckEvent using its stored method", func(t *testing.T) {
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().AckEvent(mock.Anything, triggerRegistrationID, eventID, method).Return(nil).Once()
		handle := &Handle{TriggerCapability: trigger, Method: method}

		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, handle)
		require.NoError(t, err)
	})

	t.Run("propagates the handle's AckEvent error", func(t *testing.T) {
		wantErr := errors.New("boom")
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().AckEvent(mock.Anything, triggerRegistrationID, eventID, method).Return(wantErr).Once()
		handle := &Handle{TriggerCapability: trigger, Method: method}

		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, handle)
		require.ErrorIs(t, err, wantErr)
	})
}

func TestReadLoop(t *testing.T) {
	const (
		workflowID   = "wf-id"
		triggerCapID = "trigger-cap-id"
		triggerIndex = 3
		chanBuf      = 4
	)

	tests := []struct {
		name string
		// wantDelivered are the event IDs ReadLoop is expected to hand to
		// deliver, in order. Nil means deliver must not be called at all.
		wantDelivered []string
		// run drives the scenario: it sends/closes on triggerEventCh and/or
		// cancels the reader's context.
		run func(triggerEventCh chan<- capabilities.TriggerResponse, cancel context.CancelFunc)
	}{
		{
			name:          "delivers a routed event built from the response and the clock",
			wantDelivered: []string{"event-1"},
			run: func(triggerEventCh chan<- capabilities.TriggerResponse, _ context.CancelFunc) {
				triggerEventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{ID: "event-1"}}
			},
		},
		{
			name:          "skips delivery when the response carries an error, but keeps consuming",
			wantDelivered: []string{"event-2"},
			run: func(triggerEventCh chan<- capabilities.TriggerResponse, _ context.CancelFunc) {
				triggerEventCh <- capabilities.TriggerResponse{Err: errors.New("boom")}
				triggerEventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{ID: "event-2"}}
			},
		},
		{
			name:          "delivers multiple events in order",
			wantDelivered: []string{"a", "b"},
			run: func(triggerEventCh chan<- capabilities.TriggerResponse, _ context.CancelFunc) {
				triggerEventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{ID: "a"}}
				triggerEventCh <- capabilities.TriggerResponse{Event: capabilities.TriggerEvent{ID: "b"}}
			},
		},
		{
			name: "returns when the trigger channel closes",
			run: func(triggerEventCh chan<- capabilities.TriggerResponse, _ context.CancelFunc) {
				close(triggerEventCh)
			},
		},
		{
			name: "returns when ctx is canceled",
			run: func(_ chan<- capabilities.TriggerResponse, cancel context.CancelFunc) {
				cancel()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A single deadline governs the whole subtest: every wait below
			// (deliveries, then ReadLoop's return) races against this one
			// timer instead of each getting its own fresh timeout.
			deadline, cancelDeadline := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelDeadline()

			sutCtx, sutCancel := context.WithCancel(deadline)
			defer sutCancel()

			triggerEventCh := make(chan capabilities.TriggerResponse, chanBuf)
			deliverCh := make(chan CoordinatedEvent, chanBuf)
			deliver := func(_ context.Context, event CoordinatedEvent) {
				deliverCh <- event
			}

			clock := clockwork.NewFakeClock()
			done := make(chan struct{})
			go func() {
				defer close(done)
				ReadLoop(sutCtx, logger.Test(t), newTestMetrics(t), clock, workflowID, triggerCapID, triggerIndex, triggerEventCh, deliver)
			}()

			tt.run(triggerEventCh, sutCancel)

			var gotDelivered []string
			for range tt.wantDelivered {
				select {
				case event := <-deliverCh:
					assert.Equal(t, workflowID, event.WorkflowID)
					assert.Equal(t, triggerCapID, event.TriggerCapID)
					assert.Equal(t, triggerIndex, event.TriggerIndex)
					assert.Equal(t, clock.Now(), event.ObservedAt)
					gotDelivered = append(gotDelivered, event.Event.Event.ID)
				case <-deadline.Done():
					t.Fatal("timed out waiting for delivery")
				}
			}
			assert.Equal(t, tt.wantDelivered, gotDelivered)

			// Whatever the scenario did to end ReadLoop (close/cancel), make
			// sure it actually returns before the shared deadline.
			sutCancel()
			select {
			case <-done:
			case <-deadline.Done():
				t.Fatal("ReadLoop did not return before the test deadline")
			}
		})
	}
}

func TestRegistrationID_ParseWorkflowID_RoundTrip(t *testing.T) {
	tests := []struct {
		name         string
		workflowID   string
		triggerIndex int
	}{
		{name: "trigger index zero", workflowID: validWorkflowID, triggerIndex: 0},
		{name: "large trigger index", workflowID: validWorkflowID, triggerIndex: 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registrationID := RegistrationID(tt.workflowID, tt.triggerIndex)

			got, err := ParseWorkflowID(registrationID)
			require.NoError(t, err)
			assert.Equal(t, tt.workflowID, got)
		})
	}
}
