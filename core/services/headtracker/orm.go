package headtracker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ORM struct {
	db      *gorm.DB
	chainID utils.Big
}

func NewORM(db *gorm.DB, chainID big.Int) *ORM {
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
	err := orm.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "evm_chain_id"}, {Name: "hash"}},
			DoNothing: true,
		}).Create(&h).Error

	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil
	}
	return err
}

// TrimOldHeads deletes heads such that only the top N block numbers remain
func (orm *ORM) TrimOldHeads(ctx context.Context, n uint) (err error) {
	return orm.db.WithContext(ctx).Exec(`
	DELETE FROM heads
	WHERE evm_chain_id = ? AND number < (
		SELECT min(number) FROM (
			SELECT number
			FROM heads
			WHERE evm_chain_id = ?
			ORDER BY number DESC
			LIMIT ?
		) numbers
	)`, orm.chainID, orm.chainID, n).Error
}

// Chain return the chain of heads starting at hash and up to lookback parents
// Returns RecordNotFound if no head with the given hash exists
func (orm *ORM) Chain(ctx context.Context, hash common.Hash, lookback uint) (eth.Head, error) {
	rows, err := orm.db.WithContext(ctx).Raw(`
	WITH RECURSIVE chain AS (
		SELECT * FROM heads WHERE evm_chain_id = ? AND hash = ?
	UNION
		SELECT h.* FROM heads h
		JOIN chain ON chain.parent_hash = h.hash
	) SELECT id, hash, number, parent_hash, timestamp, created_at, l1_block_number, evm_chain_id FROM chain LIMIT ?
	`, orm.chainID, hash, lookback).Rows()
	if err != nil {
		return eth.Head{}, err
	}
	defer logger.ErrorIfCalling(rows.Close)
	var firstHead *eth.Head
	var prevHead *eth.Head
	for rows.Next() {
		h := eth.Head{}
		if err = rows.Scan(&h.ID, &h.Hash, &h.Number, &h.ParentHash, &h.Timestamp, &h.CreatedAt, &h.L1BlockNumber, &h.EVMChainID); err != nil {
			return eth.Head{}, err
		}
		if firstHead == nil {
			firstHead = &h
		} else {
			prevHead.Parent = &h
		}
		prevHead = &h
	}
	if err = rows.Err(); err != nil {
		return eth.Head{}, err
	}
	if firstHead == nil {
		return eth.Head{}, gorm.ErrRecordNotFound
	}
	return *firstHead, nil
}

// LastHead returns the head with the highest number. In the case of ties (e.g.
// due to re-org) it returns the most recently seen head entry.
func (orm *ORM) LastHead(ctx context.Context) (*eth.Head, error) {
	number := &eth.Head{}
	err := orm.db.WithContext(ctx).Where("evm_chain_id = ?", orm.chainID).Order("number DESC, created_at DESC, id DESC").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// HeadByHash fetches the head with the given hash from the db, returns nil if none exists
func (orm *ORM) HeadByHash(ctx context.Context, hash common.Hash) (*eth.Head, error) {
	head := &eth.Head{}
	err := orm.db.WithContext(ctx).Where("evm_chain_id = ? AND hash = ?", orm.chainID, hash).First(head).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return head, err
}
