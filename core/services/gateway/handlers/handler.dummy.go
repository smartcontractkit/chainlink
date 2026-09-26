package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

// DummyHandler forwards each request/response without doing any checks.
type dummyHandler struct {
	router         *NodeRouter
	defaultDonID   string
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
}

type savedCallback struct {
	id string
	Callback
}

var _ Handler = (*dummyHandler)(nil)

// NewDummyHandler builds a handler that forwards each request to every node across
// every DON and shard it is given, and relays the first node response back to the
// caller.
func NewDummyHandler(shardedDONs []config.ShardedDONConfig, shardsConnMgrs [][]DON, lggr logger.Logger) (Handler, error) {
	router, err := BuildNodeRouter(shardedDONs, shardsConnMgrs)
	if err != nil {
		return nil, err
	}
	defaultDonID := ""
	if len(shardedDONs) > 0 {
		defaultDonID = shardedDONs[0].DonName
	}
	return &dummyHandler{
		router:         router,
		defaultDonID:   defaultDonID,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           logger.Named(lggr, "DummyHandler."+defaultDonID),
	}, nil
}

func (d *dummyHandler) Methods() []string {
	return []string{"dummy"}
}

func (d *dummyHandler) HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callback Callback) error {
	var msg api.Message
	if jsonRequest.Params != nil {
		if err := json.Unmarshal(*jsonRequest.Params, &msg); err != nil {
			return err
		}
	}
	msg.Body.MessageID = jsonRequest.ID
	if msg.Body.Method == "" {
		msg.Body.Method = jsonRequest.Method
	}
	if msg.Body.DonID == "" {
		msg.Body.DonID = d.defaultDonID
	}
	return d.HandleLegacyUserMessage(ctx, &msg, callback)
}

func (d *dummyHandler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callback Callback) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageID] = &savedCallback{msg.Body.MessageID, callback}
	d.mu.Unlock()
	params, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	rawParams := json.RawMessage(params)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      msg.Body.MessageID,
		Method:  msg.Body.Method,
		Params:  &rawParams,
	}
	for _, member := range d.router.Members() {
		memberDON, ok := d.router.DONFor(member.Address)
		if !ok {
			err = errors.Join(err, fmt.Errorf("no connection manager found for node %s", member.Address))
			continue
		}
		err = errors.Join(err, memberDON.SendToNode(ctx, member.Address, req))
	}
	return err
}

func (d *dummyHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var msg api.Message
	err := json.Unmarshal(*resp.Result, &msg)
	if err != nil {
		return err
	}
	msg.Body.MessageID = resp.ID
	err = msg.Validate()
	if err != nil {
		return err
	}
	if nodeAddr != msg.Body.Sender {
		return fmt.Errorf("node address %s does not match message sender %s", nodeAddr, msg.Body.Sender)
	}
	d.mu.Lock()
	savedCb, found := d.savedCallbacks[msg.Body.MessageID]
	delete(d.savedCallbacks, msg.Body.MessageID)
	d.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		codec := api.JSONRPCCodec{}
		return savedCb.SendResponse(UserCallbackPayload{RawResponse: codec.EncodeLegacyResponse(&msg), ErrorCode: api.NoError})
	}
	return nil
}

func (d *dummyHandler) Start(context.Context) error {
	return nil
}

func (d *dummyHandler) Close() error {
	return nil
}
