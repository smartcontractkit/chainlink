package keystore

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"gorm.io/gorm"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// ErrMatchingVRFKey is returned when Import attempts to import key with a
// PublicKey matching one already in the database
var ErrMatchingVRFKey = errors.New(
	`key with matching public key already stored in DB`)

// ErrAttemptToDeleteNonExistentKeyFromDB is returned when Delete is asked to
// delete a key it can't find in the DB.
var ErrAttemptToDeleteNonExistentKeyFromDB = errors.New("key is not present in DB")

// The VRF keystore tracks auxiliary VRF secret keys, and generates their VRF proofs
//
// VRF proofs need access to the actual secret key, which geth does not expose.
// Similar to the way geth's KeyStore exposes signing capability, VRF
// exposes VRF proof generation without the caller needing explicit knowledge of
// the secret key.
type VRF struct {
	lock         sync.RWMutex
	keys         InMemoryKeyStore
	orm          VRFORM
	scryptParams utils.ScryptParams
	// We store this upon first unlock to allow us
	// to create additional VRF keys via the remote CLI.
	password string
}

type InMemoryKeyStore = map[secp256k1.PublicKey]vrfkey.PrivateKey

// newVRFKeyStore returns an empty VRF KeyStore
func newVRFKeyStore(db *gorm.DB, sp utils.ScryptParams) *VRF {
	return &VRF{
		lock:         sync.RWMutex{},
		keys:         make(InMemoryKeyStore),
		orm:          NewVRFORM(db),
		scryptParams: sp,
	}
}

// GenerateProof is marshaled randomness proof given k and VRF input seed
// computed from the SeedData
//
// Key must have already been unlocked in ks, as constructing the VRF proof
// requires the secret key.
func (ks *VRF) GenerateProof(k secp256k1.PublicKey, seed *big.Int) (
	vrfkey.Proof, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	privateKey, found := ks.keys[k]
	if !found {
		return vrfkey.Proof{}, fmt.Errorf(
			"key %s has not been unlocked", k)
	}
	return privateKey.GenerateProof(seed)
}

// Unlock tries to unlock each vrf key in the db, using the given pass phrase,
// and returns any keys it manages to unlock, and any errors which result.
func (ks *VRF) Unlock(password string) (keysUnlocked []secp256k1.PublicKey,
	merr error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	keys, err := ks.get()
	if err != nil {
		return nil, errors.Wrap(err, "while retrieving vrf keys from db")
	}
	for _, k := range keys {
		key, err := vrfkey.Decrypt(k, password)
		if err != nil {
			merr = multierr.Append(merr, err)
			continue
		}
		ks.keys[key.PublicKey] = *key
		keysUnlocked = append(keysUnlocked, key.PublicKey)
	}
	ks.password = password
	return keysUnlocked, merr
}

// Forget removes the in-memory copy of the secret key of k, or errors if not
// present. Caller is responsible for taking ks.lock.
func (ks *VRF) forget(k secp256k1.PublicKey) error {
	if _, found := ks.keys[k]; !found {
		return fmt.Errorf("public key %s is not unlocked; can't forget it", k)
	}

	delete(ks.keys, k)
	return nil
}

func (ks *VRF) Forget(k secp256k1.PublicKey) error {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	return ks.forget(k)
}

// CreateKey returns a public key which is immediately unlocked in memory, and
// saved in DB encrypted with the password.
func (ks *VRF) CreateKey() (secp256k1.PublicKey, error) {
	if ks.password == "" {
		return secp256k1.PublicKey{}, errors.New("vrf keystore is not unlocked")
	}
	key := vrfkey.CreateKey()
	if err := ks.Store(key, ks.password, ks.scryptParams); err != nil {
		return secp256k1.PublicKey{}, err
	}
	return key.PublicKey, nil
}

// CreateAndUnlockWeakInMemoryEncryptedKeyXXXTestingOnly is for testing only! It returns
// an encrypted key which is fast to unlock, but correspondingly easy to brute
// force. It is not persisted to the DB, because no one should be keeping such
// keys lying around.
func (ks *VRF) CreateAndUnlockWeakInMemoryEncryptedKeyXXXTestingOnly(phrase string) (*vrfkey.EncryptedVRFKey, error) {
	key := vrfkey.CreateKey()
	encrypted, err := key.Encrypt(phrase, utils.FastScryptParams)
	if err != nil {
		return nil, errors.Wrap(err, "while creating testing key")
	}
	ks.keys[key.PublicKey] = *key
	return encrypted, nil
}

// Store saves key to ks (in memory), and to the DB, encrypted with phrase
func (ks *VRF) Store(key *vrfkey.PrivateKey, phrase string, scryptParams utils.ScryptParams) error {
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
func (ks *VRF) StoreInMemoryXXXTestingOnly(key *vrfkey.PrivateKey) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	ks.keys[key.PublicKey] = *key
}

var zeroPublicKey = secp256k1.PublicKey{}

