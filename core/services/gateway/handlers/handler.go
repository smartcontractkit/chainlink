package handlers

import (
	"context"
	"encoding/json"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// UserCallbackPayload is a response to user request sent to HandleLegacyUserMessage()/HandleJSONRPCUserMessage().
// Each message needs to receive at most one response on the provided channel.
type UserCallbackPayload struct {
	RawResponse []byte
	ErrorCode   api.ErrorCode
}

// Handler implements service-specific logic for managing messages from users and nodes.
// There is one Handler object created for each DON.
//
// The lifecycle of a Handler object is as follows:
//   - Start() call
//   - a series of HandleUserMessage/HandleNodeMessage calls, executed in parallel
//     (Handler needs to guarantee thread safety)
//   - Close() call
type Handler interface {
	job.ServiceCtx

	// Each user request is processed by a separate goroutine, which:
	//   1. calls HandleUserMessage
	//   2. waits on callbackCh with a timeout
	HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- UserCallbackPayload) error

	// Each user request is processed by a separate goroutine, which:
	//   1. calls HandleUserMessage
	//   2. waits on callbackCh with a timeout
	HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- UserCallbackPayload) error

	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage.
	// should be non-blocking
	// should validate the message inside the response
	HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error

	// The methods support by this Handler.
	// Should be globally unique across all handlers.
	Methods() []string
}

// Representation of a DON from a Handler's perspective.
type DON interface {
	// Thread-safe
	SendToNode(ctx context.Context, nodeAddress string, req *jsonrpc.Request[json.RawMessage]) error
}
