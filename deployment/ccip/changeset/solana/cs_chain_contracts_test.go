package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	solToken "github.com/gagliardetto/solana-go/programs/token"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func deployToken(t *testing.T, tenv deployment.Environment, solChain uint64) (deployment.Environment, solana.PublicKey, error) {
	e, err := commonchangeset.Apply(t, tenv, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DeploySolanaToken),
			ccipChangesetSolana.DeploySolanaTokenConfig{
				ChainSelector:    solChain,
				TokenProgramName: ccipChangeset.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      "TEST_TOKEN",
			},
		),
	)
	require.NoError(t, err)
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain].SPL2022Tokens[0]
	return e, tokenAddress, err
}

func TestAddRemoteChain(t *testing.T) {
	t.Parallel()
	// Default env just has 2 chains with all contracts
	// deployed but no lanes.
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	evmChain2 := tenv.Env.AllChainSelectors()[1]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	_, err := ccipChangeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	doTestAddRemoteChain(t, tenv.Env, evmChain, solChain, false)
	doTestAddRemoteChain(t, tenv.Env, evmChain2, solChain, true)
}

func doTestAddRemoteChain(t *testing.T, e deployment.Environment, evmChain uint64, solChain uint64, mcms bool) {
	var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
	var err error
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
			MCMS: &ccipChangeset.MCMSConfig{
				MinDelay: 1 * time.Second,
			},
			RouterOwnedByTimelock:    true,
			FeeQuoterOwnedByTimelock: true,
			OffRampOwnedByTimelock:   true,
		}
	}
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
			v1_6.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]v1_6.OnRampDestinationUpdate{
					evmChain: {
						solChain: {
							IsEnabled:        true,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.RouterConfig{
					evmChain: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
						},
					},
				},
				MCMSSolana: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
			ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.FeeQuoterConfig{
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
				MCMSSolana: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
			ccipChangesetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.OffRampConfig{
					evmChain: {
						EnabledAsSource: true,
					},
				},
				MCMSSolana: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	var offRampSourceChain solOffRamp.SourceChain
	offRampEvmSourceChainPDA, _, _ := solState.FindOfframpSourceChainPDA(evmChain, state.SolChains[solChain].OffRamp)
	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
	require.NoError(t, err)
	require.True(t, offRampSourceChain.Config.IsEnabled)

	var destChainStateAccount solRouter.DestChain
	evmDestChainStatePDA := state.SolChains[solChain].DestChainStatePDAs[evmChain]
	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
	require.True(t, destChainStateAccount.Config.AllowListEnabled)
	require.NoError(t, err)

	var destChainFqAccount solFeeQuoter.DestChain
	fqEvmDestChainPDA, _, _ := solState.FindFqDestChainPDA(evmChain, state.SolChains[solChain].FeeQuoter)
	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.Equal(t, solFeeQuoter.TimestampedPackedU224{}, destChainFqAccount.State.UsdPerUnitGas)
	require.True(t, destChainFqAccount.Config.IsEnabled)

	// Disable the chain

	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.DisableRemoteChain),
			ccipChangesetSolana.DisableRemoteChainConfig{
				ChainSelector: solChain,
				RemoteChains:  []uint64{evmChain},
				MCMSSolana:    mcmsConfig,
			},
		),
	},
	)

	require.NoError(t, err)

	state, err = ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
	require.NoError(t, err)
	require.False(t, offRampSourceChain.Config.IsEnabled)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
	require.True(t, destChainStateAccount.Config.AllowListEnabled)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.False(t, destChainFqAccount.Config.IsEnabled)

	// Re-enable the chain

	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
			ccipChangesetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.RouterConfig{
					evmChain: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: false,
						},
						IsUpdate: true,
					},
				},
				MCMSSolana: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
			ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.FeeQuoterConfig{
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
						IsUpdate: true,
					},
				},
				MCMSSolana: mcmsConfig,
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
			ccipChangesetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.OffRampConfig{
					evmChain: {
						EnabledAsSource: true,
						IsUpdate:        true,
					},
				},
				MCMSSolana: mcmsConfig,
			},
		),
	},
	)

	require.NoError(t, err)

	state, err = ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), offRampEvmSourceChainPDA, &offRampSourceChain)
	require.NoError(t, err)
	require.True(t, offRampSourceChain.Config.IsEnabled)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
	require.False(t, destChainStateAccount.Config.AllowListEnabled)

	err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.True(t, destChainFqAccount.Config.IsEnabled)
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 1)
}

