package migration1568833756

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v3"
)

// InitiatorParams is a collection of the possible parameters that different
// Initiators may require.
type InitiatorParams struct {
	Schedule   models.Cron              `json:"schedule,omitempty"`
	Time       models.AnyTime           `json:"time,omitempty"`
	Ran        bool                     `json:"ran,omitempty"`
	Address    common.Address           `json:"address,omitempty" gorm:"index"`
	Requesters models.AddressCollection `json:"requesters,omitempty" gorm:"type:text"`
	Name       string                   `json:"name,omitempty"`
	Params     string                   `json:"-"`
}

// Initiator could be thought of as a trigger, defines how a Job can be
// started, or rather, how a JobRun can be created from a Job.
// Initiators will have their own unique ID, but will be associated
// to a parent JobID.
type Initiator struct {
	ID              uint       `json:"id" gorm:"primary_key;auto_increment"`
	JobSpecID       *models.ID `json:"jobSpecId" gorm:"index;type:varchar(36) REFERENCES job_specs(id)"`
	Type            string     `json:"type" gorm:"index;not null"`
	CreatedAt       time.Time  `gorm:"index"`
	InitiatorParams `json:"params,omitempty"`
	DeletedAt       null.Time `json:"-" gorm:"index"`
}

// Migrate Initiator parameter 'Text' to support External Initaitor generic
// JSON parameters.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&Initiator{}).Error; err != nil {
		return errors.Wrap(err, "could not add 'text' field to Initiator table")
	}
	return nil
}
