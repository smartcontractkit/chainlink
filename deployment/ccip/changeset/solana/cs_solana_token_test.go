package solana_test

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
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
	solChain1 := e.AllChainSelectorsSolana()[0]
	e, err := commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{ // deployer creates token
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenProgramName: deployment.SPL2022Tokens,
				TokenDecimals:    9,
			},
		},
	})
	require.NoError(t, err)

	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	tokenAddress := state.SolChains[solChain1].SPL2022Tokens[0]
	deployerKey := e.SolChains[solChain1].DeployerKey.PublicKey()

	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()

	e, err = changeset.ApplyChangesets(t, e, nil, []changeset.ChangesetApplication{
		{ // deployer creates ATA for itself and testUser
			Changeset: changeset.WrapChangeSet(changeset_solana.CreateSolanaTokenATA),
			Config: changeset_solana.CreateSolanaTokenATAConfig{
				ChainSelector: solChain1,
				TokenPubkey:   tokenAddress,
				TokenProgram:  deployment.SPL2022Tokens,
				ATAList:       []string{deployerKey.String(), testUserPubKey.String()},
			},
		},
		{ // deployer mints token to itself and testUser
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.MintSolanaToken),
			Config: changeset_solana.MintSolanaTokenConfig{
				ChainSelector: solChain1,
				TokenPubkey:   tokenAddress,
				TokenProgram:  deployment.SPL2022Tokens,
				AmountToAddress: map[string]uint64{
					deployerKey.String():    uint64(1000),
					testUserPubKey.String(): uint64(1000),
				},
			},
		},
	})
	require.NoError(t, err)

	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.Token2022ProgramID,
		tokenAddress,
		e.SolChains[solChain1].DeployerKey.PublicKey(),
	)
	require.NoError(t, err)

	// test if minting was done correctly
	outDec, outVal, err := solTokenUtil.TokenBalance(context.Background(), e.SolChains[solChain1].Client, deployerATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	t.Logf("outDec: %d, outVal: %d", outDec, outVal)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))

	outDec, outVal, err = solTokenUtil.TokenBalance(context.Background(), e.SolChains[solChain1].Client, testUserATA, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	t.Logf("outDec: %d, outVal: %d", outDec, outVal)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))
}

func TestDeployLinkToken(t *testing.T) {
	testhelpers.DeployLinkTokenTest(t, 1)
}
