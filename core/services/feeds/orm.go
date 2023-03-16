package feeds

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

//go:generate mockery --with-expecter=true --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CountManagers() (int64, error)
	CreateManager(ms *FeedsManager, qopts ...pg.QOpt) (int64, error)
	GetManager(id int64) (*FeedsManager, error)
	ListManagers() (mgrs []FeedsManager, err error)
	ListManagersByIDs(ids []int64) ([]FeedsManager, error)
	UpdateManager(mgr FeedsManager, qopts ...pg.QOpt) error

	CreateChainConfig(cfg ChainConfig, qopts ...pg.QOpt) (int64, error)
	CreateBatchChainConfig(cfgs []ChainConfig, qopts ...pg.QOpt) ([]int64, error)
	DeleteChainConfig(id int64) (int64, error)
	GetChainConfig(id int64) (*ChainConfig, error)
	UpdateChainConfig(cfg ChainConfig) (int64, error)
	ListChainConfigsByManagerIDs(mgrIDs []int64) ([]ChainConfig, error)

	CreateJobProposal(jp *JobProposal) (int64, error)
	CountJobProposals() (int64, error)
	CountJobProposalsByStatus() (counts *JobProposalCounts, err error)
	GetJobProposal(id int64, qopts ...pg.QOpt) (*JobProposal, error)
	GetJobProposalByRemoteUUID(uuid uuid.UUID) (*JobProposal, error)
	ListJobProposals() (jps []JobProposal, err error)
	ListJobProposalsByManagersIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposal, error)
	UpdateJobProposalStatus(id int64, status JobProposalStatus, qopts ...pg.QOpt) error // NEEDED?
	UpsertJobProposal(jp *JobProposal, qopts ...pg.QOpt) (int64, error)

	ApproveSpec(id int64, externalJobID uuid.UUID, qopts ...pg.QOpt) error
	CancelSpec(id int64, qopts ...pg.QOpt) error
	CreateSpec(spec JobProposalSpec, qopts ...pg.QOpt) (int64, error)
	ExistsSpecByJobProposalIDAndVersion(jpID int64, version int32, qopts ...pg.QOpt) (exists bool, err error)
	GetLatestSpec(jpID int64) (*JobProposalSpec, error)
	GetApprovedSpec(jpID int64, qopts ...pg.QOpt) (*JobProposalSpec, error)
	GetSpec(id int64, qopts ...pg.QOpt) (*JobProposalSpec, error)
	ListSpecsByJobProposalIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposalSpec, error)
	RejectSpec(id int64, qopts ...pg.QOpt) error
	UpdateSpecDefinition(id int64, spec string, qopts ...pg.QOpt) error

	IsJobManaged(jobID int64, qopts ...pg.QOpt) (bool, error)
}

var _ ORM = &orm{}

type orm struct {
	q pg.Q
}

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *orm {
	return &orm{
		q: pg.NewQ(db, lggr, cfg),
	}
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

// CreateManager creates a feeds manager.
func (o *orm) CreateManager(ms *FeedsManager, qopts ...pg.QOpt) (id int64, err error) {
	stmt := `
INSERT INTO feeds_managers (name, uri, public_key, created_at, updated_at)
VALUES ($1,$2,$3,NOW(),NOW())
RETURNING id;
`
	err = o.q.WithOpts(qopts...).Get(&id, stmt, ms.Name, ms.URI, ms.PublicKey)

	return id, errors.Wrap(err, "CreateManager failed")
}

// CreateChainConfig creates a new chain config.
func (o *orm) CreateChainConfig(cfg ChainConfig, qopts ...pg.QOpt) (id int64, err error) {
	stmt := `
INSERT INTO feeds_manager_chain_configs (feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
RETURNING id;
`

	err = o.q.WithOpts(qopts...).Get(&id,
		stmt,
		cfg.FeedsManagerID,
		cfg.ChainID,
		cfg.ChainType,
		cfg.AccountAddress,
		cfg.AdminAddress,
		cfg.FluxMonitorConfig,
		cfg.OCR1Config,
		cfg.OCR2Config,
	)

	return id, errors.Wrap(err, "CreateChainConfig failed")
}

// CreateBatchChainConfig creates multiple chain configs.
func (o *orm) CreateBatchChainConfig(cfgs []ChainConfig, qopts ...pg.QOpt) (ids []int64, err error) {
	if len(cfgs) == 0 {
		return
	}

	stmt := `
INSERT INTO feeds_manager_chain_configs (feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at)
VALUES %s
RETURNING id;
	`

	var (
		vStrs = make([]string, 0, len(cfgs))
		vArgs = make([]interface{}, 0)
	)

	for i, cfg := range cfgs {
		// Generate the placeholders
		pnumidx := i * 8

		lo, hi := pnumidx+1, pnumidx+8
		pnums := make([]any, hi-lo+1)
		for i := range pnums {
			pnums[i] = i + lo
		}

		vStrs = append(vStrs, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, NOW(), NOW())", pnums...,
		))

		// Append the values
		vArgs = append(vArgs,
			cfg.FeedsManagerID,
			cfg.ChainID,
			cfg.ChainType,
			cfg.AccountAddress,
			cfg.AdminAddress,
			cfg.FluxMonitorConfig,
			cfg.OCR1Config,
			cfg.OCR2Config,
		)
	}

	err = o.q.WithOpts(qopts...).Select(&ids,
		fmt.Sprintf(stmt, strings.Join(vStrs, ",")),
		vArgs...,
	)

	return ids, errors.Wrap(err, "CreateBatchChainConfig failed")
}

