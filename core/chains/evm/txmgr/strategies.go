package txmgr

import (
	"context"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/common/types"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// TxStrategy controls how txes are queued and sent
//
//go:generate mockery --quiet --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable] interface {
	// Subject will be saved to eth_txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after eth_tx insertion
	PruneQueue(orm ORM[ADDR, TX_HASH], q pg.Queryer) (n int64, err error)
}

var _ TxStrategy = SendEveryStrategy{}

// NewQueueingTxStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded if a queue size is specified, and otherwise does not drop transactions.
func NewQueueingTxStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable](subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) (strategy TxStrategy[ADDR, TX_HASH]) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy[ADDR, TX_HASH](subject, queueSize, queryTimeout)
	} else {
		strategy = SendEveryStrategy[ADDR, TX_HASH]{}
	}
	return
}

// NewSendEveryStrategy creates a new TxStrategy that does not drop transactions.
func NewSendEveryStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable]() txmgrtypes.TxStrategy[ADDR, TX_HASH] {
	return SendEveryStrategy[ADDR, TX_HASH]{}

}

// SendEveryStrategy will always send the tx
type SendEveryStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable] struct{}

func (SendEveryStrategy[ADDR, TX_HASH]) Subject() uuid.NullUUID { return uuid.NullUUID{} }
func (SendEveryStrategy[ADDR, TX_HASH]) PruneQueue(pruneService txmgrtypes.UnstartedTxQueuePruner, qopt pg.QOpt) (int64, error) {

	return 0, nil
}

var _ txmgrtypes.TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable] struct {
	subject      uuid.UUID
	queueSize    uint32
	queryTimeout time.Duration
}

// NewDropOldestStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded.
func NewDropOldestStrategy[ADDR commontypes.Hashable, TX_HASH commontypes.Hashable](subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) DropOldestStrategy[ADDR, TX_HASH] {
	return DropOldestStrategy[ADDR, TX_HASH]{subject, queueSize, queryTimeout}
}

func (s DropOldestStrategy[ADDR, TX_HASH]) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy[ADDR, TX_HASH]) PruneQueue(pruneService txmgrtypes.UnstartedTxQueuePruner, qopt pg.QOpt) (n int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)

	defer cancel()
	n, err = pruneService.PruneUnstartedTxQueue(s.queueSize, s.subject, pg.WithParentCtx(ctx), qopt)
	if err != nil {
		return 0, errors.Wrap(err, "DropOldestStrategy#PruneQueue failed")
	}
	return
}
