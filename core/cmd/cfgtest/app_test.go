package cfgtest

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmcfg2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	cfg2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	configtest2 "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestDefaultConfig(t *testing.T) {
	cfgtest.Clearenv(t)
	lggr := logger.TestLogger(t)
	chainID := utils.NewBigI(42991337)

	newGeneral, err := chainlink.GeneralConfigOpts{SkipEnv: true}.New(lggr)
	require.NoError(t, err)
	legacyGeneral := config.NewGeneralConfig(lggr)

	// we expect a mismatch on some methods with redefined defaults
	redefined := []string{"BlockEmissionIdleWarningThreshold", "DatabaseLockingMode", "EVMEnabled", "EVMRPCEnabled"}

	t.Run("general", func(t *testing.T) {
		assertMethodsReturnEqual[config.GeneralConfig](t, legacyGeneral, newGeneral, redefined...)
	})
	t.Run("general-test", func(t *testing.T) {
		assertMethodsReturnEqual[config.GeneralConfig](t, configtest.NewTestGeneralConfig(t), configtest2.NewTestGeneralConfig(t),
			"KeystorePassword", // new has a dummy value to pass validation
			"DatabaseURL",      // new has a dummy value to pass validation

			// Legacy package cltest defaults that were made standard.
			"JobPipelineReaperInterval",
			"P2PEnabled",
			"P2PNetworkingStack",
			"P2PNetworkingStackRaw",
			"ShutdownGracePeriod",

			// Legacy overrides root with random, but none of these others picked that up.
			// New uses test temp dir and inherits as expected.
			"AutoPprofProfileRoot",
			"CertFile",
			"KeyFile",
			"LogFileDir",
			"RootDir",
			"TLSDir",
			"AuditLoggerEnvironment", // same problem being derived from Dev())
		)
	})
	evmCfg := evmcfg2.EVMConfig{
		ChainID: chainID,
		Chain:   evmcfg2.Defaults(chainID),
	}
	newConfig := evmcfg2.NewTOMLChainScopedConfig(newGeneral, &evmCfg, lggr)
	legacyConfig := evmcfg.NewChainScopedConfig(chainID.ToInt(), evmtypes.ChainCfg{}, nil, lggr, legacyGeneral)

	t.Run("chain-scoped", func(t *testing.T) {
		assertMethodsReturnEqual[evmcfg.ChainScopedOnlyConfig](t, legacyConfig, newConfig, redefined...)
	})
}

// assertMethodsReturnEqual calls each method from M on a and b, and asserts the returned values are equal.
// Methods which accept arguments are skipped, allowe with a few other special cases.
func assertMethodsReturnEqual[M any](t *testing.T, a, b M, redefined ...string) {
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
			}
			if slices.Contains(redefined, name) {
				t.Skip("default redefined") // see core/cmd/chainlink/TestTOMLGeneralConfig_Defaults
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
