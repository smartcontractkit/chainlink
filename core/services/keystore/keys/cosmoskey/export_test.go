package cosmoskey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
)

func TestCosmosKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	return New(), nil
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	return FromEncryptedJSON(keyJSON, password)
}
