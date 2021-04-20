package migrations

import (
	"gorm.io/gorm"
)

const (
	up21 = `ALTER TABLE initiators ADD COLUMN job_id_topic_filter uuid;`

	down21 = `ALTER TABLE initiators DROP COLUMN job_id_topic_filter;`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0021_add_job_id_topic_filter",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up21).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down21).Error
		},
	})
}
