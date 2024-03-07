package cciptypes

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Address string

func (a *Address) UnmarshalJSON(b []byte) error {
	vStr := strings.Trim(string(b), `"`)
	if !common.IsHexAddress(vStr) {
		return fmt.Errorf("invalid address: %s", vStr)
	}
	*a = Address(common.HexToAddress(vStr).String())
	return nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ToLower(string(a)) + `"`), nil
}

func (a Address) MarshalText() (text []byte, err error) {
	return []byte(strings.ToLower(string(a))), nil
}

func (a *Address) UnmarshalText(text []byte) error {
	vStr := string(text)
	if !common.IsHexAddress(vStr) {
		return fmt.Errorf("invalid address: %s", vStr)
	}
	*a = Address(common.HexToAddress(vStr).String())
	return nil
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
