package csa

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CreateCSAKey(ctx context.Context, kp *CSAKey) (uint, error)
	GetCSAKey(ctx context.Context, id uint) (*CSAKey, error)
	ListCSAKeys(ctx context.Context) ([]CSAKey, error)
	CountCSAKeys() (int64, error)
}

type orm struct {
	db *gorm.DB
}

func NewORM(db *gorm.DB) *orm {
	return &orm{db}
}

// CreateCSAKey creates a CSA key.
func (o *orm) CreateCSAKey(ctx context.Context, kp *CSAKey) (uint, error) {
	var id uint
	now := time.Now()

	// Create the CSA Key
	stmt := `
		INSERT INTO csa_keys (public_key, encrypted_private_key, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		RETURNING id;
	`

	err := o.db.Raw(stmt, kp.PublicKey, kp.EncryptedPrivateKey, now, now).Scan(&id).Error
	if err != nil {
		return id, err
	}

	return id, err
}

// ListCSAKeys lists all csa keys.
func (o *orm) ListCSAKeys(ctx context.Context) ([]CSAKey, error) {
	keys := []CSAKey{}
	stmt := `
		SELECT id, public_key, encrypted_private_key, created_at, updated_at
		FROM csa_keys;
	`

	err := o.db.Raw(stmt).Scan(&keys).Error
	if err != nil {
		return keys, err
	}

	return keys, nil
}

// GetCSAKey gets a CSA key by id
func (o *orm) GetCSAKey(ctx context.Context, id uint) (*CSAKey, error) {
	stmt := `
	SELECT id, public_key, encrypted_private_key, created_at, updated_at
		FROM csa_keys
		WHERE id = ?;
	`

	key := CSAKey{}
	result := o.db.Raw(stmt, id).Scan(&key)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &key, nil
}

// Count counts the number of csa key records.
func (o *orm) CountCSAKeys() (int64, error) {
	var count int64
	stmt := `
		SELECT COUNT(*)
		FROM csa_keys
	`

	err := o.db.Raw(stmt).Scan(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}
