package v2_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jonboulle/clockwork"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// testCoordinator is a test fixture that performs the trigger-lifecycle work the
// engine no longer does itself: it registers the engine's subscriptions with the
// capability registry, reads each trigger's event channel into engine.Put,
// forwards the engine's Ack calls to the trigger capability, and unregisters on
// test cleanup.
type testCoordinator struct {
	t      *testing.T
	cfg    *v2.EngineConfig
	capReg core.CapabilitiesRegistry
	clock  clockwork.Clock
	lggr   logger.Logger

	// engine is set before Start; readers resolve it per event.
	engine atomic.Pointer[v2.WorkflowEngine]

	mu       sync.Mutex
	triggers map[string]capabilities.TriggerCapability // triggerCapID -> capability (Ack)
	handles  map[string]*testHandle                    // registrationID -> handle (unregister)
	methods  map[string]string                         // triggerCapID -> method (Ack)
}

// testHandle is a registered trigger capability plus the registration
// payload/method needed to unregister it.
type testHandle struct {
	cap     capabilities.TriggerCapability
	payload *anypb.Any
	method  string
}

var _ v2.Acknowledger = (*testCoordinator)(nil)

// noopAcknowledger is the default no-op acknowledger. It lets newCoordinatedEngine
// distinguish "the test left the default in place" from "the test wired its own
// acknowledger" when a cfg is shared across subtests.
type noopAcknowledger struct{}

func (noopAcknowledger) Ack(context.Context, string, string, string) error { return nil }

// engineCtor constructs a workflow engine from a config. It matches the
// signatures of v2.NewEngine and v2.NewExecutionEngine.
type engineCtor func(*v2.EngineConfig) (v2.WorkflowEngine, error)

// newCoordinatedEngine wires a testCoordinator into cfg and constructs the engine
// via ctor. The caller still calls engine.Start.
//
// The fixture's role depends on the engine type, selected by executionOnly:
//   - executionOnly=true (ExecutionEngine): the fixture registers the engine's
//     subscriptions (via the OnSubscriptionsReady hook), reads each trigger's
//     event channel into engine.Put, and serves as the engine's acknowledger.
//   - executionOnly=false (Engine): the engine self-registers and
//     self-acknowledges, so the fixture stays out of the way — it leaves
//     TriggerAcknowledger nil (Engine injects itself) and does not register on
//     OnSubscriptionsReady.
func newCoordinatedEngine(t *testing.T, cfg *v2.EngineConfig, ctor engineCtor, executionOnly bool) v2.WorkflowEngine {
	t.Helper()
	clock := cfg.Clock
	if clock == nil {
		clock = clockwork.NewRealClock()
	}
	// Wrap cfg.Lggr through the same beholder logger the engine builds
	// internally, so the fixture's registration logs reach the beholder
	// observers that tests assert on.
	coordinatorLggr := logger.Logger(nil)
	if cfg.Lggr != nil {
		coordinatorLggr = logger.Sugared(custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter))
	}
	d := &testCoordinator{
		t:        t,
		cfg:      cfg,
		capReg:   cfg.CapRegistry,
		clock:    clock,
		lggr:     coordinatorLggr,
		triggers: make(map[string]capabilities.TriggerCapability),
		handles:  make(map[string]*testHandle),
		methods:  make(map[string]string),
	}

	if executionOnly {
		// The fixture is the engine's acknowledger, unless a test already wired its own.
		if _, isDefault := cfg.TriggerAcknowledger.(noopAcknowledger); cfg.TriggerAcknowledger == nil || isDefault {
			cfg.TriggerAcknowledger = d
		}

		// Register synchronously so registration errors propagate to OnInitialized,
		// then fire OnSubscribedToTriggers with the registered trigger IDs.
		//
		// OnSubscriptionsReady is overwritten, not composed: tests share one cfg
		// across subtests, and composing would re-invoke a stale fixture's
		// registration against exhausted mock expectations. OnSubscribedToTriggers
		// is set directly by tests and is preserved.
		existingSubscribed := cfg.Hooks.OnSubscribedToTriggers
		cfg.Hooks.OnSubscriptionsReady = func(subs []*sdkpb.TriggerSubscription, cre contexts.CRE) error {
			if err := d.register(subs); err != nil {
				return err
			}
			if existingSubscribed != nil {
				existingSubscribed(d.registeredTriggerIDs(subs))
			}
			return nil
		}
	}

	engine, err := ctor(cfg)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}
	d.engine.Store(&engine)

	return engine
}

