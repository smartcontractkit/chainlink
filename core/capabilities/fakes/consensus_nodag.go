package fakes

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/anypb"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
)

type fakeConsensusNoDAG struct {
	services.Service
	eng *services.Engine
}

var _ services.Service = (*fakeConsensus)(nil)
var _ commonCap.ExecutableCapability = (*fakeConsensusNoDAG)(nil)

const consensusNoDAGCapID = "consensus@1.0.0"

func NewFakeConsensusNoDAG(lggr logger.Logger) (*fakeConsensusNoDAG, error) {
	fc := &fakeConsensusNoDAG{}
	fc.Service, fc.eng = services.Config{
		Name:  "fakeConsensusNoDAG",
		Start: fc.start,
		Close: fc.close,
	}.NewServiceEngine(lggr)
	return fc, nil
}

func (fc *fakeConsensusNoDAG) start(ctx context.Context) error {
	return nil
}

func (fc *fakeConsensusNoDAG) close() error {
	return nil
}

// NOTE: This fake capability currently bounces back the request payload, ignoring everything else.
// When the real NoDAG consensus OCR plugin is ready, it should be used here, similarly to how the V1 fake works.
func (fc *fakeConsensusNoDAG) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	resp := commonCap.CapabilityResponse{}
	inputs := &pb.SimpleConsensusInputs{}
	err := request.Payload.UnmarshalTo(inputs)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal SimpleConsensusInputs err=%w, payload=%v", err, request.Payload)
	}
	fc.eng.Infow("Executing Fake Consensus NoDAG", "inputs", inputs)
	anyProto, err := anypb.New(inputs.GetValue())
	if err != nil {
		return resp, fmt.Errorf("failed to marshal SimpleConsensusInputs value err=%w, value=%v", err, inputs.GetValue())
	}
	resp.Metadata = commonCap.ResponseMetadata{}
	resp.Payload = anyProto
	return resp, nil
}

func (fc *fakeConsensusNoDAG) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fc.eng.Infow("Registering to Fake Consensus NoDAG", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *fakeConsensusNoDAG) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fc.eng.Infow("Unregistering from Fake Consensus NoDAG", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *fakeConsensusNoDAG) Info(ctx context.Context) (commonCap.CapabilityInfo, error) {
	return commonCap.CapabilityInfo{
		ID:             consensusNoDAGCapID,
		CapabilityType: commonCap.CapabilityTypeConsensus,
		Description:    "Fake OCR Consensus NoDAG",
		DON:            &commonCap.DON{},
		IsLocal:        true,
	}, nil
}
