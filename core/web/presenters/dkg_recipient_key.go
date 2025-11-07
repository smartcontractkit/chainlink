package presenters

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
)

type DKGRecipientKeyResource struct {
	JAID
	PublicKey string `json:"publicKey"`
}

// GetName implements the api2go EntityNamer interface
func (DKGRecipientKeyResource) GetName() string {
	return "dkgRecipientKeys"
}

func NewDKGRecipientKeyResource(key dkgrecipientkey.Key) *DKGRecipientKeyResource {
	return &DKGRecipientKeyResource{
		JAID:      NewJAID(key.PublicKeyString()),
		PublicKey: key.PublicKeyString(),
	}
}

func NewDKGRecipientKeyResources(keys []dkgrecipientkey.Key) []DKGRecipientKeyResource {
	rs := []DKGRecipientKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewDKGRecipientKeyResource(key))
	}

	return rs
}
