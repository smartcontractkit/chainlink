package llo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type ORM interface {
	ChannelDefinitionCacheORM
}

type PersistedDefinitions struct {
	ChainSelector uint64                      `db:"chain_selector"`
	Address       common.Address              `db:"addr"`
	Definitions   llotypes.ChannelDefinitions `db:"definitions"`
	// The block number in which the log for this definitions was emitted
	BlockNum  int64     `db:"block_num"`
	DonID     uint32    `db:"don_id"`
	Version   uint32    `db:"version"`
	UpdatedAt time.Time `db:"updated_at"`
}

var _ ORM = &orm{}

type orm struct {
	ds            sqlutil.DataSource
	chainSelector uint64
}

func NewORM(ds sqlutil.DataSource, chainSelector uint64) ORM {
	return &orm{ds, chainSelector}
}

func (o *orm) LoadChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) (pd *PersistedDefinitions, err error) {
	pd = new(PersistedDefinitions)
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
func (o *orm) StoreChannelDefinitions(ctx context.Context, addr common.Address, donID, version uint32, dfns llotypes.ChannelDefinitions, blockNum int64) error {
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

func (o *orm) CleanupChannelDefinitions(ctx context.Context, addr common.Address, donID uint32) error {
	_, err := o.ds.ExecContext(ctx, "DELETE FROM channel_definitions WHERE chain_selector = $1 AND addr = $2 AND don_id = $3", o.chainSelector, addr, donID)
	if err != nil {
		return fmt.Errorf("failed to CleanupChannelDefinitions; %w", err)
	}
	return nil
}
