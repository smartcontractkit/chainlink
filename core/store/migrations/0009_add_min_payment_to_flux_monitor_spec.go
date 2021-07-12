package migrations

import (
	"gorm.io/gorm"
)

const (
	up9 = `
ALTER TABLE flux_monitor_specs
ADD min_payment varchar(255);
`
	down9 = `
ALTER TABLE flux_monitor_specs
DROP COLUMN min_payment;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0009_add_min_payment_to_flux_monitor_spec",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up9).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down9).Error
		},
	})
}
