package migration1573812490

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type user struct {
	Email          string    `json:"email" gorm:"primary_key"`
	HashedPassword string    `json:"hashedPassword"`
	CreatedAt      time.Time `json:"createdAt" gorm:"index"`
	TokenKey       string    `json:"tokenKey"`
	TokenSecret    string    `json:"tokenSecret"`
}

// Migrate adds fields 'TokenKey' and 'TokenSecret' to support Token Authentication of a user.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&user{}).Error; err != nil {
		return errors.Wrap(err, "could not add fields 'TokenKey' and 'TokenSecreet' to table")
	}
	return nil
}
