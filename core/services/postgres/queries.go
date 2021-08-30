package postgres

import (
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

// BatchSize is the default number of DB records to access in one batch
const BatchSize uint = 1000

// BatchFunc is the function to execute on each batch of records, should return the count of records affected
type BatchFunc func(offset, limit uint) (count uint, err error)

// Batch is an iterator for batches of records
func Batch(cb BatchFunc) error {
	offset := uint(0)
	limit := BatchSize

	for {
		count, err := cb(offset, limit)
		if err != nil {
			return err
		}

		if count < limit {
			return nil
		}

		offset += limit
	}
}

// Sessions returns all sessions limited by the parameters.
func Sessions(db *gorm.DB, offset, limit int) ([]models.Session, error) {
	var sessions []models.Session
	err := db.
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}
