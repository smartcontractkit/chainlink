package core

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestUtils_GetTxBlock(t *testing.T) {
	tests := []struct {
		name         string
		txHash       common.Hash
		ethCallError error
		receipt      *types.Receipt
		status       uint64
	}{
		{
			name:   "success",
			txHash: common.HexToHash("0xc48fbf05edaf18f6aaa7de24de28528546b874bb03728d624ca407b8fed582a3"),
			receipt: &types.Receipt{
				Status:      1,
				BlockNumber: big.NewInt(2000),
			},
			status: 1,
		},
		{
			name:         "failure - eth call error",
			txHash:       common.HexToHash("0xc48fbf05edaf18f6aaa7de24de28528546b874bb03728d624ca407b8fed582a3"),
			ethCallError: errors.New("eth call failed"),
		},
		{
			name:   "failure - tx does not exist",
			txHash: common.HexToHash("0xc48fbf05edaf18f6aaa7de24de28528546b874bb03728d624ca407b8fed582a3"),
			receipt: &types.Receipt{
				Status: 0,
			},
			status: 0,
		},
	}

	for _, tt := range tests {
		client := new(clienttest.Client)
		client.On("CallContext", mock.Anything, mock.Anything, "eth_getTransactionReceipt", tt.txHash).
			Return(tt.ethCallError).Run(func(args mock.Arguments) {
			receipt := tt.receipt
			if receipt != nil {
				res := args.Get(1).(*types.Receipt)
				res.Status = receipt.Status
				res.TxHash = receipt.TxHash
				res.BlockNumber = receipt.BlockNumber
				res.BlockHash = receipt.BlockHash
			}
		})

		bn, bh, err := GetTxBlock(testutils.Context(t), client, tt.txHash)
		if tt.ethCallError != nil {
			assert.Equal(t, tt.ethCallError, err)
		} else {
			assert.Equal(t, tt.status, tt.receipt.Status)
			assert.Equal(t, tt.receipt.BlockNumber, bn)
			assert.Equal(t, tt.receipt.BlockHash, bh)
		}
	}
}
