package vrfkey

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
)

var phrase = "as3r8phu82u9ru843cdi4298yf"

var serialK = MustNewPrivateKey(big.NewInt(int64(sk)))

func TestEncryptDecryptRoundTrip(t *testing.T) {
	// Encrypt already does a roundtrip to make sure it can decrypt, anyway
	_, err := serialK.Encrypt(phrase, utils.FastScryptParams)
	assert.NoError(t, err,
		"failed to roundtrip secret key through enecryption/decryption")
}

func TestPublicKeyRoundTrip(t *testing.T) {
	pk := serialK.PublicKey
	serialized, err := pk.Value()
	require.NoError(t, err, "failed to serialize public key for db")
	var npk secp256k1.PublicKey
	require.NoError(t, npk.Scan(serialized),
		"could not deserialize serialized public key")
	assert.Equal(t, pk, npk, "should get same key back after Value/Scan roundtrip")
}
