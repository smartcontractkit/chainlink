package ocr2key

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/binary"
	"io"

	cosmosed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/blake2s"
)

var _ ocrtypes.OnchainKeyring = &terraKeyring{}

type terraKeyring struct {
	*cosmosed25519.PrivKey
	secret []byte
}

func newTerraKeyring() *terraKeyring {
	secret := make([]byte, 32)
	_, err := io.ReadFull(cryptorand.Reader, secret)
	if err != nil {
		panic(err)
	}
	return &terraKeyring{
		PrivKey: cosmosed25519.GenPrivKeyFromSecret(secret),
		secret:  secret,
	}
}

func (ok *terraKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return ok.PubKey().Bytes()
}

func (ok *terraKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	rawReportContext := evmutil.RawReportContext(reportCtx)
	h, err := blake2s.New256(nil)
	if err != nil {
		return nil, err
	}
	reportLen := make([]byte, 8)
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
	signedMsg, err := ok.PrivKey.Sign(sigData)
	if err != nil {
		return nil, err
	}
	// match on-chain parsing (first 32 bytes are for pubkey, remaining are for signature)
	return append(ok.PublicKey(), signedMsg...), nil
}

func (ok *terraKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	var cosmosPub cosmosed25519.PubKey
	err := cosmosPub.UnmarshalAmino(publicKey)
	if err != nil {
		return false
	}
	hash, err := ok.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}
	result := cosmosPub.VerifySignature(hash, signature[32:])
	return result
}

func (ok *terraKeyring) MaxSignatureLength() int {
	// Reference: https://pkg.go.dev/crypto/ed25519
	return ed25519.PublicKeySize + ed25519.SignatureSize // 32 + 64
}

func (ok *terraKeyring) marshal() ([]byte, error) {
	return ok.secret, nil
}

func (ok *terraKeyring) unmarshal(in []byte) error {
	key := cosmosed25519.GenPrivKeyFromSecret(in)
	ok.PrivKey = key
	return nil
}
