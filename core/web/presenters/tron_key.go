package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
)

// TronKeyResource represents a Tron key JSONAPI resource.
type TronKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (TronKeyResource) GetName() string {
	return "encryptedTronKeys"
}

func NewTronKeyResource(key tronkey.Key) *TronKeyResource {
	r := &TronKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewTronKeyResources(keys []tronkey.Key) []TronKeyResource {
	rs := []TronKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewTronKeyResource(key))
	}

	return rs
}
