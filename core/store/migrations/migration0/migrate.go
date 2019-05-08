package migration0

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
	if err := tx.AutoMigrate(&models.Head{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Head")
	}
	if err := tx.AutoMigrate(&models.JobSpec{}).Error; err != nil {
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
	if err := tx.AutoMigrate(&models.RunRequest{}).Error; err != nil {
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
	if err := tx.AutoMigrate(&models.TaskRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskRun")
	}
	if err := tx.AutoMigrate(&models.TaskSpec{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskSpec")
	}
	if err := tx.AutoMigrate(&models.TxAttempt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxAttempt")
	}
	if err := tx.AutoMigrate(&models.Tx{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Tx")
	}
	if err := tx.AutoMigrate(&models.User{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate User")
	}
	return nil
}
