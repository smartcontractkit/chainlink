package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestCallSetChannelDefinitions(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnv(t, false, 0)

	// Deploy a contract
	deployConf := DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{testutil.TestChain.Selector},
	}
	out, err := DeployChannelConfigStore{}.Apply(e, deployConf)
	require.NoError(t, err)

	ab, err := out.AddressBook.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	var channelConfigStoreAddr common.Address
	for addr, tv := range ab[testutil.TestChain.Selector] {
		require.Equal(t, types.ChannelConfigStore, tv.Type)
		require.Equal(t, deployment.Version1_0_0, tv.Version)

		channelConfigStoreAddr = common.HexToAddress(addr)
	}

	// Update the environment with the newly deployed contract address.
	err = e.ExistingAddresses.Merge(out.AddressBook)
	require.NoError(t, err)

	// Call the contract.
	callConf := SetChannelDefinitionsConfig{
		DefinitionsByChain: map[uint64]map[string]ChannelDefinition{
			testutil.TestChain.Selector: {
				channelConfigStoreAddr.String(): {
					ChannelConfigStore: channelConfigStoreAddr,
					DonID:              1,
					S3URL:              "https://s3.us-west-2.amazonaws.com/data-streams-channel-definitions.stage.cldev.sh/channel-definitions-staging-mainnet-5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966.json",
					Hash:               hexToByte32("5ce78acee5113c55f795984cccdaeb7b805653a1c1e2f9d0d1e3279a302f7966"),
				},
			},
		},
		MCMSConfig: nil,
	}
	_, err = CallSetChannelDefinitions(e, callConf)
	require.NoError(t, err)
}

func hexToByte32(s string) [32]byte {
	var b [32]byte
	copy(b[:], common.HexToAddress(s).Bytes())
	return b
}
