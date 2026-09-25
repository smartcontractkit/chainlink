package triggers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const _prefix = "trigger_reg"

// Handle is a registered trigger capability plus the registration
// payload/method needed to unregister and re-deliver to it.
type Handle struct {
	capabilities.TriggerCapability
	Payload *anypb.Any
	Method  string
}

// Ack acknowledges a trigger event against the given handle, logging and
// bumping the same success/failure metrics regardless of caller.
// handle may be nil (i.e., registration not found).
func Ack(
	ctx context.Context,
	lggr logger.Logger,
	metrics *monitoring.WorkflowsMetricLabeler,
	triggerCapID, triggerRegistrationID, eventID string,
	handle *Handle,
) error {
	lggr.Infow("ACKing trigger event", "triggerRegistrationID", triggerRegistrationID, "eventID", eventID)

	tm := metrics.With(platform.KeyTriggerID, triggerCapID)

	if handle == nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("failed to find trigger %s", triggerRegistrationID)
	}
	if err := handle.AckEvent(ctx, triggerRegistrationID, eventID, handle.Method); err != nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return err
	}
	tm.IncrementTriggerEventAckSuccessCounter(ctx)
	return nil
}

// CoordinatedEvent is the canonical trigger event type that flows
// through the coordinator path into the engine.
type CoordinatedEvent struct {
	WorkflowID   string
	TriggerCapID string
	TriggerIndex int

	// ObservedAt is the time the CoordinatedEvent was constructed.
	// It is used for skew metrics (queue wait time)
	// and deadline enforcement.
	ObservedAt time.Time

	// Deadline is the expiry of this event,
	// stamped once at dispatch as ObservedAt + TriggerEventQueueTimeout.
	// A settings change after dispatch does not affect already-queued events
	Deadline time.Time

	// SequenceNumber determines the execution order of trigger events across the DON.
	SequenceNumber uint64
	Event          capabilities.TriggerResponse
}

// ReadLoop consumes triggerEventCh until it closes or ctx is done,
// converting each received capabilities.TriggerResponse into a
// CoordinatedEvent and handing it to deliver.
//
// The caller owns spawning the goroutine, since goroutine lifecycle
// differs between callers.
func ReadLoop(
	ctx context.Context,
	lggr logger.Logger,
	metrics *monitoring.WorkflowsMetricLabeler,
	clock clockwork.Clock,
	workflowID, triggerCapID string,
	triggerIndex int,
	triggerEventCh <-chan capabilities.TriggerResponse,
	deliver func(ctx context.Context, event CoordinatedEvent),
) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, isOpen := <-triggerEventCh:
			if !isOpen {
				return
			}
			eventID := event.Event.ID
			metrics.With(platform.KeyTriggerID, triggerCapID).IncrementTriggerEventReceivedCounter(ctx)
			lggr.Debugw("Processing trigger event", "triggerID", triggerCapID, "eventID", eventID)
			if event.Err != nil {
				lggr.Errorw("Received a trigger event with error, dropping", "triggerID", triggerCapID, "err", event.Err)
				tm := metrics.With(platform.KeyTriggerID, triggerCapID)
				tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
				tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonTriggerResponseError)
				continue
			}

			routed := CoordinatedEvent{
				WorkflowID:   workflowID,
				TriggerCapID: triggerCapID,
				TriggerIndex: triggerIndex,
				ObservedAt:   clock.Now(),
				Event:        event,
			}

			deliver(ctx, routed)
		}
	}
}

// RegisterDeps bundles what Register needs beyond per-call metadata.
type RegisterDeps struct {
	CapRegistry  registry.CapabilitiesRegistry
	RegTimeout   limits.TimeLimiter
	ChainAllowed limits.GateLimiter
	Settings     settings.Getter
	Logger       logger.Logger
	Metrics      *monitoring.WorkflowsMetricLabeler
}