// DeleteChainConfig deletes a chain config.
func (o *orm) DeleteChainConfig(id int64) (int64, error) {
	stmt := `
DELETE FROM feeds_manager_chain_configs
WHERE id = $1
RETURNING id;
`

	var ccid int64
	err := o.q.Get(&ccid, stmt, id)

	return ccid, errors.Wrap(err, "DeleteChainConfig failed")
}

// GetChainConfig fetches a chain config.
func (o *orm) GetChainConfig(id int64) (*ChainConfig, error) {
	stmt := `
SELECT id, feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at
FROM feeds_manager_chain_configs
WHERE id = $1;
`

	var cfg ChainConfig
	err := o.q.Get(&cfg, stmt, id)

	return &cfg, errors.Wrap(err, "GetChainConfig failed")
}

// ListChainConfigsByManagerIDs fetches the chain configs matching all manager
// ids.
func (o *orm) ListChainConfigsByManagerIDs(mgrIDs []int64) ([]ChainConfig, error) {
	stmt := `
SELECT id, feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at
FROM feeds_manager_chain_configs
WHERE feeds_manager_id = ANY($1)
	`

	var cfgs []ChainConfig
	err := o.q.Select(&cfgs, stmt, mgrIDs)

	return cfgs, errors.Wrap(err, "ListJobProposalsByManagersIDs failed")
}

// UpdateChainConfig updates a chain config.
func (o *orm) UpdateChainConfig(cfg ChainConfig) (int64, error) {
	stmt := `
UPDATE feeds_manager_chain_configs
SET account_address = $1,
	admin_address = $2,
	flux_monitor_config = $3,
	ocr1_config = $4,
	ocr2_config = $5,
	updated_at = NOW()
WHERE id = $6
RETURNING id;
`

	var cfgID int64
	err := o.q.Get(&cfgID, stmt,
		cfg.AccountAddress,
		cfg.AdminAddress,
		cfg.FluxMonitorConfig,
		cfg.OCR1Config,
		cfg.OCR2Config,
		cfg.ID,
	)

	return cfgID, errors.Wrap(err, "UpdateChainConfig failed")
}

// GetManager gets a feeds manager by id.
func (o *orm) GetManager(id int64) (mgr *FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers
WHERE id = $1
`

	mgr = new(FeedsManager)
	err = o.q.Get(mgr, stmt, id)
	return mgr, errors.Wrap(err, "GetManager failed")
}

// ListManager lists all feeds managers.
func (o *orm) ListManagers() (mgrs []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers;
`

	err = o.q.Select(&mgrs, stmt)
	return mgrs, errors.Wrap(err, "ListManagers failed")
}

