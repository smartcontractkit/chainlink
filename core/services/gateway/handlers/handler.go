package handlers

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

//go:generate mockery --quiet --name Handler --output ./mocks/ --case=underscore
//go:generate mockery --quiet --name DON --output ./mocks/ --case=underscore

// UserCallbackPayload is a response to user request sent to HandleUserMessage().
// Each message needs to receive at most one response on the provided channel.
type UserCallbackPayload struct {
	Msg     *api.Message
	ErrCode api.ErrorCode
	ErrMsg  string
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
	HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- UserCallbackPayload) error

	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage
	HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error
}

// Representation of a DON from a Handler's perspective.
type DON interface {
	// Thread-safe
	SendToNode(ctx context.Context, nodeAddress string, msg *api.Message) error
}
