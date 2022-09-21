package vrfkey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
)

func TestVRFKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	return NewV2()
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	return FromEncryptedJSON(keyJSON, password)
}
