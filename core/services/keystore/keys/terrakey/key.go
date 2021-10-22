package terrakey

import (
	"github.com/tendermint/tendermint/crypto/secp256k1"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/terra-project/terra.go/key"

	"errors"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
)

type Key struct {
	ID                  uint
	PublicKey           cryptotypes.PubKey
	privateKey          []byte
	EncryptedPrivateKey crypto.EncryptedPrivateKey
	Address             TerraAddress
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// New creates a new Terra key consisting of an secp256k1 key. It encrypts the
// Key with the passphrase.
func New(passphrase string, scryptParams utils.ScryptParams) (*Key, error) {

	privKey, _ := key.PrivKeyGen(secp256k1.GenPrivKey())
	encPrivkey, err := crypto.NewEncryptedPrivateKey(privKey.Bytes(), passphrase, scryptParams)
	if err != nil {
		return nil, err
	}

	return &Key{
		PublicKey:           privKey.PubKey(),
		privateKey:          privKey.Bytes(),
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
	return KeyV2{
		// TODO: Add more? where is this used?
		Address: k.Address,
	}
}
