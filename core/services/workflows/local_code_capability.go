package workflows

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	pocWorkflow "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

type localCodeCapability struct {
	Workflow       *pocWorkflow.Spec
	CapabilityType capabilities.CapabilityType
	Id             string
}

func (l *localCodeCapability) Info(_ context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.CapabilityInfo{
		ID:             l.Id,
		CapabilityType: l.CapabilityType,
		Description:    "Run local code",
		Version:        "1",
	}, nil
}

func (l *localCodeCapability) RegisterToWorkflow(_ context.Context, _ capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (l *localCodeCapability) UnregisterFromWorkflow(_ context.Context, _ capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func (l *localCodeCapability) Execute(_ context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	s, ok := l.Workflow.LocalExecutions[request.Metadata.StepRef]
	if !ok {
		return nil, fmt.Errorf("step %s not found", request.Metadata.StepRef)
	}
	results, cont, err := s.Run(request.Inputs)
	ch := make(chan capabilities.CapabilityResponse, 1)
	if cont || err != nil {
		ch <- capabilities.CapabilityResponse{
			Value: results,
			Err:   err,
		}
	}
	close(ch)
	return ch, nil
}

var _ capabilities.CallbackCapability = &localCodeCapability{}
