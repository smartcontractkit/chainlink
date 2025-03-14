package crypto

import (
	"crypto/ed25519"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestRandomFlaky_JUST_FOR_TESTING_FLAKEGUARD(t *testing.T) {
	t.Parallel()

	// Seed random number generator with current time
	seed := time.Now().UnixNano()
	t.Logf("Using seed: %d", seed)
	r := rand.New(rand.NewSource(seed))

	// Generate a random number between 0 and 1
	randomValue := r.Float64()

	t.Logf("Random value generated: %f", randomValue)
	// Import at the top of the file:
	// import "github.com/stretchr/testify/require"

	// Using require to check that randomValue is < 0.5
	require.Less(t, randomValue, 0.2, "Random value should be less than 0.2")

	t.Log("This test randomly passed")
}

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
