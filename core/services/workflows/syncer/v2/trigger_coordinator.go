package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// ErrWorkflowNotCoordinated is returned by UnregisterTriggers/ReleaseHandles for
// a workflowID the coordinator never registered. Reachable in the ordinary case
// where a workflow's engine type flips between activations (the ExecutionOnlyEngine
// flag changed between two creations of the same workflow) — the caller decides
// whether to call the coordinator at all based on the registry's Coordinated bit,
// so this is a defensive check, not something correctness depends on.
var ErrWorkflowNotCoordinated = errors.New("workflow not registered with the trigger coordinator")

// RegistrationParams is the per-workflow metadata RegisterTriggers needs to
// build each capability's RequestMetadata. It is a method parameter, not
// coordinator state: one coordinator instance serves every execution-only
// workflow on the node — node-level dependencies (capability registry, clock,
// engine registry) live on the coordinator; everything below is per-call.
type RegistrationParams struct {
	WorkflowOwner                 string
	WorkflowName                  string // hex-encoded
	DecodedWorkflowName           string
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowDonConfigVersion      uint32
	WorkflowRegistryChainSelector string
	WorkflowRegistryAddress       string
}

// TriggerCoordinator owns trigger registration, the trigger handle map, event
// delivery, and acknowledgement for every workflow running the execution-only
// engine (CRE-6176). It is a node-level singleton, started and stopped with the
// syncer.
//
// It never holds an engine reference: at delivery time it resolves the engine
// from the EngineRegistry by workflowID, so a workflow can be popped and
// re-added (under either engine type) without a reader holding a stale pointer.
//
// Deviation from the CRE-6176 ticket text: the ticket's "transitional enqueue
// method from CRE-6174" is Engine.Put, which stayed legacy-only in this
// implementation (see engine.go) — ExecutionEngine has no queue. So this
// coordinator delivers events by calling ExecuteTrigger directly, and does not
// yet provide the execution-concurrency cap or deadline/expiry-drop that
// Engine's Put+handleAllTriggerEvents pair provides. That gap is deliberate for
// this pass (tracked for CRE-6071's CentralTriggerQueue) — do not read its
// absence as an oversight.
type TriggerCoordinator interface {
	services.Service
	v2.Acknowledger

	// RegisterTriggers registers subs with the capability registry, retains the
	// resulting handles and reader goroutines, and returns the registered
	// trigger capability IDs. On any failure it unregisters what it already
	// registered for this call and returns the error — partial registration is
	// never left behind.
	RegisterTriggers(ctx context.Context, cre contexts.CRE, params RegistrationParams, subs []*sdkpb.TriggerSubscription) ([]string, error)

	// UnregisterTriggers stops ingress for workflowID immediately (unregisters
	// with the capability registry) but RETAINS the handle map, so an execution
	// already in flight can still resolve its handle to ACK. Call ReleaseHandles
	// once the engine has drained and closed.
	//
	// Returns ErrWorkflowNotCoordinated if workflowID was never registered here —
	// see that error's doc.
	UnregisterTriggers(workflowID string) error

	// ReleaseHandles drops the retained handle map and ACK index for workflowID.
	// Returns ErrWorkflowNotCoordinated if workflowID was never registered here.
	ReleaseHandles(workflowID string) error
}

type triggerHandle struct {
	cap     capabilities.TriggerCapability
	payload *anypb.Any
	method  string
}

type workflowTriggers struct {
	wid types.WorkflowID
	cre contexts.CRE
	// cancel stops this registration's reader goroutines. Tied to the
	// registration, not to the coordinator: without it a reader outlives its
	// own registration and can deliver into a later engine for the same
	// workflowID using this registration's trigger index. Idempotent, so both
	// UnregisterTriggers and ReleaseHandles may call it.
	cancel  context.CancelFunc
	handles map[string]*triggerHandle // registrationID -> handle
}

type triggerCoordinator struct {
	services.Service
	eng *services.Engine

	capReg   registry.CapabilitiesRegistry
	registry *EngineRegistry
	clock    clockwork.Clock
	lggr     logger.Logger

	mu        sync.Mutex
	workflows map[string]*workflowTriggers // workflowID (hex) -> state
	index     map[string]string            // registrationID -> workflowID (hex), for Ack
}

