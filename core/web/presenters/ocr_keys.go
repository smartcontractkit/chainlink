package presenters

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
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

// OCR2KeysBundleResource represents a bundle of OCRs keys as JSONAPI resource
type OCR2KeysBundleResource struct {
	JAID
	ChainType         string `json:"chainType"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffChainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
}

// GetName implements the api2go EntityNamer interface
func (r OCR2KeysBundleResource) GetName() string {
	return "keyV2s"
}

func NewOCR2KeysBundleResource(key ocr2key.KeyBundle) *OCR2KeysBundleResource {
	configPublic := key.ConfigEncryptionPublicKey()
	pubKey := key.OffchainPublicKey()
	return &OCR2KeysBundleResource{
		JAID:              NewJAID(key.ID()),
		ChainType:         string(key.ChainType()),
		OnchainPublicKey:  fmt.Sprintf("ocr2on_%s_%s", key.ChainType(), key.OnChainPublicKey()),
		OffChainPublicKey: fmt.Sprintf("ocr2off_%s_%s", key.ChainType(), hex.EncodeToString(pubKey[:])),
		ConfigPublicKey:   fmt.Sprintf("ocr2cfg_%s_%s", key.ChainType(), hex.EncodeToString(configPublic[:])),
	}
}

func NewOCR2KeysBundleResources(keys []ocr2key.KeyBundle) []OCR2KeysBundleResource {
	rs := []OCR2KeysBundleResource{}
	for _, key := range keys {
		rs = append(rs, *NewOCR2KeysBundleResource(key))
	}
	// sort by chain type alphabetical, tie-break with ID
	sort.SliceStable(rs, func(i, j int) bool {
		if rs[i].ChainType == rs[j].ChainType {
			return rs[i].ID < rs[j].ID
		}
		return rs[i].ChainType < rs[j].ChainType
	})

	return rs
}
