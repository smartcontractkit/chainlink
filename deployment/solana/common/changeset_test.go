package deployment_solana_common

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestSolanaTokenChangeset(t *testing.T) {
	t.Skip("This test is flaky and should be fixed")

	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		SolChains: 1,
	})
	solChain1 := e.AllChainSelectorsSolana()[0]

	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()

	changesetInput := DeployAndMintSolanaTokenConfig{
		ChainSelector:    solChain1,
		TokenProgramName: deployment.SPL2022Tokens,
		TokenDecimals:    9,
		MintAmount:       1000,
		ATAList:          []string{testUserPubKey.String()},
		AmountToAddress: map[string]uint64{
			testUserPubKey.String(): uint64(1000),
		},
	}

	ret, err := DeployAndMintSolanaToken(e, changesetInput)
	require.NoError(t, err)

	// Check the report
	reports := e.OpEnv.Reporter.GetReports()
	require.Len(t, reports, 3)

	// Check the output
	require.Len(t, ret.Reports, 3)

}
