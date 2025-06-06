package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solToken "github.com/gagliardetto/solana-go/programs/token"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// token setup
func deployTokenAndMint(t *testing.T, tenv cldf.Environment, solChain uint64, walletPubKeys []string) (cldf.Environment, solana.PublicKey, error) {
	mintMap := make(map[string]uint64)
	for _, key := range walletPubKeys {
		mintMap[key] = uint64(1000)
	}
	e, err := commonchangeset.Apply(t, tenv,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeploySolanaToken),
			ccipChangesetSolana.DeploySolanaTokenConfig{
				ChainSelector:       solChain,
				TokenProgramName:    shared.SPLTokens,
				TokenDecimals:       9,
				TokenSymbol:         "TEST_TOKEN",
				ATAList:             walletPubKeys,
				MintAmountToAddress: mintMap,
			},
		),
	)
	if err != nil {
		return cldf.Environment{}, solana.PublicKey{}, err
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(solChain) //nolint:staticcheck // addressbook still valid
	require.NoError(t, err)
	tokenAddress := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPLTokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("TEST_TOKEN"),
		},
		addresses,
	)
	return e, tokenAddress, err
}

// remote chain setup
func TestAddRemoteChainWithMcms(t *testing.T) {
	t.Parallel()
	doTestAddRemoteChain(t, true)
}

func TestAddRemoteChainWithoutMcms(t *testing.T) {
	t.Parallel()
	doTestAddRemoteChain(t, false)
}