func TestSetOcr3(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	var err error
	evmSelectors := tenv.Env.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChainSelectors := tenv.Env.AllChainSelectorsSolana()
	_, _ = testhelpers.TransferOwnershipSolana(t, &tenv.Env, solChainSelectors[0], true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})

	tenv.Env, err = commonchangeset.ApplyChangesetsV2(t, tenv.Env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetOCR3ConfigSolana),
			v1_6.SetOCR3OffRampConfig{
				HomeChainSel:       homeChainSel,
				RemoteChainSels:    solChainSelectors,
				CCIPHomeConfigType: globals.ConfigTypeActive,
				MCMS:               &ccipChangeset.MCMSConfig{MinDelay: 1 * time.Second},
			},
		),
	})
	require.NoError(t, err)
}

func TestBilling(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Msg  string
		Mcms bool
	}{
		{
			Msg:  "with mcms",
			Mcms: true,
		},
		{
			Msg:  "without mcms",
			Mcms: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

			evmChain := tenv.Env.AllChainSelectors()[0]
			solChain := tenv.Env.AllChainSelectorsSolana()[0]

			e, tokenAddress, err := deployToken(t, tenv.Env, solChain)
			require.NoError(t, err)
			state, err := ccipChangeset.LoadOnchainStateSolana(e)
			require.NoError(t, err)
			validTimestamp := int64(100)
			value := [28]uint8{}
			bigNum, ok := new(big.Int).SetString("19816680000000000000", 10)
			require.True(t, ok)
			bigNum.FillBytes(value[:])
			var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
			testPriceUpdater := e.SolChains[solChain].DeployerKey.PublicKey()
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				}
				testPriceUpdater, err = ccipChangesetSolana.FetchTimelockSigner(e, solChain)
				require.NoError(t, err)
			}
			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddBillingTokenChangeset),
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
						MCMSSolana: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenTransferFeeForRemoteChain),
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
						MCMSSolana: mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)

			billingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenAddress, state.SolChains[solChain].FeeQuoter)
			var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
			err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
			require.NoError(t, err)
			require.True(t, token0ConfigAccount.Config.Enabled)
			require.Equal(t, tokenAddress, token0ConfigAccount.Config.Mint)
			require.Equal(t, uint64(100), token0ConfigAccount.Config.PremiumMultiplierWeiPerEth)

			remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(evmChain, tokenAddress, state.SolChains[solChain].FeeQuoter)
			var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
			err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), remoteBillingPDA, &remoteBillingAccount)
			require.NoError(t, err)
			require.Equal(t, tokenAddress, remoteBillingAccount.Mint)
			require.Equal(t, uint32(800), remoteBillingAccount.TokenTransferConfig.MinFeeUsdcents)

			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddBillingTokenChangeset),
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
						MCMSSolana: mcmsConfig,
						IsUpdate:   true,
					},
				),
			})
			require.NoError(t, err)
			err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
			require.NoError(t, err)
			require.Equal(t, uint64(200), token0ConfigAccount.Config.PremiumMultiplierWeiPerEth)
			feeAggregatorPriv, _ := solana.NewRandomPrivateKey()
			feeAggregator := feeAggregatorPriv.PublicKey()

			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToRouter),
					ccipChangesetSolana.AddRemoteChainToRouterConfig{
						ChainSelector: solChain,
						UpdatesByChain: map[uint64]ccipChangesetSolana.RouterConfig{
							evmChain: {
								RouterDestinationConfig: solRouter.DestChainConfig{
									AllowListEnabled: true,
								},
							},
						},
						MCMSSolana: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToFeeQuoter),
					ccipChangesetSolana.AddRemoteChainToFeeQuoterConfig{
						ChainSelector: solChain,
						UpdatesByChain: map[uint64]ccipChangesetSolana.FeeQuoterConfig{
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
						MCMSSolana: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToOffRamp),
					ccipChangesetSolana.AddRemoteChainToOffRampConfig{
						ChainSelector: solChain,
						UpdatesByChain: map[uint64]ccipChangesetSolana.OffRampConfig{
							evmChain: {
								EnabledAsSource: true,
							},
						},
						MCMSSolana: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.ModifyPriceUpdater),
					ccipChangesetSolana.ModifyPriceUpdaterConfig{
						ChainSelector:      solChain,
						PriceUpdater:       testPriceUpdater,
						PriceUpdaterAction: ccipChangesetSolana.AddUpdater,
						MCMSSolana:         mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.UpdatePrices),
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
						MCMSSolana:   mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.ModifyPriceUpdater),
					ccipChangesetSolana.ModifyPriceUpdaterConfig{
						ChainSelector:      solChain,
						PriceUpdater:       testPriceUpdater,
						PriceUpdaterAction: ccipChangesetSolana.RemoveUpdater,
						MCMSSolana:         mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetFeeAggregator),
					ccipChangesetSolana.SetFeeAggregatorConfig{
						ChainSelector: solChain,
						FeeAggregator: feeAggregator.String(),
						MCMSSolana:    mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.CreateSolanaTokenATA),
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
			billingSignerATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, billingSignerPDA)
			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.MintSolanaToken),
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
			_, billingResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, billingSignerATA, deployment.SolDefaultCommitment)
			require.NoError(t, err)
			require.Equal(t, 1000, billingResult)
			feeAggregatorATA, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, feeAggregator)
			_, feeAggResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, feeAggregatorATA, deployment.SolDefaultCommitment)
			require.NoError(t, err)
			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.WithdrawBilledFunds),
					ccipChangesetSolana.WithdrawBilledFundsConfig{
						ChainSelector: solChain,
						TransferAll:   true,
						TokenPubKey:   tokenAddress.String(),
						MCMSSolana:    mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)
			_, newBillingResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, billingSignerATA, deployment.SolDefaultCommitment)
			require.NoError(t, err)
			require.Equal(t, billingResult-1000, newBillingResult)
			_, newFeeAggResult, err := solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, feeAggregatorATA, deployment.SolDefaultCommitment)
			require.NoError(t, err)
			require.Equal(t, feeAggResult+1000, newFeeAggResult)
		})
	}

}

