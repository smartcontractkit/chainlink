package logpoller

import (
	"context"

	"github.com/ethereum/go-ethereum"
)

func (lp *logPoller) PollAndSaveLogs(ctx context.Context, currentBlockNumber int64) int64 {
	lp.pollAndSaveLogs(ctx, currentBlockNumber)
	lastProcessed, _ := lp.orm.SelectLatestBlock()
	return lastProcessed.BlockNumber + 1
}

func (lp *logPoller) Filter() ethereum.FilterQuery {
	return lp.filter(nil, nil, nil)
}

func (o *ORM) SelectLogsByBlockRange(start, end int64) ([]Log, error) {
	return o.selectLogsByBlockRange(start, end)
}
