package keystore

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
)

//go:generate mockery --name OCR2 --output mocks/ --case=underscore

type OCR2 interface {
	Get(id string) (ocr2key.KeyBundle, error)
	GetAll() ([]ocr2key.KeyBundle, error)
	Create() (ocr2key.KeyBundle, error)
	Delete(id string) error
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
		return ocr2key.KeyBundle{}, ErrLocked
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

func (ks ocr2) getByID(id string) (ocr2key.KeyBundle, error) {
	key, found := ks.keyRing.OCR2[id]
	if !found {
		return ocr2key.KeyBundle{}, fmt.Errorf("unable to find OCR key with id %s", id)
	}
	return key, nil
}

func (ks ocr2) Create() (ocr2key.KeyBundle, error) {
	key, err := ocr2key.NewKeyBundle()
	if err != nil {
		return ocr2key.KeyBundle{}, err
	}
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if ks.isLocked() {
		return ocr2key.KeyBundle{}, ErrLocked
	}
	return *key, ks.safeAddKey(*key)
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
