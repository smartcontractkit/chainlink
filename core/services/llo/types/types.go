package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"

	"github.com/ethereum/go-ethereum/common"
)

func ChannelDefinitionCacheFilterName(addr common.Address, donID uint32) string {
	return logpoller.FilterName("OCR3 LLO ChannelDefinitionCachePoller", addr.String(), strconv.FormatUint(uint64(donID), 10))
}

type PersistedDefinitions struct {
	ChainSelector uint64         `db:"chain_selector"`
	Address       common.Address `db:"addr"`
	Definitions   Definitions    `db:"definitions"`
	// The block number in which the log for this definitions was emitted
	BlockNum  int64     `db:"block_num"`
	DonID     uint32    `db:"don_id"`
	Version   uint32    `db:"version"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Definitions map[uint32]SourceDefinition

func (d *Definitions) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Data: value is not []byte")
	}
	if bytes == nil {
		*d = nil
		return nil
	}
	if len(bytes) == 0 {
		*d = Definitions{}
		return nil
	}

	return json.Unmarshal(bytes, d)
}

func (d Definitions) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Trigger contains the information needed to fetch channel definitions from a URL.
// It is created from on-chain events and includes the source, URL, expected SHA hash,
// block number, version (for owner sources), and transaction hash.
type Trigger struct {
	Source   uint32      `json:"source"`
	URL      string      `json:"url"`
	SHA      [32]byte    `json:"sha"`
	BlockNum int64       `json:"block_num"`
	LogIndex int64       `json:"log_index"`
	Version  uint32      `json:"version"`
	TxHash   common.Hash `json:"tx_hash"`
}

// SourceDefinition stores the channel definitions fetched from a specific source along with
// the trigger that initiated the fetch.
type SourceDefinition struct {
	Trigger     Trigger                     `json:"trigger"`
	Definitions llotypes.ChannelDefinitions `json:"definitions"`
}
