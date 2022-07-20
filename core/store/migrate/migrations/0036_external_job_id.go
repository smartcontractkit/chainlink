package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/sqlx"
)

func init() {
	goose.AddMigration(Up36, Down36)
}

const (
	up36_1 = `
	ALTER TABLE direct_request_specs DROP COLUMN on_chain_job_spec_id;
	ALTER TABLE webhook_specs DROP COLUMN on_chain_job_spec_id;
	ALTER TABLE vrf_specs ADD CONSTRAINT vrf_specs_public_key_fkey FOREIGN KEY (public_key) REFERENCES encrypted_vrf_keys(public_key) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;
	ALTER TABLE jobs ADD COLUMN external_job_id uuid; 
	`
	up36_2 = `
	ALTER TABLE jobs 
		ALTER COLUMN external_job_id SET NOT NULL,
		ADD CONSTRAINT external_job_id_uniq UNIQUE(external_job_id),
		ADD CONSTRAINT non_zero_uuid_check CHECK (external_job_id <> '00000000-0000-0000-0000-000000000000');
	`
	down36 = `
	ALTER TABLE direct_request_specs ADD COLUMN on_chain_job_spec_id bytea;
	ALTER TABLE webhook_specs ADD COLUMN on_chain_job_spec_id bytea;
	ALTER TABLE jobs DROP CONSTRAINT external_job_id_uniq;
	ALTER TABLE vrf_specs DROP CONSTRAINT vrf_specs_public_key_fkey;
    `
)

//nolint
func Up36(tx *sql.Tx) error {
	// Add the external ID column and remove type specific ones.
	if _, err := tx.Exec(up36_1); err != nil {
		return err
	}

	// Update all jobs to have an external_job_id.
	// We do this to avoid using the uuid postgres extension.
	var jobIDs []int32
	txx := sqlx.NewTx(tx, "postgres")
	if err := txx.Select(&jobIDs, "SELECT id FROM jobs"); err != nil {
		return err
	}
	if len(jobIDs) != 0 {
		stmt := `UPDATE jobs AS j SET external_job_id = vals.external_job_id FROM (values `
		for i := range jobIDs {
			if i == len(jobIDs)-1 {
				stmt += fmt.Sprintf("(uuid('%s'), %d))", uuid.NewV4(), jobIDs[i])
			} else {
				stmt += fmt.Sprintf("(uuid('%s'), %d),", uuid.NewV4(), jobIDs[i])
			}
		}
		stmt += ` AS vals(external_job_id, id) WHERE vals.id = j.id`
		if _, err := tx.Exec(stmt); err != nil {
			return err

		}
	}

	// Add constraints on the external_job_id.
	if _, err := tx.Exec(up36_2); err != nil {
		return err
	}
	return nil
}

//nolint
func Down36(tx *sql.Tx) error {
	if _, err := tx.Exec(down36); err != nil {
		return err
	}
	return nil
}
