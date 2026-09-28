package triggers

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// validWorkflowID is a syntactically valid workflow ID: a 64-character hex string (32 bytes).
var validWorkflowID = strings.Repeat("ab", 32)

func TestRegistrationID(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got := RegistrationID(tt.workflowID, tt.triggerIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseWorkflowID(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
	const (
		triggerCapID          = "trigger-cap-id"
		triggerRegistrationID = "trigger-reg-id"
		eventID               = "event-id"
		method                = "test-method"
	)

	t.Run("returns an error when the handle is not found", func(t *testing.T) {
		t.Parallel()
		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), triggerRegistrationID)
	})

	t.Run("delegates to the handle's AckEvent using its stored method", func(t *testing.T) {
		t.Parallel()
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().AckEvent(mock.Anything, triggerRegistrationID, eventID, method).Return(nil).Once()
		handle := &Handle{TriggerCapability: trigger, Method: method}

		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, handle)
		require.NoError(t, err)
	})

	t.Run("propagates the handle's AckEvent error", func(t *testing.T) {
		t.Parallel()
		wantErr := errors.New("boom")
		trigger := capmocks.NewTriggerCapability(t)
		trigger.EXPECT().AckEvent(mock.Anything, triggerRegistrationID, eventID, method).Return(wantErr).Once()
		handle := &Handle{TriggerCapability: trigger, Method: method}

		err := Ack(t.Context(), logger.Test(t), newTestMetrics(t), triggerCapID, triggerRegistrationID, eventID, handle)
		require.ErrorIs(t, err, wantErr)
	})
}

