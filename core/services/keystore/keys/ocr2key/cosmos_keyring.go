package ocr2key

import (
	"crypto/ed25519"
	"encoding/binary"
	"io"

	"github.com/hdevalence/ed25519consensus"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/blake2s"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

func (tk *cosmosKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(tk.pubKey)
}

func (tk *cosmosKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
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

func (tk *cosmosKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	sigData, err := tk.reportToSigData(reportCtx, report)
	if err != nil {
		return nil, err
	}
	signedMsg := ed25519.Sign(tk.privKey, sigData)
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return utils.ConcatBytes(tk.PublicKey(), signedMsg), nil
}

func (tk *cosmosKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	// Ed25519 signatures are always 64 bytes and the
	// public key (always prefixed, see Sign above) is always,
	// 32 bytes, so we always require the max signature length.
	if len(signature) != tk.MaxSignatureLength() {
		return false
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return false
	}
	hash, err := tk.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	return ed25519consensus.Verify(ed25519.PublicKey(publicKey), hash, signature[32:])
}

func (tk *cosmosKeyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (tk *cosmosKeyring) Marshal() ([]byte, error) {
	return tk.privKey.Seed(), nil
}

func (tk *cosmosKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	privKey := ed25519.NewKeyFromSeed(in)
	tk.privKey = privKey
	tk.pubKey = privKey.Public().(ed25519.PublicKey)
	return nil
}