func doTestAddRemoteChain(t *testing.T, mcms bool) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	e := tenv.Env
	_, err := stateview.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)
	evmChains := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	var mcmsConfig *proposalutils.TimelockConfig
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
	}
	onRampUpdates := make(map[uint64]map[uint64]v1_6.OnRampDestinationUpdate)
	for _, evmChain := range evmChains {
		onRampUpdates[evmChain] = map[uint64]v1_6.OnRampDestinationUpdate{
			solChain: {
				IsEnabled:        true,
				AllowListEnabled: false,
			},
		}
	}
	routerUpdates := make(map[uint64]*ccipChangesetSolana.RouterConfig)
	for _, evmChain := range evmChains {
		routerUpdates[evmChain] = &ccipChangesetSolana.RouterConfig{
			RouterDestinationConfig: solRouter.DestChainConfig{
				AllowListEnabled: true,
			},
		}
	}
	offRampUpdates := make(map[uint64]*ccipChangesetSolana.OffRampConfig)
	for _, evmChain := range evmChains {
		offRampUpdates[evmChain] = &ccipChangesetSolana.OffRampConfig{
			EnabledAsSource: true,
		}
	}
	feeQuoterUpdates := make(map[uint64]*ccipChangesetSolana.FeeQuoterConfig)
	for _, evmChain := range evmChains {
		feeQuoterUpdates[evmChain] = &ccipChangesetSolana.FeeQuoterConfig{
			FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
				IsEnabled:                   true,
				DefaultTxGasLimit:           200000,
				MaxPerMsgGasLimit:           3000000,
				MaxDataBytes:                30000,
				MaxNumberOfTokensPerMsg:     5,
				DefaultTokenDestGasOverhead: 5000,
				ChainFamilySelector:         [4]uint8{40, 18, 213, 44},
			},
		}
	}
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: onRampUpdates,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector:  solChain,
				UpdatesByChain: routerUpdates,
				MCMS:           mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
			ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector:  solChain,
				UpdatesByChain: feeQuoterUpdates,
				MCMS:           mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
			ccipChangesetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector:  solChain,
				UpdatesByChain: offRampUpdates,
				MCMS:           mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	var offRampSourceChain solOffRamp.SourceChain
	var destChainStateAccount solRouter.DestChain
	var destChainFqAccount solFeeQuoter.DestChain
	var offRampEvmSourceChainPDA solana.PublicKey
	var evmDestChainStatePDA solana.PublicKey
	var fqEvmDestChainPDA solana.PublicKey
	for _, evmChain := range evmChains {
		offRampEvmSourceChainPDA, _, _ = solState.FindOfframpSourceChainPDA(evmChain, state.SolChains[solChain].OffRamp)
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
		require.NoError(t, err)
		require.True(t, offRampSourceChain.Config.IsEnabled)

		evmDestChainStatePDA = state.SolChains[solChain].DestChainStatePDAs[evmChain]
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
		require.True(t, destChainStateAccount.Config.AllowListEnabled)
		require.NoError(t, err)

		fqEvmDestChainPDA, _, _ = solState.FindFqDestChainPDA(evmChain, state.SolChains[solChain].FeeQuoter)
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
		require.NoError(t, err, "failed to get account info")
		require.Equal(t, solFeeQuoter.TimestampedPackedU224{}, destChainFqAccount.State.UsdPerUnitGas)
		require.True(t, destChainFqAccount.Config.IsEnabled)
	}

	// Disable the chain

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DisableRemoteChain),
			ccipChangesetSolana.DisableRemoteChainConfig{
				ChainSelector: solChain,
				RemoteChains:  evmChains,
				MCMS:          mcmsConfig,
			},
		),
	},
	)

	require.NoError(t, err)

	state, err = stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
	require.NoError(t, err)
	require.False(t, offRampSourceChain.Config.IsEnabled)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
	require.True(t, destChainStateAccount.Config.AllowListEnabled)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.False(t, destChainFqAccount.Config.IsEnabled)

	// Re-enable the chain

	routerUpdates = make(map[uint64]*ccipChangesetSolana.RouterConfig)
	for _, evmChain := range evmChains {
		routerUpdates[evmChain] = &ccipChangesetSolana.RouterConfig{
			RouterDestinationConfig: solRouter.DestChainConfig{
				AllowListEnabled: false,
			},
		}
	}
	offRampUpdates = make(map[uint64]*ccipChangesetSolana.OffRampConfig)
	for _, evmChain := range evmChains {
		offRampUpdates[evmChain] = &ccipChangesetSolana.OffRampConfig{
			EnabledAsSource: true,
		}
	}
	feeQuoterUpdates = make(map[uint64]*ccipChangesetSolana.FeeQuoterConfig)
	for _, evmChain := range evmChains {
		feeQuoterUpdates[evmChain] = &ccipChangesetSolana.FeeQuoterConfig{
			FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
				IsEnabled:                   true,
				DefaultTxGasLimit:           200000,
				MaxPerMsgGasLimit:           3000000,
				MaxDataBytes:                30000,
				MaxNumberOfTokensPerMsg:     5,
				DefaultTokenDestGasOverhead: 5000,
				ChainFamilySelector:         [4]uint8{40, 18, 213, 44},
			},
		}
	}

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector:  solChain,
				UpdatesByChain: routerUpdates,
				MCMS:           mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
			ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector:  solChain,
				UpdatesByChain: feeQuoterUpdates,
				MCMS:           mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
			ccipChangesetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector:  solChain,
				UpdatesByChain: offRampUpdates,
				MCMS:           mcmsConfig,
			},
		),
	},
	)

	require.NoError(t, err)

	state, err = stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
	require.NoError(t, err)
	require.True(t, offRampSourceChain.Config.IsEnabled)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
	require.False(t, destChainStateAccount.Config.AllowListEnabled)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.True(t, destChainFqAccount.Config.IsEnabled)
}

