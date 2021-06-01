package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/p2pkey"
)

// P2PKeyResource represents a P2P key JSONAPI resource.
type P2PKeyResource struct {
	JAID
	PeerID    string     `json:"peerId"`
	PubKey    string     `json:"publicKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

// GetName implements the api2go EntityNamer interface
func (P2PKeyResource) GetName() string {
	return "encryptedP2PKeys"
}

func NewP2PKeyResource(key p2pkey.EncryptedP2PKey) *P2PKeyResource {
	r := &P2PKeyResource{
		JAID:      NewJAIDInt32(key.ID),
		PeerID:    key.PeerID.String(),
		PubKey:    key.PubKey.String(),
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}

	if key.DeletedAt.Valid {
		r.DeletedAt = &key.DeletedAt.Time
	}

	return r
}

func NewP2PKeyResources(keys []p2pkey.EncryptedP2PKey) []P2PKeyResource {
	rs := []P2PKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewP2PKeyResource(key))
	}

	return rs
}
