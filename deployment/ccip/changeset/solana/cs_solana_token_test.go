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
		{
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
	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()

	e, err = changeset.ApplyChangesets(t, e, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(changeset_solana.CreateSolanaTokenATA),
			Config: changeset_solana.CreateSolanaTokenATAConfig{
				ChainSelector: solChain1,
				TokenPubkey:   tokenAddress,
				TokenProgram:  deployment.SPL2022Tokens,
				ATAList:       []string{testUserPubKey.String()},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.MintSolanaToken),
			Config: changeset_solana.MintSolanaTokenConfig{
				ChainSelector: solChain1,
				TokenPubkey:   tokenAddress,
				TokenProgram:  deployment.SPL2022Tokens,
				AmountToAddress: map[string]uint64{
					testUserPubKey.String(): uint64(1000),
				},
			},
		},
	})
	require.NoError(t, err)

	ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, testUserPubKey)
	outDec, outVal, err := solTokenUtil.TokenBalance(context.Background(), e.SolChains[solChain1].Client, ata, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))
}

func TestDeployLinkToken(t *testing.T) {
	testhelpers.DeployLinkTokenTest(t, 1)
}
