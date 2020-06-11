package migration0

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

// Migrate runs the initial migration
//
// Do not reference the canonical structs from `store/models`, here. Place the
// struct corresponding to the DB table at the bottom of this file, and
// reference that in the `AutoMigrate` call. We don't want tables being
// `AutoMigrate`d whenever we change the canonical struct, because that could
// lead to data loss. Make an explicit migration in its own
// ../migrations<rough_unix_time>/migrate.go file/package.
//
// For instance, the `JobRun` struct referenced here does not have a Payment
// field. The "payment" column is added to the job_runs table later, in
// ./migration1567029116/migrate.go.
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.BridgeType{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate BridgeType")
	}
	if err := tx.AutoMigrate(&Encumbrance{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Encumbrance")
	}
	if err := tx.AutoMigrate(&ExternalInitiator{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate ExternalInitiator")
	}
	if err := tx.AutoMigrate(&Head{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Head")
	}
	if err := tx.AutoMigrate(&JobSpec{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate JobSpec")
	}
	if err := tx.AutoMigrate(&Initiator{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Initiator")
	}
	if err := tx.AutoMigrate(&JobRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate JobRun")
	}
	if err := tx.AutoMigrate(&Key{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Key")
	}
	if err := tx.AutoMigrate(&RunRequest{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunRequest")
	}
	if err := tx.AutoMigrate(&RunResult{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunResult")
	}
	if err := tx.AutoMigrate(&ServiceAgreement{}).Error; err != nil {
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
	if err := tx.AutoMigrate(&TaskSpec{}).Error; err != nil {
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
	Nonce     uint64     `gorm:"index"`
	Value     *utils.Big `gorm:"type:varchar(255)"`
	GasLimit  uint64
	Hash      common.Hash
	GasPrice  *utils.Big `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
}

// TxAttempt is a capture of the model representing TxAttempts before migration1559081901
type TxAttempt struct {
	Hash      common.Hash `gorm:"primary_key;not null"`
	TxID      uint64      `gorm:"index"`
	GasPrice  *utils.Big  `gorm:"type:varchar(255)"`
	Confirmed bool
	Hex       string `gorm:"type:text"`
	SentAt    uint64
	CreatedAt time.Time `gorm:"index"`
}

// TaskRun stores the Task and represents the status of the
// Task to be ran.
type TaskRun struct {
	ID                   string `gorm:"primary_key;not null"`
	JobRunID             string `gorm:"index;not null;type:varchar(36) REFERENCES job_runs(id) ON DELETE CASCADE"`
	ResultID             uint
	Status               string
	TaskSpecID           uint `gorm:"index;not null REFERENCES task_specs(id)"`
	MinimumConfirmations uint64
	CreatedAt            time.Time `gorm:"index"`
}

// ExternalInitiator represents a user that can initiate runs remotely
type ExternalInitiator struct {
	*gorm.Model
	AccessKey    string
	Salt         string
	HashedSecret string
}

// Head is a capture of the model before migration1560881846
type Head struct {
	HashRaw string `gorm:"primary_key;type:varchar;column:hash"`
	Number  int64  `gorm:"index;type:bigint;not null"`
}

// RunRequest  is a capture of the model before migration1565139192
type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	Requester *common.Address
	CreatedAt time.Time
}

// JobSpec is a capture of the model before migration1565139192
type JobSpec struct {
	ID        string    `gorm:"primary_key;not null"`
	CreatedAt time.Time `gorm:"index"`
	StartAt   null.Time `gorm:"index"`
	EndAt     null.Time `gorm:"index"`
	DeletedAt null.Time `gorm:"index"`
}

// JobRun is a capture of the model before migration1565210496
type JobRun struct {
	ID             string `gorm:"primary_key;not null"`
	JobSpecID      string `gorm:"index;not null;type:varchar(36) REFERENCES job_specs(id)"`
	ResultID       uint
	RunRequestID   uint
	Status         string    `gorm:"index"`
	CreatedAt      time.Time `gorm:"index"`
	FinishedAt     null.Time
	UpdatedAt      time.Time
	InitiatorID    uint
	CreationHeight string `gorm:"type:varchar(255)"`
	ObservedHeight string `gorm:"type:varchar(255)"`
	OverridesID    uint
	DeletedAt      null.Time `gorm:"index"`
}

// Initiator is a capture of the model before migration1565210496
type Initiator struct {
	ID         uint      `gorm:"primary_key;auto_increment"`
	JobSpecID  string    `gorm:"index;type:varchar(36) REFERENCES job_specs(id)"`
	Type       string    `gorm:"index;not null"`
	CreatedAt  time.Time `gorm:"index"`
	DeletedAt  null.Time `gorm:"index"`
	Schedule   string
	Time       time.Time
	Ran        bool
	Address    common.Address `gorm:"index"`
	Requesters string         `gorm:"type:text"`
}

// ServiceAgreement is a capture of the model before migration1565291711
type ServiceAgreement struct {
	ID            string    `gorm:"primary_key"`
	CreatedAt     time.Time `gorm:"index"`
	EncumbranceID uint
	RequestBody   string
	Signature     string `gorm:"type:varchar(255)"`
	JobSpecID     string `gorm:"index;not null;type:varchar(36) REFERENCES job_specs(id)"`
}

// TaskSpec is a capture of the model before migration1565291711
type TaskSpec struct {
	gorm.Model
	JobSpecID     string `gorm:"index;type:varchar(36) REFERENCES job_specs(id)"`
	Type          string `gorm:"index;not null"`
	Confirmations clnull.Uint32
	Params        string `gorm:"type:text"`
}

// RunResult is a capture of the model before migration1567029116
type RunResult struct {
	ID           uint   `gorm:"primary_key;auto_increment"`
	Data         string `gorm:"type:text"`
	Status       string
	ErrorMessage string
	Amount       *assets.Link `gorm:"type:varchar(255)"`
}

// Encumbrance is a capture of DB model before migration1568390387
type Encumbrance struct {
	ID         uint         `gorm:"primary_key;auto_increment"`
	Payment    *assets.Link `gorm:"type:varchar(255)"`
	Expiration uint64
	// This is a models.AnyTime in models.Encumbrance, but the DB column is a datetime
	EndAt time.Time
	// This is a models.EIP55AddressCollection in models.Encumbrance, but the DB column is text.
	Oracles string `gorm:"type:text"`
}

type Key struct {
	Address   models.EIP55Address `gorm:"primary_key;type:varchar(64)"`
	JSON      models.JSON         `gorm:"type:text"`
	CreatedAt time.Time           `json:"-"`
	UpdatedAt time.Time           `json:"-"`
}
