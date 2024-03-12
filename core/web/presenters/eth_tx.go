package presenters

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// EthTxResource represents a Ethereum Transaction JSONAPI resource.
type EthTxResource struct {
	JAID
	State      string          `json:"state"`
	Data       hexutil.Bytes   `json:"data"`
	From       *common.Address `json:"from"`
	GasLimit   string          `json:"gasLimit"`
	GasPrice   string          `json:"gasPrice"`
	Hash       common.Hash     `json:"hash"`
	Hex        string          `json:"rawHex"`
	Nonce      string          `json:"nonce"`
	SentAt     string          `json:"sentAt"`
	To         *common.Address `json:"to"`
	Value      string          `json:"value"`
	EVMChainID big.Big         `json:"evmChainID"`
}

// GetName implements the api2go EntityNamer interface
func (EthTxResource) GetName() string {
	return "evm_transactions"
}

// NewEthTxResource generates a EthTxResource from an Eth.Tx.
//
// For backwards compatibility, there is no id set when initializing from an
// EthTx as the id being used was the EthTxAttempt Hash.
// This should really use it's proper id
func NewEthTxResource(tx txmgr.Tx) EthTxResource {
	v := assets.Eth(tx.Value)
	r := EthTxResource{
		Data:     hexutil.Bytes(tx.EncodedPayload),
		From:     &tx.FromAddress,
		GasLimit: strconv.FormatUint(tx.FeeLimit, 10),
		State:    string(tx.State),
		To:       &tx.ToAddress,
		Value:    v.String(),
	}

	if tx.ChainID != nil {
		r.EVMChainID = *big.New(tx.ChainID)
	}
	return r
}

func NewEthTxResourceFromAttempt(txa txmgr.TxAttempt) EthTxResource {
	tx := txa.Tx

	r := NewEthTxResource(tx)
	r.JAID = NewJAID(txa.Hash.String())
	r.GasPrice = txa.TxFee.Legacy.ToInt().String()
	r.Hash = txa.Hash
	r.Hex = hexutil.Encode(txa.SignedRawTx)

	if txa.Tx.ChainID != nil {
		r.EVMChainID = *big.New(txa.Tx.ChainID)
		r.JAID = NewPrefixedJAID(r.JAID.ID, txa.Tx.ChainID.String())
	}

	if tx.Sequence != nil {
		r.Nonce = strconv.FormatUint(uint64(*tx.Sequence), 10)
	}
	if txa.BroadcastBeforeBlockNum != nil {
		r.SentAt = strconv.FormatUint(uint64(*txa.BroadcastBeforeBlockNum), 10)
	}
	return r
}
