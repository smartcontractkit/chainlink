package presenters

import (
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
)

// DKGEncryptKeyResource is just that.
type DKGEncryptKeyResource struct {
	JAID
	PublicKey string `json:"publicKey"`
}

var _ jsonapi.EntityNamer = DKGEncryptKeyResource{}

// GetName implements jsonapi.EntityNamer
func (DKGEncryptKeyResource) GetName() string {
	return "encryptedDKGEncryptKeys"
}

// NewDKGEncryptKeyResource creates a new DKGEncryptKeyResource from the given DKG sign key.
func NewDKGEncryptKeyResource(key dkgencryptkey.Key) *DKGEncryptKeyResource {
	return &DKGEncryptKeyResource{
		JAID: JAID{
			ID: key.ID(),
		},
		PublicKey: key.PublicKeyString(),
	}
}

// NewDKGEncryptKeyResources creates many DKGEncryptKeyResource objects from the given DKG sign keys.
func NewDKGEncryptKeyResources(keys []dkgencryptkey.Key) (resources []DKGEncryptKeyResource) {
	for _, key := range keys {
		resources = append(resources, *NewDKGEncryptKeyResource(key))
	}
	return
}
