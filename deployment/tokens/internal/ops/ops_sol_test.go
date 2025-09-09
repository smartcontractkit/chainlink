package ops

import (
	"testing"

	solRpc "github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
)

// Note: This test does not perform an actual token deployment, as it would require a running
// Solana node via CTF containers. Such scenarios are already thoroughly covered by integration
// tests in the changeset and sequence suites. Until there is a way to simulate or mock a Solana
// client, this is the best we can do for unit tests.
//
// This test is limited to verifying operation error handling.
func Test_OpSolDeployLinkToken(t *testing.T) {
	t.Parallel()

	var (
		chainSelector = chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	)

	tests := []struct {
		name      string
		giveDeps  OpSolDeployLinkTokenDeps
		giveInput OpSolDeployLinkTokenInput
		wantErr   string
	}{
		{
			name:     "error: invalid chain selector",
			giveDeps: OpSolDeployLinkTokenDeps{},
			giveInput: OpSolDeployLinkTokenInput{
				ChainSelector: 1, // Invalid chain selector
			},
			wantErr: "unknown chain selector 1",
		},
		{
			name: "failed to create token",
			giveDeps: OpSolDeployLinkTokenDeps{
				Client: solRpc.New("bad-url"),
			},
			giveInput: OpSolDeployLinkTokenInput{
				ChainSelector: chainSelector,
			},
			wantErr: "failed to generate instructions for link token deployment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := operations.ExecuteOperation(
				optest.NewBundle(t), OpSolDeployLinkToken, tt.giveDeps, tt.giveInput,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
