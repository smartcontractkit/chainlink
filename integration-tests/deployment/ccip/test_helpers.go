package ccipdeployment

import (
	"context"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

type DeployedTestEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	Nodes        map[string]memory.Node
}

// NewDeployedEnvironment creates a new CCIP environment
// with capreg and nodes set up.
func NewDeployedTestEnvironment(t *testing.T, lggr logger.Logger) DeployedTestEnvironment {
	ctx := Context(t)
	chains := memory.NewMemoryChains(t, 3)
	homeChainSel := uint64(0)
	homeChainEVM := uint64(0)
	// Say first chain is home chain.
	for chainSel := range chains {
		homeChainEVM, _ = chainsel.ChainIdFromSelector(chainSel)
		homeChainSel = chainSel
		break
	}
	ab, capReg, err := DeployCapReg(lggr, chains, homeChainSel)
	require.NoError(t, err)

	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, memory.RegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	})
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedTestEnvironment{
		Ab:           ab,
		Env:          e,
		HomeChainSel: homeChainSel,
		Nodes:        nodes,
	}
}
