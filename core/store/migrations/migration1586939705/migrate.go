package migration1586939705

import (
	"github.com/jinzhu/gorm"
)

// Migrate sets not null on columns that should always have a value
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE bridge_types ALTER COLUMN url SET NOT NULL;
	ALTER TABLE bridge_types ALTER COLUMN confirmations SET DEFAULT 0;
	ALTER TABLE bridge_types ALTER COLUMN confirmations SET NOT NULL;
	ALTER TABLE bridge_types ALTER COLUMN incoming_token_hash SET NOT NULL;
	ALTER TABLE bridge_types ALTER COLUMN salt SET NOT NULL;
	ALTER TABLE bridge_types ALTER COLUMN outgoing_token SET NOT NULL;

	ALTER TABLE configurations ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE configurations ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE encrypted_secret_keys ALTER COLUMN public_key SET NOT NULL;
	ALTER TABLE encrypted_secret_keys ALTER COLUMN vrf_key SET NOT NULL;

	ALTER TABLE external_initiators ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE external_initiators ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE initiators ALTER COLUMN job_spec_id SET NOT NULL;
	ALTER TABLE initiators ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE job_runs ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE job_runs ALTER COLUMN updated_at SET NOT NULL;
	ALTER TABLE job_runs ALTER COLUMN job_spec_id SET NOT NULL;

	ALTER TABLE job_specs ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE keys ALTER COLUMN json SET NOT NULL;

	ALTER TABLE run_requests ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE service_agreements ALTER COLUMN created_at SET NOT NULL;
	
	ALTER TABLE sessions ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE sync_events ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE sync_events ALTER COLUMN updated_at SET NOT NULL;
	ALTER TABLE sync_events ALTER COLUMN body SET NOT NULL;

	ALTER TABLE task_runs ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE task_runs ALTER COLUMN task_spec_id SET NOT NULL;
	ALTER TABLE task_runs ALTER COLUMN job_run_id SET NOT NULL;

	ALTER TABLE task_specs ALTER COLUMN created_at SET NOT NULL;
	ALTER TABLE task_specs ALTER COLUMN updated_at SET NOT NULL;
	ALTER TABLE task_specs ALTER COLUMN job_spec_id SET NOT NULL;

	ALTER TABLE tx_attempts ALTER COLUMN tx_id SET NOT NULL;

	ALTER TABLE users ALTER COLUMN created_at SET NOT NULL;
	`).Error
}
