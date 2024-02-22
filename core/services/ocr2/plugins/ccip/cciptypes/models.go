package cciptypes

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Address string

// TODO: make JSON marshal/unmarshal non-evm specific.
// Make sure we have casing compatibility with old versions.
func (a *Address) UnmarshalJSON(bytes []byte) error {
	vStr := strings.Trim(string(bytes), `"`)
	if !common.IsHexAddress(vStr) {
		return fmt.Errorf("invalid address: %s", vStr)
	}
	*a = Address(common.HexToAddress(vStr).String())
	return nil
}

func (a *Address) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strings.ToLower(string(*a)) + `"`), nil
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
