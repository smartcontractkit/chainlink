package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
)

// CSAKeyResource represents a CSA key JSONAPI resource.
type CSAKeyResource struct {
	JAID
	PubKey  string `json:"publicKey"`
	Version int    `json:"version"`
}

// GetName implements the api2go EntityNamer interface
func (CSAKeyResource) GetName() string {
	return "csaKeys"
}

func NewCSAKeyResource(key csakey.KeyV2) *CSAKeyResource {
	r := &CSAKeyResource{
		JAID:    NewJAID(key.ID()),
		PubKey:  key.PublicKeyString(),
		Version: 1,
	}

	return r
}

func NewCSAKeyResources(keys []csakey.KeyV2) []CSAKeyResource {
	rs := []CSAKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewCSAKeyResource(key))
	}

	return rs
}
