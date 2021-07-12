package presenters

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
)

// EthTxResource represents a Ethereum Transaction JSONAPI resource.
type EthTxResource struct {
	JAID
	State    string          `json:"state"`
	Data     hexutil.Bytes   `json:"data"`
	From     *common.Address `json:"from"`
	GasLimit string          `json:"gasLimit"`
	GasPrice string          `json:"gasPrice"`
	Hash     common.Hash     `json:"hash"`
	Hex      string          `json:"rawHex"`
	Nonce    string          `json:"nonce"`
	SentAt   string          `json:"sentAt"`
	To       *common.Address `json:"to"`
	Value    string          `json:"value"`
}

// GetName implements the api2go EntityNamer interface
func (EthTxResource) GetName() string {
	return "transactions"
}

// NewEthTxResource generates a EthTxResource from an Eth.Tx.
//
// For backwards compatibility, there is no id set when initializing from an
// EthTx as the id being used was the EthTxAttempt Hash.
// This should really use it's proper id
func NewEthTxResource(tx bulletprooftxmanager.EthTx) EthTxResource {
	return EthTxResource{
		Data:     hexutil.Bytes(tx.EncodedPayload),
		From:     &tx.FromAddress,
		GasLimit: strconv.FormatUint(tx.GasLimit, 10),
		State:    string(tx.State),
		To:       &tx.ToAddress,
		Value:    tx.Value.String(),
	}
}

func NewEthTxResourceFromAttempt(txa bulletprooftxmanager.EthTxAttempt) EthTxResource {
	tx := txa.EthTx

	r := NewEthTxResource(tx)
	r.JAID = NewJAID(txa.Hash.Hex())
	r.GasPrice = txa.GasPrice.String()
	r.Hash = txa.Hash
	r.Hex = hexutil.Encode(txa.SignedRawTx)

	if tx.Nonce != nil {
		r.Nonce = strconv.FormatUint(uint64(*tx.Nonce), 10)
	}
	if txa.BroadcastBeforeBlockNum != nil {
		r.SentAt = strconv.FormatUint(uint64(*txa.BroadcastBeforeBlockNum), 10)
	}
	return r
}
