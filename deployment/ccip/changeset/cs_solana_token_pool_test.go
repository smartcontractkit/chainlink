package changeset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
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
}
