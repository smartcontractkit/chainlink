package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const up9 = `
ALTER TABLE pipeline_task_specs ADD COLUMN bridge_name text;
ALTER TABLE pipeline_task_specs ADD CONSTRAINT fk_pipeline_task_specs_bridge_name FOREIGN KEY (bridge_name) REFERENCES bridge_types (name);
UPDATE pipeline_task_specs SET bridge_name = ts.json->>'name' FROM pipeline_task_specs ts WHERE ts.type = 'bridge';
`

const down9 = `
ALTER TABLE pipeline_task_specs DROP CONSTRAINT fk_pipeline_task_specs_bridge_name;
ALTER TABLE DROP COLUMN bridge_name;
`

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0009_bridge_fk",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up9).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down9).Error
		},
	})
}
