package ocrcommon

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	ocr2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var _ ocr3types.OnchainKeyring[[]byte] = (*OCR3OnchainKeyringAdapter)(nil)

type OCR3OnchainKeyringAdapter struct {
	o ocrtypes.OnchainKeyring
}

func NewOCR3OnchainKeyringAdapter(o ocrtypes.OnchainKeyring) *OCR3OnchainKeyringAdapter {
	return &OCR3OnchainKeyringAdapter{o}
}

func (k *OCR3OnchainKeyringAdapter) PublicKey() ocrtypes.OnchainPublicKey {
	return k.o.PublicKey()
}

func (k *OCR3OnchainKeyringAdapter) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	return k.o.Sign(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (k *OCR3OnchainKeyringAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	return k.o.Verify(opk, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, ri.Report, signature)
}

func (k *OCR3OnchainKeyringAdapter) MaxSignatureLength() int {
	return k.o.MaxSignatureLength()
}

var _ ocr3types.ContractTransmitter[[]byte] = (*OCR3ContractTransmitterAdapter)(nil)

type OCR3ContractTransmitterAdapter struct {
	ct ocrtypes.ContractTransmitter
}

func NewOCR3ContractTransmitterAdapter(ct ocrtypes.ContractTransmitter) *OCR3ContractTransmitterAdapter {
	return &OCR3ContractTransmitterAdapter{ct}
}

func (c *OCR3ContractTransmitterAdapter) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte], signatures []ocrtypes.AttributedOnchainSignature) error {
	return c.ct.Transmit(ctx, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report, signatures)
}

func (c *OCR3ContractTransmitterAdapter) FromAccount() (ocrtypes.Account, error) {
	return c.ct.FromAccount()
}

var _ ocr3types.OnchainKeyring[[]byte] = (*OCR3OnchainKeyringMultiChainAdapter)(nil)

type OCR3OnchainKeyringMultiChainAdapter struct {
	ks         keystore.OCR2
	st         ocr2.OCR2OnchainSigningStrategy
	keyBundles map[relay.Network]ocr2key.KeyBundle
	publicKey  ocrtypes.OnchainPublicKey
}

func NewOCR3OnchainKeyringMultiChainAdapter(ks keystore.OCR2, st ocr2.OCR2OnchainSigningStrategy) (*OCR3OnchainKeyringMultiChainAdapter, error) {
	var keyBundles map[relay.Network]ocr2key.KeyBundle
	for chainFamily, _ := range relay.SupportedRelays {
		kbID, err := st.KeyBundleID(chainFamily)
		if err != nil {
			return nil, err
		}
		kb, err := ks.Get(kbID)
		if err != nil {
			return nil, err
		}
		keyBundles[chainFamily] = kb
	}
	pkKeyBundleID, _ := st.PublicKey()
	kb, err := ks.Get(pkKeyBundleID)
	if err != nil {
		return nil, err
	}
	return &OCR3OnchainKeyringMultiChainAdapter{ks, st, keyBundles, kb.PublicKey()}, nil
}

func (a *OCR3OnchainKeyringMultiChainAdapter) PublicKey() ocrtypes.OnchainPublicKey {
	return a.publicKey
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	kb := a.keyBundles[relay.EVM] // TODO: how do we get the bundle name from the report info?
	return kb.Sign(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	kb := a.keyBundles[relay.EVM] // TODO: how do we get the bundle name from the report info?
	return kb.Verify(opk, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, ri.Report, signature)
}

func (a *OCR3OnchainKeyringMultiChainAdapter) MaxSignatureLength() int {
	kb := a.keyBundles[relay.EVM] // TODO: how do we get the bundle name at this point? No reporting info is available
	return kb.MaxSignatureLength()
}
