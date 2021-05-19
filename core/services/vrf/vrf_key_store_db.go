package vrf

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// CreateKey returns a public key which is immediately unlocked in memory, and
// saved in DB encrypted with phrase. If p is given, its parameters are used for
// key derivation from the phrase.
func (ks *VRFKeyStore) CreateKey(phrase string) (secp256k1.PublicKey, error) {
	key := CreateKey()
	if err := ks.Store(key, phrase, ks.scryptParams); err != nil {
		return secp256k1.PublicKey{}, err
	}
	return key.PublicKey, nil
}

// CreateWeakInMemoryEncryptedKeyXXXTestingOnly is for testing only! It returns
// an encrypted key which is fast to unlock, but correspondingly easy to brute
// force. It is not persisted to the DB, because no one should be keeping such
// keys lying around.
func (ks *VRFKeyStore) CreateWeakInMemoryEncryptedKeyXXXTestingOnly(phrase string) (*EncryptedVRFKey, error) {
	key := CreateKey()
	encrypted, err := key.Encrypt(phrase, utils.FastScryptParams)
	if err != nil {
		return nil, errors.Wrap(err, "while creating testing key")
	}
	return encrypted, nil
}

// Store saves key to ks (in memory), and to the DB, encrypted with phrase
func (ks *VRFKeyStore) Store(key *PrivateKey, phrase string, scryptParams utils.ScryptParams) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	encrypted, err := key.Encrypt(phrase, scryptParams)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt key")
	}
	if err := ks.orm.FirstOrCreateEncryptedSecretVRFKey(encrypted); err != nil {
		return errors.Wrap(err, "failed to save encrypted key to db")
	}
	ks.keys[key.PublicKey] = *key
	return nil
}

// StoreInMemoryXXXTestingOnly memorizes key, only in in-memory store.
func (ks *VRFKeyStore) StoreInMemoryXXXTestingOnly(key *PrivateKey) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	ks.keys[key.PublicKey] = *key
}

var zeroPublicKey = secp256k1.PublicKey{}

// Archive soft-deletes keys with this public key from the keystore and the DB, if present.
func (ks *VRFKeyStore) Archive(key secp256k1.PublicKey) (err error) {
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
		return errors.Wrapf(err, "while checking for existence of key %s in DB", key.String())
	} else if len(matches) == 0 {
		return AttemptToDeleteNonExistentKeyFromDB
	}
	err2 := ks.orm.ArchiveEncryptedSecretVRFKey(&EncryptedVRFKey{PublicKey: key})
	return multierr.Append(err, err2)
}

// Delete removes keys with this public key from the keystore and the DB, if present.
func (ks *VRFKeyStore) Delete(key secp256k1.PublicKey) (err error) {
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
	err2 := ks.orm.DeleteEncryptedSecretVRFKey(&EncryptedVRFKey{PublicKey: key})
	return multierr.Append(err, err2)
}

// Import adds this encrypted key to the DB and unlocks it in in-memory store
// with passphrase auth, and returns any resulting errors
func (ks *VRFKeyStore) Import(keyjson []byte, auth string) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	enckey := &EncryptedVRFKey{}
	if err := json.Unmarshal(keyjson, enckey); err != nil {
		return fmt.Errorf("could not parse %s as EncryptedVRFKey json", keyjson)
	}
	extantMatchingKeys, err := ks.get(enckey.PublicKey)
	if err != nil {
		return errors.Wrapf(err, "while checking for matching extant key in DB")
	}
	if len(extantMatchingKeys) != 0 {
		return MatchingVRFKeyError
	}
	key, err := Decrypt(enckey, auth)
	if err != nil {
		return errors.Wrapf(err,
			"while attempting to decrypt key with public key %s",
			key.PublicKey.String())
	}
	if err := ks.orm.FirstOrCreateEncryptedSecretVRFKey(enckey); err != nil {
		return errors.Wrapf(err, "while saving encrypted key to DB")
	}
	ks.keys[key.PublicKey] = *key
	return nil
}

// get retrieves all EncryptedVRFKey's associated with k, or all encrypted
// keys if k is nil, or errors. Caller is responsible for locking the store
func (ks *VRFKeyStore) get(k ...secp256k1.PublicKey) ([]*EncryptedVRFKey,
	error) {
	if len(k) > 1 {
		return nil, errors.Errorf("can get at most one secret key at a time")
	}
	var where []EncryptedVRFKey
	if len(k) == 1 { // Search for this specific public key
		where = append(where, EncryptedVRFKey{PublicKey: k[0]})
	}
	keys, err := ks.orm.FindEncryptedSecretVRFKeys(where...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find public key %s in DB", k)
	}
	return keys, nil
}

// Get retrieves all EncryptedVRFKey's associated with k, or all encrypted
// keys if k is nil, or errors
func (ks *VRFKeyStore) Get(k ...secp256k1.PublicKey) ([]*EncryptedVRFKey, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	return ks.get(k...)
}

func (ks *VRFKeyStore) GetSpecificKey(
	k secp256k1.PublicKey) (*EncryptedVRFKey, error) {
	if k == (secp256k1.PublicKey{}) {
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
		// EncryptedVRFKey table.
		panic(fmt.Errorf("found more than one key with public key %s", k.String()))
	}
	return encryptedKey[0], nil
}

// ListKeys lists the public keys contained in the db
func (ks *VRFKeyStore) ListKeys() (publicKeys []*secp256k1.PublicKey, err error) {
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
