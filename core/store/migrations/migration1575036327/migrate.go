package migration1575036327

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type user struct {
	Email             string `gorm:"primary_key"`
	HashedPassword    string
	CreatedAt         time.Time `gorm:"index"`
	TokenKey          string
	TokenSalt         string
	TokenHashedSecret string
}

// Migrate adds fields 'TokenKey' and 'TokenSecret' to support Token Authentication of a user.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&user{}).Error; err != nil {
		return errors.Wrap(err, "could not add fields 'TokenKey' and 'TokenSecreet' to table")
	}
	return nil
}
