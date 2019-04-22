package migration0

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.JobSpec{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.TaskSpec{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.JobRun{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.TaskRun{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.RunResult{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.Initiator{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.Tx{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.TxAttempt{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&migration1551816486.BridgeType{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.Head{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.User{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.Session{}).Error; err != nil {
		return err
	}
	if err := tx.AutoMigrate(&models.Encumbrance{}).Error; err != nil {
		return err
	}
	return tx.AutoMigrate(&models.ServiceAgreement{}).Error
}
