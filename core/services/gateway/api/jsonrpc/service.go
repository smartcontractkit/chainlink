package jsonrpc

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// Service implements method-specific logic for managing requests from users.
// There is one Service object created for each top-level method.
//
// The lifecycle of a Service object is as follows:
//   - Start() call
//   - a series of HandleUserRequest calls, executed in parallel
//     (Service needs to guarantee thread safety)
//   - Close() call
type Service interface {
	job.ServiceCtx

	// Each user request is processed by a separate goroutine, which:
	//   1. calls HandleUserRequest
	//   2. waits on callbackCh with a timeout
	HandleUserRequest(ctx context.Context, request *Request, callbackCh chan<- *Response) error
}
