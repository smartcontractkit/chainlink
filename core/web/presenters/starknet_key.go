package presenters

import (
	starknet "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
)

// StarknetKeyResource represents a Starknet key JSONAPI resource.
type StarknetKeyResource struct {
	JAID
	AccountAddr string `json:"accountAddr"`
	StarkKey    string `json:"starkPubKey"`
}

// GetName implements the api2go EntityNamer interface
func (StarknetKeyResource) GetName() string {
	return "encryptedStarknetKeys"
}

func NewStarknetKeyResource(key starknet.Key) *StarknetKeyResource {
	r := &StarknetKeyResource{
		JAID:        JAID{ID: key.ID()},
		AccountAddr: key.AccountAddressStr(),
		StarkKey:    key.StarkKeyStr(),
	}

	return r
}

func NewStarknetKeyResources(keys []starknet.Key) []StarknetKeyResource {
	rs := []StarknetKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewStarknetKeyResource(key))
	}

	return rs
}
