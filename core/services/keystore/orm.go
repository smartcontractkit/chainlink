package keystore

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"gorm.io/gorm"
)

func NewORM(db *gorm.DB) ksORM {
	return ksORM{
		db: db,
	}
}

type ksORM struct {
	db *gorm.DB
}

func (orm ksORM) saveEncryptedKeyRing(kr *encryptedKeyRing) error {
	err := orm.db.Exec(`
		UPDATE encrypted_key_rings
		SET encrypted_keys = ?
	`, kr.EncryptedKeys).Error
	if err != nil {
		return errors.Wrap(err, "while saving keyring")
	}
	return nil
}

func (orm ksORM) getEncryptedKeyRing() (encryptedKeyRing, error) {
	kr := encryptedKeyRing{}
	err := orm.db.First(&kr).Error
	if err == gorm.ErrRecordNotFound {
		kr = encryptedKeyRing{}
		err2 := orm.db.Create(&kr).Error
		if err2 != nil {
			return kr, err2
		}
	} else if err != nil {
		return kr, err
	}
	return kr, nil
}

func (orm ksORM) loadKeyStates() (keyStates, error) {
	ks := newKeyStates()
	var ethkeystates []ethkey.State
	if err := orm.db.Find(&ethkeystates).Error; err != nil {
		return ks, errors.Wrap(err, "error loading eth_key_states from DB")
	}
	for i := 0; i < len(ethkeystates); i++ {
		ks.Eth[ethkeystates[i].KeyID()] = &ethkeystates[i]
	}
	return ks, nil
}

// ~~~~~~~~~~~~~~~~~~~~ LEGACY FUNCTIONS FOR V1 MIGRATION ~~~~~~~~~~~~~~~~~~~~

func (orm ksORM) GetEncryptedV1CSAKeys() (retrieved []csakey.Key, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}

func (orm ksORM) GetEncryptedV1EthKeys() (retrieved []ethkey.Key, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}

func (orm ksORM) GetEncryptedV1TerraKeys() (retrieved []terrakey.Key, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}

func (orm ksORM) GetEncryptedV1OCRKeys() (retrieved []ocrkey.EncryptedKeyBundle, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}

func (orm ksORM) GetEncryptedV1P2PKeys() (retrieved []p2pkey.EncryptedP2PKey, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}

func (orm ksORM) GetEncryptedV1VRFKeys() (retrieved []vrfkey.EncryptedVRFKey, err error) {
	return retrieved, orm.db.Find(&retrieved).Error
}
