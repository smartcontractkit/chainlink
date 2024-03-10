package csakey

import (
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

type Key struct {
	ID                  uint
	PublicKey           crypto.PublicKey
	privateKey          []byte
	EncryptedPrivateKey crypto.EncryptedPrivateKey
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// New creates a new CSA key consisting of an ed25519 key. It encrypts the
// Key with the passphrase.
func New(passphrase string, scryptParams utils.ScryptParams) (*Key, error) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	encPrivkey, err := crypto.NewEncryptedPrivateKey(privkey, passphrase, scryptParams)
	if err != nil {
		return nil, err
	}

	return &Key{
		PublicKey:           crypto.PublicKey(pubkey),
		privateKey:          privkey,
		EncryptedPrivateKey: *encPrivkey,
	}, nil
}

func (k *Key) Unlock(password string) error {
	pk, err := k.EncryptedPrivateKey.Decrypt(password)
	if err != nil {
		return err
	}
	k.privateKey = pk
	return nil
}

func (k *Key) Unsafe_GetPrivateKey() ([]byte, error) {
	if k.privateKey == nil {
		return nil, errors.New("key has not been unlocked")
	}

	return k.privateKey, nil
}

func (k Key) ToV2() KeyV2 {
	pk := ed25519.PrivateKey(k.privateKey)
	return KeyV2{
		privateKey: &pk,
		PublicKey:  ed25519.PublicKey(k.PublicKey),
		Version:    1,
	}
}
