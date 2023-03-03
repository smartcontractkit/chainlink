package rpcv01

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dontpanicdao/caigo/types"
)

type CommonTransactionReceipt struct {
	TransactionHash types.Hash `json:"transaction_hash"`
	// ActualFee The fee that was charged by the sequencer
	ActualFee   string                 `json:"actual_fee"`
	Status      types.TransactionState `json:"status"`
	BlockHash   types.Hash             `json:"block_hash"`
	BlockNumber uint64                 `json:"block_number"`
	Type        TransactionType        `json:"type,omitempty"`
}

func (tr CommonTransactionReceipt) Hash() types.Hash {
	return tr.TransactionHash
}

type TransactionType string

const (
	TransactionType_Declare       TransactionType = "DECLARE"
	TransactionType_Deploy        TransactionType = "DEPLOY"
	TransactionType_DeployAccount TransactionType = "DEPLOY_ACCOUNT"
	TransactionType_Invoke        TransactionType = "INVOKE"
	TransactionType_L1Handler     TransactionType = "L1_HANDLER"
)

func (tt *TransactionType) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "DECLARE":
		*tt = TransactionType_Declare
	case "DEPLOY":
		*tt = TransactionType_Deploy
	case "DEPLOY_ACCOUNT":
		*tt = TransactionType_DeployAccount
	case "INVOKE":
		*tt = TransactionType_Invoke
	case "L1_HANDLER":
		*tt = TransactionType_L1Handler
	default:
		return fmt.Errorf("unsupported type: %s", data)
	}

	return nil
}

func (tt TransactionType) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(tt))), nil
}

type InvokeTransactionReceiptProperties struct {
	MessageSent []MsgToL1 `json:"messages_sent,omitempty"`
	// A list of events assocuated with the Invoke Transaction
	Events []Event `json:"events,omitempty"`
}

// InvokeTransactionReceipt Invoke Transaction Receipt
type InvokeTransactionReceipt struct {
	CommonTransactionReceipt
	// ActualFee The fee that was charged by the sequencer
	InvokeTransactionReceiptProperties `json:",omitempty"`
}

// DeclareTransactionReceipt Declare Transaction Receipt
type DeclareTransactionReceipt struct {
	CommonTransactionReceipt
}

// DeployTransactionReceipt Deploy Transaction Receipt
type DeployTransactionReceipt struct {
	CommonTransactionReceipt
	// ContractAddress The address of the deployed contract
	ContractAddress string `json:"contract_address"`
}

// L1HandlerTransactionReceipt L1 Handler Transaction Receipt
type L1HandlerTransactionReceipt struct {
	CommonTransactionReceipt
}

type TransactionReceipt interface {
	Hash() types.Hash
}

type MsgToL1 struct {
	// ToAddress The target L1 address the message is sent to
	ToAddress string `json:"to_address"`
	//Payload  The payload of the message
	Payload []string `json:"payload"`
}

type UnknownTransactionReceipt struct{ TransactionReceipt }

func (tr *UnknownTransactionReceipt) UnmarshalJSON(data []byte) error {
	var dec map[string]interface{}
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}

	t, err := unmarshalTransactionReceipt(dec)
	if err != nil {
		return err
	}
	*tr = UnknownTransactionReceipt{t}
	return nil
}

func unmarshalTransactionReceipt(t interface{}) (TransactionReceipt, error) {
	switch casted := t.(type) {
	case string:
		return TransactionHash{types.HexToHash(casted)}, nil
	case map[string]interface{}:
		// NOTE(tvanas): Pathfinder 0.3.3 does not return
		// transaction receipt types. We handle this by
		// naively marshalling into an invoke type. Once it
		// is supported, this condition can be removed.
		typ, ok := casted["type"]
		if !ok {
			var txn InvokeTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		}

		switch TransactionType(typ.(string)) {
		case TransactionType_Declare:
			var txn DeclareTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Deploy:
			var txn DeployTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_DeployAccount:
			var txn DeployTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_Invoke:
			var txn InvokeTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		case TransactionType_L1Handler:
			var txn L1HandlerTransactionReceipt
			remarshal(casted, &txn)
			return txn, nil
		}
	}

	return nil, fmt.Errorf("unknown transaction type: %v", t)
}
