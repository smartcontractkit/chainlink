package migrations

import (
	"gorm.io/gorm"
)

const up51 = `
ALTER TABLE webhook_specs DROP COLUMN external_initiator_name, DROP COLUMN external_initiator_spec;

CREATE TABLE external_initiator_webhook_specs (
	external_initiator_id bigint NOT NULL references external_initiators (id) ON DELETE RESTRICT DEFERRABLE INITIALLY IMMEDIATE,
	webhook_spec_id int NOT NULL references webhook_specs (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE,
	spec jsonb NOT NULL,
	PRIMARY KEY (external_initiator_id, webhook_spec_id)
);

CREATE INDEX idx_external_initiator_webhook_specs_webhook_spec_id ON external_initiator_webhook_specs (webhook_spec_id);
CREATE UNIQUE INDEX idx_jobs_unique_flux_monitor_spec_id ON jobs (flux_monitor_spec_id);
CREATE UNIQUE INDEX idx_jobs_unique_keeper_spec_id ON jobs (keeper_spec_id);
CREATE UNIQUE INDEX idx_jobs_unique_cron_spec_id ON jobs (cron_spec_id);
CREATE UNIQUE INDEX idx_jobs_unique_vrf_spec_id ON jobs (vrf_spec_id);
CREATE UNIQUE INDEX idx_jobs_unique_webhook_spec_id ON jobs (webhook_spec_id);
`

const down51 = `
DROP TABLE external_initiator_webhook_specs;
ALTER TABLE webhook_specs ADD COLUMN external_initiator_name text, ADD COLUMN external_initiator_spec text;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0051_webhook_specs_external_initiators_join",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up51).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down51).Error
		},
	})
}
