package headtracker

import (
	"context"
	"database/sql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type ORM interface {
	// IdempotentInsertHead inserts a head only if the hash is new. Will do nothing if hash exists already.
	// No advisory lock required because this is thread safe.
	IdempotentInsertHead(ctx context.Context, head *evmtypes.Head) error
	// TrimOldHeads deletes heads such that only the top N block numbers remain
	TrimOldHeads(ctx context.Context, n uint) (err error)
	// LatestHead returns the highest seen head
	LatestHead(ctx context.Context) (head *evmtypes.Head, err error)
	// LatestHeads returns the latest heads up to given limit
	LatestHeads(ctx context.Context, limit uint) (heads []*evmtypes.Head, err error)
	// HeadByHash fetches the head with the given hash from the db, returns nil if none exists
	HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error)
}

var _ ORM = &orm{}

type orm struct {
	//q       pg.Q
	chainID ubig.Big
	db      sqlutil.Queryer
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID big.Int, db sqlutil.Queryer) ORM {
	return &orm{
		chainID: ubig.Big(chainID),
		db:      db,
	}
}

/*
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, chainID big.Int) ORM {
	return &orm{pg.NewQ(db, logger.Named(lggr, "HeadTrackerORM"), cfg), ubig.Big(chainID)}
}
*/

func (orm *orm) IdempotentInsertHead(ctx context.Context, head *evmtypes.Head) error {
	// listener guarantees head.EVMChainID to be equal to orm.chainID
	//q := orm.q.WithOpts(pg.WithParentCtx(ctx))

	query := `
	INSERT INTO evm.heads (hash, number, parent_hash, created_at, timestamp, l1_block_number, evm_chain_id, base_fee_per_gas) VALUES (
	:hash, :number, :parent_hash, :created_at, :timestamp, :l1_block_number, :evm_chain_id, :base_fee_per_gas)
	ON CONFLICT (evm_chain_id, hash) DO NOTHING`

	//err := q.ExecQNamed(query, head)

	_, err := orm.db.ExecContext(ctx, query, head)

	return errors.Wrap(err, "IdempotentInsertHead failed to insert head")
}

func (orm *orm) TrimOldHeads(ctx context.Context, n uint) (err error) {
	//q := orm.q.WithOpts(pg.WithParentCtx(ctx))

	_, err = orm.db.ExecContext(ctx, `
	DELETE FROM evm.heads
	WHERE evm_chain_id = $1 AND number < (
		SELECT min(number) FROM (
			SELECT number
			FROM evm.heads
			WHERE evm_chain_id = $1
			ORDER BY number DESC
			LIMIT $2
		) numbers
	)`, orm.chainID, n)

	return err

	/*
		return q.ExecQ(`
		DELETE FROM evm.heads
		WHERE evm_chain_id = $1 AND number < (
			SELECT min(number) FROM (
				SELECT number
				FROM evm.heads
				WHERE evm_chain_id = $1
				ORDER BY number DESC
				LIMIT $2
			) numbers
		)`, orm.chainID, n)
	*/
}

func (orm *orm) LatestHead(ctx context.Context) (head *evmtypes.Head, err error) {
	head = new(evmtypes.Head)
	//q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	//err = q.Get(head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT 1`, orm.chainID)

	if err := orm.db.SelectContext(ctx, head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT 1`, orm.chainID); err != nil {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	err = errors.Wrap(err, "LatestHead failed")
	return
}

func (orm *orm) LatestHeads(ctx context.Context, limit uint) (heads []*evmtypes.Head, err error) {
	//q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	//err = q.Select(&heads, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT $2`, orm.chainID, limit)

	err = orm.db.SelectContext(ctx, &heads, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT $2`, orm.chainID, limit)
	err = errors.Wrap(err, "LatestHeads failed")
	return
}

func (orm *orm) HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error) {
	//q := orm.q.WithOpts(pg.WithParentCtx(ctx))
	head = new(evmtypes.Head)
	//err = q.Get(head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 AND hash = $2`, orm.chainID, hash)
	err = orm.db.SelectContext(ctx, head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 AND hash = $2`, orm.chainID, hash)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return head, err
}
