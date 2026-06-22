package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type Stellar interface {
	Get(id string) (stellarkey.Key, error)
	GetAll() ([]stellarkey.Key, error)
	Create(ctx context.Context) (stellarkey.Key, error)
	Add(ctx context.Context, key stellarkey.Key) error
	Delete(ctx context.Context, id string) (stellarkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (stellarkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
	Sign(ctx context.Context, id string, msg []byte) (signature []byte, err error)
}

type stellar struct {
	*keyManager
}

var _ Stellar = &stellar{}

func newStellarKeyStore(km *keyManager) *stellar {
	return &stellar{
		km,
	}
}

func (ks *stellar) Get(id string) (stellarkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return stellarkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *stellar) GetAll() (keys []stellarkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Stellar {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *stellar) Create(ctx context.Context) (stellarkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stellarkey.Key{}, ErrLocked
	}
	key, err := stellarkey.New()
	if err != nil {
		return stellarkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *stellar) Add(ctx context.Context, key stellarkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.Stellar[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *stellar) Delete(ctx context.Context, id string) (stellarkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stellarkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return stellarkey.Key{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *stellar) Import(ctx context.Context, keyJSON []byte, password string) (stellarkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return stellarkey.Key{}, ErrLocked
	}
	key, err := stellarkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return stellarkey.Key{}, errors.Wrap(err, "StellarKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Stellar[key.ID()]; found {
		return stellarkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *stellar) Export(id string, password string) ([]byte, error) {
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

func (ks *stellar) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Stellar) > 0 {
		return nil
	}

	key, err := stellarkey.New()
	if err != nil {
		return err
	}

	ks.announce(key)

	return ks.safeAddKey(ctx, key)
}

func (ks *stellar) Sign(_ context.Context, id string, msg []byte) (signature []byte, err error) {
	k, err := ks.Get(id)
	if err != nil {
		return nil, err
	}
	return k.Sign(msg)
}

func (ks *stellar) getByID(id string) (stellarkey.Key, error) {
	key, found := ks.keyRing.Stellar[id]
	if !found {
		return stellarkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Stellar"}
	}
	return key, nil
}

// StellarLooppSigner implements [github.com/smartcontractkit/chainlink-common/pkg/loop.Keystore]
// and handles signing for Stellar messages. The Stellar relayer/TXM calls
// Sign(ctx, account, data) where account is the StrKey "G..." address.
type StellarLooppSigner struct {
	Stellar
	core.UnimplementedKeystore
}

var _ loop.Keystore = &StellarLooppSigner{}

// Accounts returns a list of Stellar StrKey "G..." account addresses.
func (s *StellarLooppSigner) Accounts(ctx context.Context) (accounts []string, err error) {
	ks, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		accounts = append(accounts, k.ID())
	}
	return
}

func (s *StellarLooppSigner) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	return s.Stellar.Sign(ctx, account, data)
}
