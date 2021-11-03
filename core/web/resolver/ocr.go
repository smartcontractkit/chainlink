package resolver

import "github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"

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
