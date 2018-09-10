package migrations

import (
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/models/orm"
)

func init() {
	registerMigration(Migration1536521223{})
}

type Migration1536521223 struct{}

func (m Migration1536521223) Timestamp() string {
	return "1536521223"
}

func (m Migration1536521223) Migrate(orm *orm.ORM) error {
	orm.InitializeModel(&models.JobSpec{})
	orm.InitializeModel(&models.JobRun{})
	orm.InitializeModel(&models.Initiator{})
	orm.InitializeModel(&models.Tx{})
	orm.InitializeModel(&models.TxAttempt{})
	orm.InitializeModel(&models.BridgeType{})
	orm.InitializeModel(&models.IndexableBlockNumber{})
	orm.InitializeModel(&models.User{})
	orm.InitializeModel(&models.Session{})
	orm.InitializeModel(&models.ServiceAgreement{})
	return nil
}
