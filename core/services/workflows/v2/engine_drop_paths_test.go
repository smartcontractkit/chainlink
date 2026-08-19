package v2_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	workflowEvents "github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

const (
	v2ExecutionStartedEntity  = "workflows.v2." + workflowEvents.WorkflowExecutionStarted
	v2ExecutionFinishedEntity = "workflows.v2." + workflowEvents.WorkflowExecutionFinished
)

// failAfterNBoundLimiter succeeds (returning ok, nil) for the first n calls to Limit,
// then returns failErr for every call after that. LogEvent and ExecutionResponse are
// both peeked once during trigger-subscription init (before any execution exists) and
// again per execution, so an always-failing fake would break engine initialization
// before execution-time fallback could ever be exercised.
type failAfterNBoundLimiter[N limits.Number] struct {
	mu      sync.Mutex
	calls   int
	n       int
	ok      N
	failErr error
}

func (f *failAfterNBoundLimiter[N]) Limit(context.Context) (N, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.calls > f.n {
		var zero N
		return zero, f.failErr
	}
	return f.ok, nil
}

func (f *failAfterNBoundLimiter[N]) Check(context.Context, N) error { return nil }
func (f *failAfterNBoundLimiter[N]) Close() error                   { return nil }

// alwaysFailTimeLimiter fails every Limit/WithTimeout call. Safe to use unconditionally
// for ExecutionTime, which is not read during trigger-subscription init.
type alwaysFailTimeLimiter struct{ err error }

func (f *alwaysFailTimeLimiter) Limit(context.Context) (time.Duration, error) { return 0, f.err }

func (f *alwaysFailTimeLimiter) WithTimeout(context.Context) (context.Context, func(), error) {
	return nil, nil, f.err
}

func (f *alwaysFailTimeLimiter) Close() error { return nil }

// fakeShardResolver reports a fixed shard-ownership verdict, standing in for a real
// ring orchestrator lookup.
type fakeShardResolver struct {
	shardID uint32
	found   bool
	err     error
}

func (f *fakeShardResolver) ResolveShard(context.Context, string, string) (uint32, bool, error) {
	return f.shardID, f.found, f.err
}

func (f *fakeShardResolver) ResolveShards(context.Context, []string, []string) (map[string]uint32, error) {
	return nil, nil
}

// latestV2FinishedEvent returns the most recently emitted v2 WorkflowExecutionFinished
// event, or ok=false if none has been emitted yet.
func latestV2FinishedEvent(t *testing.T, observer beholdertest.Observer) (evt *eventsv2.WorkflowExecutionFinished, ok bool) {
	t.Helper()
	msgs := observer.Messages(t, "beholder_entity", v2ExecutionFinishedEntity)
	if len(msgs) == 0 {
		return nil, false
	}
	evt = &eventsv2.WorkflowExecutionFinished{}
	require.NoError(t, proto.Unmarshal(msgs[len(msgs)-1].Body, evt))
	return evt, true
}

// dropPathHarness is a minimal engine with a single trigger subscription (id_0)
// registered, ready to have a trigger event pushed through eventCh to drive
// startExecution. Callers customize cfg.LocalLimiters / cfg.ShardingEnabled / etc.
// via configure before the harness starts the engine, and set up the second
// Module.Execute expectation (for the per-execution call, if reached) themselves.
type dropPathHarness struct {
	module              *modulemocks.ModuleV2
	eventCh             chan capabilities.TriggerResponse
	executionFinishedCh chan string
	beholderObserver    beholdertest.Observer
}

