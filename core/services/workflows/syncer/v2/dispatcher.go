package v2

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// waitLoopService is a service that runs the queue Wait loop until context cancellation.
type waitLoopService struct {
	services.Service
	run func(context.Context) error
}

func newWaitLoopService(run func(context.Context) error, lggr logger.Logger) services.Service {
	w := &waitLoopService{run: run}
	w.Service, _ = services.Config{
		Name:  "ConsensusEventDispatcherWaitLoop",
		Start: w.run,
	}.NewServiceEngine(lggr)
	return w
}

// ConsensusEventDispatcher routes trigger events from a shared queue (or OCR callback) to the
// correct workflow engine. Uses EngineRegistry to look up engines by workflowID.
//
// For standard queue: runs a Wait loop, popping events and routing to engines.
// For OCR queue: receives events via OnConsensusEvent from the Transmitter.
type ConsensusEventDispatcher struct {
	services.Service
	eng *services.Engine

	lggr           logger.Logger
	registry       *EngineRegistry
	queue          limits.QueueLimiter[v2.EnqueuedTriggerEvent]
	engineLimiters *v2.EngineLimiters
}

// NewConsensusEventDispatcher creates a dispatcher that routes events to engines.
// queue is used for the standard Wait loop; may be nil when using OCR callback mode only.
func NewConsensusEventDispatcher(
	lggr logger.Logger,
	registry *EngineRegistry,
	queue limits.QueueLimiter[v2.EnqueuedTriggerEvent],
	engineLimiters *v2.EngineLimiters,
) (*ConsensusEventDispatcher, error) {
	if registry == nil {
		return nil, fmt.Errorf("engine registry is required")
	}
	d := &ConsensusEventDispatcher{
		lggr:           lggr,
		registry:       registry,
		queue:          queue,
		engineLimiters: engineLimiters,
	}
	runLoop := d.runWaitLoop
	d.Service, d.eng = services.Config{
		Name: "ConsensusEventDispatcher",
		NewSubServices: func(lggr logger.Logger) []services.Service {
			if queue != nil {
				return []services.Service{newWaitLoopService(runLoop, lggr)}
			}
			return nil
		},
	}.NewServiceEngine(lggr)
	return d, nil
}

func (d *ConsensusEventDispatcher) runWaitLoop(ctx context.Context) error {
	for {
		event, err := d.queue.Wait(ctx)
		if err != nil {
			return err
		}
		if routeErr := d.routeEvent(ctx, event); routeErr != nil {
			d.lggr.Warnw("Failed to route consensus event", "workflowID", event.WorkflowID(), "err", routeErr)
		}
	}
}

// OnConsensusEvent implements ConsensusEventReceiver. Routes the event to the correct engine.
// Called by the OCR Transmitter when using callback mode.
func (d *ConsensusEventDispatcher) OnConsensusEvent(ctx context.Context, event v2.EnqueuedTriggerEvent) error {
	return d.routeEvent(ctx, event)
}

func (d *ConsensusEventDispatcher) routeEvent(ctx context.Context, event v2.EnqueuedTriggerEvent) error {
	wid, err := types.WorkflowIDFromHex(event.WorkflowID())
	if err != nil {
		d.lggr.Warnw("Invalid workflowID in consensus event", "workflowID", event.WorkflowID(), "err", err)
		return err
	}
	entry, ok := d.registry.Get(wid)
	if !ok {
		d.lggr.Debugw("No engine for workflow, skipping consensus event", "workflowID", event.WorkflowID())
		return nil
	}
	receiver, ok := entry.Service.(v2.ConsensusEventReceiver)
	if !ok {
		d.lggr.Warnw("Engine does not implement ConsensusEventReceiver", "workflowID", event.WorkflowID())
		return nil
	}
	return receiver.OnConsensusEvent(ctx, event)
}
