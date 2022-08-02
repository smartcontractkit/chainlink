package presenters

import (
	starknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
)

// StarkNetKeyResource represents a StarkNet key JSONAPI resource.
type StarkNetKeyResource struct {
	JAID
	AccountAddr string `json:"accountAddr"`
	StarkKey    string `json:"starkPubKey"`
}

// GetName implements the api2go EntityNamer interface
func (StarkNetKeyResource) GetName() string {
	return "encryptedStarkNetKeys"
}

func NewStarkNetKeyResource(key starknet.Key) *StarkNetKeyResource {
	r := &StarkNetKeyResource{
		JAID:        JAID{ID: key.ID()},
		AccountAddr: key.AccountAddressStr(),
		StarkKey:    key.StarkKeyStr(),
	}

	return r
}

func NewStarkNetKeyResources(keys []starknet.Key) []StarkNetKeyResource {
	rs := []StarkNetKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStarkNetKeyResource(key))
	}

	return rs
}
