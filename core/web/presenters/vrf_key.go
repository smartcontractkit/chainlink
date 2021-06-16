package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
)

type VRFKeyResource struct {
	JAID
	Compressed   string     `json:"compressed"`
	Uncompressed string     `json:"uncompressed"`
	Hash         string     `json:"hash"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

// GetName implements the api2go EntityNamer interface
func (VRFKeyResource) GetName() string {
	return "encryptedVRFKeys"
}

func NewVRFKeyResource(key vrfkey.EncryptedVRFKey) *VRFKeyResource {
	uncompressed, err := key.PublicKey.StringUncompressed()
	if err != nil {
		logger.Error("unable to get uncompressed pk", "err", err)
	}
	r := &VRFKeyResource{
		JAID:         NewJAID(key.PublicKey.String()),
		Compressed:   key.PublicKey.String(),
		Uncompressed: uncompressed,
		Hash:         key.PublicKey.MustHash().String(),
		CreatedAt:    key.CreatedAt,
		UpdatedAt:    key.UpdatedAt,
	}

	if key.DeletedAt.Valid {
		r.DeletedAt = &key.DeletedAt.Time
	}

	return r
}

func NewVRFKeyResources(keys []*vrfkey.EncryptedVRFKey) []VRFKeyResource {
	rs := []VRFKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewVRFKeyResource(*key))
	}

	return rs
}
