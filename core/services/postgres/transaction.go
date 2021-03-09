package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// NOTE: In an ideal world the timeouts below would be set to something sane in
// the postgres configuration by the user. Since we do not live in an ideal
// world, it is necessary to override them here.
//
// They cannot easily be set at a session level due to how Go's connection
// pooling works.
const (
	// LockTimeout controls the max time we will wait for any kind of database lock.
	// It's good to set this to _something_ because waiting for locks forever is really bad.
	LockTimeout = 15 * time.Second
	// IdleInTxSessionTimeout controls the max time we leave a transaction open and idle.
	// It's good to set this to _something_ because leaving transactions open forever is really bad.
	IdleInTxSessionTimeout = 1 * time.Hour
)

func GormTransaction(ctx context.Context, db *gorm.DB, fc func(tx *gorm.DB) error, txOptss ...sql.TxOptions) (err error) {
	var txOpts sql.TxOptions
	if len(txOptss) > 0 {
		txOpts = txOptss[0]
	} else {
		txOpts = DefaultSqlTxOptions
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = %v; SET LOCAL idle_in_transaction_session_timeout = %v;`, LockTimeout.Milliseconds(), IdleInTxSessionTimeout.Milliseconds())).Error
		if err != nil {
			return errors.Wrap(err, "error setting transaction timeouts")
		}
		return fc(tx)
	}, &txOpts)
}
