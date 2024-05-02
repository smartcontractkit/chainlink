package ocrcommon

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
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
	ks keystore.OCR2
	st ocr2.OCR2OnchainSigningStrategy
	// keep a map of chain-family -> key-bundle
}

func NewOCR3OnchainKeyringMultiChainAdapter(ks keystore.OCR2, st ocr2.OCR2OnchainSigningStrategy) *OCR3OnchainKeyringMultiChainAdapter {
	for a, _ := range relay.SupportedRelays {
		// go through all the key-bundles and create a map of chain-family -> key-bundle, return an error when failing to get a key-bundle
	}
	return &OCR3OnchainKeyringMultiChainAdapter{ks, st}
}

func (a *OCR3OnchainKeyringMultiChainAdapter) PublicKey() ocrtypes.OnchainPublicKey {
	// TODO: how do we handle errors if we cannot bubble them up? Do we use a logger?
	pkKeyBundleID, _ := a.st.PublicKey()
	kb, _ := a.ks.Get(pkKeyBundleID)
	return kb.PublicKey()
}

func (a *OCR3OnchainKeyringMultiChainAdapter) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	kbID, _ := a.st.KeyBundleID("") // TODO: how do we get the bundle name from the report info?
	kb, err := a.ks.Get(kbID)
	if err != nil {
		return nil, err
	}
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
	kbID, _ := a.st.KeyBundleID("") // TODO: how do we get the bundle name from the report info?
	kb, err := a.ks.Get(kbID)
	if err != nil {
		return false
	}
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
	kbID, _ := a.st.KeyBundleID("") // TODO: how do we get the bundle name in this case?
	kb, err := a.ks.Get(kbID)
	if err != nil {
		return -1
	}
	return kb.MaxSignatureLength()
}
