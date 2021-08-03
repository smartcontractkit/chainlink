package bulletprooftxmanager

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

// TxStrategy controls how txes are queued and sent
//go:generate mockery --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved to eth_txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after eth_tx insertion
	PruneQueue(tx postgres.Queryer) (n int64, err error)
}

var _ TxStrategy = SendEveryStrategy{}

func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32) (strategy TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize)
	} else {
		strategy = SendEveryStrategy{}
	}
	return
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct{}

func (SendEveryStrategy) Subject() uuid.NullUUID                     { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(postgres.Queryer) (int64, error) { return 0, nil }

var _ TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy struct {
	subject   uuid.UUID
	queueSize uint32
}

func NewDropOldestStrategy(subject uuid.UUID, queueSize uint32) DropOldestStrategy {
	return DropOldestStrategy{subject, queueSize}
}

func (s DropOldestStrategy) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy) PruneQueue(tx postgres.Queryer) (n int64, err error) {
	res, err := tx.Exec(`
DELETE FROM eth_txes
WHERE state = 'unstarted' AND subject = ? AND
id < (
	SELECT min(id) FROM (
		SELECT id
		FROM eth_txes
		WHERE state = 'unstarted' AND subject = $1
		ORDER BY id DESC
		LIMIT $2
	) numbers
)`, s.subject, s.subject, s.queueSize)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
