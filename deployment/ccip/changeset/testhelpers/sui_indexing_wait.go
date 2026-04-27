package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
)

// ErrSuiTxIndexingWaitTimeout is returned when the fullnode never served sui_getTransactionBlock
// for the digest within SuiTxIndexingWaitTimeout. Mirrors bind.ErrTxIndexingTimeout semantics;
// replace this helper with bind.WaitForTransactionIndexed once chainlink's minimum chainlink-sui
// version includes that API.
var ErrSuiTxIndexingWaitTimeout = errors.New("sui tx not visible on fullnode within timeout")

var (
	// SuiTxIndexingWaitTimeout bounds polling after a successful sui_executeTransactionBlock.
	SuiTxIndexingWaitTimeout = 150 * time.Second
	suiTxIndexingInitial     = 200 * time.Millisecond
	suiTxIndexingMax         = 8 * time.Second
)

// WaitForSuiFullnodeTransaction polls sui_getTransactionBlock until digest is returned.
// Use after any mutating Sui tx when the next RPC read must see updated owned-object versions
// (gas coin, pool state, etc.) on JSON-RPC fullnodes that ignore WaitForLocalExecution.
//
// ctx is reserved for future cancellation/propagation; the poll deadline is intentionally
// NOT derived from ctx. cldf.Environment.GetContext() typically carries the whole test's
// deadline; nesting context.WithTimeout(ctx, SuiTxIndexingWaitTimeout) caps the wait at
// whatever time remains on that parent and causes spurious "context deadline exceeded" on
// sui_getTransactionBlock long before SuiTxIndexingWaitTimeout elapses.
func WaitForSuiFullnodeTransaction(ctx context.Context, client sui.ISuiAPI, digest string) error {
	_ = ctx // API stability; poll uses an independent wall-clock budget (see doc above).
	if digest == "" {
		return nil
	}
	pollCtx, cancel := context.WithTimeout(context.Background(), SuiTxIndexingWaitTimeout)
	defer cancel()

	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	}

	backoff := suiTxIndexingInitial
	var lastErr error
	for {
		resp, err := client.SuiGetTransactionBlock(pollCtx, req)
		if err == nil && resp.Digest == digest {
			return nil
		}
		lastErr = err

		select {
		case <-pollCtx.Done():
			if lastErr != nil {
				return fmt.Errorf("%w (digest=%s): %w", ErrSuiTxIndexingWaitTimeout, digest, lastErr)
			}
			return fmt.Errorf("%w (digest=%s)", ErrSuiTxIndexingWaitTimeout, digest)
		case <-time.After(backoff):
		}

		if backoff < suiTxIndexingMax {
			backoff *= 2
			if backoff > suiTxIndexingMax {
				backoff = suiTxIndexingMax
			}
		}
	}
}
