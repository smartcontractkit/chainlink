package llo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
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
// It only updates if the new version is greater than the existing record
func (o *chainScopedORM) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) error {
	_, err := o.ds.ExecContext(ctx, `
INSERT INTO channel_definitions (chain_selector, addr, don_id, definitions, block_num, version, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (chain_selector, addr, don_id) DO UPDATE
SET definitions = $4, block_num = $5, version = $6, updated_at = NOW()
WHERE EXCLUDED.version > channel_definitions.version
`, o.chainSelector, addr, donID, dfns, blockNum, version)
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
