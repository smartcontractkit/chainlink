package types

import (
	"encoding/json"
)

type TransactionHash Felt

func BytesToTransactionHash(b []byte) TransactionHash {
	return TransactionHash(BytesToFelt(b))
}

func HexToTransactionHash(s string) TransactionHash {
	return TransactionHash(HexToFelt(s))
}

func (t TransactionHash) Felt() Felt {
	return Felt(t)
}

func (t TransactionHash) Bytes() []byte {
	return t.Felt().Bytes()
}

func (t TransactionHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(Felt(t))
}

func (t TransactionHash) String() string {
	return Felt(t).String()
}

func (t *TransactionHash) UnmarshalJSON(data []byte) error {
	f := Felt{}
	err := f.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	*t = TransactionHash(f)
	return nil
}

type IsTransaction interface {
	GetHash() TransactionHash
}

type TransactionDeploy struct {
	Hash                TransactionHash
	ContractAddress     Address
	ConstructorCallData []Felt
}

func (tx *TransactionDeploy) GetHash() TransactionHash {
	return tx.Hash
}

type TransactionInvoke struct {
	Hash               TransactionHash `json:"txn_hash"`
	ContractAddress    Address         `json:"contract_address"`
	EntryPointSelector Felt            `json:"entry_point_selector"`
	CallData           []Felt          `json:"calldata"`
	Signature          []Felt          `json:"-"`
	MaxFee             Felt            `json:"max_fee"`
}

func (tx *TransactionInvoke) GetHash() TransactionHash {
	return tx.Hash
}

type TransactionStatus int64

const (
	TxStatusUnknown TransactionStatus = iota
	TxStatusNotReceived
	TxStatusReceived
	TxStatusPending
	TxStatusRejected
	TxStatusAcceptedOnL2
	TxStatusAcceptedOnL1
)

var (
	TxStatusName = map[TransactionStatus]string{
		TxStatusUnknown:      "UNKNOWN",
		TxStatusNotReceived:  "NOT_RECEIVED",
		TxStatusReceived:     "RECEIVED",
		TxStatusPending:      "PENDING",
		TxStatusRejected:     "REJECTED",
		TxStatusAcceptedOnL2: "ACCEPTED_ON_L2",
		TxStatusAcceptedOnL1: "ACCEPTED_ON_L1",
	}
	TxStatusValue = map[string]TransactionStatus{
		"UNKNOWN":        TxStatusUnknown,
		"NOT_RECEIVED":   TxStatusNotReceived,
		"RECEIVED":       TxStatusReceived,
		"PENDING":        TxStatusPending,
		"REJECTED":       TxStatusRejected,
		"ACCEPTED_ON_L2": TxStatusAcceptedOnL2,
		"ACCEPTED_ON_L1": TxStatusAcceptedOnL1,
	}
)

func (s TransactionStatus) String() string {
	// notest
	return TxStatusName[s]
}

type TransactionReceipt struct {
	TxHash          TransactionHash
	ActualFee       Felt
	Status          TransactionStatus
	StatusData      string
	MessagesSent    []MessageL2ToL1
	L1OriginMessage *MessageL1ToL2
	Events          []Event
}

type TransactionWithReceipt struct {
	IsTransaction
	TransactionReceipt
}
