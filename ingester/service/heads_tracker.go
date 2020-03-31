package service

import (
	"database/sql"

	"chainlink/ingester/client"
	"chainlink/ingester/logger"

	"github.com/ethereum/go-ethereum/core/types"
)

// HeadsTracker is an interface for subscribing to new Chainlink aggregator feeds
type HeadsTracker interface {
	Start() error
	Stop()
}

type headsTracker struct {
	db  *sql.DB
	eth client.ETH
}

// NewHeadsTracker returns an instantiated instance of a HeadsTracker implementation
func NewHeadsTracker(db *sql.DB, eth client.ETH) HeadsTracker {
	ht := &headsTracker{
		db:  db,
		eth: eth,
	}
	return ht
}

// Start subscribes to new heads and records them in the DB
func (ht *headsTracker) Start() error {
	headChan := make(chan types.Header)
	if _, err := ht.eth.SubscribeToNewHeads(headChan); err != nil {
		return err
	}

	go func() {
		logger.Debug("Listening for heads")
		for head := range headChan {
			nonce := make([]byte, 8)
			copy(nonce, head.Nonce[:])

			logger.Debugw("Observed new head", "blockHeight", head.Number, "blockHash", head.Hash())
			_, err := ht.db.Exec(`INSERT INTO "ethereum_head" ("blockHash", "parentHash", "uncleHash", "coinbase", "root", "txHash", "receiptHash", "bloom", "difficulty", "number", "gasLimit", "gasUsed", "time", "extra", "mixDigest", "nonce") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
				head.Hash().Bytes(),
				head.ParentHash,
				head.UncleHash,
				head.Coinbase,
				head.Root,
				head.TxHash,
				head.ReceiptHash,
				head.Bloom.Bytes(),
				head.Difficulty.String(),
				head.Number.String(),
				head.GasLimit,
				head.GasUsed,
				head.Time,
				head.Extra,
				head.MixDigest,
				nonce)
			if err != nil {
				logger.Errorw("Insert failed", "error", err)
			}
		}
	}()

	return nil
}

// Stop
func (ht *headsTracker) Stop() {
}
