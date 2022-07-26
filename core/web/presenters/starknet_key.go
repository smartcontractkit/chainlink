package presenters

import (
	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
)

// StarkNetKeyResource represents a Solana key JSONAPI resource.
type StarkNetKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (StarkNetKeyResource) GetName() string {
	return "encryptedStarkNetKeys"
}

func NewStarkNetKeyResource(key starkkey.StarkKey) *StarkNetKeyResource {
	r := &StarkNetKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewStarkNetKeyResources(keys []starkkey.StarkKey) []StarkNetKeyResource {
	rs := []StarkNetKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStarkNetKeyResource(key))
	}

	return rs
}
