package v2

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// TriggerRegistrationDeps bundles what RegisterWorkflowTriggers needs beyond
// per-call metadata — dependencies each caller already owns in a different
// shape (Engine via EngineLimiters, TriggerCoordinator via its own fields).
//
// ChainAllowed and Settings are both nil-safe: a nil ChainAllowed skips the
// chain-access check entirely (TriggerCoordinator has none yet — the check
// moves with the limiter split, CRE-6177), and a nil Settings falls back to
// each setting's default — the same fallback an always-nil settings.Getter
// already produced before this was shared.
type TriggerRegistrationDeps struct {
	CapRegistry  core.CapabilitiesRegistry
	RegTimeout   limits.TimeLimiter
	ChainAllowed limits.GateLimiter
	Settings     settings.Getter
	Logger       logger.Logger
	Metrics      *monitoring.WorkflowsMetricLabeler
}

// TriggerRegistrationMetadata is everything that goes into every trigger's
// capabilities.RequestMetadata for one workflow.
type TriggerRegistrationMetadata struct {
	WorkflowID    string // hex-encoded workflow ID
	WorkflowOwner string
	// WorkflowName is nil-safe: if unset, both WorkflowName and
	// DecodedWorkflowName are sent empty rather than panicking on a nil
	// interface call.
	WorkflowName                  types.WorkflowName
	WorkflowTag                   string
	WorkflowDonID                 uint32
	WorkflowRegistryChainSelector string
	WorkflowRegistryAddress       string
	// OrgID is the resolved org ID, if any. Whether it's actually sent is
	// still gated by cresettings.Default.PropagateOrgIDInRequestMetadata,
	// evaluated here against deps.Settings — same as before this was shared.
	OrgID string
}

// RegisterWorkflowTriggers concurrently registers subs against their trigger
// capabilities (resolved via deps.CapRegistry), after validating each
// subscription's chain selector and, if deps.ChainAllowed is set, checking
// chain access for it. On success it returns the registered trigger
// capability IDs (in subscription order), one TriggerHandle per registration
// keyed by registration ID, and each registration's event channel (also in
// subscription order, ready for RunTriggerReader). On any registration
// failure it unregisters every handle that DID succeed and returns the first
// error — nothing is left half-registered.
//
// Shared by Engine.runTriggerSubscriptionPhase and
// TriggerCoordinator.RegisterTriggers.
func RegisterWorkflowTriggers(
	ctx context.Context,
	deps TriggerRegistrationDeps,
	meta TriggerRegistrationMetadata,
	subs []*sdkpb.TriggerSubscription,
) ([]string, map[string]*TriggerHandle, []<-chan capabilities.TriggerResponse, error) {
	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subs))
	for _, sub := range subs {
		_, labels, _ := capabilities.ParseID(sub.Id)
		chainSelector, err := capabilities.ChainSelectorLabel(labels)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("invalid chain selector for ID %s: %w", sub.Id, err)
		}
		if chainSelector != nil && deps.ChainAllowed != nil {
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
		triggers = append(triggers, triggerCap)
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
		triggerCap := triggers[i]
		g.Go(func() error {
			registrationID := TriggerRegistrationID(meta.WorkflowID, i)
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
				WorkflowDonConfigVersion:      PinnedWorkflowDonConfigVersion,
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

	handles := make(map[string]*TriggerHandle, len(subs))
	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs))
	triggerCapIDs := make([]string, len(subs))
	for result := range resultsCh {
		handles[result.registrationID] = &TriggerHandle{
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
		unregisterOnFailure(ctx, deps.Logger, meta.WorkflowID, meta.WorkflowDonID, handles)
		return nil, nil, nil, registrationErr
	}

	deps.Logger.Infow("All triggers registered successfully", "numTriggers", len(subs), "triggerIDs", triggerCapIDs)
	deps.Metrics.IncrementWorkflowRegisteredCounter(ctx)
	return triggerCapIDs, handles, eventChans, nil
}

// unregisterOnFailure rolls back every handle that registered successfully
// when at least one registration in the same batch failed. It logs (rather
// than returns) individual UnregisterTrigger failures, matching
// Engine.unregisterAllTriggers and TriggerCoordinator.unregisterAll's
// existing tolerance for a handle that fails to unregister cleanly.
func unregisterOnFailure(ctx context.Context, lggr logger.Logger, workflowID string, workflowDonID uint32, handles map[string]*TriggerHandle) {
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
		}
	}
}
