package ocrimpls

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

// OCR3SignerVerifierExtra is an extension of OCR3SignerVerifier that
// also exposes the public key and max signature length which are required by the ocr3Keyring adapter.
type OCR3SignerVerifierExtra interface {
	ocr2key.OCR3SignerVerifier
	PublicKey() types.OnchainPublicKey
	MaxSignatureLength() int
}

var _ ocr3types.OnchainKeyring[[]byte] = &ocr3Keyring{}

// ocr3Keyring is an adapter that exposes ocr3 onchain keyring.
type ocr3Keyring struct {
	core OCR3SignerVerifierExtra
	lggr logger.Logger
}

func NewOnchainKeyring(keyring OCR3SignerVerifierExtra, lggr logger.Logger) *ocr3Keyring {
	return &ocr3Keyring{
		core: keyring,
		lggr: lggr.Named("OCR3Keyring"),
	}
}

func (w *ocr3Keyring) PublicKey() types.OnchainPublicKey {
	return w.core.PublicKey()
}

func (w *ocr3Keyring) MaxSignatureLength() int {
	return w.core.MaxSignatureLength()
}

func (w *ocr3Keyring) Sign(configDigest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte]) (signature []byte, err error) {
	w.lggr.Debugw(
		"signing report",
		"configDigest", configDigest.Hex(),
		"seqNr", seqNr,
		"report", hexutil.Encode(r.Report),
	)
	return w.core.Sign3(configDigest, seqNr, r.Report)
}

func (w *ocr3Keyring) Verify(key types.OnchainPublicKey, configDigest types.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte], signature []byte) bool {
	w.lggr.Debugw("verifying report",
		"configDigest", configDigest.Hex(),
		"seqNr", seqNr,
		"report", hexutil.Encode(r.Report),
	)
	return w.core.Verify3(key, configDigest, seqNr, r.Report, signature)
}
