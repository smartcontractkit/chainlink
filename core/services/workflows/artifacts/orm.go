package artifacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

type WorkflowSecretsDS interface {
	// GetSecretsURLByID returns the secrets URL for the given ID.
	GetSecretsURLByID(ctx context.Context, id int64) (string, error)

	// GetSecretsURLByID returns the secrets URL for the given ID.
	GetSecretsURLByHash(ctx context.Context, hash string) (string, error)

	// GetContents returns the contents of the secret at the given plain URL.
	GetContents(ctx context.Context, url string) (string, error)

	// GetContentsByHash returns the contents of the secret at the given hashed URL.
	GetContentsByHash(ctx context.Context, hash string) (string, error)

	// GetContentsByWorkflowID returns the contents and secrets_url of the secret for the given workflow.
	GetContentsByWorkflowID(ctx context.Context, workflowID string) (string, string, error)

	// GetSecretsURLHash returns the keccak256 hash of the owner and secrets URL.
	GetSecretsURLHash(owner, secretsURL []byte) ([]byte, error)

	// Update updates the contents of the secrets at the given plain URL or inserts a new record if not found.
	Update(ctx context.Context, secretsURL, contents string) (int64, error)

	Create(ctx context.Context, secretsURL, hash, contents string) (int64, error)
}

type WorkflowSpecsDS interface {
	// UpsertWorkflowSpec inserts or updates a workflow spec.  Updates on conflict of workflow name
	// and owner
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)

	// UpsertWorkflowSpecWithSecrets inserts or updates a workflow spec with secrets in a transaction.
	// Updates on conflict of workflow name and owner.
	UpsertWorkflowSpecWithSecrets(ctx context.Context, spec *job.WorkflowSpec, url, hash, contents string) (int64, error)

	// GetWorkflowSpec returns the workflow spec for the given owner and name.
	GetWorkflowSpec(ctx context.Context, owner, name string) (*job.WorkflowSpec, error)

	// DeleteWorkflowSpec deletes the workflow spec for the given owner and name.
	DeleteWorkflowSpec(ctx context.Context, owner, name string) error

	// GetWorkflowSpecByID returns the workflow spec for the given workflowID.
	GetWorkflowSpecByID(ctx context.Context, id string) (*job.WorkflowSpec, error)
}

