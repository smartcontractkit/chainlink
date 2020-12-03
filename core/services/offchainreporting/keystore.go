package offchainreporting

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type KeyStore struct {
	*gorm.DB
	password     string
	p2pkeys      map[models.PeerID]p2pkey.Key
	ocrkeys      map[models.Sha256Hash]ocrkey.KeyBundle
	scryptParams utils.ScryptParams
	mu           *sync.RWMutex
}

func NewKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *KeyStore {
	return &KeyStore{
		DB:           db,
		p2pkeys:      make(map[models.PeerID]p2pkey.Key),
		ocrkeys:      make(map[models.Sha256Hash]ocrkey.KeyBundle),
		scryptParams: scryptParams,
		mu:           new(sync.RWMutex),
	}
}

func (ks *KeyStore) Unlock(password string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var errs error

	p2pkeys, err := ks.FindEncryptedP2PKeys()
	errs = multierr.Append(errs, err)
	ocrkeys, err := ks.FindEncryptedOCRKeyBundles()
	errs = multierr.Append(errs, err)

	for _, ek := range p2pkeys {
		k, err := ek.Decrypt(password)
		errs = multierr.Append(errs, err)
		peerID, err := k.GetPeerID()
		errs = multierr.Append(errs, err)
		ks.p2pkeys[models.PeerID(peerID)] = k
		logger.Debugw("Unlocked P2P key", "peerID", peerID)
	}
	for _, ek := range ocrkeys {
		k, err := ek.Decrypt(password)
		errs = multierr.Append(errs, err)
		if k != nil {
			ks.ocrkeys[k.ID] = *k
			logger.Debugw("Unlocked OCR key", "hash", k.ID)
		}
	}
	ks.password = password
	return errs
}

func (ks KeyStore) DecryptedP2PKey(peerID peer.ID) (p2pkey.Key, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	k, exists := ks.p2pkeys[models.PeerID(peerID)]
	return k, exists
}

func (ks KeyStore) DecryptedOCRKey(hash models.Sha256Hash) (ocrkey.KeyBundle, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	k, exists := ks.ocrkeys[hash]
	return k, exists
}

func (ks KeyStore) GenerateEncryptedP2PKey() (p2pkey.Key, p2pkey.EncryptedP2PKey, error) {
	key, err := p2pkey.CreateKey()
	if err != nil {
		return p2pkey.Key{}, p2pkey.EncryptedP2PKey{}, errors.Wrapf(err, "while generating new p2p key")
	}
	enc, err := key.ToEncryptedP2PKey(ks.password, ks.scryptParams)
	if err != nil {
		return p2pkey.Key{}, p2pkey.EncryptedP2PKey{}, errors.Wrapf(err, "while encrypting p2p key")
	}
	err = ks.UpsertEncryptedP2PKey(&enc)
	if err != nil {
		return p2pkey.Key{}, p2pkey.EncryptedP2PKey{}, err
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.p2pkeys[enc.PeerID] = key
	return key, enc, nil
}

func (ks KeyStore) UpsertEncryptedP2PKey(k *p2pkey.EncryptedP2PKey) error {
	err := ks.
		Set(
			"gorm:insert_option",
			"ON CONFLICT (pub_key) DO UPDATE SET "+
				"updated_at=NOW(),"+
				"deleted_at=null",
		).
		Create(k).
		Error
	if err != nil {
		return errors.Wrapf(err, "while inserting p2p key")
	}
	return nil
}

func (ks KeyStore) FindEncryptedP2PKeys() (keys []p2pkey.EncryptedP2PKey, err error) {
	return keys, ks.Order("created_at asc, id asc").Find(&keys).Error
}

func (ks KeyStore) FindEncryptedP2PKeyByID(id int32) (*p2pkey.EncryptedP2PKey, error) {
	var key p2pkey.EncryptedP2PKey
	err := ks.Where("id = ?", id).First(&key).Error
	return &key, err
}

func (ks KeyStore) ArchiveEncryptedP2PKey(key *p2pkey.EncryptedP2PKey) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	err := ks.Delete(key).Error
	if err != nil {
		return err
	}
	delete(ks.p2pkeys, key.PeerID)
	return nil
}

func (ks KeyStore) DeleteEncryptedP2PKey(key *p2pkey.EncryptedP2PKey) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	err := ks.Unscoped().Delete(key).Error
	if err != nil {
		return err
	}
	delete(ks.p2pkeys, key.PeerID)
	return nil
}

func (ks KeyStore) GenerateEncryptedOCRKeyBundle() (ocrkey.KeyBundle, ocrkey.EncryptedKeyBundle, error) {
	key, err := ocrkey.NewKeyBundle()
	if err != nil {
		return ocrkey.KeyBundle{}, ocrkey.EncryptedKeyBundle{}, errors.Wrapf(err, "while generating the new OCR key bundle")
	}
	enc, err := key.Encrypt(ks.password, ks.scryptParams)
	if err != nil {
		return ocrkey.KeyBundle{}, ocrkey.EncryptedKeyBundle{}, errors.Wrapf(err, "while encrypting the new OCR key bundle")
	}
	err = ks.CreateEncryptedOCRKeyBundle(enc)
	if err != nil {
		return ocrkey.KeyBundle{}, ocrkey.EncryptedKeyBundle{}, err
	}
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.ocrkeys[enc.ID] = *key
	return *key, *enc, nil
}

