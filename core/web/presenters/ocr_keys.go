package presenters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
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

// OCRKeysBundleResource represents a bundle of OCRs keys as JSONAPI resource
type OCR2KeysBundleResource struct {
	JAID
	OnChainSigningAddress common.Address                     `json:"onChainSigningAddress"`
	OffChainPublicKey     ocrtypes.OffchainPublicKey         `json:"offChainPublicKey"`
	ConfigPublicKey       ocrtypes.ConfigEncryptionPublicKey `json:"configPublicKey"`
}

// GetName implements the api2go EntityNamer interface
func (r OCR2KeysBundleResource) GetName() string {
	return "keyV2s"
}

func NewOCR2KeysBundleResource(key ocr2key.KeyBundle) *OCR2KeysBundleResource {
	return &OCR2KeysBundleResource{
		JAID:                  NewJAID(key.ID()),
		OnChainSigningAddress: key.OnchainKeyring.SigningAddress(),
		OffChainPublicKey:     key.OffchainKeyring.OffchainPublicKey(),
		ConfigPublicKey:       key.PublicKeyConfig(),
	}
}

func NewOCR2KeysBundleResources(keys []ocr2key.KeyBundle) []OCR2KeysBundleResource {
	rs := []OCR2KeysBundleResource{}
	for _, key := range keys {
		rs = append(rs, *NewOCR2KeysBundleResource(key))
	}

	return rs
}
