package llo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/channeldefinitions"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/types"
)

type ChainScopedORM interface {
	channeldefinitions.ChannelDefinitionCacheORM
}

var _ ChainScopedORM = &chainScopedORM{}

type chainScopedORM struct {
	ds            sqlutil.DataSource
	chainSelector uint64
}

func NewChainScopedORM(ds sqlutil.DataSource, chainSelector uint64) ChainScopedORM {
	return &chainScopedORM{ds, chainSelector}
}

func (o *chainScopedORM) LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *types.PersistedDefinitions, err error) {
	pd = new(types.PersistedDefinitions)
	err = o.ds.GetContext(ctx, pd, "SELECT * FROM channel_definitions WHERE chain_selector = $1 AND addr = $2 AND don_id = $3", o.chainSelector, addr, donID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to LoadChannelDefinitions; %w", err)
	}

	return pd, nil
}

// StoreChannelDefinitions will store a ChannelDefinitions list for a given chain_selector, addr, don_id
// It updates if the new version is greater than the existing record OR if the block number has changed
// (indicating definitions were updated even if version hasn't progressed)
func (o *chainScopedORM) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns json.RawMessage, blockNum int64, format uint32) error {
	_, err := o.ds.ExecContext(ctx, `
INSERT INTO channel_definitions (chain_selector, addr, don_id, definitions, block_num, version, updated_at, format)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
ON CONFLICT (chain_selector, addr, don_id) DO UPDATE
SET definitions = $4, block_num = $5, version = $6, updated_at = NOW(), format = $7
WHERE EXCLUDED.don_id = channel_definitions.don_id AND EXCLUDED.chain_selector = channel_definitions.chain_selector
AND (EXCLUDED.version >= channel_definitions.version OR EXCLUDED.block_num >= channel_definitions.block_num)`,
		o.chainSelector, addr, donID, dfns, blockNum, version, format)
	if err != nil {
		return fmt.Errorf("StoreChannelDefinitions failed: %w", err)
	}
	return nil
}

func (o *chainScopedORM) CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) error {
	_, err := o.ds.ExecContext(ctx, "DELETE FROM channel_definitions WHERE chain_selector = $1 AND addr = $2 AND don_id = $3", o.chainSelector, addr, donID)
	if err != nil {
		return fmt.Errorf("failed to CleanupChannelDefinitions; %w", err)
	}
	return nil
}
