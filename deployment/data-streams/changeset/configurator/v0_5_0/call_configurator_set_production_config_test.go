package v0_5_0

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

func TestCallSetProductionConfig(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	chainSelector := e.AllChainSelectors()[0]

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployConfiguratorChangeset,
			DeployConfiguratorConfig{
				ChainsToDeploy: []uint64{chainSelector},
			},
		),
	)

	require.NoError(t, err)

	configuratorAddrHex, err := dsutil.GetContractAddress(e.DataStore.Addresses(), types.Configurator)
	require.NoError(t, err)

	configuratorAddr := common.HexToAddress(configuratorAddrHex)

	onchainConfigHex := "0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	onchainConfig, err := hex.DecodeString(onchainConfigHex)
	require.NoError(t, err)
	require.Len(t, onchainConfig, 64)

	prodCfg := ConfiguratorConfig{
		ConfiguratorAddress: configuratorAddr,
		ConfigID:            [32]byte{},
		Signers: [][]byte{
			{0x01}, {0x02}, {0x03}, {0x04},
		},
		OffchainTransmitters: [][32]byte{
			{}, {}, {}, {},
		},
		F:                     1,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain config"),
	}

	callConf := SetProductionConfigConfig{
		ConfigurationsByChain: map[uint64][]ConfiguratorConfig{
			chainSelector: {prodCfg},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetProductionConfigChangeset,
			callConf,
		),
	)

	require.NoError(t, err)

	t.Run("VerifyMetadata", func(t *testing.T) {
		// Use View To Confirm Data
		_, outputs, err := commonChangesets.ApplyChangesetsV2(t, e,
			[]commonChangesets.ConfiguredChangeSet{
				commonChangesets.Configure(
					changeset.SaveContractViews,
					changeset.SaveContractViewsConfig{
						Chains: []uint64{testutil.TestChain.Selector},
					},
				),
			},
		)
		require.NoError(t, err)
		require.Len(t, outputs, 1)
		output := outputs[0]

		client := e.Chains[testutil.TestChain.Selector].Client
		contract, err := configurator.NewConfigurator(configuratorAddr, client)
		require.NoError(t, err)

		stagingIter, err := contract.FilterProductionConfigSet(nil, nil)
		require.NoError(t, err)
		defer stagingIter.Close()

		var digest [32]byte
		for stagingIter.Next() {
			event := stagingIter.Event
			digest = event.ConfigDigest
		}

		VerifyConfiguratorState(t,
			output.DataStore,
			testutil.TestChain.Selector,
			configuratorAddr,
			digest,
			prodCfg,
			1)
	})
}
