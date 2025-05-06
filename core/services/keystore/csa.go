package keystore

import (
	"context"
	"crypto"
	"crypto/rand"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
)

// ErrCSAKeyExists describes the error when the CSA key already exists
var ErrCSAKeyExists = errors.New("can only have 1 CSA key")

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

var _ loop.Keystore = &CSASigner{}

type CSASigner struct {
	CSA
}

func (c CSASigner) Accounts(ctx context.Context) (accounts []string, err error) {
	keys, err := c.CSA.GetAll()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		accounts = append(accounts, key.ID())
	}
	return
}

func (c CSASigner) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	k, err := c.CSA.Get(account)
	if err != nil {
		return nil, err
	}
	// loopp spec requires passing nil hash to check existence of id
	if data == nil {
		return nil, nil
	}
	return k.Sign(rand.Reader, data, crypto.Hash(0))
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
