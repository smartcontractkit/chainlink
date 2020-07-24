package store

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
)

// CreateKey returns a public key which is immediately unlocked in memory, and
// saved in DB encrypted with phrase. If p is given, its parameters are used for
// key derivation from the phrase.
func (ks *VRFKeyStore) CreateKey(phrase string, p ...vrfkey.ScryptParams,
) (vrfkey.PublicKey, error) {
	key := vrfkey.CreateKey()
	if err := ks.Store(key, phrase, p...); err != nil {
		return vrfkey.PublicKey{}, err
	}
	return key.PublicKey, nil
}

// CreateWeakInMemoryEncryptedKeyXXXTestingOnly is for testing only! It returns
// an encrypted key which is fast to unlock, but correspondingly easy to brute
// force. It is not persisted to the DB, because no one should be keeping such
// keys lying around.
func (ks *VRFKeyStore) CreateWeakInMemoryEncryptedKeyXXXTestingOnly(
	phrase string) (*vrfkey.EncryptedSecretKey, error) {
	key := vrfkey.CreateKey()
	encrypted, err := key.Encrypt(phrase, vrfkey.FastScryptParams)
	if err != nil {
		return nil, errors.Wrap(err, "while creating testing key")
	}
	return encrypted, nil
}

// Store saves key to ks (in memory), and to the DB, encrypted with phrase
func (ks *VRFKeyStore) Store(key *vrfkey.PrivateKey, phrase string,
	p ...vrfkey.ScryptParams) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	encrypted, err := key.Encrypt(phrase, p...)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt key")
	}
	if err := ks.store.FirstOrCreateEncryptedSecretVRFKey(encrypted); err != nil {
		return errors.Wrap(err, "failed to save encrypted key to db")
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
func (ks *VRFKeyStore) Delete(key vrfkey.PublicKey) (err error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	if key == zeroPublicKey {
		return fmt.Errorf("cannot delete the empty public key")
	}
	if _, found := ks.keys[key]; found {
		err = ks.forget(key) // Destroy in-memory representation of key
		delete(ks.keys, key)
	}
	matches, err := ks.get(key)
	if err != nil {
		return errors.Wrapf(err, "while checking for existence of key %s in DB",
			key.String())
	}
	if len(matches) == 0 {
		return AttemptToDeleteNonExistentKeyFromDB
	}
	return multierr.Append(err, ks.store.ORM.DeleteEncryptedSecretVRFKey(
		&vrfkey.EncryptedSecretKey{PublicKey: key}))
}

// Import adds this encrypted key to the DB and unlocks it in in-memory store
// with passphrase auth, and returns any resulting errors
func (ks *VRFKeyStore) Import(keyjson []byte, auth string) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	enckey := &vrfkey.EncryptedSecretKey{}
	if err := json.Unmarshal(keyjson, enckey); err != nil {
		return fmt.Errorf("could not parse %s as EncryptedSecretKey json", keyjson)
	}
	extantMatchingKeys, err := ks.get(enckey.PublicKey)
	if err != nil {
		return errors.Wrapf(err, "while checking for matching extant key in DB")
	}
	if len(extantMatchingKeys) != 0 {
		return MatchingVRFKeyError
	}
	key, err := enckey.Decrypt(auth)
	if err != nil {
		return errors.Wrapf(err,
			"while attempting to decrypt key with public key %s",
			key.PublicKey.String())
	}
	if err := ks.store.FirstOrCreateEncryptedSecretVRFKey(enckey); err != nil {
		return errors.Wrapf(err, "while saving encrypted key to DB")
	}
	ks.keys[key.PublicKey] = *key
	return nil
}

// get retrieves all EncryptedSecretKey's associated with k, or all encrypted
// keys if k is nil, or errors. Caller is responsible for locking the store
func (ks *VRFKeyStore) get(k ...vrfkey.PublicKey) ([]*vrfkey.EncryptedSecretKey,
	error) {
	if len(k) > 1 {
		return nil, errors.Errorf("can get at most one secret key at a time")
	}
	var where []vrfkey.EncryptedSecretKey
	if len(k) == 1 { // Search for this specific public key
		where = append(where, vrfkey.EncryptedSecretKey{PublicKey: k[0]})
	}
	keys, err := ks.store.FindEncryptedSecretVRFKeys(where...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find public key %s in DB", k)
	}
	return keys, nil
}

// Get retrieves all EncryptedSecretKey's associated with k, or all encrypted
// keys if k is nil, or errors
func (ks *VRFKeyStore) Get(k ...vrfkey.PublicKey) ([]*vrfkey.EncryptedSecretKey, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	return ks.get(k...)
}

func (ks *VRFKeyStore) GetSpecificKey(
	k vrfkey.PublicKey) (*vrfkey.EncryptedSecretKey, error) {
	if k == (vrfkey.PublicKey{}) {
		return nil, fmt.Errorf("can't retrieve zero key")
	}
	encryptedKey, err := ks.Get(k)
	if err != nil {
		return nil, errors.Wrapf(err, "could not retrieve %s from db", k)
	}
	if len(encryptedKey) == 0 {
		return nil, fmt.Errorf("could not find any keys with public key %s",
			k.String())
	}
	if len(encryptedKey) > 1 {
		// This is impossible, as long as the public key is the primary key on the
		// EncryptedSecretKey table.
		panic(fmt.Errorf("found more than one key with public key %s", k.String()))
	}
	return encryptedKey[0], nil
}

// ListKeys lists the public keys contained in the db
func (ks *VRFKeyStore) ListKeys() (publicKeys []*vrfkey.PublicKey, err error) {
	enc, err := ks.Get()
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
