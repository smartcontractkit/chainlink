package offchainreporting

import (
	"context"
	"database/sql"
	"encoding/hex"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type transmitter struct {
	db                         *sql.DB
	fromAddress                gethCommon.Address
	gasLimit                   uint64
	maxUnconfirmedTransactions uint64
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(sqldb *sql.DB, fromAddress gethCommon.Address, gasLimit, maxUnconfirmedTransactions uint64) Transmitter {
	return &transmitter{
		db:                         sqldb,
		fromAddress:                fromAddress,
		gasLimit:                   gasLimit,
		maxUnconfirmedTransactions: maxUnconfirmedTransactions,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress gethCommon.Address, payload []byte) error {
	err := utils.CheckOKToTransmit(ctx, t.db, t.fromAddress, t.maxUnconfirmedTransactions)
	if err != nil {
		return errors.Wrap(err, "transmitter#CreateEthTransaction")
	}

	value := 0
	res, err := t.db.ExecContext(ctx, `
INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
SELECT $1,$2,$3,$4,$5,'unstarted',NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM eth_tx_attempts
	JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
	WHERE eth_txes.from_address = $1
		AND eth_txes.state = 'unconfirmed'
		AND eth_tx_attempts.state = 'insufficient_eth'
);
`, t.fromAddress, toAddress, payload, value, t.gasLimit)
	if err != nil {
		return errors.Wrap(err, "transmitter failed to insert eth_tx")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "transmitter failed to get RowsAffected on eth_tx insert")
	}
	if rowsAffected == 0 {
		err := errors.Errorf("Skipped OCR transmission because wallet is out of eth: %s", t.fromAddress.Hex())
		logger.Warnw(err.Error(),
			"fromAddress", t.fromAddress,
			"toAddress", toAddress,
			"payload", "0x"+hex.EncodeToString(payload),
			"value", value,
			"gasLimit", t.gasLimit,
		)
		return err
	}
	return nil
}

func (t *transmitter) FromAddress() gethCommon.Address {
	return t.fromAddress
}
