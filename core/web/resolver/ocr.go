package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

type OCRKeyBundleResolver struct {
	key ocrkey.KeyV2
}

func NewOCRKeyBundleResolver(key ocrkey.KeyV2) OCRKeyBundleResolver {
	return OCRKeyBundleResolver{key: key}
}

func (k OCRKeyBundleResolver) ID() graphql.ID {
	return graphql.ID(k.key.ID())
}

func (k OCRKeyBundleResolver) ConfigPublicKey() string {
	return ocrkey.ConfigPublicKey(k.key.PublicKeyConfig()).String()
}

func (k OCRKeyBundleResolver) OffChainPublicKey() string {
	return k.key.OffChainSigning.PublicKey().String()
}

func (k OCRKeyBundleResolver) OnChainSigningAddress() string {
	return k.key.OnChainSigning.Address().String()
}

type OCRKeyBundlesPayloadResolver struct {
	keys []ocrkey.KeyV2
}

func NewOCRKeyBundlesPayloadResolver(keys []ocrkey.KeyV2) *OCRKeyBundlesPayloadResolver {
	return &OCRKeyBundlesPayloadResolver{keys: keys}
}

func (r *OCRKeyBundlesPayloadResolver) Results() []OCRKeyBundleResolver {
	var bundles []OCRKeyBundleResolver
	for _, k := range r.keys {
		bundles = append(bundles, NewOCRKeyBundleResolver(k))
	}
	return bundles
}

// -- CreateOCRKeyBundle Mutation --

type CreateOCRKeyBundlePayloadResolver struct {
	key *ocrkey.KeyV2
}

func NewCreateOCRKeyBundlePayload(key *ocrkey.KeyV2) *CreateOCRKeyBundlePayloadResolver {
	return &CreateOCRKeyBundlePayloadResolver{key: key}
}

func (r *CreateOCRKeyBundlePayloadResolver) ToCreateOCRKeyBundleSuccess() (*CreateOCRKeyBundleSuccessResolver, bool) {
	if r.key != nil {
		return NewCreateOCRKeyBundleSuccess(r.key), true
	}

	return nil, false
}

type CreateOCRKeyBundleSuccessResolver struct {
	key *ocrkey.KeyV2
}

func NewCreateOCRKeyBundleSuccess(key *ocrkey.KeyV2) *CreateOCRKeyBundleSuccessResolver {
	return &CreateOCRKeyBundleSuccessResolver{key: key}
}

func (r *CreateOCRKeyBundleSuccessResolver) Bundle() *OCRKeyBundleResolver {
	return &OCRKeyBundleResolver{key: *r.key}
}

// -- Delete --

type DeleteOCRKeyBundleSuccessResolver struct {
	key ocrkey.KeyV2
}

func NewDeleteOCRKeyBundleSuccessResolver(key ocrkey.KeyV2) *DeleteOCRKeyBundleSuccessResolver {
	return &DeleteOCRKeyBundleSuccessResolver{key: key}
}

func (r *DeleteOCRKeyBundleSuccessResolver) Bundle() OCRKeyBundleResolver {
	return OCRKeyBundleResolver{key: r.key}
}

type DeleteOCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
	NotFoundErrorUnionType
}

func NewDeleteOCRKeyBundlePayloadResolver(key ocrkey.KeyV2, err error) *DeleteOCRKeyBundlePayloadResolver {
	var e NotFoundErrorUnionType

	if err != nil {
		e = NotFoundErrorUnionType{err: err, message: err.Error(), isExpectedErrorFn: func(err error) bool {
			return errors.As(err, &keystore.KeyNotFoundError{})
		}}
	}

	return &DeleteOCRKeyBundlePayloadResolver{key: key, NotFoundErrorUnionType: e}
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToDeleteOCRKeyBundleSuccess() (*DeleteOCRKeyBundleSuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteOCRKeyBundleSuccessResolver(r.key), true
	}
	return nil, false
}