func newDropPathHarness(t *testing.T, billingClient *metmocks.BillingClient, configure func(cfg *v2.EngineConfig)) *dropPathHarness {
	t.Helper()

	module := modulemocks.NewModuleV2(t)
	module.EXPECT().Start()
	module.EXPECT().Close()
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error, 1)
	subscribedToTriggersCh := make(chan []string, 1)
	executionFinishedCh := make(chan string, 1)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.BillingClient = billingClient
	cfg.OrgResolver = &mockOrgResolver{orgID: "test-org-123"}
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
		OnSubscribedToTriggers: func(triggerIDs []string) {
			subscribedToTriggersCh <- triggerIDs
		},
		OnExecutionFinished: func(executionID string, _ string) {
			executionFinishedCh <- executionID
		},
	}
	beholderObserver := beholdertest.NewObserver(t)
	cfg.BeholderEmitter = custmsg.NewLabeler()

	if configure != nil {
		configure(cfg)
	}

	// Trigger-subscription init always calls Module.Execute once for the Subscribe
	// request, regardless of what's under test below.
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(newTriggerSubs(1), nil).Once()
	trigger := capmocks.NewTriggerCapability(t)
	capreg.EXPECT().GetTrigger(matches.AnyContext, "id_0").Return(trigger, nil)
	eventCh := make(chan capabilities.TriggerResponse)
	trigger.EXPECT().RegisterTrigger(matches.AnyContext, mock.Anything).Return(eventCh, nil).Once()
	trigger.EXPECT().UnregisterTrigger(matches.AnyContext, mock.Anything).Return(nil).Once()
	// Reserve-failure drop paths return before the trigger event is ever
	// ACKed, so this must not be a hard requirement across all harness users.
	trigger.EXPECT().AckEvent(matches.AnyContext, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)
	require.NoError(t, engine.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, engine.Close()) })

	require.NoError(t, <-initDoneCh)
	require.Equal(t, []string{"id_0"}, <-subscribedToTriggersCh)

	return &dropPathHarness{
		module:              module,
		eventCh:             eventCh,
		executionFinishedCh: executionFinishedCh,
		beholderObserver:    beholderObserver,
	}
}

// TestEngine_LimitReadFallback_ExecutesAnyway covers: an execution proceeds to
// completion (instead of being silently dropped) when a limit read fails,
// falling back to the setting's static default.
func TestEngine_LimitReadFallback_ExecutesAnyway(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	testCases := map[string]func(cfg *v2.EngineConfig){
		"ExecutionTime": func(cfg *v2.EngineConfig) {
			cfg.LocalLimiters.ExecutionTime = &alwaysFailTimeLimiter{err: errors.New("limit read boom")}
		},
		"LogEvent": func(cfg *v2.EngineConfig) {
			cfg.LocalLimiters.LogEvent = &failAfterNBoundLimiter[int]{n: 1, ok: 1000, failErr: errors.New("limit read boom")}
		},
		"ExecutionResponse": func(cfg *v2.EngineConfig) {
			cfg.LocalLimiters.ExecutionResponse = &failAfterNBoundLimiter[config.Size]{n: 1, ok: config.Size(10 * 1024 * 1024), failErr: errors.New("limit read boom")}
		},
	}

	for name, breakLimiter := range testCases { //nolint:paralleltest // shares the package-level beholder singleton
		t.Run(name, func(t *testing.T) {
			billingClient := setupMockBillingClient(t)
			harness := newDropPathHarness(t, billingClient, breakLimiter)

			harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).Return(nil, nil).Once()
			harness.eventCh <- capabilities.TriggerResponse{
				Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_" + name},
			}

			executionID := <-harness.executionFinishedCh
			require.NotEmpty(t, executionID)

			startedMsgs := harness.beholderObserver.Messages(t, "beholder_entity", v2ExecutionStartedEntity)
			assert.NotEmpty(t, startedMsgs, "execution should still emit a Started event despite the limit read failure")

			evt, ok := latestV2FinishedEvent(t, harness.beholderObserver)
			require.True(t, ok, "execution should still emit a Finished event despite the limit read failure")
			assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED, evt.Status)
		})
	}
}

