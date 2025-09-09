package v2_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
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

	testWorkflowNameA       = "my-best-workflow"
	hashedTestWorkflowNameA = "36363037306133663637"
	testWorkflowTagA        = "test-tag"
)

func TestEngineConfig_Validate(t *testing.T) {
	t.Parallel()
	cfg := defaultTestConfig(t, nil)

	t.Run("nil module", func(t *testing.T) {
		cfg.Module = nil
		require.Error(t, cfg.Validate())
	})

	t.Run("success", func(t *testing.T) {
		cfg.Module = modulemocks.NewModuleV2(t)
		require.NoError(t, cfg.Validate())
		require.NotEqual(t, 0, cfg.LocalLimits.HeartbeatFrequencyMs)
		require.NotEqual(t, 0, cfg.LocalLimits.ShutdownTimeoutMs)
		require.NotNil(t, cfg.Hooks.OnInitialized)
	})

	t.Run("empty workflow tag is allowed", func(t *testing.T) {
		cfg.Module = modulemocks.NewModuleV2(t)
		cfg.WorkflowTag = "" // V1 workflows don't have tags
		require.NoError(t, cfg.Validate())
	})
}

// defaultTestConfig returns a default v2.EngineConfig. CRE settings can optionally be configured by cfgFn.
func defaultTestConfig(t *testing.T, cfgFn func(*cresettings.Workflows)) *v2.EngineConfig {
	lf := limits.Factory{Logger: logger.TestLogger(t)}
	name, err := types.NewWorkflowName(testWorkflowNameA)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	sLimiter, err := syncerlimiter.NewWorkflowLimits(lggr, syncerlimiter.Config{}, lf)
	require.NoError(t, err)
	rateLimiter, err := ratelimiter.NewRateLimiter(ratelimiter.Config{
		GlobalRPS:      10.0,
		GlobalBurst:    100,
		PerSenderRPS:   10.0,
		PerSenderBurst: 100,
	}, lf)
	require.NoError(t, err)
	limiters, err := v2.NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, limiters.Close()) })

	return &v2.EngineConfig{
		Lggr:                              lggr,
		Module:                            modulemocks.NewModuleV2(t),
		CapRegistry:                       regmocks.NewCapabilitiesRegistry(t),
		DonTimeStore:                      dontime.NewStore(dontime.DefaultRequestTimeout),
		UseLocalTimeProvider:              true,
		ExecutionsStore:                   store.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		WorkflowID:                        testWorkflowID,
		WorkflowOwner:                     testWorkflowOwnerA,
		WorkflowName:                      name,
		WorkflowTag:                       testWorkflowTagA,
		WorkflowEncryptionKey:             workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		LocalLimits:                       v2.EngineLimits{},
		LocalLimiters:                     limiters,
		GlobalExecutionConcurrencyLimiter: sLimiter,
		GlobalExecutionRateLimiter:        rateLimiter,
		BeholderEmitter:                   &noopBeholderEmitter{},
		BillingClient:                     metmocks.NewBillingClient(t),
		WorkflowRegistryAddress:           "0x123",
		WorkflowRegistryChainSelector:     "11155111", // Sepolia chain ID
	}
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