func NewTriggerCoordinator(capReg registry.CapabilitiesRegistry, engineRegistry *EngineRegistry, clock clockwork.Clock, lggr logger.Logger) TriggerCoordinator {
	c := &triggerCoordinator{
		capReg:    capReg,
		registry:  engineRegistry,
		clock:     clock,
		lggr:      logger.Named(lggr, "TriggerCoordinator"),
		workflows: make(map[string]*workflowTriggers),
		index:     make(map[string]string),
	}
	c.Service, c.eng = services.Config{
		Name:  "TriggerCoordinator",
		Start: func(context.Context) error { return nil },
		Close: func() error { return nil },
	}.NewServiceEngine(c.lggr)
	return c
}

func (c *triggerCoordinator) RegisterTriggers(ctx context.Context, cre contexts.CRE, params RegistrationParams, subs []*sdkpb.TriggerSubscription) ([]string, error) {
	workflowID := cre.Workflow
	wid, err := types.WorkflowIDFromHex(workflowID)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow id %q: %w", workflowID, err)
	}

	type registered struct {
		index int
		cap   capabilities.TriggerCapability
		ch    <-chan capabilities.TriggerResponse
		sub   *sdkpb.TriggerSubscription
	}
	results := make([]registered, 0, len(subs))

	rollback := func() {
		for _, r := range results {
			_ = r.cap.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
				TriggerID: v2.TriggerRegistrationID(workflowID, r.index),
				Metadata:  capabilities.RequestMetadata{WorkflowID: workflowID},
				Payload:   r.sub.Payload,
				Method:    r.sub.Method,
			})
		}
	}

	for i, sub := range subs {
		triggerCap, err := c.capReg.GetTrigger(ctx, sub.Id)
		if err != nil {
			rollback()
			return nil, fmt.Errorf("trigger capability not found: %w", err)
		}

		registrationID := v2.TriggerRegistrationID(workflowID, i)
		c.lggr.Infow("Registering trigger", "workflowID", workflowID, "triggerID", sub.Id, "method", sub.Method)

		eventCh, err := triggerCap.RegisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:                    workflowID,
				WorkflowOwner:                 params.WorkflowOwner,
				WorkflowName:                  params.WorkflowName,
				DecodedWorkflowName:           params.DecodedWorkflowName,
				WorkflowTag:                   params.WorkflowTag,
				WorkflowDonID:                 params.WorkflowDonID,
				WorkflowDonConfigVersion:      params.WorkflowDonConfigVersion,
				ReferenceID:                   fmt.Sprintf("trigger_%d", i),
				WorkflowRegistryChainSelector: params.WorkflowRegistryChainSelector,
				WorkflowRegistryAddress:       params.WorkflowRegistryAddress,
				EngineVersion:                 platform.ValueWorkflowVersionV2,
				OrgID:                         cre.Org,
			},
			Payload: sub.Payload,
			Method:  sub.Method,
		})
		if err != nil {
			rollback()
			return nil, fmt.Errorf("failed to register trigger %s: %w", sub.Id, err)
		}
		results = append(results, registered{index: i, cap: triggerCap, ch: eventCh, sub: sub})
	}

	readerCtx, cancel := context.WithCancel(context.Background())

	c.mu.Lock()
	// A prior registration for this workflowID should already have been torn
	// down by the syncer, but never leave its readers running if it wasn't:
	// they would race this one to deliver events to the same engine.
	if prev, ok := c.workflows[workflowID]; ok {
		prev.cancel()
	}
	wt := &workflowTriggers{wid: wid, cre: cre, cancel: cancel, handles: make(map[string]*triggerHandle, len(results))}
	c.workflows[workflowID] = wt
	triggerIDs := make([]string, len(results))
	for _, r := range results {
		registrationID := v2.TriggerRegistrationID(workflowID, r.index)
		wt.handles[registrationID] = &triggerHandle{cap: r.cap, payload: r.sub.Payload, method: r.sub.Method}
		c.index[registrationID] = workflowID
		triggerIDs[r.index] = r.sub.Id
	}
	c.mu.Unlock()

	for _, r := range results {
		c.startReader(readerCtx, wt, r.index, r.sub.Id, r.ch)
	}

	c.lggr.Infow("All triggers registered successfully", "workflowID", workflowID, "numTriggers", len(subs), "triggerIDs", triggerIDs)
	return triggerIDs, nil
}

