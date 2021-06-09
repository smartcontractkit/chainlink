package csa

import (
	"crypto/ed25519"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

type CSAKey struct {
	ID                  uint
	PublicKey           crypto.PublicKey
	EncryptedPrivateKey crypto.EncryptedPrivateKey
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewCSAKey creates a new CSA key consisting of an ed25519 key. It encrypts the
// CSAKey with the passphrase.
func NewCSAKey(passphrase string, scryptParams utils.ScryptParams) (*CSAKey, error) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	encPrivkey, err := crypto.NewEncryptedPrivateKey(privkey, passphrase, scryptParams)
	if err != nil {
		return nil, err
	}

	return &CSAKey{
		PublicKey:           crypto.PublicKey(pubkey),
		EncryptedPrivateKey: *encPrivkey,
	}, nil
}
