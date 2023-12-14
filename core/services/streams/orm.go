package streams

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type ORM interface {
	StreamCacheORM
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

func (o *orm) LoadStreams(ctx context.Context, lggr logger.Logger, runner Runner, m map[commontypes.StreamID]Stream) error {
	rows, err := o.q.QueryContext(ctx, "SELECT s.id, ps.id, ps.dot_dag_source, ps.max_task_duration FROM streams s JOIN pipeline_specs ps ON ps.id = s.pipeline_spec_id")
	if err != nil {
		// TODO: retries?
		return err
	}

	for rows.Next() {
		var strm stream
		if err := rows.Scan(&strm.id, &strm.spec.ID, &strm.spec.DotDagSource, &strm.spec.MaxTaskDuration); err != nil {
			return err
		}
		strm.lggr = lggr.Named("Stream").With("streamID", strm.id)
		strm.runner = runner

		m[strm.id] = &strm
	}
	return rows.Err()
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
