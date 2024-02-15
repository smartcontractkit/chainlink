package ccip

import (
	"encoding/hex"
)

type Address string

type Hash [32]byte

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}

type TxMeta struct {
	BlockTimestampUnixMilli int64
	BlockNumber             uint64
	TxHash                  string
	LogIndex                uint64
}
