package types

import (
	"encoding/json"
)

type BlockStatus int32

const (
	BlockStatusUnknown BlockStatus = iota
	BlockStatusPending
	BlockStatusProven
	BlockStatusAcceptedOnL2
	BlockStatusAcceptedOnL1
	BlockStatusRejected
)

var (
	BlockStatusName = map[BlockStatus]string{
		BlockStatusUnknown:      "UNKNOWN",
		BlockStatusPending:      "PENDING",
		BlockStatusProven:       "PROVEN",
		BlockStatusAcceptedOnL2: "ACCEPTED_ON_L2",
		BlockStatusAcceptedOnL1: "ACCEPTED_ON_L1",
		BlockStatusRejected:     "REJECTED",
	}
	BlockStatusValue = map[string]BlockStatus{
		"UNKNOWN":        BlockStatusUnknown,
		"PENDING":        BlockStatusPending,
		"PROVEN":         BlockStatusProven,
		"ACCEPTED_ON_L2": BlockStatusAcceptedOnL2,
		"ACCEPTED_ON_L1": BlockStatusAcceptedOnL1,
		"REJECTED":       BlockStatusRejected,
	}
)

func StringToBlockStatus(s string) BlockStatus {
	blockStatus, ok := BlockStatusValue[s]
	if !ok {
		// notest
		return BlockStatusUnknown
	}
	return blockStatus
}

func (b BlockStatus) String() string {
	return BlockStatusName[b]
}

func (b BlockStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

type BlockHash Felt

func BytesToBlockHash(b []byte) BlockHash {
	return BlockHash(BytesToFelt(b))
}

func HexToBlockHash(s string) BlockHash {
	return BlockHash(HexToFelt(s))
}

func (b BlockHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Felt())
}

func (b *BlockHash) Bytes() []byte {
	return b.Felt().Bytes()
}

func (b *BlockHash) Felt() Felt {
	return Felt(*b)
}

func (b BlockHash) Hex() string {
	return Felt(b).Hex()
}

type BlockTag string

type Block struct {
	BlockHash    BlockHash   `json:"bloch_hash"`
	ParentHash   BlockHash   `json:"parent_hash"`
	BlockNumber  uint64      `json:"block_number"`
	Status       BlockStatus `json:"status"`
	Sequencer    Address     `json:"sequencer"`
	NewRoot      Felt        `json:"new_root,omitempty"`
	OldRoot      Felt        `json:"old_root"`
	AcceptedTime int64       `json:"accepted_time"`
	TimeStamp    int64       `json:"time_stamp"`

	TxCount      uint64 `json:"tx_count"`
	TxCommitment Felt   `json:"tx_commitment"`
	TxHashes     []TransactionHash

	EventCount      uint64 `json:"event_count"`
	EventCommitment Felt   `json:"event_commitment"`
}
