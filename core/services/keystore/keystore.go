package keystore

import (
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

func NewKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *KeyStore {
	return &KeyStore{
		OCR: NewOCRKeyStore(db, scryptParams),
	}
}

type KeyStore struct {
	OCR *OCRKeyStore
}
