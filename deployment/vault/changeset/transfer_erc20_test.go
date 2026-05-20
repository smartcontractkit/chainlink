package changeset

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestTransferERC20Validation(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name      string
		config    types.TransferERC20Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty transfers_by_chain",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{},
				MCMSConfig:       &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "transfers_by_chain must not be empty",
		},
		{
			name: "missing MCMS config",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(1),
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "MCMSConfig is required",
		},
		{
			name: "payee not whitelisted",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(1),
						},
					},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "is not whitelisted",
		},
		{
			name: "zero amount",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(0),
						},
					},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "amount must be positive",
		},
		{
			name: "chain with no transfers",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "has no transfers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateTransferERC20Config(env.GetContext(), *env, tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
