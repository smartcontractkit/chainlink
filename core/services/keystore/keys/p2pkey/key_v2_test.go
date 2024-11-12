package p2pkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeys_Raw(t *testing.T) {
	_, pk, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	r := Raw(pk)

	assert.Equal(t, r.String(), r.GoString())
	assert.Equal(t, "<P2P Raw Private Key>", r.String())
}