func TestTokenAdminRegistry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Msg  string
		Mcms bool
	}{
		{
			Msg:  "with mcms",
			Mcms: true,
		},
		{
			Msg:  "without mcms",
			Mcms: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			solChain := tenv.Env.AllChainSelectorsSolana()[0]
			e, tokenAddress, err := deployToken(t, tenv.Env, solChain)
			require.NoError(t, err)
			state, err := ccipChangeset.LoadOnchainStateSolana(e)
			require.NoError(t, err)
			linkTokenAddress := state.SolChains[solChain].LinkToken
			newAdminNonTimelock, _ := solana.NewRandomPrivateKey()
			newAdmin := newAdminNonTimelock.PublicKey()
			newTokenAdmin := e.SolChains[solChain].DeployerKey.PublicKey()

			var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				}
				timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
				require.NoError(t, err)
				newAdmin = timelockSignerPDA
				newTokenAdmin = timelockSignerPDA
			}
			timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
			require.NoError(t, err)

			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					// register token admin registry for tokenAddress via admin instruction
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
					ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
						ChainSelector:           solChain,
						TokenPubKey:             tokenAddress.String(),
						TokenAdminRegistryAdmin: newAdmin.String(),
						RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
						MCMSSolana:              mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetTokenAuthority),
					ccipChangesetSolana.SetTokenAuthorityConfig{
						ChainSelector: solChain,
						AuthorityType: solToken.AuthorityMintTokens,
						TokenPubkey:   linkTokenAddress,
						NewAuthority:  newTokenAdmin,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetTokenAuthority),
					ccipChangesetSolana.SetTokenAuthorityConfig{
						ChainSelector: solChain,
						AuthorityType: solToken.AuthorityFreezeAccount,
						TokenPubkey:   linkTokenAddress,
						NewAuthority:  newTokenAdmin,
					},
				),
				commonchangeset.Configure(
					// register token admin registry for linkToken via owner instruction
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
					ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
						ChainSelector:           solChain,
						TokenPubKey:             linkTokenAddress.String(),
						TokenAdminRegistryAdmin: newAdmin.String(),
						RegisterType:            ccipChangesetSolana.ViaOwnerInstruction,
						MCMSSolana:              mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)

			tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)
			var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
			err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
			require.NoError(t, err)
			require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.Administrator)
			// pending administrator should be the proposed admin key
			require.Equal(t, newAdmin, tokenAdminRegistryAccount.PendingAdministrator)

			linkTokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(linkTokenAddress, state.SolChains[solChain].Router)
			var linkTokenAdminRegistryAccount solCommon.TokenAdminRegistry
			err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, linkTokenAdminRegistryPDA, &linkTokenAdminRegistryAccount)
			require.NoError(t, err)
			require.Equal(t, newAdmin, linkTokenAdminRegistryAccount.PendingAdministrator)

			// While we can assign the admin role arbitrarily regardless of mcms, we can only accept it as timelock
			if test.Mcms {
				e, err = commonchangeset.Apply(t, e, nil,
					commonchangeset.Configure(
						// accept admin role for tokenAddress
						deployment.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
						ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
							ChainSelector: solChain,
							TokenPubKey:   tokenAddress.String(),
							MCMSSolana:    mcmsConfig,
						},
					),
				)
				require.NoError(t, err)
				err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
				require.NoError(t, err)
				// confirm that the administrator is the deployer key
				require.Equal(t, timelockSignerPDA, tokenAdminRegistryAccount.Administrator)
				require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.PendingAdministrator)

				e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						// transfer admin role for tokenAddress
						deployment.CreateLegacyChangeSet(ccipChangesetSolana.TransferAdminRoleTokenAdminRegistry),
						ccipChangesetSolana.TransferAdminRoleTokenAdminRegistryConfig{
							ChainSelector:             solChain,
							TokenPubKey:               tokenAddress.String(),
							NewRegistryAdminPublicKey: newAdminNonTimelock.PublicKey().String(),
							MCMSSolana:                mcmsConfig,
						},
					),
				},
				)
				require.NoError(t, err)
				err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
				require.NoError(t, err)
				require.Equal(t, newAdminNonTimelock.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)
			}
		})
	}
}

func TestPoolLookupTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Msg  string
		Mcms bool
	}{
		{
			Msg:  "with mcms",
			Mcms: true,
		},
		{
			Msg:  "without mcms",
			Mcms: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			ctx := testcontext.Get(t)
			tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
			solChain := tenv.Env.AllChainSelectorsSolana()[0]

			var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
			newAdmin := tenv.Env.SolChains[solChain].DeployerKey.PublicKey()
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &tenv.Env, solChain, true,
					ccipChangesetSolana.CCIPContractsToTransfer{
						Router:    true,
						FeeQuoter: true,
						OffRamp:   true,
					})
				mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
					MCMS: &ccipChangeset.MCMSConfig{
						MinDelay: 1 * time.Second,
					},
					RouterOwnedByTimelock:    true,
					FeeQuoterOwnedByTimelock: true,
					OffRampOwnedByTimelock:   true,
				}
				timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(tenv.Env, solChain)
				require.NoError(t, err)
				newAdmin = timelockSignerPDA
			}

			e, tokenAddress, err := deployToken(t, tenv.Env, solChain)
			require.NoError(t, err)
			e, err = commonchangeset.Apply(t, e, nil,
				commonchangeset.Configure(
					// add token pool lookup table
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolLookupTable),
					ccipChangesetSolana.TokenPoolLookupTableConfig{
						ChainSelector: solChain,
						TokenPubKey:   tokenAddress.String(),
					},
				),
			)
			require.NoError(t, err)
			state, err := ccipChangeset.LoadOnchainStateSolana(e)
			require.NoError(t, err)
			lookupTablePubKey := state.SolChains[solChain].TokenPoolLookupTable[tokenAddress]

			lookupTableEntries0, err := solCommonUtil.GetAddressLookupTable(ctx, e.SolChains[solChain].Client, lookupTablePubKey)
			require.NoError(t, err)
			require.Equal(t, lookupTablePubKey, lookupTableEntries0[0])
			require.Equal(t, tokenAddress, lookupTableEntries0[7])

			e, err = commonchangeset.Apply(t, e, nil,
				commonchangeset.Configure(
					// register token admin registry for linkToken via owner instruction
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
					ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
						ChainSelector:           solChain,
						TokenPubKey:             tokenAddress.String(),
						TokenAdminRegistryAdmin: newAdmin.String(),
						RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
						MCMSSolana:              mcmsConfig,
					},
				),
				commonchangeset.Configure(
					// accept admin role for tokenAddress
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
					ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
						ChainSelector: solChain,
						TokenPubKey:   tokenAddress.String(),
						MCMSSolana:    mcmsConfig,
					},
				),
				commonchangeset.Configure(
					// set pool -> this updates tokenAdminRegistryPDA, hence above changeset is required
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetPool),
					ccipChangesetSolana.SetPoolConfig{
						ChainSelector:   solChain,
						TokenPubKey:     tokenAddress.String(),
						WritableIndexes: []uint8{3, 4, 7},
						MCMSSolana:      mcmsConfig,
					},
				),
			)
			require.NoError(t, err)
			tokenAdminRegistry := solCommon.TokenAdminRegistry{}
			tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)

			err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistry)
			require.NoError(t, err)
			require.Equal(t, newAdmin, tokenAdminRegistry.Administrator)
			require.Equal(t, lookupTablePubKey, tokenAdminRegistry.LookupTable)
		})
	}
}
