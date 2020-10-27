package postgres

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func GormTransaction(ctx context.Context, db *gorm.DB, fc func(tx *gorm.DB) error, txOptss ...sql.TxOptions) (err error) {
	var txOpts sql.TxOptions
	if len(txOptss) > 0 {
		txOpts = txOptss[0]
	} else {
		txOpts = DefaultSqlTxOptions
	}
	tx := db.BeginTx(ctx, &txOpts)
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("%+v", r)
			tx.Rollback()
			// Rethrow the panic in case the calling code finds that desirable
			panic(r)
		}
	}()

	err = fc(tx)

	if err == nil {
		err = errors.WithStack(tx.Commit().Error)
	}

	// Make sure to rollback in case of a Block error or Commit error
	if err != nil {
		tx.Rollback()
	}
	return
}
