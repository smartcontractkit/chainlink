package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

// TerraKeyResource represents a Terra key JSONAPI resource.
type TerraKeyResource struct {
	JAID
	PubKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (TerraKeyResource) GetName() string {
	return "terraKeys"
}

func NewTerraKeyResource(key terrakey.Key) *TerraKeyResource {
	r := &TerraKeyResource{
		JAID:   JAID{ID: key.ID()},
		PubKey: key.PublicKeyStr(),
	}

	return r
}

func NewTerraKeyResources(keys []terrakey.Key) []TerraKeyResource {
	rs := []TerraKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewTerraKeyResource(key))
	}

	return rs
}
