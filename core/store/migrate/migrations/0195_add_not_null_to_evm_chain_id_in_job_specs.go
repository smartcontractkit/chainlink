package migrations

import (
	"context"
	"database/sql"
	"os"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
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
	chainID, set := os.LookupEnv(env.EVMChainIDNotNullMigration0195)
	if set {
		updateQueries := []string{
			`UPDATE direct_request_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE flux_monitor_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE ocr_oracle_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE keeper_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE vrf_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE blockhash_store_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
			`UPDATE block_header_feeder_specs SET evm_chain_id = $1 WHERE evm_chain_id IS NULL;`,
		}
		for i := range updateQueries {
			_, err := tx.Exec(updateQueries[i], chainID)
			if err != nil {
				return errors.Wrap(err, "failed to set missing evm chain ids")
			}
		}
	}

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
