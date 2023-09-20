package legacygasstation

import (
	"context"

	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

var _ legacygasstation.LogPoller = &logPollerAdapter{}

type logPollerAdapter struct {
	lp logpoller.LogPoller
}

func NewLogPollerAdapter(lp logpoller.LogPoller) *logPollerAdapter {
	return &logPollerAdapter{lp: lp}
}

func (a *logPollerAdapter) LatestBlock(ctx context.Context) (int64, error) {
	return a.lp.LatestBlock(pg.WithParentCtx(ctx))
}

func (a logPollerAdapter) FilterName(id string, args ...any) string {
	return logpoller.FilterName(id, args)
}

func (a *logPollerAdapter) RegisterFilter(ctx context.Context, name string, eventSigs []gethcommon.Hash, addresses []gethcommon.Address) error {
	return a.lp.RegisterFilter(logpoller.Filter{
		Name:      name,
		EventSigs: eventSigs,
		Addresses: addresses,
	}, pg.WithParentCtx(ctx))
}

func (a *logPollerAdapter) IndexedLogsByBlockRange(ctx context.Context, start, end int64, eventSig gethcommon.Hash, address gethcommon.Address, topicIndex int, topicValues []gethcommon.Hash) ([]gethtypes.Log, error) {
	logs, err := a.lp.IndexedLogsByBlockRange(start, end, eventSig, address, topicIndex, topicValues, pg.WithParentCtx(ctx))
	if err != nil {
		return nil, err
	}
	var gethlogs []gethtypes.Log
	for _, l := range logs {
		gethlogs = append(gethlogs, l.ToGethLog())
	}
	return gethlogs, nil
}
