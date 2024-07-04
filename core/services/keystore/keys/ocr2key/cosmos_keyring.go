package ocr2key

import (
	"crypto/ed25519"
	"encoding/binary"
	"io"

	"github.com/hdevalence/ed25519consensus"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2s"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OnchainKeyring = &cosmosKeyring{}

type cosmosKeyring struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func newCosmosKeyring(material io.Reader) (*cosmosKeyring, error) {
	pubKey, privKey, err := ed25519.GenerateKey(material)
	if err != nil {
		return nil, err
	}
	return &cosmosKeyring{pubKey: pubKey, privKey: privKey}, nil
}

func (ckr *cosmosKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(ckr.pubKey)
}

func (ckr *cosmosKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h, err := blake2s.New256(nil)
	if err != nil {
		return nil, err
	}
	reportLen := make([]byte, 4)
	binary.BigEndian.PutUint32(reportLen[0:], uint32(len(report)))
	h.Write(reportLen[:])
	h.Write(report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	h.Write(rawReportContext[2][:])
	return h.Sum(nil), nil
}

func (ckr *cosmosKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData, err := ckr.reportToSigData(reportCtx, report)
	if err != nil {
		return nil, err
	}
	return ckr.signBlob(sigData)
}

func (ckr *cosmosKeyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	return nil, errors.New("not implemented")
}

func (ckr *cosmosKeyring) signBlob(b []byte) ([]byte, error) {
	signedMsg := ed25519.Sign(ckr.privKey, b)
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return utils.ConcatBytes(ckr.PublicKey(), signedMsg), nil
}

func (ckr *cosmosKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash, err := ckr.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	return ckr.verifyBlob(publicKey, hash, signature)
}

func (ckr *cosmosKeyring) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	return false
}

func (ckr *cosmosKeyring) verifyBlob(pubkey ocrtypes.OnchainPublicKey, b, sig []byte) bool {
	// Ed25519 signatures are always 64 bytes and the
	// public key (always prefixed, see Sign above) is always,
	// 32 bytes, so we always require the max signature length.
	if len(sig) != ckr.MaxSignatureLength() {
		return false
	}
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(pubkey), b, sig[32:])
}

func (ckr *cosmosKeyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (ckr *cosmosKeyring) Marshal() ([]byte, error) {
	return ckr.privKey.Seed(), nil
}

func (ckr *cosmosKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	privKey := ed25519.NewKeyFromSeed(in)
	ckr.privKey = privKey
	pubKey, ok := privKey.Public().(ed25519.PublicKey)
	if !ok {
		return errors.New("failed to cast public key to ed25519.PublicKey")
	}
	ckr.pubKey = pubKey
	return nil
}
