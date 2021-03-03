package models

import (
	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
)

type FluxMonitorRoundStats struct {
	ID              uint64         `gorm:"primary key;not null;auto_increment"`
	JobRunID        uuid.NullUUID  `gorm:"default:null;foreignkey:JobRunID"`
	Aggregator      common.Address `gorm:"not null"`
	RoundID         uint32         `gorm:"not null"`
	NumNewRoundLogs uint64         `gorm:"not null;default 0"`
	NumSubmissions  uint64         `gorm:"not null;default 0"`
}
