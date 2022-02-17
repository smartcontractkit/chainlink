package ocr2key

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffchainKeyring(t *testing.T) {
	kr, err := newOffchainKeyring(cryptorand.Reader, cryptorand.Reader)
	require.NoError(t, err)
	pubKey := kr.OffchainPublicKey()
	assert.True(t, bytes.Equal(kr.signingKey.Public().(ed25519.PublicKey), pubKey[:]))
}
