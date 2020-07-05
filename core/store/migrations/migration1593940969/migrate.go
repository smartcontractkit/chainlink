package migration1591603775

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds the events_oracle_request table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE events_oracle_request (
			spec_id uuid NOT NULL REFERENCES job_specs (id),
			requester bytea NOT NULL,
			request_id bytea PRIMARY KEY,
			payment numeric(78, 0) NOT NULL,
			callback_addr bytea NOT NULL,
			callback_function_id bytea NOT NULL,
			data_version  numeric(78, 0) NOT NULL,
			data bytea NOT NULL,
			created_at timestamptz NOT NULL
		);

		CREATE INDEX events_oracle_request_spec_id ON events_oracle_request (spec_id);

		ALTER TABLE events_oracle_request_request_id ADD CONSTRAINT chk_request_id_len CHECK (
			octet_length(request_id) = 32
		);

		ALTER TABLE events_oracle_request_requester ADD CONSTRAINT chk_requester_len CHECK (
			octet_length(requester) = 20
		);

		ALTER TABLE events_oracle_request_callback_addr ADD CONSTRAINT chk_callback_addr_len CHECK (
			octet_length(callback_addr) = 20
		);

		ALTER TABLE events_oracle_request_callback_function_id ADD CONSTRAINT chk_callback_function_id CHECK (
			octet_length(callback_function_id) = 4
		);

		CREATE INDEX idx_events_oracle_request_created_at ON events_oracle_request USING BRIN (created_at);
    `).Error
}
