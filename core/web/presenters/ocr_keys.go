package presenters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

// OCRKeysBundleResource represents a bundle of OCRs keys as JSONAPI resource
type OCRKeysBundleResource struct {
	JAID
	OnChainSigningAddress ocrkey.OnChainSigningAddress `json:"onChainSigningAddress"`
	OffChainPublicKey     ocrkey.OffChainPublicKey     `json:"offChainPublicKey"`
	ConfigPublicKey       ocrkey.ConfigPublicKey       `json:"configPublicKey"`
}

// GetName implements the api2go EntityNamer interface
func (r OCRKeysBundleResource) GetName() string {
	return "keyV2s"
}

func NewOCRKeysBundleResource(key ocrkey.KeyV2) *OCRKeysBundleResource {
	return &OCRKeysBundleResource{
		JAID:                  NewJAID(key.ID()),
		OnChainSigningAddress: key.OnChainSigning.Address(),
		OffChainPublicKey:     key.OffChainSigning.PublicKey(),
		ConfigPublicKey:       key.PublicKeyConfig(),
	}
}

func NewOCRKeysBundleResources(keys []ocrkey.KeyV2) []OCRKeysBundleResource {
	rs := []OCRKeysBundleResource{}
	for _, key := range keys {
		rs = append(rs, *NewOCRKeysBundleResource(key))
	}

	return rs
}
