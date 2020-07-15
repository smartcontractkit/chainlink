package migration1594815494

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the oracle_requests table which will eventually supercede run_requests
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE oracle_requests (
			log_consumption_id bigint REFERENCES log_consumptions (id) NOT NULL,
			spec_id uuid REFERENCES job_specs (id) NOT NULL,
			requester bytea NOT NULL CHECK octet_length(requester) = 20,
			request_id bytea PRIMARY KEY CHECK octet_length(request_id) = 32,
			payment numeric(78,0) NOT NULL,
			callback_addr bytea NOT NULL CHECK octet_length(callback_addr) = 20,
			callback_function_id bytea NOT NULL CHECK octet_length(callback_addr) = 4,
			cancel_expiration timestamptz NOT NULL,
			data_version numeric(78,0) NOT NULL,
			data bytea NOT NULL,
			created_at timestamptz NOT NULL
		)
	`).Error
}
