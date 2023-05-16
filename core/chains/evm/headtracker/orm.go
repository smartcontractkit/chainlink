package headtracker

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Chain Agnostic HeadStore
type HeadStore[H commontypes.Head[BLOCK_HASH], BLOCK_HASH commontypes.Hashable] interface {
	// Insert Head into the DB
	IdempotentInsertHead(ctx context.Context, head H) (err error)
	// Delete heads from the DB beyond n blocks
	TrimOldHeads(ctx context.Context, n uint) (err error)
	// Get the latest head from DB
	LatestHead(ctx context.Context) (head H, err error)
	// LatestHeads returns the latest heads up to given limit
	LatestHeads(ctx context.Context, limit uint) (heads []*evmtypes.Head, err error)
	// Find head by Hash from DB
	HeadByHash(ctx context.Context, hash BLOCK_HASH) (head H, err error)
}

type orm struct {
	q       pg.Q
	chainID utils.Big
}

var _ HeadStore[*evmtypes.Head, common.Hash] = (*orm)(nil)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, chainID big.Int) HeadStore[*evmtypes.Head, common.Hash] {
	return &orm{pg.NewQ(db, lggr.Named("HeadTrackerORM"), cfg), utils.Big(chainID)}
}

func (orm *orm) IdempotentInsertHead(ctx context.Context, head *evmtypes.Head) error {
	// listener guarantees head.EVMChainID to be equal to orm.chainID
	q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	query := `
	INSERT INTO evm_heads (hash, number, parent_hash, created_at, timestamp, l1_block_number, evm_chain_id, base_fee_per_gas) VALUES (
	:hash, :number, :parent_hash, :created_at, :timestamp, :l1_block_number, :evm_chain_id, :base_fee_per_gas)
	ON CONFLICT (evm_chain_id, hash) DO NOTHING`
	err := q.ExecQNamed(query, head)
	return errors.Wrap(err, "IdempotentInsertHead failed to insert head")
}

func (orm *orm) TrimOldHeads(ctx context.Context, n uint) (err error) {
	q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	return q.ExecQ(`
	DELETE FROM evm_heads
	WHERE evm_chain_id = $1 AND number < (
		SELECT min(number) FROM (
			SELECT number
			FROM evm_heads
			WHERE evm_chain_id = $1
			ORDER BY number DESC
			LIMIT $2
		) numbers
	)`, orm.chainID, n)
}

func (orm *orm) LatestHead(ctx context.Context) (head *evmtypes.Head, err error) {
	head = new(evmtypes.Head)
	q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	err = q.Get(head, `SELECT * FROM evm_heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT 1`, orm.chainID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	err = errors.Wrap(err, "LatestHead failed")
	return
}

func (orm *orm) LatestHeads(ctx context.Context, limit uint) (heads []*evmtypes.Head, err error) {
	q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	err = q.Select(&heads, `SELECT * FROM evm_heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT $2`, orm.chainID, limit)
	err = errors.Wrap(err, "LatestHeads failed")
	return
}

func (orm *orm) HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error) {
	q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	head = new(evmtypes.Head)
	err = q.Get(head, `SELECT * FROM evm_heads WHERE evm_chain_id = $1 AND hash = $2`, orm.chainID, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return head, err
}
