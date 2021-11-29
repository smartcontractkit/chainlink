package presenters

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
)

type VRFKeyResource struct {
	JAID
	Compressed   string `json:"compressed"`
	Uncompressed string `json:"uncompressed"`
	Hash         string `json:"hash"`
}

// GetName implements the api2go EntityNamer interface
func (VRFKeyResource) GetName() string {
	return "encryptedVRFKeys"
}

func NewVRFKeyResource(key vrfkey.KeyV2, lggr logger.Logger) *VRFKeyResource {
	uncompressed, err := key.PublicKey.StringUncompressed()
	if err != nil {
		lggr.Errorw("Unable to get uncompressed pk", "err", err)
	}
	return &VRFKeyResource{
		JAID:         NewJAID(key.PublicKey.String()),
		Compressed:   key.PublicKey.String(),
		Uncompressed: uncompressed,
		Hash:         key.PublicKey.MustHash().String(),
	}
}

func NewVRFKeyResources(keys []vrfkey.KeyV2, lggr logger.Logger) []VRFKeyResource {
	rs := []VRFKeyResource{}
	for _, key := range keys {
		rs = append(rs, *NewVRFKeyResource(key, lggr))
	}

	return rs
}
