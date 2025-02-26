package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
)

type Tron interface {
	Get(id string) (tronkey.Key, error)
	GetAll() ([]tronkey.Key, error)
	Create(ctx context.Context) (tronkey.Key, error)
	Add(ctx context.Context, key tronkey.Key) error
	Delete(ctx context.Context, id string) (tronkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (tronkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
	Sign(ctx context.Context, id string, msg []byte) (signature []byte, err error)
}

type tron struct {
	*keyManager
}

var _ Tron = &tron{}

func newTronKeyStore(km *keyManager) *tron {
	return &tron{
		km,
	}
}

func (ks *tron) Get(id string) (tronkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return tronkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *tron) GetAll() (keys []tronkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Tron {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *tron) Create(ctx context.Context) (tronkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tronkey.Key{}, ErrLocked
	}
	key, err := tronkey.New()
	if err != nil {
		return tronkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *tron) Add(ctx context.Context, key tronkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Tron[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *tron) Delete(ctx context.Context, id string) (tronkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tronkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return tronkey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *tron) Import(ctx context.Context, keyJSON []byte, password string) (tronkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return tronkey.Key{}, ErrLocked
	}
	key, err := tronkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return tronkey.Key{}, errors.Wrap(err, "TronKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Tron[key.ID()]; found {
		return tronkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *tron) Export(id string, password string) ([]byte, error) {
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

func (ks *tron) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.Tron) > 0 {
		return nil
	}

	key, err := tronkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Tron key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *tron) getByID(id string) (tronkey.Key, error) {
	key, found := ks.keyRing.Tron[id]
	if !found {
		return tronkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Tron"}
	}
	return key, nil
}

func (ks *tron) Sign(_ context.Context, id string, msg []byte) (signature []byte, err error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}
	// loopp spec requires passing nil hash to check existence of id
	if msg == nil {
		return nil, nil
	}
	return k.Sign(msg)
}

// TronLOOPSigner implements the [github.com/smartcontractkit/chainlink-common/pkg/loop.Keystore] interface and
// handles signing for Tron messages.
type TronLOOPSigner struct {
	Tron
}

var _ loop.Keystore = &TronLOOPSigner{}

func (lk *TronLOOPSigner) Accounts(ctx context.Context) ([]string, error) {
	keys, err := lk.GetAll()
	if err != nil {
		return nil, err
	}

	accounts := []string{}
	for _, k := range keys {
		accounts = append(accounts, k.ID())
	}

	return accounts, nil
}
