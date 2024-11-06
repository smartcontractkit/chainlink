package secrets

import (
	"context"
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"golang.org/x/crypto/sha3"
)

type ORM interface {
	// GetSecretsURL returns the original URL for a given hash value.  Fails if hash does not exist.
	GetSecretsURL(ctx context.Context, hash string) (string, error)

	// Update updates the contents of the secret at the given URL hash or inserts a new record if not found.
	Update(ctx context.Context, secretsURL, contents string) (int64, error)

	// SecretsFor returns a map of secrets for the given workflowOwner and workflowName.
	SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error)
}

type orm struct {
	ds   sqlutil.DataSource
	lggr logger.Logger
}

var _ ORM = (*orm)(nil)

func (orm *orm) GetSecretsURL(ctx context.Context, hash string) (string, error) {
	var secretsURL string
	if err := orm.ds.GetContext(ctx, &secretsURL,
		`SELECT secrets_url FROM workflow_artifacts WHERE workflow_artifacts.secrets_url_hash = $1`,
		hash,
	); err != nil {
		return secretsURL, err
	}
	return secretsURL, nil
}

type Artifact struct {
	SecretsURL string
	Contents   string
}

// GetArtifactByHash retrieves the artifact with the given secrets_url_hash,
// returning the secrets URL and contents.
func (orm *orm) GetArtifactByHash(ctx context.Context, hash string) (Artifact, error) {
	var artifact Artifact

	err := orm.ds.GetContext(ctx, &artifact,
		`SELECT secrets_url, contents 
         FROM workflow_artifacts 
         WHERE secrets_url_hash = $1`,
		hash,
	)

	if err != nil {
		return Artifact{}, err // Return an empty Artifact struct and the error
	}

	return artifact, nil // Return the populated Artifact struct
}

// Update updates the contents of the secret at the given URL hash or inserts a new record if not found.
func (orm *orm) Update(ctx context.Context, secretsURL, contents string) (int64, error) {
	var id int64
	err := orm.ds.QueryRowxContext(ctx,
		`INSERT INTO workflow_artifacts (secrets_url_hash, secrets_url, contents)
         VALUES ($1, $2, $3)
         ON CONFLICT (secrets_url_hash) DO UPDATE
         SET secrets_url = EXCLUDED.secrets_url, contents = EXCLUDED.contents
         RETURNING id`,
		keccak256Hash(secretsURL), secretsURL, contents,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (orm *orm) SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error) {
	return map[string]string{}, nil
}

func keccak256Hash(data string) string {
	// Create a new Keccak-256 hash instance
	hasher := sha3.NewLegacyKeccak256()

	// Write the input string to the hasher
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
