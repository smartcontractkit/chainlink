package gateway

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// UserCallbackPayload is a response to user request sent to HandleUserMessage().
// Each message needs to receive at most one response on the provided channel.
type UserCallbackPayload struct {
	Msg     *Message
	ErrCode ErrorCode
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
//go:generate mockery --quiet --name Handler --output ./mocks/ --case=underscore

type Handler interface {
	job.ServiceCtx

	// Each user request is processed by a separate goroutine, which:
	//   1. calls HandleUserMessage
	//   2. waits on callbackCh with a timeout
	HandleUserMessage(ctx context.Context, msg *Message, callbackCh chan<- UserCallbackPayload) error

	// Handlers should not make any assumptions about goroutines calling HandleNodeMessage
	HandleNodeMessage(ctx context.Context, msg *Message, nodeAddr string) error
}

type HandlerType = string

const (
	Dummy HandlerType = "dummy"
)

func NewHandler(handlerType HandlerType, donConfig *DONConfig, connMgr DONConnectionManager) (Handler, error) {
	switch handlerType {
	case Dummy:
		return NewDummyHandler(donConfig, connMgr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
