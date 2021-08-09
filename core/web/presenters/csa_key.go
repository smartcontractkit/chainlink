package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
)

// CSAKeyResource represents a CSA key JSONAPI resource.
type CSAKeyResource struct {
	JAID
	PubKey    string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (CSAKeyResource) GetName() string {
	return "csaKeys"
}

func NewCSAKeyResource(key csakey.Key) *CSAKeyResource {
	r := &CSAKeyResource{
		JAID:      NewJAIDUint(key.ID),
		PubKey:    key.PublicKey.String(),
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}

	return r
}

func NewCSAKeyResources(keys []csakey.Key) []CSAKeyResource {
	rs := []CSAKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewCSAKeyResource(key))
	}

	return rs
}
