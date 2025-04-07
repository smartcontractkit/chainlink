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

func TestCallPromoteStagingConfig(t *testing.T) {
	e := testutil.NewMemoryEnv(t, true, 0)
	ctx := testcontext.Get(t)

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
	require.NoError(t, err)

	configuratorAddr := common.HexToAddress(configuratorAddrHex)

	// Prepare a production onchain config (64 bytes):
	// - First 32 bytes: version (with 1 at the last byte)
	// - Next 32 bytes: all zeros (predecessor must be zero for production)
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

	configAbi, err := abi.JSON(strings.NewReader(configurator.ConfiguratorMetaData.ABI))
	require.NoError(t, err)

	prodFilterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   nil,
		Addresses: []common.Address{configuratorAddr},
		Topics:    [][]common.Hash{{configAbi.Events["ProductionConfigSet"].ID}},
	}
	prodLogs, err := e.Chains[testutil.TestChain.Selector].Client.FilterLogs(ctx, prodFilterQuery)
	require.NoError(t, err)
	require.NotEmpty(t, prodLogs)

	var productionDigest [32]byte
	// ABI.Unpack of log.Data returns a slice where index 1 is the configDigest.
	for _, lg := range prodLogs {
		decoded, err := configAbi.Unpack("ProductionConfigSet", lg.Data)
		require.NoError(t, err)
		digestBytes := decoded[1].([32]byte)
		copy(productionDigest[:], digestBytes[:])
		break
	}

	// 3. Build the staging onchain config:
	// First 32 bytes: version (1)
	// Next 32 bytes: productionDigest (as predecessor)
	stgConfig := make([]byte, 64)
	stgConfig[31] = 1 // version = 1
	copy(stgConfig[32:64], productionDigest[:])

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

	callPromote := PromoteStagingConfigConfig{
		PromotionsByChain: map[uint64][]PromoteStagingConfig{
			testutil.TestChain.Selector: {
				{
					ConfiguratorAddress: configuratorAddr,
					ConfigID:            [32]byte{},
					// currentState = false (matches production state)
					IsGreenProduction: false,
				},
			},
		},
		MCMSConfig: nil,
	}

	e, err = commonChangesets.Apply(t, e, nil,
		commonChangesets.Configure(
			PromoteStagingConfigChangeset,
			callPromote,
		),
	)

	require.NoError(t, err)

	promoFilterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   nil,
		Addresses: []common.Address{configuratorAddr},
		Topics:    [][]common.Hash{{configAbi.Events["PromoteStagingConfig"].ID}},
	}
	promoLogs, err := e.Chains[testutil.TestChain.Selector].Client.FilterLogs(ctx, promoFilterQuery)
	require.NoError(t, err)
	require.NotEmpty(t, promoLogs)

	// Process the first promotion event log.
	promoLog := promoLogs[0]

	decodedPromo, err := configAbi.Unpack("PromoteStagingConfig", promoLog.Data)
	require.NoError(t, err)
	require.Len(t, decodedPromo, 1)
	isGreenProduction := decodedPromo[0].(bool)

	// The indexed parameters are in Topics:
	// Topic[0] is the event signature,
	// Topic[1] is configId, and Topic[2] is retiredConfigDigest.
	require.GreaterOrEqual(t, len(promoLog.Topics), 3)
	var eventConfigID [32]byte
	copy(eventConfigID[:], promoLog.Topics[1].Bytes())
	var eventRetiredDigest [32]byte
	copy(eventRetiredDigest[:], promoLog.Topics[2].Bytes())

	var expectedConfigID [32]byte
	require.Equal(t, expectedConfigID, eventConfigID)
	// retiredConfigDigest should equal the production config digest (as production becomes retired).
	require.Equal(t, productionDigest, eventRetiredDigest)
	// isGreenProduction flag is flipped to true.
	require.True(t, isGreenProduction)
}
