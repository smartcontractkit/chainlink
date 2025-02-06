package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func TestCallSetChannelDefinitions(t *testing.T) {
	e := newMemoryEnv(t)
	someChannelConfigStore := common.HexToAddress("0xc2C56B48C2225Cbb18758a0AE1968E66285EE90f")
	cc := SetChannelDefinitionsConfig{
		DefinitionsByChain: map[uint64]map[string]ChannelDefinition{
			TestChain.Selector: {
				someChannelConfigStore.String(): {
					ChannelConfigStore: someChannelConfigStore,
					DonID:              1,
					S3URL:              "https://s3.us-west-2.amazonaws.com/data-streams-channel-definitions.stage.cldev.sh/channel-definitions-staging-mainnet-5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966.json",
					Hash:               hexToByte32("5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966"),
				},
			},
		},
		MCMSConfig: nil,
	}
	out, err := CallSetChannelDefinitions(e, cc)
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

func hexToByte32(s string) [32]byte {
	var b [32]byte
	copy(b[:], common.HexToAddress(s).Bytes())
	return b
}
