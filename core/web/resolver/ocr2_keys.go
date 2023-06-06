package resolver

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

// OCR2ChainType defines OCR2 Chain Types accepted on this resolver
type OCR2ChainType string

// These constants map to the enum type OCR2ChainType in ocr2_keys.graphql
const (
	// OCR2ChainTypeEVM defines OCR2 EVM Chain Type
	OCR2ChainTypeEVM = "EVM"
	// OCR2ChainTypeCosmos defines OCR2 Cosmos Chain Type
	OCR2ChainTypeCosmos = "COSMOS"
	// OCR2ChainTypeSolana defines OCR2 Solana Chain Type
	OCR2ChainTypeSolana = "SOLANA"
	// OCR2ChainTypeStarkNet defines OCR2 StarkNet Chain Type
	OCR2ChainTypeStarkNet = "STARKNET"
)

// ToOCR2ChainType turns a valid string into a OCR2ChainType
func ToOCR2ChainType(s string) (OCR2ChainType, error) {
	switch s {
	case string(chaintype.EVM):
		return OCR2ChainTypeEVM, nil
	case string(chaintype.Cosmos):
		return OCR2ChainTypeCosmos, nil
	case string(chaintype.Solana):
		return OCR2ChainTypeSolana, nil
	case string(chaintype.StarkNet):
		return OCR2ChainTypeStarkNet, nil
	default:
		return "", errors.New("unknown ocr2 chain type")
	}
}

// FromOCR2ChainType returns the string (lowercased) value from a OCR2ChainType
func FromOCR2ChainType(ct OCR2ChainType) string {
	switch ct {
	case OCR2ChainTypeEVM:
		return string(chaintype.EVM)
	case OCR2ChainTypeCosmos:
		return string(chaintype.Cosmos)
	case OCR2ChainTypeSolana:
		return string(chaintype.Solana)
	case OCR2ChainTypeStarkNet:
		return string(chaintype.StarkNet)
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

// ID returns the OCR2 Key bundle ID
func (r OCR2KeyBundleResolver) ID() graphql.ID {
	return graphql.ID(r.key.ID())
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
	pubKey := r.key.OffchainPublicKey()
	return fmt.Sprintf("ocr2off_%s_%s", r.key.ChainType(), hex.EncodeToString(pubKey[:]))
}

// ConfigPublicKey returns the OCR2 Key bundle config public key
func (r OCR2KeyBundleResolver) ConfigPublicKey() string {
	configPublic := r.key.ConfigEncryptionPublicKey()
	return fmt.Sprintf("ocr2cfg_%s_%s", r.key.ChainType(), hex.EncodeToString(configPublic[:]))
}

// -- OCR2KeyBundles Query --

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

// -- CreateOCR2KeyBundle Mutation --

// CreateOCR2KeyBundlePayloadResolver defines the create OCR2 Key bundle mutation resolver
type CreateOCR2KeyBundlePayloadResolver struct {
	key *ocr2key.KeyBundle
}

// NewCreateOCR2KeyBundlePayload returns the create OCR2 key bundle resolver
func NewCreateOCR2KeyBundlePayload(key *ocr2key.KeyBundle) *CreateOCR2KeyBundlePayloadResolver {
	return &CreateOCR2KeyBundlePayloadResolver{key: key}
}

// ToCreateOCR2KeyBundleSuccess resolves the create OCR2 key bundle success
func (r *CreateOCR2KeyBundlePayloadResolver) ToCreateOCR2KeyBundleSuccess() (*CreateOCR2KeyBundleSuccessResolver, bool) {
	if r.key == nil {
		return nil, false
	}

	return NewCreateOCR2KeyBundleSuccess(r.key), true
}

// CreateOCR2KeyBundleSuccessResolver defines the create OCR2 key bundle success resolver
type CreateOCR2KeyBundleSuccessResolver struct {
	key *ocr2key.KeyBundle
}

// NewCreateOCR2KeyBundleSuccess returns the create OCR2 key bundle success resolver
func NewCreateOCR2KeyBundleSuccess(key *ocr2key.KeyBundle) *CreateOCR2KeyBundleSuccessResolver {
	return &CreateOCR2KeyBundleSuccessResolver{key: key}
}

// Bundle resolves the creates OCR2 key bundle
func (r *CreateOCR2KeyBundleSuccessResolver) Bundle() *OCR2KeyBundleResolver {
	return NewOCR2KeyBundle(*r.key)
}

// -- DeleteOCR2KeyBundle mutation --

// DeleteOCR2KeyBundlePayloadResolver defines the delete OCR2 Key bundle mutation resolver
type DeleteOCR2KeyBundlePayloadResolver struct {
	key *ocr2key.KeyBundle
	NotFoundErrorUnionType
}

// NewDeleteOCR2KeyBundlePayloadResolver returns the delete OCR2 key bundle payload resolver
func NewDeleteOCR2KeyBundlePayloadResolver(key *ocr2key.KeyBundle, err error) *DeleteOCR2KeyBundlePayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			// returning true since the only error triggered by the search is a not found error
			// and we don't want the default check to happen, since it is a SQL Not Found error check
			return true
		}}
	}

	return &DeleteOCR2KeyBundlePayloadResolver{key: key, NotFoundErrorUnionType: e}
}

// ToDeleteOCR2KeyBundleSuccess resolves the delete OCR2 key bundle success
func (r *DeleteOCR2KeyBundlePayloadResolver) ToDeleteOCR2KeyBundleSuccess() (*DeleteOCR2KeyBundleSuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteOCR2KeyBundleSuccessResolver(r.key), true
	}

	return nil, false
}

// DeleteOCR2KeyBundleSuccessResolver defines the delete OCR2 key bundle success resolver
type DeleteOCR2KeyBundleSuccessResolver struct {
	key *ocr2key.KeyBundle
}

// NewDeleteOCR2KeyBundleSuccessResolver returns the delete OCR2 key bundle success resolver
func NewDeleteOCR2KeyBundleSuccessResolver(key *ocr2key.KeyBundle) *DeleteOCR2KeyBundleSuccessResolver {
	return &DeleteOCR2KeyBundleSuccessResolver{key: key}
}

// Bundle resolves the creates OCR2 key bundle
func (r *DeleteOCR2KeyBundleSuccessResolver) Bundle() *OCR2KeyBundleResolver {
	return NewOCR2KeyBundle(*r.key)
}
