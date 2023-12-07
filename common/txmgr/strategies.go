package txmgr

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
)

var _ txmgrtypes.TxStrategy = SendEveryStrategy{}

// NewQueueingTxStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded if a queue size is specified, and otherwise does not drop transactions.
func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) (strategy txmgrtypes.TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize, queryTimeout)
	} else {
		strategy = SendEveryStrategy{}
	}
	return
}

// NewSendEveryStrategy creates a new TxStrategy that does not drop transactions.
func NewSendEveryStrategy() txmgrtypes.TxStrategy {
	return SendEveryStrategy{}
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct{}

func (SendEveryStrategy) Subject() uuid.NullUUID { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(ctx context.Context, pruneService txmgrtypes.UnstartedTxQueuePruner) (int64, error) {
	return 0, nil
}

var _ txmgrtypes.TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy struct {
	subject      uuid.UUID
	queueSize    uint32
	queryTimeout time.Duration
}

// NewDropOldestStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded.
func NewDropOldestStrategy(subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) DropOldestStrategy {
	return DropOldestStrategy{subject, queueSize, queryTimeout}
}

func (s DropOldestStrategy) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy) PruneQueue(ctx context.Context, pruneService txmgrtypes.UnstartedTxQueuePruner) (n int64, err error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	n, err = pruneService.PruneUnstartedTxQueue(ctx, s.queueSize, s.subject)
	if err != nil {
		return 0, fmt.Errorf("DropOldestStrategy#PruneQueue failed: %w", err)
	}
	return
}
