package migration0

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v3"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.BridgeType{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate BridgeType")
	}
	if err := tx.AutoMigrate(&models.Encumbrance{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Encumbrance")
	}
	if err := tx.AutoMigrate(&models.ExternalInitiator{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate ExternalInitiator")
	}
	if err := tx.AutoMigrate(&Head{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Head")
	}
	if err := tx.AutoMigrate(JobSpec{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate JobSpec")
	}
	if err := tx.AutoMigrate(&models.Initiator{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Initiator")
	}
	if err := tx.AutoMigrate(&models.JobRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate JobRun")
	}
	if err := tx.AutoMigrate(&models.Key{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Key")
	}
	if err := tx.AutoMigrate(&RunRequest{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunRequest")
	}
	if err := tx.AutoMigrate(&models.RunResult{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunResult")
	}
	if err := tx.AutoMigrate(&models.ServiceAgreement{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate ServiceAgreement")
	}
	if err := tx.AutoMigrate(&models.Session{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Session")
	}
	if err := tx.AutoMigrate(&models.SyncEvent{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate SyncEvent")
	}
	if err := tx.AutoMigrate(&TaskRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskRun")
	}
	if err := tx.AutoMigrate(&models.TaskSpec{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskSpec")
	}
	if err := tx.AutoMigrate(&TxAttempt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxAttempt")
	}
	if err := tx.AutoMigrate(&Tx{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Tx")
	}
	if err := tx.AutoMigrate(&models.User{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate User")
	}
	return nil
}

// Tx is a capture of the model representing Txes before migration1559081901
type Tx struct {
	ID        uint64         `gorm:"primary_key;auto_increment"`
	From      common.Address `gorm:"index;not null"`
	To        common.Address `gorm:"not null"`
	Data      []byte
	Nonce     uint64      `gorm:"index"`
	Value     *models.Big `gorm:"type:varchar(255)"`
	GasLimit  uint64
	Hash      common.Hash
	GasPrice  *models.Big `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
}

// TxAttempt is a capture of the model representing TxAttempts before migration1559081901
type TxAttempt struct {
	Hash      common.Hash `gorm:"primary_key;not null"`
	TxID      uint64      `gorm:"index"`
	GasPrice  *models.Big `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
	CreatedAt time.Time `gorm:"index"`
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                   string    `json:"id" gorm:"primary_key;not null"`
	JobRunID             string    `json:"-" gorm:"index;not null;type:varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE"`
	ResultID             uint      `json:"-"`
	Status               string    `json:"status"`
	TaskSpecID           uint      `json:"-" gorm:"index;not null REFERENCES task_specs(id)"`
	MinimumConfirmations uint64    `json:"minimumConfirmations"`
	CreatedAt            time.Time `json:"-" gorm:"index"`
}

// Head is a capture of the model representing Head before migration1560881846
type Head struct {
	HashRaw string `gorm:"primary_key;type:varchar;column:hash"`
	Number  int64  `gorm:"index;type:bigint;not null"`
}

// RunRequest stores the fields used to initiate the parent job run.
type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	Requester *common.Address
	CreatedAt time.Time
}

// JobSpec is a capture of the model representing Head before migration1565139192
type JobSpec struct {
	ID         string             `json:"id,omitempty" gorm:"primary_key;not null"`
	CreatedAt  time.Time          `json:"createdAt" gorm:"index"`
	Initiators []models.Initiator `json:"initiators"`
	Tasks      []models.TaskSpec  `json:"tasks"`
	StartAt    null.Time          `json:"startAt" gorm:"index"`
	EndAt      null.Time          `json:"endAt" gorm:"index"`
	DeletedAt  null.Time          `json:"-" gorm:"index"`
}