// Archive soft-deletes keys with this public key from the keystore and the DB, if present.
func (ks *VRF) Archive(key secp256k1.PublicKey) (err error) {
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
		return ErrAttemptToDeleteNonExistentKeyFromDB
	}
	err2 := ks.orm.ArchiveEncryptedSecretVRFKey(&vrfkey.EncryptedVRFKey{PublicKey: key})
	return multierr.Append(err, err2)
}

// Delete removes keys with this public key from the keystore and the DB, if present.
func (ks *VRF) Delete(key secp256k1.PublicKey) (err error) {
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
		return ErrAttemptToDeleteNonExistentKeyFromDB
	}
	err2 := ks.orm.DeleteEncryptedSecretVRFKey(&vrfkey.EncryptedVRFKey{PublicKey: key})
	return multierr.Append(err, err2)
}

// Import adds this encrypted key to the DB and unlocks it in in-memory store
// with passphrase auth, and returns any resulting errors
func (ks *VRF) Import(keyjson []byte, auth string) (vrfkey.EncryptedVRFKey, error) {
	ks.lock.Lock()
	defer ks.lock.Unlock()
	enckey := vrfkey.EncryptedVRFKey{}
	if err := json.Unmarshal(keyjson, &enckey); err != nil {
		return enckey, fmt.Errorf("could not parse %s as vrfkey.EncryptedVRFKey json", keyjson)
	}
	extantMatchingKeys, err := ks.get(enckey.PublicKey)
	if err != nil {
		return enckey, errors.Wrapf(err, "while checking for matching extant key in DB")
	}
	if len(extantMatchingKeys) != 0 {
		return enckey, ErrMatchingVRFKey
	}
	key, err := vrfkey.Decrypt(&enckey, auth)
	if err != nil {
		return enckey, errors.Wrapf(err,
			"while attempting to decrypt key with public key %s",
			enckey.PublicKey.String())
	}
	if err := ks.orm.FirstOrCreateEncryptedSecretVRFKey(&enckey); err != nil {
		return enckey, errors.Wrapf(err, "while saving encrypted key to DB")
	}
	ks.keys[key.PublicKey] = *key
	return enckey, nil
}

func (ks *VRF) Export(pk secp256k1.PublicKey, newPassword string) ([]byte, error) {
	keys, err := ks.get(pk)
	if err != nil {
		return nil, err
	}
	privateKey, err := vrfkey.Decrypt(keys[0], ks.password)
	if err != nil {
		return nil, err
	}
	encKey, err := privateKey.Encrypt(newPassword, ks.scryptParams)
	if err != nil {
		return nil, err
	}
	return encKey.JSON()
}

// get retrieves all vrfkey.EncryptedVRFKey's associated with k, or all encrypted
// keys if k is nil, or errors. Caller is responsible for locking the store
func (ks *VRF) get(k ...secp256k1.PublicKey) ([]*vrfkey.EncryptedVRFKey,
	error) {
	if len(k) > 1 {
		return nil, errors.Errorf("can get at most one secret key at a time")
	}
	var where []vrfkey.EncryptedVRFKey
	if len(k) == 1 { // Search for this specific public key
		where = append(where, vrfkey.EncryptedVRFKey{PublicKey: k[0]})
	}
	keys, err := ks.orm.FindEncryptedSecretVRFKeys(where...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find public key %s in DB", k)
	}
	return keys, nil
}

// Get retrieves all vrfkey.EncryptedVRFKey's associated with k, or all encrypted
// keys if k is nil, or errors
func (ks *VRF) Get(k ...secp256k1.PublicKey) ([]*vrfkey.EncryptedVRFKey, error) {
	ks.lock.RLock()
	defer ks.lock.RUnlock()
	return ks.get(k...)
}

func (ks *VRF) GetSpecificKey(
	k secp256k1.PublicKey) (*vrfkey.EncryptedVRFKey, error) {
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
		// vrfkey.EncryptedVRFKey table.
		panic(fmt.Errorf("found more than one key with public key %s", k.String()))
	}
	return encryptedKey[0], nil
}

// ListKeys lists the public keys contained in the db
func (ks *VRF) ListKeys() (publicKeys []*secp256k1.PublicKey, err error) {
	enc, err := ks.Get()
	if err != nil {
		return nil, errors.Wrapf(err, "while listing db keys")
	}
	for _, enckey := range enc {
		publicKeys = append(publicKeys, &enckey.PublicKey)
	}
	return publicKeys, nil
}

// ListKeysIncludingArchived lists the public keys contained in the db
func (ks *VRF) ListKeysIncludingArchived() (publicKeys []*secp256k1.PublicKey, err error) {
	allKeys, err := ks.orm.FindEncryptedSecretVRFKeysIncludingArchived()
	if err != nil {
		return nil, err
	}
	for _, enckey := range allKeys {
		publicKeys = append(publicKeys, &enckey.PublicKey)
	}
	return publicKeys, nil
}
