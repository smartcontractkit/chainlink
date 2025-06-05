package v2_test

import (
	"context"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

const (
	testWorkflowID = "ffffaabbccddeeff00112233aabbccddeeff00112233aabbccddeeff00112233"

	testWorkflowOwnerA = "1100000000000000000000000000000000000000"
	testWorkflowOwnerB = "2200000000000000000000000000000000000000"
	testWorkflowOwnerC = "3300000000000000000000000000000000000000"

	testWorkflowNameA = "my-best-workflow"
)

func TestEngineConfig_Validate(t *testing.T) {
	t.Parallel()
	cfg := defaultTestConfig(t)

	t.Run("nil module", func(t *testing.T) {
		cfg.Module = nil
		require.Error(t, cfg.Validate())
	})

	t.Run("success", func(t *testing.T) {
		cfg.Module = modulemocks.NewModuleV2(t)
		require.NoError(t, cfg.Validate())
		require.NotEqual(t, 0, cfg.LocalLimits.TriggerEventQueueSize)
		require.NotNil(t, cfg.Hooks.OnInitialized)
	})
}

func defaultTestConfig(t *testing.T) *v2.EngineConfig {
	name, err := types.NewWorkflowName(testWorkflowNameA)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	sLimiter, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{})
	require.NoError(t, err)
	rateLimiter, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      10.0,
		GlobalBurst:    100,
		PerSenderRPS:   10.0,
		PerSenderBurst: 100,
	})
	require.NoError(t, err)

	return &v2.EngineConfig{
		Lggr:                 lggr,
		Module:               modulemocks.NewModuleV2(t),
		CapRegistry:          regmocks.NewCapabilitiesRegistry(t),
		ExecutionsStore:      store.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		WorkflowID:           testWorkflowID,
		WorkflowOwner:        testWorkflowOwnerA,
		WorkflowName:         name,
		LocalLimits:          v2.EngineLimits{},
		GlobalLimits:         sLimiter,
		ExecutionRateLimiter: rateLimiter,
		BeholderEmitter:      &noopBeholderEmitter{},
		BillingClient:        &mockBillingClient{},
	}
}

type mockBillingClient struct {
	mock.Mock
}

func (_m *mockBillingClient) SubmitWorkflowReceipt(ctx context.Context, req *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error) {
	args := _m.Called(ctx, req)

	var a0 *billing.SubmitWorkflowReceiptResponse
	if arg, ok := args.Get(0).(*billing.SubmitWorkflowReceiptResponse); ok {
		a0 = arg
	}

	return a0, args.Error(1)
}

type noopBeholderEmitter struct {
}

func (m *noopBeholderEmitter) Emit(_ context.Context, _ string) error {
	return nil
}

func (m *noopBeholderEmitter) WithMapLabels(labels map[string]string) custmsg.MessageEmitter {
	return m
}

func (m *noopBeholderEmitter) With(kvs ...string) custmsg.MessageEmitter {
	return m
}

func (m *noopBeholderEmitter) Labels() map[string]string {
	return map[string]string{}
}
