package v2_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// TestEngine_LimitersAlwaysSeeTenant guards against CRE context propagation regressions.
// A limiter call without the tenant its scope needs either fails closed (required scopes: the
// execution is dropped) or silently fails open (org scope), and neither shows up in a normal
// unit test. Here one execution goes through every limiter site, with the limiters logging to
// an observed logger, and the test fails if any of them reports a missing tenant.
//
// Known gap: the scoped resource pool limiter (the org layer of ExecutionConcurrency) does not
// use the factory logger yet, so a missing org there is not detected until chainlink-common
// wires it.
//
// When adding a limiter call to the engine, make sure this execution reaches it.
func TestEngine_LimitersAlwaysSeeTenant(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	limitsLggr, limitsLogs := logger.TestObserved(t, zapcore.WarnLevel)
	engineLggr, engineLogs := logger.TestObserved(t, zapcore.WarnLevel)

	// Registered before the harness, so it runs after the engine has closed and also covers
	// the per-tenant cleanup done on shutdown.
	t.Cleanup(func() {
		for _, entry := range limitsLogs.FilterMessageSnippet("missing tenant").All() {
			t.Errorf("limiter reported a missing tenant: %s %v", entry.Message, entry.ContextMap())
		}
		for _, snippet := range []string{"with no usable value", "could not be evaluated"} {
			for _, entry := range engineLogs.FilterMessageSnippet(snippet).All() {
				t.Errorf("engine hit a limit it could not evaluate: %s %v", entry.Message, entry.ContextMap())
			}
		}
	})

	harness := newDropPathHarness(t, setupMockBillingClient(t), func(cfg *v2.EngineConfig) {
		lf := limits.Factory{Logger: limitsLggr}

		limiters, err := v2.NewLimiters(lf, nil)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, limiters.Close()) })

		featureFlags, err := v2.NewFeatureFlags(lf, nil)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, featureFlags.Close()) })

		workflowLimit, err := syncerlimiter.NewWorkflowLimits(engineLggr, syncerlimiter.Config{}, lf)
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, workflowLimit.Close()) })

		cfg.LocalLimiters = limiters
		cfg.FeatureFlags = featureFlags
		cfg.GlobalWorkflowLimit = workflowLimit
		cfg.Lggr = engineLggr
	})

	// The capability and secrets limits are enforced before the capability is resolved, so
	// failing the lookup is enough to exercise them without a real capability.
	harness.capreg.EXPECT().GetExecutable(matches.AnyContext, mock.Anything).
		Return(nil, errors.New("capability not available in this test")).Maybe()

	const userLog = "tenant propagation check"
	harness.module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, _ *sdkpb.ExecuteRequest, helper host.ExecutionHelper) {
			// An allowlisted chain selector, so the call gets past ChainAllowed to the per-call
			// limit and the capability concurrency pool.
			// Both calls are expected to fail on the capability lookup. Required scopes report a
			// missing tenant through the returned error rather than a log, so check for it here.
			// assert, not require: this runs on the engine's goroutine, where FailNow would kill
			// that goroutine instead of failing the test.
			_, err := helper.CallCapability(ctx, &sdkpb.CapabilityRequest{
				Id:     "evm:ChainSelector:3379446385462418246@1.0.0",
				Method: "Write",
			})
			assert.NotErrorIs(t, err, limits.ErrMissingTenant{}, "capability call limits") //nolint:testifylint // see above
			_, err = helper.GetSecrets(ctx, &sdkpb.GetSecretsRequest{
				Requests: []*sdkpb.SecretRequest{{Id: "secret", Namespace: "main"}},
			})
			assert.NotErrorIs(t, err, limits.ErrMissingTenant{}, "secrets call limit") //nolint:testifylint // see above
			assert.NoError(t, helper.EmitUserLog(userLog))
		}).
		Return(nil, nil).Once()

	harness.eventCh <- capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_tenant_propagation"},
	}
	select {
	case <-harness.executionFinishedCh:
	case <-time.After(10 * time.Second):
		t.Fatal("execution never finished; a limit read with no usable value may have dropped it")
	}

	evt, ok := latestV2FinishedEvent(t, harness.beholderObserver)
	require.True(t, ok)
	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED, evt.Status)

	// The log line is drained asynchronously; wait for it so the LogEvent/LogLine checks have run.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.NotEmpty(c, harness.beholderObserver.Messages(t, "beholder_entity", v1UserLogsEntity))
	}, 5*time.Second, 50*time.Millisecond)
}
