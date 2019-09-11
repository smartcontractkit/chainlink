package migration1566915476

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Migrate adds the name parameter in support of the new JobSpec's intiator
// 'external'.
func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`ALTER TABLE initiators ADD name varchar(255)`).Error; err != nil {
		return errors.Wrap(err, "failed to add name to Initiator")
	}
	return nil
}
