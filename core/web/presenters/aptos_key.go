package presenters

import "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"

// AptosKeyResource represents a Aptos key JSONAPI resource.
type AptosKeyResource struct {
	JAID
	Account string `json:"account"`
	PubKey  string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (AptosKeyResource) GetName() string {
	return "encryptedAptosKeys"
}

func NewAptosKeyResource(key aptoskey.Key) *AptosKeyResource {
	r := &AptosKeyResource{
		JAID:    JAID{ID: key.ID()},
		Account: key.Account(),
		PubKey:  key.PublicKeyStr(),
	}

	return r
}

func NewAptosKeyResources(keys []aptoskey.Key) []AptosKeyResource {
	rs := []AptosKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewAptosKeyResource(key))
	}

	return rs
}
