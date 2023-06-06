package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

// CosmosKeyResource represents a Cosmos key JSONAPI resource.
type CosmosKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (CosmosKeyResource) GetName() string {
	return "cosmosKeys"
}

func NewCosmosKeyResource(key cosmoskey.Key) *CosmosKeyResource {
	r := &CosmosKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewCosmosKeyResources(keys []cosmoskey.Key) []CosmosKeyResource {
	rs := []CosmosKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewCosmosKeyResource(key))
	}

	return rs
}
