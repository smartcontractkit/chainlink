package presenters

import "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"

// TonKeyResource represents a Ton key JSONAPI resource.
type TonKeyResource struct {
	JAID
	UserFriendlyAddress string `json:"userFriendlyAddress"`
	RawAddress          string `json:"rawAddress"`
	PubKey              string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (TonKeyResource) GetName() string {
	return "encryptedTonKeys"
}

func NewTonKeyResource(key tonkey.Key) *TonKeyResource {
	r := &TonKeyResource{
		JAID:                JAID{ID: key.ID()},
		UserFriendlyAddress: key.UserFriendlyAddress(),
		RawAddress:          key.RawAddress(),
		PubKey:              key.PublicKeyStr(),
	}

	return r
}

func NewTonKeyResources(keys []tonkey.Key) []TonKeyResource {
	rs := []TonKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewTonKeyResource(key))
	}

	return rs
}
