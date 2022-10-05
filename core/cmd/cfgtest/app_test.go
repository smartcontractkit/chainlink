package cfgtest

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmcfg2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	cfg2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestDefaultConfig(t *testing.T) {
	cfgtest.Clearenv(t)
	lggr := logger.TestLogger(t)
	chainID := utils.NewBigI(42991337)

	newGeneral, err := chainlink.NewTOMLGeneralConfig(lggr, "", "", nil, nil)
	require.NoError(t, err)
	legacyGeneral := config.NewGeneralConfig(lggr)

	t.Run("general", func(t *testing.T) {
		assertMethodsReturnEqual[config.GeneralConfig](t, legacyGeneral, newGeneral)
	})
	newChain, _ := evmcfg2.Defaults(chainID)
	evmCfg := evmcfg2.EVMConfig{
		ChainID: chainID,
		Chain:   newChain,
	}
	newConfig := evmcfg2.NewTOMLChainScopedConfig(newGeneral, &evmCfg, lggr)
	legacyConfig := evmcfg.NewChainScopedConfig(chainID.ToInt(), evmtypes.ChainCfg{}, nil, lggr, legacyGeneral)

	t.Run("chain-scoped", func(t *testing.T) {
		assertMethodsReturnEqual[evmcfg.ChainScopedOnlyConfig](t, legacyConfig, newConfig)
	})
}

// assertMethodsReturnEqual calls each method from M on a and b, and asserts the returned values are equal.
// Methods which accept arguments are skipped, allowe with a few other special cases.
func assertMethodsReturnEqual[M any](t *testing.T, a, b M) {
	av, bv := reflect.ValueOf(a), reflect.ValueOf(b)
	to := reflect.TypeOf((*M)(nil)).Elem()
	for i := 0; i < to.NumMethod(); i++ {
		m := to.Method(i)
		name := m.Name
		t.Run(name, func(t *testing.T) {
			if m.Type.NumIn() > 0 {
				t.Skip("has arguments")
			}
			switch name {
			case "Validate", "PersistedConfig", "SetEvmGasPriceDefault":
				t.Skip("irrelevant")
			case "P2PListenPort", "AppID":
				t.Skip("randomized")
			case "EVMEnabled", "P2PNetworkingStack", "P2PNetworkingStackRaw", "P2PEnabled", "BlockEmissionIdleWarningThreshold":
				t.Skip("default redefined")
			}

			defer func() {
				r := recover()
				if r == nil {
					return
				}
				err, ok := r.(error)
				if !ok {
					t.Fatalf("panic: %v", r)
				}
				if errors.Is(err, cfg2.ErrUnsupported) {
					return // no problem - expected mismatch
				}
				t.Fatalf("panic: %v", err)
			}()
			ar := av.MethodByName(m.Name).Call(nil)
			br := bv.MethodByName(m.Name).Call(nil)
			for i := range ar {
				ae, be := ar[i], br[i]
				t.Logf("a: %v (%s)", ae, ae.Type())
				t.Logf("b: %v (%s)", be, be.Type())
				var ai, bi any
				if ae.Kind() == reflect.Ptr {
					if ae.IsNil() {
						assert.True(t, be.IsNil(), "a is nil but b is not")
						continue
					}
					if !assert.False(t, be.IsNil(), "b is nil but a is not") {
						continue
					}
					ai, bi = ae.Elem().Interface(), be.Elem().Interface()
				} else {
					ai, bi = ae.Interface(), be.Interface()
				}
				if aerr, ok := ai.(error); ok {
					berr, ok := bi.(error)
					if assert.True(t, ok, "a is error but b is not") {
						aroot, broot := unwrapFully(aerr), unwrapFully(berr)
						assert.Equal(t, aroot, broot, "different root errors")
					}
					continue
				}
				assert.Equal(t, ai, bi, "%dth return arg", i)

			}
		})
	}
}

func unwrapFully(e error) error {
	uw := errors.Unwrap
	for u := uw(e); u != nil; u = uw(e) {
		e = u
	}
	return e
}
