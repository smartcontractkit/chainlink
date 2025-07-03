package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func Test_OpEVMDeployLinkToken(t *testing.T) {
	t.Parallel()

	var (
		chainID       uint64 = 11155111
		chainSelector uint64 = 16015286601757825753
	)

	tests := []struct {
		name    string
		give    OpEVMDeployLinkTokenInput
		want    OpEvmDeployLinkTokenOutput
		wantErr string
	}{
		{
			name: "deploys LinkToken on EVM chain",
			give: OpEVMDeployLinkTokenInput{
				ChainSelector: chainSelector,
			},
			want: OpEvmDeployLinkTokenOutput{
				Type:    LinkTokenTypeAndVersion1.Type.String(),
				Version: LinkTokenTypeAndVersion1.Version.String(),
			},
		},
		{
			name: "error: invalid chain selector",
			give: OpEVMDeployLinkTokenInput{
				ChainSelector: 1, // Invalid chain selector
			},
			wantErr: "unknown chain selector 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				chains = cldf_chain.NewBlockChainsFromSlice(
					memory.NewMemoryChainsEVMWithChainIDs(t, []uint64{chainID}, 1),
				).EVMChains()
				chain = chains[chainSelector]
				auth  = chain.DeployerKey
				deps  = OpEVMDeployLinkTokenDeps{
					Auth:        auth,
					Backend:     chain.Client,
					ConfirmFunc: chain.Confirm,
				}
			)

			got, err := operations.ExecuteOperation(
				optest.NewBundle(t), OpEVMDeployLinkToken, deps, tt.give,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, got.Output.Address.String())
				assert.Equal(t, tt.want.Type, got.Output.Type)
				assert.Equal(t, tt.want.Version, got.Output.Version)
			}
		})
	}
}
