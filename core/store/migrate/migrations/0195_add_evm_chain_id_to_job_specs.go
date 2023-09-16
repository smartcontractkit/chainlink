package migrations

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pressly/goose/v3"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func init() {
	goose.AddMigrationContext(Up195, Down195)
}

const (
	up195_1 = `
	UPDATE direct_request_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE flux_monitor_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE ocr_oracle_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE keeper_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE vrf_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE blockhash_store_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	UPDATE block_header_feeder_specs SET evm_chain_id = '%d' WHERE evm_chain_id IS NULL;
	`
	up195_2 = `
	ALTER TABLE direct_request_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE flux_monitor_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE ocr_oracle_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE keeper_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE vrf_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE blockhash_store_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	ALTER TABLE block_header_feeder_specs ALTER COLUMN evm_chain_id SET NOT NULL;
	`

	down195 = `
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

	url := cfg.Database().URL()
	if cfg.EVMEnabled() && !strings.Contains(url.String(), "_test") {
		chainID := cfg.EVMConfigs()[0].ChainID
		_, err = tx.ExecContext(ctx, up195_1, chainID)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, up195_2)
	return err
}

// nolint
func Down195(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, down195); err != nil {
		return err
	}
	return nil
}
