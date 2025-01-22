package changeset_test

import (
	"context"
	"encoding/binary"
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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
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
	testhelpers.SavePreloadedSolAddresses(t, e, solChain1)

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
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]changeset.ChainContractParams{
					solChain1: {
						FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
						OffRampParams:   changeset.DefaultOffRampParams(),
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddTokenPool),
			Config: changeset.AddTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "LinkToken",
				TokenProgramName: "spl-token-2022",
				PoolType:         "LockAndRelease",
				RampAuthority:    e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				Authority:        e.SolChains[solChain1].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetupTokenPoolForRemoteChain),
			Config: changeset.SetupTokenPoolForRemoteChainConfig{
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
	})
	require.NoError(t, err)

	// solana test
	tokenPubKey, err := deployment.FindTokenAddress(e, solChain1, "LinkToken")
	require.NoError(t, err)

	// pool stuff
	poolConfig, err := solTokenUtil.TokenPoolConfigAddress(tokenPubKey)
	require.NoError(t, err)
	poolSigner, err := solTokenUtil.TokenPoolSignerAddress(tokenPubKey)
	require.NoError(t, err)
	var configAccount token_pool.Config
	require.NoError(t, solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, poolConfig, solRpc.CommitmentConfirmed, &configAccount))
	poolTokenAccount, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenPubKey, poolSigner)
	require.Equal(t, poolTokenAccount, configAccount.PoolTokenAccount)
}

func TestBilling(t *testing.T) {
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
	testhelpers.SavePreloadedSolAddresses(t, e, solChain1)

	// Any nonzero timestamp is valid (for now)
	validTimestamp := int64(100)
	value := [28]uint8{}
	bigNum, ok := new(big.Int).SetString("19816680000000000000", 10)
	require.True(t, ok)
	bigNum.FillBytes(value[:])

	e, err = commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployLinkToken),
			Config:    []uint64{solChain1},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]changeset.ChainContractParams{
					solChain1: {
						FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
						OffRampParams:   changeset.DefaultOffRampParams(),
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddBillingToken),
			Config: changeset.AddBillingTokenPoolConfig{
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
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddBillingToken),
			Config: changeset.AddBillingTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "",
				TokenProgramName: "spl-token",
				TokenPubKey:      solana.SolMint.String(),
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
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddBillingTokenForRemoteChain),
			Config: changeset.BillingTokenForRemoteChainConfig{
				ChainSelector:       solChain1,
				RemoteChainSelector: homeChainSel,
				TokenName:           "LinkToken",
				TokenProgramName:    "spl-token-2022",
				Config:              ccip_router.TokenBilling{},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddBillingTokenForRemoteChain),
			Config: changeset.BillingTokenForRemoteChainConfig{
				ChainSelector:       solChain1,
				RemoteChainSelector: homeChainSel,
				TokenName:           "",
				TokenProgramName:    "spl-token",
				TokenPubKey:         solana.SolMint.String(),
				Config:              ccip_router.TokenBilling{},
			},
		},
	})
	require.NoError(t, err)

	// solana test
	tokenPubKey, err := deployment.FindTokenAddress(e, solChain1, "LinkToken")
	require.NoError(t, err)

	state, _ := changeset.LoadOnchainStateSolana(e)
	chainState := state.SolChains[solChain1]
	linkTokenBillingPDA, _, _ := solana.FindProgramAddress([][]byte{solTestConfig.BillingTokenConfigPrefix, tokenPubKey.Bytes()}, chainState.Router)
	var linkTokenConfigAccountPDA ccip_router.BillingTokenConfigWrapper
	aerr := solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, linkTokenBillingPDA, solRpc.CommitmentConfirmed, &linkTokenConfigAccountPDA)
	require.NoError(t, aerr)

	solTokenBillingPDA, _, _ := solana.FindProgramAddress([][]byte{solTestConfig.BillingTokenConfigPrefix, solana.SolMint.Bytes()}, chainState.Router)
	var solTokenConfigAccountPDA ccip_router.BillingTokenConfigWrapper
	aerr = solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, solTokenBillingPDA, solRpc.CommitmentConfirmed, &solTokenConfigAccountPDA)
	require.NoError(t, aerr)

	linkTokenRemoteBillingPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("ccip_tokenpool_billing"), binary.LittleEndian.AppendUint64([]byte{}, homeChainSel), tokenPubKey.Bytes()}, chainState.Router)
	var linkTokenRemoteConfigAccountPDA ccip_router.PerChainPerTokenConfig
	aerr = solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, linkTokenRemoteBillingPDA, solRpc.CommitmentConfirmed, &linkTokenRemoteConfigAccountPDA)
	require.NoError(t, aerr)

	solTokenRemoteBillingPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("ccip_tokenpool_billing"), binary.LittleEndian.AppendUint64([]byte{}, homeChainSel), solana.SolMint.Bytes()}, chainState.Router)
	var solTokenRemoteConfigAccountPDA ccip_router.PerChainPerTokenConfig
	aerr = solCommonUtil.GetAccountDataBorshInto(context.Background(), e.SolChains[solChain1].Client, solTokenRemoteBillingPDA, solRpc.CommitmentConfirmed, &solTokenRemoteConfigAccountPDA)
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
	testhelpers.SavePreloadedSolAddresses(t, e, solChain1)
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
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChainChangeset),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  testhelpers.NewTestRMNStaticConfig(),
				RMNDynamicConfig: testhelpers.NewTestRMNDynamicConfig(),
				NodeOperators:    testhelpers.NewTestNodeOperator(e.Chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": p2pIds,
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContractsChangeset),
			Config: changeset.DeployChainContractsConfig{
				HomeChainSelector: homeChainSel,
				ContractParamsPerChain: map[uint64]changeset.ChainContractParams{
					solChain1: {
						FeeQuoterParams: changeset.DefaultFeeQuoterParams(),
						OffRampParams:   changeset.DefaultOffRampParams(),
					},
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.AddTokenPool),
			Config: changeset.AddTokenPoolConfig{
				ChainSelector:    solChain1,
				TokenName:        "LinkToken",
				TokenProgramName: "spl-token-2022",
				PoolType:         "LockAndRelease",
				RampAuthority:    e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				Authority:        e.SolChains[solChain1].DeployerKey.PublicKey().String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.SetupTokenPoolForRemoteChain),
			Config: changeset.SetupTokenPoolForRemoteChainConfig{
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
			Changeset: commonchangeset.WrapChangeSet(changeset.RegisterTokenAdminRegistry),
			Config: changeset.RegisterTokenAdminRegistryConfig{
				ChainSelector:       solChain1,
				TokenName:           "LinkToken",
				TokenPoolAdmin:      tokenAdmin1.PublicKey().String(),
				AuthorityPrivateKey: e.SolChains[solChain1].DeployerKey.String(),
				RegisterType:        changeset.ViaOwnerInstruction,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.TransferAndAcceptAdminRoleTokenAdminRegistry),
			Config: changeset.TransferAndAcceptAdminRoleTokenAdminRegistryConfig{
				ChainSelector:               solChain1,
				TokenName:                   "LinkToken",
				TokenPoolAdminPrivateKey:    tokenAdmin1.String(),
				NewTokenPoolAdminPrivateKey: tokenAdmin2.String(),
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.UpdateTokenPool),
			Config: changeset.UpdateTokenPoolConfig{
				ChainSelector:       solChain1,
				TokenName:           "LinkToken",
				AuthorityPrivateKey: tokenAdmin2.String(),
				PoolLookupTable:     poolLookup.PublicKey().String(),
			},
		},
	})
	require.NoError(t, err)
}
