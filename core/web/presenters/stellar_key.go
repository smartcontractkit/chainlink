package presenters

import (
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
)

// StellarKeyResource represents a Stellar key JSONAPI resource.
type StellarKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (StellarKeyResource) GetName() string {
	return "encryptedStellarKeys"
}

func NewStellarKeyResource(key stellarkey.Key) *StellarKeyResource {
	r := &StellarKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewStellarKeyResources(keys []stellarkey.Key) []StellarKeyResource {
	rs := []StellarKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStellarKeyResource(key))
	}

	return rs
}
