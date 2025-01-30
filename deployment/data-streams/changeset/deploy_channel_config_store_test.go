package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams"
)

func TestDeployChannelConfigStore(t *testing.T) {
	env := newMemoryEnv(t)
	cc := DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{TestChain.Selector},
	}
	out, err := DeployChannelConfigStore(env, cc)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, data_streams.ChannelConfigStore, tv.Type)
		}
	}
}
