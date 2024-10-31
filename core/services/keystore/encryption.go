package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/encryptionkey"
)

// ErrEncryptionKeyExists describes the error when the Encryption key already exists
var ErrEncryptionKeyExists = errors.New("can only have 1 Encryption key")

type Encryption interface {
	Get(id string) (encryptionkey.Key, error)
	GetAll() ([]encryptionkey.Key, error)
	Create(ctx context.Context) (encryptionkey.Key, error)
	Add(ctx context.Context, key encryptionkey.Key) error
	Delete(ctx context.Context, id string) (encryptionkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (encryptionkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type encryption struct {
	*keyManager
}

var _ Encryption = &encryption{}

func newEncryptionKeyStore(km *keyManager) *encryption {
	return &encryption{
		km,
	}
}

func (ks *encryption) Get(id string) (encryptionkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return encryptionkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *encryption) GetAll() (keys []encryptionkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Encryption {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *encryption) Create(ctx context.Context) (encryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return encryptionkey.Key{}, ErrLocked
	}
	// Ensure you can only have one Encryption at a time.
	if len(ks.keyRing.Encryption) > 0 {
		return encryptionkey.Key{}, ErrEncryptionKeyExists
	}

	key, err := encryptionkey.New()
	if err != nil {
		return encryptionkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *encryption) Add(ctx context.Context, key encryptionkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Encryption) > 0 {
		return ErrEncryptionKeyExists
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *encryption) Delete(ctx context.Context, id string) (encryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return encryptionkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return encryptionkey.Key{}, err
	}

	err = ks.safeRemoveKey(ctx, key)

	return key, err
}

func (ks *encryption) Import(ctx context.Context, keyJSON []byte, password string) (encryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return encryptionkey.Key{}, ErrLocked
	}
	key, err := encryptionkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return encryptionkey.Key{}, errors.Wrap(err, "EncryptionKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Encryption[key.ID()]; found {
		return encryptionkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *encryption) Export(id string, password string) ([]byte, error) {
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

// EnsureKey verifies whether the Encryption key has been seeded, if not, it creates it.
func (ks *encryption) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.Encryption) > 0 {
		return nil
	}

	key, err := encryptionkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Encryption key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *encryption) getByID(id string) (encryptionkey.Key, error) {
	key, found := ks.keyRing.Encryption[id]
	if !found {
		return encryptionkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Encryption"}
	}
	return key, nil
}
