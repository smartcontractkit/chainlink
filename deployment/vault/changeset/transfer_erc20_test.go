package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const (
	transferERC20Payee    = "0x742d35cc64ca395db82e2e3e8fa8bc6d1b7c0832"
	transferERC20Token    = "0x1234567890123456789012345678901234567890"
	transferERC20ZeroAddr = "0x0000000000000000000000000000000000000000"
)

func TestValidateTransferERC20Config(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(), environment.WithEVMSimulated(t, []uint64{selector}))
	require.NoError(t, err)

	validMCMSConfig := &proposalutils.TimelockConfig{MinDelay: 0}

	tests := []struct {
		name      string
		config    types.TransferERC20Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty transfers",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers:          []types.ERC20Transfer{},
				MCMSConfig:         validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "transfers must not be empty",
		},
		{
			name: "zero payee",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20ZeroAddr, Token: transferERC20Token, Amount: big.NewInt(100)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "payee cannot be zero address",
		},
		{
			name: "zero token",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20ZeroAddr, Amount: big.NewInt(100)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "token address cannot be zero",
		},
		{
			name: "zero amount",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: big.NewInt(0)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "amount must be positive",
		},
		{
			name: "nil amount",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: nil},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "amount must be positive",
		},
		{
			name: "missing MCMSConfig",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: big.NewInt(100)},
				},
				MCMSConfig: nil,
			},
			wantError: true,
			errorMsg:  "MCMSConfig is required",
		},
		{
			name: "invalid chain selector",
			config: types.TransferERC20Config{
				ChainSelector:      999999,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: big.NewInt(100)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "invalid chain selector",
		},
		{
			name: "chain not in environment",
			config: types.TransferERC20Config{
				ChainSelector:      chainselectors.TEST_90000002.Selector,
				TimelockIdentifier: "",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: big.NewInt(100)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "not found in environment",
		},
		{
			name: "timelock not found for qualifier",
			config: types.TransferERC20Config{
				ChainSelector:      selector,
				TimelockIdentifier: "nonexistent_timelock",
				Transfers: []types.ERC20Transfer{
					{Payee: transferERC20Payee, Token: transferERC20Token, Amount: big.NewInt(100)},
				},
				MCMSConfig: validMCMSConfig,
			},
			wantError: true,
			errorMsg:  "timelock not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateTransferERC20Config(*env, tt.config)
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

func TestEncodeERC20TransferCalldata(t *testing.T) {
	t.Parallel()

	payee := common.HexToAddress(transferERC20Payee)
	amount := big.NewInt(1000)

	calldata, err := encodeERC20TransferCalldata(types.ERC20Transfer{
		Payee:  payee.Hex(),
		Token:  transferERC20Token,
		Amount: amount,
	})
	require.NoError(t, err)

	require.Len(t, calldata, 4+32+32) // selector + address + uint256
	require.Equal(t, erc20TransferSelector, calldata[:4])
}
