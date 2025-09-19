package ccip_attestation_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gagliardetto/solana-go"
	signerRegistry "github.com/smartcontractkit/ccip-base/chains/solana/go_bindings"
	"go.uber.org/zap/zapcore"

	ccip_attestation "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/ccip_attestation"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestSignerRegistryInitialization(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	// We download the artifacts before spinning out the environment so the program can be preloaded.
	_, currentFile, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	targetPath := filepath.Join(rootDir, "ccip/changeset/internal", "solana_contracts")
	// Replace with the desired run to test
	err := ccip_attestation.DownloadReleaseArtifactsFromGithubWorkflowRun(context.Background(), "17459947443", "3925592769", targetPath)
	require.NoError(t, err)

	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		SolChains: 1,
	})
	solChain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]

	owner, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	e, err = commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			// deployer creates token
			cldf.CreateLegacyChangeSet(ccip_attestation.InitializeBaseSignerRegistryContractChangeset),
			ccip_attestation.InitalizeBaseSignerRegistryContractConfig{
				ChainSelector: solChain1,
				Owner:         owner.PublicKey(),
			},
		),
	)
	require.NoError(t, err)

	programID := solana.MustPublicKeyFromBase58(memory.SolanaProgramIDs["ccip_signer_registry"])
	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, programID)

	var configAccount signerRegistry.Config
	chain := e.BlockChains.SolanaChains()[solChain1]
	err = chain.GetAccountDataBorshInto(context.Background(), configPda, &configAccount)
	require.NoError(t, err)

	require.Equal(t, owner.PublicKey(), configAccount.Owner, "owner not set correctly in config PDA")
}