type ORM interface {
	WorkflowSecretsDS
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

func (orm *orm) GetSecretsURLByID(ctx context.Context, id int64) (string, error) {
	var secretsURL string
	err := orm.ds.GetContext(ctx, &secretsURL,
		`SELECT secrets_url FROM workflow_secrets WHERE workflow_secrets.id = $1`,
		id,
	)

	return secretsURL, err
}

func (orm *orm) GetSecretsURLByHash(ctx context.Context, hash string) (string, error) {
	var secretsURL string
	err := orm.ds.GetContext(ctx, &secretsURL,
		`SELECT secrets_url FROM workflow_secrets WHERE workflow_secrets.secrets_url_hash = $1`,
		hash,
	)

	return secretsURL, err
}

func (orm *orm) GetContentsByHash(ctx context.Context, hash string) (string, error) {
	var contents string
	err := orm.ds.GetContext(ctx, &contents,
		`SELECT contents 
         FROM workflow_secrets 
         WHERE secrets_url_hash = $1`,
		hash,
	)

	if err != nil {
		return "", err // Return an empty Artifact struct and the error
	}

	return contents, nil // Return the populated Artifact struct
}

func (orm *orm) GetContents(ctx context.Context, url string) (string, error) {
	var contents string
	err := orm.ds.GetContext(ctx, &contents,
		`SELECT contents 
         FROM workflow_secrets 
         WHERE secrets_url = $1`,
		url,
	)

	if err != nil {
		return "", err // Return an empty Artifact struct and the error
	}

	return contents, nil // Return the populated Artifact struct
}

type Int struct {
	sql.NullInt64
}

type joinRecord struct {
	SecretsID      sql.NullString `db:"wspec_secrets_id"`
	SecretsURLHash sql.NullString `db:"wsec_secrets_url_hash"`
	Contents       sql.NullString `db:"wsec_contents"`
}

var ErrEmptySecrets = errors.New("secrets field is empty")

// GetContentsByWorkflowID joins the workflow_secrets on the workflow_specs table and gets
// the associated secrets contents.
func (orm *orm) GetContentsByWorkflowID(ctx context.Context, workflowID string) (string, string, error) {
	var jr joinRecord
	err := orm.ds.GetContext(
		ctx,
		&jr,
		`SELECT wsec.secrets_url_hash AS wsec_secrets_url_hash, wsec.contents AS wsec_contents, wspec.secrets_id AS wspec_secrets_id
	FROM workflow_specs AS wspec
	LEFT JOIN
		workflow_secrets AS wsec ON wspec.secrets_id = wsec.id
	WHERE wspec.workflow_id = $1`,
		workflowID,
	)
	if err != nil {
		return "", "", err
	}

	if !jr.SecretsID.Valid {
		return "", "", ErrEmptySecrets
	}

	if jr.Contents.String == "" {
		return "", "", ErrEmptySecrets
	}

	return jr.SecretsURLHash.String, jr.Contents.String, nil
}

// Update updates the secrets content at the given hash or inserts a new record if not found.
func (orm *orm) Update(ctx context.Context, hash, contents string) (int64, error) {
	var id int64
	err := orm.ds.QueryRowxContext(ctx,
		`INSERT INTO workflow_secrets (secrets_url_hash, contents)
         VALUES ($1, $2)
         ON CONFLICT (secrets_url_hash) DO UPDATE
         SET secrets_url_hash = EXCLUDED.secrets_url_hash, contents = EXCLUDED.contents
         RETURNING id`,
		hash, contents,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update updates the secrets content at the given hash or inserts a new record if not found.
func (orm *orm) Create(ctx context.Context, url, hash, contents string) (int64, error) {
	var id int64
	err := orm.ds.QueryRowxContext(ctx,
		`INSERT INTO workflow_secrets (secrets_url, secrets_url_hash, contents)
         VALUES ($1, $2, $3)
         RETURNING id`,
		url, hash, contents,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (orm *orm) GetSecretsURLHash(owner, secretsURL []byte) ([]byte, error) {
	return crypto.Keccak256(append(owner, secretsURL...))
}

func (orm *orm) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	var id int64
	err := sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		txErr := tx.QueryRowxContext(
			ctx,
			`DELETE FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_id != $3`,
			spec.WorkflowOwner,
			spec.WorkflowName,
			spec.WorkflowID,
		).Scan(nil)
		if txErr != nil && !errors.Is(txErr, sql.ErrNoRows) {
			return fmt.Errorf("failed to clean up previous workflow specs: %w", txErr)
		}

		query := `
			INSERT INTO workflow_specs (
				workflow,
				config,
				workflow_id,
				workflow_owner,
				workflow_name,
				status,
				binary_url,
				config_url,
				secrets_id,
				created_at,
				updated_at,
				spec_type
			) VALUES (
				:workflow,
				:config,
				:workflow_id,
				:workflow_owner,
				:workflow_name,
				:status,
				:binary_url,
				:config_url,
				:secrets_id,
				:created_at,
				:updated_at,
				:spec_type
			) ON CONFLICT (workflow_owner, workflow_name) DO UPDATE
			SET
				workflow = EXCLUDED.workflow,
				config = EXCLUDED.config,
				workflow_id = EXCLUDED.workflow_id,
				workflow_owner = EXCLUDED.workflow_owner,
				workflow_name = EXCLUDED.workflow_name,
				status = EXCLUDED.status,
				binary_url = EXCLUDED.binary_url,
				config_url = EXCLUDED.config_url,
				secrets_id = EXCLUDED.secrets_id,
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

func (orm *orm) UpsertWorkflowSpecWithSecrets(
	ctx context.Context,
	spec *job.WorkflowSpec, url, hash, contents string) (int64, error) {
	var id int64
	err := sqlutil.TransactDataSource(ctx, orm.ds, nil, func(tx sqlutil.DataSource) error {
		var sid int64
		txErr := tx.QueryRowxContext(ctx,
			`INSERT INTO workflow_secrets (secrets_url, secrets_url_hash, contents)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (secrets_url_hash) DO UPDATE
         	 SET 
			 	secrets_url_hash = EXCLUDED.secrets_url_hash, 
				contents = EXCLUDED.contents,
				secrets_url = EXCLUDED.secrets_url
			 RETURNING id`,
			url, hash, contents,
		).Scan(&sid)

		if txErr != nil {
			return fmt.Errorf("failed to create workflow secrets: %w", txErr)
		}

		txErr = tx.QueryRowxContext(
			ctx,
			`DELETE FROM workflow_specs WHERE workflow_owner = $1 AND workflow_name = $2 AND workflow_id != $3`,
			spec.WorkflowOwner,
			spec.WorkflowName,
			spec.WorkflowID,
		).Scan(nil)
		if txErr != nil && !errors.Is(txErr, sql.ErrNoRows) {
			return fmt.Errorf("failed to clean up previous workflow specs: %w", txErr)
		}

		spec.SecretsID = sql.NullInt64{Int64: sid, Valid: true}

		query := `
			INSERT INTO workflow_specs (
				workflow,
				config,
				workflow_id,
				workflow_owner,
				workflow_name,
				status,
				binary_url,
				config_url,
				secrets_id,
				created_at,
				updated_at,
				spec_type
			) VALUES (
				:workflow,
				:config,
				:workflow_id,
				:workflow_owner,
				:workflow_name,
				:status,
				:binary_url,
				:config_url,
				:secrets_id,
				:created_at,
				:updated_at,
				:spec_type
			) ON CONFLICT (workflow_owner, workflow_name) DO UPDATE
			SET
				workflow = EXCLUDED.workflow,
				config = EXCLUDED.config,
				workflow_id = EXCLUDED.workflow_id,
				workflow_owner = EXCLUDED.workflow_owner,
				workflow_name = EXCLUDED.workflow_name,
				status = EXCLUDED.status,
				binary_url = EXCLUDED.binary_url,
				config_url = EXCLUDED.config_url,
				created_at = EXCLUDED.created_at,
				updated_at = EXCLUDED.updated_at,
				spec_type = EXCLUDED.spec_type,
				secrets_id = EXCLUDED.secrets_id
			RETURNING id
		`

		stmt, txErr := tx.PrepareNamedContext(ctx, query)
		if txErr != nil {
			return txErr
		}
		defer stmt.Close()

		spec.UpdatedAt = time.Now()
		return stmt.QueryRowxContext(ctx, spec).Scan(&id)
	})
	return id, err
}

func (orm *orm) GetWorkflowSpec(ctx context.Context, owner, name string) (*job.WorkflowSpec, error) {
	query := `
		SELECT *
		FROM workflow_specs
		WHERE workflow_owner = $1 AND workflow_name = $2
	`

	var spec job.WorkflowSpec
	err := orm.ds.GetContext(ctx, &spec, query, owner, name)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func (orm *orm) GetWorkflowSpecByID(ctx context.Context, id string) (*job.WorkflowSpec, error) {
	query := `
		SELECT *
		FROM workflow_specs
		WHERE workflow_id = $1
	`

	var spec job.WorkflowSpec
	err := orm.ds.GetContext(ctx, &spec, query, id)
	if err != nil {
		return nil, err
	}

	return &spec, nil
}

func (orm *orm) DeleteWorkflowSpec(ctx context.Context, owner, name string) error {
	query := `
		DELETE FROM workflow_specs
		WHERE workflow_owner = $1 AND workflow_name = $2
	`

	result, err := orm.ds.ExecContext(ctx, query, owner, name)
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
