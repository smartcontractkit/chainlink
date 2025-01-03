package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.AllChainSelectors()[0]
	e, err := changeset.ApplyChangesets(t, e, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(changeset.DeployLinkToken),
			Config:    []uint64{chain1},
		},
	})
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)
	state, err := changeset.MaybeLoadLinkTokenChainState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	// View itself already unit tested
	_, err = state.GenerateLinkView()
	require.NoError(t, err)
}

// TODO: fails sometimes with
// (string) (len=96) "Program TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb failed: invalid account data for instruction"
func TestDeployLinkTokenSol(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := deployment.SolanaChainSelector
	e, err := changeset.ApplyChangesets(t, e, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(changeset.DeployLinkTokenSolana),
			Config:    []uint64{chain1},
		},
	})
	require.NoError(t, err)
	addrs, err := e.ExistingAddresses.AddressesForChain(chain1)
	require.NoError(t, err)
	t.Log(addrs)
	// // View itself already unit tested
	// _, err = state.GenerateLinkView()
	// require.NoError(t, err)
}
