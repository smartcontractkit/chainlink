package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TODO: This is to test the implementation of Aptos chains in memory environment
// To be deleted after changesets tests are added
func TestAptosMemoryEnv(t *testing.T) {
	lggr := logger.TestLogger(t)
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		AptosChains: 1,
	})
	aptosChainSelectors := env.AllChainSelectorsAptos()
	require.Len(t, aptosChainSelectors, 1)
	require.NotEqual(t, 0, env.AptosChains[0].Selector)
}

// TODO: This is to test the implementation of Aptos chains in memory environment
// To be deleted after changesets tests are added
func TestAptosHelperMemoryEnv(t *testing.T) {
	depEvn, testEnv := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
		testhelpers.WithNoJobsAndContracts(), // currently not supporting jobs and contracts
	)
	aptosChainSelectors := depEvn.Env.AllChainSelectorsAptos()
	require.Len(t, aptosChainSelectors, 1)
	aptosChainSelectors2 := testEnv.DeployedEnvironment().Env.AllChainSelectorsAptos()
	require.Len(t, aptosChainSelectors2, 1)
}
