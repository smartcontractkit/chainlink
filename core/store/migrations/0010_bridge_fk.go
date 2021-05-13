package migrations

import (
	"gorm.io/gorm"
)

const (
	up10 = `
ALTER TABLE pipeline_task_specs ADD COLUMN bridge_name text;
ALTER TABLE pipeline_task_specs ADD CONSTRAINT fk_pipeline_task_specs_bridge_name FOREIGN KEY (bridge_name) REFERENCES bridge_types (name);
UPDATE pipeline_task_specs SET bridge_name = ts.json->>'name' FROM pipeline_task_specs ts WHERE ts.type = 'bridge';
`
	down10 = `
ALTER TABLE pipeline_task_specs DROP CONSTRAINT fk_pipeline_task_specs_bridge_name, DROP COLUMN bridge_name;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0010_bridge_fk",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up10).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down10).Error
		},
	})
}
