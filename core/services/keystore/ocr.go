package keystore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
)

type OCR interface {
	Get(id string) (ocrkey.KeyV2, error)
	GetAll() ([]ocrkey.KeyV2, error)
	Create(ctx context.Context) (ocrkey.KeyV2, error)
	Add(ctx context.Context, key ocrkey.KeyV2) error
	Delete(ctx context.Context, id string) (ocrkey.KeyV2, error)
	Import(ctx context.Context, keyJSON []byte, password string) (ocrkey.KeyV2, error)
	Export(id string, password string) ([]byte, error)
	EnsureKey(ctx context.Context) error
}

// KeyNotFoundError is returned when we don't find a requested key
type KeyNotFoundError struct {
	ID      string
	KeyType string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("unable to find %s key with id %s", e.KeyType, e.ID)
}

type ocr struct {
	*keyManager
}

var _ OCR = &ocr{}

func newOCRKeyStore(km *keyManager) *ocr {
	return &ocr{
		km,
	}
}

func (ks *ocr) Get(id string) (ocrkey.KeyV2, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return ocrkey.KeyV2{}, ErrLocked
	}
	return ks.getByID(id)
}

func (ks *ocr) GetAll() (keys []ocrkey.KeyV2, _ error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	for _, key := range ks.keyRing.OCR {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks *ocr) Create(ctx context.Context) (ocrkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ocrkey.KeyV2{}, ErrLocked
	}
	key, err := ocrkey.NewV2()
	if err != nil {
		return ocrkey.KeyV2{}, err
	}
	return key, ks.safeAddKey(ctx, key)
}

func (ks *ocr) Add(ctx context.Context, key ocrkey.KeyV2) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.OCR[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(ctx, key)
}

func (ks *ocr) Delete(ctx context.Context, id string) (ocrkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ocrkey.KeyV2{}, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return ocrkey.KeyV2{}, err
	}
	err = ks.safeRemoveKey(ctx, key)
	return key, err
}

func (ks *ocr) Import(ctx context.Context, keyJSON []byte, password string) (ocrkey.KeyV2, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ocrkey.KeyV2{}, ErrLocked
	}
	key, err := ocrkey.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return ocrkey.KeyV2{}, errors.Wrap(err, "OCRKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.OCR[key.ID()]; found {
		return ocrkey.KeyV2{}, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(ctx, key)
}

func (ks *ocr) Export(id string, password string) ([]byte, error) {
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

// EnsureKey verifies whether the OCR key has been seeded, if not, it creates it.
func (ks *ocr) EnsureKey(ctx context.Context) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	if len(ks.keyRing.OCR) > 0 {
		return nil
	}

	key, err := ocrkey.NewV2()
	if err != nil {
		return err
	}

	ks.logger.Infof("Created OCR key with ID %s", key.ID())

	return ks.safeAddKey(ctx, key)
}

func (ks *ocr) getByID(id string) (ocrkey.KeyV2, error) {
	key, found := ks.keyRing.OCR[id]
	if !found {
		return ocrkey.KeyV2{}, KeyNotFoundError{ID: id, KeyType: "OCR"}
	}
	return key, nil
}
