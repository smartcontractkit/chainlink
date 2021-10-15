package headtracker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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

	if ctx.Err() != nil {
		return nil
	} else if err != nil && err.Error() == "sql: no rows in result set" {
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

// LatestHead returns the highest seen head
func (orm *ORM) LatestHead(ctx context.Context) (head *eth.Head, err error) {
	head = new(eth.Head)
	err = orm.db.WithContext(ctx).Where("evm_chain_id = ?", orm.chainID).Order("number DESC, created_at DESC, id DESC").First(&head).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	err = errors.Wrap(err, "LatestHead failed")
	return
}

// LatestHeads returns the latest heads up to given limit
func (orm *ORM) LatestHeads(ctx context.Context, limit int) (heads []*eth.Head, err error) {
	err = orm.db.WithContext(ctx).Where("evm_chain_id = ?", orm.chainID).Order("number DESC, created_at DESC, id DESC").Limit(limit).Find(&heads).Error
	err = errors.Wrap(err, "LatestHeads failed")
	return
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