// billing test
func doTestBilling(t *testing.T, mcms bool) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]

	e, tokenAddress, err := deployTokenAndMint(t, tenv.Env, solChain, []string{})
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	validTimestamp := int64(100)
	value := [28]uint8{}
	bigNum, ok := new(big.Int).SetString("19816680000000000000", 10)
	require.True(t, ok)
	bigNum.FillBytes(value[:])
	var mcmsConfig *proposalutils.TimelockConfig
	testPriceUpdater := e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
		testPriceUpdater, err = ccipChangesetSolana.FetchTimelockSigner(e, solChain)
		require.NoError(t, err)
	}
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddBillingTokenChangeset),
			ccipChangesetSolana.BillingTokenConfig{
				ChainSelector: solChain,
				TokenPubKey:   tokenAddress.String(),
				Config: solFeeQuoter.BillingTokenConfig{
					Enabled: true,
					Mint:    tokenAddress,
					UsdPerToken: solFeeQuoter.TimestampedPackedU224{
						Timestamp: validTimestamp,
						Value:     value,
					},
					PremiumMultiplierWeiPerEth: 100,
				},
				MCMS: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenTransferFeeForRemoteChain),
			ccipChangesetSolana.TokenTransferFeeForRemoteChainConfig{
				ChainSelector:       solChain,
				RemoteChainSelector: evmChain,
				TokenPubKey:         tokenAddress.String(),
				Config: solFeeQuoter.TokenTransferFeeConfig{
					MinFeeUsdcents:    800,
					MaxFeeUsdcents:    1600,
					DeciBps:           0,
					DestGasOverhead:   100,
					DestBytesOverhead: 100,
					IsEnabled:         true,
				},
				MCMS: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	billingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenAddress, state.SolChains[solChain].FeeQuoter)
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
	require.NoError(t, err)
	require.True(t, token0ConfigAccount.Config.Enabled)
	require.Equal(t, tokenAddress, token0ConfigAccount.Config.Mint)
	require.Equal(t, uint64(100), token0ConfigAccount.Config.PremiumMultiplierWeiPerEth)

	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(evmChain, tokenAddress, state.SolChains[solChain].FeeQuoter)
	var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), remoteBillingPDA, &remoteBillingAccount)
	require.NoError(t, err)
	require.Equal(t, tokenAddress, remoteBillingAccount.Mint)
	require.Equal(t, uint32(800), remoteBillingAccount.TokenTransferConfig.MinFeeUsdcents)

	// test update
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddBillingTokenChangeset),
			ccipChangesetSolana.BillingTokenConfig{
				ChainSelector: solChain,
				TokenPubKey:   tokenAddress.String(),
				Config: solFeeQuoter.BillingTokenConfig{
					Enabled: true,
					Mint:    tokenAddress,
					UsdPerToken: solFeeQuoter.TimestampedPackedU224{
						Timestamp: validTimestamp,
						Value:     value,
					},
					PremiumMultiplierWeiPerEth: 200,
				},
				MCMS: mcmsConfig,
			},
		),
	})
	require.NoError(t, err)
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
	require.NoError(t, err)
	require.Equal(t, uint64(200), token0ConfigAccount.Config.PremiumMultiplierWeiPerEth)
	feeAggregatorPriv, _ := solana.NewRandomPrivateKey()
	feeAggregator := feeAggregatorPriv.PublicKey()

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]*ccipChangesetSolana.RouterConfig{
					evmChain: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
						},
					},
				},
				MCMS: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
			ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]*ccipChangesetSolana.FeeQuoterConfig{
					evmChain: {
						FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
							IsEnabled:                   true,
							DefaultTxGasLimit:           200000,
							MaxPerMsgGasLimit:           3000000,
							MaxDataBytes:                30000,
							MaxNumberOfTokensPerMsg:     5,
							DefaultTokenDestGasOverhead: 5000,
							ChainFamilySelector:         [4]uint8{40, 18, 213, 44},
						},
					},
				},
				MCMS: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
			ccipChangesetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]*ccipChangesetSolana.OffRampConfig{
					evmChain: {
						EnabledAsSource: true,
					},
				},
				MCMS: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.ModifyPriceUpdater),
			ccipChangesetSolana.ModifyPriceUpdaterConfig{
				ChainSelector:      solChain,
				PriceUpdater:       testPriceUpdater,
				PriceUpdaterAction: ccipChangesetSolana.AddUpdater,
				MCMS:               mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.UpdatePrices),
			ccipChangesetSolana.UpdatePricesConfig{
				ChainSelector: solChain,
				TokenPriceUpdates: []solFeeQuoter.TokenPriceUpdate{
					{
						SourceToken: tokenAddress,
						UsdPerToken: solCommonUtil.To28BytesBE(123),
					},
				},
				GasPriceUpdates: []solFeeQuoter.GasPriceUpdate{
					{
						DestChainSelector: evmChain,
						UsdPerUnitGas:     solCommonUtil.To28BytesBE(345),
					},
				},
				PriceUpdater: testPriceUpdater,
				MCMS:         mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.ModifyPriceUpdater),
			ccipChangesetSolana.ModifyPriceUpdaterConfig{
				ChainSelector:      solChain,
				PriceUpdater:       testPriceUpdater,
				PriceUpdaterAction: ccipChangesetSolana.RemoveUpdater,
				MCMS:               mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
			ccipChangesetSolana.SetFeeAggregatorConfig{
				ChainSelector: solChain,
				FeeAggregator: feeAggregator.String(),
				MCMS:          mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.CreateSolanaTokenATA),
			ccipChangesetSolana.CreateSolanaTokenATAConfig{
				ChainSelector: solChain,
				TokenPubkey:   tokenAddress,
				ATAList:       []string{feeAggregator.String()}, // create ATA for the fee aggregator
			},
		),
	},
	)
	require.NoError(t, err)

	// just send funds to the router manually rather than run e2e
	billingSignerPDA, _, _ := solState.FindFeeBillingSignerPDA(state.SolChains[solChain].Router)
	billingSignerATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, tokenAddress, billingSignerPDA)
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.MintSolanaToken),
			ccipChangesetSolana.MintSolanaTokenConfig{
				ChainSelector: solChain,
				TokenPubkey:   tokenAddress.String(),
				AmountToAddress: map[string]uint64{
					billingSignerPDA.String(): uint64(1000),
				},
			},
		),
	},
	)
	require.NoError(t, err)
	// check that the billing account has the right amount
	_, billingResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, billingSignerATA, cldf_solana.SolDefaultCommitment)
	require.NoError(t, err)
	require.Equal(t, 1000, billingResult)
	feeAggregatorATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, tokenAddress, feeAggregator)
	_, feeAggResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, feeAggregatorATA, cldf_solana.SolDefaultCommitment)
	require.NoError(t, err)
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.WithdrawBilledFunds),
			ccipChangesetSolana.WithdrawBilledFundsConfig{
				ChainSelector: solChain,
				TransferAll:   true,
				TokenPubKey:   tokenAddress.String(),
				MCMS:          mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)
	_, newBillingResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, billingSignerATA, cldf_solana.SolDefaultCommitment)
	require.NoError(t, err)
	require.Equal(t, billingResult-1000, newBillingResult)
	_, newFeeAggResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, feeAggregatorATA, cldf_solana.SolDefaultCommitment)
	require.NoError(t, err)
	require.Equal(t, feeAggResult+1000, newFeeAggResult)
}

