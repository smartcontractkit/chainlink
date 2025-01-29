package solana_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/stretchr/testify/require"

	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

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

	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	_, err = commonchangeset.ApplyChangesets(t, tenv.Env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateOnRampsDestsChangeset),
			Config: changeset.UpdateOnRampDestsConfig{
				UpdatesByChain: map[uint64]map[uint64]changeset.OnRampDestinationUpdate{
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
				UpdatesByChain: map[uint64]map[uint64]changeset_solana.RemoteChainConfigSolana{
					solChain: {
						evmChain: {
							EnabledAsSource:          true,
							RemoteChainOnRampAddress: state.Chains[evmChain].OnRamp.Address().String(),
							DestinationConfig: solRouter.DestChainConfig{
								IsEnabled:                   true,
								DefaultTxGasLimit:           200000,
								MaxPerMsgGasLimit:           3000000,
								MaxDataBytes:                30000,
								MaxNumberOfTokensPerMsg:     5,
								DefaultTokenDestGasOverhead: 5000,
								// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
								// TODO: do a similar test for other chain families
								ChainFamilySelector: [4]uint8{40, 18, 213, 44},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	state, err = changeset.LoadOnchainStateSolana(tenv.Env)
	require.NoError(t, err)

	var sourceChainStateAccount solRouter.SourceChain
	evmSourceChainStatePDA, _ := solState.FindSourceChainStatePDA(evmChain, state.SolChains[solChain].Router)
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, evmSourceChainStatePDA, &sourceChainStateAccount)
	require.NoError(t, err)
	require.Equal(t, uint64(1), sourceChainStateAccount.State.MinSeqNr)
	require.True(t, sourceChainStateAccount.Config.IsEnabled)

	var destChainStateAccount solRouter.DestChain
	evmDestChainStatePDA, _ := solState.FindDestChainStatePDA(evmChain, state.SolChains[solChain].Router)
	err = tenv.Env.SolChains[solChain].GetAccountDataBorshInto(ctx, evmDestChainStatePDA, &destChainStateAccount)
	require.NoError(t, err)
}

func TestDeployCCIPContracts(t *testing.T) {
	t.Parallel()
	testhelpers.DeployCCIPContractsTest(t, 1)
}

func TestAddTokenPool(t *testing.T) {
	t.Parallel()

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
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.AddTokenPool),
			Config: changeset_solana.TokenPoolConfig{
				ChainSelector:    solChain,
				TokenPubKey:      tokenAddress.String(),
				TokenProgramName: deployment.SPL2022Tokens,
				PoolType:         "LockAndRelease",
				// this works for testing, but if we really want some other authority we need to pass in a private key for signing purposes
				Authority: e.SolChains[solChain].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.SetupTokenPoolForRemoteChain),
			Config: changeset_solana.RemoteChainTokenPoolConfig{
				ChainSelector:       solChain,
				RemoteChainSelector: evmChain,
				TokenPubKey:         tokenAddress.String(),
				RemoteConfig: token_pool.RemoteConfig{
					PoolAddresses: []token_pool.RemoteAddress{{Address: []byte{1, 2, 3}}},
					TokenAddress:  token_pool.RemoteAddress{Address: []byte{4, 5, 6}},
					Decimals:      9,
				},
				InboundRateLimit: token_pool.RateLimitConfig{
					Enabled:  true,
					Capacity: uint64(1000),
					Rate:     1,
				},
				OutboundRateLimit: token_pool.RateLimitConfig{
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
	var configAccount token_pool.Config
	err = e.SolChains[solChain].GetAccountDataBorshInto(context.Background(), poolConfigPDA, &configAccount)
	t.Logf("configAccount: %+v", configAccount)
	require.NoError(t, err)
	require.Equal(t, configAccount.PoolType, token_pool.LockAndRelease_PoolType)
	require.Equal(t, configAccount.Mint, tokenAddress)
	// try minting after this and see if the pool or the deployer key is the authority

	// test SetupTokenPoolForRemoteChain results
	remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(evmChain, tokenAddress, state.SolChains[solChain].TokenPool)
	var remoteChainConfigAccount token_pool.RemoteConfig
	err = e.SolChains[solChain].GetAccountDataBorshInto(context.Background(), remoteChainConfigPDA, &remoteChainConfigAccount)
	t.Logf("remoteChainConfigAccount: %+v", remoteChainConfigAccount)
	t.Logf("err: %+v", err)
	// TODO: this is returning empty config
	// require.NoError(t, err)
	// require.Equal(t, remoteChainConfigAccount.Decimals, uint8(9))
}

func TestBilling(t *testing.T) {
	t.Parallel()

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
			Config: changeset_solana.BillingTokenPoolConfig{
				ChainSelector:    solChain,
				TokenPubKey:      tokenAddress.String(),
				TokenProgramName: deployment.SPL2022Tokens,
				Config: solRouter.BillingTokenConfig{
					Enabled: true,
					Mint:    tokenAddress,
					UsdPerToken: solRouter.TimestampedPackedU224{
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
				Config: solRouter.TokenBilling{
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

	billingConfigPDA, _, _ := solState.FindFeeBillingTokenConfigPDA(tokenAddress, state.SolChains[solChain].Router)
	var token0ConfigAccount solRouter.BillingTokenConfigWrapper
	err = e.SolChains[solChain].GetAccountDataBorshInto(context.Background(), billingConfigPDA, &token0ConfigAccount)
	require.NoError(t, err)
	require.Equal(t, token0ConfigAccount.Config.Enabled, true)
	require.Equal(t, token0ConfigAccount.Config.Mint, tokenAddress)

	remoteBillingPDA, _, _ := solState.FindCcipTokenpoolBillingPDA(evmChain, tokenAddress, state.SolChains[solChain].Router)
	var remoteBillingAccount solRouter.PerChainPerTokenConfig
	err = e.SolChains[solChain].GetAccountDataBorshInto(context.Background(), remoteBillingPDA, &remoteBillingAccount)
	require.NoError(t, err)
	require.Equal(t, remoteBillingAccount.Mint, tokenAddress)
	require.Equal(t, remoteBillingAccount.Billing.MinFeeUsdcents, uint32(800))
}
