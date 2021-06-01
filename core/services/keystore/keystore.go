package keystore

import (
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

func NewKeyStore(db *gorm.DB, scryptParams utils.ScryptParams) *KeyStore {
	return &KeyStore{
		OCR: offchainreporting.NewKeyStore(db, scryptParams),
	}
}

type KeyStore struct {
	OCR *offchainreporting.KeyStore
}
