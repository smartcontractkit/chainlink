package utils

import (
	"context"
	"database/sql"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// CheckOKToTransmit returns an error if the transaction is not OK to transmit
// based on existing eth_txes in the database.
//
// NOTE: This is in the utils package to avoid import cycles, since it is used
// in both offchainreporting and adapters. Tests can be found in
// bulletprooftxmanager_test.go
func CheckOKToTransmit(ctx context.Context, db *sql.DB, fromAddress gethCommon.Address, maxUnconfirmedTransactions uint64) (err error) {
	if maxUnconfirmedTransactions == 0 {
		return nil
	}
	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, `SELECT count(*) FROM eth_txes WHERE from_address = $1 AND state IN ('unstarted', 'in_progress', 'unconfirmed')`, fromAddress)
	if err != nil {
		err = errors.Wrap(err, "bulletprooftxmanager.CheckOKToTransmit query failed")
		return
	}
	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()
	var count uint64
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			err = errors.Wrap(err, "bulletprooftxmanager.CheckOKToTransmit scan failed")
			return
		}
	}

	if count > maxUnconfirmedTransactions {
		err = errors.Errorf("cannot transmit eth transaction; there are currently %v unconfirmed transactions in the queue which exceeds the configured maximum of %v", count, maxUnconfirmedTransactions)
	}
	return
}
