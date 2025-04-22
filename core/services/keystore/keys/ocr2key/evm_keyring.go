package ocr2key

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocrtypes.OnchainKeyring = &evmKeyring{}

type evmKeyring struct {
	privateKey func() *ecdsa.PrivateKey
}

func newEVMKeyring(material io.Reader) (*evmKeyring, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, material)
	if err != nil {
		return nil, err
	}
	return &evmKeyring{privateKey: func() *ecdsa.PrivateKey { return ecdsaKey }}, nil
}

// XXX: PublicKey returns the address of the public key not the public key itself
func (ekr *evmKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	address := ekr.signingAddress()
	return address[:]
}

func (ekr *evmKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return ekr.SignBlob(ekr.reportToSigData(reportCtx, report))
}

func (ekr *evmKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	sigData := crypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	return crypto.Keccak256(sigData)
}

func (ekr *evmKeyring) Sign3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) (signature []byte, err error) {
	return ekr.SignBlob(ReportToSigData3(digest, seqNr, r))
}

func ReportToSigData3(digest types.ConfigDigest, seqNr uint64, r ocrtypes.Report) []byte {
	rawReportContext := RawReportContext3(digest, seqNr)
	sigData := crypto.Keccak256(r)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	return crypto.Keccak256(sigData)
}

func RawReportContext3(digest types.ConfigDigest, seqNr uint64) [2][32]byte {
	seqNrBytes := [32]byte{}
	binary.BigEndian.PutUint64(seqNrBytes[24:], seqNr)
	return [2][32]byte{
		digest,
		seqNrBytes,
	}
}

func (ekr *evmKeyring) SignBlob(b []byte) (sig []byte, err error) {
	return crypto.Sign(b, ekr.privateKey())
}

func (ekr *evmKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash := ekr.reportToSigData(reportCtx, report)
	return ekr.VerifyBlob(publicKey, hash, signature)
}

func (ekr *evmKeyring) Verify3(publicKey ocrtypes.OnchainPublicKey, cd ocrtypes.ConfigDigest, seqNr uint64, r ocrtypes.Report, signature []byte) bool {
	hash := ReportToSigData3(cd, seqNr, r)
	return ekr.VerifyBlob(publicKey, hash, signature)
}

func (ekr *evmKeyring) VerifyBlob(pubkey types.OnchainPublicKey, b, sig []byte) bool {
	authorPubkey, err := crypto.SigToPub(b, sig)
	if err != nil {
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	// no need for constant time compare since neither arg is sensitive
	return bytes.Equal(pubkey[:], authorAddress[:])
}

func (ekr *evmKeyring) MaxSignatureLength() int {
	return 65
}

func (ekr *evmKeyring) signingAddress() common.Address {
	return crypto.PubkeyToAddress(ekr.privateKey().PublicKey)
}

func (ekr *evmKeyring) Marshal() ([]byte, error) {
	return crypto.FromECDSA(ekr.privateKey()), nil
}

func (ekr *evmKeyring) Unmarshal(in []byte) error {
	privateKey, err := crypto.ToECDSA(in)
	if err != nil {
		return err
	}
	ekr.privateKey = func() *ecdsa.PrivateKey { return privateKey }
	return nil
}
