package ocr2keepers

import (
	"encoding/json"
	"fmt"
	"strings"
)

type UpkeepIdentifier []byte

type BlockKey string

type UpkeepKey []byte

type UpkeepResult interface{}

func upkeepKeysToString(keys []UpkeepKey) string {
	keysStr := make([]string, len(keys))
	for i, key := range keys {
		keysStr[i] = string(key)
	}

	return strings.Join(keysStr, ", ")
}

type PerformLog struct {
	Key             UpkeepKey
	TransmitBlock   BlockKey
	Confirmations   int64
	TransactionHash string
}

type StaleReportLog struct {
	Key             UpkeepKey
	TransmitBlock   BlockKey
	Confirmations   int64
	TransactionHash string
}

type BlockHistory []BlockKey

func (bh BlockHistory) Latest() (BlockKey, error) {
	if len(bh) == 0 {
		return BlockKey(""), fmt.Errorf("empty block history")
	}

	return bh[0], nil
}

func (bh BlockHistory) Keys() []BlockKey {
	return bh
}

func (bh *BlockHistory) UnmarshalJSON(b []byte) error {
	var raw []string

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	output := make([]BlockKey, len(raw))
	for i, value := range raw {
		output[i] = BlockKey(value)
	}

	*bh = output

	return nil
}
