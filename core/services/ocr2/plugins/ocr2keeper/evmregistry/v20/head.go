package evm

import (
	"context"
	"strconv"

	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"

	"github.com/smartcontractkit/chainlink-evm/pkg/heads"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

type HeadProvider struct {
	ht         heads.Tracker
	hb         heads.Broadcaster
	chHead     chan ocr2keepers.BlockKey
	subscribed bool
}

// HeadTicker provides external access to the heads channel
func (hw *HeadProvider) HeadTicker() chan ocr2keepers.BlockKey {
	if !hw.subscribed {
		_, _ = hw.hb.Subscribe(&headWrapper{c: hw.chHead})
		hw.subscribed = true
	}
	return hw.chHead
}

func (hw *HeadProvider) LatestBlock() int64 {
	lc := hw.ht.LatestChain()
	if lc == nil {
		return 0
	}
	return lc.Number
}

// send does a non-blocking send of the key on c.
func send(c chan ocr2keepers.BlockKey, k ocr2keepers.BlockKey) {
	select {
	case c <- k:
	default:
	}
}

type headWrapper struct {
	c chan ocr2keepers.BlockKey
}

func (w *headWrapper) OnNewLongestChain(ctx context.Context, head *types.Head) {
	var bl int64
	if head != nil {
		bl = head.Number
	}

	send(w.c, ocr2keepers.BlockKey(strconv.FormatInt(bl, 10)))
}
