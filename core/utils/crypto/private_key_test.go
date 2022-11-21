package crypto

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func Test_EncryptedPrivateKey(t *testing.T) {
	t.Parallel()

	privatekey := []byte("privatekey")
	passphrase := "passphrase"
	ecp, err := NewEncryptedPrivateKey(privatekey, passphrase, utils.FastScryptParams)
	require.NoError(t, err)

	actual, err := ecp.Decrypt(passphrase)
	require.NoError(t, err)

	assert.Equal(t, privatekey, actual)
}

func Test_EncryptedPrivateKey_Decrypt(t *testing.T) {
	t.Parallel()

	passphrase := []byte("passphrase")
	_, privkey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	encprivkey, err := keystore.EncryptDataV3(privkey, passphrase, 2, 1)
	require.NoError(t, err)

	ecp := EncryptedPrivateKey{CryptoJSON: encprivkey}

	actual, err := ecp.Decrypt(string(passphrase))
	require.NoError(t, err)

	assert.Equal(t, []byte(privkey), actual)
}

func Test_EncryptedPrivateKey_Scan(t *testing.T) {
	t.Parallel()

	_, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	encPrivkey, err := keystore.EncryptDataV3(privKey, []byte("passphrase"), 2, 1)
	require.NoError(t, err)
	b, err := json.Marshal(encPrivkey)
	require.NoError(t, err)

	actual := &EncryptedPrivateKey{}

	// Error if not bytes
	err = actual.Scan("not bytes")
	assert.Error(t, err)

	// Bytes
	err = actual.Scan(b)
	require.NoError(t, err)

	// Unmarshaling bytes into a struct results in numbers being stored as a
	// float64 which prevents us from asserting against the generated public key
	// which uses ints. Instead we do a JSON string comparison
	expPrivKey, err := json.Marshal(EncryptedPrivateKey{CryptoJSON: encPrivkey})
	require.NoError(t, err)
	actPrivKey, err := json.Marshal(actual)
	require.NoError(t, err)
	assert.JSONEq(t, string(expPrivKey), string(actPrivKey))
}

func Test_EncryptedPrivateKey_Value(t *testing.T) {
	t.Parallel()

	_, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	cryptoJSON, err := keystore.EncryptDataV3(privKey, []byte("passphrase"), 2, 1)
	require.NoError(t, err)

	encPrivkey := EncryptedPrivateKey{CryptoJSON: cryptoJSON}

	dv, err := encPrivkey.Value()
	require.NoError(t, err)

	expected, err := json.Marshal(EncryptedPrivateKey{CryptoJSON: cryptoJSON})
	require.NoError(t, err)

	assert.Equal(t, expected, dv)
}
