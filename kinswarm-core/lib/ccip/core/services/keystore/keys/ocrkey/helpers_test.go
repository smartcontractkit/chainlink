package ocrkey

func (kb *KeyV2) ExportedOnChainSigning() *onChainPrivateKey {
	return kb.OnChainSigning
}

func (kb *KeyV2) ExportedOffChainSigning() *offChainPrivateKey {
	return kb.OffChainSigning
}

func (kb *KeyV2) ExportedOffChainEncryption() *[32]byte {
	return kb.OffChainEncryption
}
