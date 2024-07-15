package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
)

type Aptos interface {
	Get(id string) (aptoskey.Key, error)
	GetAll() ([]aptoskey.Key, error)
	Create(ctx context.Context) (aptoskey.Key, error)
	Add(ctx context.Context, key aptoskey.Key) error
	Delete(ctx context.Context, id string) (aptoskey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (aptoskey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
	Sign(ctx context.Context, id string, msg []byte) (signature []byte, err error)
}

type aptos struct {
	*keyManager
}

var _ Aptos = &aptos{}

func newAptosKeyStore(km *keyManager) *aptos {
	return &aptos{
		km,
	}
}

func (ks *aptos) Get(id string) (aptoskey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return aptoskey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *aptos) GetAll() (keys []aptoskey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Aptos {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *aptos) Create(ctx context.Context) (aptoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return aptoskey.Key{}, ErrLocked
	}
	key, err := aptoskey.New()
	if err != nil {
		return aptoskey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *aptos) Add(ctx context.Context, key aptoskey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Aptos[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *aptos) Delete(ctx context.Context, id string) (aptoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return aptoskey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return aptoskey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *aptos) Import(ctx context.Context, keyJSON []byte, password string) (aptoskey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return aptoskey.Key{}, ErrLocked
	}
	key, err := aptoskey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return aptoskey.Key{}, errors.Wrap(err, "AptosKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Aptos[key.ID()]; found {
		return aptoskey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *aptos) Export(id string, password string) ([]byte, error) {
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

func (ks *aptos) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Aptos) > 0 {
		return nil
	}

	key, err := aptoskey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Aptos key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *aptos) Sign(_ context.Context, id string, msg []byte) (signature []byte, err error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}
	return k.Sign(msg)
}

func (ks *aptos) getByID(id string) (aptoskey.Key, error) {
	key, found := ks.keyRing.Aptos[id]
	if !found {
		return aptoskey.Key{}, KeyNotFoundError{ID: id, KeyType: "Aptos"}
	}
	return key, nil
}

// AptosSigner implements [github.com/smartcontractkit/chainlink-common/pkg/loop.Keystore] interface and the requirements
// Handles signing for Apots Messages
type AptosLooppSigner struct {
	Aptos
}

var _ loop.Keystore = &AptosLooppSigner{}

// Returns a list of Aptos Public Keys
func (s *AptosLooppSigner) Accounts(ctx context.Context) (accounts []string, err error) {
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.ID())
	}
	return
}
