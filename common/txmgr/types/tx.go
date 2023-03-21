package types

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// TxStrategy controls how txes are queued and sent
//
//go:generate mockery --quiet --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after tx insertion
	// It accepts the service responsible for deleting
	// unstarted txs and deletion options
	PruneQueue(pruneService UnstartedTxQueuePruner, qopt pg.QOpt) (n int64, err error)
}

// No constants are defined because states may vary based on chain
type TxAttemptState = string
