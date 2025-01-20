package changeset

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solTestConfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/token_pool"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestAddTokenPool(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChain1 := e.AllChainSelectorsSolana()[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	SavePreloadedSolAddresses(t, e, solChain1)
	// tokenAdmin1, err := solana.NewRandomPrivateKey()
	// require.NoError(t, err)
	// tokenAdmin2, err := solana.NewRandomPrivateKey()
	// require.NoError(t, err)

	// Any nonzero timestamp is valid (for now)
	validTimestamp := int64(100)
	value := [28]uint8{}
	bigNum, ok := new(big.Int).SetString("19816680000000000000", 10)
	require.True(t, ok)
	bigNum.FillBytes(value[:])

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		// I CANNOT LOAD STATE IF I DEPLOY a random token, because load token expects to understand every address ?
		// {
		// 	Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeploySolanaToken),
		// 	Config: &commonchangeset.DeploySolanaTokenConfig{
		// 		ChainSelector:    solChain1,
		// 		TokenName:        "spl-token-2022",
		// 		TokenProgramName: "spl-token-2022",
		// 		ATAList: []string{
		// 			e.SolChains[solChain1].DeployerKey.PublicKey().String(),
		// 		},
		// 	},
		// },
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    []uint64{solChain1},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployHomeChain),
			Config: DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    []uint64{solChain1},
				HomeChainSelector: homeChainSel,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(AddTokenPool),
			Config: AddTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "LinkToken",
				TokenProgramName: "spl-token-2022",
				PoolType:         "LockAndRelease",
				RampAuthority:    e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				Authority:        e.SolChains[solChain1].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(SetupTokenPoolForChain),
			Config: SetupTokenPoolForChainConfig{
				ChainSelector:       solChain1,
				RemoteChainSelector: homeChainSel,
				TokenName:           "LinkToken",
				TokenProgramName:    "spl-token-2022",
				RemoteConfig: token_pool.RemoteConfig{
					PoolAddress:  []byte{1, 2, 3},
					TokenAddress: []byte{4, 5, 6},
					Decimals:     9,
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
		{
			Changeset: commonchangeset.WrapChangeSet(AddBillingTokenPool),
			Config: AddBillingTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "LinkToken",
				TokenProgramName: "spl-token-2022",
				Config: ccip_router.BillingTokenConfig{
					Enabled: true,
					// Mint:    token2022.mint,
					UsdPerToken: ccip_router.TimestampedPackedU224{
						Value:     value,
						Timestamp: validTimestamp,
					},
					PremiumMultiplierWeiPerEth: 11000000,
				}},
		},
		// {
		// 	Changeset: commonchangeset.WrapChangeSet(RegisterTokenAdminRegistry),
		// 	Config: RegisterTokenAdminRegistryConfig{
		// 		ChainSelector:       solChain1,
		// 		TokenName:           "LinkToken",
		// 		TokenPoolAdmin:      tokenAdmin1.PublicKey().String(),
		// 		AuthorityPrivateKey: e.SolChains[solChain1].DeployerKey.String(),
		// 		RegisterType:        ViaGetCcipAdminInstruction,
		// 	},
		// },
		// {
		// 	Changeset: commonchangeset.WrapChangeSet(TransferAndAcceptAdminRoleTokenAdminRegistry),
		// 	Config: TransferAndAcceptAdminRoleTokenAdminRegistryConfig{
		// 		ChainSelector:               solChain1,
		// 		TokenName:                   "LinkToken",
		// 		TokenPoolAdminPrivateKey:    tokenAdmin1.String(),
		// 		NewTokenPoolAdminPrivateKey: tokenAdmin2.String(),
		// 	},
		// },
	})
	require.NoError(t, err)

	// solana test
	tokenPubKey, err := deployment.FindTokenAddress(e, solChain1, "LinkToken")
	require.NoError(t, err)
	poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey)
	require.NoError(t, err)
	poolSigner, err := solTokenUtil.TokenPoolSignerAddress(tokenPubKey)
	require.NoError(t, err)
	var configAccount token_pool.Config
	require.NoError(t, solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, poolConfig, solRpc.CommitmentConfirmed, &configAccount))
	poolTokenAccount, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenPubKey, poolSigner)
	require.Equal(t, poolTokenAccount, configAccount.PoolTokenAccount)

	state, _ := LoadOnchainStateSolana(e)
	chainState := state.SolChains[solChain1]
	tokenBillingPDA, _, _ := solana.FindProgramAddress([][]byte{solTestConfig.BillingTokenConfigPrefix, tokenPubKey.Bytes()}, chainState.SolCcipRouter)
	var token0ConfigAccount ccip_router.BillingTokenConfigWrapper
	aerr := solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, tokenBillingPDA, solRpc.CommitmentConfirmed, &token0ConfigAccount)
	require.NoError(t, aerr)
}

func TestTokenAdminRegistry(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		SolChains:  1,
		Nodes:      4,
	})
	evmSelectors := e.AllChainSelectors()
	homeChainSel := evmSelectors[0]
	solChain1 := e.AllChainSelectorsSolana()[0]
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	p2pIds := nodes.NonBootstraps().PeerIDs()
	SavePreloadedSolAddresses(t, e, solChain1)
	tokenAdmin1, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	tokenAdmin2, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	poolLookup, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    []uint64{solChain1},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployHomeChain),
			Config: DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  NewTestRMNStaticConfig(),
				RMNDynamicConfig: NewTestRMNDynamicConfig(),
				NodeOperators:    NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(DeployChainContracts),
			Config: DeployChainContractsConfig{
				ChainSelectors:    []uint64{solChain1},
				HomeChainSelector: homeChainSel,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(AddTokenPool),
			Config: AddTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "LinkToken",
				TokenProgramName: "spl-token-2022",
				PoolType:         "LockAndRelease",
				RampAuthority:    e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				Authority:        e.SolChains[solChain1].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(SetupTokenPoolForChain),
			Config: SetupTokenPoolForChainConfig{
				ChainSelector:       solChain1,
				RemoteChainSelector: homeChainSel,
				TokenName:           "LinkToken",
				TokenProgramName:    "spl-token-2022",
				RemoteConfig: token_pool.RemoteConfig{
					PoolAddress:  []byte{1, 2, 3},
					TokenAddress: []byte{4, 5, 6},
					Decimals:     9,
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
		{
			Changeset: commonchangeset.WrapChangeSet(RegisterTokenAdminRegistry),
			Config: RegisterTokenAdminRegistryConfig{
				ChainSelector:       solChain1,
				TokenName:           "LinkToken",
				TokenPoolAdmin:      tokenAdmin1.PublicKey().String(),
				AuthorityPrivateKey: e.SolChains[solChain1].DeployerKey.String(),
				RegisterType:        ViaOwnerInstruction,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(TransferAndAcceptAdminRoleTokenAdminRegistry),
			Config: TransferAndAcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:               solChain1,
				TokenName:                   "LinkToken",
				TokenPoolAdminPrivateKey:    tokenAdmin1.String(),
				NewTokenPoolAdminPrivateKey: tokenAdmin2.String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(UpdateTokenPool),
			Config: UpdateTokenPoolConfig{
				ChainSelector:       solChain1,
				TokenName:           "LinkToken",
				AuthorityPrivateKey: tokenAdmin2.String(),
				PoolLookupTable:     poolLookup.PublicKey().String(),
			},
		},
	})
	require.NoError(t, err)
}
