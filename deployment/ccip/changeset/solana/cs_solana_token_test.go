package solana_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeploySolanaToken(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		SolChains: 1,
	})
	solChain1 := e.AllChainSelectorsSolana()[0]
	e, err := commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.DeploySolanaToken),
			Config: &changeset_solana.DeploySolanaTokenConfig{
				ChainSelector:    solChain1,
				TokenName:        "spl-token-2022",
				TokenProgramName: "spl-token-2022",
				ATAList: []string{
					e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset_solana.MintSolanaToken),
			Config: &changeset_solana.MintSolanaTokenConfig{
				ChainSelector: solChain1,
				TokenName:     "spl-token-2022",
				TokenProgram:  "spl-token-2022",
				Amount:        uint64(1000),
				ToAddressList: []string{
					e.SolChains[solChain1].DeployerKey.PublicKey().String(),
				},
			},
		},
	})
	require.NoError(t, err)

	// solana test
	tokenAddress, err := deployment.FindTokenAddress(e, solChain1, "spl-token-2022")
	require.NoError(t, err)
	toAddressBase58 := solana.MustPublicKeyFromBase58(e.SolChains[solChain1].DeployerKey.PublicKey().String())
	ata, _, _ := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, tokenAddress, toAddressBase58)
	outDec, outVal, err := solTokenUtil.TokenBalance(context.Background(), e.SolChains[solChain1].Client, ata, solRpc.CommitmentConfirmed)
	require.NoError(t, err)
	require.Equal(t, int(1000), outVal)
	require.Equal(t, 9, int(outDec))
}

func TestDeployLinkToken(t *testing.T) {
	testhelpers.DeployLinkTokenTest(t, 1)
}
