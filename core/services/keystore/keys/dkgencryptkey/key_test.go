package dkgencryptkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestNew(t *testing.T) {
	key, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, key.privateKey)
	assert.NotNil(t, key.PublicKey)
	assert.NotNil(t, key.publicKeyBytes)
}

func TestStringers(t *testing.T) {
	key := MustNewXXXTestingOnly(big.NewInt(1337))
	assert.Equal(t, "26578c46722826d18dc5f5a954c65c5c78e0d215a465356502ff8f002aff36ef", key.PublicKeyString())
	assert.Equal(t,
		"DKGEncryptKey{PrivateKey: <redacted>, PublicKey: 26578c46722826d18dc5f5a954c65c5c78e0d215a465356502ff8f002aff36ef",
		key.String())
	assert.Equal(t,
		"DKGEncryptKey{PrivateKey: <redacted>, PublicKey: 26578c46722826d18dc5f5a954c65c5c78e0d215a465356502ff8f002aff36ef",
		key.GoString())
	assert.Equal(t,
		"26578c46722826d18dc5f5a954c65c5c78e0d215a465356502ff8f002aff36ef",
		key.ID())
}

func TestRaw(t *testing.T) {
	key := MustNewXXXTestingOnly(big.NewInt(1337))
	rawFromKey := key.Raw()
	scalar := g1.Scalar().SetBytes(rawFromKey)
	assert.True(t, scalar.Equal(key.privateKey))

	keyFromRaw := rawFromKey.Key()
	assert.True(t, keyFromRaw.privateKey.Equal(key.privateKey))

	assert.Equal(t, "<DKGEncrypt Raw Private Key>", rawFromKey.GoString())
	assert.Equal(t, "<DKGEncrypt Raw Private Key>", rawFromKey.String())
}

func TestExportImport(t *testing.T) {
	password := "helloworld"
	key := MustNewXXXTestingOnly(big.NewInt(1337))
	encryptedJSON, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
	assert.NoError(t, err)

	decryptedKey, err := FromEncryptedJSON(encryptedJSON, password)
	assert.NoError(t, err)
	assert.True(t, decryptedKey.privateKey.Equal(key.privateKey))
	assert.True(t, decryptedKey.PublicKey.Equal(key.PublicKey))
	assert.ElementsMatch(t, decryptedKey.publicKeyBytes, key.publicKeyBytes)
}
