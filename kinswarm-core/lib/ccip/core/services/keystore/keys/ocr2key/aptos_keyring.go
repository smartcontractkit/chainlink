package ocr2key

import (
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

var _ ocrtypes.OnchainKeyring = &aptosKeyring{}

type aptosKeyring struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func newAptosKeyring(material io.Reader) (*aptosKeyring, error) {
	pubKey, privKey, err := ed25519.GenerateKey(material)
	if err != nil {
		return nil, err
	}
	return &aptosKeyring{pubKey: pubKey, privKey: privKey}, nil
}

func (akr *aptosKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(akr.pubKey)
}

func (akr *aptosKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
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

func (akr *aptosKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData, err := akr.reportToSigData(reportCtx, report)
	if err != nil {
		return nil, err
	}
	return akr.signBlob(sigData)
}

func (akr *aptosKeyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	return nil, errors.New("not implemented")
}

func (akr *aptosKeyring) signBlob(b []byte) ([]byte, error) {
	signedMsg := ed25519.Sign(akr.privKey, b)
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return utils.ConcatBytes(akr.PublicKey(), signedMsg), nil
}

func (akr *aptosKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash, err := akr.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	return akr.verifyBlob(publicKey, hash, signature)
}

func (akr *aptosKeyring) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	return false
}

func (akr *aptosKeyring) verifyBlob(pubkey ocrtypes.OnchainPublicKey, b, sig []byte) bool {
	// Ed25519 signatures are always 64 bytes and the
	// public key (always prefixed, see Sign above) is always,
	// 32 bytes, so we always require the max signature length.
	if len(sig) != akr.MaxSignatureLength() {
		return false
	}
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(pubkey), b, sig[32:])
}

func (akr *aptosKeyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (akr *aptosKeyring) Marshal() ([]byte, error) {
	return akr.privKey.Seed(), nil
}

func (akr *aptosKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	privKey := ed25519.NewKeyFromSeed(in)
	akr.privKey = privKey
	pubKey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ed25519.PublicKey")
	}
	akr.pubKey = pubKey
	return nil
}
