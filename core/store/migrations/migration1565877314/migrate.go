package migration1565877314

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type ExternalInitiatorOld struct {
	*gorm.Model
	AccessKey    string
	Salt         string
	HashedSecret string
}

func (ExternalInitiatorOld) TableName() string {
	return "external_initiators"
}

type ExternalInitiator struct {
	*gorm.Model
	Name         string        `gorm:"not null,unique"`
	URL          models.WebURL `gorm:"not null"`
	AccessKey    string
	Salt         string
	HashedSecret string
}

// newExternalInitiator creates a new row, setting the Name to the AccessKey
func newExternalInitiator(arg ExternalInitiatorOld) ExternalInitiator {
	url, _ := models.NewWebURL("https://unset.url")
	return ExternalInitiator{
		Model: arg.Model,

		Name:         arg.AccessKey,
		URL:          url,
		AccessKey:    arg.AccessKey,
		Salt:         arg.Salt,
		HashedSecret: arg.HashedSecret,
	}
}

// Migrate adds External Initiator Name and URL fields.
func Migrate(tx *gorm.DB) error {
	var exis []ExternalInitiatorOld
	if err := tx.Find(&exis).Error; err != nil {
		return errors.Wrap(err, "could not load all External Intitiators")
	}

	// Make new table
	if err := tx.DropTable(ExternalInitiatorOld{}).Error; err != nil {
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