// startReader reads one trigger's event channel and delivers to the engine by
// calling ExecuteTrigger directly (see the type doc's deviation note — this
// bypasses Engine.Put's queue, admission, and concurrency limiting, none of
// which exist on ExecutionEngine). It resolves the engine from the registry on
// every event, never caching a reference: exits quietly once the channel
// closes, the engine is no longer in the registry, or readerCtx is cancelled
// by UnregisterTriggers.
//
// readerCtx is scoped to one registration, so it does not depend on the
// capability closing its channel on unregister — which not every capability
// does. In-flight deliveries are unaffected: the ExecuteTrigger call below
// carries its own context.
func (c *triggerCoordinator) startReader(readerCtx context.Context, wt *workflowTriggers, idx int, triggerCapID string, ch <-chan capabilities.TriggerResponse) {
	c.eng.GoCtx(readerCtx, func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case event, isOpen := <-ch:
				if !isOpen {
					return
				}
				if event.Err != nil {
					c.lggr.Errorw("Received a trigger event with error, dropping", "workflowID", wt.cre.Workflow, "triggerID", triggerCapID, "err", event.Err)
					continue
				}

				svc, found := c.registry.Get(wt.wid)
				if !found {
					c.lggr.Infow("Engine gone, dropping trigger event", "workflowID", wt.cre.Workflow, "triggerID", triggerCapID)
					return
				}
				engine, ok := svc.Service.(v2.EventSink)
				if !ok {
					c.lggr.Errorw("Registry entry does not implement EventSink, dropping trigger event", "workflowID", wt.cre.Workflow, "triggerID", triggerCapID)
					continue
				}

				routed := v2.RoutedTriggerEvent{
					WorkflowID:   wt.cre.Workflow,
					TriggerCapID: triggerCapID,
					TriggerIndex: idx,
					ObservedAt:   c.clock.Now(),
					Event:        event,
				}
				deliveryCtx := contexts.WithCRE(context.Background(), wt.cre)
				if err := engine.ExecuteTrigger(deliveryCtx, routed); err != nil {
					c.lggr.Errorw("Failed to execute trigger event", "workflowID", wt.cre.Workflow, "triggerID", triggerCapID, "err", err)
				}
			}
		}
	})
}

func (c *triggerCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	c.mu.Lock()
	workflowID, ok := c.index[triggerRegistrationID]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("failed to find workflow for registration %s", triggerRegistrationID)
	}
	wt, ok := c.workflows[workflowID]
	if !ok {
		c.mu.Unlock()
		return fmt.Errorf("failed to find workflow %s", workflowID)
	}
	handle, ok := wt.handles[triggerRegistrationID]
	c.mu.Unlock()
	if !ok {
		return fmt.Errorf("failed to find trigger handle %s", triggerRegistrationID)
	}
	return handle.cap.AckEvent(ctx, triggerRegistrationID, eventID, handle.method)
}

func (c *triggerCoordinator) UnregisterTriggers(workflowID string) error {
	ctx := context.Background()

	c.mu.Lock()
	wt, ok := c.workflows[workflowID]
	c.mu.Unlock()
	if !ok {
		return ErrWorkflowNotCoordinated
	}

	// Stop the readers before unregistering with the capabilities: once this
	// returns, no further event can be delivered for this workflow, whether or
	// not the capability closes its channel.
	wt.cancel()

	failCount := 0
	for registrationID, handle := range wt.handles {
		if err := handle.cap.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata:  capabilities.RequestMetadata{WorkflowID: workflowID},
			Payload:   handle.payload,
			Method:    handle.method,
		}); err != nil {
			c.lggr.Errorw("Failed to unregister trigger", "workflowID", workflowID, "registrationID", registrationID, "err", err)
			failCount++
		}
	}
	c.lggr.Infow("Unregistered triggers, retaining handles for in-flight executions", "workflowID", workflowID, "numTriggers", len(wt.handles), "failed", failCount)
	// Handles are intentionally retained here (AC 6): an execution already in
	// flight still needs to resolve its handle to ACK. The caller releases them
	// via ReleaseHandles once the engine has drained and closed.
	return nil
}

func (c *triggerCoordinator) ReleaseHandles(workflowID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	wt, ok := c.workflows[workflowID]
	if !ok {
		return ErrWorkflowNotCoordinated
	}
	// Normally already cancelled by UnregisterTriggers; repeated here so a
	// direct ReleaseHandles cannot strand the readers.
	wt.cancel()
	for registrationID := range wt.handles {
		delete(c.index, registrationID)
	}
	delete(c.workflows, workflowID)
	return nil
}
