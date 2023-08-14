package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Status uint8

const (
	Submitted       Status = iota // meta-transaction is submitted to tx manager
	Confirmed                     // meta-transaction has 1 block confirmation on the source chain
	SourceFinalized               // cross-chain meta-transaction is finalized on the source chain
	Finalized                     // same-chain meta-transaction is finalized on the source chain
	Failure                       // same-chain or cross-chain meta-transaction failed
)

func (s Status) String() string {
	switch s {
	case Submitted:
		return "submitted"
	case Confirmed:
		return "confirmed"
	case SourceFinalized:
		return "sourceFinalized"
	case Finalized:
		return "finalized"
	case Failure:
		return "failure"
	}
	return "unknown"
}

func (s Status) ToExternalStatus() (string, error) {
	switch s {
	case Submitted:
		return "SUBMITTED", nil
	case Confirmed:
		return "CONFIRMED", nil
	case SourceFinalized:
		return "SOURCE_FINALIZED", nil
	case Finalized:
		return "FINALIZED", nil
	case Failure:
		return "FAILURE", nil
	}
	return "", fmt.Errorf("unknown status: %s", s)
}

func (s *Status) Scan(value interface{}) error {
	status, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to set %v of %T to Status enum", value, value)
	}
	switch status {
	case "submitted":
		*s = Submitted
	case "confirmed":
		*s = Confirmed
	case "sourceFinalized":
		*s = SourceFinalized
	case "finalized":
		*s = Finalized
	case "failure":
		*s = Failure
	default:
		return fmt.Errorf("string to enum conversion not found for \"%s\"", status)
	}
	return nil
}

// LegacyGaslessTx encapsulates data related to a meta-transaction
type LegacyGaslessTx struct {
	ID                 string         `db:"legacy_gasless_tx_id"` // UUID
	Forwarder          common.Address `db:"forwarder_address"`    // forwarder contract
	From               common.Address `db:"from_address"`         // token sender
	Target             common.Address `db:"target_address"`       // token contract
	Receiver           common.Address `db:"receiver_address"`     // token receiver
	Nonce              *utils.Big     // forwarder nonce
	Amount             *utils.Big     // token amount to be transferred
	SourceChainID      uint64         // meta-transaction source chain ID. This is CCIP chain selector instead of EVM chain ID.
	DestinationChainID uint64         // meta-transaction destination chain ID. This is CCIP chain selector instead of EVM chain ID.
	ValidUntilTime     *utils.Big     // unix timestamp of meta-transaction expiry in seconds
	Signature          []byte         `db:"tx_signature"` // EIP712 signature
	Status             Status         `db:"tx_status"`    // status of meta-transaction
	FailureReason      *string        // failure reason of meta-transaction
	TokenName          string         // name of token used to generate EIP712 domain separator hash
	TokenVersion       string         // version of token used to generate EIP712 domain separator hash
	CCIPMessageID      *common.Hash   `db:"ccip_message_id"` // CCIP message ID
	EthTxID            string         `db:"eth_tx_id"`       // tx ID in transaction manager
	TxHash             *common.Hash   `db:"tx_hash"`         // transaction hash on source chain
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (gt *LegacyGaslessTx) Key() (*string, error) {
	if gt.Forwarder == utils.ZeroAddress {
		return nil, errors.New("empty forwarder address")
	}
	if gt.From == utils.ZeroAddress {
		return nil, errors.New("empty from address")
	}
	if gt.Nonce == nil {
		return nil, errors.New("nil nonce")
	}

	key := fmt.Sprintf("%s/%s/%s", gt.Forwarder, gt.From, gt.Nonce.String())
	return &key, nil
}

// LegacyGaslessTxPlus has additional fieds from eth_txes and eth_tx_attempts table
type LegacyGaslessTxPlus struct {
	LegacyGaslessTx
	EthTxStatus txmgrtypes.TxState `db:"etx_state"`
	EthTxHash   *common.Hash       `db:"etx_hash"`
	EthTxError  *string            `db:"etx_error"`
}

func (gt *LegacyGaslessTx) ToSendTransactionStatusRequest() (*SendTransactionStatusRequest, error) {
	status, err := gt.Status.ToExternalStatus()
	if err != nil {
		return nil, err
	}
	return &SendTransactionStatusRequest{
		RequestID:     gt.ID,
		Status:        status,
		TxHash:        gt.TxHash,
		CCIPMessageID: gt.CCIPMessageID,
		FailureReason: gt.FailureReason,
	}, nil
}
