package migrations

import (
	"gorm.io/gorm"
)

const up18 = `
CREATE TABLE IF NOT EXISTS "node_versions" (
    "version" TEXT PRIMARY KEY,
    "created_at" timestamp without time zone NOT NULL
);
`

const down18 = `
DROP TABLE IF EXISTS "node_versions";
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0018_add_node_version_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up18).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down18).Error
		},
	})
}