// RegisterMetadata is everything that goes into every trigger's
// capabilities.RequestMetadata for one workflow.
type RegisterMetadata struct {
	WorkflowID    string // hex-encoded workflow ID
	WorkflowOwner string
	// WorkflowName is nil-safe: if unset, both WorkflowName and
	// DecodedWorkflowName are sent empty rather than panicking on a nil
	// interface call.
	WorkflowName  types.WorkflowName
	WorkflowTag   string
	WorkflowDonID uint32
	// WorkflowDonConfigVersion is the caller's pinned DON config version
	// (e.g. engine.go's pinnedWorkflowDonConfigVersion constant). triggers
	// does not own this value since it's also used for unrelated things in
	// v2 (buildLabels, localNodeSync, capability_executor.go).
	WorkflowDonConfigVersion      uint32
	WorkflowRegistryChainSelector string
	WorkflowRegistryAddress       string
	// OrgID is the resolved org ID.
	OrgID string
}

// Register concurrently registers subs against their trigger capabilities
// (resolved via deps.CapRegistry), after validating each subscription's
// chain selector and checking chain access for it via deps.ChainAllowed.
//
// On success it returns the registered trigger capability IDs (in
// subscription order), one Handle per registration keyed by registration ID,
// and each registration's event channel.
//
// On any registration failure it unregisters every handle that did succeed
// and returns the first error.
func Register(
	ctx context.Context,
	deps RegisterDeps,
	meta RegisterMetadata,
	subs []*sdkpb.TriggerSubscription,
) ([]string, map[string]*Handle, []<-chan capabilities.TriggerResponse, error) {
	// check if all requested triggers exist in the registry
	tcs := make([]capabilities.TriggerCapability, 0, len(subs))
	for _, sub := range subs {
		_, labels, _ := capabilities.ParseID(sub.Id)
		chainSelector, err := capabilities.ChainSelectorLabel(labels)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid chain selector for ID %s: %w", sub.Id, err)
		}
		if chainSelector != nil {
			if err := deps.ChainAllowed.AllowErr(contexts.WithChainSelector(ctx, *chainSelector)); err != nil {
				if errors.Is(err, limits.ErrorNotAllowed{}) {
					return nil, nil, nil, fmt.Errorf("unable to subscribe to capability %s: ChainSelector %d: %w", sub.Id, *chainSelector, err)
				}
				return nil, nil, nil, fmt.Errorf("failed to check access for ChainSelector %d: %w", *chainSelector, err)
			}
		}
		triggerCap, err := deps.CapRegistry.GetTrigger(ctx, sub.Id)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("trigger capability not found: %w", err)
		}
		tcs = append(tcs, triggerCap)
	}

	// register to all triggers concurrently
	regCtx, regCancel, err := deps.RegTimeout.WithTimeout(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	defer regCancel()

	// trigger registration results for use in concurrent trigger subscriptions
	type triggerRegResult struct {
		index          int
		registrationID string
		triggerCap     capabilities.TriggerCapability
		eventCh        <-chan capabilities.TriggerResponse
		payload        *anypb.Any
		method         string
		triggerCapID   string
	}

	resultsCh := make(chan triggerRegResult, len(subs))
	g, gCtx := errgroup.WithContext(regCtx)

	var workflowNameHex, decodedWorkflowName string
	if meta.WorkflowName != nil {
		workflowNameHex = meta.WorkflowName.Hex()
		decodedWorkflowName = meta.WorkflowName.String()
	}

	// Launch concurrent trigger registrations
	for i, sub := range subs {
		triggerCap := tcs[i]
		g.Go(func() error {
			registrationID := RegistrationID(meta.WorkflowID, i)
			args := []any{"triggerID", sub.Id, "method", sub.Method}
			if sub.Payload != nil {
				args = append(args, "payload", protojson.Format(sub.Payload))
			}
			deps.Logger.Infow("Registering trigger", args...)
			metadata := capabilities.RequestMetadata{
				WorkflowID:                    meta.WorkflowID,
				WorkflowOwner:                 meta.WorkflowOwner,
				WorkflowName:                  workflowNameHex,
				DecodedWorkflowName:           decodedWorkflowName,
				WorkflowTag:                   meta.WorkflowTag,
				WorkflowDonID:                 meta.WorkflowDonID,
				WorkflowDonConfigVersion:      meta.WorkflowDonConfigVersion,
				ReferenceID:                   fmt.Sprintf("trigger_%d", i),
				WorkflowRegistryChainSelector: meta.WorkflowRegistryChainSelector,
				WorkflowRegistryAddress:       meta.WorkflowRegistryAddress,
				EngineVersion:                 platform.ValueWorkflowVersionV2,
				// no WorkflowExecutionID needed (or available at this stage)
			}
			propagateOrgIDMeta, _ := cresettings.Default.PropagateOrgIDInRequestMetadata.GetOrDefault(gCtx, deps.Settings)
			if propagateOrgIDMeta && meta.OrgID != "" {
				metadata.OrgID = meta.OrgID
			}
			triggerEventCh, regErr := triggerCap.RegisterTrigger(gCtx, capabilities.TriggerRegistrationRequest{
				TriggerID: registrationID,
				Metadata:  metadata,
				Payload:   sub.Payload,
				Method:    sub.Method,
				// no Config needed - NoDAG uses Payload
			})
			if regErr != nil {
				deps.Logger.Errorw("Trigger registration failed", "triggerID", sub.Id, "err", regErr)
				deps.Metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(gCtx)
				return fmt.Errorf("failed to register trigger %s: %w", sub.Id, regErr)
			}
			// Send successful result
			resultsCh <- triggerRegResult{
				index:          i,
				registrationID: registrationID,
				triggerCap:     triggerCap,
				eventCh:        triggerEventCh,
				payload:        sub.Payload,
				method:         sub.Method,
				triggerCapID:   sub.Id,
			}
			return nil
		})
	}

	// wait for all registrations to complete.
	// returns first non-nil error.
	registrationErr := g.Wait()
	close(resultsCh)

	handles := make(map[string]*Handle, len(subs))
	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs))
	triggerCapIDs := make([]string, len(subs))
	for result := range resultsCh {
		handles[result.registrationID] = &Handle{
			TriggerCapability: result.triggerCap,
			Payload:           result.payload,
			Method:            result.method,
		}
		eventChans[result.index] = result.eventCh
		triggerCapIDs[result.index] = result.triggerCapID
	}

	// If any registration failed, unregister successful ones and return error
	if registrationErr != nil {
		deps.Logger.Errorw("One or more trigger registrations failed - reverting all", "err", registrationErr)
		Unregister(ctx, deps.Logger, meta.WorkflowID, meta.WorkflowDonID, handles)
		return nil, nil, nil, registrationErr
	}

	deps.Logger.Infow("All triggers registered successfully", "numTriggers", len(subs), "triggerIDs", triggerCapIDs)
	deps.Metrics.IncrementWorkflowRegisteredCounter(ctx)
	return triggerCapIDs, handles, eventChans, nil
}

