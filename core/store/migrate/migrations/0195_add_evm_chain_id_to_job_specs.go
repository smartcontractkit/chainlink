package migrations

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up195, Down195)
}

const (
	addNullConstraintsToSpecs = `
	ALTER TABLE direct_request_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE flux_monitor_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE ocr_oracle_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE keeper_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE vrf_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE blockhash_store_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE block_header_feeder_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	`

	dropNullConstraintsFromSpecs = `
	ALTER TABLE direct_request_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE flux_monitor_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE ocr_oracle_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE keeper_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE vrf_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE blockhash_store_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
	ALTER TABLE block_header_feeder_specs ALTER COLUMN evm_chain_id DROP NOT NULL;
    `
)

// nolint
func Up195(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, addNullConstraintsToSpecs)
	return errors.Wrap(err, "failed to add null constraints")
}

// nolint
func Down195(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, dropNullConstraintsFromSpecs); err != nil {
		return err
	}
	return nil
}
