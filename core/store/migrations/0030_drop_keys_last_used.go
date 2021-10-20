package migrations

import "gorm.io/gorm"

const (
	up30 = `ALTER TABLE keys DROP COLUMN last_used`

	down30 = `ALTER TABLE keys ADD COLUMN last_used timestamptz`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0030_drop_keys_last_used",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up30).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down30).Error
		},
	})
}
