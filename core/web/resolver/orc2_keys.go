package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
)

// OCR2ChainType defines OCR2 Chain Types accepted on this resolver
type OCR2ChainType string

const (
	// OCR2ChainTypeEVM defines OCR2 EVM Chain Type
	OCR2ChainTypeEVM OCR2ChainType = "EMV"
	// OCR2ChainTypeSolana defines OCR2 Solana Chain Type
	OCR2ChainTypeSolana OCR2ChainType = "SOLANA"
	// OCR2ChainTypeTerra defines OCR2 Terra Chain Type
	OCR2ChainTypeTerra OCR2ChainType = "TERRA"
)

// ToOCR2ChainType turns a valid string into a OCR2ChainType
func ToOCR2ChainType(s string) (OCR2ChainType, error) {
	switch s {
	case "evm":
		return OCR2ChainTypeEVM, nil
	case "solana":
		return OCR2ChainTypeSolana, nil
	case "terra":
		return OCR2ChainTypeTerra, nil
	default:
		return "", errors.New("invalid ocr2 chain type")
	}
}

// FromOCR2ChainType returns the string (lowercased) value from a OCR2ChainType
func FromOCR2ChainType(ct OCR2ChainType) string {
	switch ct {
	case OCR2ChainTypeEVM:
		return "evm"
	case OCR2ChainTypeSolana:
		return "solana"
	case OCR2ChainTypeTerra:
		return "terra"
	default:
		return strings.ToLower(string(ct))
	}
}

// OCR2KeyBundleResolver defines the OCR2 Key bundle on GQL
type OCR2KeyBundleResolver struct {
	key ocr2key.KeyBundle
}

// NewOCR2KeyBundle creates a new GQL OCR2 key bundle resolver
func NewOCR2KeyBundle(key ocr2key.KeyBundle) *OCR2KeyBundleResolver {
	return &OCR2KeyBundleResolver{key: key}
}

// ChainType returns the OCR2 Key bundle chain type
func (r OCR2KeyBundleResolver) ChainType() *OCR2ChainType {
	ct, err := ToOCR2ChainType(string(r.key.ChainType()))
	if err != nil {
		return nil
	}

	return &ct
}

// OnChainPublicKey returns the OCR2 Key bundle on-chain public key
func (r OCR2KeyBundleResolver) OnChainPublicKey() string {
	return fmt.Sprintf("ocr2on_%s_%s", r.key.ChainType(), r.key.OnChainPublicKey())
}

// OffChainPublicKey returns the OCR2 Key bundle off-chain public key
func (r OCR2KeyBundleResolver) OffChainPublicKey() string {
	return fmt.Sprintf("ocr2off_%s_%s", r.key.ChainType(), hex.EncodeToString(r.key.OffchainPublicKey()))
}

// ConfigPublicKey returns the OCR2 Key bundle config public key
func (r OCR2KeyBundleResolver) ConfigPublicKey() string {
	configPublic := r.key.ConfigEncryptionPublicKey()
	return fmt.Sprintf("ocr2cfg_%s_%s", r.key.ChainType(), hex.EncodeToString(configPublic[:]))
}

// -- OCR2Keys Query --

// OCR2KeyBundlesPayloadResolver defines the OCR2 Key bundles query resolver
type OCR2KeyBundlesPayloadResolver struct {
	keys []ocr2key.KeyBundle
}

// NewOCR2KeyBundlesPayload returns the OCR2 key bundles resolver
func NewOCR2KeyBundlesPayload(keys []ocr2key.KeyBundle) *OCR2KeyBundlesPayloadResolver {
	return &OCR2KeyBundlesPayloadResolver{keys: keys}
}

// Results resolves the list of OCR2 key bundles
func (r *OCR2KeyBundlesPayloadResolver) Results() []OCR2KeyBundleResolver {
	var results []OCR2KeyBundleResolver

	for _, k := range r.keys {
		results = append(results, *NewOCR2KeyBundle(k))
	}

	return results
}
