package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

// P2PKeyResource represents a P2P key JSONAPI resource.
type P2PKeyResource struct {
	JAID
	PeerID string `json:"peerId"`
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (P2PKeyResource) GetName() string {
	return "encryptedP2PKeys"
}

func NewP2PKeyResource(key p2pkey.KeyV2) *P2PKeyResource {
	r := &P2PKeyResource{
		JAID:   JAID{ID: key.ID()},
		PeerID: key.PeerID().String(),
		PubKey: key.PublicKeyHex(),
	}

	return r
}

func NewP2PKeyResources(keys []p2pkey.KeyV2) []P2PKeyResource {
	rs := []P2PKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewP2PKeyResource(key))
	}

	return rs
}
