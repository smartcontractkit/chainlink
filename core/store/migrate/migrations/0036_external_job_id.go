package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

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

func Up36(ctx context.Context, tx *sql.Tx) error {
	// Add the external ID column and remove type specific ones.
	if _, err := tx.ExecContext(ctx, up36_1); err != nil {
		return err
	}

	// Update all jobs to have an external_job_id.
	// We do this to avoid using the uuid postgres extension.
	var jobIDs []int32
	txx := sqlx.Tx{Tx: tx}
	if err := txx.SelectContext(ctx, &jobIDs, "SELECT id FROM jobs"); err != nil {
		return err
	}
	if len(jobIDs) != 0 {
		stmt := `UPDATE jobs AS j SET external_job_id = vals.external_job_id FROM (values `
		for i := range jobIDs {
			if i == len(jobIDs)-1 {
				stmt += fmt.Sprintf("(uuid('%s'), %d))", uuid.New(), jobIDs[i])
			} else {
				stmt += fmt.Sprintf("(uuid('%s'), %d),", uuid.New(), jobIDs[i])
			}
		}
		stmt += ` AS vals(external_job_id, id) WHERE vals.id = j.id`
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	// Add constraints on the external_job_id.
	if _, err := tx.ExecContext(ctx, up36_2); err != nil {
		return err
	}
	return nil
}

func Down36(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, down36); err != nil {
		return err
	}
	return nil
}

var Migration36 = goose.NewGoMigration(36, &goose.GoFunc{RunTx: Up36}, &goose.GoFunc{RunTx: Down36})
