package migration1560881846

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`CREATE INDEX sync_events_id_created_at_idx ON sync_events ("id", "created_at")`).Error
	return errors.Wrap(err, "failed to create sync events id + created at index")
}
