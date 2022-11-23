package txmgr

import (
	"context"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// TxStrategy controls how txes are queued and sent
//
//go:generate mockery --quiet --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved to eth_txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after eth_tx insertion
	PruneQueue(q pg.Queryer) (n int64, err error)
}

var _ TxStrategy = SendEveryStrategy{}

// NewQueueingTxStrategy creates a new TxStrategy that drops the oldest transactions after the
// queue size is exceeded if a queue size is specified, and otherwise does not drop transactions.
func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32, queryTimeout time.Duration) (strategy TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize, queryTimeout)
	} else {
		strategy = SendEveryStrategy{}
	}
	return
}

// NewSendEveryStrategy creates a new TxStrategy that does not drop transactions.
func NewSendEveryStrategy() TxStrategy {
	return SendEveryStrategy{}
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct{}

func (SendEveryStrategy) Subject() uuid.NullUUID               { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(pg.Queryer) (int64, error) { return 0, nil }

var _ TxStrategy = DropOldestStrategy{}

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

func (s DropOldestStrategy) PruneQueue(q pg.Queryer) (n int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()
	res, err := q.ExecContext(ctx, `
DELETE FROM eth_txes
WHERE state = 'unstarted' AND subject = $1 AND
id < (
	SELECT min(id) FROM (
		SELECT id
		FROM eth_txes
		WHERE state = 'unstarted' AND subject = $2
		ORDER BY id DESC
		LIMIT $3
	) numbers
)`, s.subject, s.subject, s.queueSize)
	if err != nil {
		return 0, errors.Wrap(err, "DropOldestStrategy#PruneQueue failed")
	}
	return res.RowsAffected()
}
