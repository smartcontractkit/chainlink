package dkgsignkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	assert.Equal(t, "becd7a86af89b2f3ffd11fabe897de820b74cd2956c6e047a14e35d090ade17d", key.PublicKeyString())
	assert.Equal(t,
		"DKGSignKey{PrivateKey: <redacted>, PublicKey: becd7a86af89b2f3ffd11fabe897de820b74cd2956c6e047a14e35d090ade17d",
		key.String())
	assert.Equal(t,
		"DKGSignKey{PrivateKey: <redacted>, PublicKey: becd7a86af89b2f3ffd11fabe897de820b74cd2956c6e047a14e35d090ade17d",
		key.GoString())
	assert.Equal(t,
		"becd7a86af89b2f3ffd11fabe897de820b74cd2956c6e047a14e35d090ade17d",
		key.ID())
}

func TestRaw(t *testing.T) {
	key := MustNewXXXTestingOnly(big.NewInt(1337))
	rawFromKey := key.Raw()
	scalar := suite.Scalar().SetBytes(rawFromKey)
	assert.True(t, scalar.Equal(key.privateKey))

	keyFromRaw := rawFromKey.Key()
	assert.True(t, keyFromRaw.privateKey.Equal(key.privateKey))

	assert.Equal(t, "<DKGSign Raw Private Key>", rawFromKey.GoString())
	assert.Equal(t, "<DKGSign Raw Private Key>", rawFromKey.String())
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
