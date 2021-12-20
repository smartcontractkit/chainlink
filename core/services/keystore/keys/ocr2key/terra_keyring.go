package ocr2key

import (
	"bytes"
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/minio/sha256-simd"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ ocrtypes.OnchainKeyring = &terraKeyring{}

type terraKeyring struct {
	privateKey ecdsa.PrivateKey
}

func newTerraKeyring(material io.Reader) (*terraKeyring, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, material)
	if err != nil {
		return nil, err
	}
	return &terraKeyring{privateKey: *ecdsaKey}, nil
}

// XXX: PublicKey returns the evm-style address of the public key not the public key itself
func (ok *terraKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	address := crypto.PubkeyToAddress(*(&ok.privateKey).Public().(*ecdsa.PublicKey))
	return address[:]
}

func (ok *terraKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) []byte {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h := sha256.New()
	h.Write([]byte{uint8(len(report))})
	h.Write(report)
	h.Write(rawReportContext[0][:])
	h.Write(rawReportContext[1][:])
	h.Write(rawReportContext[2][:])
	return h.Sum(nil)
}

func (ok *terraKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	return crypto.Sign(ok.reportToSigData(reportCtx, report), &ok.privateKey)

}

func (ok *terraKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	hash := ok.reportToSigData(reportCtx, report)
	authorPubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	return bytes.Equal(publicKey[:], authorAddress[:])
}

func (ok *terraKeyring) MaxSignatureLength() int {
	return 65
}

func (ok *terraKeyring) marshal() ([]byte, error) {
	return crypto.FromECDSA(&ok.privateKey), nil
}

func (ok *terraKeyring) unmarshal(in []byte) error {
	privateKey, err := crypto.ToECDSA(in)
	ok.privateKey = *privateKey
	return err
}
