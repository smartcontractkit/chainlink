package feeds

import (
	"context"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	Count() (int64, error)
	CreateManagerService(ctx context.Context, ms *ManagerService) (int32, error)
	GetManagerService(ctx context.Context, id int32) (*ManagerService, error)
	ListManagerServices(ctx context.Context) ([]ManagerService, error)
}

type orm struct {
	db *gorm.DB
}

func NewORM(db *gorm.DB) *orm {
	return &orm{
		db: db,
	}
}

// CreateFeedsManager creates a feeds manager record.
func (o *orm) CreateManagerService(ctx context.Context, ms *ManagerService) (int32, error) {
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

func (o *orm) ListManagerServices(ctx context.Context) ([]ManagerService, error) {
	mss := []ManagerService{}
	stmt := `
		SELECT id, name, uri, public_key, job_types, network, created_at, updated_at
		FROM feeds_managers;
	`

	rows, err := o.db.Raw(stmt).Rows()
	if err != nil {
		return mss, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		ms := ManagerService{}
		err := rows.Scan(&ms.ID, &ms.Name, &ms.URI, &ms.PublicKey, &ms.JobTypes, &ms.Network, &ms.CreatedAt, &ms.UpdatedAt)
		if err != nil {
			return mss, err
		}
		mss = append(mss, ms)
	}

	return mss, nil
}

func (o *orm) GetManagerService(ctx context.Context, id int32) (*ManagerService, error) {
	stmt := `
		SELECT id, name, uri, public_key, job_types, network, created_at, updated_at
		FROM feeds_managers
		WHERE id = ?;
	`

	row := o.db.Raw(stmt, id).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}

	ms := ManagerService{}
	err := row.Scan(&ms.ID, &ms.Name, &ms.URI, &ms.PublicKey, &ms.JobTypes, &ms.Network, &ms.CreatedAt, &ms.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ms, nil
}

// Count counts the number of manager service records.
func (o *orm) Count() (int64, error) {
	var count int64
	stmt := `
		SELECT COUNT(*)
		FROM feeds_managers
	`

	row := o.db.Raw(stmt).Row()
	if row.Err() != nil {
		return count, row.Err()
	}

	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
