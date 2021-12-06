package ocr2key

import (
	"bytes"
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/minio/sha256-simd"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ ocrtypes.OnchainKeyring = &solanaKeyring{}

type solanaKeyring struct {
	privateKey ecdsa.PrivateKey
}

func newSolanaKeyring(material io.Reader) (*solanaKeyring, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, material)
	if err != nil {
		return nil, err
	}
	return &solanaKeyring{privateKey: *ecdsaKey}, nil
}

// XXX: PublicKey returns the address of the public key not the public key itself
func (ok *solanaKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	address := ok.SigningAddress()
	return address[:]
}

func (ok *solanaKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	sigData := make([]byte, len(report))
	sigData = append(sigData, report...)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	result := sha256.Sum256(sigData)
	return result[:]
}

func (ok *solanaKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return crypto.Sign(ok.reportToSigData(reportCtx, report), &ok.privateKey)

}

func (ok *solanaKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash := ok.reportToSigData(reportCtx, report)
	authorPubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	return bytes.Equal(publicKey[:], authorAddress[:])
}

func (ok *solanaKeyring) MaxSignatureLength() int {
	return 65
}

func (ok *solanaKeyring) SigningAddress() common.Address {
	return crypto.PubkeyToAddress(*(&ok.privateKey).Public().(*ecdsa.PublicKey))
}

func (ok *solanaKeyring) marshal() ([]byte, error) {
	return crypto.FromECDSA(&ok.privateKey), nil
}

func (ok *solanaKeyring) unmarshal(in []byte) error {
	privateKey, err := crypto.ToECDSA(in)
	ok.privateKey = *privateKey
	return err
}
