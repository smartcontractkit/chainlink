package handlers

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// UserCallbackPayload is a response to user request sent to HandleUserMessage().
// Each message needs to receive at most one response on the provided channel.
type UserCallbackPayload struct {
	Msg     *api.Message
	ErrCode api.ErrorCode
	ErrMsg  string
}

// UserMessageHandler implements service-specific logic for managing messages from users and nodes.
// There is one UserMessageHandler object created for each DON.
//
// The lifecycle of a UserMessageHandler object is as follows:
//   - Start() call
//   - a series of HandleUserMessage/HandleNodeMessage calls, executed in parallel
//     (UserMessageHandler needs to guarantee thread safety)
//   - Close() call
type UserMessageHandler interface {
	job.ServiceCtx

	// Each user request is processed by a separate goroutine, which:
	//   1. calls HandleUserMessage
	//   2. waits on callbackCh with a timeout
	HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- UserCallbackPayload) error
}

// NodeMessageHandler implements service-specific logic for managing messages from nodes.
type NodeMessageHandler interface {
	job.ServiceCtx

	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage.
	// should be non-blocking
	HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error
}

// Representation of a DON from a UserMessageHandler's perspective.
type DON interface {
	// Thread-safe
	SendToNode(ctx context.Context, nodeAddress string, msg *api.Message) error
}
