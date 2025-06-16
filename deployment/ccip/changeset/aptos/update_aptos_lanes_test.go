package aptos_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptosfeequoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_router"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func TestAddAptosLanes_Apply(t *testing.T) {
	// Setup environment and config
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	emvSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	emvSelector2 := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[1]
	aptosSelector := uint64(4457093679053095497)

	// Get chain selectors
	aptosChainSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))
	require.Len(t, aptosChainSelectors, 1, "Expected exactly 1 Aptos chain ")
	chainSelector := aptosChainSelectors[0]
	t.Log("Deployer: ", env.BlockChains.AptosChains()[chainSelector].DeployerSigner)

	// Deploy Lane
	cfg := getMockUpdateConfig(t, emvSelector, emvSelector2, aptosSelector)

	// Apply the changeset
	env, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.AddAptosLanes{}, cfg),
	})
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	// bind ccip aptos
	aptosCCIPAddr := state.AptosChains[aptosSelector].CCIPAddress
	aptosOnRamp := ccip_onramp.Bind(aptosCCIPAddr, env.BlockChains.AptosChains()[aptosSelector].Client)
	aptosOffRamp := ccip_offramp.Bind(aptosCCIPAddr, env.BlockChains.AptosChains()[aptosSelector].Client)
	aptosRouter := ccip_router.Bind(aptosCCIPAddr, env.BlockChains.AptosChains()[aptosSelector].Client)

	dynCfg, err := aptosOffRamp.Offramp().GetDynamicConfig(&bind.CallOpts{})
	require.NoError(t, err)
	require.Positive(t, dynCfg.PermissionlessExecutionThresholdSeconds)

	isSupported, err := aptosOnRamp.Onramp().IsChainSupported(&bind.CallOpts{}, emvSelector)
	require.NoError(t, err)
	require.True(t, isSupported)

	_, _, router, err := aptosOnRamp.Onramp().GetDestChainConfig(&bind.CallOpts{}, emvSelector)
	require.NoError(t, err)
	require.NotEqual(t, aptos.AccountAddress{}, router)

	_, _, router2, err := aptosOnRamp.Onramp().GetDestChainConfig(&bind.CallOpts{}, emvSelector2)
	require.NoError(t, err)
	require.NotEqual(t, aptos.AccountAddress{}, router2)

	versions, err := aptosRouter.Router().GetOnRampVersions(&bind.CallOpts{}, []uint64{emvSelector, emvSelector2})
	require.NoError(t, err)
	require.ElementsMatch(t, versions, [][]byte{{1, 6, 1}, {1, 6, 0}})
}

func getMockUpdateConfig(
	t *testing.T,
	emvSelector,
	emvSelector2,
	aptosSelector uint64,
) config.UpdateAptosLanesConfig {
	linkToken := aptos.AccountAddress{}
	_ = linkToken.ParseStringRelaxed("0x3b17dad1bdd88f337712cc2f6187bb741d56da467320373fd9198262cc93de76")
	otherToken := aptos.AccountAddress{}
	_ = linkToken.ParseStringRelaxed("0xa")

	return config.UpdateAptosLanesConfig{
		EVMMCMSConfig: nil,
		AptosMCMSConfig: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
		// Aptos1 <> EVM1 + Aptos1 -> EVM2
		Lanes: []config.LaneConfig{
			{
				Source: config.AptosChainDefinition{
					Selector:                 aptosSelector,
					GasPrice:                 big.NewInt(1e17),
					FeeQuoterDestChainConfig: aptosTestDestFeeQuoterConfig(t),
					ConnectionConfig: v1_6.ConnectionConfig{
						RMNVerificationDisabled: true,
						AllowListEnabled:        false,
					},
					AddTokenTransferFeeConfigs: []aptosfeequoter.TokenTransferFeeConfigAdded{
						{
							DestChainSelector: emvSelector,
							Token:             linkToken,
							TokenTransferFeeConfig: aptosfeequoter.TokenTransferFeeConfig{
								MinFeeUsdCents:    1,
								MaxFeeUsdCents:    10000,
								DeciBps:           0,
								DestGasOverhead:   1000,
								DestBytesOverhead: 1000,
								IsEnabled:         true,
							},
						},
						{
							DestChainSelector: emvSelector,
							Token:             otherToken,
							TokenTransferFeeConfig: aptosfeequoter.TokenTransferFeeConfig{
								MinFeeUsdCents:    1,
								MaxFeeUsdCents:    10000,
								DeciBps:           0,
								DestGasOverhead:   1000,
								DestBytesOverhead: 1000,
								IsEnabled:         true,
							},
						},
					},
				},
				Dest: config.EVMChainDefinition{
					ChainDefinition: v1_6.ChainDefinition{
						Selector:                 emvSelector,
						GasPrice:                 big.NewInt(1e17),
						TokenPrices:              map[common.Address]*big.Int{},
						FeeQuoterDestChainConfig: v1_6.DefaultFeeQuoterDestChainConfig(true),
						ConnectionConfig: v1_6.ConnectionConfig{
							RMNVerificationDisabled: true,
							AllowListEnabled:        false,
						},
					},
					OnRampVersion: []byte{1, 6, 1},
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
						ConnectionConfig: v1_6.ConnectionConfig{
							RMNVerificationDisabled: true,
							AllowListEnabled:        false,
						},
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
		TestRouter: false,
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
		GasMultiplierWeiPerEth:            1e7,
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
