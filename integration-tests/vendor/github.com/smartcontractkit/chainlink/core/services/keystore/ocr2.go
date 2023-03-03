package keystore

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
)

//go:generate mockery --quiet --name OCR2 --output mocks/ --case=underscore

type OCR2 interface {
	Get(id string) (ocr2key.KeyBundle, error)
	GetAll() ([]ocr2key.KeyBundle, error)
	GetAllOfType(chaintype.ChainType) ([]ocr2key.KeyBundle, error)
	Create(chaintype.ChainType) (ocr2key.KeyBundle, error)
	Add(key ocr2key.KeyBundle) error
	Delete(id string) error
	Import(keyJSON []byte, password string) (ocr2key.KeyBundle, error)
	Export(id string, password string) ([]byte, error)
	EnsureKeys() error
}

type ocr2 struct {
	*keyManager
}

var _ OCR2 = ocr2{}

func newOCR2KeyStore(km *keyManager) ocr2 {
	return ocr2{
		km,
	}
}

func (ks ocr2) Get(id string) (ocr2key.KeyBundle, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.getByID(id)
}

func (ks ocr2) GetAll() ([]ocr2key.KeyBundle, error) {
	keys := []ocr2key.KeyBundle{}
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return keys, ErrLocked
	}
	for _, key := range ks.keyRing.OCR2 {
		keys = append(keys, key)
	}
	return keys, nil
}

func (ks ocr2) GetAllOfType(chainType chaintype.ChainType) ([]ocr2key.KeyBundle, error) {
	keys := []ocr2key.KeyBundle{}
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return keys, ErrLocked
	}
	return ks.getAllOfType(chainType)
}

func (ks ocr2) Create(chainType chaintype.ChainType) (ocr2key.KeyBundle, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	return ks.create(chainType)
}

func (ks ocr2) Add(key ocr2key.KeyBundle) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	if _, found := ks.keyRing.OCR2[key.ID()]; found {
		return fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return ks.safeAddKey(key)
}

func (ks ocr2) Delete(id string) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return err
	}
	err = ks.safeRemoveKey(key)
	return err
}

func (ks ocr2) Import(keyJSON []byte, password string) (ocr2key.KeyBundle, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ocr2key.FromEncryptedJSON(keyJSON, password)
	if err != nil {
		return nil, errors.Wrap(err, "OCRKeyStore#ImportKey failed to decrypt key")
	}
	if _, found := ks.keyRing.OCR[key.ID()]; found {
		return nil, fmt.Errorf("key with ID %s already exists", key.ID())
	}
	return key, ks.keyManager.safeAddKey(key)
}

func (ks ocr2) Export(id string, password string) ([]byte, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	if ks.isLocked() {
		return nil, ErrLocked
	}
	key, err := ks.getByID(id)
	if err != nil {
		return nil, err
	}
	return ocr2key.ToEncryptedJSON(key, password, ks.scryptParams)
}

func (ks ocr2) EnsureKeys() error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ErrLocked
	}

	for _, chainType := range chaintype.SupportedChainTypes {
		keys, err := ks.getAllOfType(chainType)
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			continue
		}

		created, err := ks.create(chainType)
		if err != nil {
			return err
		}

		ks.logger.Infof("Created OCR2 key with ID %s", created.ID())
	}

	return nil
}

func (ks ocr2) getByID(id string) (ocr2key.KeyBundle, error) {
	key, found := ks.keyRing.OCR2[id]
	if !found {
		return nil, fmt.Errorf("unable to find OCR key with id %s", id)
	}
	return key, nil
}

func (ks ocr2) getAllOfType(chainType chaintype.ChainType) ([]ocr2key.KeyBundle, error) {
	keys := []ocr2key.KeyBundle{}
	for _, key := range ks.keyRing.OCR2 {
		if key.ChainType() == chainType {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (ks ocr2) create(chainType chaintype.ChainType) (ocr2key.KeyBundle, error) {
	if !chaintype.IsSupportedChainType(chainType) {
		return nil, chaintype.NewErrInvalidChainType(chainType)
	}
	key, err := ocr2key.New(chainType)
	if err != nil {
		return nil, err
	}
	return key, ks.safeAddKey(key)
}
