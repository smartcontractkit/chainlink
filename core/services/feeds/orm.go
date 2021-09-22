package feeds

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	ApproveJobProposal(ctx context.Context, id int64, externalJobID uuid.UUID, status JobProposalStatus) error
	CancelJobProposal(ctx context.Context, id int64) error
	CountJobProposals(ctx context.Context) (int64, error)
	CountManagers(ctx context.Context) (int64, error)
	CreateJobProposal(ctx context.Context, jp *JobProposal) (int64, error)
	CreateManager(ctx context.Context, ms *FeedsManager) (int64, error)
	GetJobProposal(ctx context.Context, id int64) (*JobProposal, error)
	GetJobProposalByRemoteUUID(ctx context.Context, uuid uuid.UUID) (*JobProposal, error)
	GetManager(ctx context.Context, id int64) (*FeedsManager, error)
	IsJobManaged(ctx context.Context, jobID int64) (bool, error)
	ListJobProposals(ctx context.Context) ([]JobProposal, error)
	ListManagers(ctx context.Context) ([]FeedsManager, error)
	UpdateJobProposalSpec(ctx context.Context, id int64, spec string) error
	UpdateJobProposalStatus(ctx context.Context, id int64, status JobProposalStatus) error
	UpdateManager(ctx context.Context, mgr FeedsManager) error
	UpsertJobProposal(ctx context.Context, jp *JobProposal) (int64, error)
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
INSERT INTO feeds_managers (name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;
`

	row := o.db.WithContext(ctx).Raw(stmt,
		ms.Name,
		ms.URI,
		ms.PublicKey,
		ms.JobTypes,
		ms.IsOCRBootstrapPeer,
		ms.OCRBootstrapPeerMultiaddr,
		now,
		now,
	).Row()
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
SELECT id, name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at
FROM feeds_managers;
`

	err := o.db.WithContext(ctx).Raw(stmt).Scan(&mgrs).Error
	if err != nil {
		return mgrs, err
	}

	return mgrs, nil
}

// GetManager gets a feeds manager by id
func (o *orm) GetManager(ctx context.Context, id int64) (*FeedsManager, error) {
	stmt := `
SELECT id, name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at
FROM feeds_managers
WHERE id = ?;
`

	mgr := FeedsManager{}
	result := o.db.WithContext(ctx).Raw(stmt, id).Scan(&mgr)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &mgr, nil
}

func (o *orm) UpdateManager(ctx context.Context, mgr FeedsManager) error {
	tx := postgres.TxFromContext(ctx, o.db)
	now := time.Now()

	stmt := `
UPDATE feeds_managers
SET name = ?,
	uri = ?,
	public_key = ?,
	job_types = ?,
	is_ocr_bootstrap_peer = ?,
	ocr_bootstrap_peer_multiaddr = ?,
	updated_at = ?
WHERE id = ?;
`

	result := tx.Exec(stmt, mgr.Name, mgr.URI, mgr.PublicKey, mgr.JobTypes, mgr.IsOCRBootstrapPeer, mgr.OCRBootstrapPeerMultiaddr, now, mgr.ID)
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}

	return nil

}

