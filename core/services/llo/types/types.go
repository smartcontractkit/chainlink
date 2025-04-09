package types

import (
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
	ChainSelector uint64                      `db:"chain_selector"`
	Address       common.Address              `db:"addr"`
	Definitions   llotypes.ChannelDefinitions `db:"definitions"`
	// The block number in which the log for this definitions was emitted
	BlockNum  int64     `db:"block_num"`
	DonID     uint32    `db:"don_id"`
	Version   uint32    `db:"version"`
	UpdatedAt time.Time `db:"updated_at"`
}
