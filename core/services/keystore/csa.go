package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

//go:generate mockery --quiet --name CSA --output mocks/ --case=underscore

// ErrCSAKeyExists describes the error when the CSA key already exists
var ErrCSAKeyExists = errors.New("can only have 1 CSA key")

// type CSAKeystoreInterface interface {
type CSA interface {
	Get(id string) (csakey.KeyV2, error)
	GetAll() ([]csakey.KeyV2, error)
	Create(ctx context.Context) (csakey.KeyV2, error)
	Add(ctx context.Context, key csakey.KeyV2) error
	Delete(ctx context.Context, id string) (csakey.KeyV2, error)
	Import(ctx context.Context, keyJSON []byte, password string) (csakey.KeyV2, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type csa struct {
	*keyManager
}

var _ CSA = &csa{}

func newCSAKeyStore(km *keyManager) *csa {
	return &csa{
		km,
	}
}

func (ks *csa) Get(id string) (csakey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return csakey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *csa) GetAll() (keys []csakey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.CSA {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *csa) Create(ctx context.Context) (csakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return csakey.KeyV2{}, ErrLocked
	}
	// Ensure you can only have one CSA at a time. This is a temporary
	// restriction until we are able to handle multiple CSA keys in the
	// communication channel
	if len(ks.keyRing.CSA) > 0 {
		return csakey.KeyV2{}, ErrCSAKeyExists
	}
	key, err := csakey.NewV2()
	if err != nil {
		return csakey.KeyV2{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *csa) Add(ctx context.Context, key csakey.KeyV2) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.CSA) > 0 {
		return ErrCSAKeyExists
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *csa) Delete(ctx context.Context, id string) (csakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return csakey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return csakey.KeyV2{}, err
	}

	err = ks.safeRemoveKey(ctx, key)

	return key, err
}

func (ks *csa) Import(ctx context.Context, keyJSON []byte, password string) (csakey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return csakey.KeyV2{}, ErrLocked
	}
	key, err := csakey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return csakey.KeyV2{}, errors.Wrap(err, "CSAKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.CSA[key.ID()]; found {
		return csakey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *csa) Export(id string, password string) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return key.ToEncryptedJSON(password, ks.scryptParams)
}

// EnsureKey verifies whether the CSA key has been seeded, if not, it creates it.
func (ks *csa) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.CSA) > 0 {
		return nil
	}

	key, err := csakey.NewV2()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created CSA key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *csa) getByID(id string) (csakey.KeyV2, error) {
	key, found := ks.keyRing.CSA[id]
	if !found {
		return csakey.KeyV2{}, KeyNotFoundError{ID: id, KeyType: "CSA"}
	}
	return key, nil
}
