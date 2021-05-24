package vrf

import (
	"gorm.io/gorm"
)

type ORM interface {
	FirstOrCreateEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	ArchiveEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	DeleteEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	FindEncryptedSecretVRFKeys(where ...EncryptedVRFKey) ([]*EncryptedVRFKey, error)
}

type orm struct {
	db *gorm.DB
}

var _ ORM = &orm{}

func NewORM(db *gorm.DB) ORM {
	return &orm{
		db: db,
	}
}

// FirstOrCreateEncryptedVRFKey returns the first key found or creates a new one in the orm.
func (orm *orm) FirstOrCreateEncryptedSecretVRFKey(k *EncryptedVRFKey) error {
	return orm.db.FirstOrCreate(k).Error
}

// ArchiveEncryptedVRFKey soft-deletes k from the encrypted keys table, or errors
func (orm *orm) ArchiveEncryptedSecretVRFKey(k *EncryptedVRFKey) error {
	return orm.db.Delete(k).Error
}

// DeleteEncryptedVRFKey deletes k from the encrypted keys table, or errors
func (orm *orm) DeleteEncryptedSecretVRFKey(k *EncryptedVRFKey) error {
	return orm.db.Unscoped().Delete(k).Error
}

// FindEncryptedVRFKeys retrieves matches to where from the encrypted keys table, or errors
func (orm *orm) FindEncryptedSecretVRFKeys(where ...EncryptedVRFKey) (
	retrieved []*EncryptedVRFKey, err error) {
	var anonWhere []interface{} // Find needs "where" contents coerced to interface{}
	for _, constraint := range where {
		c := constraint
		anonWhere = append(anonWhere, &c)
	}
	return retrieved, orm.db.Find(&retrieved, anonWhere...).Error
}
