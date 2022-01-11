package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
)

type OCR2ChainType string

const (
	OCR2ChainTypeEVM    OCR2ChainType = "EMV"
	OCR2ChainTypeSolana OCR2ChainType = "SOLANA"
	OCR2ChainTypeTerra  OCR2ChainType = "TERRA"
)

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

type OCR2KeyBundleResolver struct {
	key ocr2key.KeyBundle
}

func NewOCR2KeyBundle(key ocr2key.KeyBundle) *OCR2KeyBundleResolver {
	return &OCR2KeyBundleResolver{key: key}
}

func (r OCR2KeyBundleResolver) ChainType() *OCR2ChainType {
	ct, err := ToOCR2ChainType(string(r.key.ChainType()))
	if err != nil {
		return nil
	}

	return &ct
}

func (r OCR2KeyBundleResolver) OnChainPublicKey() string {
	return fmt.Sprintf("ocr2on_%s_%s", r.key.ChainType(), r.key.OnChainPublicKey())
}

func (r OCR2KeyBundleResolver) OffChainPublicKey() string {
	return fmt.Sprintf("ocr2off_%s_%s", r.key.ChainType(), hex.EncodeToString(r.key.OffchainPublicKey()))
}

func (r OCR2KeyBundleResolver) ConfigPublicKey() string {
	configPublic := r.key.ConfigEncryptionPublicKey()
	return fmt.Sprintf("ocr2cfg_%s_%s", r.key.ChainType(), hex.EncodeToString(configPublic[:]))
}

// -- OCR2Keys Query --

type OCR2KeyBundlesPayloadResolver struct {
	keys []ocr2key.KeyBundle
}

func NewOCR2KeyBundlesPayload(keys []ocr2key.KeyBundle) *OCR2KeyBundlesPayloadResolver {
	return &OCR2KeyBundlesPayloadResolver{keys: keys}
}

func (r *OCR2KeyBundlesPayloadResolver) Results() []OCR2KeyBundleResolver {
	var results []OCR2KeyBundleResolver

	for _, k := range r.keys {
		results = append(results, *NewOCR2KeyBundle(k))
	}

	return results
}
