package ethkey

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
)

func TestEthKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, func(keyJSON []byte, password string) (kt keys.KeyType, err error) {
		t.SkipNow()
		return kt, err
	})
}

func createKey() (keys.KeyType, error) {
	return NewV2()
}
