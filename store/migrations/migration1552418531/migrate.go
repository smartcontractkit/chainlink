package migration1552418531

import (
	"github.com/smartcontractkit/chainlink/store/orm"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v3"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1552418531"
}

// Migrate creates a new bridge_types table with the correct primary key
// because sqlite does not allow you to modify the primary key
// after table creation.
func (m Migration) Migrate(orm *orm.ORM) error {
	return multierr.Combine(
		orm.DB.AutoMigrate(&initiator{}).Error,
		orm.DB.AutoMigrate(&jobSpec{}).Error,
		orm.DB.AutoMigrate(&jobRun{}).Error,
	)
}

type jobSpec struct {
	ID        string    `json:"id,omitempty" gorm:"primary_key;not null"`
	DeletedAt null.Time `json:"-" gorm:"index"`
}

type jobRun struct {
	ID        string    `json:"id" gorm:"primary_key;not null"`
	DeletedAt null.Time `json:"-" gorm:"index"`
}

type initiator struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment"`
	DeletedAt null.Time `json:"=" gorm:"index"`
}
