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
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	pb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
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
	v1UserLogsEntity          = "workflows.v1." + workflowEvents.UserLogs
)

// failAfterNBoundLimiter succeeds (returning ok, nil) for the first n calls to Limit,
// then returns (ok, failErr) after that, mirroring the real chainlink-common
// BoundLimiter. LogEvent and ExecutionResponse are both peeked once during
// trigger-subscription init and again per execution, so an always-failing fake would
// break engine initialization before the value-on-error path could be exercised.
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
		return f.ok, f.failErr // resolved value, not zero; err is advisory
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

// alwaysErrCheckLimiter fails every Check with a plain error, standing in for a settings
// read failure rather than the bound actually being exceeded (which would be an
// ErrorBoundLimited). Used to exercise the fail-open branches in processLogLine.
//
// Limit deliberately still succeeds: LogEvent.Limit is peeked during trigger-subscription
// init to size the user-log channel, so failing it too would break engine startup before
// the Check path under test is ever reached.
type alwaysErrCheckLimiter[N limits.Number] struct {
	ok  N
	err error
}

func (f *alwaysErrCheckLimiter[N]) Limit(context.Context) (N, error) { return f.ok, nil }
func (f *alwaysErrCheckLimiter[N]) Check(context.Context, N) error   { return f.err }
func (f *alwaysErrCheckLimiter[N]) Close() error                     { return nil }

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

// userLogLines returns every user log line emitted so far via the v1 UserLogs event.
func userLogLines(t *testing.T, observer beholdertest.Observer) []string {
	t.Helper()
	var lines []string
	for _, msg := range observer.Messages(t, "beholder_entity", v1UserLogsEntity) {
		var payload pb.UserLogs
		if err := proto.Unmarshal(msg.Body, &payload); err != nil {
			continue
		}
		for _, l := range payload.LogLines {
			lines = append(lines, l.Message)
		}
	}
	return lines
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

// TestEngine_LimitReadFallback_UsesLimiterValue: on a read failure, the engine must use
// whatever value the limiter returns, not substitute its own compiled default.
func TestEngine_LimitReadFallback_UsesLimiterValue(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	const resolvedValue = config.Size(12345)

	harness := newDropPathHarness(t, setupMockBillingClient(t), func(cfg *v2.EngineConfig) {
		cfg.LocalLimiters.ExecutionResponse = &failAfterNBoundLimiter[config.Size]{n: 1, ok: resolvedValue, failErr: errors.New("limit read boom")}
	})

	var gotMaxResponseSize uint64
	harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(_ context.Context, req *sdkpb.ExecuteRequest, _ host.ExecutionHelper) {
			gotMaxResponseSize = req.MaxResponseSize
		}).
		Return(nil, nil).Once()
	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_limiter_value"},
	}

	executionID := <-harness.executionFinishedCh
	require.NotEmpty(t, executionID)

	assert.Equal(t, uint64(resolvedValue), gotMaxResponseSize,
		"engine must use the value the limiter returned, not fall back to its own compiled default")
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

// TestEngine_UserLog_LogLineCheck covers the LogLine Limit()->Check() swap: the bound is
// now enforced by the limiter (so it records usage/denied), truncation is driven by the
// bound carried on the returned ErrorBoundLimited, and a settings read failure fails open
// by emitting the line untruncated instead of falling back to a hardcoded default.
func TestEngine_UserLog_LogLineCheck(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	t.Run("over-long message truncated to the bound", func(t *testing.T) { //nolint:paralleltest // shares the package-level beholder singleton
		harness := newDropPathHarness(t, setupMockBillingClient(t), func(cfg *v2.EngineConfig) {
			cfg.LocalLimiters.LogLine = limits.NewUpperBoundLimiter[config.Size](10)
		})

		harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(_ context.Context, _ *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				require.NoError(t, executor.EmitUserLog("0123456789ABCDEF"))
			}).
			Return(nil, nil).Once()
		harness.eventCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_truncated"},
		}

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Contains(c, userLogLines(t, harness.beholderObserver), "0123456789 ...(truncated)")
		}, 5*time.Second, 50*time.Millisecond)
	})

	t.Run("settings read failure emits untruncated", func(t *testing.T) { //nolint:paralleltest // shares the package-level beholder singleton
		harness := newDropPathHarness(t, setupMockBillingClient(t), func(cfg *v2.EngineConfig) {
			cfg.LocalLimiters.LogLine = &alwaysErrCheckLimiter[config.Size]{err: errors.New("settings unavailable")}
		})

		const message = "this message is not over any real bound"
		harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
			Run(func(_ context.Context, _ *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
				require.NoError(t, executor.EmitUserLog(message))
			}).
			Return(nil, nil).Once()
		harness.eventCh <- capabilities.TriggerResponse{
			Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_untruncated"},
		}

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Contains(c, userLogLines(t, harness.beholderObserver), message)
		}, 5*time.Second, 50*time.Millisecond)
	})
}

// TestEngine_UserLog_LogEventCheckFailure_KeepsDraining covers the fail-open fix in the
// LogEvent.Check branch: a settings read failure (a plain error, not ErrorBoundLimited)
// used to return false from processLogLine, which ended the emitUserLogs drain goroutine
// and silently discarded every remaining user log for the execution. Both lines must land.
func TestEngine_UserLog_LogEventCheckFailure_KeepsDraining(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	harness := newDropPathHarness(t, setupMockBillingClient(t), func(cfg *v2.EngineConfig) {
		cfg.LocalLimiters.LogEvent = &alwaysErrCheckLimiter[int]{ok: 1000, err: errors.New("settings unavailable")}
	})

	harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(_ context.Context, _ *sdkpb.ExecuteRequest, executor host.ExecutionHelper) {
			require.NoError(t, executor.EmitUserLog("first line"))
			require.NoError(t, executor.EmitUserLog("second line"))
		}).
		Return(nil, nil).Once()
	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_log_event_check_fails"},
	}

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		lines := userLogLines(t, harness.beholderObserver)
		assert.Contains(c, lines, "first line")
		assert.Contains(c, lines, "second line")
	}, 5*time.Second, 50*time.Millisecond)
}
