package solana_test

import (
	"math/big"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestAddRemoteChain(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	// Default env just has 2 chains with all contracts
	// deployed but no lanes.
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	_, err := ccipChangeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	tenv.Env, err = commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(ccipChangeset.UpdateOnRampsDestsChangeset),
			Config: ccipChangeset.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]ccipChangeset.OnRampDestinationUpdate{
					evmChain: {
						solChain: {
							IsEnabled:        true,
							TestRouter:       false,
							AllowListEnabled: false,
						},
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddRemoteChainToSolana),
			Config: changeset_solana.AddRemoteChainToSolanaConfig{
				ChainSelector: solChain,
				UpdatesByChain: map[uint64]changeset_solana.RemoteChainConfigSolana{
					evmChain: {
						EnabledAsSource:         true,
						RouterDestinationConfig: solRouter.DestChainConfig{},
						FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
							IsEnabled:                   true,
							DefaultTxGasLimit:           200000,
							MaxPerMsgGasLimit:           3000000,
							MaxDataBytes:                30000,
							MaxNumberOfTokensPerMsg:     5,
							DefaultTokenDestGasOverhead: 5000,
							// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
							// TODO: do a similar test for other chain families
							// https://smartcontract-it.atlassian.net/browse/INTAUTO-438
							ChainFamilySelector: [4]uint8{40, 18, 213, 44},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	var destChainStateAccount solRouter.DestChain
	evmDestChainStatePDA := state.SolChains[solChain].DestChainStatePDAs[evmChain]
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)

	var destChainFqAccount solFeeQuoter.DestChain
	fqEvmDestChainPDA, _, _ := solState.FindFqDestChainPDA(evmChain, state.SolChains[solChain].FeeQuoter)
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, fqEvmDestChainPDA, &destChainFqAccount)
	require.NoError(t, err, "failed to get account info")
	require.Equal(t, solFeeQuoter.TimestampedPackedU224{}, destChainFqAccount.State.UsdPerUnitGas)
	require.True(t, destChainFqAccount.Config.IsEnabled)
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 1)
}

func TestAddTokenPool(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	e, err := commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain].SPL2022Tokens[0]

	// TODO: can test this with solana.SolMint as well (WSOL)
	// https://smartcontract-it.atlassian.net/browse/INTAUTO-440
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddTokenPool),
			Config: changeset_solana.TokenPoolConfig{
				ChainSelector:    solChain,
				TokenPubKey:      tokenAddress.String(),
				TokenProgramName: deployment.SPL2022Tokens,
				PoolType:         solTestTokenPool.LockAndRelease_PoolType,
				// this works for testing, but if we really want some other authority we need to pass in a private key for signing purposes
				Authority: e.SolChains[solChain].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.SetupTokenPoolForRemoteChain),
			Config: changeset_solana.RemoteChainTokenPoolConfig{
				SolChainSelector:    solChain,
				RemoteChainSelector: evmChain,
				SolTokenPubKey:      tokenAddress.String(),
				RemoteConfig: solTestTokenPool.RemoteConfig{
					// TODO:this can be potentially read from the state if we are given the token symbol
					PoolAddresses: []solTestTokenPool.RemoteAddress{{Address: []byte{1, 2, 3}}},
					TokenAddress:  solTestTokenPool.RemoteAddress{Address: []byte{4, 5, 6}},
					Decimals:      9,
				},
				InboundRateLimit: solTestTokenPool.RateLimitConfig{
					Enabled:  true,
					Capacity: uint64(1000),
					Rate:     1,
				},
				OutboundRateLimit: solTestTokenPool.RateLimitConfig{
					Enabled:  false,
					Capacity: 0,
					Rate:     0,
				},
			},
		},
	})
	require.NoError(t, err)

	// test AddTokenPool results
	poolConfigPDA, err := solTokenUtil.TokenPoolConfigAddress(tokenAddress, state.SolChains[solChain].TokenPool)
	require.NoError(t, err)
	var configAccount solTestTokenPool.State
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, poolConfigPDA, &configAccount)
	require.NoError(t, err)
	require.Equal(t, solTestTokenPool.LockAndRelease_PoolType, configAccount.PoolType)
	require.Equal(t, tokenAddress, configAccount.Config.Mint)
	// try minting after this and see if the pool or the deployer key is the authority

	// test SetupTokenPoolForRemoteChain results
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(evmChain, tokenAddress, state.SolChains[solChain].TokenPool)
	var remoteChainConfigAccount solTestTokenPool.ChainConfig
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, remoteChainConfigPDA, &remoteChainConfigAccount)
	require.NoError(t, err)
	require.Equal(t, uint8(9), remoteChainConfigAccount.Base.Remote.Decimals)
}

