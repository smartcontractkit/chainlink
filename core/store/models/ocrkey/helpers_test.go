package ocrkey

func (kb *KeyBundle) ExportedOnChainSigning() *onChainPrivateKey {
	return kb.onChainSigning
}

func (kb *KeyBundle) ExportedOffChainSigning() *offChainPrivateKey {
	return kb.offChainSigning
}

func (kb *KeyBundle) ExportedOffChainEncryption() *[32]byte {
	return kb.offChainEncryption
}

type KeyBundleRawData = keyBundleRawData
