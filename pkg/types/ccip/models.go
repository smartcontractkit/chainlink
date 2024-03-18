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
}
