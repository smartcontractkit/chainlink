package aptos_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	aptosfeequoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"

	// "github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	seq "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/sequence"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/operations"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInputSerializable(t *testing.T) {
	input := seq.UpdateAptosLanesSeqInput{
		UpdateFeeQuoterDestsConfig: operation.UpdateFeeQuoterDestsInput{
			MCMSAddress: aptos.AccountAddress{},
			Updates:     map[uint64]aptosfeequoter.DestChainConfig{},
		},
		UpdateFeeQuoterPricesConfig: operation.UpdateFeeQuoterPricesInput{
			MCMSAddress: aptos.AccountAddress{},
			Prices:      operation.FeeQuoterPriceUpdatePerSource{},
		},
		UpdateOnRampDestsConfig: operation.UpdateOnRampDestsInput{
			MCMSAddress: aptos.AccountAddress{},
			Updates:     map[uint64]v1_6.OnRampDestinationUpdate{},
		},
		UpdateOffRampSourcesConfig: operation.UpdateOffRampSourcesInput{
			MCMSAddress: aptos.AccountAddress{},
			Updates:     map[uint64]v1_6.OffRampSourceUpdate{},
		},
	}
	isSer := operations.IsSerializable(logger.TestLogger(t), input)
	fmt.Println(isSer)
}

func TestAddAptosLanes_Apply(t *testing.T) {
	// Setup environment and config
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	emvSelector := env.AllChainSelectors()[0]
	emvSelector2 := env.AllChainSelectors()[1]
	aptosSelector := uint64(4457093679053095497)

	// Get chain selectors
	aptosChainSelectors := env.AllChainSelectorsAptos()
	require.Equal(t, 1, len(aptosChainSelectors), "Expected exactly 1 Aptos chain ")
	chainSelector := aptosChainSelectors[0]
	t.Log("Deployer: ", env.AptosChains[chainSelector].DeployerSigner)

	// Deploy Lane
	cfg := getMockUpdateConfig(t, emvSelector, emvSelector2, aptosSelector)

	// Apply the changeset
	env, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.AddAptosLanes{}, cfg),
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	// bind ccip aptos
	aptosCCIPAddr := state.AptosChains[aptosSelector].CCIPAddress
	aptosCCIP := ccip.Bind(aptosCCIPAddr, env.AptosChains[aptosSelector].Client)

	is_enabled1, sequence_number1, allowlist_enabled1, err := aptosCCIP.Onramp().GetDestChainConfig(nil, emvSelector)
	require.NoError(t, err)
	require.True(t, is_enabled1)
	require.True(t, sequence_number1 > 0)
	require.True(t, allowlist_enabled1)

	is_enabled2, sequence_number2, allowlist_enabled2, err := aptosCCIP.Onramp().GetDestChainConfig(nil, emvSelector2)
	require.NoError(t, err)
	require.True(t, is_enabled2)
	require.True(t, sequence_number2 > 0)
	require.True(t, allowlist_enabled2)
}

func getMockUpdateConfig(
	t *testing.T,
	emvSelector,
	emvSelector2,
	aptosSelector uint64,
) config.UpdateAptosLanesConfig {
	return config.UpdateAptosLanesConfig{
		MCMSConfig: nil,
		// Aptos1 <> EVM1 | Aptos1 -> EVM2
		Lanes: []config.LaneConfig{
			{
				Source: config.AptosChainDefinition{
					Selector:                 aptosSelector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: aptosTestDestFeeQuoterConfig(t),
				},
				Dest: config.EVMChainDefinition{
					ChainDefinition: v1_6.ChainDefinition{
						Selector:                 emvSelector,
						GasPrice:                 big.NewInt(1e17),
						TokenPrices:              map[common.Address]*big.Int{},
						FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
					},
				},
				IsDisabled: false,
			},
			{
				Source: config.EVMChainDefinition{
					ChainDefinition: v1_6.ChainDefinition{
						Selector:                 emvSelector,
						GasPrice:                 big.NewInt(1e17),
						TokenPrices:              map[common.Address]*big.Int{},
						FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
					},
				},
				Dest: config.AptosChainDefinition{
					Selector:                 aptosSelector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: aptosTestDestFeeQuoterConfig(t),
				},
				IsDisabled: false,
			},
			{
				Source: config.AptosChainDefinition{
					Selector:                 aptosSelector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: aptosTestDestFeeQuoterConfig(t),
				},
				Dest: config.EVMChainDefinition{
					ChainDefinition: v1_6.ChainDefinition{
						Selector:                 emvSelector2,
						GasPrice:                 big.NewInt(1e17),
						TokenPrices:              map[common.Address]*big.Int{},
						FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
					},
				},
				IsDisabled: false,
			},
		},
		TestRouter: true,
	}
}

// TODO: Deduplicate these test helpers
func aptosTestDestFeeQuoterConfig(t *testing.T) aptosfeequoter.DestChainConfig {
	return aptosfeequoter.DestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           11,
		MaxDataBytes:                      40_000,
		MaxPerMsgGasLimit:                 4_000_000,
		DestGasOverhead:                   ccipevm.DestGasOverhead,
		DefaultTokenFeeUsdCents:           30,
		DestGasPerPayloadByteBase:         ccipevm.CalldataGasPerByteBase,
		DestGasPerPayloadByteHigh:         ccipevm.CalldataGasPerByteHigh,
		DestGasPerPayloadByteThreshold:    ccipevm.CalldataGasPerByteThreshold,
		DestDataAvailabilityOverheadGas:   700,
		DestGasPerDataAvailabilityByte:    17,
		DestDataAvailabilityMultiplierBps: 2,
		DefaultTokenDestGasOverhead:       100_000,
		DefaultTxGasLimit:                 100_000,
		GasMultiplierWeiPerEth:            12e17,
		NetworkFeeUsdCents:                20,
		ChainFamilySelector:               hexMustDecode(t, v1_6.AptosFamilySelector),
		EnforceOutOfOrder:                 false,
		GasPriceStalenessThreshold:        2,
	}
}

func hexMustDecode(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}
