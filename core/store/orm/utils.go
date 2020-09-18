package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func GormTransaction(db *gorm.DB, fc func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("%s", r)
			tx.Rollback()
			return
		}
	}()

	err = fc(tx)

	if err == nil {
		err = errors.WithStack(tx.Commit().Error)
	}

	// Makesure rollback when Block error or Commit error
	if err != nil {
		tx.Rollback()
	}
	return
}
