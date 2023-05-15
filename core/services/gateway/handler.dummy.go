package gateway

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

// DummyHandler forwards each request/response without doing any checks.
type dummyHandler struct {
	donConfig      *DONConfig
	connMgr        DONConnectionManager
	savedCallbacks map[string]chan UserCallbackPayload
	mu             sync.Mutex
}

var _ Handler = (*dummyHandler)(nil)

func NewDummyHandler(donConfig *DONConfig, connMgr DONConnectionManager) (Handler, error) {
	return &dummyHandler{
		donConfig:      donConfig,
		connMgr:        connMgr,
		savedCallbacks: make(map[string]chan UserCallbackPayload),
	}, nil
}

func (d *dummyHandler) HandleUserMessage(ctx context.Context, msg *Message, callbackChan chan UserCallbackPayload) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = callbackChan
	connMgr := d.connMgr
	d.mu.Unlock()

	var err error
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, connMgr.SendToNode(ctx, member.Address, msg))
	}
	return err
}

func (d *dummyHandler) HandleNodeMessage(ctx context.Context, msg *Message, nodeAddr string) error {
	d.mu.Lock()
	callbackChan := d.savedCallbacks[msg.Body.MessageId]
	delete(d.savedCallbacks, msg.Body.MessageId)
	d.mu.Unlock()

	if callbackChan != nil {
		callbackChan <- UserCallbackPayload{Msg: msg, ErrCode: NoError, ErrMsg: ""}
		close(callbackChan)
	}
	return nil
}

func (d *dummyHandler) Start(context.Context) error {
	return nil
}

func (d *dummyHandler) Close() error {
	return nil
}
