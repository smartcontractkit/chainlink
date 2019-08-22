package migration1565877314

import (
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ExternalInitiator struct {
	*gorm.Model
	Name           string        `gorm:"not null,unique"`
	URL            models.WebURL `gorm:"not null"`
	AccessKey      string        `gorm:"not null"`
	Salt           string        `gorm:"not null"`
	HashedSecret   string        `gorm:"not null"`
	OutgoingSecret string        `gorm:"not null"`
	OutgoingToken  string        `gorm:"not null"`
}

// newExternalInitiator creates a new row, setting the Name to the AccessKey
func newExternalInitiator(arg migration0.ExternalInitiator) ExternalInitiator {
	url, _ := url.ParseRequestURI("https://unset.url")
	return ExternalInitiator{
		Model: arg.Model,

		Name:           arg.AccessKey,
		URL:            models.WebURL(*url),
		AccessKey:      arg.AccessKey,
		Salt:           arg.Salt,
		HashedSecret:   arg.HashedSecret,
		OutgoingSecret: utils.NewSecret(48),
		OutgoingToken:  utils.NewSecret(48),
	}
}

// Migrate adds External Initiator Name and URL fields.
func Migrate(tx *gorm.DB) error {
	var exis []migration0.ExternalInitiator
	if err := tx.Find(&exis).Error; err != nil {
		return errors.Wrap(err, "could not load all External Intitiators")
	}

	// Make new table
	if err := tx.DropTable(migration0.ExternalInitiator{}).Error; err != nil {
		return errors.Wrap(err, "could not drop old External Intitiator table")
	}
	if err := tx.AutoMigrate(&ExternalInitiator{}).Error; err != nil {
		return errors.Wrap(err, "could not create new External Intitiator table")
	}

	// Copy
	for _, old := range exis {
		exi := newExternalInitiator(old)
		if err := tx.Save(exi).Error; err != nil {
			return errors.Wrap(err, "could not save migrated version of External Initiator")
		}
	}

	return nil
}
