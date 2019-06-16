package migration1559767168

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v3"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(JobSpec{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate JobSpec")
	}
	return nil
}

type JobSpec struct {
	ID         string             `json:"id,omitempty" gorm:"primary_key;not null"`
	CreatedAt  time.Time          `json:"createdAt" gorm:"index"`
	Initiators []models.Initiator `json:"initiators"`
	MinPayment *assets.Link       `json:"minPayment" gorm:"type:varchar(255)"`
	Tasks      []models.TaskSpec  `json:"tasks"`
	StartAt    null.Time          `json:"startAt" gorm:"index"`
	EndAt      null.Time          `json:"endAt" gorm:"index"`
	DeletedAt  null.Time          `json:"-" gorm:"index"`
}
