package cre

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func BuildMinimalEnvironment(t *testing.T, lggr logger.Logger) (cldf.Environment, uint64) {
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  0,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	t.Logf("Environment created with operations bundle")

	// Get the chain selector for the deployment
	chainSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.NotEmpty(t, chainSelectors, "should have at least one EVM chain")
	chainSelector := chainSelectors[0]
	t.Logf("Using chain selector: %d", chainSelector)

	return env, chainSelector
}
