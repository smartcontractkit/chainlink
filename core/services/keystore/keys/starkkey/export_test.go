package starkkey

import (
	"testing"

	starknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestStarkNetKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	key, err := starknet.New()
	return TestWrapped{key}, err
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	key, err := FromEncryptedJSON(keyJSON, password)
	return TestWrapped{key}, err
}

// wrap key to conform to desired test interface
type TestWrapped struct {
	starknet.Key
}

func (w TestWrapped) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) ([]byte, error) {
	return ToEncryptedJSON(w.Key, password, scryptParams)
}
