package store

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"chainlink/core/store/models/vrfkey"
)

// CreateKey an immediately unlocked key, saved in DB encrypted with phrase.
func (ks *VRFKeyStore) CreateKey(phrase string, p ...vrfkey.ScryptParams,
) (*vrfkey.PublicKey, error) {
	key := vrfkey.CreateKey()
	if err := ks.Store(key, phrase, p...); err != nil {
		return nil, err
	}
	return &key.PublicKey, nil
}

func (ks *VRFKeyStore) CreateWeakInMemoryEncryptedKeyXXXTestingOnly(
	phrase string) (*vrfkey.EncryptedSecretKey, error) {
	key := vrfkey.CreateKey()
	encrypted, err := key.Encrypt(phrase, vrfkey.FastScryptParams)
	if err != nil {
		return nil, errors.Wrap(err, "while creating testing key")
	}
	return encrypted, nil
}

// Store saves key to this keystore, and to the DB encrypted with phrase.
func (ks *VRFKeyStore) Store(key *vrfkey.PrivateKey, phrase string,
	p ...vrfkey.ScryptParams) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	encrypted, err := key.Encrypt(phrase, p...)
	if err != nil {
		return errors.Wrap(err, "failed to serialize encrypted key")
	}
	if err := ks.store.FirstOrCreateEncryptedSecretVRFKey(encrypted); err != nil {
		return errors.Wrap(err, "failed to  save encrypted key to db")
	}
	ks.keys[key.PublicKey] = *key
	return nil
}

// StoreInMemoryXXXTestingOnly memorizes key, only in in-memory store.
func (ks *VRFKeyStore) StoreInMemoryXXXTestingOnly(key *vrfkey.PrivateKey) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	ks.keys[key.PublicKey] = *key
}

var zeroPublicKey = vrfkey.PublicKey{}

// Delete removes keys with this public key from the keystore and the DB, if present.
func (ks *VRFKeyStore) Delete(key *vrfkey.PublicKey) (err error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if *key == zeroPublicKey {
		return fmt.Errorf("cannot delete the empty public key")
	}
	if _, found := ks.keys[*key]; found {
		err = ks.forget(key) // Destroy in-memory representation of key
		delete(ks.keys, *key)
	}
	matches, err := ks.get(key)
	if err != nil {
		return errors.Wrapf(err, "while checking for existence of key %s", key.String())
	}
	if len(matches) == 0 {
		return AttemptToDeleteNonExistentKeyFromDB
	}
	return multierr.Append(err, ks.store.ORM.DeleteEncryptedSecretVRFKey(
		&vrfkey.EncryptedSecretKey{PublicKey: *key}))
}

// Import adds this encrypted key to the DB and in-memory store, or errors
func (ks *VRFKeyStore) Import(keyjson []byte, auth string) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	enckey := &vrfkey.EncryptedSecretKey{}
	if err := json.Unmarshal(keyjson, enckey); err != nil {
		return fmt.Errorf("could not parse %s as EncryptedSecretKey json", keyjson)
	}
	extantMatchingKeys, err := ks.get(&enckey.PublicKey)
	if err != nil {
		return errors.Wrapf(err, "while checking for matching extant key in DB")
	}
	if len(extantMatchingKeys) != 0 {
		return MatchingVRFKeyError
	}
	key, err := enckey.Decrypt(auth)
	if err != nil {
		return err
	}
	if err := ks.store.FirstOrCreateEncryptedSecretVRFKey(enckey); err != nil {
		return err
	}
	ks.keys[key.PublicKey] = *key
	return nil
}

// get retrieves all EncryptedSecretKey's associated with k, or all encrypted
// keys if k is nil, or errors. Caller is responsible for locking the store
func (ks *VRFKeyStore) get(k *vrfkey.PublicKey) ([]*vrfkey.EncryptedSecretKey, error) {
	var where []vrfkey.EncryptedSecretKey
	if k != nil { // Search for this specific public key
		where = append(where, vrfkey.EncryptedSecretKey{PublicKey: *k})
	}
	keys, err := ks.store.FindEncryptedSecretVRFKeys(where...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find public key %s in DB", k)
	}
	return keys, nil
}

// Get retrieves all EncryptedSecretKey's associated with k, or all encrypted
// keys if k is nil, or errors
//
// (There could be more than one match to a public key, if saved multiple times
// or under different passwords.)
func (ks *VRFKeyStore) Get(k *vrfkey.PublicKey) ([]*vrfkey.EncryptedSecretKey, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	return ks.get(k)
}

// Export returns the encrypted key data for the given public key. (Could be
// more than one export, if key has been imported more than once.) enckeys and
// merr can both be non-trivial, if some keys were retrievable and some not.
func (ks *VRFKeyStore) Export(k *vrfkey.PublicKey) (enckeys [][]byte, merr error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	enc, err := ks.get(k)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve %s from db", k)
	}
	for _, enckey := range enc {
		keyjson, err := json.Marshal(enckey)
		if err != nil {
			merr = multierr.Append(err, errors.Wrapf(err, "could not marshal %+v to json", enckey))
		} else {
			enckeys = append(enckeys, keyjson)
		}
	}
	return enckeys, merr
}

// ListKey lists the public keys contained in the db
func (ks *VRFKeyStore) ListKeys() (publicKeys []*vrfkey.PublicKey, err error) {
	enc, err := ks.Get(nil)
	if err != nil {
		return nil, errors.Wrapf(err, "while listing db keys")
	}
	for _, enckey := range enc {
		publicKeys = append(publicKeys, &enckey.PublicKey)
	}
	return publicKeys, nil
}

// MatchingVRFKeyError is returned when Import attempts to import key with a
// PublicKey matching one already in the database
var MatchingVRFKeyError = errors.New(
	`key with matching public key already stored in DB`)

// AttemptToDeleteNonExistentKeyFromDB is returned when Delete is asked to
// delete a key it can't find in the DB.
var AttemptToDeleteNonExistentKeyFromDB = errors.New("key is not present in DB")
