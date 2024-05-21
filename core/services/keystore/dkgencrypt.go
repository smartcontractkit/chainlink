package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
)

//go:generate mockery --quiet --name DKGEncrypt --output mocks/ --case=underscore

// DKGEncrypt provides encryption keys for the DKG.
type DKGEncrypt interface {
	Get(id string) (dkgencryptkey.Key, error)
	GetAll() ([]dkgencryptkey.Key, error)
	Create(ctx context.Context) (dkgencryptkey.Key, error)
	Add(ctx context.Context, key dkgencryptkey.Key) error
	Delete(ctx context.Context, id string) (dkgencryptkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (dkgencryptkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type dkgEncrypt struct {
	*keyManager
}

func newDKGEncryptKeyStore(km *keyManager) *dkgEncrypt {
	return &dkgEncrypt{
		keyManager: km,
	}
}

var _ DKGEncrypt = &dkgEncrypt{}

// Add implements DKGEncrypt
func (d *dkgEncrypt) Add(ctx context.Context, key dkgencryptkey.Key) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isLocked() {
		return ErrLocked
	}
	return d.safeAddKey(ctx, key)
}

// Create implements DKGEncrypt
func (d *dkgEncrypt) Create(ctx context.Context) (dkgencryptkey.Key, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isLocked() {
		return dkgencryptkey.Key{}, ErrLocked
	}
	key, err := dkgencryptkey.New()
	if err != nil {
		return dkgencryptkey.Key{}, errors.Wrap(err, "dkgencryptkey.New()")
	}
	return key, d.safeAddKey(ctx, key)
}

// Delete implements DKGEncrypt
func (d *dkgEncrypt) Delete(ctx context.Context, id string) (dkgencryptkey.Key, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isLocked() {
		return dkgencryptkey.Key{}, ErrLocked
	}
	key, err := d.getByID(id)
	if err != nil {
		return dkgencryptkey.Key{}, err
	}

	err = d.safeRemoveKey(ctx, key)
	return key, errors.Wrap(err, "safe remove key")
}

// EnsureKey implements DKGEncrypt
func (d *dkgEncrypt) EnsureKey(ctx context.Context) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isLocked() {
		return ErrLocked
	}
	if len(d.keyRing.DKGEncrypt) > 0 {
		return nil
	}

	key, err := dkgencryptkey.New()
	if err != nil {
		return errors.Wrap(err, "dkgencryptkey. New()")
	}

	d.logger.Infof("Created DKGEncrypt key with ID %s", key.ID())

	return d.safeAddKey(ctx, key)
}

// Export implements DKGEncrypt
func (d *dkgEncrypt) Export(id string, password string) ([]byte, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.isLocked() {
		return nil, ErrLocked
	}
	key, err := d.getByID(id)
	if err != nil {
		return nil, err
	}
	return key.ToEncryptedJSON(password, d.scryptParams)
}

// Get implements DKGEncrypt
func (d *dkgEncrypt) Get(id string) (keys dkgencryptkey.Key, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.isLocked() {
		return dkgencryptkey.Key{}, ErrLocked
	}
	return d.getByID(id)
}

// GetAll implements DKGEncrypt
func (d *dkgEncrypt) GetAll() (keys []dkgencryptkey.Key, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range d.keyRing.DKGEncrypt {
		keys = append(keys, key)
	}
	return keys, nil
}

// Import implements DKGEncrypt
func (d *dkgEncrypt) Import(ctx context.Context, keyJSON []byte, password string) (dkgencryptkey.Key, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isLocked() {
		return dkgencryptkey.Key{}, ErrLocked
	}
	key, err := dkgencryptkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return dkgencryptkey.Key{}, errors.Wrap(err, "from encrypted json")
	}
	_, err = d.getByID(key.ID())
	if err == nil {
		return dkgencryptkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, d.keyManager.safeAddKey(ctx, key)
}

// caller must hold lock
func (d *dkgEncrypt) getByID(id string) (dkgencryptkey.Key, error) {
	key, found := d.keyRing.DKGEncrypt[id]
	if !found {
		return dkgencryptkey.Key{}, KeyNotFoundError{
			ID:      id,
			KeyType: "DKGEncrypt",
		}
	}
	return key, nil
}
