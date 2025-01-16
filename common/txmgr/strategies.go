package txmgr

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
)

var _ types.TxStrategy = SendEveryStrategy{}

// NewQueueingTxStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded if a queue size is specified, and otherwise does not drop transactions.
func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32) (strategy types.TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize)
	} else {
		strategy = SendEveryStrategy{}
	}
	return
}

// NewSendEveryStrategy creates a new TxStrategy that does not drop transactions.
func NewSendEveryStrategy() types.TxStrategy {
	return SendEveryStrategy{}
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct{}

func (SendEveryStrategy) Subject() uuid.NullUUID { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(ctx context.Context, pruneService types.UnstartedTxQueuePruner) ([]int64, error) {
	return nil, nil
}

var _ types.TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy struct {
	subject   uuid.UUID
	queueSize uint32
}

// NewDropOldestStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded.
func NewDropOldestStrategy(subject uuid.UUID, queueSize uint32) DropOldestStrategy {
	return DropOldestStrategy{subject, queueSize}
}

func (s DropOldestStrategy) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy) PruneQueue(ctx context.Context, pruneService types.UnstartedTxQueuePruner) (ids []int64, err error) {
	// NOTE: We prune one less than the queue size to prevent the queue from exceeding the max queue size. Which could occur if a new transaction is added to the queue right after we prune.
	ids, err = pruneService.PruneUnstartedTxQueue(ctx, s.queueSize-1, s.subject)
	if err != nil {
		return ids, fmt.Errorf("DropOldestStrategy#PruneQueue failed: %w", err)
	}
	return
}
