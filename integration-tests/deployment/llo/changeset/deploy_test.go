package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/llo"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployChannelConfigStoreChangeSet(t *testing.T) {
	lggr := logger.TestLogger(t)
	tenv := llo.NewMemoryEnvironment(t, lggr)
	e := tenv.Env

	c := llo.DeployLLOContractConfig{
		ChainsToDeploy: []uint64{llo.TestChain.Selector},
	}
	out, err := DeployChannelConfigStoreChangeSet(e, c)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, llo.TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, tv.Type, llo.ChannelConfigStore)
		}
	}
}
