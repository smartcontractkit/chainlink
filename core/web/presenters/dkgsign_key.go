package presenters

import (
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"
)

// DKGSignKeyResource is just that.
type DKGSignKeyResource struct {
	JAID
	PublicKey string `json:"publicKey"`
}

var _ jsonapi.EntityNamer = DKGSignKeyResource{}

// GetName implements jsonapi.EntityNamer
func (DKGSignKeyResource) GetName() string {
	return "encryptedDKGSignKeys"
}

// NewDKGSignKeyResource creates a new DKGSignKeyResource from the given DKG sign key.
func NewDKGSignKeyResource(key dkgsignkey.Key) *DKGSignKeyResource {
	return &DKGSignKeyResource{
		JAID: JAID{
			ID: key.ID(),
		},
		PublicKey: key.PublicKeyString(),
	}
}

// NewDKGSignKeyResources creates many DKGSignKeyResource objects from the given DKG sign keys.
func NewDKGSignKeyResources(keys []dkgsignkey.Key) (resources []DKGSignKeyResource) {
	for _, key := range keys {
		resources = append(resources, *NewDKGSignKeyResource(key))
	}
	return
}
