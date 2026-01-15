package products

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// WaitMinedFast is a method for Anvil's instant blocks mode to ovecrome bind.WaitMined ticker hardcode.
func WaitMinedFast(ctx context.Context, b bind.DeployBackend, txHash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(5 * time.Millisecond)
	defer queryTicker.Stop()
	for {
		receipt, err := b.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}
