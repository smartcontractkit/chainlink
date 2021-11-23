package feeds

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	ApproveJobProposal(id int64, externalJobID uuid.UUID, status JobProposalStatus, qopts ...pg.QOpt) error
	CancelJobProposal(id int64, qopts ...pg.QOpt) error
	CountJobProposals() (int64, error)
	CountManagers() (int64, error)
	CreateJobProposal(jp *JobProposal) (int64, error)
	CreateManager(ms *FeedsManager) (int64, error)
	GetJobProposal(id int64, qopts ...pg.QOpt) (*JobProposal, error)
	GetJobProposalByManagersIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposal, error)
	GetJobProposalByRemoteUUID(uuid uuid.UUID) (*JobProposal, error)
	GetManager(id int64) (*FeedsManager, error)
	GetManagers(ids []int64) ([]FeedsManager, error)
	IsJobManaged(jobID int64, qopts ...pg.QOpt) (bool, error)
	ListJobProposals() (jps []JobProposal, err error)
	ListManagers() (mgrs []FeedsManager, err error)
	UpdateJobProposalSpec(id int64, spec string, qopts ...pg.QOpt) error
	UpdateJobProposalStatus(id int64, status JobProposalStatus, qopts ...pg.QOpt) error
	UpdateManager(mgr FeedsManager, qopts ...pg.QOpt) error
	UpsertJobProposal(jp *JobProposal) (int64, error)
}

var _ ORM = &orm{}

type orm struct {
	q pg.Q
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) *orm {
	return &orm{
		q: pg.NewQ(db, lggr, cfg),
	}
}

// CreateManager creates a feeds manager.
func (o *orm) CreateManager(ms *FeedsManager) (id int64, err error) {
	stmt := `
INSERT INTO feeds_managers (name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
RETURNING id;
`
	err = o.q.Get(&id, stmt,
		ms.Name, ms.URI, ms.PublicKey, ms.JobTypes, ms.IsOCRBootstrapPeer, ms.OCRBootstrapPeerMultiaddr,
	)
	if err != nil {
		return id, errors.Wrap(err, "CreateManager failed")
	}
	return id, err
}

// ListManager lists all feeds managers
func (o *orm) ListManagers() (mgrs []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at
FROM feeds_managers;
`

	err = o.q.Select(&mgrs, stmt)
	return mgrs, errors.Wrap(err, "ListManagers failed")
}

// GetManager gets a feeds manager by id
func (o *orm) GetManager(id int64) (mgr *FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at
FROM feeds_managers
WHERE id = $1
`

	mgr = new(FeedsManager)
	err = o.q.Get(mgr, stmt, id)
	return mgr, errors.Wrap(err, "GetManager failed")
}

// GetManagers gets feeds managers by ids
func (o *orm) GetManagers(ids []int64) (managers []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, job_types, is_ocr_bootstrap_peer, ocr_bootstrap_peer_multiaddr, created_at, updated_at
FROM feeds_managers
WHERE id = ANY($1)
ORDER BY created_at, id;`

	mgrIds := pq.Array(ids)
	err = o.q.Select(&managers, stmt, mgrIds)

	return managers, errors.Wrap(err, "GetManagers failed")
}

func (o *orm) UpdateManager(mgr FeedsManager, qopts ...pg.QOpt) (err error) {
	stmt := `
UPDATE feeds_managers
SET name = $1, uri = $2, public_key = $3, job_types = $4, is_ocr_bootstrap_peer = $5, ocr_bootstrap_peer_multiaddr = $6, updated_at = NOW()
WHERE id = $7;
`

	res, err := o.q.WithOpts(qopts...).Exec(stmt, mgr.Name, mgr.URI, mgr.PublicKey, mgr.JobTypes, mgr.IsOCRBootstrapPeer, mgr.OCRBootstrapPeerMultiaddr, mgr.ID)
	if err != nil {
		return errors.Wrap(err, "UpdateManager failed to update feeds_managers")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "UpdateManager failed to get RowsAffected")
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil

}

// Count counts the number of feeds manager records.
func (o *orm) CountManagers() (count int64, err error) {
	stmt := `
SELECT COUNT(*)
FROM feeds_managers
	`

	err = o.q.Get(&count, stmt)
	return count, errors.Wrap(err, "CountManagers failed")
}

// CreateJobProposal creates a job proposal.
func (o *orm) CreateJobProposal(jp *JobProposal) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (remote_uuid, spec, status, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), NOW())
RETURNING id;
`

	err = o.q.Get(&id, stmt, jp.RemoteUUID, jp.Spec, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "CreateJobProposal failed")
}

