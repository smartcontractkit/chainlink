package v2_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncaps "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	regmocks "github.com/smartcontractkit/chainlink-common/pkg/types/core/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"

	capmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils/matches"
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

	t.Run("inserts no-ops for unset life cycle hooks", func(t *testing.T) {
		cfg := defaultTestConfig(t, nil)
		cfg.Module = modulemocks.NewModuleV2(t)
		cfg.Hooks = v2.LifecycleHooks{}

		require.NoError(t, cfg.Validate())

		// Using reflection to verify all hooks set ensures none were missed by the test
		hooksValue := reflect.ValueOf(cfg.Hooks)
		hooksType := hooksValue.Type()
		for i := 0; i < hooksType.NumField(); i++ {
			field := hooksType.Field(i)
			fieldValue := hooksValue.Field(i)

			if !fieldValue.CanInterface() {
				continue
			}

			require.NotNil(t, fieldValue.Interface(), "hook field %s should not be nil after validation", field.Name)
		}
	})

	t.Run("does not override existing life cycle hooks", func(t *testing.T) {
		cfg := defaultTestConfig(t, nil)
		cfg.Module = modulemocks.NewModuleV2(t)
		called := make(map[string]bool)

		// Using reflection to verify all hooks set ensures none were missed by the test
		hooksValue := reflect.ValueOf(&cfg.Hooks).Elem()
		hooksType := hooksValue.Type()

		for i := 0; i < hooksType.NumField(); i++ {
			field := hooksType.Field(i)
			fieldValue := hooksValue.Field(i)

			if !fieldValue.CanSet() || !fieldValue.CanInterface() {
				continue
			}

			funcType := field.Type
			if funcType.Kind() != reflect.Func {
				continue
			}

			fieldName := field.Name
			callback := reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
				called[fieldName] = true
				var results []reflect.Value
				for j := 0; j < funcType.NumOut(); j++ {
					results = append(results, reflect.Zero(funcType.Out(j)))
				}
				return results
			})

			fieldValue.Set(callback)
		}

		require.NoError(t, cfg.Validate())

		// Use reflection to verify all hooks are still the same (not overridden)
		for i := 0; i < hooksType.NumField(); i++ {
			field := hooksType.Field(i)
			fieldValue := hooksValue.Field(i)

			if !fieldValue.CanInterface() {
				continue
			}

			fieldName := field.Name
			funcValue := fieldValue
			if funcValue.Kind() == reflect.Func {
				// Prepare zero-value arguments based on the function signature
				funcType := funcValue.Type()
				var args []reflect.Value
				for j := 0; j < funcType.NumIn(); j++ {
					args = append(args, reflect.Zero(funcType.In(j)))
				}
				funcValue.Call(args)

				assert.True(t, called[fieldName], "hook field %s should have been called", fieldName)
			}
		}
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
	limiters, err := v2.NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	subscriberMock := capmocks.NewDonSubscriber(t)
	subscriberMock.EXPECT().Subscribe(matches.AnyContext).Return(make(<-chan commoncaps.DON), func() {}, nil).Maybe()
	t.Cleanup(func() { assert.NoError(t, limiters.Close()) })

	return &v2.EngineConfig{
		Lggr:                              lggr,
		Module:                            modulemocks.NewModuleV2(t),
		CapRegistry:                       regmocks.NewCapabilitiesRegistry(t),
		DonTimeStore:                      dontime.NewStore(dontime.DefaultRequestTimeout),
		UseLocalTimeProvider:              true,
		DonSubscriber:                     subscriberMock,
		ExecutionsStore:                   store.NewInMemoryStore(lggr, clockwork.NewRealClock()),
		WorkflowID:                        testWorkflowID,
		WorkflowOwner:                     testWorkflowOwnerA,
		WorkflowName:                      name,
		WorkflowTag:                       testWorkflowTagA,
		WorkflowEncryptionKey:             workflowkey.MustNewXXXTestingOnly(big.NewInt(1)),
		LocalLimits:                       v2.EngineLimits{},
		LocalLimiters:                     limiters,
		GlobalExecutionConcurrencyLimiter: sLimiter,
		BeholderEmitter:                   &noopBeholderEmitter{},
		BillingClient:                     metmocks.NewBillingClient(t),
		WorkflowRegistryAddress:           "0x123",
		WorkflowRegistryChainSelector:     "11155111", // Sepolia chain ID
	}
}

type noopBeholderEmitter struct{}

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
