package vrf

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

type ORM interface {
	FirstOrCreateEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	ArchiveEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	DeleteEncryptedSecretVRFKey(k *EncryptedVRFKey) error
	FindEncryptedSecretVRFKeys(where ...EncryptedVRFKey) ([]*EncryptedVRFKey, error)
	CreateEthTransaction(
		db *gorm.DB,
		fromAddress common.Address,
		toAddress common.Address,
		payload []byte,
		gasLimit uint64,
		maxUnconfirmedTransactions uint64) (*models.EthTx, error)
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

// CreateEthTransaction creates an ethereum transaction for the BPTXM to pick up
func (o *orm) CreateEthTransaction(
	db *gorm.DB,
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	maxUnconfirmedTransactions uint64,
) (*models.EthTx, error) {
	var etx models.EthTx
	err := utils.CheckOKToTransmit(postgres.MustSQLDB(db), fromAddress, maxUnconfirmedTransactions)
	if err != nil {
		return nil, errors.Wrap(err, "orm#CreateEthTransaction")
	}

	value := 0
	err = db.Raw(`
		INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at)
		SELECT $1,$2,$3,$4,$5,'unstarted',NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM eth_tx_attempts
			JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id
			WHERE eth_txes.from_address = $1
				AND eth_txes.state = 'unconfirmed'
				AND eth_tx_attempts.state = 'insufficient_eth'
		) RETURNING id;`,
		fromAddress,
		toAddress,
		payload,
		value,
		gasLimit,
	).Scan(&etx.ID).Error
	if err != nil {
		return nil, errors.Wrap(err, "keeper failed to insert eth_tx")
	}
	if etx.ID == 0 {
		return nil, errors.New("a keeper eth_tx with insufficient eth is present, not creating a new eth_tx")
	}
	err = db.First(&etx).Error
	if err != nil {
		return nil, errors.Wrap(err, "keeper find eth_tx after inserting")
	}
	return &etx, nil
}
