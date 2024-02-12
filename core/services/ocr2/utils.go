package ocr2

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.OnchainKeyring[any] = (*keyBundleOCR3Wrapper)(nil)

type keyBundleOCR3Wrapper struct {
	kb ocrtypes.OnchainKeyring
}

func (k *keyBundleOCR3Wrapper) PublicKey() ocrtypes.OnchainPublicKey {
	return k.kb.PublicKey()
}

func (k *keyBundleOCR3Wrapper) Sign(digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[any]) (signature []byte, err error) {
	return k.kb.Sign(ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report)
}

func (k *keyBundleOCR3Wrapper) Verify(opk ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, ri ocr3types.ReportWithInfo[any], signature []byte) bool {
	return k.kb.Verify(opk, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, ri.Report, signature)
}

func (k *keyBundleOCR3Wrapper) MaxSignatureLength() int {
	return k.kb.MaxSignatureLength()
}

var _ ocr3types.ContractTransmitter[any] = (*contractTransmitterOCR3Wrapper)(nil)

type contractTransmitterOCR3Wrapper struct {
	ct ocrtypes.ContractTransmitter
}

func (c contractTransmitterOCR3Wrapper) Transmit(ctx context.Context, digest ocrtypes.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[any], signatures []ocrtypes.AttributedOnchainSignature) error {
	return c.ct.Transmit(ctx, ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        uint32(seqNr),
			Round:        0,
		},
		ExtraHash: [32]byte(make([]byte, 32)),
	}, r.Report, signatures)
}

func (c contractTransmitterOCR3Wrapper) FromAccount() (ocrtypes.Account, error) {
	return c.ct.FromAccount()
}