func TestBillingWithMcms(t *testing.T) {
	t.Parallel()
	doTestBilling(t, true)
}

func TestBillingWithoutMcms(t *testing.T) {
	t.Parallel()
	doTestBilling(t, false)
}

// token admin registry test
func doTestTokenAdminRegistry(t *testing.T, mcms bool) {
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	e, tokenAddress, err := deployTokenAndMint(t, tenv.Env, solChain, []string{})
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	linkTokenAddress := state.SolChains[solChain].LinkToken
	newAdminNonTimelock, _ := solana.NewRandomPrivateKey()
	newAdmin := newAdminNonTimelock.PublicKey()
	newTokenAdmin := e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()

	var mcmsConfig *proposalutils.TimelockConfig
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
		timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
		require.NoError(t, err)
		newAdmin = timelockSignerPDA
		newTokenAdmin = timelockSignerPDA
	}
	timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
	require.NoError(t, err)

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// register token admin registry for tokenAddress via admin instruction
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
			ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             tokenAddress,
				TokenAdminRegistryAdmin: newAdmin.String(),
				RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
				MCMS:                    mcmsConfig,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetTokenAuthority),
			ccipChangesetSolana.SetTokenAuthorityConfig{
				ChainSelector: solChain,
				AuthorityType: solToken.AuthorityMintTokens,
				TokenPubkey:   linkTokenAddress,
				NewAuthority:  newTokenAdmin,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetTokenAuthority),
			ccipChangesetSolana.SetTokenAuthorityConfig{
				ChainSelector: solChain,
				AuthorityType: solToken.AuthorityFreezeAccount,
				TokenPubkey:   linkTokenAddress,
				NewAuthority:  newTokenAdmin,
			},
		),
		commonchangeset.Configure(
			// register token admin registry for linkToken via owner instruction
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
			ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             linkTokenAddress,
				TokenAdminRegistryAdmin: newAdmin.String(),
				RegisterType:            ccipChangesetSolana.ViaOwnerInstruction,
				MCMS:                    mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)
	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.Administrator)
	// pending administrator should be the proposed admin key
	require.Equal(t, newAdmin, tokenAdminRegistryAccount.PendingAdministrator)

	linkTokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(linkTokenAddress, state.SolChains[solChain].Router)
	var linkTokenAdminRegistryAccount solCommon.TokenAdminRegistry
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, linkTokenAdminRegistryPDA, &linkTokenAdminRegistryAccount)
	require.NoError(t, err)
	require.Equal(t, newAdmin, linkTokenAdminRegistryAccount.PendingAdministrator)

	// While we can assign the admin role arbitrarily regardless of mcms, we can only accept it as timelock
	if mcms {
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				// accept admin role for tokenAddress
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
				ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
					ChainSelector: solChain,
					TokenPubKey:   tokenAddress,
					MCMS:          mcmsConfig,
				},
			),
		)
		require.NoError(t, err)
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
		require.NoError(t, err)
		// confirm that the administrator is the deployer key
		require.Equal(t, timelockSignerPDA, tokenAdminRegistryAccount.Administrator)
		require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.PendingAdministrator)

		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				// transfer admin role for tokenAddress
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferAdminRoleTokenAdminRegistry),
				ccipChangesetSolana.TransferAdminRoleTokenAdminRegistryConfig{
					ChainSelector:             solChain,
					TokenPubKey:               tokenAddress.String(),
					NewRegistryAdminPublicKey: newAdminNonTimelock.PublicKey().String(),
					MCMS:                      mcmsConfig,
				},
			),
		},
		)
		require.NoError(t, err)
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
		require.NoError(t, err)
		require.Equal(t, newAdminNonTimelock.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)
	}
}

