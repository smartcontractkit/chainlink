package presenters

import "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"

// TONKeyResource represents a TON key JSONAPI resource.
type TONKeyResource struct {
	JAID
	AddressBase64 string `json:"addressBase64"`
	RawAddress    string `json:"rawAddress"`
	PubKey        string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (TONKeyResource) GetName() string {
	return "encryptedTONKeys"
}

func NewTONKeyResource(key tonkey.Key) *TONKeyResource {
	r := &TONKeyResource{
		JAID:          JAID{ID: key.ID()},
		AddressBase64: key.AddressBase64(),
		RawAddress:    key.RawAddress(),
		PubKey:        key.PublicKeyStr(),
	}

	return r
}

func NewTONKeyResources(keys []tonkey.Key) []TONKeyResource {
	rs := []TONKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewTONKeyResource(key))
	}

	return rs
}
