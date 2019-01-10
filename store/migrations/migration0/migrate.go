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
	// orm.DB.LogMode(true) to trace sql commands
	var err error
	err = multierr.Append(err, migrationHelper(orm, &models.JobSpec{}))
	err = multierr.Append(err, migrationHelper(orm, &models.TaskSpec{}))
	err = multierr.Append(err, migrationHelper(orm, &models.JobRun{}))
	err = multierr.Append(err, migrationHelper(orm, &models.TaskRun{}))
	err = multierr.Append(err, migrationHelper(orm, &models.RunResult{}))
	err = multierr.Append(err, migrationHelper(orm, &models.Initiator{}))
	err = multierr.Append(err, migrationHelper(orm, &models.Tx{}))
	err = multierr.Append(err, migrationHelper(orm, &models.TxAttempt{}))
	err = multierr.Append(err, migrationHelper(orm, &models.BridgeType{}))
	err = multierr.Append(err, migrationHelper(orm, &models.IndexableBlockNumber{}))
	err = multierr.Append(err, migrationHelper(orm, &models.User{}))
	err = multierr.Append(err, migrationHelper(orm, &models.Session{}))
	err = multierr.Append(err, migrationHelper(orm, &models.Encumbrance{}))
	err = multierr.Append(err, migrationHelper(orm, &models.ServiceAgreement{}))
	err = multierr.Append(err, migrationHelper(orm, &models.BulkDeleteRunTask{}))
	err = multierr.Append(err, migrationHelper(orm, &models.BulkDeleteRunRequest{}))
	return err
}

func migrationHelper(orm *orm.ORM, model interface{}) error {
	db := orm.DB
	err := db.AutoMigrate(model).Error
	if err != nil {
		err = multierr.Append(fmt.Errorf("Migration for %T failed", model), err)
	}
	return err
}
