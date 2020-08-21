package migration1597695690

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the offchain_reporting_job_specs table
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`
        CREATE TABLE offchain_reporting_job_specs (
            id uuid PRIMARY KEY,
            contract_address bytea NOT NULL,
            p2p_node_id text NOT NULL,
            p2p_bootstrap_nodes jsonb NOT NULL,
            key_bundle text NOT NULL,
            monitoring_endpoint text NOT NULL,
            node_address bytea NOT NULL,
            observation_timeout integer NOT NULL,
            observation_source jsonb NOT NULL
        );
    `).Error
	if err != nil {
		return err
	}

	err = tx.Exec(`
        CREATE TABLE offchain_reporting_persistent_states (
            id SERIAL PRIMARY KEY,
            job_spec_id uuid NOT NULL,
            group_id bytea NOT NULL,
            epoch integer NOT NULL,
            highest_sent_epoch integer NOT NULL,
            highest_received_epoch integer[31] NOT NULL
        );
        ALTER TABLE offchain_reporting_persistent_states ADD CONSTRAINT "ocr_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES offchain_reporting_job_specs ("id") ON DELETE CASCADE;
        ALTER TABLE offchain_reporting_persistent_states ADD CONSTRAINT chk_group_id_length CHECK (
            octet_length(group_id) = 32
        );
        CREATE UNIQUE INDEX ocr_persistent_states_unique_idx ON offchain_reporting_persistent_states ("job_id", "group_id");
    `).Error
	if err != nil {
		return err
	}

	err = tx.Exec(`
        CREATE TABLE offchain_reporting_configs (
            id SERIAL PRIMARY KEY,
            job_spec_id uuid NOT NULL,
            group_id bytea NOT NULL,
            oracles jsonb NOT NULL,
            secret bytea NOT NULL,
            f integer NOT NULL,
            delta_progress integer NOT NULL,
            delta_resend integer NOT NULL,
            delta_round integer NOT NULL,
            delta_observe integer NOT NULL,
            delta_c integer NOT NULL,
            alpha float NOT NULL,
            r_max integer NOT NULL,
            delta_stage integer NOT NULL,
            schedule integer[] NOT NULL
        );
        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT "ocr_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES offchain_reporting_job_specs ("id") ON DELETE CASCADE;
        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT chk_group_id_length CHECK (
            octet_length(group_id) = 27
        );
        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT chk_secret_length CHECK (
            octet_length(secret) = 16
        );
        CREATE UNIQUE INDEX ocr_configs_unique_idx ON offchain_reporting_configs ("job_id", "group_id");
   `).Error
	if err != nil {
		return err
	}
	return nil
}
