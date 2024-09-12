package integration_tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

var (
	_ capabilities.ActionCapability = &logTarget{}
)

const triggerIDm = "log-target@1.0.0"

type sink struct {
	services.StateMachine
	targets []logTarget

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func newSink() *sink {
	return &sink{
		stopCh: make(services.StopChan),
	}
}

func (r *sink) Start(ctx context.Context) error {
	return r.StartOnce("sink", func() error {
		return nil
	})
}

func (r *sink) Close() error {
	return r.StopOnce("sink", func() error {
		close(r.stopCh)
		r.wg.Wait()
		return nil
	})
}

// func (r *sink) getLogs(reportList []*datastreams.FeedReport) {
// 	for _, trigger := range r.triggers {
// 		resp, err := wrapReports(reportList, "1", 12, datastreams.Metadata{})
// 		if err != nil {
// 			panic(err)
// 		}
// 		trigger.sendResponse(resp)
// 	}
// }

func (r *sink) getNewTarget(t *testing.T) *logTarget {
	target := logTarget{t: t, toSend: make(chan capabilities.TriggerResponse, 1000),
		wg: &r.wg, stopCh: r.stopCh}
	r.targets = append(r.targets, target)
	return &target
}

type logTarget struct {
	t      *testing.T
	cancel context.CancelFunc
	toSend chan capabilities.TriggerResponse

	wg     *sync.WaitGroup
	stopCh services.StopChan
}

func (lt *logTarget) Execute(ctx context.Context, rawRequest capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	fmt.Println("##########################################")
	return capabilities.CapabilityResponse{}, nil
}

func (lt *logTarget) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return capabilities.MustNewCapabilityInfo(
		triggerIDm,
		capabilities.CapabilityTypeTarget,
		"issues a trigger when a report is received.",
	), nil
}

func (lt *logTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (lt *logTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
