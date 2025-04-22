package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// KeyType represents a key type for keys testing
type KeyType interface {
	ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error)
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

	json, err := key.ToEncryptedJSON("password", utils.FastScryptParams)
	require.NoError(t, err)

	assert.NotEmpty(t, json)

	imported, err := decrypt(json, "password")
	require.NoError(t, err)

	require.Equal(t, key.ID(), imported.ID())
	require.Equal(t, internal.RawBytes(key), internal.RawBytes(imported))
}
