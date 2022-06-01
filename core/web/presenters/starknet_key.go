package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
)

// StarknetKeyResource represents a Solana key JSONAPI resource.
type StarknetKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (StarknetKeyResource) GetName() string {
	return "encryptedSolanaKeys"
}

func NewStarknetKeyResource(key starkkey.Key) *StarknetKeyResource {
	r := &StarknetKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewStarknetKeyResources(keys []starkkey.Key) []StarknetKeyResource {
	rs := []StarknetKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStarknetKeyResource(key))
	}

	return rs
}
