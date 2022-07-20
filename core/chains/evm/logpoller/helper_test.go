package logpoller

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (lp *logPoller) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
	return lp.pollAndSaveLogs(ctx, currentBlockNumber)
}

func (lp *logPoller) FilterAddresses() []common.Address {
	return lp.filterAddresses()
}

func (lp *logPoller) FilterTopics() [][]common.Hash {
	return lp.filterTopics()
}

func (o *ORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
	return o.selectLogsByBlockRange(start, end)
}
