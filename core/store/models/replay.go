package models

import (
	"fmt"
)

type ReplayBlocksRequest struct {
	BlockNumber int64 `json:"blockNumber" gorm:"type:text"`
}

func ValidateReplayRequest(request *ReplayBlocksRequest) error {
	if request.BlockNumber < 0 {
		return fmt.Errorf("cannot replay from a negative block number: %s", request.BlockNumber)
	}

	return nil
}
