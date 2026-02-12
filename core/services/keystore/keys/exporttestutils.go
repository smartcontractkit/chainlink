package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// KeyType represents a key type for keys testing
type KeyType interface {
	ToEncryptedJSON(password string, scryptParams commonkeystore.ScryptParams) (export []byte, err error)
	Raw() internal.Raw
	ID() string
}

// CreateKeyFunc represents a function to create a key
type CreateKeyFunc func() (KeyType, error)

// DecryptFunc represents a function to decrypt a key
type DecryptFunc func(keyJSON []byte, password string) (KeyType, error)

// RunKeyExportImportTestcase executes a testcase to validate keys import/export functionality
func RunKeyExportImportTestcase(t *testing.T, createKey CreateKeyFunc, decrypt DecryptFunc) {
	key, err := createKey()
	require.NoError(t, err)

	json, err := key.ToEncryptedJSON("password", commonkeystore.FastScryptParams)
	require.NoError(t, err)

	assert.NotEmpty(t, json)

	imported, err := decrypt(json, "password")
	require.NoError(t, err)

	require.Equal(t, key.ID(), imported.ID())
	require.Equal(t, internal.RawBytes(key), internal.RawBytes(imported))
}

func RequireEqualKeys(t *testing.T, a, b interface {
	ID() string
	Raw() internal.Raw
}) {
	t.Helper()
	require.Equal(t, a.ID(), b.ID(), "ids be equal")
	require.Equal(t, a.Raw(), b.Raw(), "raw bytes must be equal")
	require.EqualExportedValues(t, a, b)
}
