package types

import (
	"fmt"
	"math/big"
	"strconv"
)

type NumAsHex string

type AddInvokeTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
}

type AddDeclareResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

type AddDeployResponse struct {
	Code            string `json:"code"`
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"address"`
}

type DeployRequest struct {
	Type                string        `json:"type"`
	ContractAddressSalt string        `json:"contract_address_salt"`
	ConstructorCalldata []string      `json:"constructor_calldata"`
	ContractDefinition  ContractClass `json:"contract_definition"`
}

type DeployAccountRequest struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version uint64 `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`

	Type                string   `json:"type"`
	ContractAddressSalt string   `json:"contract_address_salt"`
	ConstructorCalldata []string `json:"constructor_calldata"`
	ClassHash           string   `json:"class_hash"`
}

// FunctionCall function call information
type FunctionCall struct {
	ContractAddress    Hash   `json:"contract_address"`
	EntryPointSelector string `json:"entry_point_selector,omitempty"`

	// Calldata The parameters passed to the function
	Calldata []string `json:"calldata"`
}

type Signature []*big.Int

type FunctionInvoke struct {
	MaxFee *big.Int `json:"max_fee"`
	// Version of the transaction scheme, should be set to 0 or 1
	Version *big.Int `json:"version"`
	// Signature
	Signature Signature `json:"signature"`
	// Nonce should only be set with Transaction V1
	Nonce *big.Int `json:"nonce,omitempty"`
	// Defines the transaction type to invoke
	Type string `json:"type,omitempty"`

	FunctionCall
}

type FeeEstimate struct {
	GasConsumed NumAsHex `json:"gas_consumed"`
	GasPrice    NumAsHex `json:"gas_price"`
	OverallFee  NumAsHex `json:"overall_fee"`
}

// ExecuteDetails provides some details about the execution.
type ExecuteDetails struct {
	MaxFee *big.Int
	Nonce  *big.Int
}

type TransactionState string

const (
	TransactionAcceptedOnL1 TransactionState = "ACCEPTED_ON_L1"
	TransactionAcceptedOnL2 TransactionState = "ACCEPTED_ON_L2"
	TransactionNotReceived  TransactionState = "NOT_RECEIVED"
	TransactionPending      TransactionState = "PENDING"
	TransactionReceived     TransactionState = "RECEIVED"
	TransactionRejected     TransactionState = "REJECTED"
)

func (ts *TransactionState) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	switch unquoted {
	case "ACCEPTED_ON_L2":
		*ts = TransactionAcceptedOnL2
	case "ACCEPTED_ON_L1":
		*ts = TransactionAcceptedOnL1
	case "NOT_RECEIVED":
		*ts = TransactionNotReceived
	case "PENDING":
		*ts = TransactionPending
	case "RECEIVED":
		*ts = TransactionReceived
	case "REJECTED":
		*ts = TransactionRejected
	default:
		return fmt.Errorf("unsupported status: %s", data)
	}
	return nil
}

func (ts TransactionState) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(ts))), nil
}

func (s TransactionState) String() string {
	return string(s)
}

func (s TransactionState) IsTransactionFinal() bool {
	if s == TransactionAcceptedOnL2 ||
		s == TransactionAcceptedOnL1 ||
		s == TransactionRejected {
		return true
	}
	return false
}
