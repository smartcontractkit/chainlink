package feeds

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

//go:generate mockery --with-expecter=true --quiet --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	CountManagers(ctx context.Context) (int64, error)
	CreateManager(ctx context.Context, ms *FeedsManager) (int64, error)
	GetManager(ctx context.Context, id int64) (*FeedsManager, error)
	ListManagers(ctx context.Context) (mgrs []FeedsManager, err error)
	ListManagersByIDs(ctx context.Context, ids []int64) ([]FeedsManager, error)
	UpdateManager(ctx context.Context, mgr FeedsManager) error

	CreateBatchChainConfig(ctx context.Context, cfgs []ChainConfig) ([]int64, error)
	CreateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error)
	DeleteChainConfig(ctx context.Context, id int64) (int64, error)
	GetChainConfig(ctx context.Context, id int64) (*ChainConfig, error)
	ListChainConfigsByManagerIDs(ctx context.Context, mgrIDs []int64) ([]ChainConfig, error)
	UpdateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error)

	CountJobProposals(ctx context.Context) (int64, error)
	CountJobProposalsByStatus(ctx context.Context) (counts *JobProposalCounts, err error)
	CreateJobProposal(ctx context.Context, jp *JobProposal) (int64, error)
	DeleteProposal(ctx context.Context, id int64) error
	GetJobProposal(ctx context.Context, id int64) (*JobProposal, error)
	GetJobProposalByRemoteUUID(ctx context.Context, uuid uuid.UUID) (*JobProposal, error)
	ListJobProposals(ctx context.Context) (jps []JobProposal, err error)
	ListJobProposalsByManagersIDs(ctx context.Context, ids []int64) ([]JobProposal, error)
	UpdateJobProposalStatus(ctx context.Context, id int64, status JobProposalStatus) error // NEEDED?
	UpsertJobProposal(ctx context.Context, jp *JobProposal) (int64, error)

	ApproveSpec(ctx context.Context, id int64, externalJobID uuid.UUID) error
	CancelSpec(ctx context.Context, id int64) error
	CreateSpec(ctx context.Context, spec JobProposalSpec) (int64, error)
	ExistsSpecByJobProposalIDAndVersion(ctx context.Context, jpID int64, version int32) (exists bool, err error)
	GetApprovedSpec(ctx context.Context, jpID int64) (*JobProposalSpec, error)
	GetLatestSpec(ctx context.Context, jpID int64) (*JobProposalSpec, error)
	GetSpec(ctx context.Context, id int64) (*JobProposalSpec, error)
	ListSpecsByJobProposalIDs(ctx context.Context, ids []int64) ([]JobProposalSpec, error)
	RejectSpec(ctx context.Context, id int64) error
	RevokeSpec(ctx context.Context, id int64) error
	UpdateSpecDefinition(ctx context.Context, id int64, spec string) error

	IsJobManaged(ctx context.Context, jobID int64) (bool, error)

	Transact(context.Context, func(ORM) error) error
	WithDataSource(sqlutil.DataSource) ORM
}

var _ ORM = &orm{}

type orm struct {
	ds sqlutil.DataSource
}

func NewORM(ds sqlutil.DataSource) *orm {
	return &orm{ds: ds}
}

func (o *orm) Transact(ctx context.Context, fn func(ORM) error) error {
	return sqlutil.Transact(ctx, o.WithDataSource, o.ds, nil, fn)
}

func (o *orm) WithDataSource(ds sqlutil.DataSource) ORM { return &orm{ds} }

// Count counts the number of feeds manager records.
func (o *orm) CountManagers(ctx context.Context) (count int64, err error) {
	stmt := `
SELECT COUNT(*)
FROM feeds_managers
	`

	err = o.ds.GetContext(ctx, &count, stmt)
	return count, errors.Wrap(err, "CountManagers failed")
}

// CreateManager creates a feeds manager.
func (o *orm) CreateManager(ctx context.Context, ms *FeedsManager) (id int64, err error) {
	stmt := `
INSERT INTO feeds_managers (name, uri, public_key, created_at, updated_at)
VALUES ($1,$2,$3,NOW(),NOW())
RETURNING id;
`
	err = o.ds.GetContext(ctx, &id, stmt, ms.Name, ms.URI, ms.PublicKey)

	return id, errors.Wrap(err, "CreateManager failed")
}

