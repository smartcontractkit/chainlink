package vrfkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var phrase = "as3r8phu82u9ru843cdi4298yf"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	_, err := k.Encrypt(phrase, FastScryptParams)
	require.NoError(t, err)
}

func TestPublicKeyRoundTrip(t *testing.T) {
	pk := k.PublicKey
	serialized, err := pk.Value()
	require.NoError(t, err)
	var npk PublicKey
	require.NoError(t, npk.Scan(serialized))
	assert.Equal(t, pk, npk)
}
