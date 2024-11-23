package syncer

import (
	"context"
	"errors"

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

	// GetSecretsURLHash returns the keccak256 hash of the owner and secrets URL.
	GetSecretsURLHash(owner, secretsURL []byte) ([]byte, error)

	// Update updates the contents of the secrets at the given plain URL or inserts a new record if not found.
	Update(ctx context.Context, secretsURL, contents string) (int64, error)

	Create(ctx context.Context, secretsURL, hash, contents string) (int64, error)
}

type WorkflowSpecsDS interface {
	CreateWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)
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

func (orm *orm) CreateWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	return 0, errors.New("not implemented")
}
