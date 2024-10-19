package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/llo"
)

func TestDeployChannelConfigStoreChangeSet(t *testing.T) {
	e := newMemoryEnv(t)
	c := llo.DeployLLOContractConfig{
		ChainsToDeploy: []uint64{TestChain.Selector},
	}
	out, err := DeployChannelConfigStoreChangeSet(e, c)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, tv.Type, llo.ChannelConfigStore)
		}
	}
}
