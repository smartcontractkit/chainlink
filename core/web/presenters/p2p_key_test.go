package presenters

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func TestP2PKeyResource(t *testing.T) {
	key := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(1))
	peerID := key.PeerID()
	peerIDStr := peerID.String()

	r := NewP2PKeyResource(key)
	b, err := jsonapi.Marshal(r)
	require.NoError(t, err)

	expected := fmt.Sprintf(`
	{
		"data":{
			"type":"encryptedP2PKeys",
			"id":"%s",
			"attributes":{
				"peerId":"%s",
				"publicKey": "%s"
			}
		}
	}`, key.ID(), peerIDStr, key.PublicKeyHex())

	assert.JSONEq(t, expected, string(b))

	r = NewP2PKeyResource(key)
	b, err = jsonapi.Marshal(r)
	require.NoError(t, err)

	expected = fmt.Sprintf(`
	{
		"data": {
			"type":"encryptedP2PKeys",
			"id":"%s",
			"attributes":{
				"peerId":"%s",
				"publicKey": "%s"
			}
		}
	}`, key.ID(), peerIDStr, key.PublicKeyHex())

	assert.JSONEq(t, expected, string(b))
}
