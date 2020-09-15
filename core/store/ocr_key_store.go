package store

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
)

// OCRKeyStore tracks auxiliary OCR secret keys, and generates their OCR proofs
//
// OCR proofs need access to the actual secret key, which geth does not expose.
// Similar to the way geth's KeyStore exposes signing capability, OCRKeyStore
// exposes OCR proof generation without the caller needing explicit knowledge of
// the secret key.
type OCRKeyStore struct {
	lock  sync.RWMutex
	keys  InMemoryOCRKeyStore
	store *Store
}

// InMemoryOCRKeyStore holds all the unlocked OCRPrivateKey, mapped by the IDs of their
// encrypted DB records
type InMemoryOCRKeyStore = map[int32]*ocrkey.OCRPrivateKey

// NewOCRKeyStore returns an empty OCRKeyStore
func NewOCRKeyStore(store *Store) *OCRKeyStore {
	return &OCRKeyStore{
		lock:  sync.RWMutex{},
		keys:  make(InMemoryOCRKeyStore),
		store: store,
	}
}

// CreateKey creates a new private key, persists it to the DB, and adds it
// to the in-memory keystore
func (ks *OCRKeyStore) CreateKey(auth string) (*ocrkey.OCRPrivateKey, error) {
	return ks.createKey(auth, ocrkey.DefaultScryptParams)
}

// CreateFastKeyXXXTestingOnly creates a new private key with weak encryption parameters and should
// only be used for testing
func (ks *OCRKeyStore) CreateWeakKeyXXXTestingOnly(auth string) (*ocrkey.OCRPrivateKey, error) {
	return ks.createKey(auth, ocrkey.FastScryptParams)
}

func (ks *OCRKeyStore) createKey(auth string, scryptParams ocrkey.ScryptParams) (*ocrkey.OCRPrivateKey, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	key, err := ocrkey.NewOCRPrivateKey()
	if err != nil {
		return nil, err
	}
	if err = ks.saveToDB(key, auth, scryptParams); err != nil {
		return nil, err
	}
	if err = ks.addToKeystore(key); err != nil {
		return nil, err
	}
	return key, nil
}

// Unlock tries to unlock each ocr key in the db, using the given pass phrase,
// and returns any keys it manages to unlock. It returns an error if none were unlocked.
func (ks *OCRKeyStore) Unlock(phrase string) (keysUnlocked []int32, merr error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	keys, err := ks.store.FindEncryptedOCRKey()
	if err != nil {
		return nil, errors.Wrap(err, "while retrieving ocr keys from db")
	}
	for _, k := range keys {
		key, err := k.Decrypt(phrase)
		if err != nil {
			merr = multierr.Append(merr, err)
			continue
		}
		ks.addToKeystore(key)
		keysUnlocked = append(keysUnlocked, key.ID)
	}
	if len(keysUnlocked) == 0 {
		return keysUnlocked, errors.Wrap(merr, "unable to unlock any keys")
	}
	return keysUnlocked, nil
}

// Has returns a bool indicating whether an ID is present in the keystore or not
func (ks *OCRKeyStore) Has(id int32) bool {
	_, found := ks.keys[id]
	return found
}

// Get gets a key from the keystore by it's ID, errors if not found
func (ks *OCRKeyStore) Get(id int32) (*ocrkey.OCRPrivateKey, error) {
	key, found := ks.keys[id]
	if !found {
		return nil, fmt.Errorf("public key %s is not unlocked; can't forget it", key)
	}
	return key, nil
}

// Forget removes the in-memory copy of the secret key of k, or errors if not
// present.
func (ks *OCRKeyStore) Forget(key *ocrkey.OCRPrivateKey) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if _, found := ks.keys[key.ID]; !found {
		return fmt.Errorf("public key %s is not unlocked; can't forget it", key)
	}
	delete(ks.keys, key.ID)
	return nil
}

// Delete removes key from in-memory keystore and from DB
func (ks *OCRKeyStore) Delete(key *ocrkey.OCRPrivateKey) (err error) {
	if err = ks.Forget(key); err != nil {
		return err
	}
	encryptedKey := ocrkey.EncryptedOCRPrivateKey{ID: key.ID} // don't actually encrypt, we just need the ID to delete
	if err = ks.store.DeleteEncryptedOCRKeys(&encryptedKey); err != nil {
		return err
	}
	return nil
}

// saveToDB encrypts an OPCR private key and persists it to the DB
func (ks *OCRKeyStore) saveToDB(key *ocrkey.OCRPrivateKey, auth string, scryptParams ocrkey.ScryptParams) error {
	encryptedKey, err := key.Encrypt(auth, scryptParams)
	if err != nil {
		return errors.Wrap(err, "encrypting new OCR key")
	}
	err = ks.store.CreateEncryptedOCRKeys(encryptedKey)
	if err != nil {
		return errors.Wrap(err, "persisting new OCR key")
	}
	key.ID = encryptedKey.ID
	return nil
}

// addToKeystore adds an OCR private key to the in-memory keystore
func (ks *OCRKeyStore) addToKeystore(key *ocrkey.OCRPrivateKey) error {
	if key.ID == 0 {
		return errors.New("key is not yet saved to DB - cannot add to in-memory key store")
	}
	ks.keys[key.ID] = key
	return nil
}
