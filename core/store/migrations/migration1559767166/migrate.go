package migration1559767166

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&TaskRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskRun")
	}
	return nil
}

type TaskRun struct {
	ID                   string    `json:"id" gorm:"primary_key;not null"`
	JobRunID             string    `json:"-" gorm:"index;not null;type:varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE"`
	ResultID             uint      `json:"-"`
	Status               string    `json:"status"`
	TaskSpecID           uint      `json:"-" gorm:"index;not null REFERENCES task_specs(id)"`
	MinimumConfirmations uint64    `json:"minimumConfirmations"`
	Confirmations        uint64    `json:"confirmations" gorm:"default: 0;not null"`
	CreatedAt            time.Time `json:"-" gorm:"index"`
}
