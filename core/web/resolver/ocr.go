package resolver

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
)

type OCRKeyBundle struct {
	id                    string
	configPublicKey       ocrkey.ConfigPublicKey
	offChainPublicKey     ocrkey.OffChainPublicKey
	onChainSigningAddress ocrkey.OnChainSigningAddress
}

func (k OCRKeyBundle) ID() string {
	return k.id
}

func (k OCRKeyBundle) ConfigPublicKey() string {
	return k.configPublicKey.String()
}

func (k OCRKeyBundle) OffChainPublicKey() string {
	return k.offChainPublicKey.String()
}

func (k OCRKeyBundle) OnChainSigningAddress() string {
	return k.onChainSigningAddress.String()
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
		bundles = append(bundles, OCRKeyBundle{
			id:                    k.ID(),
			configPublicKey:       k.PublicKeyConfig(),
			offChainPublicKey:     k.OffChainSigning.PublicKey(),
			onChainSigningAddress: k.OnChainSigning.Address(),
		})
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
	return OCRKeyBundle{
		id:                    r.key.ID(),
		configPublicKey:       r.key.PublicKeyConfig(),
		offChainPublicKey:     r.key.OffChainSigning.PublicKey(),
		onChainSigningAddress: r.key.OnChainSigning.Address(),
	}
}

type DeleteOCRKeyBundleSuccessResolver struct {
	key ocrkey.KeyV2
}

func NewDeleteOCRKeyBundleSuccessResolver(key ocrkey.KeyV2) *DeleteOCRKeyBundleSuccessResolver {
	return &DeleteOCRKeyBundleSuccessResolver{key}
}

func (r *DeleteOCRKeyBundleSuccessResolver) Bundle() OCRKeyBundle {
	return OCRKeyBundle{
		id:                    r.key.ID(),
		configPublicKey:       r.key.PublicKeyConfig(),
		offChainPublicKey:     r.key.OffChainSigning.PublicKey(),
		onChainSigningAddress: r.key.OnChainSigning.Address(),
	}
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
