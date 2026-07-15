package fakes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	confidentialhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestHasEncryptionSecret(t *testing.T) {
	t.Run("returns true when magic key exists", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: AESGCMEncryptionKeyName},
		}
		assert.True(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false when magic key does not exist", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: "another-key"},
		}
		assert.False(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false for empty secrets", func(t *testing.T) {
		assert.False(t, hasEncryptionSecret(nil))
		assert.False(t, hasEncryptionSecret([]*confidentialhttp.SecretIdentifier{}))
	})
}

func TestDirectConfidentialHTTPAction_SecretsLoading(t *testing.T) {
	// Create a temporary secrets.yaml
	tmpDir := t.TempDir()
	secretsFile := filepath.Join(tmpDir, "secrets.yaml")
	secretsContent := `
secretsNames:
  API_KEY:
    - API_KEY_ALL
`
	err := os.WriteFile(secretsFile, []byte(secretsContent), 0600)
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Remove(secretsFile)) }()

	// Set environment variables
	t.Setenv("API_KEY_ALL", "resolved-api-key")

	lggr := logger.Test(t)
	action := NewDirectConfidentialHTTPAction(lggr, secretsFile)

	// Verify secrets were resolved
	secrets, ok := action.secretsConfig.SecretsNames["API_KEY"]
	require.True(t, ok)
	require.Len(t, secrets, 1)
	assert.Equal(t, "resolved-api-key", secrets[0])
}

func TestDirectConfidentialHTTPAction_SecretsLoading_MissingEnv(t *testing.T) {
	// Create a temporary secrets.yaml
	tmpDir := t.TempDir()
	secretsFile := filepath.Join(tmpDir, "secrets.yaml")
	secretsContent := `
secretsNames:
  API_KEY:
    - MISSING_ENV_VAR
`
	err := os.WriteFile(secretsFile, []byte(secretsContent), 0600)
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Remove(secretsFile)) }()

	lggr := logger.Test(t)
	action := NewDirectConfidentialHTTPAction(lggr, secretsFile)

	// Verify secret remains as the placeholder when env var is missing
	secrets, ok := action.secretsConfig.SecretsNames["API_KEY"]
	require.True(t, ok)
	require.Len(t, secrets, 1)
	assert.Equal(t, "MISSING_ENV_VAR", secrets[0])
}

func TestDirectConfidentialHTTPAction_Redirects(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/somewhere-else", http.StatusFound)
	}))
	t.Cleanup(srv.Close)

	lggr := logger.Test(t)
	action := NewDirectConfidentialHTTPAction(lggr, filepath.Join(t.TempDir(), "secrets.yaml"))

	input := &confidentialhttp.ConfidentialHTTPRequest{
		Request: &confidentialhttp.HTTPRequest{
			Url:    srv.URL,
			Method: "GET",
		},
	}

	result, err := action.SendRequest(context.Background(), commonCap.RequestMetadata{}, input)
	require.Error(t, err)
	assert.Equal(t, caperrors.InvalidArgument, err.Code())
	assert.Contains(t, err.Error(), "redirects are not allowed")
	assert.Nil(t, result)
}
