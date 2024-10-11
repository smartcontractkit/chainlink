package llo

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type DeployedEnv struct {
	Env          deployment.Environment
	Ab           deployment.AddressBook
	HomeChainSel uint64
	FeedChainSel uint64
	ReplayBlocks map[uint64]uint64
}

var (
	// TestChain is the chain used by the in-memory environment.
	TestChain = chain_selectors.Chain{
		EvmChainID: 90000001,
		Selector:   909606746561742123,
		Name:       "Test Chain",
		VarName:    "",
	}
)

// NewMemoryEnvironment creates a new LLO environment with capreg, fee tokens, feeds and nodes set up.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger) DeployedEnv {
	chains := memory.NewMemoryChains(t, 1)
	rc := deployment.CapabilityRegistryConfig{}
	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, rc)
	ctx := testcontext.Get(t)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedEnv{
		Ab:  deployment.NewMemoryAddressBook(),
		Env: e,
	}
}
