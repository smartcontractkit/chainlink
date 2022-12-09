package evm

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

type HeadWatcher struct {
	client     client.Client
	ctx        context.Context
	mu         sync.Mutex
	latest     int64
	chReady    chan struct{}
	dataInChan bool
}

// OnNewHead should continue running until the context ends
func (hw *HeadWatcher) OnNewHead(ctx context.Context, f func(blockKey types.BlockKey)) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error encountered while running head func: %s", err)
			debug.PrintStack()
		}
	}()

	for {
		select {
		case <-hw.chReady:
			hw.mu.Lock()
			hw.dataInChan = false
			val := hw.latest
			hw.mu.Unlock()
			f(types.BlockKey(fmt.Sprintf("%d", val)))
		case <-ctx.Done():
			return nil
		}
	}
}

func (hw *HeadWatcher) Watch(ctx context.Context) error {
	hw.mu.Lock()
	hw.ctx = ctx
	hw.mu.Unlock()

	// subscribe to new heads, set latest, and call the headFunc
	chHead := make(chan *evmtypes.Head)
	sub, err := hw.client.SubscribeNewHead(hw.ctx, chHead)
	if err != nil {
		return err
	}

	go func(ctx context.Context, ch chan *evmtypes.Head, f func(int64), s ethereum.Subscription) {
		for {
			select {
			case h := <-ch:
				f(h.Number)
			case <-ctx.Done():
				s.Unsubscribe()
				return
			}
		}
	}(hw.ctx, chHead, hw.update, sub)

	return nil
}

func (hw *HeadWatcher) LatestBlock() int64 {
	hw.mu.Lock()
	defer hw.mu.Unlock()
	return hw.latest
}

func (hw *HeadWatcher) update(block int64) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	hw.latest = block
	if !hw.dataInChan {
		hw.dataInChan = true
		hw.chReady <- struct{}{}
	}
}
