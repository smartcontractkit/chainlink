package resolver

import "github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"

type OCRKey struct {
	id                    string
	configPublicKey       ocrkey.ConfigPublicKey
	offChainPublicKey     ocrkey.OffChainPublicKey
	onChainSigningAddress ocrkey.OnChainSigningAddress
}

func (k OCRKey) ID() string {
	return k.id
}

func (k OCRKey) ConfigPublicKey() string {
	return k.configPublicKey.String()
}

func (k OCRKey) OffChainPublicKey() string {
	return k.offChainPublicKey.String()
}

func (k OCRKey) OnChainSigningAddress() string {
	return k.onChainSigningAddress.String()
}

type OCRKeysPayloadResolver struct {
	keys []ocrkey.KeyV2
}

func NewOCRKeysPayloadResolver(keys []ocrkey.KeyV2) *OCRKeysPayloadResolver {
	return &OCRKeysPayloadResolver{keys}
}

func (r *OCRKeysPayloadResolver) Results() []OCRKey {
	viewKeys := []OCRKey{}
	for _, k := range r.keys {
		viewKeys = append(viewKeys, OCRKey{
			id:                    k.ID(),
			configPublicKey:       k.PublicKeyConfig(),
			offChainPublicKey:     k.OffChainSigning.PublicKey(),
			onChainSigningAddress: k.OnChainSigning.Address(),
		})
	}
	return viewKeys
}
