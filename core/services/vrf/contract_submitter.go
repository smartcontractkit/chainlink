package vrf

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

//go:generate mockery --name ContractSubmitter --output ./mocks/ --case=underscore

type ContractSubmitter interface {
	// TODO: Should be using the generic keystore for this once available
	CreateEthTransaction(
		db *gorm.DB,
		meta models.EthTxMetaV2,
		fromAddress common.Address,
		toAddress common.Address,
		payload []byte,
		gasLimit uint64,
		maxUnconfirmedTransactions uint64) (*models.EthTx, error)
}

type contractSubmitter struct {
}

func NewContractSubmitter() *contractSubmitter {
	return &contractSubmitter{}
}

// CreateEthTransaction creates an ethereum transaction for the BPTXM to pick up
func (*contractSubmitter) CreateEthTransaction(
	db *gorm.DB,
	meta models.EthTxMetaV2,
	fromAddress common.Address,
	toAddress common.Address,
	payload []byte,
	gasLimit uint64,
	maxUnconfirmedTransactions uint64,
) (*models.EthTx, error) {
	var etx models.EthTx
	err := bulletprooftxmanager.CheckEthTxQueueCapacity(db, fromAddress, maxUnconfirmedTransactions)
	if err != nil {
		return nil, errors.Wrap(err, "VRFListener: failed to check if ok to transmit")
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return nil, errors.Wrap(err, "VRFListener: failed to marshal ethtx metadata")
	}

	value := 0
	err = db.Raw(`
		INSERT INTO eth_txes (from_address, to_address, encoded_payload, value, gas_limit, state, created_at, meta)
		SELECT ?,?,?,?,?,'unstarted',NOW(),?
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
		b,
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
