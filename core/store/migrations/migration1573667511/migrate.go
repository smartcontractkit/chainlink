package migration1573667511

import (
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1568280052"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type ExternalInitiator struct {
	*gorm.Model
	Name           string `gorm:"not null;unique"`
	URL            models.WebURL
	AccessKey      string `gorm:"not null"`
	Salt           string `gorm:"not null"`
	HashedSecret   string `gorm:"not null"`
	OutgoingSecret string `gorm:"not null"`
	OutgoingToken  string `gorm:"not null"`
}

// newExternalInitiator returns a copy of the old struct with the fields untouched.
func newExternalInitiator(arg migration1568280052.ExternalInitiator) ExternalInitiator {
	return ExternalInitiator{
		Model:          arg.Model,
		Name:           arg.AccessKey,
		URL:            arg.URL,
		AccessKey:      arg.AccessKey,
		Salt:           arg.Salt,
		HashedSecret:   arg.HashedSecret,
		OutgoingSecret: arg.OutgoingSecret,
		OutgoingToken:  arg.OutgoingToken,
	}
}

// Migrate adds External Initiator Name and URL fields.
func Migrate(tx *gorm.DB) error {
	var exis []migration1568280052.ExternalInitiator
	if err := tx.Find(&exis).Error; err != nil {
		return errors.Wrap(err, "could not load all External Intitiators")
	}

	// Make new table
	if err := tx.DropTable(migration1568280052.ExternalInitiator{}).Error; err != nil {
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
