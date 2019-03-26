package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1549496047"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551895034"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1552418531"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1553029703"
	gormigrate "gopkg.in/gormigrate.v1"
)

type migration interface {
	Migrate(tx *gorm.DB) error
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "0",
			Migrate: func(tx *gorm.DB) error {
				return migration0.Migration{}.Migrate(tx)
			},
		},
		{
			ID: "1549496047",
			Migrate: func(tx *gorm.DB) error {
				return migration1549496047.Migration{}.Migrate(tx)
			},
		},
		{
			ID: "1551816486",
			Migrate: func(tx *gorm.DB) error {
				return migration1551816486.Migration{}.Migrate(tx)
			},
		},
		{
			ID: "1551895034",
			Migrate: func(tx *gorm.DB) error {
				return migration1551895034.Migration{}.Migrate(tx)
			},
		},
		{
			ID: "1552418531",
			Migrate: func(tx *gorm.DB) error {
				return migration1552418531.Migration{}.Migrate(tx)
			},
		},
		{
			ID: "1553029703",
			Migrate: func(tx *gorm.DB) error {
				return migration1553029703.Migration{}.Migrate(tx)
			},
		},
	})

	return m.Migrate()
}
