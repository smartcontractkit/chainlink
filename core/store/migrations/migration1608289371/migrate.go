package migration1608289371

import "github.com/jinzhu/gorm"

// Migrate renames eth_request_event_specs to direct_request_specs and adds on_chain_job_spec_id
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE eth_request_event_specs RENAME TO direct_request_specs;
		ALTER TABLE direct_request_specs ADD COLUMN on_chain_job_spec_id bytea NOT NULL CHECK (octet_length(on_chain_job_spec_id) = 32);

		CREATE UNIQUE INDEX idx_direct_request_specs_unique_job_spec_id ON direct_request_specs (on_chain_job_spec_id);

		ALTER TABLE jobs RENAME COLUMN eth_request_event_spec_id TO direct_request_spec_id;

		ALTER INDEX idx_jobs_unique_eth_request_event_spec_id RENAME TO idx_jobs_unique_direct_request_spec_id;
		ALTER INDEX eth_request_event_specs_pkey RENAME TO direct_request_specs_pkey;
		ALTER TABLE jobs RENAME CONSTRAINT "jobs_eth_request_event_spec_id_fkey" TO "jobs_direct_request_spec_id_fkey";
    `).Error
}