// Count counts the number of feeds manager records.
func (o *orm) CountManagers(ctx context.Context) (int64, error) {
	var count int64
	stmt := `
SELECT COUNT(*)
FROM feeds_managers
	`

	err := o.db.WithContext(ctx).Raw(stmt).Scan(&count).Error
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
INSERT INTO job_proposals (remote_uuid, spec, status, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;
`

	row := o.db.WithContext(ctx).Raw(stmt, jp.RemoteUUID, jp.Spec, jp.Status, jp.FeedsManagerID, jp.Multiaddrs, now, now, now).Row()
	if row.Err() != nil {
		return id, row.Err()
	}

	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, err
}

// UpsertJobProposal creates a job proposal if it does not exist. If it does exist,
// then we update the details of the existing job proposal only if the provided
// feeds manager id exists.
func (o *orm) UpsertJobProposal(ctx context.Context, jp *JobProposal) (int64, error) {
	var id int64
	now := time.Now()

	stmt := `
INSERT INTO job_proposals (remote_uuid, spec, status, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (remote_uuid)
DO
	UPDATE SET
		spec = excluded.spec,
		status = excluded.status,
		multiaddrs = excluded.multiaddrs,
		proposed_at = excluded.proposed_at,
		updated_at = excluded.updated_at
RETURNING id;
`

	row := o.db.WithContext(ctx).Raw(stmt,
		jp.RemoteUUID, jp.Spec, jp.Status, jp.FeedsManagerID, jp.Multiaddrs, now, now, now,
	).Row()
	if row.Err() != nil {
		return id, row.Err()
	}

	err := row.Scan(&id)
	return id, err
}

// ListJobProposals lists all job proposals
func (o *orm) ListJobProposals(ctx context.Context) ([]JobProposal, error) {
	jps := []JobProposal{}
	stmt := `
SELECT remote_uuid, id, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals;
`

	err := o.db.WithContext(ctx).Raw(stmt).Scan(&jps).Error
	if err != nil {
		return jps, err
	}

	return jps, nil
}

// GetJobProposal gets a job proposal by id
func (o *orm) GetJobProposal(ctx context.Context, id int64) (*JobProposal, error) {
	stmt := `
SELECT id, remote_uuid, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals
WHERE id = ?;
`

	return o.getJobProposal(ctx, stmt, id)
}

// GetJobProposalByRemoteUUID gets a job proposal by the remote FMS uuid
func (o *orm) GetJobProposalByRemoteUUID(ctx context.Context, id uuid.UUID) (*JobProposal, error) {
	stmt := `
SELECT id, remote_uuid, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals
WHERE remote_uuid = ?;
`

	return o.getJobProposal(ctx, stmt, id)
}

// getJobProposal performs the db call to fetch a single job proposal record
func (o *orm) getJobProposal(ctx context.Context, stmt string, values ...interface{}) (*JobProposal, error) {
	jp := JobProposal{}
	result := o.db.WithContext(ctx).Raw(stmt, values...).Scan(&jp)
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &jp, nil
}

// UpdateJobProposalStatus updates the status of a job proposal by id.
func (o *orm) UpdateJobProposalStatus(ctx context.Context, id int64, status JobProposalStatus) error {
	tx := postgres.TxFromContext(ctx, o.db)

	now := time.Now()

	stmt := `
UPDATE job_proposals
SET status = ?,
	updated_at = ?
WHERE id = ?;
`

	result := tx.Exec(stmt, status, now, id)
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateJobProposalSpec updates the spec of a job proposal by id.
func (o *orm) UpdateJobProposalSpec(ctx context.Context, id int64, spec string) error {
	tx := postgres.TxFromContext(ctx, o.db)
	now := time.Now()

	stmt := `
UPDATE job_proposals
SET spec = ?,
	updated_at = ?
WHERE id = ?;
`

	result := tx.Exec(stmt, spec, now, id)
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ApproveJobProposal updates the job proposal as approved.
func (o *orm) ApproveJobProposal(ctx context.Context, id int64, externalJobID uuid.UUID, status JobProposalStatus) error {
	tx := postgres.TxFromContext(ctx, o.db)
	now := time.Now()

	stmt := `
UPDATE job_proposals
SET status = ?,
	external_job_id = ?,
	updated_at = ?
WHERE id = ?;
`

	result := tx.Exec(stmt, status, externalJobID, now, id)
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CancelJobProposal cancels a job proposal.
func (o *orm) CancelJobProposal(ctx context.Context, id int64) error {
	tx := postgres.TxFromContext(ctx, o.db)
	now := time.Now()
	stmt := `
UPDATE job_proposals
SET status = ?,
	external_job_id = ?,
	updated_at = ?
WHERE id = ?;
`

	result := tx.Exec(stmt, JobProposalStatusCancelled, nil, now, id)
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposals(ctx context.Context) (int64, error) {
	var count int64
	stmt := `
SELECT COUNT(*)
FROM job_proposals
`

	err := o.db.WithContext(ctx).Raw(stmt).Scan(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}

// IsJobManaged determines if a job is managed by the feeds manager.
func (o *orm) IsJobManaged(ctx context.Context, jobID int64) (bool, error) {
	stmt := `
SELECT exists (
	SELECT 1
	FROM job_proposals
	INNER JOIN jobs ON job_proposals.external_job_id = jobs.external_job_id
	WHERE jobs.id = ?
);
`

	var exists bool
	err := o.db.WithContext(ctx).Raw(stmt, jobID).Scan(&exists).Error

	return exists, err
}
