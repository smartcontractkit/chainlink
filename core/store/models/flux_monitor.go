package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type FluxMonitorRoundStats struct {
	ID              uint64         `gorm:"primary key;not null;auto_increment"`
	JobRunID        *ID            `gorm:"default:null;foreignkey:JubRunID;association_autoupdate:false;association_autocreate:false"`
	Aggregator      common.Address `gorm:"not null"`
	RoundID         uint32         `gorm:"not null"`
	NumNewRoundLogs uint64         `gorm:"not null;default 0"`
	NumSubmissions  uint64         `gorm:"not null;default 0"`
}
