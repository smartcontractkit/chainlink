package ocr2key

import (
	cryptorand "crypto/rand"
	"io"

	starksig "github.com/NethermindEth/juno/pkg/crypto/signature"
	"github.com/NethermindEth/juno/pkg/crypto/weierstrass"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ ocrtypes.OnchainKeyring = &starknetKeyring{}

var starkCurve = weierstrass.Stark()

type starknetKeyring struct {
	privateKey starksig.PrivateKey
}

func newStarknetKeyring(material io.Reader) (*starknetKeyring, error) {
	privKey, err := starksig.GenerateKey(starkCurve, material)
	if err != nil {
		return nil, err
	}
	return &starknetKeyring{privateKey: *privKey}, err
}

func (sk *starknetKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return weierstrass.Marshal(starkCurve, sk.privateKey.PublicKey.X, sk.privateKey.PublicKey.Y)
}

func (sk *starknetKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {

	digest := []byte{}
	return starksig.SignASN1(cryptorand.Reader, &sk.privateKey, digest)
}

func (sk *starknetKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {

	pubKey := starksig.PublicKey{Curve: starkCurve}
	pubKey.X, pubKey.Y = weierstrass.Unmarshal(starkCurve, publicKey)

  digest := []byte{}
	return starksig.VerifyASN1(&pubKey, digest, signature)
}

func (sk *starknetKeyring) MaxSignatureLength() int {
	return 0
}
