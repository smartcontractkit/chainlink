package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SendTransactionRequest struct {
	From               common.Address `json:"from"`                 // address of sender/signer of meta-transaction
	Target             common.Address `json:"target"`               // ERC20 token contract address. Forwarder calls the target contract
	TargetName         string         `json:"target_name"`          // Name of the target contract
	Version            string         `json:"version"`              // Version of the target contract
	Nonce              *big.Int       `json:"nonce"`                // Forwarder nonce
	Receiver           common.Address `json:"receiver"`             // Recipient of the token-transfer. Must be an address in the destination chain
	Amount             *big.Int       `json:"amount"`               // Amount of ERC20 tokens to deliver
	SourceChainID      uint64         `json:"chain_id"`             // Source chain ID
	DestinationChainID uint64         `json:"destination_chain_id"` // Destination chain ID. Same as source chain ID for same-chain transfer
	ValidUntilTime     *big.Int       `json:"valid_until_time"`     // Meta-transaction expires if it is not written on-chain by validUntilTime (unix timestamp)
	Signature          []byte         `json:"signature"`            // EIP712 signature for the meta-transaction
}

type SendTransactionResponse struct {
	RequestID string `json:"request_id"` // UUID
}

type SendTransactionStatusRequest struct {
	RequestID     string       `json:"chnlnk_req_id"` // UUID
	Status        string       `json:"tx_status"`     // Status of request
	TxHash        *common.Hash `json:"tx_hash"`       // Transaction hash on source chain
	CCIPMessageID *common.Hash `json:"ccip_msg_id"`   // CCIP message ID
	FailureReason *string      `json:"tx_error"`      // failure reason of meta-transaction
}