// UpsertJobProposal creates a job proposal if it does not exist. If it does exist,
// then we update the details of the existing job proposal only if the provided
// feeds manager id exists.
func (o *orm) UpsertJobProposal(jp *JobProposal) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (remote_uuid, spec, status, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), NOW())
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

	err = o.q.Get(&id, stmt, jp.RemoteUUID, jp.Spec, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "UpsertJobProposal")
}

// ListJobProposals lists all job proposals
func (o *orm) ListJobProposals() (jps []JobProposal, err error) {
	stmt := `
SELECT remote_uuid, id, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals;
`

	err = o.q.Select(&jps, stmt)
	return jps, errors.Wrap(err, "ListJobProposals failed")
}

// GetJobProposal gets a job proposal by id
func (o *orm) GetJobProposal(id int64, qopts ...pg.QOpt) (jp *JobProposal, err error) {
	stmt := `
SELECT id, remote_uuid, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals
WHERE id = $1
`
	jp = new(JobProposal)
	err = o.q.WithOpts(qopts...).Get(jp, stmt, id)
	return jp, errors.Wrap(err, "GetJobProposal failed")
}

// GetJobProposalByManagersIDs gets job proposals by feeds managers IDs
func (o *orm) GetJobProposalByManagersIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposal, error) {
	stmt := `
SELECT id, remote_uuid, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals
WHERE feeds_manager_id = ANY($1)
`
	var jps []JobProposal
	err := o.q.WithOpts(qopts...).Select(&jps, stmt, ids)
	return jps, errors.Wrap(err, "GetJobProposalByManagersIDs failed")
}

// GetJobProposalByRemoteUUID gets a job proposal by the remote FMS uuid
func (o *orm) GetJobProposalByRemoteUUID(id uuid.UUID) (jp *JobProposal, err error) {
	stmt := `
SELECT id, remote_uuid, spec, status, external_job_id, feeds_manager_id, multiaddrs, proposed_at, created_at, updated_at
FROM job_proposals
WHERE remote_uuid = $1;
`

	jp = new(JobProposal)
	err = o.q.Get(jp, stmt, id)
	return jp, errors.Wrap(err, "GetJobProposalByRemoteUUID failed")
}

// UpdateJobProposalStatus updates the status of a job proposal by id.
func (o *orm) UpdateJobProposalStatus(id int64, status JobProposalStatus, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposals
SET status = $1,
	updated_at = NOW()
WHERE id = $2;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateJobProposalSpec updates the spec of a job proposal by id.
func (o *orm) UpdateJobProposalSpec(id int64, spec string, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposals
SET spec = $1,
	updated_at = NOW()
WHERE id = $2;
`

	res, err := o.q.WithOpts(qopts...).Exec(stmt, spec, id)
	if err != nil {
		return errors.Wrap(err, "UpdateJobProposalSpec failed to update job_proposals")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "UpdateJobProposalSpec failed to get RowsAffected")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ApproveJobProposal updates the job proposal as approved.
func (o *orm) ApproveJobProposal(id int64, externalJobID uuid.UUID, status JobProposalStatus, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposals
SET status = $1,
	external_job_id = $2,
	updated_at = NOW()
WHERE id = $3;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, status, externalJobID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// CancelJobProposal cancels a job proposal.
func (o *orm) CancelJobProposal(id int64, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposals
SET status = $1,
	external_job_id = $2,
	updated_at = NOW()
WHERE id = $3;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, JobProposalStatusCancelled, nil, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposals() (count int64, err error) {
	stmt := `SELECT COUNT(*) FROM job_proposals`

	err = o.q.Get(&count, stmt)
	return count, errors.Wrap(err, "CountJobProposals failed")
}

// IsJobManaged determines if a job is managed by the feeds manager.
func (o *orm) IsJobManaged(jobID int64, qopts ...pg.QOpt) (exists bool, err error) {
	stmt := `
SELECT exists (
	SELECT 1
	FROM job_proposals
	INNER JOIN jobs ON job_proposals.external_job_id = jobs.external_job_id
	WHERE jobs.id = $1
);
`

	err = o.q.WithOpts(qopts...).Get(&exists, stmt, jobID)
	return exists, errors.Wrap(err, "IsJobManaged failed")
}
