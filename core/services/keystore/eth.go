package keystore

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ErrKeyStoreLocked is returned if you call a method that requires unlocked keys before you unlocked the keystore
var ErrKeyStoreLocked = errors.New("keystore is locked (HINT: did you forget to call keystore.Unlock?)")

// EthKeyStoreInterface is the external interface for EthKeyStore
//go:generate mockery --name EthKeyStoreInterface --output mocks/ --case=underscore
type EthKeyStoreInterface interface {
	Unlock(password string) error

	// Requires Unlock
	CreateNewKey() (ethkey.Key, error)
	EnsureFundingKey() (key ethkey.Key, didExist bool, err error)
	ImportKey(keyJSON []byte, oldPassword string) (ethkey.Key, error)
	ExportKey(address common.Address, newPassword string) ([]byte, error)
	AddKey(key *ethkey.Key) error
	RemoveKey(address common.Address, hardDelete bool) (deletedKey ethkey.Key, err error)
	SubscribeToKeyChanges() (ch chan struct{}, unsub func())

	SignTx(fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	AllKeys() (keys []ethkey.Key, err error)
	SendingKeys() (keys []ethkey.Key, err error)
	FundingKeys() (keys []ethkey.Key, err error)
	KeyByAddress(address common.Address) (ethkey.Key, error)
	HasSendingKeyWithAddress(address common.Address) (bool, error)
	GetRoundRobinAddress(addresses ...common.Address) (address common.Address, err error)

	// Does not require Unlock
	HasDBSendingKeys() (bool, error)
	ImportKeyFileToDB(keyPath string) (ethkey.Key, error)
}

var _ EthKeyStoreInterface = &Eth{}

type combinedKey struct {
	DBKey        ethkey.Key
	DecryptedKey keystore.Key
	lastUsed     time.Time
}

// EthKeyStore manages an in-memory key list backed by a database table
// It never exposes private keys to consumers
type Eth struct {
	db           *gorm.DB
	password     string
	scryptParams utils.ScryptParams
	keys         []combinedKey
	mu           *sync.RWMutex

	subscribers   [](chan struct{})
	subscribersMu *sync.RWMutex
}

func newEthKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *Eth {
	return &Eth{db, "", scryptParams, make([]combinedKey, 0), new(sync.RWMutex), make([](chan struct{}), 0), new(sync.RWMutex)}
}

// Unlock loads keys from the database, and uses the given password to try to
// unlock all of them
// If any key fails to decrypt, returns an error
// Trying to unlock the keystore multiple times with different passwords will panic
func (ks *Eth) Unlock(password string) (merr error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	if ks.password != "" {
		if password == ks.password {
			return nil
		}
		return errors.New("may not unlock keystore more than once with different passwords")
	}
	var keys []ethkey.Key
	keys, merr = ks.loadDBKeys()
	if merr != nil {
		return errors.Wrap(merr, "EthKeyStore failed to load keys from database")
	}
	for _, k := range keys {
		dKey, err := keystore.DecryptKey(k.JSON, password)
		if err != nil {
			merr = multierr.Combine(merr, errors.Errorf("invalid password for account %s", k.Address.Hex()), err)
			continue

		}
		logger.Infow(fmt.Sprint("Unlocked account ", k.Address.Hex()), "address", k.Address.Hex(), "type", k.Type())
		cKey := combinedKey{k, *dKey, time.Time{}}
		ks.keys = append(ks.keys, cKey)
	}
	if merr != nil {
		return merr
	}
	ks.password = password
	return nil
}

func (ks *Eth) isLocked() bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	return ks.password == ""
}

// CreateNewKey adds an account to the underlying geth keystore (which
// writes the file to disk) and inserts the new key to the database
func (ks *Eth) CreateNewKey() (k ethkey.Key, err error) {
	if ks.isLocked() {
		return k, ErrKeyStoreLocked
	}
	return ks.createNewKey(false)
}