// ListManagersByIDs gets feeds managers by ids.
func (o *orm) ListManagersByIDs(ids []int64) (managers []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers
WHERE id = ANY($1)
ORDER BY created_at, id;`

	mgrIds := pq.Array(ids)
	err = o.q.Select(&managers, stmt, mgrIds)

	return managers, errors.Wrap(err, "GetManagers failed")
}

// UpdateManager updates the manager details.
func (o *orm) UpdateManager(mgr FeedsManager, qopts ...pg.QOpt) (err error) {
	stmt := `
UPDATE feeds_managers
SET name = $1, uri = $2, public_key = $3, updated_at = NOW()
WHERE id = $4;
`

	res, err := o.q.WithOpts(qopts...).Exec(stmt, mgr.Name, mgr.URI, mgr.PublicKey, mgr.ID)
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

// CreateJobProposal creates a job proposal.
func (o *orm) CreateJobProposal(jp *JobProposal) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (name, remote_uuid, status, feeds_manager_id, multiaddrs, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id;
`

	err = o.q.Get(&id, stmt, jp.Name, jp.RemoteUUID, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "CreateJobProposal failed")
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposals() (count int64, err error) {
	stmt := `SELECT COUNT(*) FROM job_proposals`

	err = o.q.Get(&count, stmt)
	return count, errors.Wrap(err, "CountJobProposals failed")
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposalsByStatus() (counts *JobProposalCounts, err error) {
	stmt := `
SELECT 
	COUNT(*) filter (where job_proposals.status = 'pending' OR job_proposals.pending_update = TRUE) as pending,
	COUNT(*) filter (where job_proposals.status = 'approved' AND job_proposals.pending_update = FALSE) as approved,
	COUNT(*) filter (where job_proposals.status = 'rejected' AND job_proposals.pending_update = FALSE) as rejected,
	COUNT(*) filter (where job_proposals.status = 'cancelled' AND job_proposals.pending_update = FALSE) as cancelled
FROM job_proposals;
	`

	counts = new(JobProposalCounts)
	err = o.q.Get(counts, stmt)
	return counts, errors.Wrap(err, "CountJobProposalsByStatus failed")
}

// GetJobProposal gets a job proposal by id.
func (o *orm) GetJobProposal(id int64, qopts ...pg.QOpt) (jp *JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE id = $1
`
	jp = new(JobProposal)
	err = o.q.WithOpts(qopts...).Get(jp, stmt, id)
	return jp, errors.Wrap(err, "GetJobProposal failed")
}

// GetJobProposalByRemoteUUID gets a job proposal by the remote FMS uuid.
func (o *orm) GetJobProposalByRemoteUUID(id uuid.UUID) (jp *JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE remote_uuid = $1;
`

	jp = new(JobProposal)
	err = o.q.Get(jp, stmt, id)
	return jp, errors.Wrap(err, "GetJobProposalByRemoteUUID failed")
}

// ListJobProposals lists all job proposals.
func (o *orm) ListJobProposals() (jps []JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals;
`

	err = o.q.Select(&jps, stmt)
	return jps, errors.Wrap(err, "ListJobProposals failed")
}

// ListJobProposalsByManagersIDs gets job proposals by feeds managers IDs.
func (o *orm) ListJobProposalsByManagersIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposal, error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE feeds_manager_id = ANY($1)
`
	var jps []JobProposal
	err := o.q.WithOpts(qopts...).Select(&jps, stmt, ids)
	return jps, errors.Wrap(err, "ListJobProposalsByManagersIDs failed")
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

// UpsertJobProposal creates a job proposal if it does not exist. If it does exist,
// then we update the details of the existing job proposal only if the provided
// feeds manager id exists.
func (o *orm) UpsertJobProposal(jp *JobProposal, qopts ...pg.QOpt) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (name, remote_uuid, status, feeds_manager_id, multiaddrs, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (remote_uuid)
DO
	UPDATE SET
		pending_update = TRUE,
		name = EXCLUDED.name,
		multiaddrs = EXCLUDED.multiaddrs,
		updated_at = EXCLUDED.updated_at
RETURNING id;
`

	err = o.q.WithOpts(qopts...).Get(&id, stmt, jp.Name, jp.RemoteUUID, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "UpsertJobProposal")
}

// ApproveSpec approves the spec and sets the external job ID on the associated
// job proposal.
func (o *orm) ApproveSpec(id int64, externalJobID uuid.UUID, qopts ...pg.QOpt) error {
	// Update the status of the approval
	stmt := `
UPDATE job_proposal_specs
SET status = $1,
	status_updated_at = NOW(),
	updated_at = NOW()
WHERE id = $2
RETURNING job_proposal_id;
`

	var jpID int64
	if err := o.q.WithOpts(qopts...).Get(&jpID, stmt, JobProposalStatusApproved, id); err != nil {
		return err
	}

	// Update the job proposal external id
	stmt = `
UPDATE job_proposals
SET status = $1,
	external_job_id = $2,
	pending_update = FALSE,
	updated_at = NOW()
WHERE id = $3;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, JobProposalStatusApproved, externalJobID, jpID)
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

// CancelSpec cancels the spec and removes the external job id from the
// associated job proposal.
func (o *orm) CancelSpec(id int64, qopts ...pg.QOpt) error {
	// Update the status of the approval
	stmt := `
UPDATE job_proposal_specs
SET status = $1,
	status_updated_at = NOW(),
	updated_at = NOW()
WHERE id = $2
RETURNING job_proposal_id;
`

	var jpID int64
	if err := o.q.WithOpts(qopts...).Get(&jpID, stmt, JobProposalStatusCancelled, id); err != nil {
		return err
	}

	stmt = `
UPDATE job_proposals
SET status = $1,
	external_job_id = $2,
	updated_at = NOW()
WHERE id = $3;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, JobProposalStatusCancelled, nil, jpID)
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

// CreateSpec creates a new job proposal spec
func (o *orm) CreateSpec(spec JobProposalSpec, qopts ...pg.QOpt) (int64, error) {
	stmt := `
INSERT INTO job_proposal_specs (definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
RETURNING id;
`

	var id int64
	err := o.q.WithOpts(qopts...).Get(&id, stmt, spec.Definition, spec.Version, spec.Status, spec.JobProposalID)

	return id, errors.Wrap(err, "CreateJobProposalSpec failed")
}

// ExistsSpecByJobProposalIDAndVersion checks if a job proposal spec exists for
// a specific job proposal and version.
func (o *orm) ExistsSpecByJobProposalIDAndVersion(jpID int64, version int32, qopts ...pg.QOpt) (exists bool, err error) {
	stmt := `
SELECT exists (
	SELECT 1
	FROM job_proposal_specs
	WHERE job_proposal_id = $1 AND version = $2
);
`

	err = o.q.WithOpts(qopts...).Get(&exists, stmt, jpID, version)
	return exists, errors.Wrap(err, "JobProposalSpecVersionExists failed")
}

// GetSpec fetches the job proposal spec by id
func (o *orm) GetSpec(id int64, qopts ...pg.QOpt) (*JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE id = $1;
`
	var spec JobProposalSpec
	err := o.q.WithOpts(qopts...).Get(&spec, stmt, id)

	return &spec, errors.Wrap(err, "CreateJobProposalSpec failed")
}

// GetApprovedSpec gets the approved spec for a job proposal
func (o *orm) GetApprovedSpec(jpID int64, qopts ...pg.QOpt) (*JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE status = $1
AND job_proposal_id = $2
`

	var spec JobProposalSpec
	err := o.q.WithOpts(qopts...).Get(&spec, stmt, SpecStatusApproved, jpID)

	return &spec, errors.Wrap(err, "GetApprovedSpec failed")
}

// GetLatestSpec gets the latest spec for a job proposal.
func (o *orm) GetLatestSpec(jpID int64) (*JobProposalSpec, error) {
	stmt := `
	SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE (job_proposal_id, version) IN
(
	SELECT job_proposal_id, MAX(version)
	FROM job_proposal_specs
	GROUP BY job_proposal_id
)
AND job_proposal_id = $1
`

	var spec JobProposalSpec
	err := o.q.Get(&spec, stmt, jpID)

	return &spec, errors.Wrap(err, "GetLatestSpec failed")
}

// ListSpecsByJobProposalIDs lists the specs which belong to any of job proposal
// ids.
func (o *orm) ListSpecsByJobProposalIDs(ids []int64, qopts ...pg.QOpt) ([]JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE job_proposal_id = ANY($1)
`
	var specs []JobProposalSpec
	err := o.q.WithOpts(qopts...).Select(&specs, stmt, ids)
	return specs, errors.Wrap(err, "GetJobProposalsByManagersIDs failed")
}

// RejectSpec rejects the spec and updates the job proposal
func (o *orm) RejectSpec(id int64, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposal_specs
SET status = $1,
	status_updated_at = NOW(),
	updated_at = NOW()
WHERE id = $2
RETURNING job_proposal_id;
`

	var jpID int64
	if err := o.q.WithOpts(qopts...).Get(&jpID, stmt, JobProposalStatusRejected, id); err != nil {
		return err
	}

	stmt = `
UPDATE job_proposals
SET status = subquery.updateStatus,
	pending_update = FALSE,
	updated_at = NOW()
FROM (
	SELECT (CASE WHEN status = 'approved' THEN 'approved'::job_proposal_status ELSE 'rejected'::job_proposal_status END) as updateStatus
	FROM job_proposals
	WHERE id = $1
) as subquery
WHERE id = $2;
`

	result, err := o.q.WithOpts(qopts...).Exec(stmt, jpID, jpID)
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

// UpdateSpecDefinition updates the definition of a job proposal spec by id.
func (o *orm) UpdateSpecDefinition(id int64, spec string, qopts ...pg.QOpt) error {
	stmt := `
UPDATE job_proposal_specs
SET definition = $1,
	updated_at = NOW()
WHERE id = $2;
`

	res, err := o.q.WithOpts(qopts...).Exec(stmt, spec, id)
	if err != nil {
		return errors.Wrap(err, "UpdateSpecDefinition failed to update definition")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "UpdateSpecDefinition failed to get RowsAffected")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
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
