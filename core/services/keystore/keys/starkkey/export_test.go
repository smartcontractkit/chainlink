package starkkey

import (
	"testing"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
)

func TestStarkNetKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	key, err := New()
	return TestWrapped{key}, err
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	key, err := FromEncryptedJSON(keyJSON, password)
	return TestWrapped{key}, err
}

// wrap key to conform to desired test interface
type TestWrapped struct {
	Key
}

func (w TestWrapped) ToEncryptedJSON(password string, scryptParams commonkeystore.ScryptParams) ([]byte, error) {
	return ToEncryptedJSON(w.Key, password, scryptParams)
}
