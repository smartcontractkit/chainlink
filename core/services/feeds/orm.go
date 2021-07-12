package feeds

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CountJobProposals() (int64, error)
	CountManagers() (int64, error)
	CreateJobProposal(ctx context.Context, jp *JobProposal) (int64, error)
	CreateManager(ctx context.Context, ms *FeedsManager) (int64, error)
	GetJobProposal(ctx context.Context, id int64) (*JobProposal, error)
	GetManager(ctx context.Context, id int64) (*FeedsManager, error)
	ListJobProposals(ctx context.Context) ([]JobProposal, error)
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
func (o *orm) CreateManager(ctx context.Context, ms *FeedsManager) (int64, error) {
	var id int64
	now := time.Now()

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
func (o *orm) GetManager(ctx context.Context, id int64) (*FeedsManager, error) {
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

// CreateJobProposal creates a job proposal.
func (o *orm) CreateJobProposal(ctx context.Context, jp *JobProposal) (int64, error) {
	var id int64
	now := time.Now()

	stmt := `
		INSERT INTO job_proposals (spec, status, feeds_manager_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id;
	`

	row := o.db.Raw(stmt, jp.Spec, jp.Status, jp.FeedsManagerID, now, now).Row()
	if row.Err() != nil {
		return id, row.Err()
	}

	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, err
}

// ListJobProposals lists all job proposals
func (o *orm) ListJobProposals(ctx context.Context) ([]JobProposal, error) {
	jps := []JobProposal{}
	stmt := `
		SELECT id, spec, status, job_id, feeds_manager_id, created_at, updated_at
		FROM job_proposals;
	`

	err := o.db.Raw(stmt).Scan(&jps).Error
	if err != nil {
		return jps, err
	}

	return jps, nil
}

// GetJobProposal gets a job proposal by id
func (o *orm) GetJobProposal(ctx context.Context, id int64) (*JobProposal, error) {
	stmt := `
		SELECT id, spec, status, job_id, feeds_manager_id, created_at, updated_at
		FROM job_proposals
		WHERE id = ?;
	`

	jp := JobProposal{}
	result := o.db.Raw(stmt, id).Scan(&jp)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &jp, nil
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposals() (int64, error) {
	var count int64
	stmt := `
		SELECT COUNT(*)
		FROM job_proposals
	`

	err := o.db.Raw(stmt).Scan(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}
