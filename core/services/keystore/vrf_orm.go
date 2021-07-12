package keystore

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"gorm.io/gorm"
)

type VRFORM interface {
	FirstOrCreateEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error
	ArchiveEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error
	DeleteEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error
	FindEncryptedSecretVRFKeys(where ...vrfkey.EncryptedVRFKey) ([]*vrfkey.EncryptedVRFKey, error)
	FindEncryptedSecretVRFKeysIncludingArchived(where ...vrfkey.EncryptedVRFKey) ([]*vrfkey.EncryptedVRFKey, error)
}

type vrfORM struct {
	db *gorm.DB
}

var _ VRFORM = &vrfORM{}

func NewVRFORM(db *gorm.DB) VRFORM {
	return &vrfORM{
		db: db,
	}
}

// FirstOrCreateEncryptedVRFKey returns the first key found or creates a new one in the orm.
func (orm *vrfORM) FirstOrCreateEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error {
	return orm.db.FirstOrCreate(k).Error
}

// ArchiveEncryptedVRFKey soft-deletes k from the encrypted keys table, or errors
func (orm *vrfORM) ArchiveEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error {
	return orm.db.Delete(k).Error
}

// DeleteEncryptedVRFKey deletes k from the encrypted keys table, or errors
func (orm *vrfORM) DeleteEncryptedSecretVRFKey(k *vrfkey.EncryptedVRFKey) error {
	return orm.db.Unscoped().Delete(k).Error
}

// FindEncryptedVRFKeys retrieves matches to where from the encrypted keys table, or errors
func (orm *vrfORM) FindEncryptedSecretVRFKeys(where ...vrfkey.EncryptedVRFKey) (
	retrieved []*vrfkey.EncryptedVRFKey, err error) {
	var anonWhere []interface{} // Find needs "where" contents coerced to interface{}
	for _, constraint := range where {
		c := constraint
		anonWhere = append(anonWhere, &c)
	}
	return retrieved, orm.db.Find(&retrieved, anonWhere...).Order("created_at DESC").Error
}

func (orm *vrfORM) FindEncryptedSecretVRFKeysIncludingArchived(where ...vrfkey.EncryptedVRFKey) (
	retrieved []*vrfkey.EncryptedVRFKey, err error) {
	var anonWhere []interface{} // Find needs "where" contents coerced to interface{}
	for _, constraint := range where {
		c := constraint
		anonWhere = append(anonWhere, &c)
	}
	return retrieved, orm.db.Unscoped().Find(&retrieved, anonWhere...).Order("created_at DESC").Error
}
