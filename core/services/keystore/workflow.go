package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
)

// ErrWorkflowKeyExists describes the error when the workflow key already exists
var ErrWorkflowKeyExists = errors.New("can only have 1 Workflow key")

type Workflow interface {
	Get(id string) (workflowkey.Key, error)
	GetAll() ([]workflowkey.Key, error)
	Create(ctx context.Context) (workflowkey.Key, error)
	Add(ctx context.Context, key workflowkey.Key) error
	Delete(ctx context.Context, id string) (workflowkey.Key, error)
	Import(ctx context.Context, keyJSON []byte, password string) (workflowkey.Key, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

type workflow struct {
	*keyManager
}

var _ Workflow = &workflow{}

func newWorkflowKeyStore(km *keyManager) *workflow {
	return &workflow{
		km,
	}
}

func (ks *workflow) Get(id string) (workflowkey.Key, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return workflowkey.Key{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *workflow) GetAll() (keys []workflowkey.Key, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.Workflow {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *workflow) Create(ctx context.Context) (workflowkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowkey.Key{}, ErrLocked
	}
	// Ensure you can only have one Workflow at a time.
	if len(ks.keyRing.Workflow) > 0 {
		return workflowkey.Key{}, ErrWorkflowKeyExists
	}

	key, err := workflowkey.New()
	if err != nil {
		return workflowkey.Key{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *workflow) Add(ctx context.Context, key workflowkey.Key) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if len(ks.keyRing.Workflow) > 0 {
		return ErrWorkflowKeyExists
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *workflow) Delete(ctx context.Context, id string) (workflowkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowkey.Key{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return workflowkey.Key{}, err
	}

	err = ks.safeRemoveKey(ctx, key)

	return key, err
}

func (ks *workflow) Import(ctx context.Context, keyJSON []byte, password string) (workflowkey.Key, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return workflowkey.Key{}, ErrLocked
	}

	key, err := workflowkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return workflowkey.Key{}, errors.Wrap(err, "WorkflowKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.Workflow[key.ID()]; found {
		return workflowkey.Key{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *workflow) Export(id string, password string) ([]byte, error) {
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

// EnsureKey verifies whether the Workflow key has been seeded, if not, it creates it.
func (ks *workflow) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.Workflow) > 0 {
		return nil
	}

	key, err := workflowkey.New()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created Workflow key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *workflow) getByID(id string) (workflowkey.Key, error) {
	key, found := ks.keyRing.Workflow[id]
	if !found {
		return workflowkey.Key{}, KeyNotFoundError{ID: id, KeyType: "Encryption"}
	}
	return key, nil
}
