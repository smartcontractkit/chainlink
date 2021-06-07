package feeds

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CountManagers() (int64, error)
	CreateManager(ctx context.Context, ms *FeedsManager) (int32, error)
	GetManager(ctx context.Context, id int32) (*FeedsManager, error)
	ListManagers(ctx context.Context) ([]FeedsManager, error)
}

type orm struct {
	db *gorm.DB
}

func NewORM(db *gorm.DB) *orm {
	return &orm{
		db: db,
	}
}

// CreateManager creates a feeds manager.
func (o *orm) CreateManager(ctx context.Context, ms *FeedsManager) (int32, error) {
	var id int32
	now := time.Now()

	// Create the ManagerService
	stmt := `
		INSERT INTO feeds_managers (name, uri, public_key, job_types, network, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	row := o.db.Raw(stmt, ms.Name, ms.URI, ms.PublicKey, ms.JobTypes, ms.Network, now, now).Row()
	if row.Err() != nil {
		return id, row.Err()
	}

	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, err
}

// ListManager lists all feeds managers
func (o *orm) ListManagers(ctx context.Context) ([]FeedsManager, error) {
	mgrs := []FeedsManager{}
	stmt := `
		SELECT id, name, uri, public_key, job_types, network, created_at, updated_at
		FROM feeds_managers;
	`

	err := o.db.Raw(stmt).Scan(&mgrs).Error
	if err != nil {
		return mgrs, err
	}

	return mgrs, nil
}

// GetManager gets a feeds manager by id
func (o *orm) GetManager(ctx context.Context, id int32) (*FeedsManager, error) {
	stmt := `
		SELECT id, name, uri, public_key, job_types, network, created_at, updated_at
		FROM feeds_managers
		WHERE id = ?;
	`

	mgr := FeedsManager{}
	result := o.db.Raw(stmt, id).Scan(&mgr)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &mgr, nil
}

// Count counts the number of feeds manager records.
func (o *orm) CountManagers() (int64, error) {
	var count int64
	stmt := `
		SELECT COUNT(*)
		FROM feeds_managers
	`

	err := o.db.Raw(stmt).Scan(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}
