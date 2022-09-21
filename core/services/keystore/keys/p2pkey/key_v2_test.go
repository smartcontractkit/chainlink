package p2pkey

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeys_Raw(t *testing.T) {
	pk, _, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	pkr, err := pk.Raw()
	require.NoError(t, err)

	r := Raw(pkr)

	assert.Equal(t, r.String(), r.GoString())
	assert.Equal(t, "<P2P Raw Private Key>", r.String())
}

func TestP2PKeys_KeyV2(t *testing.T) {
	pk, _, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	k := Key{PrivKey: pk}
	kv2 := k.ToV2()
	pkv2, err := kv2.GetPublic().Raw()
	require.NoError(t, err)

	assert.Equal(t, kv2.String(), kv2.GoString())
	assert.Equal(t, peer.ID(k.PeerID()).String(), kv2.ID())
	assert.Equal(t, hex.EncodeToString(pkv2), kv2.PublicKeyHex())
}
