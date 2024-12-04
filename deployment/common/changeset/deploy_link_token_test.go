package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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
	state, err := changeset.LoadLinkTokenState(e.Chains[chain1], addrs)
	require.NoError(t, err)
	view, err := state.GenerateLinkView()
	require.NoError(t, err)
	assert.Equal(t, view.Owner, e.Chains[chain1].DeployerKey.From)
	assert.Equal(t, view.TypeAndVersion, "LinkToken 1.0.0")
	// Initially nothing minted.
	assert.Equal(t, view.Supply.String(), "0")
}
