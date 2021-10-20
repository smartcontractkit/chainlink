package keystore

import (
	"context"
	"database/sql"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"gorm.io/gorm"
)

type csaORM struct {
	db *gorm.DB
}

func NewCSAORM(db *gorm.DB) csaORM {
	return csaORM{db}
}

// CreateCSAKey creates a CSA key.
func (o csaORM) CreateCSAKey(ctx context.Context, kp *csakey.Key) (uint, error) {
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
func (o csaORM) ListCSAKeys(ctx context.Context) ([]csakey.Key, error) {
	keys := []csakey.Key{}
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
func (o csaORM) GetCSAKey(ctx context.Context, id uint) (*csakey.Key, error) {
	stmt := `
	SELECT id, public_key, encrypted_private_key, created_at, updated_at
		FROM csa_keys
		WHERE id = ?;
	`

	key := csakey.Key{}
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
func (o csaORM) CountCSAKeys() (int64, error) {
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
