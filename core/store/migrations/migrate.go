package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559081901"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559767166"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560433987"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560791143"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881846"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560886530"
	gormigrate "gopkg.in/gormigrate.v1"
)

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(db *gorm.DB) error {
	options := *gormigrate.DefaultOptions
	options.UseTransaction = true

	m := gormigrate.New(db, &options, []*gormigrate.Migration{
		{
			ID:      "0",
			Migrate: migration0.Migrate,
		},
		{
			ID:      "1559081901",
			Migrate: migration1559081901.Migrate,
		},
		{
			ID:      "1559767166",
			Migrate: migration1559767166.Migrate,
		},
		{
			ID:      "1560433987",
			Migrate: migration1560433987.Migrate,
		},
		{
			ID:      "1560791143",
			Migrate: migration1560791143.Migrate,
		},
		{
			ID:      "1560881846",
			Migrate: migration1560881846.Migrate,
		},
		{
			ID:      "1560886530",
			Migrate: migration1560886530.Migrate,
		},
	})

	db.LogMode(true)
	err := m.Migrate()
	if err != nil {
		return errors.Wrap(err, "error running migrations")
	}
	return nil
}
