package headtracker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type ORM interface {
	SaveBlock(ctx context.Context, block *Block) error
	GetBlock(ctx context.Context, hash common.Hash) (*Block, error)
}

type orm struct {
	db *gorm.DB
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB) *orm {
	return &orm{
		db: db,
	}
}

func (orm *orm) Close() error {
	return nil
}

func (orm *orm) SaveBlock(ctx context.Context, block *Block) error {
	trans, err := json.Marshal(block.Transactions)
	if err != nil {
		return errors.Wrap(err, "Failed to convert transactions to json")
	}
	return orm.db.Exec(`
		INSERT INTO blocks (
				hash, number, parent_hash, transactions, created_at, "timestamp"
		) VALUES (
				?, ?, ?, ?, ?, ?
		) ON CONFLICT (hash) DO NOTHING
    `, block.Hash, block.Number, block.ParentHash, trans, time.Now(), time.Now()).Error
}

func (orm *orm) GetBlock(ctx context.Context, hash common.Hash) (*Block, error) {
	rows, err := orm.db.WithContext(ctx).Raw(`
	SELECT hash, number, parent_hash, transactions, timestamp, created_at 
  FROM blocks
  WHERE hash = ?
	`, hash).Rows()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get block by hash: %v", hash)
	}

	if rows.Next() {
		block := Block{}
		var transactionsBytes []byte
		if err = rows.Scan(&block.Hash, &block.Number, &block.ParentHash, &transactionsBytes); err != nil {
			return nil, errors.Wrap(err, "unexpected error scanning row")
		}

		err = json.Unmarshal(transactionsBytes, &block.Transactions)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal transactions from json")
		}
		return &block, nil
	}

	return nil, errors.New("block not found")
}
