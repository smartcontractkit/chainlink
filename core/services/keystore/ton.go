package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tonkey"
)

type Ton interface {
	Get(id string) (tonkey.Key, error)
	GetAll() ([]tonkey.Key, error)
	Create(ctx context.Context) (tonkey.Key, error)
	Add(ctx context.Context, key tonkey.Key) error
	Delete(ctx context.Context, id string) (tonkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (tonkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
	Sign(ctx context.Context, id string, msg []byte) (signature []byte, err error)
}

type ton struct {
	*keyManager
}

var _ Ton = &ton{}

func newTonKeyStore(km *keyManager) *ton {
	return &ton{
		km,
	}
}

func (ks *ton) Get(id string) (tonkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return tonkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *ton) GetAll() (keys []tonkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Ton {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *ton) Create(ctx context.Context) (tonkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tonkey.Key{}, ErrLocked
	}
	key, err := tonkey.New()
	if err != nil {
		return tonkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *ton) Add(ctx context.Context, key tonkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Ton[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *ton) Delete(ctx context.Context, id string) (tonkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tonkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return tonkey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *ton) Import(ctx context.Context, keyJSON []byte, password string) (tonkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tonkey.Key{}, ErrLocked
	}
	key, err := tonkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return tonkey.Key{}, errors.Wrap(err, "TonKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Ton[key.ID()]; found {
		return tonkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *ton) Export(id string, password string) ([]byte, error) {
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

func (ks *ton) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Ton) > 0 {
		return nil
	}

	key, err := tonkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Ton key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *ton) Sign(_ context.Context, id string, msg []byte) (signature []byte, err error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}
	return k.Sign(msg)
}

func (ks *ton) getByID(id string) (tonkey.Key, error) {
	key, found := ks.keyRing.Ton[id]
	if !found {
		return tonkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Ton"}
	}
	return key, nil
}

// TonLooppSigner implements the [github.com/smartcontractkit/chainlink-common/pkg/loop.Keystore] interface and
// handles signing for Ton messages.
type TonLooppSigner struct {
	Ton
}

var _ loop.Keystore = &TonLooppSigner{}

// Returns a list of Ton Public Keys
func (s *TonLooppSigner) Accounts(ctx context.Context) (accounts []string, err error) {
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.ID())
	}
	return
}
