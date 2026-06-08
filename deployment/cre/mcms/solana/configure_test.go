package solana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	testenv "github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	cresolmcms "github.com/smartcontractkit/chainlink/deployment/cre/mcms/solana"
)

func TestConfigureSolanaMCMS_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	loader := testenv.NewLoader()
	env, err := loader.Load(t.Context(), testenv.WithEVMSimulatedN(t, 1))
	assert.NoError(t, err)

	cs := cresolmcms.ConfigureSolanaMCMS{}
	stagingCfg := testStagingMCMSConfig()

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{})
	assert.ErrorContains(t, err, "no chain selectors provided")

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{
		ChainSelectors:         []uint64{5009297550715157269},
		MCMSWithTimelockConfig: stagingCfg,
	})
	assert.ErrorContains(t, err, "not a solana chain")

	err = cs.VerifyPreconditions(*env, cresolmcms.ConfigureSolanaMCMSConfig{
		ChainSelectors:         []uint64{16423721717087811551},
		MCMSWithTimelockConfig: stagingCfg,
	})
	assert.ErrorContains(t, err, "solana chain not found")
}
