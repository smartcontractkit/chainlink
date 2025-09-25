package presenters

import "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/suikey"

// SuiKeyResource represents a Sui key JSONAPI resource.
type SuiKeyResource struct {
	JAID
	Account string `json:"account"`
	PubKey  string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (SuiKeyResource) GetName() string {
	return "encryptedSuiKeys"
}

func NewSuiKeyResource(key suikey.Key) *SuiKeyResource {
	r := &SuiKeyResource{
		JAID:    JAID{ID: key.ID()},
		Account: key.Account(),
		PubKey:  key.PublicKeyStr(),
	}

	return r
}

func NewSuiKeyResources(keys []suikey.Key) []SuiKeyResource {
	rs := []SuiKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewSuiKeyResource(key))
	}

	return rs
}
