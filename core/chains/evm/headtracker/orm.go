package headtracker

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

type ORM interface {
	// IdempotentInsertHead inserts a head only if the hash is new. Will do nothing if hash exists already.
	// No advisory lock required because this is thread safe.
	IdempotentInsertHead(ctx context.Context, head *evmtypes.Head) error
	// TrimOldHeads deletes heads such that only blocks >= minBlockNumber remain
	TrimOldHeads(ctx context.Context, minBlockNumber int64) (err error)
	// LatestHead returns the highest seen head
	LatestHead(ctx context.Context) (head *evmtypes.Head, err error)
	// LatestHeads returns the latest heads with blockNumbers >= minBlockNumber
	LatestHeads(ctx context.Context, minBlockNumber int64) (heads []*evmtypes.Head, err error)
	// HeadByHash fetches the head with the given hash from the db, returns nil if none exists
	HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error)
}

var _ ORM = &DbORM{}

type DbORM struct {
	chainID ubig.Big
	ds      sqlutil.DataSource
}

// NewORM creates an ORM scoped to chainID.
func NewORM(chainID big.Int, ds sqlutil.DataSource) *DbORM {
	return &DbORM{
		chainID: ubig.Big(chainID),
		ds:      ds,
	}
}

func (orm *DbORM) IdempotentInsertHead(ctx context.Context, head *evmtypes.Head) error {
	// listener guarantees head.EVMChainID to be equal to DbORM.chainID
	query := `
	INSERT INTO evm.heads (hash, number, parent_hash, created_at, timestamp, l1_block_number, evm_chain_id, base_fee_per_gas) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (evm_chain_id, hash) DO NOTHING`
	_, err := orm.ds.ExecContext(ctx, query, head.Hash, head.Number, head.ParentHash, head.CreatedAt, head.Timestamp, head.L1BlockNumber, orm.chainID, head.BaseFeePerGas)
	return pkgerrors.Wrap(err, "IdempotentInsertHead failed to insert head")
}

func (orm *DbORM) TrimOldHeads(ctx context.Context, minBlockNumber int64) (err error) {
	query := `DELETE FROM evm.heads WHERE evm_chain_id = $1 AND number < $2`
	_, err = orm.ds.ExecContext(ctx, query, orm.chainID, minBlockNumber)
	return err
}

func (orm *DbORM) LatestHead(ctx context.Context) (head *evmtypes.Head, err error) {
	head = new(evmtypes.Head)
	err = orm.ds.GetContext(ctx, head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 ORDER BY number DESC, created_at DESC, id DESC LIMIT 1`, orm.chainID)
	if pkgerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	err = pkgerrors.Wrap(err, "LatestHead failed")
	return
}

func (orm *DbORM) LatestHeads(ctx context.Context, minBlockNumer int64) (heads []*evmtypes.Head, err error) {
	err = orm.ds.SelectContext(ctx, &heads, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 AND number >= $2 ORDER BY number DESC, created_at DESC, id DESC`, orm.chainID, minBlockNumer)
	err = pkgerrors.Wrap(err, "LatestHeads failed")
	return
}

func (orm *DbORM) HeadByHash(ctx context.Context, hash common.Hash) (head *evmtypes.Head, err error) {
	head = new(evmtypes.Head)
	err = orm.ds.GetContext(ctx, head, `SELECT * FROM evm.heads WHERE evm_chain_id = $1 AND hash = $2`, orm.chainID, hash)
	if pkgerrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return head, err
}
