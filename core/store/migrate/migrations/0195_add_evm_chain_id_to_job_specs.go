package migrations

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func init() {
	goose.AddMigrationContext(Up195, Down195)
}

const (
	setMissingEvmChainIDsInSpecs = `
	UPDATE direct_request_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE flux_monitor_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE ocr_oracle_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE keeper_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE vrf_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE blockhash_store_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE block_header_feeder_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	`
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
	var opts chainlink.GeneralConfigOpts
	cfg, err := opts.New()
	if err != nil {
		return err
	}

	if cfg.EVMEnabled() {
		chainID := cfg.EVMConfigs()[0].ChainID
		_, err = tx.ExecContext(ctx, setMissingEvmChainIDsInSpecs, chainID)
		if err != nil {
			return errors.Wrap(err, "failed to set missing evm chain ids")
		}
	}

	_, err = tx.ExecContext(ctx, addNullConstraintsToSpecs)
	return errors.Wrap(err, "failed to add null constraints")
}

// nolint
func Down195(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, dropNullConstraintsFromSpecs); err != nil {
		return err
	}
	return nil
}
