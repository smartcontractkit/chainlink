package tonkey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
)

func TestTONKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	return New()
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	return FromEncryptedJSON(keyJSON, password)
}