func TestBilling(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	e, err := commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain].SPL2022Tokens[0]
	validTimestamp := int64(100)
	value := [28]uint8{}
	bigNum, ok := new(big.Int).SetString("19816680000000000000", 10)
	require.True(t, ok)
	bigNum.FillBytes(value[:])
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddBillingToken),
			Config: changeset_solana.BillingTokenConfig{
				ChainSelector:    solChain,
				TokenPubKey:      tokenAddress.String(),
				TokenProgramName: deployment.SPL2022Tokens,
				Config: solFeeQuoter.BillingTokenConfig{
					Enabled: true,
					Mint:    tokenAddress,
					UsdPerToken: solFeeQuoter.TimestampedPackedU224{
						Timestamp: validTimestamp,
						Value:     value,
					},
					PremiumMultiplierWeiPerEth: 100,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddBillingTokenForRemoteChain),
			Config: changeset_solana.BillingTokenForRemoteChainConfig{
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
			},
		},
	})
	require.NoError(t, err)

	billingConfigPDA, _, _ := solState.FindFqBillingTokenConfigPDA(tokenAddress, state.SolChains[solChain].FeeQuoter)
	var token0ConfigAccount solFeeQuoter.BillingTokenConfigWrapper
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, billingConfigPDA, &token0ConfigAccount)
	require.NoError(t, err)
	require.True(t, token0ConfigAccount.Config.Enabled)
	require.Equal(t, tokenAddress, token0ConfigAccount.Config.Mint)

	remoteBillingPDA, _, _ := solState.FindFqPerChainPerTokenConfigPDA(evmChain, tokenAddress, state.SolChains[solChain].FeeQuoter)
	var remoteBillingAccount solFeeQuoter.PerChainPerTokenConfig
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, remoteBillingPDA, &remoteBillingAccount)
	require.NoError(t, err)
	require.Equal(t, tokenAddress, remoteBillingAccount.Mint)
	require.Equal(t, uint32(800), remoteBillingAccount.TokenTransferConfig.MinFeeUsdcents)
}

func TestTokenAdminRegistry(t *testing.T) {
	t.Parallel()
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	solChain := tenv.Env.AllChainSelectorsSolana()[0]

	e, err := commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain].SPL2022Tokens[0]
	tokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()

	// We have to do run the ViaOwnerInstruction testcase for linkToken as we already register a PDA for tokenAddress in the previous testcase
	linkTokenAddress := state.SolChains[solChain].LinkToken

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // register token admin registry for tokenAddress via admin instruction
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.RegisterTokenAdminRegistry),
			Config: changeset_solana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             tokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            changeset_solana.ViaGetCcipAdminInstruction,
			},
		},
		{ // register token admin registry for linkToken via owner instruction
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.RegisterTokenAdminRegistry),
			Config: changeset_solana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             linkTokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            changeset_solana.ViaOwnerInstruction,
			},
		},
	})
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

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // accept admin role for tokenAddress
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AcceptAdminRoleTokenAdminRegistry),
			Config: changeset_solana.AcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:              solChain,
				TokenPubKey:                tokenAddress.String(),
				NewRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		},
	})
	require.NoError(t, err)
	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	// confirm that the administrator is the deployer key
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistryAccount.Administrator)
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.PendingAdministrator)

	// TODO: transfer and accept is breaking
	newTokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // transfer admin role for tokenAddress
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.TransferAdminRoleTokenAdminRegistry),
			Config: changeset_solana.TransferAdminRoleTokenAdminRegistryConfig{
				ChainSelector:                  solChain,
				TokenPubKey:                    tokenAddress.String(),
				NewRegistryAdminPublicKey:      newTokenAdminRegistryAdminPrivKey.PublicKey().String(),
				CurrentRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		},
	})
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

	e, err := commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{ // deploy token
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain].SPL2022Tokens[0]

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // add token pool lookup table
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddTokenPoolLookupTable),
			Config: changeset_solana.TokenPoolLookupTableConfig{
				ChainSelector: solChain,
				TokenPubKey:   tokenAddress.String(),
				TokenProgram:  deployment.SPL2022Tokens,
			},
		},
	})
	require.NoError(t, err)

	state, err = ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	lookupTablePubKey := state.SolChains[solChain].TokenPoolLookupTable[tokenAddress]

	lookupTableEntries0, err := solCommonUtil.GetAddressLookupTable(ctx, e.SolChains[solChain].Client, lookupTablePubKey)
	require.NoError(t, err)
	require.Equal(t, lookupTablePubKey, lookupTableEntries0[0])
	require.Equal(t, tokenAddress, lookupTableEntries0[7])

	tokenAdminRegistryAdminPrivKey, _ := solana.NewRandomPrivateKey()

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // register token admin registry for linkToken via owner instruction
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.RegisterTokenAdminRegistry),
			Config: changeset_solana.RegisterTokenAdminRegistryConfig{
				ChainSelector:           solChain,
				TokenPubKey:             tokenAddress.String(),
				TokenAdminRegistryAdmin: tokenAdminRegistryAdminPrivKey.PublicKey().String(),
				RegisterType:            changeset_solana.ViaGetCcipAdminInstruction,
			},
		},
		{ // accept admin role for tokenAddress
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AcceptAdminRoleTokenAdminRegistry),
			Config: changeset_solana.AcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:              solChain,
				TokenPubKey:                tokenAddress.String(),
				NewRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
			},
		},
		{ // set pool -> this updates tokenAdminRegistryPDA, hence above changeset is required
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.SetPool),
			Config: changeset_solana.SetPoolConfig{
				ChainSelector:                     solChain,
				TokenPubKey:                       tokenAddress.String(),
				PoolLookupTable:                   lookupTablePubKey.String(),
				TokenAdminRegistryAdminPrivateKey: tokenAdminRegistryAdminPrivKey.String(),
				WritableIndexes:                   []uint8{3, 4, 7},
			},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistry := solRouter.TokenAdminRegistry{}
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)

	err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistry)
	require.NoError(t, err)
	require.Equal(t, tokenAdminRegistryAdminPrivKey.PublicKey(), tokenAdminRegistry.Administrator)
	require.Equal(t, lookupTablePubKey, tokenAdminRegistry.LookupTable)
}
