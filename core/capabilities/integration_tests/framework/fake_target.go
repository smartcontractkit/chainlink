package framework

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var (
	_ capabilities.ActionCapability = &target{}
)

type TargetSink struct {
	services.StateMachine
	targetID   string
	targetName string
	version    string

	targets []target
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
	trg := target{
		t:        t,
		targetID: ts.targetID,
		ch:       ts.Sink,
	}
	ts.targets = append(ts.targets, trg)
	return &trg
}

type target struct {
	t        *testing.T
	targetID string
	ch       chan capabilities.CapabilityRequest
}

func (mt *target) Execute(ctx context.Context, rawRequest capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	mt.ch <- rawRequest
	return capabilities.CapabilityResponse{}, nil
}

func (mt *target) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		mt.targetID,
		capabilities.CapabilityTypeTarget,
		"fake target for target ID "+mt.targetID,
	), nil
}

func (mt *target) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (mt *target) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
