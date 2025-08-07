package v2

import (
	"context"
	"database/sql"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type WorkflowSpecsDS interface {
	// UpsertWorkflowSpec inserts or updates a workflow spec. Multiple workflow specs can exist per owner/name combination; unique by workflow ID.
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)

	// GetWorkflowSpecByID returns the workflow spec for the given workflowID.
	GetWorkflowSpec(ctx context.Context, id string) (*job.WorkflowSpec, error)

	// DeleteWorkflowSpec deletes the workflow spec for the given workflow ID.
	DeleteWorkflowSpec(ctx context.Context, id string) error
}

type ORM interface {
	WorkflowSpecsDS
}

type WorkflowRegistryDS = ORM

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ WorkflowRegistryDS = (*orm)(nil)

func NewWorkflowRegistryDS(ds sqlutil.DataSource, lggr logger.Logger) *orm {
	return &orm{
		ds:   ds,
		lggr: lggr,
	}
}

// UpsertWorkflowSpec inserts or updates a workflow spec. Unique by workflow ID. Multiple workflow specs can exists per owner/name combination.
func (orm *orm) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	var id int64
	err := sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		query := `
			INSERT INTO workflow_specs_v2 (
				workflow,
				config,
				workflow_id,
				workflow_owner,
				workflow_name,
				workflow_tag,
				status,
				binary_url,
				config_url,
				created_at,
				updated_at,
				spec_type
			) VALUES (
				:workflow,
				:config,
				:workflow_id,
				:workflow_owner,
				:workflow_name,
				:workflow_tag,
				:status,
				:binary_url,
				:config_url,
				:created_at,
				:updated_at,
				:spec_type
			) ON CONFLICT (workflow_id) DO UPDATE
			SET
				workflow = EXCLUDED.workflow,
				config = EXCLUDED.config,
				workflow_owner = EXCLUDED.workflow_owner,
				workflow_name = EXCLUDED.workflow_name,
				workflow_tag = EXCLUDED.workflow_tag,
				status = EXCLUDED.status,
				binary_url = EXCLUDED.binary_url,
				config_url = EXCLUDED.config_url,
				created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at,
				spec_type = EXCLUDED.spec_type
			RETURNING id
		`

		stmt, err := orm.ds.PrepareNamedContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		spec.UpdatedAt = time.Now()
		return stmt.QueryRowxContext(ctx, spec).Scan(&id)
	})

	return id, err
}

func (orm *orm) GetWorkflowSpec(ctx context.Context, id string) (*job.WorkflowSpec, error) {
	query := `
		SELECT *
		FROM workflow_specs_v2
		WHERE workflow_id = $1
	`

	var spec job.WorkflowSpec
	err := orm.ds.GetContext(ctx, &spec, query, id)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func (orm *orm) DeleteWorkflowSpec(ctx context.Context, id string) error {
	query := `
		DELETE FROM workflow_specs_v2
		WHERE workflow_id = $1
	`

	result, err := orm.ds.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No spec deleted
	}

	return nil
}