// unregisterOnCleanup unregisters every stored handle via t.Cleanup so the
// tests' UnregisterTrigger expectations are met. It is registered from register
// (not newCoordinatedEngine) so that, under LIFO cleanup ordering, it runs before
// the trigger mocks' own AssertExpectations cleanup.
func (d *testCoordinator) unregisterOnCleanup() {
	d.t.Cleanup(func() {
		d.mu.Lock()
		handles := d.handles
		d.handles = make(map[string]*testHandle)
		d.mu.Unlock()
		ctx := context.Background()
		for registrationID, h := range handles {
			if unregErr := h.cap.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
				TriggerID: registrationID,
				Metadata: capabilities.RequestMetadata{
					WorkflowID: d.cfg.WorkflowID,
				},
				Payload: h.payload,
				Method:  h.method,
			}); unregErr != nil {
				d.t.Logf("failed to unregister trigger %s: %v", registrationID, unregErr)
			}
		}
	})
}

// register validates and registers each subscription with the capability
// registry, stores the handles, and starts one reader goroutine per
// subscription. On any failure it rolls back successful registrations.
func (d *testCoordinator) register(subs []*sdkpb.TriggerSubscription) error {
	ctx := context.Background()
	type regResult struct {
		index int
		cap   capabilities.TriggerCapability
		ch    <-chan capabilities.TriggerResponse
		sub   *sdkpb.TriggerSubscription
	}
	results := make([]regResult, 0, len(subs))
	var regErr error
	for i, sub := range subs {
		triggerCap, err := d.capReg.GetTrigger(ctx, sub.Id)
		if err != nil {
			regErr = fmt.Errorf("trigger capability not found: %w", err)
			break
		}
		registrationID := v2.TriggerRegistrationID(d.cfg.WorkflowID, i)
		if d.lggr != nil {
			args := []any{"triggerID", sub.Id, "method", sub.Method}
			if sub.Payload != nil {
				args = append(args, "payload", protojson.Format(sub.Payload))
			}
			d.lggr.Infow("Registering trigger", args...)
		}
		metadata := capabilities.RequestMetadata{
			WorkflowID:                    d.cfg.WorkflowID,
			WorkflowOwner:                 d.cfg.WorkflowOwner,
			WorkflowName:                  d.cfg.WorkflowName.Hex(),
			DecodedWorkflowName:           d.cfg.WorkflowName.String(),
			WorkflowTag:                   d.cfg.WorkflowTag,
			WorkflowDonID:                 0, // tests use a zero-ID DON
			WorkflowDonConfigVersion:      1,
			ReferenceID:                   fmt.Sprintf("trigger_%d", i),
			WorkflowRegistryChainSelector: d.cfg.WorkflowRegistryChainSelector,
			WorkflowRegistryAddress:       d.cfg.WorkflowRegistryAddress,
			EngineVersion:                 platform.ValueWorkflowVersionV2,
		}
		var creGetter settings.Getter
		if d.cfg.LocalLimiters != nil {
			creGetter = d.cfg.LocalLimiters.Settings
		}
		propagateOrgIDMeta, _ := cresettings.Default.PropagateOrgIDInRequestMetadata.GetOrDefault(ctx, creGetter)
		if propagateOrgIDMeta && d.cfg.OrgResolver != nil {
			if orgID, orgErr := d.cfg.OrgResolver.Get(ctx, d.cfg.WorkflowOwner); orgErr == nil && orgID != "" {
				metadata.OrgID = orgID
			}
		}
		eventCh, err := triggerCap.RegisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata:  metadata,
			Payload:   sub.Payload,
			Method:    sub.Method,
		})
		if err != nil {
			regErr = fmt.Errorf("failed to register trigger %s: %w", sub.Id, err)
			break
		}
		results = append(results, regResult{index: i, cap: triggerCap, ch: eventCh, sub: sub})
	}

	if regErr != nil {
		// Roll back successful registrations, mirroring the previous engine
		// behavior the tests assert ("failed trigger registration and rollback").
		for _, r := range results {
			_ = r.cap.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
				TriggerID: v2.TriggerRegistrationID(d.cfg.WorkflowID, r.index),
				Metadata: capabilities.RequestMetadata{
					WorkflowID: d.cfg.WorkflowID,
				},
				Payload: r.sub.Payload,
				Method:  r.sub.Method,
			})
		}
		return regErr
	}

	d.mu.Lock()
	for _, r := range results {
		registrationID := v2.TriggerRegistrationID(d.cfg.WorkflowID, r.index)
		d.triggers[r.sub.Id] = r.cap
		d.methods[r.sub.Id] = r.sub.Method
		d.handles[registrationID] = &testHandle{cap: r.cap, payload: r.sub.Payload, method: r.sub.Method}
	}
	d.mu.Unlock()

	d.unregisterOnCleanup()

	for _, r := range results {
		d.startReader(r.index, r.sub, r.ch)
	}
	if d.lggr != nil {
		triggerCapIDs := make([]string, len(results))
		for _, r := range results {
			triggerCapIDs[r.index] = r.sub.Id
		}
		d.lggr.Infow("All triggers registered successfully", "numTriggers", len(subs), "triggerIDs", triggerCapIDs)
	}
	return nil
}