func TestReadLoop(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
			lggr := logger.Test(t)
			testMetrics := newTestMetrics(t)
			done := make(chan struct{})
			go func() {
				defer close(done)
				ReadLoop(sutCtx, lggr, testMetrics, clock, workflowID, triggerCapID, triggerIndex, triggerEventCh, deliver)
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

func newRegisterDeps(t *testing.T, capReg registry.CapabilitiesRegistry, chainAllowed limits.GateLimiter) RegisterDeps {
	t.Helper()
	return RegisterDeps{
		CapRegistry:  capReg,
		RegTimeout:   limits.NewTimeLimiter(5 * time.Second),
		ChainAllowed: chainAllowed,
		Logger:       logger.Test(t),
		Metrics:      newTestMetrics(t),
	}
}

// erroringTimeLimiter is a limits.TimeLimiter whose WithTimeout always fails,
// for exercising Register's early-return path when the registration timeout
// itself can't be established.
type erroringTimeLimiter struct{}

func (erroringTimeLimiter) Close() error { return nil }
func (erroringTimeLimiter) Limit(context.Context) (time.Duration, error) {
	return 0, nil
}
func (erroringTimeLimiter) WithTimeout(context.Context) (context.Context, func(), error) {
	return nil, nil, errors.New("timeout limiter broken")
}

func TestRegister(t *testing.T) {
	t.Parallel()
	const (
		workflowID    = "wf-id"
		workflowOwner = "wf-owner"
		workflowDonID = uint32(7)
		donConfigVer  = uint32(3)
	)

	baseMeta := RegisterMetadata{
		WorkflowID:               workflowID,
		WorkflowOwner:            workflowOwner,
		WorkflowDonID:            workflowDonID,
		WorkflowDonConfigVersion: donConfigVer,
	}

	t.Run("registers all subscriptions and returns handles, capability IDs, and event channels in subscription order", func(t *testing.T) {
		t.Parallel()
		sub0 := &sdkpb.TriggerSubscription{Id: "trigger-a@1.0.0", Method: "method-a"}
		sub1 := &sdkpb.TriggerSubscription{Id: "trigger-b@1.0.0", Method: "method-b"}

		ch0 := make(chan capabilities.TriggerResponse)
		ch1 := make(chan capabilities.TriggerResponse)

		var mu sync.Mutex
		gotMetadata := map[string]capabilities.RequestMetadata{}
		captureMetadata := func(triggerCapID string) func(context.Context, capabilities.TriggerRegistrationRequest) {
			return func(_ context.Context, req capabilities.TriggerRegistrationRequest) {
				mu.Lock()
				gotMetadata[triggerCapID] = req.Metadata
				mu.Unlock()
			}
		}

		trigger0 := capmocks.NewTriggerCapability(t)
		trigger0.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Run(captureMetadata(sub0.Id)).
			Return((<-chan capabilities.TriggerResponse)(ch0), nil).Once()
		trigger1 := capmocks.NewTriggerCapability(t)
		trigger1.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Run(captureMetadata(sub1.Id)).
			Return((<-chan capabilities.TriggerResponse)(ch1), nil).Once()

		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetTrigger(mock.Anything, sub0.Id).Return(trigger0, nil).Once()
		capReg.EXPECT().GetTrigger(mock.Anything, sub1.Id).Return(trigger1, nil).Once()

		deps := newRegisterDeps(t, capReg, limits.NewGateLimiter(true))
		triggerCapIDs, handles, eventChans, err := Register(t.Context(), deps, baseMeta, []*sdkpb.TriggerSubscription{sub0, sub1})
		require.NoError(t, err)

		assert.Equal(t, []string{sub0.Id, sub1.Id}, triggerCapIDs)
		require.Len(t, eventChans, 2)
		assert.Equal(t, (<-chan capabilities.TriggerResponse)(ch0), eventChans[0])
		assert.Equal(t, (<-chan capabilities.TriggerResponse)(ch1), eventChans[1])

		reg0, reg1 := RegistrationID(workflowID, 0), RegistrationID(workflowID, 1)
		require.Contains(t, handles, reg0)
		require.Contains(t, handles, reg1)
		assert.Equal(t, sub0.Method, handles[reg0].Method)
		assert.Equal(t, sub1.Method, handles[reg1].Method)

		mu.Lock()
		defer mu.Unlock()
		md0 := gotMetadata[sub0.Id]
		assert.Equal(t, workflowID, md0.WorkflowID)
		assert.Equal(t, workflowOwner, md0.WorkflowOwner)
		assert.Equal(t, workflowDonID, md0.WorkflowDonID)
		assert.Equal(t, donConfigVer, md0.WorkflowDonConfigVersion)
		assert.Equal(t, "trigger_0", md0.ReferenceID)
		assert.Empty(t, md0.OrgID, "OrgID propagation is off by default")
	})

	t.Run("returns an error without registering when the chain gate is closed", func(t *testing.T) {
		t.Parallel()
		sub := &sdkpb.TriggerSubscription{Id: "trigger-a:ChainSelector_12345@1.0.0", Method: "method-a"}
		// No GetTrigger expectation: the gate must deny before the registry is consulted.
		capReg := regmocks.NewCapabilitiesRegistry(t)
		deps := newRegisterDeps(t, capReg, limits.NewGateLimiter(false))

		_, _, _, err := Register(t.Context(), deps, baseMeta, []*sdkpb.TriggerSubscription{sub})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ChainSelector")
	})

	t.Run("returns an error when a trigger capability isn't found in the registry", func(t *testing.T) {
		t.Parallel()
		sub := &sdkpb.TriggerSubscription{Id: "trigger-a@1.0.0", Method: "method-a"}
		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetTrigger(mock.Anything, sub.Id).Return(nil, errors.New("not found")).Once()
		deps := newRegisterDeps(t, capReg, limits.NewGateLimiter(true))

		_, _, _, err := Register(t.Context(), deps, baseMeta, []*sdkpb.TriggerSubscription{sub})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "trigger capability not found")
	})

	t.Run("unregisters already-succeeded handles when one registration fails", func(t *testing.T) {
		t.Parallel()
		sub0 := &sdkpb.TriggerSubscription{Id: "trigger-a@1.0.0", Method: "method-a"}
		sub1 := &sdkpb.TriggerSubscription{Id: "trigger-b@1.0.0", Method: "method-b"}

		ch0 := make(chan capabilities.TriggerResponse)
		trigger0 := capmocks.NewTriggerCapability(t)
		trigger0.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Return((<-chan capabilities.TriggerResponse)(ch0), nil).Once()
		trigger0.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

		trigger1 := capmocks.NewTriggerCapability(t)
		trigger1.EXPECT().RegisterTrigger(mock.Anything, mock.Anything).
			Return(nil, errors.New("registration failed")).Once()

		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetTrigger(mock.Anything, sub0.Id).Return(trigger0, nil).Once()
		capReg.EXPECT().GetTrigger(mock.Anything, sub1.Id).Return(trigger1, nil).Once()

		deps := newRegisterDeps(t, capReg, limits.NewGateLimiter(true))
		triggerCapIDs, handles, eventChans, err := Register(t.Context(), deps, baseMeta, []*sdkpb.TriggerSubscription{sub0, sub1})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register trigger")
		assert.Nil(t, triggerCapIDs)
		assert.Nil(t, handles)
		assert.Nil(t, eventChans)
		// trigger0.UnregisterTrigger's expectation being satisfied (verified by
		// capmocks.NewTriggerCapability's t.Cleanup) confirms the rollback ran.
	})

	t.Run("returns an error when the registration timeout limiter fails", func(t *testing.T) {
		t.Parallel()
		sub := &sdkpb.TriggerSubscription{Id: "trigger-a@1.0.0", Method: "method-a"}
		trigger0 := capmocks.NewTriggerCapability(t)
		capReg := regmocks.NewCapabilitiesRegistry(t)
		capReg.EXPECT().GetTrigger(mock.Anything, sub.Id).Return(trigger0, nil).Once()

		deps := newRegisterDeps(t, capReg, limits.NewGateLimiter(true))
		deps.RegTimeout = erroringTimeLimiter{}

		_, _, _, err := Register(t.Context(), deps, baseMeta, []*sdkpb.TriggerSubscription{sub})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout limiter broken")
	})
}

func TestUnregister(t *testing.T) {
	t.Parallel()
	const (
		workflowID    = "wf-id"
		workflowDonID = uint32(9)
	)

	t.Run("unregisters every handle and returns zero failures on success", func(t *testing.T) {
		t.Parallel()
		trigger0 := capmocks.NewTriggerCapability(t)
		trigger0.EXPECT().UnregisterTrigger(mock.Anything, mock.MatchedBy(func(req capabilities.TriggerRegistrationRequest) bool {
			return req.TriggerID == "reg-0" && req.Metadata.WorkflowID == workflowID && req.Metadata.WorkflowDonID == workflowDonID
		})).Return(nil).Once()
		trigger1 := capmocks.NewTriggerCapability(t)
		trigger1.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(nil).Once()

		handles := map[string]*Handle{
			"reg-0": {TriggerCapability: trigger0, Method: "m0"},
			"reg-1": {TriggerCapability: trigger1, Method: "m1"},
		}

		failCount := Unregister(t.Context(), logger.Test(t), workflowID, workflowDonID, handles)
		assert.Equal(t, 0, failCount)
	})

	t.Run("counts and logs individual unregister failures", func(t *testing.T) {
		t.Parallel()
		trigger0 := capmocks.NewTriggerCapability(t)
		trigger0.EXPECT().UnregisterTrigger(mock.Anything, mock.Anything).Return(errors.New("boom")).Once()

		handles := map[string]*Handle{"reg-0": {TriggerCapability: trigger0, Method: "m0"}}

		failCount := Unregister(t.Context(), logger.Test(t), workflowID, workflowDonID, handles)
		assert.Equal(t, 1, failCount)
	})

	t.Run("returns zero for an empty handle map", func(t *testing.T) {
		t.Parallel()
		failCount := Unregister(t.Context(), logger.Test(t), workflowID, workflowDonID, map[string]*Handle{})
		assert.Equal(t, 0, failCount)
	})
}

func TestRegistrationID_ParseWorkflowID_RoundTrip(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			registrationID := RegistrationID(tt.workflowID, tt.triggerIndex)

			got, err := ParseWorkflowID(registrationID)
			require.NoError(t, err)
			assert.Equal(t, tt.workflowID, got)
		})
	}
}
