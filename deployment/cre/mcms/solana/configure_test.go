package solana_test

import (
	"testing"

	testenv "github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/stretchr/testify/require"

	cresolmcms "github.com/smartcontractkit/chainlink/deployment/cre/mcms/solana"
)

func TestConfigureSolanaMCMS_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	loader := testenv.NewLoader()
	env, err := loader.Load(t.Context(), testenv.WithEVMSimulatedN(t, 1))
	require.NoError(t, err)

	cs := cresolmcms.ConfigureSolanaMCMS{}
	stagingCfg := testStagingMCMSConfig()

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{})
	require.ErrorContains(t, err, "no chain selectors provided")

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{
		ChainSelectors:         []uint64{5009297550715157269},
		MCMSWithTimelockConfig: stagingCfg,
	})
	require.ErrorContains(t, err, "not a solana chain")

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{
		ChainSelectors:         []uint64{16423721717087811551},
		MCMSWithTimelockConfig: stagingCfg,
	})
	require.ErrorContains(t, err, "solana chain not found")
}
