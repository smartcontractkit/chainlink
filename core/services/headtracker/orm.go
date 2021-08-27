package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ORM struct {
	db *gorm.DB
}

func NewORM(db *gorm.DB) *ORM {
	return &ORM{db}
}

// IdempotentInsertHead inserts a head only if the hash is new. Will do nothing if hash exists already.
// No advisory lock required because this is thread safe.
func (orm *ORM) IdempotentInsertHead(ctx context.Context, h models.Head) error {
	err := orm.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "hash"}},
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
	WHERE number < (
		SELECT min(number) FROM (
			SELECT number
			FROM heads
			ORDER BY number DESC
			LIMIT ?
		) numbers
	)`, n).Error
}

// Chain return the chain of heads starting at hash and up to lookback parents
// Returns RecordNotFound if no head with the given hash exists
func (orm *ORM) Chain(ctx context.Context, hash common.Hash, lookback uint) (models.Head, error) {
	rows, err := orm.db.WithContext(ctx).Raw(`
	WITH RECURSIVE chain AS (
		SELECT * FROM heads WHERE hash = ?
	UNION
		SELECT h.* FROM heads h
		JOIN chain ON chain.parent_hash = h.hash
	) SELECT id, hash, number, parent_hash, timestamp, created_at FROM chain LIMIT ?
	`, hash, lookback).Rows()
	if err != nil {
		return models.Head{}, err
	}
	defer logger.ErrorIfCalling(rows.Close)
	var firstHead *models.Head
	var prevHead *models.Head
	for rows.Next() {
		h := models.Head{}
		if err = rows.Scan(&h.ID, &h.Hash, &h.Number, &h.ParentHash, &h.Timestamp, &h.CreatedAt); err != nil {
			return models.Head{}, err
		}
		if firstHead == nil {
			firstHead = &h
		} else {
			prevHead.Parent = &h
		}
		prevHead = &h
	}
	if err = rows.Err(); err != nil {
		return models.Head{}, err
	}
	if firstHead == nil {
		return models.Head{}, gorm.ErrRecordNotFound
	}
	return *firstHead, nil
}

// LastHead returns the head with the highest number. In the case of ties (e.g.
// due to re-org) it returns the most recently seen head entry.
func (orm *ORM) LastHead(ctx context.Context) (*models.Head, error) {
	number := &models.Head{}
	err := orm.db.WithContext(ctx).Order("number DESC, created_at DESC, id DESC").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// HeadByHash fetches the head with the given hash from the db, returns nil if none exists
func (orm *ORM) HeadByHash(ctx context.Context, hash common.Hash) (*models.Head, error) {
	head := &models.Head{}
	err := orm.db.WithContext(ctx).Where("hash = ?", hash).First(head).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return head, err
}
