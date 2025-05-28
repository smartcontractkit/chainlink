package ocr2key

import (
	"crypto/ed25519"
	"crypto/sha256"
	"io"

	"github.com/hdevalence/ed25519consensus"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OnchainKeyring = &tonKeyring{}

type tonKeyring struct {
	privKey func() ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func newTONKeyring(material io.Reader) (*tonKeyring, error) {
	pubKey, privKey, err := ed25519.GenerateKey(material)
	if err != nil {
		return nil, err
	}
	return &tonKeyring{pubKey: pubKey, privKey: func() ed25519.PrivateKey { return privKey }}, nil
}

func (tkr *tonKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(tkr.pubKey)
}

func (tkr *tonKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h := sha256.New()
	h.Write([]byte{uint8(len(report))}) //nolint:gosec // assumes len(report) < 256
	h.Write(report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	h.Write(rawReportContext[2][:])
	return h.Sum(nil)
}

func (tkr *tonKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData := tkr.reportToSigData(reportCtx, report)
	return tkr.SignBlob(sigData)
}

func (tkr *tonKeyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) ([]byte, error) {
	bytes := tkr.reportToSigData3(digest, seqNr, r)
	return tkr.SignBlob(bytes)
}

func (tkr *tonKeyring) reportToSigData3(digest types.ConfigDigest, seqNr uint64, report ocrtypes.Report) []byte {
	rawReportContext := RawReportContext3(digest, seqNr)
	h := sha256.New()
	h.Write([]byte{uint8(len(report))}) //nolint:gosec // assumes len(report) < 256
	h.Write(report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	return h.Sum(nil)
}

func (tkr *tonKeyring) SignBlob(b []byte) ([]byte, error) {
	sig := ed25519.Sign(tkr.privKey(), b)
	return utils.ConcatBytes(tkr.PublicKey(), sig), nil
}

func (tkr *tonKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash := tkr.reportToSigData(reportCtx, report)
	return tkr.VerifyBlob(publicKey, hash, signature)
}

func (tkr *tonKeyring) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	hash := tkr.reportToSigData3(cd, seqNr, r)
	return tkr.VerifyBlob(publicKey, hash, signature)
}

func (tkr *tonKeyring) VerifyBlob(pubkey ocrtypes.OnchainPublicKey, b, sig []byte) bool {
	if len(sig) != tkr.MaxSignatureLength() {
		return false
	}
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(pubkey), b, sig[32:])
}

func (tkr *tonKeyring) MaxSignatureLength() int {
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (tkr *tonKeyring) Marshal() ([]byte, error) {
	return tkr.privKey().Seed(), nil
}

func (tkr *tonKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	privKey := ed25519.NewKeyFromSeed(in)
	tkr.privKey = func() ed25519.PrivateKey { return privKey }
	pubKey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ed25519.PublicKey")
	}
	tkr.pubKey = pubKey
	return nil
}
