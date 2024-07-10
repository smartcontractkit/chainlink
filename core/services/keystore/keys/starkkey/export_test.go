package starkkey

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaintext(t *testing.T) {
	// replace with the contents of encrypted-key.txt
	keyJSON := []byte(`<ENCRYPTED_STRING>`)
	// replace with your encryption password
	key, _ := FromEncryptedJSON(keyJSON, "<PASSWORD>")
	privateKey := fmt.Sprintf("0x0%x", key.ToPrivKey())
	assert.Equal(t, privateKey, "^ your decrypted private key")
}

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

func (w TestWrapped) ToEncryptedJSON(password string, scryptParams utils.ScryptParams) ([]byte, error) {
	return ToEncryptedJSON(w.Key, password, scryptParams)
}
