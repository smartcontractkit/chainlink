package headtracker

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

type ORM struct {
	db      *sqlx.DB
	chainID utils.Big
}

func NewORM(db *sqlx.DB, chainID big.Int) *ORM {
	if db == nil {
		panic("db may not be nil")
	}
	return &ORM{db, utils.Big(chainID)}
}

// IdempotentInsertHead inserts a head only if the hash is new. Will do nothing if hash exists already.
// No advisory lock required because this is thread safe.
func (orm *ORM) IdempotentInsertHead(ctx context.Context, h eth.Head) error {
	if h.EVMChainID == nil {
		h.EVMChainID = &orm.chainID
	} else if ((*big.Int)(h.EVMChainID)).Cmp((*big.Int)(&orm.chainID)) != 0 {
		return errors.Errorf("head chain ID %s does not match orm chain ID %s", h.EVMChainID.String(), orm.chainID.String())
	}
	q := postgres.NewQ(orm.db, postgres.WithParentCtx(ctx))
	query := `
INSERT INTO heads (hash, number, parent_hash, created_at, timestamp, l1_block_number, evm_chain_id, base_fee_per_gas) VALUES (
:hash, :number, :parent_hash, :created_at, :timestamp, :l1_block_number, :evm_chain_id, :base_fee_per_gas)
ON CONFLICT (evm_chain_id, hash) DO NOTHING
`
	err := q.ExecQNamed(query, h)
	return errors.Wrap(err, "IdempotentInsertHead failed to insert head")
}

// TrimOldHeads deletes heads such that only the top N block numbers remain
func (orm *ORM) TrimOldHeads(ctx context.Context, n uint) (err error) {
	return postgres.NewQ(orm.db, postgres.WithParentCtx(ctx)).ExecQ(`
	DELETE FROM heads
	WHERE evm_chain_id = $1 AND number < (
		SELECT min(number) FROM (
			SELECT number
			FROM heads
			WHERE evm_chain_id = $2
			ORDER BY number DESC
			LIMIT $3
		) numbers
	)`, orm.chainID, orm.chainID, n)
}

// LatestHead returns the highest seen head
func (orm *ORM) LatestHead(ctx context.Context) (head *eth.Head, err error) {
	head = new(eth.Head)
	err = postgres.NewQ(orm.db, postgres.WithParentCtx(ctx)).
		Get(head, `SELECT * FROM heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT 1`, orm.chainID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	err = errors.Wrap(err, "LatestHead failed")
	return
}

// LatestHeads returns the latest heads up to given limit
func (orm *ORM) LatestHeads(ctx context.Context, limit int) (heads []*eth.Head, err error) {
	err = postgres.NewQ(orm.db, postgres.WithParentCtx(ctx)).
		Select(&heads, `SELECT * FROM heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT $2`, orm.chainID, limit)
	err = errors.Wrap(err, "LatestHeads failed")
	return
}

// HeadByHash fetches the head with the given hash from the db, returns nil if none exists
func (orm *ORM) HeadByHash(ctx context.Context, hash common.Hash) (head *eth.Head, err error) {
	head = new(eth.Head)
	err = postgres.NewQ(orm.db, postgres.WithParentCtx(ctx)).Get(head, `SELECT * FROM heads WHERE evm_chain_id = $1 AND hash = $2`, orm.chainID, hash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return head, err
}
