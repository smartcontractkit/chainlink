package keystore

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	ErrLocked      = errors.New("Keystore is locked")
	ErrKeyNotFound = errors.New("Key not found")
	ErrKeyExists   = errors.New("Key already exists")
)

// DefaultEVMChainIDFunc is a func for getting a default evm chain ID -
// necessary because it is lazily evaluated
type DefaultEVMChainIDFunc func() (defaultEVMChainID *big.Int, err error)

//go:generate mockery --quiet --name Master --output ./mocks/ --case=underscore

type Master interface {
	CSA() CSA
	DKGSign() DKGSign
	DKGEncrypt() DKGEncrypt
	Eth() Eth
	OCR() OCR
	OCR2() OCR2
	P2P() P2P
	Solana() Solana
	Cosmos() Cosmos
	StarkNet() StarkNet
	VRF() VRF
	Unlock(password string) error
	IsEmpty() (bool, error)
}

type master struct {
	*keyManager
	cosmos     *cosmos
	csa        *csa
	eth        *eth
	ocr        *ocr
	ocr2       ocr2
	p2p        *p2p
	solana     *solana
	starknet   *starknet
	vrf        *vrf
	dkgSign    *dkgSign
	dkgEncrypt *dkgEncrypt
}

func New(db *sqlx.DB, scryptParams utils.ScryptParams, lggr logger.Logger, cfg pg.QConfig) Master {
	return newMaster(db, scryptParams, lggr, cfg)
}

func newMaster(db *sqlx.DB, scryptParams utils.ScryptParams, lggr logger.Logger, cfg pg.QConfig) *master {
	orm := NewORM(db, lggr, cfg)
	km := &keyManager{
		orm:          orm,
		keystateORM:  orm,
		scryptParams: scryptParams,
		lock:         &sync.RWMutex{},
		logger:       lggr.Named("KeyStore"),
	}

	return &master{
		keyManager: km,
		cosmos:     newCosmosKeyStore(km),
		csa:        newCSAKeyStore(km),
		eth:        newEthKeyStore(km, orm, orm.q),
		ocr:        newOCRKeyStore(km),
		ocr2:       newOCR2KeyStore(km),
		p2p:        newP2PKeyStore(km),
		solana:     newSolanaKeyStore(km),
		starknet:   newStarkNetKeyStore(km),
		vrf:        newVRFKeyStore(km),
		dkgSign:    newDKGSignKeyStore(km),
		dkgEncrypt: newDKGEncryptKeyStore(km),
	}
}

func (ks *master) DKGEncrypt() DKGEncrypt {
	return ks.dkgEncrypt
}

func (ks master) DKGSign() DKGSign {
	return ks.dkgSign
}

func (ks master) CSA() CSA {
	return ks.csa
}

func (ks *master) Eth() Eth {
	return ks.eth
}

func (ks *master) OCR() OCR {
	return ks.ocr
}

func (ks *master) OCR2() OCR2 {
	return ks.ocr2
}

func (ks *master) P2P() P2P {
	return ks.p2p
}

func (ks *master) Solana() Solana {
	return ks.solana
}

func (ks *master) Cosmos() Cosmos {
	return ks.cosmos
}

func (ks *master) StarkNet() StarkNet {
	return ks.starknet
}

func (ks *master) VRF() VRF {
	return ks.vrf
}

type ORM interface {
	isEmpty() (bool, error)
	saveEncryptedKeyRing(*encryptedKeyRing, ...func(pg.Queryer) error) error
	getEncryptedKeyRing() (encryptedKeyRing, error)
}

type keystateORM interface {
	loadKeyStates() (*keyStates, error)
}

type keyManager struct {
	orm          ORM
	keystateORM  keystateORM
	scryptParams utils.ScryptParams
	keyRing      *keyRing
	keyStates    *keyStates
	lock         *sync.RWMutex
	password     string
	logger       logger.Logger
}

