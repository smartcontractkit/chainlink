package keystore

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

var (
	ErrCSAKeyExists = errors.New("a csa key already exists")
)

//go:generate mockery --name CSAKeystoreInterface --output mocks/ --case=underscore

type CSAKeystoreInterface interface {
	CreateCSAKey() (*csakey.Key, error)
	ListCSAKeys() ([]csakey.Key, error)
	Unsafe_GetUnlockedPrivateKey(pubkey crypto.PublicKey) ([]byte, error)
}

type CSA struct {
	mu           *sync.RWMutex
	orm          csaORM
	password     string
	keys         map[string]*csakey.Key // Maps the public key hex value to the CSA Key
	scryptParams utils.ScryptParams
}

func newCSAKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *CSA {
	return &CSA{
		orm:          NewCSAORM(db),
		keys:         make(map[string]*csakey.Key),
		scryptParams: scryptParams,
		mu:           new(sync.RWMutex),
	}
}

// CreateCSAKey creates a new CSA key
func (ks *CSA) CreateCSAKey() (*csakey.Key, error) {
	// Ensure you can only have one CSA at a time. This is a temporary
	// restriction until we are able to handle multiple CSA keys in the
	// communication channel
	ks.mu.Lock()
	defer ks.mu.Unlock()

	count, err := ks.orm.CountCSAKeys()
	if err != nil {
		return nil, err
	}

	if count >= 1 {
		return nil, ErrCSAKeyExists
	}

	key, err := csakey.New(ks.password, ks.scryptParams)
	if err != nil {
		return nil, err
	}

	id, err := ks.orm.CreateCSAKey(context.Background(), key)
	if err != nil {
		return nil, err
	}

	key, err = ks.orm.GetCSAKey(context.Background(), id)
	if err != nil {
		return nil, err
	}

	err = ks.unlockAndAddKey(key, ks.password)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// ListCSAKeys lists all CSA keys.
func (ks *CSA) ListCSAKeys() ([]csakey.Key, error) {
	return ks.orm.ListCSAKeys(context.Background())
}

// CountCSAKeys counts the total number of CSA keys.
func (ks *CSA) CountCSAKeys() (int64, error) {
	return ks.orm.CountCSAKeys()
}

// Unsafe_GetUnlockedKey gets the unlocked private key in the keystore.
//
// Ideally we do not want to expose private keys outside of the keystore,
// however we need to pass this priv key to the wsrpc library in order to dial
// the server. When wsrpc is updated to allow an interface to be passed in, we
// can implement that interface here to provide the private key.
func (ks *CSA) Unsafe_GetUnlockedPrivateKey(pubkey crypto.PublicKey) ([]byte, error) {
	return ks.keys[pubkey.String()].Unsafe_GetPrivateKey()
}

func (ks *CSA) Unlock(password string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var errs error

	keys, err := ks.ListCSAKeys()
	if err != nil {
		return err
	}

	for i := range keys {
		logger.Debugw("Unlocked CSA Key", "publicKey", keys[i].PublicKey)
		err := ks.unlockAndAddKey(&keys[i], password)
		errs = multierr.Append(errs, err)
	}

	ks.password = password
	return nil
}

func (ks *CSA) unlockAndAddKey(key *csakey.Key, password string) error {
	// DEV: caller must hold lock
	err := key.Unlock(password)
	if err != nil {
		return err
	}

	ks.keys[key.PublicKey.String()] = key
	return nil
}