// CreateEncryptedOCRKeyBundle creates an encrypted OCR private key record
func (ks KeyStore) CreateEncryptedOCRKeyBundle(encryptedKey *ocrkey.EncryptedKeyBundle) error {
	err := ks.Create(encryptedKey).Error
	return errors.Wrapf(err, "while persisting the new encrypted OCR key bundle")
}

func (ks KeyStore) UpsertEncryptedOCRKeyBundle(encryptedKey *ocrkey.EncryptedKeyBundle) error {
	fmt.Println("encryptedKey.ID", encryptedKey.ID)
	err := ks.
		Set(
			"gorm:insert_option",
			"ON CONFLICT (id) DO UPDATE SET "+
				"updated_at=NOW(),"+
				"deleted_at=null",
		).
		Create(encryptedKey).
		Error
	if err != nil {
		return errors.Wrapf(err, "while upserting ocr key")
	}
	return nil
}

// FindEncryptedOCRKeyBundles finds all the encrypted OCR key records
func (ks KeyStore) FindEncryptedOCRKeyBundles() (keys []ocrkey.EncryptedKeyBundle, err error) {
	err = ks.Order("created_at asc, id asc").Find(&keys).Error
	return keys, err
}

// FindEncryptedOCRKeyBundleByID finds an EncryptedKeyBundle bundle by its ID
func (ks KeyStore) FindEncryptedOCRKeyBundleByID(id models.Sha256Hash) (ocrkey.EncryptedKeyBundle, error) {
	var key ocrkey.EncryptedKeyBundle
	err := ks.Where("id = ?", id).First(&key).Error
	return key, err
}

// ArchiveEncryptedOCRKeyBundle deletes the provided encrypted OCR key bundle
func (ks KeyStore) ArchiveEncryptedOCRKeyBundle(key *ocrkey.EncryptedKeyBundle) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	err := ks.Delete(key).Error
	if err != nil {
		return err
	}
	delete(ks.ocrkeys, key.ID)
	return nil
}

// DeleteEncryptedOCRKeyBundle deletes the provided encrypted OCR key bundle
func (ks KeyStore) DeleteEncryptedOCRKeyBundle(key *ocrkey.EncryptedKeyBundle) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	err := ks.Unscoped().Delete(key).Error
	if err != nil {
		return err
	}
	delete(ks.ocrkeys, key.ID)
	return nil
}

// ImportP2PKey imports a p2p key to the database
func (ks KeyStore) ImportP2PKey(keyJSON []byte, oldPassword string) (*p2pkey.EncryptedP2PKey, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var encryptedExport p2pkey.EncryptedP2PKeyExport
	err := json.Unmarshal(keyJSON, &encryptedExport)
	if err != nil {
		return nil, errors.Wrap(err, "invalid p2p key json")
	}
	privateKey, err := encryptedExport.DecryptPrivateKey(oldPassword)
	if err != nil {
		return nil, err
	}
	encryptedKey, err := privateKey.ToEncryptedP2PKey(ks.password, utils.DefaultScryptParams)
	if err != nil {
		return nil, err
	}
	err = ks.UpsertEncryptedP2PKey(&encryptedKey)
	if err != nil {
		return nil, err
	}
	ks.p2pkeys[encryptedKey.PeerID] = *privateKey

	return &encryptedKey, nil
}

// ExportP2PKey exports a p2p key from the database
func (ks KeyStore) ExportP2PKey(ID int32, newPassword string) ([]byte, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	emptyExport := []byte{}
	encryptedP2PKey, err := ks.FindEncryptedP2PKeyByID(ID)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to find p2p key with given ID")
	}
	decryptedP2PKey, err := encryptedP2PKey.Decrypt(ks.password)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to decrypt p2p key with given keystore password")
	}
	encryptedExport, err := decryptedP2PKey.ToEncryptedExport(newPassword, utils.DefaultScryptParams)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to encrypt p2p key for export with provided password")
	}

	return encryptedExport, nil
}

// ImportOCRKeyBundle imports an OCR key bundle to the database
func (ks KeyStore) ImportOCRKeyBundle(keyJSON []byte, oldPassword string) (*ocrkey.EncryptedKeyBundle, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	var encryptedExport ocrkey.EncryptedOCRKeyExport
	err := json.Unmarshal(keyJSON, &encryptedExport)
	if err != nil {
		return nil, errors.Wrap(err, "invalid OCR key json")
	}
	privateKey, err := encryptedExport.DecryptPrivateKey(oldPassword)
	if err != nil {
		return nil, err
	}
	encryptedKey, err := privateKey.Encrypt(ks.password, utils.DefaultScryptParams)
	if err != nil {
		return nil, err
	}
	err = ks.UpsertEncryptedOCRKeyBundle(encryptedKey)
	if err != nil {
		return nil, err
	}
	ks.ocrkeys[privateKey.ID] = *privateKey

	return encryptedKey, nil
}

// ExportOCRKeyBundle exports an OCR key bundle from the database
func (ks KeyStore) ExportOCRKeyBundle(id models.Sha256Hash, newPassword string) ([]byte, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	emptyExport := []byte{}
	encryptedP2PKey, err := ks.FindEncryptedOCRKeyBundleByID(id)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to find OCR key with given ID")
	}
	decryptedP2PKey, err := encryptedP2PKey.Decrypt(ks.password)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to decrypt p2p key with given keystore password")
	}
	encryptedExport, err := decryptedP2PKey.ToEncryptedExport(newPassword, utils.DefaultScryptParams)
	if err != nil {
		return emptyExport, errors.Wrap(err, "unable to encrypt p2p key for export with provided password")
	}

	return encryptedExport, nil
}
