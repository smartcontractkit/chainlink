package migrations

import "gorm.io/gorm"

const (
	up33 = `
               CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
               ALTER TABLE direct_request_specs DROP COLUMN on_chain_job_spec_id;
               ALTER TABLE webhook_specs DROP COLUMN on_chain_job_spec_id;
               ALTER TABLE jobs ADD COLUMN external_job_id uuid NOT NULL DEFAULT uuid_generate_v4();
               ALTER TABLE jobs ADD CONSTRAINT external_job_id_uniq UNIQUE(external_job_id);
               ALTER TABLE jobs ADD CONSTRAINT non_zero_uuid_check CHECK (external_job_id <> '00000000-0000-0000-0000-000000000000');
               ALTER TABLE vrf_specs ADD CONSTRAINT vrf_specs_public_key_fkey FOREIGN KEY (public_key) REFERENCES encrypted_vrf_keys(public_key);
	`
	down33 = `
               ALTER TABLE direct_request_specs ADD COLUMN on_chain_job_spec_id bytea;
               ALTER TABLE webhook_specs ADD COLUMN on_chain_job_spec_id;
               ALTER TABLE jobs DROP CONSTRAINT external_job_id_uniq;
               ALTER TABLE vrf_specs DROP CONSTRAINT vrf_specs_public_key_fkey;
       `
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0033_external_job_id",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up33).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down33).Error
		},
	})
}
