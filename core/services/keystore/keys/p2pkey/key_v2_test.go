package p2pkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
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

func TestP2PKeys_KeyV2(t *testing.T) {
	_, pk, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	k := Key{PrivKey: pk}
	kv2 := k.ToV2()

	pkv2 := kv2.PrivKey.Public().(ed25519.PublicKey)

	assert.Equal(t, kv2.String(), kv2.GoString())
	assert.Equal(t, ragep2ptypes.PeerID(k.PeerID()).String(), kv2.ID())
	assert.Equal(t, hex.EncodeToString(pkv2), kv2.PublicKeyHex())
}
