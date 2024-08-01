package rpc

import "github.com/NethermindEth/juno/core/felt"

// AddDeclareTransactionResponse provides the output for AddDeclareTransaction.
type AddDeclareTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ClassHash       *felt.Felt `json:"class_hash"`
}

// AddDeployTransactionResponse provides the output for AddDeployTransaction.
type AddDeployAccountTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
	ContractAddress *felt.Felt `json:"contract_address"`
}

// AddInvokeTransactionResponse provides the output for AddInvokeTransaction.
type AddInvokeTransactionResponse struct {
	TransactionHash *felt.Felt `json:"transaction_hash"`
}
