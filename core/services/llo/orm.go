package llo

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

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
	// TODO: make sure to scope by chain ID everywhere
}

func (o *orm) LoadChannelDefinitions(ctx context.Context, addr common.Address) (cd commontypes.ChannelDefinitions, blockNum int64, err error) {
	type scd struct {
		definitions []byte
		blockNum    int64
	}
	var scanned scd
	if err = o.q.GetContext(ctx, scanned, "SELECT definitions, block_num FROM streams_channel_definitions WHERE evm_chain_id = $1 AND addr = $2", o.evmChainID.String(), addr); err != nil {
		return nil, 0, err
	}

	if err = json.Unmarshal(scanned.definitions, &cd); err != nil {
		return nil, 0, err
	}

	return cd, scanned.blockNum, nil
}

func (o *orm) StoreChannelDefinitions(ctx context.Context, cd commontypes.ChannelDefinitions) error {
	panic("TODO")
}
