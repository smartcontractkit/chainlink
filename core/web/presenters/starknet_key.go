package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
)

// StarkNetKeyResource represents a StarkNet key JSONAPI resource.
type StarkNetKeyResource struct {
	JAID
	StarkKey string `json:"starkPubKey"`
}

// GetName implements the api2go EntityNamer interface
func (StarkNetKeyResource) GetName() string {
	return "encryptedStarkNetKeys"
}

func NewStarkNetKeyResource(key starkkey.Key) *StarkNetKeyResource {
	r := &StarkNetKeyResource{
		JAID:     JAID{ID: key.ID()},
		StarkKey: key.StarkKeyStr(),
	}

	return r
}

func NewStarkNetKeyResources(keys []starkkey.Key) []StarkNetKeyResource {
	rs := []StarkNetKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStarkNetKeyResource(key))
	}

	return rs
}
