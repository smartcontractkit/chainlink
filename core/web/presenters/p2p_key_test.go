package presenters

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	cryptop2p "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services/keystore/p2pkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestP2PKeyResource(t *testing.T) {
	timestamp := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	_, pubKey, err := cryptop2p.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	pubKeyBytes, err := pubKey.Raw()
	require.NoError(t, err)

	peerIDStr := "12D3KooWApUJaQB2saFjyEUfq6BmysnsSnhLnY5CF9tURYVKgoXK"
	p2pPeerID, err := peer.Decode(peerIDStr)
	require.NoError(t, err)
	peerID := p2pkey.PeerID(p2pPeerID)

	key := p2pkey.EncryptedP2PKey{
		ID:        1,
		PeerID:    peerID,
		PubKey:    pubKeyBytes,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	r := NewP2PKeyResource(key)
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"encryptedP2PKeys",
			"id":"1",
			"attributes":{
				"peerId":"p2p_%s",
				"publicKey": "%s",
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"deletedAt":null
			}
		}
	}`, peerIDStr, hex.EncodeToString(pubKeyBytes))

	assert.JSONEq(t, expected, string(b))

	// With a deleted field
	key.DeletedAt = gorm.DeletedAt(sql.NullTime{Time: timestamp, Valid: true})

	r = NewP2PKeyResource(key)
	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = fmt.Sprintf(`
	{
		"data": {
			"type":"encryptedP2PKeys",
			"id":"1",
			"attributes":{
				"peerId":"p2p_%s",
				"publicKey": "%s",
				"createdAt":"2000-01-01T00:00:00Z",
				"updatedAt":"2000-01-01T00:00:00Z",
				"deletedAt":"2000-01-01T00:00:00Z"
			}
		}
	}`, peerIDStr, hex.EncodeToString(pubKeyBytes))

	assert.JSONEq(t, expected, string(b))
}
