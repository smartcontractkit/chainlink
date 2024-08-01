package ccip

import (
	"encoding/hex"
)

type Address string

type Addresses []Address

func (a Addresses) Strings() []string {
	s := make([]string, len(a))
	for i, addr := range a {
		s[i] = string(addr)
	}
	return s
}

func MakeAddresses(s []string) Addresses {
	addresses := make(Addresses, len(s))
	for i, a := range s {
		addresses[i] = Address(a)
	}
	return addresses
}

type Hash [32]byte

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

type TxMeta struct {
	BlockTimestampUnixMilli int64
	BlockNumber             uint64
	TxHash                  string
	LogIndex                uint64
	Finalized               FinalizedStatus
}

func (t *TxMeta) IsFinalized() bool {
	return t.Finalized == FinalizedStatusFinalized
}

// WithFinalityStatus accepts finalizedBlockNumber and based on that sets the Finalized status
// It's immutable so it creates new struct instead of in-place modification of the existing one
func (t *TxMeta) WithFinalityStatus(finalizedBlockNumber uint64) TxMeta {
	txMeta := TxMeta{
		BlockTimestampUnixMilli: t.BlockTimestampUnixMilli,
		BlockNumber:             t.BlockNumber,
		TxHash:                  t.TxHash,
		LogIndex:                t.LogIndex,
		Finalized:               FinalizedStatusNotFinalized,
	}
	if txMeta.BlockNumber <= finalizedBlockNumber {
		txMeta.Finalized = FinalizedStatusFinalized
	}
	return txMeta
}

type FinalizedStatus int

const (
	FinalizedStatusUnknown FinalizedStatus = iota
	FinalizedStatusFinalized
	FinalizedStatusNotFinalized
)
