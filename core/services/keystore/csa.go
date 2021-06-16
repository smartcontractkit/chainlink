package keystore

import (
	"context"
	"crypto"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type CSA struct {
	orm          csaORM
	password     string
	keys         map[crypto.PublicKey]*csakey.Key
	scryptParams utils.ScryptParams
	mu           *sync.RWMutex
}

func newCSAKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *CSA {
	return &CSA{
		orm:          NewCSAORM(db),
		keys:         make(map[crypto.PublicKey]*csakey.Key),
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
		return nil, errors.New("can only have 1 CSA key")
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

func (ks *CSA) Unlock(password string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var errs error

	keys, err := ks.ListCSAKeys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := ks.unlockAndAddKey(&key, password)
		errs = multierr.Append(errs, err)
	}

	ks.password = password
	return nil
}

func (ks *CSA) unlockAndAddKey(key *csakey.Key, password string) error {
	// DEV: caller must hold lock
	err := key.Unlock(password)
	if err != nil {
		ks.keys[key.PublicKey] = key
	}
	return err
}