func TestTokenAdminRegistryWithMcms(t *testing.T) {
	t.Parallel()
	doTestTokenAdminRegistry(t, true)
}

func TestTokenAdminRegistryWithoutMcms(t *testing.T) {
	t.Parallel()
	doTestTokenAdminRegistry(t, false)
}

// pool lookup table test
func doTestPoolLookupTable(t *testing.T, e cldf.Environment, mcms bool, tokenMetadata string) {
	ctx := testcontext.Get(t)

	solChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]

	var mcmsConfig *proposalutils.TimelockConfig
	newAdmin := e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
		timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
		require.NoError(t, err)
		newAdmin = timelockSignerPDA
	}

	e, tokenAddress, err := deployTokenAndMint(t, e, solChain, []string{})
	require.NoError(t, err)
	pool := solTestTokenPool.LockAndRelease_PoolType
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// add token pool lookup table
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolLookupTable),
			ccipChangesetSolana.TokenPoolLookupTableConfig{
				ChainSelector: solChain,
				TokenPubKey:   tokenAddress,
				PoolType:      &pool,
				Metadata:      tokenMetadata,
			},
		),
	)
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	lookupTablePubKey := state.SolChains[solChain].TokenPoolLookupTable[tokenAddress][pool][tokenMetadata]

	lookupTableEntries0, err := solCommonUtil.GetAddressLookupTable(ctx, e.BlockChains.SolanaChains()[solChain].Client, lookupTablePubKey)
	require.NoError(t, err)
	require.Equal(t, lookupTablePubKey, lookupTableEntries0[0])
	require.Equal(t, tokenAddress, lookupTableEntries0[7])

	e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
		// register token admin registry for linkToken via owner instruction
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
		ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
			ChainSelector:           solChain,
			TokenPubKey:             tokenAddress,
			TokenAdminRegistryAdmin: newAdmin.String(),
			RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
			MCMS:                    mcmsConfig,
		},
	), commonchangeset.Configure(
		// accept admin role for tokenAddress
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
		ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
			ChainSelector: solChain,
			TokenPubKey:   tokenAddress,
			MCMS:          mcmsConfig,
		},
	), commonchangeset.Configure(
		// set pool -> this updates tokenAdminRegistryPDA, hence above changeset is required
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetPool),
		ccipChangesetSolana.SetPoolConfig{
			ChainSelector:   solChain,
			TokenPubKey:     tokenAddress,
			PoolType:        &pool,
			Metadata:        tokenMetadata,
			WritableIndexes: []uint8{3, 4, 7},
			MCMS:            mcmsConfig,
		},
	))
	require.NoError(t, err)
	tokenAdminRegistry := solCommon.TokenAdminRegistry{}
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)

	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistry)
	require.NoError(t, err)
	require.Equal(t, newAdmin, tokenAdminRegistry.Administrator)
	require.Equal(t, lookupTablePubKey, tokenAdminRegistry.LookupTable)
}

