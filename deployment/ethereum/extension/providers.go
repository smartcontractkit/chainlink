package deployment_ethereum

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func ConfirmTx(client bind.DeployBackend, hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receipt, err := waitMined(ctx, client, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmed receipt for chain: %w", err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("receipt was nil for tx: %s", hash.Hex())
	}
	return receipt, nil
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func waitMined(ctx context.Context, b bind.DeployBackend, hash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	for {
		receipt, err := b.TransactionReceipt(ctx, hash)
		if err == nil {
			return receipt, nil
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func ConfirmTxWithMemoryBackend(backend *memory.Backend, hash common.Hash) (*types.Receipt, error) {
	for {
		backend.Commit()
		receipt, err := func() (*types.Receipt, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			return waitMined(ctx, backend, hash)
		}()

		if err != nil {
			return nil, fmt.Errorf("failed to get confirmed receipt for chain: %w", err)
		}
		if receipt == nil {
			return nil, fmt.Errorf("receipt was nil for tx: %s", hash.Hex())
		}
		return receipt, nil
	}
}
