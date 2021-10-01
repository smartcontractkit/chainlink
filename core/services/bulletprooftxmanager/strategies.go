package bulletprooftxmanager

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// TxStrategy controls how txes are queued and sent
//go:generate mockery --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved to eth_txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after eth_tx insertion
	PruneQueue(tx *gorm.DB) (n int64, err error)
	// Simulate indicates whether this transaction can be safely simulated using eth_call
	// Simulating transactions before send and aborting on revert can save gas
	// BE CAREFUL - not all transaction types are safe to simulate, e.g. if
	// they can call arbitrary user-specified code, because there could be a case where
	// it would erroneously fail during simulation but would succeed for real
	Simulate() bool
}

var _ TxStrategy = SendEveryStrategy{}

func NewQueueingTxStrategy(subject uuid.UUID, queueSize uint32, simulate bool) (strategy TxStrategy) {
	if queueSize > 0 {
		strategy = NewDropOldestStrategy(subject, queueSize, simulate)
	} else {
		strategy = SendEveryStrategy{simulate}
	}
	return
}

func NewSendEveryStrategy(simulate bool) TxStrategy {
	return SendEveryStrategy{simulate}
}

// SendEveryStrategy will always send the tx
type SendEveryStrategy struct {
	simulate bool
}

func (SendEveryStrategy) Subject() uuid.NullUUID             { return uuid.NullUUID{} }
func (SendEveryStrategy) PruneQueue(*gorm.DB) (int64, error) { return 0, nil }
func (s SendEveryStrategy) Simulate() bool                   { return s.simulate }

var _ TxStrategy = DropOldestStrategy{}

// DropOldestStrategy will send the newest N transactions, older ones will be
// removed from the queue
type DropOldestStrategy struct {
	subject   uuid.UUID
	queueSize uint32
	simulate  bool
}

func NewDropOldestStrategy(subject uuid.UUID, queueSize uint32, simulate bool) DropOldestStrategy {
	return DropOldestStrategy{subject, queueSize, simulate}
}

func (s DropOldestStrategy) Subject() uuid.NullUUID {
	return uuid.NullUUID{UUID: s.subject, Valid: true}
}

func (s DropOldestStrategy) PruneQueue(tx *gorm.DB) (n int64, err error) {
	res := tx.Exec(`
DELETE FROM eth_txes
WHERE state = 'unstarted' AND subject = ? AND
id < (
	SELECT min(id) FROM (
		SELECT id
		FROM eth_txes
		WHERE state = 'unstarted' AND subject = ?
		ORDER BY id DESC
		LIMIT ?
	) numbers
)`, s.subject, s.subject, s.queueSize)
	return res.RowsAffected, res.Error
}

func (s DropOldestStrategy) Simulate() bool {
	return s.simulate
}
