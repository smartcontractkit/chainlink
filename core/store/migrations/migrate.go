package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	gormigrate "gopkg.in/gormigrate.v1"
)

type migration interface {
	Migrate(tx *gorm.DB) error
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(db *gorm.DB) error {
	options := gormigrate.DefaultOptions
	options.UseTransaction = true

	m := gormigrate.New(db, options, []*gormigrate.Migration{
		{
			ID:      "0",
			Migrate: migration0.Migrate,
		},
	})

	err := m.Migrate()
	if err != nil {
		return errors.Wrap(err, "error running migrations")
	}
	return nil
}
