package syncer

import (
	"context"
	"encoding/base64"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ORM interface {
	// GetSecretsURLByID returns the secrets URL for the given ID.
	GetSecretsURLByID(ctx context.Context, id int64) (string, error)

	// GetContents returns the contents of the secret at the given plain URL.
	GetContents(ctx context.Context, url string) (string, error)

	// Update updates the contents of the secret at the given plain URL or inserts a new record if not found.
	Update(ctx context.Context, secretsURL, contents string) (int64, error)

	// SecretsFor returns a map of secrets for the given workflowOwner and workflowName.
	SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error)
}

type WorkflowRegistryDS = ORM

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func NewWorkflowRegistryDS(ds sqlutil.DataSource, lggr logger.Logger) *orm {
	return &orm{
		ds:   ds,
		lggr: lggr,
	}
}

func (orm *orm) GetSecretsURLByID(ctx context.Context, id int64) (string, error) {
	var secretsURL string
	if err := orm.ds.GetContext(ctx, &secretsURL,
		`SELECT secrets_url FROM workflow_secrets WHERE workflow_secrets.id = $1`,
		id,
	); err != nil {
		return secretsURL, err
	}

	url, err := base64.StdEncoding.DecodeString(secretsURL)
	if err != nil {
		return secretsURL, err
	}
	return string(url), nil
}

func (orm *orm) GetContents(ctx context.Context, url string) (string, error) {
	encoded := base64.URLEncoding.EncodeToString([]byte(url))
	var contents string
	err := orm.ds.GetContext(ctx, &contents,
		`SELECT contents 
         FROM workflow_secrets 
         WHERE secrets_url = $1`,
		encoded,
	)

	if err != nil {
		return "", err // Return an empty Artifact struct and the error
	}

	return contents, nil // Return the populated Artifact struct
}

// Update updates the contents of the secret at the given URL hash or inserts a new record if not found.
func (orm *orm) Update(ctx context.Context, secretsURL, contents string) (int64, error) {
	encoded := base64.URLEncoding.EncodeToString([]byte(secretsURL))
	var id int64
	err := orm.ds.QueryRowxContext(ctx,
		`INSERT INTO workflow_secrets (secrets_url, contents)
         VALUES ($1, $2)
         ON CONFLICT (secrets_url) DO UPDATE
         SET secrets_url = EXCLUDED.secrets_url, contents = EXCLUDED.contents
         RETURNING id`,
		encoded, contents,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (orm *orm) SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error) {
	return map[string]string{}, nil
}
