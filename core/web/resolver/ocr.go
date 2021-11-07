package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

type OCRKeyBundle struct {
	key ocrkey.KeyV2
}

func (k OCRKeyBundle) ID() graphql.ID {
	return graphql.ID(k.key.ID())
}

func (k OCRKeyBundle) ConfigPublicKey() string {
	return ocrkey.ConfigPublicKey(k.key.PublicKeyConfig()).String()
}

func (k OCRKeyBundle) OffChainPublicKey() string {
	return k.key.OffChainSigning.PublicKey().String()
}

func (k OCRKeyBundle) OnChainSigningAddress() string {
	return k.key.OnChainSigning.Address().String()
}

type OCRKeyBundlesPayloadResolver struct {
	keys []ocrkey.KeyV2
}

func NewOCRKeyBundlesPayloadResolver(keys []ocrkey.KeyV2) *OCRKeyBundlesPayloadResolver {
	return &OCRKeyBundlesPayloadResolver{keys}
}

func (r *OCRKeyBundlesPayloadResolver) Results() []OCRKeyBundle {
	bundles := []OCRKeyBundle{}
	for _, k := range r.keys {
		bundles = append(bundles, OCRKeyBundle{k})
	}
	return bundles
}

type CreateOCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
}

func NewCreateOCRKeyBundlePayloadResolver(key ocrkey.KeyV2) *CreateOCRKeyBundlePayloadResolver {
	return &CreateOCRKeyBundlePayloadResolver{key}
}

func (r *CreateOCRKeyBundlePayloadResolver) Bundle() OCRKeyBundle {
	return OCRKeyBundle{r.key}
}

type DeleteOCRKeyBundleSuccessResolver struct {
	key ocrkey.KeyV2
}

func NewDeleteOCRKeyBundleSuccessResolver(key ocrkey.KeyV2) *DeleteOCRKeyBundleSuccessResolver {
	return &DeleteOCRKeyBundleSuccessResolver{key}
}

func (r *DeleteOCRKeyBundleSuccessResolver) Bundle() OCRKeyBundle {
	return OCRKeyBundle{r.key}
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
		return &DeleteOCRKeyBundleSuccessResolver{r.key}, true
	}
	return nil, false
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError(r.err.Error()), true
	}
	return nil, false
}
