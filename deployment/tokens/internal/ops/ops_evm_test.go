package ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainevm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	chainevmprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
)

func Test_OpEVMDeployLinkToken(t *testing.T) {
	t.Parallel()

	var (
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

			prov := chainevmprovider.NewSimChainProvider(
				t, chainSelector, chainevmprovider.SimChainProviderConfig{},
			)
			blockchain, err := prov.Initialize(t.Context())
			require.NoError(t, err)

			chain, ok := blockchain.(chainevm.Chain)
			require.True(t, ok)

			var (
				auth = chain.DeployerKey
				deps = OpEVMDeployLinkTokenDeps{
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
