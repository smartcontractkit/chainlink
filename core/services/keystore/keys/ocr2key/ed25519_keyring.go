package ocr2key

import (
	"bytes"
	"crypto/ed25519"
	"io"

	"github.com/hdevalence/ed25519consensus"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OnchainKeyring = &ed25519Keyring{}

type ed25519Keyring struct {
	privKey func() ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func newEd25519Keyring(material io.Reader) (*ed25519Keyring, error) {
	pubKey, privKey, err := ed25519.GenerateKey(material)
	if err != nil {
		return nil, err
	}
	return &ed25519Keyring{pubKey: pubKey, privKey: func() ed25519.PrivateKey { return privKey }}, nil
}

func (akr *ed25519Keyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(akr.pubKey)
}

func (akr *ed25519Keyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}
	// blake2b_256(report_context | report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	h.Write(rawReportContext[2][:])
	h.Write(report)
	return h.Sum(nil), nil
}

func (akr *ed25519Keyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData, err := akr.reportToSigData(reportCtx, report)
	if err != nil {
		return nil, err
	}
	return akr.SignBlob(sigData)
}

func (akr *ed25519Keyring) reportToSigData3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) ([]byte, error) {
	rawReportContext := RawReportContext3(digest, seqNr)
	h, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}
	h.Write(r)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	return h.Sum(nil), nil
}

func (akr *ed25519Keyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	sigData, err := akr.reportToSigData3(digest, seqNr, r)
	if err != nil {
		return nil, err
	}
	return akr.SignBlob(sigData)
}

func (akr *ed25519Keyring) SignBlob(b []byte) ([]byte, error) {
	signedMsg := ed25519.Sign(akr.privKey(), b)
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return utils.ConcatBytes(akr.PublicKey(), signedMsg), nil
}

func (akr *ed25519Keyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash, err := akr.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	return akr.VerifyBlob(publicKey, hash, signature)
}

func (akr *ed25519Keyring) Verify3(publicKey ocrtypes.OnchainPublicKey, digest ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	sigData, err := akr.reportToSigData3(digest, seqNr, r)
	if err != nil {
		return false
	}
	return akr.VerifyBlob(publicKey, sigData, signature)
}

func (akr *ed25519Keyring) VerifyBlob(pubkey ocrtypes.OnchainPublicKey, b, sig []byte) bool {
	// Ed25519 signatures are always 64 bytes and the
	// public key (always prefixed, see Sign above) is always,
	// 32 bytes, so we always require the max signature length.
	if len(sig) != akr.MaxSignatureLength() {
		return false
	}
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	if !bytes.Equal(pubkey, sig[:ed25519.PublicKeySize]) {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(pubkey), b, sig[ed25519.PublicKeySize:])
}

func (akr *ed25519Keyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (akr *ed25519Keyring) Marshal() ([]byte, error) {
	return akr.privKey().Seed(), nil
}

func (akr *ed25519Keyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	privKey := ed25519.NewKeyFromSeed(in)
	akr.privKey = func() ed25519.PrivateKey { return privKey }
	pubKey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ed25519.PublicKey")
	}
	akr.pubKey = pubKey
	return nil
}
