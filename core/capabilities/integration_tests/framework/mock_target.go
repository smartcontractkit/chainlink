package framework

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var (
	_ capabilities.ActionCapability = &mockTarget{}
)

type TargetSink struct {
	services.StateMachine
	targetID   string
	targetName string
	version    string

	targets []mockTarget
	Sink    chan capabilities.CapabilityRequest
}

func NewTargetSink(targetName string, version string) *TargetSink {
	return &TargetSink{
		targetID:   targetName + "@" + version,
		targetName: targetName,
		version:    version,
		Sink:       make(chan capabilities.CapabilityRequest, 1000),
	}
}

func (ts *TargetSink) GetTargetVersion() string {
	return ts.version
}

func (ts *TargetSink) GetTargetName() string {
	return ts.targetName
}

func (ts *TargetSink) GetTargetID() string {
	return ts.targetID
}

func (ts *TargetSink) Start(ctx context.Context) error {
	return ts.StartOnce("TargetSinkService", func() error {
		return nil
	})
}

func (ts *TargetSink) Close() error {
	return ts.StopOnce("TargetSinkService", func() error {
		return nil
	})
}

func (ts *TargetSink) CreateNewTarget(t *testing.T) capabilities.TargetCapability {
	target := mockTarget{
		t:        t,
		targetID: ts.targetID,
		ch:       ts.Sink,
	}
	ts.targets = append(ts.targets, target)
	return &target
}

type mockTarget struct {
	t        *testing.T
	targetID string
	ch       chan capabilities.CapabilityRequest
}

func (mt *mockTarget) Execute(ctx context.Context, rawRequest capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	mt.ch <- rawRequest
	return capabilities.CapabilityResponse{}, nil
}

func (mt *mockTarget) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		mt.targetID,
		capabilities.CapabilityTypeTarget,
		"mock target for target ID "+mt.targetID,
	), nil
}

func (mt *mockTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (mt *mockTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