// EnsureFundingKey ensures that a funding account exists, and returns it
func (ks *Eth) EnsureFundingKey() (k ethkey.Key, didExist bool, err error) {
	if ks.isLocked() {
		return k, false, ErrKeyStoreLocked
	}
	existing, err := ks.getFundingKey()
	if err != nil {
		return k, false, err
	} else if existing != nil {
		return *existing, true, nil
	}
	k, err = ks.createNewKey(true)
	return k, false, err
}

func (ks *Eth) createNewKey(isFunding bool) (k ethkey.Key, err error) {
	dKey, err := newKey()
	if err != nil {
		return
	}
	exportedJSON, err := ks.encryptKey(&dKey, ks.password)
	if err != nil {
		return k, err
	}
	key := ethkey.Key{
		Address:   ethkey.EIP55AddressFromAddress(dKey.Address),
		IsFunding: isFunding,
		JSON:      datatypes.JSON(exportedJSON),
	}
	if err = ks.insertKeyIfNotExists(&key); err != nil {
		return k, err
	}
	cKey := combinedKey{key, dKey, time.Time{}}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys = append(ks.keys, cKey)
	ks.notify()
	return key, nil
}

func (ks *Eth) encryptKey(dKey *keystore.Key, newPassword string) ([]byte, error) {
	return keystore.EncryptKey(dKey, newPassword, ks.scryptParams.N, ks.scryptParams.P)
}

func (ks *Eth) getFundingKey() (*ethkey.Key, error) {
	fundingKeys, err := ks.FundingKeys()
	if err != nil {
		return nil, err
	}
	if len(fundingKeys) > 0 {
		return &fundingKeys[0], nil
	}
	return nil, nil
}

// SignTx uses the unlocked account to sign the given transaction.
func (ks *Eth) SignTx(fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if ks.isLocked() {
		return nil, ErrKeyStoreLocked
	}
	signer := types.LatestSignerForChainID(chainID)

	dKey := ks.getDecryptedKeyForAddress(fromAddress)
	if dKey == nil {
		return nil, newNoKeyError(fromAddress)
	}

	return types.SignTx(tx, signer, dKey.PrivateKey)
}

func (ks *Eth) getDecryptedKeyForAddress(addr common.Address) *keystore.Key {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if cKey.DecryptedKey.Address == addr {
			return &cKey.DecryptedKey
		}
	}
	return nil
}

// HasSendingKeyWithAddress returns true if keystore has an account with the given address
func (ks *Eth) HasSendingKeyWithAddress(address common.Address) (bool, error) {
	if ks.isLocked() {
		return false, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if !cKey.DBKey.IsFunding && cKey.DecryptedKey.Address == address {
			return true, nil
		}
	}
	return false, nil
}

// GetKeyByAddress returns the account matching the address provided, or an error if it is missing
func (ks *Eth) GetKeyByAddress(address common.Address) (ethkey.Key, error) {
	if ks.isLocked() {
		return ethkey.Key{}, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if cKey.DBKey.Address.Address() == address {
			return cKey.DBKey, nil
		}
	}
	return ethkey.Key{}, newNoKeyError(address)
}

// ImportKey adds a new key to the keystore and inserts to DB
func (ks *Eth) ImportKey(keyJSON []byte, oldPassword string) (key ethkey.Key, err error) {
	if ks.isLocked() {
		return key, ErrKeyStoreLocked
	}
	dKey, err := keystore.DecryptKey(keyJSON, oldPassword)
	if err != nil {
		return key, errors.Wrap(err, "EthKeyStore#ImportKey failed to decrypt key")
	}
	exportedJSON, err := ks.encryptKey(dKey, ks.password)
	if err != nil {
		return key, err
	}
	key = ethkey.Key{
		Address:   ethkey.EIP55AddressFromAddress(dKey.Address),
		IsFunding: false,
		JSON:      datatypes.JSON(exportedJSON),
	}
	if err := ks.insertKeyIfNotExists(&key); err != nil {
		return key, err
	}
	cKey := combinedKey{key, *dKey, time.Time{}}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys = append(ks.keys, cKey)
	ks.notify()
	return key, nil
}

