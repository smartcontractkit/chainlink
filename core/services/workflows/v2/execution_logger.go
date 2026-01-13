package v2

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/executionlog"
)

type ExecutionCallbackHelper struct {
	ExecutionHelper  host.ExecutionHelper
	OnExecutionEvent func(event *executionlog.ExecutionEvent)
}

var _ host.ExecutionHelper = (*ExecutionCallbackHelper)(nil)

func (e ExecutionCallbackHelper) CallCapability(ctx context.Context, request *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	response, err := e.ExecutionHelper.CallCapability(ctx, request)

	event := &executionlog.ExecutionEvent{
		Event: &executionlog.ExecutionEvent_CapabilityEvent{
			CapabilityEvent: &executionlog.CapabilityEvent{
				Request:  request,
				Response: response,
			},
		},
	}
	if err != nil {
		event.GetCapabilityEvent().Error = err.Error()
	}
	if e.OnExecutionEvent != nil {
		e.OnExecutionEvent(event)
	}

	return response, err
}

func (e ExecutionCallbackHelper) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	responses, err := e.ExecutionHelper.GetSecrets(ctx, request)

	event := &executionlog.ExecutionEvent{
		Event: &executionlog.ExecutionEvent_SecretsEvent{
			SecretsEvent: &executionlog.SecretsEvent{
				Request: request,
			},
		},
	}
	if err != nil {
		event.GetSecretsEvent().Error = err.Error()
	} else if len(responses) > 0 {
		// Use the first response as the representative response
		event.GetSecretsEvent().Response = responses[0]
	}
	if e.OnExecutionEvent != nil {
		e.OnExecutionEvent(event)
	}

	return responses, err
}

func (e ExecutionCallbackHelper) GetWorkflowExecutionID() string {
	workflowID := e.ExecutionHelper.GetWorkflowExecutionID()

	if e.OnExecutionEvent != nil {
		event := &executionlog.ExecutionEvent{
			Event: &executionlog.ExecutionEvent_WorkflowIdEvent{
				WorkflowIdEvent: workflowID,
			},
		}
		e.OnExecutionEvent(event)
	}

	return workflowID
}

func (e ExecutionCallbackHelper) GetNodeTime() time.Time {
	nodeTime := e.ExecutionHelper.GetNodeTime()

	if e.OnExecutionEvent != nil {
		event := &executionlog.ExecutionEvent{
			Event: &executionlog.ExecutionEvent_NodeTimeEvent{
				NodeTimeEvent: &executionlog.NodeTimeEvent{
					Response: timestamppb.New(nodeTime),
				},
			},
		}
		e.OnExecutionEvent(event)
	}

	return nodeTime
}

func (e ExecutionCallbackHelper) GetDONTime() (time.Time, error) {
	donTime, err := e.ExecutionHelper.GetDONTime()

	event := &executionlog.ExecutionEvent{
		Event: &executionlog.ExecutionEvent_DonTimeEvent{
			DonTimeEvent: &executionlog.DonTimeEvent{},
		},
	}
	if err != nil {
		event.GetDonTimeEvent().Error = err.Error()
	} else {
		event.GetDonTimeEvent().Response = timestamppb.New(donTime)
	}
	if e.OnExecutionEvent != nil {
		e.OnExecutionEvent(event)
	}

	return donTime, err
}

func (e ExecutionCallbackHelper) EmitUserLog(log string) error {
	err := e.ExecutionHelper.EmitUserLog(log)

	event := &executionlog.ExecutionEvent{
		Event: &executionlog.ExecutionEvent_EmitLogEvent{
			EmitLogEvent: &executionlog.EmitLogEvent{
				Log: log,
			},
		},
	}
	if err != nil {
		event.GetEmitLogEvent().Error = err.Error()
	}
	if e.OnExecutionEvent != nil {
		e.OnExecutionEvent(event)
	}

	return err
}
