package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func TestUtil_GetTxBlock(t *testing.T) {
	tests := []struct {
		name         string
		txHash       common.Hash
		ethCallError error
		receipt      *types.Receipt
	}{
		{
			name:   "success",
			txHash: common.HexToHash("0xc48fbf05edaf18f6aaa7de24de28528546b874bb03728d624ca407b8fed582a3"),
			receipt: &types.Receipt{
				Status:      1,
				BlockNumber: big.NewInt(2000),
			},
		},
		{
			name:         "failure - eth call error",
			txHash:       common.HexToHash("0xc48fbf05edaf18f6aaa7de24de28528546b874bb03728d624ca407b8fed582a3"),
			ethCallError: fmt.Errorf("eth call failed"),
		},
	}

	for _, tt := range tests {
		client := new(evmClientMocks.Client)
		var h [32]byte
		copy(h[:], tt.txHash.Bytes())
		client.On("CallContext", mock.Anything, mock.Anything, "eth_getTransactionReceipt", h).
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

		bn, bh, err := GetTxBlock(client, tt.txHash)
		if tt.ethCallError != nil {
			assert.Equal(t, tt.ethCallError, err)
		} else {
			assert.Equal(t, tt.receipt.BlockNumber, bn)
			assert.Equal(t, tt.receipt.BlockHash, bh)
		}
	}
}
