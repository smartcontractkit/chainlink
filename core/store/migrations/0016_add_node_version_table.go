package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up16 = `
		CREATE TABLE "node_versions" (
        "version" TEXT PRIMARY KEY,
        "created_at" timestamp without time zone NOT NULL
    );
`

const down16 = `
    DROP TABLE IF EXISTS "node_version";
    DROP TABLE IF EXISTS "node_versions";
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0016_add_node_version_table",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up16).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down16).Error
		},
	})
}