// CreateChainConfig creates a new chain config.
func (o *orm) CreateChainConfig(ctx context.Context, cfg ChainConfig) (id int64, err error) {
	stmt := `
INSERT INTO feeds_manager_chain_configs (feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
RETURNING id;
`

	err = o.ds.GetContext(ctx,
		&id,
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
func (o *orm) CreateBatchChainConfig(ctx context.Context, cfgs []ChainConfig) (ids []int64, err error) {
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

	err = o.ds.SelectContext(ctx,
		&ids,
		fmt.Sprintf(stmt, strings.Join(vStrs, ",")),
		vArgs...,
	)

	return ids, errors.Wrap(err, "CreateBatchChainConfig failed")
}

// DeleteChainConfig deletes a chain config.
func (o *orm) DeleteChainConfig(ctx context.Context, id int64) (int64, error) {
	stmt := `
DELETE FROM feeds_manager_chain_configs
WHERE id = $1
RETURNING id;
`

	var ccid int64
	err := o.ds.GetContext(ctx, &ccid, stmt, id)

	return ccid, errors.Wrap(err, "DeleteChainConfig failed")
}

// GetChainConfig fetches a chain config.
func (o *orm) GetChainConfig(ctx context.Context, id int64) (*ChainConfig, error) {
	stmt := `
SELECT id, feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at
FROM feeds_manager_chain_configs
WHERE id = $1;
`

	var cfg ChainConfig
	err := o.ds.GetContext(ctx, &cfg, stmt, id)

	return &cfg, errors.Wrap(err, "GetChainConfig failed")
}

// ListChainConfigsByManagerIDs fetches the chain configs matching all manager
// ids.
func (o *orm) ListChainConfigsByManagerIDs(ctx context.Context, mgrIDs []int64) ([]ChainConfig, error) {
	stmt := `
SELECT id, feeds_manager_id, chain_id, chain_type, account_address, admin_address, flux_monitor_config, ocr1_config, ocr2_config, created_at, updated_at
FROM feeds_manager_chain_configs
WHERE feeds_manager_id = ANY($1)
	`

	var cfgs []ChainConfig
	err := o.ds.SelectContext(ctx, &cfgs, stmt, mgrIDs)

	return cfgs, errors.Wrap(err, "ListJobProposalsByManagersIDs failed")
}

// UpdateChainConfig updates a chain config.
func (o *orm) UpdateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error) {
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
	err := o.ds.GetContext(ctx, &cfgID, stmt,
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
func (o *orm) GetManager(ctx context.Context, id int64) (mgr *FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers
WHERE id = $1
`

	mgr = new(FeedsManager)
	err = o.ds.GetContext(ctx, mgr, stmt, id)
	return mgr, errors.Wrap(err, "GetManager failed")
}

// ListManager lists all feeds managers.
func (o *orm) ListManagers(ctx context.Context) (mgrs []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers;
`

	err = o.ds.SelectContext(ctx, &mgrs, stmt)
	return mgrs, errors.Wrap(err, "ListManagers failed")
}

// ListManagersByIDs gets feeds managers by ids.
func (o *orm) ListManagersByIDs(ctx context.Context, ids []int64) (managers []FeedsManager, err error) {
	stmt := `
SELECT id, name, uri, public_key, created_at, updated_at
FROM feeds_managers
WHERE id = ANY($1)
ORDER BY created_at, id;`

	mgrIds := pq.Array(ids)
	err = o.ds.SelectContext(ctx, &managers, stmt, mgrIds)

	return managers, errors.Wrap(err, "GetManagers failed")
}

// UpdateManager updates the manager details.
func (o *orm) UpdateManager(ctx context.Context, mgr FeedsManager) (err error) {
	stmt := `
UPDATE feeds_managers
SET name = $1, uri = $2, public_key = $3, updated_at = NOW()
WHERE id = $4;
`

	res, err := o.ds.ExecContext(ctx, stmt, mgr.Name, mgr.URI, mgr.PublicKey, mgr.ID)
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
func (o *orm) CreateJobProposal(ctx context.Context, jp *JobProposal) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (name, remote_uuid, status, feeds_manager_id, multiaddrs, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id;
`

	err = o.ds.GetContext(ctx, &id, stmt, jp.Name, jp.RemoteUUID, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "CreateJobProposal failed")
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposals(ctx context.Context) (count int64, err error) {
	stmt := `SELECT COUNT(*) FROM job_proposals`

	err = o.ds.GetContext(ctx, &count, stmt)
	return count, errors.Wrap(err, "CountJobProposals failed")
}

// CountJobProposals counts the number of job proposal records.
func (o *orm) CountJobProposalsByStatus(ctx context.Context) (counts *JobProposalCounts, err error) {
	stmt := `
SELECT 
	COUNT(*) filter (where job_proposals.status = 'pending' OR job_proposals.pending_update = TRUE) as pending,
	COUNT(*) filter (where job_proposals.status = 'approved' AND job_proposals.pending_update = FALSE) as approved,
	COUNT(*) filter (where job_proposals.status = 'rejected' AND job_proposals.pending_update = FALSE) as rejected,
	COUNT(*) filter (where job_proposals.status = 'revoked' AND job_proposals.pending_update = FALSE) as revoked,
	COUNT(*) filter (where job_proposals.status = 'deleted' AND job_proposals.pending_update = FALSE) as deleted,
	COUNT(*) filter (where job_proposals.status = 'cancelled' AND job_proposals.pending_update = FALSE) as cancelled
FROM job_proposals;
	`

	counts = new(JobProposalCounts)
	err = o.ds.GetContext(ctx, counts, stmt)
	return counts, errors.Wrap(err, "CountJobProposalsByStatus failed")
}

// GetJobProposal gets a job proposal by id.
func (o *orm) GetJobProposal(ctx context.Context, id int64) (jp *JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE id = $1
`
	jp = new(JobProposal)
	err = o.ds.GetContext(ctx, jp, stmt, id)
	return jp, errors.Wrap(err, "GetJobProposal failed")
}

// GetJobProposalByRemoteUUID gets a job proposal by the remote FMS uuid. This
// method will filter out the deleted job proposals. To get all job proposals,
// use the GetJobProposal get by id method.
func (o *orm) GetJobProposalByRemoteUUID(ctx context.Context, id uuid.UUID) (jp *JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE remote_uuid = $1
AND status <> $2;
`

	jp = new(JobProposal)
	err = o.ds.GetContext(ctx, jp, stmt, id, JobProposalStatusDeleted)
	return jp, errors.Wrap(err, "GetJobProposalByRemoteUUID failed")
}

// ListJobProposals lists all job proposals.
func (o *orm) ListJobProposals(ctx context.Context) (jps []JobProposal, err error) {
	stmt := `
SELECT *
FROM job_proposals;
`

	err = o.ds.SelectContext(ctx, &jps, stmt)
	return jps, errors.Wrap(err, "ListJobProposals failed")
}

// ListJobProposalsByManagersIDs gets job proposals by feeds managers IDs.
func (o *orm) ListJobProposalsByManagersIDs(ctx context.Context, ids []int64) ([]JobProposal, error) {
	stmt := `
SELECT *
FROM job_proposals
WHERE feeds_manager_id = ANY($1)
`
	var jps []JobProposal
	err := o.ds.SelectContext(ctx, &jps, stmt, ids)
	return jps, errors.Wrap(err, "ListJobProposalsByManagersIDs failed")
}

// UpdateJobProposalStatus updates the status of a job proposal by id.
func (o *orm) UpdateJobProposalStatus(ctx context.Context, id int64, status JobProposalStatus) error {
	stmt := `
UPDATE job_proposals
SET status = $1,
	updated_at = NOW()
WHERE id = $2;
`

	result, err := o.ds.ExecContext(ctx, stmt, status, id)
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
func (o *orm) UpsertJobProposal(ctx context.Context, jp *JobProposal) (id int64, err error) {
	stmt := `
INSERT INTO job_proposals (name, remote_uuid, status, feeds_manager_id, multiaddrs, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (remote_uuid)
DO
	UPDATE SET
		pending_update = TRUE,
		name = EXCLUDED.name,
		status = (
			CASE
				WHEN job_proposals.status = 'deleted' THEN 'deleted'::job_proposal_status
				WHEN job_proposals.status = 'approved' THEN 'approved'::job_proposal_status
				ELSE EXCLUDED.status
			END
		),
		multiaddrs = EXCLUDED.multiaddrs,
		updated_at = EXCLUDED.updated_at
RETURNING id;
`

	err = o.ds.GetContext(ctx, &id, stmt, jp.Name, jp.RemoteUUID, jp.Status, jp.FeedsManagerID, jp.Multiaddrs)
	return id, errors.Wrap(err, "UpsertJobProposal")
}

// ApproveSpec approves the spec and sets the external job ID on the associated
// job proposal.
func (o *orm) ApproveSpec(ctx context.Context, id int64, externalJobID uuid.UUID) error {
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
	if err := o.ds.GetContext(ctx, &jpID, stmt, JobProposalStatusApproved, id); err != nil {
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

	result, err := o.ds.ExecContext(ctx, stmt, JobProposalStatusApproved, externalJobID, jpID)
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

// CancelSpec cancels the spec and removes the external job id from the associated job proposal. It
// sets the status of the spec and the proposal to cancelled, except in the case of deleted
// proposals.
func (o *orm) CancelSpec(ctx context.Context, id int64) error {
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
	if err := o.ds.GetContext(ctx, &jpID, stmt, SpecStatusCancelled, id); err != nil {
		return err
	}

	stmt = `
UPDATE job_proposals
SET status = (
		CASE
			WHEN status = 'deleted' THEN 'deleted'::job_proposal_status
			ELSE 'cancelled'::job_proposal_status
		END
	),
	pending_update = FALSE,
	external_job_id = $2,
	updated_at = NOW()
WHERE id = $1;
`
	result, err := o.ds.ExecContext(ctx, stmt, jpID, nil)
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
func (o *orm) CreateSpec(ctx context.Context, spec JobProposalSpec) (int64, error) {
	stmt := `
INSERT INTO job_proposal_specs (definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
RETURNING id;
`

	var id int64
	err := o.ds.GetContext(ctx, &id, stmt, spec.Definition, spec.Version, spec.Status, spec.JobProposalID)

	return id, errors.Wrap(err, "CreateJobProposalSpec failed")
}

// ExistsSpecByJobProposalIDAndVersion checks if a job proposal spec exists for a specific job
// proposal and version.
func (o *orm) ExistsSpecByJobProposalIDAndVersion(ctx context.Context, jpID int64, version int32) (exists bool, err error) {
	stmt := `
SELECT exists (
	SELECT 1
	FROM job_proposal_specs
	WHERE job_proposal_id = $1 AND version = $2
);
`

	err = o.ds.GetContext(ctx, &exists, stmt, jpID, version)
	return exists, errors.Wrap(err, "JobProposalSpecVersionExists failed")
}

// DeleteProposal performs a soft delete of the job proposal by setting the status to deleted
func (o *orm) DeleteProposal(ctx context.Context, id int64) error {
	// Get the latest spec for the proposal.
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
	err := o.ds.GetContext(ctx, &spec, stmt, id)
	if err != nil {
		return err
	}

	// Set pending update to true only if the latest proposal is approved so that any running jobs
	// are reminded to be cancelled.
	pendingUpdate := spec.Status == SpecStatusApproved
	stmt = `
UPDATE job_proposals
SET status = $1,
    pending_update = $3,
    updated_at = NOW()
WHERE id = $2;
`

	result, err := o.ds.ExecContext(ctx, stmt, JobProposalStatusDeleted, id, pendingUpdate)
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

// GetSpec fetches the job proposal spec by id
func (o *orm) GetSpec(ctx context.Context, id int64) (*JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE id = $1;
`
	var spec JobProposalSpec
	err := o.ds.GetContext(ctx, &spec, stmt, id)

	return &spec, errors.Wrap(err, "CreateJobProposalSpec failed")
}

// GetApprovedSpec gets the approved spec for a job proposal
func (o *orm) GetApprovedSpec(ctx context.Context, jpID int64) (*JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE status = $1
AND job_proposal_id = $2
`

	var spec JobProposalSpec
	err := o.ds.GetContext(ctx, &spec, stmt, SpecStatusApproved, jpID)

	return &spec, errors.Wrap(err, "GetApprovedSpec failed")
}

// GetLatestSpec gets the latest spec for a job proposal.
func (o *orm) GetLatestSpec(ctx context.Context, jpID int64) (*JobProposalSpec, error) {
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
	err := o.ds.GetContext(ctx, &spec, stmt, jpID)

	return &spec, errors.Wrap(err, "GetLatestSpec failed")
}

// ListSpecsByJobProposalIDs lists the specs which belong to any of job proposal
// ids.
func (o *orm) ListSpecsByJobProposalIDs(ctx context.Context, ids []int64) ([]JobProposalSpec, error) {
	stmt := `
SELECT id, definition, version, status, job_proposal_id, status_updated_at, created_at, updated_at
FROM job_proposal_specs
WHERE job_proposal_id = ANY($1)
`
	var specs []JobProposalSpec
	err := o.ds.SelectContext(ctx, &specs, stmt, ids)
	return specs, errors.Wrap(err, "GetJobProposalsByManagersIDs failed")
}

// RejectSpec rejects the spec and updates the job proposal
func (o *orm) RejectSpec(ctx context.Context, id int64) error {
	stmt := `
UPDATE job_proposal_specs
SET status = $1,
	status_updated_at = NOW(),
	updated_at = NOW()
WHERE id = $2
RETURNING job_proposal_id;
`

	var jpID int64
	if err := o.ds.GetContext(ctx, &jpID, stmt, SpecStatusRejected, id); err != nil {
		return err
	}

	stmt = `
UPDATE job_proposals
SET status = (
		CASE
			WHEN status = 'approved' THEN 'approved'::job_proposal_status
			WHEN status = 'deleted' THEN 'deleted'::job_proposal_status
			ELSE 'rejected'::job_proposal_status
		END
	),
	pending_update = FALSE,
	updated_at = NOW()
WHERE id = $1
`

	result, err := o.ds.ExecContext(ctx, stmt, jpID)
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

// RevokeSpec revokes a job proposal with a pending job spec. An approved
// proposal cannot be revoked. A revoked proposal's job spec cannot be approved
// or edited, but the job can be reproposed by FMS.
func (o *orm) RevokeSpec(ctx context.Context, id int64) error {
	// Update the status of the spec
	stmt := `
UPDATE job_proposal_specs
SET status = (
		CASE
			WHEN status = 'approved' THEN 'approved'::job_proposal_spec_status
			ELSE $2
		END
	),
	status_updated_at = NOW(),
	updated_at = NOW()
WHERE id = $1
RETURNING job_proposal_id;
`

	var jpID int64
	if err := o.ds.GetContext(ctx, &jpID, stmt, id, SpecStatusRevoked); err != nil {
		return err
	}

	stmt = `
UPDATE job_proposals
SET status = (
		CASE
			WHEN status = 'deleted' THEN 'deleted'::job_proposal_status
			WHEN status = 'approved' THEN 'approved'::job_proposal_status
			ELSE $3
		END
	),
	pending_update = FALSE,
	external_job_id = (
		CASE
			WHEN status <> 'approved' THEN $2
			ELSE job_proposals.external_job_id
		END
	),
	updated_at = NOW()
WHERE id = $1
	`

	result, err := o.ds.ExecContext(ctx, stmt, jpID, nil, JobProposalStatusRevoked)
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
func (o *orm) UpdateSpecDefinition(ctx context.Context, id int64, spec string) error {
	stmt := `
UPDATE job_proposal_specs
SET definition = $1,
	updated_at = NOW()
WHERE id = $2;
`

	res, err := o.ds.ExecContext(ctx, stmt, spec, id)
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
func (o *orm) IsJobManaged(ctx context.Context, jobID int64) (exists bool, err error) {
	stmt := `
SELECT exists (
	SELECT 1
	FROM job_proposals
	INNER JOIN jobs ON job_proposals.external_job_id = jobs.external_job_id
	WHERE jobs.id = $1
);
`

	err = o.ds.GetContext(ctx, &exists, stmt, jobID)
	return exists, errors.Wrap(err, "IsJobManaged failed")
}
