package v2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// newTestMeter returns a private meter backed by a manual meter reader, so
// tests can inject it via EngineConfig.Meter and then collect and inspect the
// engine's counters. It never touches the process-global beholder client, so
// it is safe in parallel tests. Mirrors newTestMeter in the internal test
// package; every increment on the returned reader comes from this test's engine.
func newTestMeter(t *testing.T) (metric.Meter, *sdkmetric.ManualReader) {
	t.Helper()

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	return mp.Meter("beholder"), reader
}

// collectCounterByAttr sums the named counter and returns the values observed
// for the given attribute key.
func collectCounterByAttr(t *testing.T, reader *sdkmetric.ManualReader, name, attrKey string) (int64, []string) {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	var total int64
	var values []string
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			data, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, m.Data)
			for _, dp := range data.DataPoints {
				total += dp.Value
				if v, ok := dp.Attributes.Value(attribute.Key(attrKey)); ok {
					values = append(values, v.AsString())
				}
			}
		}
	}
	return total, values
}

// collectInitFailureCounter sums the named counter and returns the values
// observed for the "reason" attribute.
func collectInitFailureCounter(t *testing.T, reader *sdkmetric.ManualReader, name string) (int64, []string) {
	t.Helper()
	return collectCounterByAttr(t, reader, name, "reason")
}

// TestEngine_InitFailureCounter_DisallowedSecretsCall verifies end to end
// that a module failing the Subscribe phase by attempting a secrets call
// during trigger subscription fails engine initialization and increments the
// initialization failure counter with the disallowed-secrets reason.
func TestEngine_InitFailureCounter_DisallowedSecretsCall(t *testing.T) {
	t.Parallel()

	meter, reader := newTestMeter(t)

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error, 1)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Meter = meter
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	module.EXPECT().Start().Once()
	module.EXPECT().Close().Once()
	// the module reports the host's secrets rejection back as its error result
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{
			Result: &sdkpb.ExecutionResult_Error{
				Error: "failed to get secrets for call 1: secrets calls cannot be made during trigger subscription",
			},
		}, nil).Once()

	servicetest.Run(t, engine)

	initErr := <-initDoneCh
	require.ErrorContains(t, initErr, "secrets calls cannot be made during trigger subscription")

	// the meter is private to this test, so the count is exact
	gotCount, gotReasons := collectInitFailureCounter(t, reader, "platform_engine_workflow_initialization_failures_total")
	require.Equal(t, int64(1), gotCount)
	require.Equal(t, []string{"disallowed_secrets_call_during_subscription"}, gotReasons)
}

// TestEngine_InitFailureCounter_ZeroTriggerSubscriptions verifies end to end
// that a module reporting zero trigger subscriptions fails engine
// initialization as a permanent error and increments the initialization
// failure counter with the zero-trigger-subscriptions reason, so this
// condition can be alerted on.
func TestEngine_InitFailureCounter_ZeroTriggerSubscriptions(t *testing.T) {
	t.Parallel()

	meter, reader := newTestMeter(t)

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error, 1)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
	cfg.Meter = meter
	cfg.Hooks = v2.LifecycleHooks{
		OnInitialized: func(err error) {
			initDoneCh <- err
		},
	}

	engine, err := v2.NewEngine(cfg)
	require.NoError(t, err)

	module.EXPECT().Start().Once()
	module.EXPECT().Close().Once()
	module.EXPECT().Execute(matches.AnyContext, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{
			Result: &sdkpb.ExecutionResult_TriggerSubscriptions{
				TriggerSubscriptions: &sdkpb.TriggerSubscriptionRequest{},
			},
		}, nil).Once()

	servicetest.Run(t, engine)

	initErr := <-initDoneCh
	require.ErrorIs(t, initErr, v2.ErrZeroTriggerSubscriptions)

	// the meter is private to this test, so the count is exact
	gotCount, gotReasons := collectInitFailureCounter(t, reader, "platform_engine_workflow_initialization_failures_total")
	require.Equal(t, int64(1), gotCount)
	require.Equal(t, []string{"zero_trigger_subscriptions"}, gotReasons)
}
