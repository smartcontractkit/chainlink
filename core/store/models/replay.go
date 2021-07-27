package models

type ReplayBlocksRequest struct {
	BlockNumber int64 `json:"blockNumber" gorm:"type:text"`
}
