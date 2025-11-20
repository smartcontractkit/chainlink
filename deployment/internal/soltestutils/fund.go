package soltestutils

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

// FundSignerPDAs funds the timelock signer and MCMS signer PDAs with 1 SOL for testing
func FundSignerPDAs(
	t *testing.T, chain cldfsolana.Chain, mcmsState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()

	timelockSignerPDA := state.GetTimelockSignerPDA(mcmsState.TimelockProgram, mcmsState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(mcmsState.McmProgram, mcmsState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	err := solutils.FundAccounts(t.Context(), chain.Client, signerPDAs, 1)
	require.NoError(t, err)
}
