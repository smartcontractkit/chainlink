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

type OCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
}

func NewOCRKeyBundleResolver(key ocrkey.KeyV2) *OCRKeyBundlePayloadResolver {
	return &OCRKeyBundlePayloadResolver{key}
}

func (r *OCRKeyBundlePayloadResolver) Bundle() OCRKeyBundle {
	return OCRKeyBundle{
		id:                    r.key.ID(),
		configPublicKey:       r.key.PublicKeyConfig(),
		offChainPublicKey:     r.key.OffChainSigning.PublicKey(),
		onChainSigningAddress: r.key.OnChainSigning.Address(),
	}
}

// OCRErrorResolver resolves errors that could arise from interacting with the OCR
// resource. Possible errors may include key not found errors or unprocessable errors
// (i.e, the keystore is locked), but there could be many more failure modes.
type OCRErrorResolver struct {
	message string
	code    ErrorCode
}

func NewOCRErrorResolver(message string, code ErrorCode) *OCRErrorResolver {
	return &OCRErrorResolver{message, code}
}

func (r *OCRErrorResolver) Message() string { return r.message }
func (r *OCRErrorResolver) Code() ErrorCode { return r.code }

type DeleteOCRKeyBundlePayloadResolver struct {
	key ocrkey.KeyV2
	err error
}

func NewDeleteOCRKeyBundlePayloadResolver(key ocrkey.KeyV2, err error) *DeleteOCRKeyBundlePayloadResolver {
	return &DeleteOCRKeyBundlePayloadResolver{key, err}
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToOCRKeyBundlePayload() (*OCRKeyBundlePayloadResolver, bool) {
	if r.err == nil {
		return &OCRKeyBundlePayloadResolver{r.key}, true
	}
	return nil, false
}

func (r *DeleteOCRKeyBundlePayloadResolver) ToOCRError() (*OCRErrorResolver, bool) {
	if r.err != nil {
		return NewOCRErrorResolver(r.err.Error(), ErrorCodeUnprocessable), true
	}
	return nil, false
}
