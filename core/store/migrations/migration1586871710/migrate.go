package migration1586871710

import (
	"github.com/jinzhu/gorm"
)

// Migrate converts all 32-bit primary keys into 64 bit integers.
// In the (hopefully likely!) scenario that Chainlink is around for a long time, it's actually possible we might experience ID wraparounds.
// If a node processes more than 2.147B jobs it could have undefined behaviour.
// This is not strictly necessary for small tables e.g. configurations, but the cost is small and for consistency it is simply easier to use 64-bit integers as ID everywhere.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE configurations ALTER COLUMN id TYPE bigint;
	ALTER TABLE encumbrances ALTER COLUMN id TYPE bigint;
	ALTER TABLE external_initiators ALTER COLUMN id TYPE bigint;
	ALTER TABLE initiators ALTER COLUMN id TYPE bigint;
	ALTER TABLE job_runs ALTER COLUMN result_id TYPE bigint;
	ALTER TABLE job_runs ALTER COLUMN initiator_id TYPE bigint;
	ALTER TABLE job_runs ALTER COLUMN run_request_id TYPE bigint;
	ALTER TABLE run_requests ALTER COLUMN id TYPE bigint;
	ALTER TABLE run_results ALTER COLUMN id TYPE bigint;
	ALTER TABLE service_agreements ALTER COLUMN encumbrance_id TYPE bigint;
	ALTER TABLE sync_events ALTER COLUMN id TYPE bigint;
	ALTER TABLE task_runs ALTER COLUMN result_id TYPE bigint;
	ALTER TABLE task_runs ALTER COLUMN task_spec_id TYPE bigint;
	ALTER TABLE task_specs ALTER COLUMN id TYPE bigint;
	`).Error
}
