package ocr2key

import (
	"bytes"
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ ocrtypes.OnchainKeyring = &EthereumKeyring{}

type EthereumKeyring struct {
	privateKey ecdsa.PrivateKey
}

func newEthereumKeyring(material io.Reader) (*EthereumKeyring, error) {
	ecdsaKey, err := ecdsa.GenerateKey(curve, material)
	if err != nil {
		return nil, err
	}
	return &EthereumKeyring{privateKey: *ecdsaKey}, nil
}

// XXX: PublicKey returns the address of the public key not the public key itself
func (ok *EthereumKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	publicKey := (*ecdsa.PrivateKey)(&ok.privateKey).Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	return address[:]
}

func (ok *EthereumKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	sigData := crypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	return crypto.Sign(crypto.Keccak256(sigData), &ok.privateKey)

}

func (ok *EthereumKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	sigData := crypto.Keccak256(report)
	sigData = append(sigData, rawReportContext[0][:]...)
	sigData = append(sigData, rawReportContext[1][:]...)
	sigData = append(sigData, rawReportContext[2][:]...)
	hash := crypto.Keccak256(sigData)
	authorPubkey, err := crypto.SigToPub(hash, signature)

	if err != nil {
		return false
	}
	authorAddress := crypto.PubkeyToAddress(*authorPubkey)
	return bytes.Equal(publicKey[:], authorAddress[:])
}

func (ok *EthereumKeyring) MaxSignatureLength() int {
	return 65
}

func (ok *EthereumKeyring) SigningAddress() common.Address {
	return crypto.PubkeyToAddress(*(*ecdsa.PrivateKey)(&ok.privateKey).Public().(*ecdsa.PublicKey))
}

func (ok *EthereumKeyring) marshal() ([]byte, error) {
	return crypto.FromECDSA(&ok.privateKey), nil
}

func (ok *EthereumKeyring) unmarshal(in []byte) error {
	privateKey, err := crypto.ToECDSA(in)
	ok.privateKey = *privateKey
	return err
}
