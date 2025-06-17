package ocr2key

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OnchainKeyring = &solanaKeyring{}

type solanaKeyring struct {
	privateKey func() *ecdsa.PrivateKey
}

func newSolanaKeyring(material io.Reader) (*solanaKeyring, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, material)
	if err != nil {
		return nil, err
	}
	return &solanaKeyring{privateKey: func() *ecdsa.PrivateKey { return ecdsaKey }}, nil
}

// XXX: PublicKey returns the evm-style address of the public key not the public key itself
func (skr *solanaKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	address := crypto.PubkeyToAddress(skr.privateKey().PublicKey)
	return address[:]
}

func (skr *solanaKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h := sha256.New()
	h.Write([]byte{uint8(len(report))})
	h.Write(report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	h.Write(rawReportContext[2][:])
	return h.Sum(nil)
}

func (skr *solanaKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return skr.SignBlob(skr.reportToSigData(reportCtx, report))
}

func (skr *solanaKeyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	bytes, err := skr.reportToSigData3(digest, seqNr, r)
	if err != nil {
		return nil, err
	}
	return skr.SignBlob(bytes)
}

func (skr *solanaKeyring) reportToSigData3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) ([]byte, error) {
	rawReportContext := RawReportContext3(digest, seqNr)
	h := sha3.NewLegacyKeccak256()
	reportLen := uint16(len(r)) //nolint:gosec // max U16 larger than solana transaction size
	err := binary.Write(h, binary.LittleEndian, reportLen)
	h.Write(r)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	return h.Sum(nil), err
}

func (skr *solanaKeyring) SignBlob(b []byte) (sig []byte, err error) {
	return crypto.Sign(b, skr.privateKey())
}

func (skr *solanaKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash := skr.reportToSigData(reportCtx, report)
	return skr.VerifyBlob(publicKey, hash, signature)
}

func (skr *solanaKeyring) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	hash, err := skr.reportToSigData3(cd, seqNr, r)
	if err != nil {
		return false
	}
	return skr.VerifyBlob(publicKey, hash, signature)
}

func (skr *solanaKeyring) VerifyBlob(pubkey types.OnchainPublicKey, b, sig []byte) bool {
	authorPubkey, err := crypto.SigToPub(b, sig)
	if err != nil {
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	// no need for constant time compare since neither arg is sensitive
	return bytes.Equal(pubkey[:], authorAddress[:])
}

func (skr *solanaKeyring) MaxSignatureLength() int {
	return 65
}

func (skr *solanaKeyring) Marshal() ([]byte, error) {
	return crypto.FromECDSA(skr.privateKey()), nil
}

func (skr *solanaKeyring) Unmarshal(in []byte) error {
	privateKey, err := crypto.ToECDSA(in)
	skr.privateKey = func() *ecdsa.PrivateKey { return privateKey }
	return err
}
