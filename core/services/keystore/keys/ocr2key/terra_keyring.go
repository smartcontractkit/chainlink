package ocr2key

import (
	"crypto/ed25519"
	"encoding/binary"
	"io"

	"github.com/hdevalence/ed25519consensus"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/blake2s"
)

var _ ocrtypes.OnchainKeyring = &terraKeyring{}

type terraKeyring struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func newTerraKeyring(material io.Reader) (*terraKeyring, error) {
	pubKey, privKey, err := ed25519.GenerateKey(material)
	if err != nil {
		return nil, err
	}
	return &terraKeyring{pubKey: pubKey, privKey: privKey}, nil
}

func (ok *terraKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(ok.pubKey)
}

func (ok *terraKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
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

func (ok *terraKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData, err := ok.reportToSigData(reportCtx, report)
	if err != nil {
		return nil, err
	}
	signedMsg := ed25519.Sign(ok.privKey, sigData)
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return utils.ConcatBytes(ok.PublicKey(), signedMsg), nil
}

// Note signature is prefixed with the public key.
func (ok *terraKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	if len(signature) != ok.MaxSignatureLength() {
		return false
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	hash, err := ok.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(publicKey), hash, signature[32:])
}

func (ok *terraKeyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (ok *terraKeyring) marshal() ([]byte, error) {
	return ok.privKey.Seed(), nil
}

func (ok *terraKeyring) unmarshal(in []byte) error {
	privKey := ed25519.NewKeyFromSeed(in)
	ok.privKey = privKey
	ok.pubKey = ed25519.PublicKey(privKey[32:])
	return nil
}
