package migration0

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "0"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	return multierr.Combine(
		migrationHelper(orm, &models.JobSpec{}),
		migrationHelper(orm, &models.TaskSpec{}),
		migrationHelper(orm, &models.JobRun{}),
		migrationHelper(orm, &models.TaskRun{}),
		migrationHelper(orm, &models.RunResult{}),
		migrationHelper(orm, &models.Initiator{}),
		migrationHelper(orm, &models.Tx{}),
		migrationHelper(orm, &models.TxAttempt{}),
		migrationHelper(orm, &models.BridgeType{}),
		migrationHelper(orm, &models.IndexableBlockNumber{}),
		migrationHelper(orm, &models.User{}),
		migrationHelper(orm, &models.Session{}),
		migrationHelper(orm, &models.Encumbrance{}),
		migrationHelper(orm, &models.ServiceAgreement{}),
		migrationHelper(orm, &models.BulkDeleteRunTask{}),
		migrationHelper(orm, &models.BulkDeleteRunRequest{}))
}

func migrationHelper(orm *orm.ORM, model interface{}) error {
	db := orm.DB
	err := db.AutoMigrate(model).Error
	if err != nil {
		err = multierr.Append(fmt.Errorf("Migration for %T failed", model), err)
	}
	return err
}