// startReader reads one trigger's event channel and feeds the engine through its
// transitional Put method. The reader outlives the registering test context; it
// exits when the channel closes or the engine is closed and drained.
func (d *testCoordinator) startReader(idx int, sub *sdkpb.TriggerSubscription, ch <-chan capabilities.TriggerResponse) {
	ctx := context.WithoutCancel(d.t.Context())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, isOpen := <-ch:
				if !isOpen {
					return
				}
				// Events carrying a trigger-response error are dropped before delivery.
				if event.Err != nil {
					continue
				}
				engine := d.engine.Load()
				if engine == nil {
					return
				}
				// The engine's scoped limiters (queue, execution concurrency) key
				// on the workflow tenant, so the delivery context must carry it.
				putCtx := contexts.WithCRE(ctx, contexts.CRE{
					Owner:    d.cfg.WorkflowOwner,
					Workflow: d.cfg.WorkflowID,
				})
				routed := v2.RoutedTriggerEvent{
					WorkflowID:   d.cfg.WorkflowID,
					TriggerCapID: sub.Id,
					TriggerIndex: idx,
					ObservedAt:   d.clock.Now(),
					Event:        event,
				}
				// Put is transitional and not on the WorkflowEngine interface; the
				// fixture only reads event channels for ExecutionEngine, which has it.
				putter, ok := (*engine).(interface {
					Put(context.Context, v2.RoutedTriggerEvent) error
				})
				if !ok {
					return
				}
				if err := putter.Put(putCtx, routed); err != nil {
					// Draining is expected during workflow deletion.
					d.t.Logf("Put failed for event %s: %v", event.Event.ID, err)
				}
			}
		}
	}()
}

// registeredTriggerIDs returns the trigger capability IDs for the given
// subscriptions, in subscription order.
func (d *testCoordinator) registeredTriggerIDs(subs []*sdkpb.TriggerSubscription) []string {
	ids := make([]string, len(subs))
	for i, sub := range subs {
		ids[i] = sub.Id
	}
	return ids
}

// Ack forwards the engine's acknowledgement to the trigger capability.
func (d *testCoordinator) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	d.mu.Lock()
	triggerCap, ok := d.triggers[triggerCapID]
	method := d.methods[triggerCapID]
	d.mu.Unlock()
	if !ok {
		return fmt.Errorf("failed to find trigger %s", triggerCapID)
	}
	return triggerCap.AckEvent(ctx, triggerRegistrationID, eventID, method)
}
