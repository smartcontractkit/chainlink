package ocrkey

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOCR_OffchainPublicKey_MarshalJSON(t *testing.T) {
	t.Parallel()
	rawBytes := make([]byte, 32)
	rawBytes[31] = 1
	pubKey := OffChainPublicKey(rawBytes)

	pubKeyString := "0000000000000000000000000000000000000000000000000000000000000001"
	pubKeyJSON := fmt.Sprintf(`"%s"`, pubKeyString)

	result, err := json.Marshal(pubKey)
	assert.NoError(t, err)
	assert.Equal(t, pubKeyJSON, string(result))
}

func TestOCR_OffchainPublicKey_UnmarshalJSON_Happy(t *testing.T) {
	t.Parallel()

	pubKeyString := "918a65a518c005d6367309bec4b26805f8afabef72cbf9940d9a0fd04ec80b38"
	pubKeyJSON := fmt.Sprintf(`"%s"`, pubKeyString)
	pubKey := OffChainPublicKey{}

	err := json.Unmarshal([]byte(pubKeyJSON), &pubKey)
	assert.NoError(t, err)
	assert.Equal(t, pubKeyString, pubKey.String())
}

func TestOCR_OffchainPublicKey_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()

	pubKeyString := "hello world"
	pubKeyJSON := fmt.Sprintf(`"%s"`, pubKeyString)
	pubKey := OffChainPublicKey{}

	err := json.Unmarshal([]byte(pubKeyJSON), &pubKey)
	assert.Error(t, err)
}