// Unregister unregisters every handle in handles with the capability
// registry. It logs individual UnregisterTrigger failures and returns how
// many occurred, so callers can log their own summary.
func Unregister(ctx context.Context, lggr logger.Logger, workflowID string, workflowDonID uint32, handles map[string]*Handle) (failCount int) {
	for registrationID, handle := range handles {
		if err := handle.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    workflowID,
				WorkflowDonID: workflowDonID,
			},
			Payload: handle.Payload,
			Method:  handle.Method,
		}); err != nil {
			lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
			failCount++
		}
	}
	return failCount
}

// RegistrationID constructs a trigger registration ID from a workflow ID and trigger index.
func RegistrationID(workflowID string, triggerIndex int) string {
	return fmt.Sprintf("%s_%s_%d", _prefix, workflowID, triggerIndex)
}

// ParseWorkflowID extracts the workflow ID from a registration ID produced by RegistrationID.
func ParseWorkflowID(registrationID string) (string, error) {
	rest, ok := strings.CutPrefix(registrationID, _prefix+"_")
	if !ok {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing prefix %q", registrationID, _prefix)
	}

	idx := strings.LastIndex(rest, "_")
	if idx == -1 {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing trigger index", registrationID)
	}

	workflowID, triggerIndexStr := rest[:idx], rest[idx+1:]
	if _, err := strconv.Atoi(triggerIndexStr); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid trigger index %q: %w", registrationID, triggerIndexStr, err)
	}

	if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid workflow ID %q: %w", registrationID, workflowID, err)
	}

	return workflowID, nil
}
