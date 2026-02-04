package products

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// SethKeyHelper wraps Seth's NonceManager to provide tx hash tracking and
// Anvil-specific transaction dropping when all keys are stuck.
// This is a thin wrapper that reuses Seth's AnySyncedKey() logic.
type SethKeyHelper struct {
	mu         sync.Mutex
	client     *seth.Client
	rpcClient  *rpc.Client
	logger     zerolog.Logger
	pendingTxs map[int]common.Hash
	isAnvil    bool
	rpcTimeout time.Duration
}

// NewSethKeyHelper creates a new helper that wraps Seth for tx tracking and recovery.
func NewSethKeyHelper(logger zerolog.Logger, client *seth.Client, rpcURL string, isAnvil bool) (*SethKeyHelper, error) {
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC: %w", err)
	}

	logger.Info().
		Bool("isAnvil", isAnvil).
		Int("keyCount", len(client.Addresses)).
		Msg("Initialized SethTxHelper")

	return &SethKeyHelper{
		client:     client,
		rpcClient:  rpcClient,
		logger:     logger,
		pendingTxs: make(map[int]common.Hash),
		isAnvil:    isAnvil,
		rpcTimeout: 10 * time.Second,
	}, nil
}

// GetKey wraps Seth's AnySyncedKey() and handles the TimeoutKeyNum case
// by dropping pending transactions (on Anvil) and retrying.
// Returns the key number or an error if no key is available.
func (h *SethKeyHelper) GetKey() (int, error) {
	keyNum := h.client.AnySyncedKey()

	if keyNum == seth.TimeoutKeyNum {
		h.logger.Warn().Msg("All keys have pending transactions, attempting recovery")

		if h.isAnvil {
			// Calculate timeout: rpcTimeout per pending tx + 10s buffer, minimum 30s
			dropCtx, dropCancel := context.WithTimeout(context.Background(), h.calculateDropTimeout())
			defer dropCancel()

			dropped, err := h.DropPendingTxs(dropCtx)
			if err != nil {
				h.logger.Error().Err(err).Msg("Error dropping pending transactions")
			} else {
				h.logger.Info().Int("dropped", dropped).Msg("Dropped pending transactions")
			}

			// Retry once after dropping
			keyNum = h.client.AnySyncedKey()
		}
	}

	if keyNum == seth.TimeoutKeyNum {
		return seth.TimeoutKeyNum, fmt.Errorf("all keys have pending transactions")
	}

	h.logger.Debug().Int("keyNum", keyNum).Msg("Got key from pool")
	return keyNum, nil
}

// calculateDropTimeout returns a reasonable timeout for dropping all pending txs
// Formula: (pendingCount * rpcTimeout) + 10s buffer, minimum 30s
func (h *SethKeyHelper) calculateDropTimeout() time.Duration {
	h.mu.Lock()
	pendingCount := len(h.pendingTxs)
	h.mu.Unlock()

	timeout := max(time.Duration(pendingCount)*h.rpcTimeout+10*time.Second, 30*time.Second)
	return timeout
}

// RecordPendingTx stores the transaction hash for a key.
// Call this after successfully sending a transaction.
func (h *SethKeyHelper) RecordPendingTx(keyNum int, txHash common.Hash) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pendingTxs[keyNum] = txHash
	h.logger.Debug().
		Int("keyNum", keyNum).
		Str("txHash", txHash.Hex()).
		Int("totalPending", len(h.pendingTxs)).
		Msg("Recorded pending tx")
}

// DropPendingTxs drops all tracked pending transactions from Anvil's mempool.
// Returns the number of transactions dropped.
func (h *SethKeyHelper) DropPendingTxs(ctx context.Context) (int, error) {
	if !h.isAnvil {
		h.logger.Debug().Msg("DropPendingTxs called but not running on Anvil, skipping")
		return 0, nil
	}

	h.mu.Lock()
	// Copy the map to avoid holding the lock during RPC calls
	toDropMap := make(map[int]common.Hash, len(h.pendingTxs))
	maps.Copy(toDropMap, h.pendingTxs)
	h.mu.Unlock()

	if len(toDropMap) == 0 {
		h.logger.Debug().Msg("No pending transactions to drop")
		return 0, nil
	}

	h.logger.Info().
		Int("pendingTxCount", len(toDropMap)).
		Msg("Dropping pending transactions from Anvil mempool")

	var dropped int
	var errors []string

	for keyNum, txHash := range toDropMap {
		if err := h.dropTransaction(ctx, txHash); err != nil {
			errors = append(errors, fmt.Sprintf("key %d (%s): %v", keyNum, txHash.Hex(), err))
			continue
		}

		h.mu.Lock()
		delete(h.pendingTxs, keyNum)
		h.mu.Unlock()

		dropped++
		h.logger.Debug().
			Int("keyNum", keyNum).
			Str("txHash", txHash.Hex()).
			Msg("Dropped pending transaction")
	}

	if len(errors) > 0 {
		h.logger.Warn().
			Int("dropped", dropped).
			Int("failed", len(errors)).
			Msg("Finished dropping pending transactions with some failures")
		return dropped, fmt.Errorf("failed to drop some transactions: %v", errors)
	}

	h.logger.Info().
		Int("dropped", dropped).
		Msg("Successfully dropped all pending transactions")
	return dropped, nil
}

// dropTransaction drops a single transaction from Anvil's mempool
func (h *SethKeyHelper) dropTransaction(ctx context.Context, txHash common.Hash) error {
	ctx, cancel := context.WithTimeout(ctx, h.rpcTimeout)
	defer cancel()

	var result any
	if err := h.rpcClient.CallContext(ctx, &result, "anvil_dropTransaction", txHash); err != nil {
		return fmt.Errorf("anvil_dropTransaction failed: %w", err)
	}
	return nil
}
