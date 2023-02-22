package ocr2key

import (
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/crypto/ed25519"
)

type legacyEdDSAKeyring struct {
	privateKey ed25519.PrivateKey
}

// PublicKey returns the ed25519.PublicKey
func (l *legacyEdDSAKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return []byte(l.privateKey.Public().(ed25519.PublicKey))
}

// Sign always returns an error, cannot use a legacy keyring to sign,
// legacy keyring only offer limited support
//
// Deprecated: legacy keyring cannot be used to sign
func (l *legacyEdDSAKeyring) Sign(_ ocrtypes.ReportContext, _ ocrtypes.Report) (signature []byte, err error) {
	return nil, errors.New("cannot use a legacy key to sign")
}

// Verify always returns false, cannot use a legacy keyring to verify a signature,
// legacy keyring only offer limited support
//
// Deprecated: legacy keyring cannot be used to verify
func (l *legacyEdDSAKeyring) Verify(_ ocrtypes.OnchainPublicKey, _ ocrtypes.ReportContext, _ ocrtypes.Report, _ []byte) bool {
	return false
}

func (l *legacyEdDSAKeyring) MaxSignatureLength() int {
	return ed25519.PublicKeySize + ed25519.SignatureSize
}

// Marshal will return the ed25519 private key seed
func (l *legacyEdDSAKeyring) Marshal() ([]byte, error) {
	return l.privateKey.Seed(), nil
}

// Unmarshal adds
func (l *legacyEdDSAKeyring) Unmarshal(in []byte) error {
	if len(in) != ed25519.SeedSize {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), ed25519.SeedSize)
	}
	l.privateKey = ed25519.NewKeyFromSeed(in)

	return nil
}
