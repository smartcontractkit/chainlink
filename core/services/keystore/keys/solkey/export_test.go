package solkey

import (
	"testing"

	solkeys "github.com/smartcontractkit/chainlink-solana/pkg/solana/keys"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestSolanaKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	key, err := solkeys.New()
	return TestWrapped{key}, err
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	key, err := FromEncryptedJSON(keyJSON, password)
	return TestWrapped{key}, err
}

type TestWrapped struct {
	solkeys.Key
}

func (w TestWrapped) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) ([]byte, error) {
	return ToEncryptedJSON(w.Key, password, scryptParams)
}
