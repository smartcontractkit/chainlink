package core

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// GetTxBlock calls eth_getTransactionReceipt on the eth client to obtain a tx receipt
func GetTxBlock(client client.Client, txHash common.Hash) (*big.Int, common.Hash, error) {
	receipt := types.Receipt{}
	var h [32]byte
	copy(h[:], txHash.Bytes())
	err := client.CallContext(context.Background(), &receipt, "eth_getTransactionReceipt", h)
	if err != nil {
		return nil, common.Hash{}, err
	}

	if receipt.Status != 1 {
		return nil, common.Hash{}, nil
	}

	return receipt.GetBlockNumber(), receipt.GetBlockHash(), nil
}
