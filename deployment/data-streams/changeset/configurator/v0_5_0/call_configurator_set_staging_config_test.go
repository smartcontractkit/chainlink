package v0_5_0

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

func TestCallSetStagingConfig(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)

	e, err := commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			DeployConfiguratorChangeset,
			DeployConfiguratorConfig{
				ChainsToDeploy: []uint64{testutil.TestChain.Selector},
			},
		),
	)

	require.NoError(t, err)

	ab, err := e.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Len(t, ab, 1)

	configuratorAddrHex, err := deployment.SearchAddressBook(e.ExistingAddresses, testutil.TestChain.Selector, types.Configurator)
	configuratorAddr := common.HexToAddress(configuratorAddrHex)

	require.NoError(t, err)

	onchainConfigHex := "0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	onchainConfigProd, err := hex.DecodeString(onchainConfigHex)
	require.NoError(t, err)
	require.Len(t, onchainConfigProd, 64)

	prodCfg := SetProductionConfig{
		ConfiguratorAddress:   configuratorAddr,
		ConfigID:              [32]byte{},
		Signers:               [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
		OffchainTransmitters:  [][32]byte{{}, {}, {}, {}},
		F:                     1,
		OnchainConfig:         onchainConfigProd,
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain config prod"),
	}

	callProd := SetProductionConfigConfig{
		ConfigurationsByChain: map[uint64][]SetProductionConfig{
			testutil.TestChain.Selector: {prodCfg},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetProductionConfigChangeset,
			callProd,
		),
	)

	require.NoError(t, err)

	ctx := testcontext.Get(t)

	configAbi, err := abi.JSON(strings.NewReader(configurator.ConfiguratorMetaData.ABI))
	require.NoError(t, err)

	filterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   nil,
		Addresses: []common.Address{configuratorAddr},
		Topics:    [][]common.Hash{{configAbi.Events["ProductionConfigSet"].ID}},
	}
	prodLogs, err := e.Chains[testutil.TestChain.Selector].Client.FilterLogs(ctx, filterQuery)
	require.NoError(t, err)
	require.NotEmpty(t, prodLogs)

	var productionDigest [32]byte
	for _, lg := range prodLogs {
		decoded, err := configAbi.Unpack("ProductionConfigSet", lg.Data)
		require.NoError(t, err)

		digestBytes := decoded[1].([32]byte)
		copy(productionDigest[:], digestBytes[:])
	}

	stgConfig := make([]byte, 64)
	stgConfig[31] = 1                           // version
	copy(stgConfig[32:64], productionDigest[:]) // predecessor = productionDigest

	stagingCfg := SetStagingConfig{
		ConfiguratorAddress:   configuratorAddr,
		ConfigID:              [32]byte{},
		Signers:               [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
		OffchainTransmitters:  [][32]byte{{}, {}, {}, {}},
		F:                     1,
		OnchainConfig:         stgConfig,
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain config staging"),
	}

	callStaging := SetStagingConfigConfig{
		ConfigurationsByChain: map[uint64][]SetStagingConfig{
			testutil.TestChain.Selector: {stagingCfg},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			SetStagingConfigChangeset,
			callStaging,
		),
	)

	require.NoError(t, err)
}
