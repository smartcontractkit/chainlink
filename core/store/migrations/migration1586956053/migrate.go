package migration1586956053

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds timestamps to tables that ought to have them, but don't
// Also add some indexes where appropriate
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE bridge_types ADD COLUMN IF NOT EXISTS created_at timestamptz;
	UPDATE bridge_types SET created_at = '2019-01-01' WHERE created_at IS NULL;
	CREATE INDEX idx_bridge_types_created_at ON bridge_types USING BRIN (created_at);
	ALTER TABLE bridge_types ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE bridge_types ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE bridge_types SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	CREATE INDEX idx_bridge_types_updated_at ON bridge_types USING BRIN (updated_at);
	ALTER TABLE bridge_types ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE encrypted_secret_keys ADD COLUMN IF NOT EXISTS created_at timestamptz;
	UPDATE encrypted_secret_keys SET created_at = '2019-01-01' WHERE created_at IS NULL;
	ALTER TABLE encrypted_secret_keys ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE encrypted_secret_keys ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE encrypted_secret_keys SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	ALTER TABLE encrypted_secret_keys ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE encumbrances ADD COLUMN created_at timestamptz;
	UPDATE encumbrances SET created_at = '2019-01-01';
	CREATE INDEX idx_encumbrances_created_at ON encumbrances USING BRIN (created_at);
	ALTER TABLE encumbrances ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE encumbrances ADD COLUMN updated_at timestamptz;
	UPDATE encumbrances SET updated_at = '2019-01-01';
	CREATE INDEX idx_encumbrances_updated_at ON encumbrances USING BRIN (updated_at);
	ALTER TABLE encumbrances ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE initiators ADD COLUMN updated_at timestamptz;
	UPDATE initiators SET updated_at = '2019-01-01';
	CREATE INDEX idx_initiators_updated_at ON initiators USING BRIN (updated_at);
	ALTER TABLE initiators ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE job_specs ADD COLUMN updated_at timestamptz;
	UPDATE job_specs SET updated_at = '2019-01-01';
	CREATE INDEX idx_job_specs_updated_at ON job_specs USING BRIN (updated_at);
	ALTER TABLE job_specs ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE keys ADD COLUMN IF NOT EXISTS created_at timestamptz;
	UPDATE keys SET created_at = '2019-01-01' WHERE created_at IS NULL;
	ALTER TABLE keys ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE keys ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE keys SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	ALTER TABLE keys ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE run_results ADD COLUMN created_at timestamptz;
	UPDATE run_results SET created_at = '2019-01-01';
	CREATE INDEX idx_run_results_created_at ON run_results USING BRIN (created_at);
	ALTER TABLE run_results ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE run_results ADD COLUMN updated_at timestamptz;
	UPDATE run_results SET updated_at = '2019-01-01';
	CREATE INDEX idx_run_results_updated_at ON run_results USING BRIN (updated_at);
	ALTER TABLE run_results ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE service_agreements ADD COLUMN updated_at timestamptz;
	UPDATE service_agreements SET updated_at = '2019-01-01';
	CREATE INDEX idx_service_agreements_updated_at ON service_agreements USING BRIN (updated_at);
	ALTER TABLE service_agreements ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE task_runs ADD COLUMN updated_at timestamptz;
	UPDATE task_runs SET updated_at = '2019-01-01';
	CREATE INDEX idx_task_runs_updated_at ON task_runs USING BRIN (updated_at);
	ALTER TABLE task_runs ALTER COLUMN updated_at SET NOT NULL;
	
	ALTER TABLE tx_attempts ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE tx_attempts SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	CREATE INDEX idx_tx_attempts_updated_at ON tx_attempts USING BRIN (updated_at);
	ALTER TABLE tx_attempts ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE txes ADD COLUMN IF NOT EXISTS created_at timestamptz;
	UPDATE txes SET created_at = '2019-01-01' WHERE created_at IS NULL;
	CREATE INDEX idx_txes_created_at ON txes USING BRIN (created_at);
	ALTER TABLE txes ALTER COLUMN created_at SET NOT NULL;

	ALTER TABLE txes ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE txes SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	CREATE INDEX idx_txes_updated_at ON txes USING BRIN (updated_at);
	ALTER TABLE txes ALTER COLUMN updated_at SET NOT NULL;

	ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at timestamptz;
	UPDATE users SET updated_at = '2019-01-01' WHERE updated_at IS NULL;
	CREATE INDEX idx_users_updated_at ON users USING BRIN (updated_at);
	ALTER TABLE users ALTER COLUMN updated_at SET NOT NULL;
	`).Error
}