func TestPoolLookupTableWithMcms(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	doTestPoolLookupTable(t, tenv.Env, true, shared.CLLMetadata)
}

func TestPoolLookupTableWithoutMcms(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	doTestPoolLookupTable(t, tenv.Env, false, shared.CLLMetadata)
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 1)
}

// ocr3 test
func TestSetOcr3Active(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithNumOfNodes(16),
		testhelpers.WithNumOfBootstrapNodes(3),
		testhelpers.WithSolChains(1))
	var err error
	evmSelectors := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	homeChainSel := evmSelectors[0]
	solChainSelectors := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	_, _ = testhelpers.TransferOwnershipSolana(t, &tenv.Env, solChainSelectors[0], true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})

	tenv.Env, _, err = commonchangeset.ApplyChangesets(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetOCR3ConfigSolana),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       homeChainSel,
				RemoteChainSels:    solChainSelectors,
				CCIPHomeConfigType: globals.ConfigTypeActive,
				MCMS:               &proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
			},
		),
	})
	require.NoError(t, err)
}

func TestSetOcr3Candidate(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t,
		testhelpers.WithSolChains(1))
	var err error
	evmSelectors := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	homeChainSel := evmSelectors[0]
	solChainSelectors := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))
	_, _ = testhelpers.TransferOwnershipSolana(t, &tenv.Env, solChainSelectors[0], true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})

	tenv.Env, _, err = commonchangeset.ApplyChangesets(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetOCR3ConfigSolana),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       homeChainSel,
				RemoteChainSels:    solChainSelectors,
				CCIPHomeConfigType: globals.ConfigTypeCandidate,
				MCMS:               &proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
			},
		),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid OCR3 config state, expected candidate config")
}
