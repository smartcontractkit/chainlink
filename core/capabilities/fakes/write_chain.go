package fakes

import (
	"context"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type fakeWriteChain struct {
	services.Service
	eng      *services.Engine
	targetID string
}

var _ services.Service = (*fakeWriteChain)(nil)
var _ commonCap.TargetCapability = (*fakeWriteChain)(nil)

func (wc *fakeWriteChain) Info(ctx context.Context) (commonCap.CapabilityInfo, error) {
	return commonCap.CapabilityInfo{
		ID:             wc.targetID,
		CapabilityType: commonCap.CapabilityTypeTarget,
		Description:    "Fake Write Chain",
		DON:            &commonCap.DON{},
		IsLocal:        true,
	}, nil
}

func (wc *fakeWriteChain) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	wc.eng.Infow("Registered to Fake Write Chain", "targetID", wc.targetID, "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (wc *fakeWriteChain) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	wc.eng.Infow("Unregistered from Fake Write Chain", "targetID", wc.targetID, "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (wc *fakeWriteChain) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	wc.eng.Infow("Executed Fake Write Chain", "targetID", wc.targetID, "workflowID", request.Metadata.WorkflowID, "executionID", request.Metadata.WorkflowExecutionID)
	return commonCap.CapabilityResponse{Value: &values.Map{}}, nil
}

func NewFakeWriteChain(lggr logger.Logger, targetID string) *fakeWriteChain {
	wc := &fakeWriteChain{targetID: targetID}
	wc.Service, wc.eng = services.Config{
		Name:  "fakeWriteChain",
		Start: wc.start,
	}.NewServiceEngine(lggr)
	return wc
}

func (wc *fakeWriteChain) start(_ context.Context) error {
	wc.eng.Infow("Fake Write Chain started", "targetID", wc.targetID)
	return nil
}
