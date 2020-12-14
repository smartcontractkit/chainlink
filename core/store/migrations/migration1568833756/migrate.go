package migration1568833756

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

// Initiator could be thought of as a trigger, defines how a Job can be
// started, or rather, how a JobRun can be created from a Job.
// Initiators will have their own unique ID, but will be associated
// to a parent JobID.
type Initiator struct {
	ID         uint       `gorm:"primary_key;auto_increment"`
	JobSpecID  *models.ID `gorm:"index;type:varchar(36) REFERENCES job_specs(id)"`
	Type       string     `gorm:"index;not null"`
	CreatedAt  time.Time  `gorm:"index"`
	Schedule   models.Cron
	Time       models.AnyTime
	Ran        bool
	Address    common.Address           `gorm:"index"`
	Requesters models.AddressCollection `gorm:"type:text"`
	Name       string
	Params     string
	DeletedAt  null.Time `gorm:"index"`
}

// Migrate Initiator parameter 'Text' to support External Initaitor generic
// JSON parameters.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&Initiator{}).Error; err != nil {
		return errors.Wrap(err, "could not add fields 'names' and 'params' to the Initiator table")
	}
	return nil
}
