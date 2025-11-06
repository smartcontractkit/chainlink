package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgrecipientkey"
)

// ErrDKGRecipientKeyExists describes the error when the DKG recipient key already exists
var ErrDKGRecipientKeyExists = errors.New("can only have 1 DKG recipient key")

type DKGRecipient interface {
	Get(id string) (dkgrecipientkey.Key, error)
	GetAll() ([]dkgrecipientkey.Key, error)
	Create(ctx context.Context) (dkgrecipientkey.Key, error)
	Add(ctx context.Context, key dkgrecipientkey.Key) error
	Delete(ctx context.Context, id string) (dkgrecipientkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (dkgrecipientkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type dkgRecipient struct {
	*keyManager
}

var _ DKGRecipient = &dkgRecipient{}

func newDKGRecipientKeyStore(km *keyManager) *dkgRecipient {
	return &dkgRecipient{
		keyManager: km,
	}
}

func (ks *dkgRecipient) Get(id string) (dkgrecipientkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return dkgrecipientkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *dkgRecipient) GetAll() (keys []dkgrecipientkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.DKGRecipient {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *dkgRecipient) Create(ctx context.Context) (dkgrecipientkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return dkgrecipientkey.Key{}, ErrLocked
	}
	// Ensure you can only have one DKGRecipient at a time.
	if len(ks.keyRing.DKGRecipient) > 0 {
		return dkgrecipientkey.Key{}, ErrDKGRecipientKeyExists
	}

	key, err := dkgrecipientkey.New()
	if err != nil {
		return dkgrecipientkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *dkgRecipient) Add(ctx context.Context, key dkgrecipientkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.DKGRecipient) > 0 {
		return ErrDKGRecipientKeyExists
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *dkgRecipient) Delete(ctx context.Context, id string) (dkgrecipientkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return dkgrecipientkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return dkgrecipientkey.Key{}, err
	}

	err = ks.safeRemoveKey(ctx, key)

	return key, err
}

func (ks *dkgRecipient) Import(ctx context.Context, keyJSON []byte, password string) (dkgrecipientkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return dkgrecipientkey.Key{}, ErrLocked
	}

	key, err := dkgrecipientkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return dkgrecipientkey.Key{}, errors.Wrap(err, "dkgRecipient#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.DKGRecipient[key.ID()]; found {
		return dkgrecipientkey.Key{}, fmt.Errorf("%w: key with ID %s already exists", ErrKeyExists, key.ID())
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *dkgRecipient) Export(id string, password string) ([]byte, error) {
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

// EnsureKey verifies whether the DKGRecipient key has been seeded, if not, it creates it.
func (ks *dkgRecipient) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.DKGRecipient) > 0 {
		return nil
	}

	key, err := dkgrecipientkey.New()
	if err != nil {
		return err
	}

	ks.announce(key)

	return ks.safeAddKey(ctx, key)
}

func (ks *dkgRecipient) getByID(id string) (dkgrecipientkey.Key, error) {
	key, found := ks.keyRing.DKGRecipient[id]
	if !found {
		return dkgrecipientkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Encryption"}
	}
	return key, nil
}
