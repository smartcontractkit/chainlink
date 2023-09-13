package core

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// GetTxBlock calls eth_getTransactionReceipt on the eth client to obtain a tx receipt
func GetTxBlock(ctx context.Context, client client.Client, txHash common.Hash) (*big.Int, common.Hash, error) {
	receipt := types.Receipt{}
	err := client.CallContext(ctx, &receipt, "eth_getTransactionReceipt", txHash)
	if err != nil {
		if strings.Contains(err.Error(), "not yet been implemented") {
			// workaround for simulated chains
			// Exploratory: fix this properly (e.g. in the simulated backend)
			r, err1 := client.TransactionReceipt(ctx, txHash)
			if err1 != nil {
				return nil, common.Hash{}, err1
			}
			if r.Status != 1 {
				return nil, common.Hash{}, nil
			}
			return r.BlockNumber, r.BlockHash, nil
		}
		return nil, common.Hash{}, err
	}

	if receipt.Status != 1 {
		return nil, common.Hash{}, nil
	}

	return receipt.GetBlockNumber(), receipt.GetBlockHash(), nil
}
