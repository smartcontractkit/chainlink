package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowencryptionkey"
)

// ErrWorkflowEncryptionKeyExists describes the error when the workflow encryption key already exists
var ErrWorkflowEncryptionKeyExists = errors.New("can only have 1 Workflow encryption key")

type WorkflowEncryption interface {
	Get(id string) (workflowencryptionkey.Key, error)
	GetAll() ([]workflowencryptionkey.Key, error)
	Create(ctx context.Context) (workflowencryptionkey.Key, error)
	Add(ctx context.Context, key workflowencryptionkey.Key) error
	Delete(ctx context.Context, id string) (workflowencryptionkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (workflowencryptionkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type workflowEncryption struct {
	*keyManager
}

var _ WorkflowEncryption = &workflowEncryption{}

func newWorkflowEncryptionKeyStore(km *keyManager) *workflowEncryption {
	return &workflowEncryption{
		km,
	}
}

func (ks *workflowEncryption) Get(id string) (workflowencryptionkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return workflowencryptionkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *workflowEncryption) GetAll() (keys []workflowencryptionkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.WorkflowEncryption {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *workflowEncryption) Create(ctx context.Context) (workflowencryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowencryptionkey.Key{}, ErrLocked
	}
	// Ensure you can only have one WorkflowEncryption at a time.
	if len(ks.keyRing.WorkflowEncryption) > 0 {
		return workflowencryptionkey.Key{}, ErrWorkflowEncryptionKeyExists
	}

	key, err := workflowencryptionkey.New()
	if err != nil {
		return workflowencryptionkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *workflowEncryption) Add(ctx context.Context, key workflowencryptionkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.WorkflowEncryption) > 0 {
		return ErrWorkflowEncryptionKeyExists
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *workflowEncryption) Delete(ctx context.Context, id string) (workflowencryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowencryptionkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return workflowencryptionkey.Key{}, err
	}

	err = ks.safeRemoveKey(ctx, key)

	return key, err
}

func (ks *workflowEncryption) Import(ctx context.Context, keyJSON []byte, password string) (workflowencryptionkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowencryptionkey.Key{}, ErrLocked
	}
	key, err := workflowencryptionkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return workflowencryptionkey.Key{}, errors.Wrap(err, "WorkflowEncryptionKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.WorkflowEncryption[key.ID()]; found {
		return workflowencryptionkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *workflowEncryption) Export(id string, password string) ([]byte, error) {
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

// EnsureKey verifies whether the WorkflowEncryption key has been seeded, if not, it creates it.
func (ks *workflowEncryption) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.WorkflowEncryption) > 0 {
		return nil
	}

	key, err := workflowencryptionkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Encryption key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *workflowEncryption) getByID(id string) (workflowencryptionkey.Key, error) {
	key, found := ks.keyRing.WorkflowEncryption[id]
	if !found {
		return workflowencryptionkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Encryption"}
	}
	return key, nil
}
