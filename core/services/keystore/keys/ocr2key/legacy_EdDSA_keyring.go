package ocr2key

import (
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/ed25519"
)

type legacyEdDSAKeyring struct {
	privateKey ed25519.PrivateKey
}

func (l *legacyEdDSAKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(l.privateKey.Public().(ed25519.PublicKey))
}

func (l *legacyEdDSAKeyring) Sign(_ ocrtypes.ReportContext, _ ocrtypes.Report) (signature []byte, err error) {
	return nil, errors.New("cannot use a legacy key to sign")
}

func (l *legacyEdDSAKeyring) Verify(_ ocrtypes.OnchainPublicKey, _ ocrtypes.ReportContext, _ ocrtypes.Report, _ []byte) bool {
	return false
}

func (l *legacyEdDSAKeyring) MaxSignatureLength() int {
	return ed25519.PublicKeySize + ed25519.SignatureSize
}

func (l *legacyEdDSAKeyring) Marshal() ([]byte, error) {
	return l.privateKey.Seed(), nil
}

func (l *legacyEdDSAKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	l.privateKey = ed25519.NewKeyFromSeed(in)

	return nil
}
