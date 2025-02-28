package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestDeployChannelConfigStore(t *testing.T) {
	e := newMemoryEnv(t)
	cc := DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{TestChain.Selector},
	}
	out, err := DeployChannelConfigStore{}.Apply(e, cc)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, types.ChannelConfigStore, tv.Type)
		}
	}
}
