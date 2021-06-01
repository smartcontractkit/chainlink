package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/ocrkey"
)

// OCRKeysBundleResource represents a bundle of OCRs keys as JSONAPI resource
type OCRKeysBundleResource struct {
	JAID
	OnChainSigningAddress ocrkey.OnChainSigningAddress `json:"onChainSigningAddress"`
	OffChainPublicKey     ocrkey.OffChainPublicKey     `json:"offChainPublicKey"`
	ConfigPublicKey       ocrkey.ConfigPublicKey       `json:"configPublicKey"`
	CreatedAt             time.Time                    `json:"createdAt"`
	UpdatedAt             time.Time                    `json:"updatedAt"`
	DeletedAt             *time.Time                   `json:"deletedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r OCRKeysBundleResource) GetName() string {
	return "encryptedKeyBundles"
}

func NewOCRKeysBundleResource(bundle ocrkey.EncryptedKeyBundle) *OCRKeysBundleResource {
	r := &OCRKeysBundleResource{
		JAID:                  NewJAID(bundle.ID.String()),
		OnChainSigningAddress: bundle.OnChainSigningAddress,
		OffChainPublicKey:     bundle.OffChainPublicKey,
		ConfigPublicKey:       bundle.ConfigPublicKey,
		CreatedAt:             bundle.CreatedAt,
		UpdatedAt:             bundle.UpdatedAt,
	}

	if bundle.DeletedAt.Valid {
		r.DeletedAt = &bundle.DeletedAt.Time
	}

	return r
}

func NewOCRKeysBundleResources(keys []ocrkey.EncryptedKeyBundle) []OCRKeysBundleResource {
	rs := []OCRKeysBundleResource{}
	for _, key := range keys {
		rs = append(rs, *NewOCRKeysBundleResource(key))
	}

	return rs
}
