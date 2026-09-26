package v2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
)

// setupInitFailureTestMeter swaps the global beholder client for one backed
// by a manual meter reader, so the engine's counters can be collected. Mirrors
// setupTestMeter in the internal test package.
//
// The swap is global: parallel full-engine tests may register their counters
// on this meter too, so assertions here must tolerate foreign increments
// (reason-unique labels and lower bounds, not exact totals).
func setupInitFailureTestMeter(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	prevClient := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prevClient) })

	client := beholder.NoopClientConfig{Lggr: logger.Test(t)}.New()
	client.Meter = mp.Meter("beholder")
	client.MeterProvider = mp
	beholder.SetClient(client)

	return reader
}

// collectInitFailureCounter sums the named counter and returns the values
// observed for the "reason" attribute.
func collectInitFailureCounter(t *testing.T, reader *sdkmetric.ManualReader, name string) (int64, []string) {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	var total int64
	var reasons []string
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			data, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, m.Data)
			for _, dp := range data.DataPoints {
				total += dp.Value
				if r, ok := dp.Attributes.Value("reason"); ok {
					reasons = append(reasons, r.AsString())
				}
			}
		}
	}
	return total, reasons
}

// TestEngine_InitFailureCounter_DisallowedSecretsCall verifies end to end
// that a module failing the Subscribe phase by attempting a secrets call
// during trigger subscription fails engine initialization and increments the
// initialization failure counter with the disallowed-secrets reason.
//
//nolint:paralleltest // swaps the global beholder client
func TestEngine_InitFailureCounter_DisallowedSecretsCall(t *testing.T) {
	reader := setupInitFailureTestMeter(t)

	module := modulemocks.NewModuleV2(t)
	capreg := regmocks.NewCapabilitiesRegistry(t)
	capreg.EXPECT().LocalNode(matches.AnyContext).Return(newNode(t), nil)

	initDoneCh := make(chan error, 1)

	cfg := defaultTestConfig(t, nil)
	cfg.Module = module
	cfg.CapRegistry = capreg
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

	// count may include foreign increments from parallel engine tests sharing
	// the swapped meter; the disallowed-secrets reason is unique to this test
	gotCount, gotReasons := collectInitFailureCounter(t, reader, "platform_engine_workflow_initialization_failures_total")
	require.GreaterOrEqual(t, gotCount, int64(1))
	require.Contains(t, gotReasons, "disallowed_secrets_call_during_subscription")
}
