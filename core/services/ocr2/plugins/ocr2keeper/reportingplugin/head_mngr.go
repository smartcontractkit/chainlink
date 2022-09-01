package reportingplugin

import (
	"context"
	"sync"

	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type headsMngr struct {
	logger           logger.Logger
	headBroadcaster  httypes.HeadBroadcasterRegistry
	unsubscribeHeads func()
	mailbox          *utils.Mailbox[*evmtypes.Head]
	chStop           chan struct{}

	currentHead     *evmtypes.Head
	currentHeadLock sync.Mutex
}

func newHeadsMngr(logger logger.Logger, headBroadcaster httypes.HeadBroadcaster) *headsMngr {
	return &headsMngr{
		logger:          logger,
		headBroadcaster: headBroadcaster,
		mailbox:         utils.NewMailbox[*evmtypes.Head](1),
		chStop:          make(chan struct{}),
	}
}

func (h *headsMngr) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	h.mailbox.Deliver(head)
}

func (h *headsMngr) setCurrentHead(head *evmtypes.Head) {
	h.currentHeadLock.Lock()
	h.currentHead = head
	h.currentHeadLock.Unlock()
}

func (h *headsMngr) getCurrentHead() *evmtypes.Head {
	h.currentHeadLock.Lock()
	defer h.currentHeadLock.Unlock()

	return h.currentHead
}

func (h *headsMngr) start() {
	go func() {
		for {
			select {
			case <-h.chStop:
				return
			case <-h.mailbox.Notify():
				head, exists := h.mailbox.Retrieve()
				if !exists {
					h.logger.Info("No head to retrieve. It might have been skipped")
					continue
				}

				h.setCurrentHead(head)
			}
		}
	}()

	latestHead, unsubscribeHeads := h.headBroadcaster.Subscribe(h)
	if latestHead != nil {
		h.mailbox.Deliver(latestHead)
	}

	h.unsubscribeHeads = unsubscribeHeads
}

func (h *headsMngr) stop() {
	close(h.chStop)
	h.unsubscribeHeads()
}
