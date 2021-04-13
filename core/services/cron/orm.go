package cron

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

// ORM defines an interface for database commands related to Flux Monitor v2
type ORM interface {
	CreateEthTransaction(ctx context.Context, fromAddress common.Address, toAddress common.Address, payload []byte, gasLimit uint64, maxUnconfirmedTransactions uint64) error
}

type orm struct {
	db *gorm.DB
}

// NewORM initializes a new ORM
func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

func (orm *orm) CreateEthTransaction(ctx context.Context,
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	maxUnconfirmedTransactions uint64) error {

	db, err := orm.db.DB()
	if err != nil {
		return errors.Wrap(err, "orm#CreateEthTransaction")
	}

	err = utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
	if err != nil {
		return errors.Wrap(err, "transmitter#CreateEthTransaction")
	}

	value := 0
	res, err := db.ExecContext(ctx, `
	INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
	SELECT $1,$2,$3,$4,$5,'unstarted',NOW()
	WHERE NOT EXISTS (
		SELECT 1 FROM eth_tx_attempts
		JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
		WHERE eth_txes.from_address = $1
			AND eth_txes.state = 'unconfirmed'
			AND eth_tx_attempts.state = 'insufficient_eth'
	);
`, fromAddress, toAddress, payload, value, gasLimit)
	if err != nil {
		return fmt.Errorf("error inserting transaction for cron job: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected from cron: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("error inserting transaction: no row inserted")
	}

	return nil
}