func (km *keyManager) IsEmpty() (bool, error) {
	return km.orm.isEmpty()
}

func (km *keyManager) Unlock(password string) error {
	km.lock.Lock()
	defer km.lock.Unlock()
	// DEV: allow Unlock() to be idempotent - this is especially useful in tests,
	if km.password != "" {
		if password != km.password {
			return errors.New("attempting to unlock keystore again with a different password")
		}
		return nil
	}
	ekr, err := km.orm.getEncryptedKeyRing()
	if err != nil {
		return errors.Wrap(err, "unable to get encrypted key ring")
	}
	kr, err := ekr.Decrypt(password)
	if err != nil {
		return errors.Wrap(err, "unable to decrypt encrypted key ring")
	}
	kr.logPubKeys(km.logger)
	km.keyRing = kr

	ks, err := km.keystateORM.loadKeyStates()
	if err != nil {
		return errors.Wrap(err, "unable to load key states")
	}
	km.keyStates = ks

	km.password = password
	return nil
}

// caller must hold lock!
func (km *keyManager) save(callbacks ...func(pg.Queryer) error) error {
	ekb, err := km.keyRing.Encrypt(km.password, km.scryptParams)
	if err != nil {
		return errors.Wrap(err, "unable to encrypt keyRing")
	}
	return km.orm.saveEncryptedKeyRing(&ekb, callbacks...)
}

// caller must hold lock!
func (km *keyManager) safeAddKey(unknownKey Key, callbacks ...func(pg.Queryer) error) error {
	fieldName, err := GetFieldNameForKey(unknownKey)
	if err != nil {
		return err
	}
	// add key to keyring
	id := reflect.ValueOf(unknownKey.ID())
	key := reflect.ValueOf(unknownKey)
	keyRing := reflect.Indirect(reflect.ValueOf(km.keyRing))
	keyMap := keyRing.FieldByName(fieldName)
	keyMap.SetMapIndex(id, key)
	// save keyring to DB
	err = km.save(callbacks...)
	// if save fails, remove key from keyring
	if err != nil {
		keyMap.SetMapIndex(id, reflect.Value{})
		return err
	}
	return nil
}

// caller must hold lock!
func (km *keyManager) safeRemoveKey(unknownKey Key, callbacks ...func(pg.Queryer) error) (err error) {
	fieldName, err := GetFieldNameForKey(unknownKey)
	if err != nil {
		return err
	}
	id := reflect.ValueOf(unknownKey.ID())
	key := reflect.ValueOf(unknownKey)
	keyRing := reflect.Indirect(reflect.ValueOf(km.keyRing))
	keyMap := keyRing.FieldByName(fieldName)
	keyMap.SetMapIndex(id, reflect.Value{})
	// save keyring to DB
	err = km.save(callbacks...)
	// if save fails, add key back to keyRing
	if err != nil {
		keyMap.SetMapIndex(id, key)
		return err
	}
	return nil
}

// caller must hold lock!
func (km *keyManager) isLocked() bool {
	return len(km.password) == 0
}

func GetFieldNameForKey(unknownKey Key) (string, error) {
	switch unknownKey.(type) {
	case cosmoskey.Key:
		return "Cosmos", nil
	case csakey.KeyV2:
		return "CSA", nil
	case ethkey.KeyV2:
		return "Eth", nil
	case ocrkey.KeyV2:
		return "OCR", nil
	case ocr2key.KeyBundle:
		return "OCR2", nil
	case p2pkey.KeyV2:
		return "P2P", nil
	case solkey.Key:
		return "Solana", nil
	case starkkey.Key:
		return "StarkNet", nil
	case vrfkey.KeyV2:
		return "VRF", nil
	case dkgsignkey.Key:
		return "DKGSign", nil
	case dkgencryptkey.Key:
		return "DKGEncrypt", nil
	}
	return "", fmt.Errorf("unknown key type: %T", unknownKey)
}

type Key interface {
	ID() string
}
