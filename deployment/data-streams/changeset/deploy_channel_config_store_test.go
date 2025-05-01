package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestDeployChannelConfigStore(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{}).Environment

	cc := DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{testutil.TestChain.Selector},
	}
	out, err := DeployChannelConfigStore{}.Apply(e, cc)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	for sel, addrMap := range ab {
		require.Equal(t, testutil.TestChain.Selector, sel)
		for _, tv := range addrMap {
			require.Equal(t, types.ChannelConfigStore, tv.Type)
		}
	}
}
