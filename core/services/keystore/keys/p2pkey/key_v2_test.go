package p2pkey

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeys_KeyV2(t *testing.T) {
	kv2, err := NewV2()
	require.NoError(t, err)

	pkv2 := kv2.Public().(ed25519.PublicKey)

	assert.Equal(t, ragep2ptypes.PeerID(kv2.PeerID()).String(), kv2.ID())
	assert.Equal(t, hex.EncodeToString(pkv2), kv2.PublicKeyHex())
}
