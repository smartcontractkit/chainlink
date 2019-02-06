package migration1549496047

import (
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1549496047"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	return orm.DB.AutoMigrate(&models.Key{}).Error
}
