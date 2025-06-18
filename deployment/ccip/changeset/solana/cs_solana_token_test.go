package solana_test

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSolanaTokenOps(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		SolChains: 1,
	})
	solChain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      "TEST_TOKEN",
			},
		),
	)
	require.NoError(t, err)

	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPLTokens,
				MintPrivateKey:   privKey,
				TokenDecimals:    9,
				TokenSymbol:      "SPL_TEST_TOKEN",
			},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(solChain1) //nolint:staticcheck // addressbook still valid
	require.NoError(t, err)
	tokenAddress := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPL2022Tokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("TEST_TOKEN"),
		},
		addresses,
	)

	deployerKey := e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey()

	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()

	e, err = commonchangeset.Apply(t, e, commonchangeset.Configure(
		// deployer creates ATA for itself and testUser
		cldf.CreateLegacyChangeSet(changeset_solana.CreateSolanaTokenATA),
		changeset_solana.CreateSolanaTokenATAConfig{
			ChainSelector: solChain1,
			TokenPubkey:   tokenAddress,
			ATAList:       []string{deployerKey.String(), testUserPubKey.String()},
		},
	), commonchangeset.Configure(
		// deployer mints token to itself and testUser
		cldf.CreateLegacyChangeSet(changeset_solana.MintSolanaToken),
		changeset_solana.MintSolanaTokenConfig{
			ChainSelector: solChain1,
			TokenPubkey:   tokenAddress.String(),
			AmountToAddress: map[string]uint64{
				deployerKey.String():    uint64(1000),
				testUserPubKey.String(): uint64(1000),
			},
		},
	))
	require.NoError(t, err)

	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.Token2022ProgramID,
		tokenAddress,
		e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey(),
	)
	require.NoError(t, err)

	// test if minting was done correctly
	outDec, outVal, err := solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, deployerATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, testUserATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	// now lets do it altogether
	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
			changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: shared.SPLTokens,
				TokenDecimals:    9,
				TokenSymbol:      "TEST_TOKEN_2",
				ATAList:          []string{deployerKey.String(), testUserPubKey.String()},
				MintAmountToAddress: map[string]uint64{
					deployerKey.String():    uint64(1000),
					testUserPubKey.String(): uint64(1000),
				},
			},
		),
	)
	require.NoError(t, err)
	addresses, err = e.ExistingAddresses.AddressesForChain(solChain1) //nolint:staticcheck // addressbook still valid
	require.NoError(t, err)
	tokenAddress2 := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPLTokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("TEST_TOKEN_2"),
		},
		addresses,
	)
	testUserATA2, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, tokenAddress2, testUserPubKey)
	require.NoError(t, err)
	deployerATA2, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.TokenProgramID,
		tokenAddress2,
		e.BlockChains.SolanaChains()[solChain1].DeployerKey.PublicKey(),
	)
	require.NoError(t, err)
	// test if minting was done correctly
	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, deployerATA2, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.BlockChains.SolanaChains()[solChain1].Client, testUserATA2, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))
}

func TestDeployLinkToken(t *testing.T) {
	commonchangeset.DeployLinkTokenTest(t, memory.MemoryEnvironmentConfig{
		Chains:    1,
		SolChains: 1,
	})
}
