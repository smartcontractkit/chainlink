package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

type OCRKeyBundleResolver struct {
	key ocrkey.KeyV2
}

func NewOCRKeyBundleResolver(key ocrkey.KeyV2) OCRKeyBundleResolver {
	return OCRKeyBundleResolver{key}
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
	return &OCRKeyBundlesPayloadResolver{keys}
}

func (r *OCRKeyBundlesPayloadResolver) Results() []OCRKeyBundleResolver {
	var bundles []OCRKeyBundleResolver
	for _, k := range r.keys {
		bundles = append(bundles, NewOCRKeyBundleResolver(k))
	}
	return bundles
}

type CreateOCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
}

func NewCreateOCRKeyBundlePayloadResolver(key ocrkey.KeyV2) *CreateOCRKeyBundlePayloadResolver {
	return &CreateOCRKeyBundlePayloadResolver{key}
}

func (r *CreateOCRKeyBundlePayloadResolver) Bundle() OCRKeyBundleResolver {
	return OCRKeyBundleResolver{r.key}
}

type DeleteOCRKeyBundleSuccessResolver struct {
	key ocrkey.KeyV2
}

func NewDeleteOCRKeyBundleSuccessResolver(key ocrkey.KeyV2) *DeleteOCRKeyBundleSuccessResolver {
	return &DeleteOCRKeyBundleSuccessResolver{key}
}

func (r *DeleteOCRKeyBundleSuccessResolver) Bundle() OCRKeyBundleResolver {
	return OCRKeyBundleResolver{r.key}
}

type DeleteOCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
	err error
}

func NewDeleteOCRKeyBundlePayloadResolver(key ocrkey.KeyV2, err error) *DeleteOCRKeyBundlePayloadResolver {
	return &DeleteOCRKeyBundlePayloadResolver{key, err}
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToDeleteOCRKeyBundleSuccess() (*DeleteOCRKeyBundleSuccessResolver, bool) {
	if r.err == nil {
		return NewDeleteOCRKeyBundleSuccessResolver(r.key), true
	}
	return nil, false
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError(r.err.Error()), true
	}
	return nil, false
}
