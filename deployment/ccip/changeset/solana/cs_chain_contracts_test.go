package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
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
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true, true, true, true)
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
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
			},
		),
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToSolana),
			ccipChangesetSolana.AddRemoteChainToSolanaConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.RemoteChainConfigSolana{
					evmChain: {
						EnabledAsSource: true,
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
						},
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
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddRemoteChainToSolana),
			ccipChangesetSolana.AddRemoteChainToSolanaConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]ccipChangesetSolana.RemoteChainConfigSolana{
					evmChain: {
						EnabledAsSource: true,
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: false,
						},
						FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
							IsEnabled:                   true,
							DefaultTxGasLimit:           30000,
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

func TestBilling(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Msg  string
		Mcms bool
	}{
		{
			Msg:  "TestBilling with mcms",
			Mcms: true,
		},
		{
			Msg:  "TestBilling without mcms",
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
			if test.Mcms {
				_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true, true, true, true)
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
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddBillingTokenForRemoteChain),
					ccipChangesetSolana.BillingTokenForRemoteChainConfig{
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
			},
			)
			require.NoError(t, err)
			err = e.SolChains[solChain].GetAccountDataBorshInto(e.GetContext(), billingConfigPDA, &token0ConfigAccount)
			require.NoError(t, err)
			require.Equal(t, uint64(200), token0ConfigAccount.Config.PremiumMultiplierWeiPerEth)
		})
	}

}

func TestTokenAdminRegistry(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	solChain := tenv.Env.AllChainSelectorsSolana()[0]
	e, tokenAddress, err := deployToken(t, tenv.Env, solChain)
	require.NoError(t, err)
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	linkTokenAddress := state.SolChains[solChain].LinkToken

	tokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			// register token admin registry for tokenAddress via admin instruction
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
			ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             tokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
			},
		),
		commonchangeset.Configure(
			// register token admin registry for linkToken via owner instruction
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
			ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             linkTokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            ccipChangesetSolana.ViaOwnerInstruction,
			},
		),
	)
	require.NoError(t, err)

	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)
	var tokenAdminRegistryAccount solRouter.TokenAdminRegistry
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.Administrator)
	// pending administrator should be the proposed admin key
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)

	linkTokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(linkTokenAddress, state.SolChains[solChain].Router)
	var linkTokenAdminRegistryAccount solRouter.TokenAdminRegistry
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, linkTokenAdminRegistryPDA, &linkTokenAdminRegistryAccount)
	require.NoError(t, err)
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), linkTokenAdminRegistryAccount.PendingAdministrator)

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			// accept admin role for tokenAddress
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
			ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:              solChain,
				TokenPubKey:                tokenAddress.String(),
				NewRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		),
	)
	require.NoError(t, err)
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	// confirm that the administrator is the deployer key
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistryAccount.Administrator)
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.PendingAdministrator)

	newTokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()
	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			// transfer admin role for tokenAddress
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.TransferAdminRoleTokenAdminRegistry),
			ccipChangesetSolana.TransferAdminRoleTokenAdminRegistryConfig{
				ChainSelector:                  solChain,
				TokenPubKey:                    tokenAddress.String(),
				NewRegistryAdminPublicKey:      newTokenAdminRegistryAdminPrivKey.PublicKey().String(),
				CurrentRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		),
	)
	require.NoError(t, err)
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	require.Equal(t, newTokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)
}

func TestPoolLookupTable(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

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

	tokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()

	e, err = commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			// register token admin registry for linkToken via owner instruction
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.RegisterTokenAdminRegistry),
			ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             tokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
			},
		),
		commonchangeset.Configure(
			// accept admin role for tokenAddress
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistry),
			ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:              solChain,
				TokenPubKey:                tokenAddress.String(),
				NewRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		),
		commonchangeset.Configure(
			// set pool -> this updates tokenAdminRegistryPDA, hence above changeset is required
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetPool),
			ccipChangesetSolana.SetPoolConfig{
				ChainSelector:                     solChain,
				TokenPubKey:                       tokenAddress.String(),
				TokenAdminRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
				WritableIndexes:                   []uint8{3, 4, 7},
			},
		),
	)
	require.NoError(t, err)
	tokenAdminRegistry := solRouter.TokenAdminRegistry{}
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)

	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistry)
	require.NoError(t, err)
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistry.Administrator)
	require.Equal(t, lookupTablePubKey, tokenAdminRegistry.LookupTable)
}
