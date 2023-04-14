package gateway

import (
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type requestState struct {
	msg      *Message
	cb       Callback
	ackCount int
}

type functionsHandler struct {
	pendingCallbacks map[string]*requestState
	mu               sync.Mutex
	connManager      ConnectionManager
	donConfig        *GatewayDONConfig
	lggr             logger.SugaredLogger
}

func NewFunctionsHandler(lggr logger.SugaredLogger) Handler {
	return &functionsHandler{pendingCallbacks: make(map[string]*requestState), lggr: lggr}
}

func (h *functionsHandler) Init(connMgr ConnectionManager, donConfig *GatewayDONConfig) {
	h.connManager = connMgr
	h.donConfig = donConfig
}

func (h *functionsHandler) HandleUserMessage(msg *Message, cb Callback) {
	h.lggr.Info("Gateway: Functions handler, user message")
	reqId := msg.Payload["request_id"].(string)
	// Save the callback and forward request to all nodes in the DON
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pendingCallbacks[reqId] = &requestState{msg: msg, cb: cb, ackCount: 0}
	if msg.Method == "secrets_upload" {
		h.lggr.Info("Gateway: Functions handler, sending secrets_upload to all nodes")
		h.connManager.SendToAll(msg)
	}
}

func (h *functionsHandler) HandleNodeMessage(msg *Message, nodeAddr string) {
	reqId := msg.Payload["request_id"].(string)
	if msg.Method == "secrets_upload_ack" {
		h.mu.Lock()
		defer h.mu.Unlock()
		state := h.pendingCallbacks[reqId]
		state.ackCount += 1
		h.lggr.Info("Gateway: Functions handler, received secrets_upload_ack from ", nodeAddr, " ", state.ackCount)
		if state.ackCount == 2 {
			// Return success to caller after receiving confirmations from F+1 nodes
			h.lggr.Info("Gateway: Functions handler, enough node ACKs to return an ACK to user")
			resp := Message{DonId: state.msg.DonId, MessageId: state.msg.MessageId, Method: "ack"}
			state.cb.SendResponse(&resp)
		}
	}
}