// ExportKey exports as a JSON key, encrypted with newPassword
func (ks *Eth) ExportKey(address common.Address, newPassword string) ([]byte, error) {
	if ks.isLocked() {
		return nil, ErrKeyStoreLocked
	}
	var dKey keystore.Key
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, k := range ks.keys {
		if k.DecryptedKey.Address == address {
			dKey = k.DecryptedKey
		}
	}
	if dKey.Address == utils.ZeroAddress {
		return nil, newNoKeyError(address)
	}
	return ks.encryptKey(&dKey, newPassword)
}

// AddKey inserts the key to the database and adds it to the keystore's memory keys
// It modifies the given key (adding created_at etc)
func (ks *Eth) AddKey(key *ethkey.Key) error {
	if ks.isLocked() {
		return ErrKeyStoreLocked
	}
	dKey, err := keystore.DecryptKey(key.JSON, ks.password)
	if err != nil {
		return errors.Wrap(err, "unable to decrypt key JSON with keystore password")
	}
	if err := ks.insertKeyIfNotExists(key); err != nil {
		return errors.Wrap(err, "unable to insert key")
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	cKey := combinedKey{DBKey: *key, DecryptedKey: *dKey}
	ks.keys = append(ks.keys, cKey)
	ks.notify()
	return nil
}

// RemoveKey removes a key from the keystore
// If hard delete is set to true, removes the key from the database. If false, the key has its deleted_at set to a non-null value.
func (ks *Eth) RemoveKey(address common.Address, hardDelete bool) (removedKey ethkey.Key, err error) {
	if ks.isLocked() {
		return removedKey, ErrKeyStoreLocked
	}

	ks.mu.Lock()
	for i, cKey := range ks.keys {
		if cKey.DecryptedKey.Address == address {
			removedKey = cKey.DBKey
			ks.keys = append(ks.keys[:i], ks.keys[i+1:]...)
			ks.notify()
		}
	}
	ks.mu.Unlock()

	if removedKey.Address.Address() == utils.ZeroAddress {
		return removedKey, newNoKeyError(address)
	}

	var sql string
	if hardDelete {
		sql = `DELETE FROM keys WHERE address = ?`
	} else {
		sql = `UPDATE keys SET deleted_at = NOW() WHERE address = ?`
	}
	err = postgres.DBWithDefaultContext(ks.db, func(db *gorm.DB) error {
		return db.Exec(sql, address).Error
	})

	return
}

// SubscribeToKeyChanges returns a channel that will fire if a key is added or removed
// Consumers should call unsubscribe when they are done to close the channel
func (ks *Eth) SubscribeToKeyChanges() (ch chan struct{}, unsubscribe func()) {
	ch = make(chan struct{}, 1)
	ks.subscribersMu.Lock()
	defer ks.subscribersMu.Unlock()
	ks.subscribers = append(ks.subscribers, ch)
	return ch, func() {
		ks.subscribersMu.Lock()
		defer ks.subscribersMu.Unlock()
		for i, sub := range ks.subscribers {
			if sub == ch {
				ks.subscribers = append(ks.subscribers[:i], ks.subscribers[i+1:]...)
				close(ch)
			}
		}
	}
}

func (ks *Eth) notify() {
	ks.subscribersMu.RLock()
	defer ks.subscribersMu.RUnlock()
	for _, ch := range ks.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// AllKeys returns all keys
func (ks *Eth) AllKeys() (keys []ethkey.Key, err error) {
	if ks.isLocked() {
		return nil, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	keys = make([]ethkey.Key, len(ks.keys))
	for i, cKey := range ks.keys {
		keys[i] = cKey.DBKey
	}
	return keys, nil
}

// SendingKeys will return only the keys that are is_funding=false
func (ks *Eth) SendingKeys() (keys []ethkey.Key, err error) {
	if ks.isLocked() {
		return nil, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if !cKey.DBKey.IsFunding {
			keys = append(keys, cKey.DBKey)
		}
	}
	return keys, nil
}

// FundingKeys will return only the keys that are is_funding=true
func (ks *Eth) FundingKeys() (keys []ethkey.Key, err error) {
	if ks.isLocked() {
		return nil, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if cKey.DBKey.IsFunding {
			keys = append(keys, cKey.DBKey)
		}
	}
	return keys, nil
}

// KeyByAddress returns the key matching provided address
func (ks *Eth) KeyByAddress(address common.Address) (ethkey.Key, error) {
	if ks.isLocked() {
		return ethkey.Key{}, ErrKeyStoreLocked
	}
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	for _, cKey := range ks.keys {
		if cKey.DecryptedKey.Address == address {
			return cKey.DBKey, nil
		}
	}
	return ethkey.Key{}, newNoKeyError(address)
}

// GetRoundRobinAddress gets the address of the "next" available sending key (i.e. the least recently used key)
// This takes an optional param for a slice of addresses it should pick from. Leave empty to pick from all
// addresses in the keystore.
func (ks *Eth) GetRoundRobinAddress(whitelist ...common.Address) (address common.Address, err error) {
	if ks.isLocked() {
		return common.Address{}, ErrKeyStoreLocked
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()

	var keys []combinedKey
	for _, cKey := range ks.keys {
		if !cKey.DBKey.IsFunding {
			if len(whitelist) == 0 {
				keys = append(keys, cKey)
			} else {
				for _, addr := range whitelist {
					if addr == cKey.DecryptedKey.Address {
						keys = append(keys, cKey)
					}
				}
			}
		}
	}

	if len(keys) == 0 {
		return common.Address{}, errors.New("no keys available")
	}

	var leastRecentlyUsed combinedKey

	for i, cKey := range keys {
		if i == 0 {
			leastRecentlyUsed = cKey
		} else if cKey.lastUsed.Before(leastRecentlyUsed.lastUsed) {
			leastRecentlyUsed = cKey
		}
	}

	for i, cKey := range ks.keys {
		if cKey.DecryptedKey.Address == leastRecentlyUsed.DecryptedKey.Address {
			ks.keys[i].lastUsed = time.Now()
		}
	}

	return leastRecentlyUsed.DecryptedKey.Address, nil
}

// HasDBSendingKeys returns true if any key in the database is a sending key
func (ks *Eth) HasDBSendingKeys() (exists bool, err error) {
	err = postgres.DBWithDefaultContext(ks.db, func(db *gorm.DB) error {
		return db.Raw(`SELECT EXISTS(SELECT 1 FROM keys WHERE is_funding=false)`).Scan(&exists).Error
	})
	return
}

// ImportKeyFileToDB reads a file and writes the key to the database
func (ks *Eth) ImportKeyFileToDB(keyPath string) (k ethkey.Key, err error) {
	k, err = ethkey.NewKeyFromFile(keyPath)
	if err != nil {
		return k, errors.Wrap(err, "could not import key from file")
	}
	err = ks.insertKeyIfNotExists(&k)
	return
}

// loadDBKeys returns a map of all of the keys saved in the database
// including the funding key.
func (ks *Eth) loadDBKeys() (keys []ethkey.Key, err error) {
	err = postgres.GormTransactionWithDefaultContext(ks.db, func(db *gorm.DB) error {
		return db.Order("created_at ASC, address ASC").Where("deleted_at IS NULL").Find(&keys).Error
	})
	return
}

// insertKeyIfNotExists inserts a key if a key with that address doesn't exist already
// If a key with this address exists, it does nothing
func (ks *Eth) insertKeyIfNotExists(k *ethkey.Key) error {
	err := ks.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil}),
	}).Create(k).Error
	if err == nil || err.Error() == "sql: no rows in result set" {
		return nil
	}
	return errors.Wrap(err, "insertKeyIfNotExists failed")
}

// newKey pulled from geth (sadly not exported)
func newKey() (dKey keystore.Key, err error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), crand.Reader)
	if err != nil {
		return dKey, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return dKey, errors.Errorf("Could not create random uuid: %v", err)
	}
	dKey = keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return dKey, nil
}

func newNoKeyError(address common.Address) error {
	return errors.Errorf("address %s not in keystore", address.Hex())
}
