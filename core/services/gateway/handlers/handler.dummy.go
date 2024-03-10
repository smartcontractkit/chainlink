package handlers

import (
	"context"
	"sync"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// DummyHandler forwards each request/response without doing any checks.
type dummyHandler struct {
	donConfig      *config.DONConfig
	don            DON
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
}

type savedCallback struct {
	id         string
	callbackCh chan<- UserCallbackPayload
}

var _ Handler = (*dummyHandler)(nil)

func NewDummyHandler(donConfig *config.DONConfig, don DON, lggr logger.Logger) (Handler, error) {
	return &dummyHandler{
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           lggr.Named("DummyHandler." + donConfig.DonId),
	}, nil
}

func (d *dummyHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- UserCallbackPayload) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()

	var err error
	// Send to all nodes.
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
	return err
}

func (d *dummyHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	d.mu.Lock()
	savedCb, found := d.savedCallbacks[msg.Body.MessageId]
	delete(d.savedCallbacks, msg.Body.MessageId)
	d.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		savedCb.callbackCh <- UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (d *dummyHandler) Start(context.Context) error {
	return nil
}

func (d *dummyHandler) Close() error {
	return nil
}