// setupInsufficientFundingBilling returns a billing mock whose ReserveCredits denies
// the reservation (Success: false), the only Reserve failure reachable in practice
// (every other Reserve error fails open inside Reserve itself). expectReceipt controls
// whether SubmitWorkflowReceipt (called from Reports.End) is asserted as called.
func setupInsufficientFundingBilling(t *testing.T, expectReceipt bool) *metmocks.BillingClient {
	t.Helper()
	billingClient := metmocks.NewBillingClient(t)
	billingClient.EXPECT().GetWorkflowExecutionRates(mock.Anything, mock.Anything).Return(&billing.GetWorkflowExecutionRatesResponse{
		RateCards: []*billing.RateCard{
			{
				ResourceType:    billing.ResourceType_RESOURCE_TYPE_COMPUTE,
				MeasurementUnit: billing.MeasurementUnit_MEASUREMENT_UNIT_MILLISECONDS,
				UnitsPerCredit:  "0.0001",
			},
		},
	}, nil)
	billingClient.EXPECT().ReserveCredits(mock.Anything, mock.Anything).Return(&billing.ReserveCreditsResponse{Success: false}, nil)
	receiptCall := billingClient.EXPECT().SubmitWorkflowReceipt(mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)
	if expectReceipt {
		receiptCall.Once()
	} else {
		receiptCall.Maybe()
	}
	return billingClient
}

// TestEngine_MeteringReserveFailure_SurfacesAsUserError covers: a metering
// Reserve failure (insufficient funding) used to drop the execution before the
// Started/Finished emit points were reached, so the failure never reached the UI.
// It now publishes the Started/Finished pair, attributed to the user (their own
// out-of-credits state), not the platform.
func TestEngine_MeteringReserveFailure_SurfacesAsUserError(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	billingClient := setupInsufficientFundingBilling(t, false)
	harness := newDropPathHarness(t, billingClient, nil)

	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_reserve_failure"},
	}

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		startedMsgs := harness.beholderObserver.Messages(t, "beholder_entity", v2ExecutionStartedEntity)
		assert.NotEmpty(c, startedMsgs)
	}, 5*time.Second, 50*time.Millisecond)

	var evt *eventsv2.WorkflowExecutionFinished
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		got, ok := latestV2FinishedEvent(t, harness.beholderObserver)
		if !assert.True(c, ok) {
			return
		}
		evt = got
	}, 5*time.Second, 50*time.Millisecond)

	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED, evt.Status)
	assert.NotEmpty(t, evt.Error)
	assert.Equal(t, eventsv2.ClassifiedExecutionStatus_CLASSIFIED_EXECUTION_STATUS_USER_ERROR, evt.ClassifiedStatus)
}

// TestEngine_MeteringReserveFailure_ReportReleased covers: before this fix,
// meterReports.End was only reached on the happy/normal-error path, so a Reserve
// failure leaked the report; a redelivered trigger event for the same executionID
// would then hit ErrReportExists and run unmetered. SubmitWorkflowReceipt (only
// reachable via Reports.End) being called proves the report was released.
func TestEngine_MeteringReserveFailure_ReportReleased(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	billingClient := setupInsufficientFundingBilling(t, true)
	harness := newDropPathHarness(t, billingClient, nil)

	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_report_released"},
	}

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		_, ok := latestV2FinishedEvent(t, harness.beholderObserver)
		assert.True(c, ok)
	}, 5*time.Second, 50*time.Millisecond)
	// billingClient's t.Cleanup (registered by metmocks.NewBillingClient) asserts
	// SubmitWorkflowReceipt was actually called.
}

// TestEngine_ShardDenial_StaysSilent is a regression guard: with sharding on, every
// node outside the owning shard denies every execution, so emitting Started/Finished
// there (unlike the other drop paths) would publish a DON-wide failure for a run
// that actually succeeded on the owner. This must stay silent.
func TestEngine_ShardDenial_StaysSilent(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	billingClient := metmocks.NewBillingClient(t) // no billing calls expected: shard check runs before metering

	harness := newDropPathHarness(t, billingClient, func(cfg *v2.EngineConfig) {
		cfg.ShardingEnabled = true
		cfg.MyShardID = 0
		cfg.ShardResolver = &fakeShardResolver{shardID: 1, found: true} // some other shard owns it
	})

	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_shard_denied"},
	}

	require.Never(t, func() bool {
		started := harness.beholderObserver.Messages(t, "beholder_entity", v2ExecutionStartedEntity)
		finished := harness.beholderObserver.Messages(t, "beholder_entity", v2ExecutionFinishedEntity)
		return len(started) > 0 || len(finished) > 0
	}, 500*time.Millisecond, 25*time.Millisecond)
}
