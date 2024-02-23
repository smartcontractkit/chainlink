package llo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type ORM interface {
	ChannelDefinitionCacheORM
}

var _ ORM = &orm{}

type orm struct {
	q          pg.Queryer
	evmChainID *big.Int
}

func NewORM(q pg.Queryer, evmChainID *big.Int) ORM {
	return &orm{q, evmChainID}
}

func (o *orm) LoadChannelDefinitions(ctx context.Context, addr common.Address) (dfns llotypes.ChannelDefinitions, blockNum int64, err error) {
	type scd struct {
		Definitions []byte `db:"definitions"`
		BlockNum    int64  `db:"block_num"`
	}
	var scanned scd
	err = o.q.GetContext(ctx, &scanned, "SELECT definitions, block_num FROM channel_definitions WHERE evm_chain_id = $1 AND addr = $2", o.evmChainID.String(), addr)
	if errors.Is(err, sql.ErrNoRows) {
		return dfns, blockNum, nil
	} else if err != nil {
		return nil, 0, fmt.Errorf("failed to LoadChannelDefinitions; %w", err)
	}

	if err = json.Unmarshal(scanned.Definitions, &dfns); err != nil {
		return nil, 0, fmt.Errorf("failed to LoadChannelDefinitions; JSON Unmarshal failure; %w", err)
	}

	return dfns, scanned.BlockNum, nil
}

// TODO: Test this method
// https://smartcontract-it.atlassian.net/jira/software/c/projects/MERC/issues/MERC-3653
func (o *orm) StoreChannelDefinitions(ctx context.Context, addr common.Address, dfns llotypes.ChannelDefinitions, blockNum int64) error {
	_, err := o.q.ExecContext(ctx, `
INSERT INTO channel_definitions (evm_chain_id, addr, definitions, block_num, updated_at)
VALUES ($1, $2, $3, $4, NOW())
ON CONFLICT (evm_chain_id, addr) DO UPDATE
SET definitions = $3, block_num = $4, updated_at = NOW()
`, o.evmChainID.String(), addr, dfns, blockNum)
	if err != nil {
		return fmt.Errorf("StoreChannelDefinitions failed: %w", err)
	}
	return nil
}
