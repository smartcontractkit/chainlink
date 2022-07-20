package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
)

// SolanaKeyResource represents a Solana key JSONAPI resource.
type SolanaKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (SolanaKeyResource) GetName() string {
	return "encryptedSolanaKeys"
}

func NewSolanaKeyResource(key solkey.Key) *SolanaKeyResource {
	r := &SolanaKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewSolanaKeyResources(keys []solkey.Key) []SolanaKeyResource {
	rs := []SolanaKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewSolanaKeyResource(key))
	}

	return rs
}
