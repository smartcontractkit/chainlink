package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/null"
)

// FluxMonitorRoundStatsV2 defines the stats for a round
type FluxMonitorRoundStatsV2 struct {
	ID              uint64         `gorm:"primary key;not null;auto_increment"`
	PipelineRunID   null.Int64     `gorm:"default:null"`
	Aggregator      common.Address `gorm:"not null"`
	RoundID         uint32         `gorm:"not null"`
	NumNewRoundLogs uint64         `gorm:"not null;default 0"`
	NumSubmissions  uint64         `gorm:"not null;default 0"`
}
