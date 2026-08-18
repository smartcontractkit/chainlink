package v2

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type WorkflowSpecsDS interface {
	// UpsertWorkflowSpec inserts or updates a workflow spec. Multiple workflow specs can exist per owner/name combination; unique by workflow ID.
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)

	// GetWorkflowSpecByID returns the workflow spec for the given workflowID.
	GetWorkflowSpec(ctx context.Context, id string) (*job.WorkflowSpec, error)

	// ListWorkflowSpecs returns the persisted workflow specs. It projects only
	// identity columns (workflow_id, workflow_owner, registered_at);
	// other fields are left zero to keep returned batch size small
	ListWorkflowSpecs(ctx context.Context) ([]*job.WorkflowSpec, error)

	// DeleteWorkflowSpec deletes the spec row for id. Idempotent: a missing row
	// returns (nil, nil), not an error. The deleted row is returned so the
	// caller can derive the generation-scoped metering event_id from the row's
	// own persisted registered_at (every node emits the identical id regardless
	// of local staleness).
	DeleteWorkflowSpec(ctx context.Context, id string) (*job.WorkflowSpec, error)

	// PauseWorkflowSpec tombstones the spec row: status becomes paused and the
	// heavy artifact payload (hex binary + config) is cleared in the same
	// statement, freeing storage while the row itself remains the durable record
	// that this registration generation is still held (level-neutral for
	// metering). Idempotent; pausing an absent or already-paused row is a no-op.
	PauseWorkflowSpec(ctx context.Context, id string) error

	// DeleteWorkflowSpecs deletes workflow specs for the given workflow IDs in a single query.
	DeleteWorkflowSpecs(ctx context.Context, ids []string) error
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
				spec_type,
				attributes,
				registered_at,
				source
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
				:spec_type,
				:attributes,
				:registered_at,
				:source
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
				spec_type = EXCLUDED.spec_type,
				attributes = EXCLUDED.attributes,
				registered_at = EXCLUDED.registered_at,
				source = EXCLUDED.source
			RETURNING id
		`

		now := time.Now().UTC()
		spec.UpdatedAt = now
		if spec.CreatedAt.IsZero() {
			spec.CreatedAt = now
		}
		q, args, namedErr := sqlx.Named(query, spec)
		if namedErr != nil {
			return namedErr
		}
		q = sqlx.Rebind(sqlx.DOLLAR, q)
		return tx.QueryRowxContext(ctx, q, args...).Scan(&id)
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
	// Note: "Get will return sql.ErrNoRows like row.Scan would" - sqlx@v1.4.0
	err := orm.ds.GetContext(ctx, &spec, query, id)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

// ListWorkflowSpecs returns all persisted workflow specs, projecting the
// identity columns needed by the orphan sweep and the metering snapshot path.
func (orm *orm) ListWorkflowSpecs(ctx context.Context) ([]*job.WorkflowSpec, error) {
	query := `
		SELECT workflow_id, workflow_owner, registered_at, source
		FROM workflow_specs_v2
	`

	var specs []*job.WorkflowSpec
	if err := orm.ds.SelectContext(ctx, &specs, query); err != nil {
		return nil, err
	}

	return specs, nil
}

func (orm *orm) DeleteWorkflowSpec(ctx context.Context, id string) (*job.WorkflowSpec, error) {
	query := `DELETE FROM workflow_specs_v2 WHERE workflow_id = $1
		RETURNING id, workflow, config, workflow_id, workflow_owner, workflow_name,
			workflow_tag, status, binary_url, config_url, created_at, updated_at,
			spec_type, attributes, registered_at, source`
	var spec job.WorkflowSpec
	err := orm.ds.GetContext(ctx, &spec, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func (orm *orm) PauseWorkflowSpec(ctx context.Context, id string) error {
	query := `UPDATE workflow_specs_v2
		SET status = $2, workflow = '', config = '', updated_at = now()
		WHERE workflow_id = $1`
	_, err := orm.ds.ExecContext(ctx, query, id, job.WorkflowSpecStatusPaused)
	return err
}

func (orm *orm) DeleteWorkflowSpecs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query := `DELETE FROM workflow_specs_v2 WHERE workflow_id = ANY($1)`
	_, err := orm.ds.ExecContext(ctx, query, pq.Array(ids))
	return err
}
