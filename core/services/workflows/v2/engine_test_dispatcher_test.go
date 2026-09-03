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

// testDispatcher plays the role the TriggerDispatcher plays in production for
// engine unit tests: it registers the engine's subscriptions with the
// capability registry (satisfying the tests' GetTrigger/RegisterTrigger mock
// expectations), reads each trigger's event channel into engine.Put, forwards
// the engine's Ack calls to the trigger capability (satisfying AckEvent
// expectations), and unregisters on test cleanup.
//
// It exists because the engine no longer registers triggers, reads event
// channels, or acknowledges by itself — those responsibilities moved to the
// syncer-owned dispatcher, which these tests do not construct.
type testDispatcher struct {
	t      *testing.T
	cfg    *v2.EngineConfig
	capReg core.CapabilitiesRegistry
	clock  clockwork.Clock
	lggr   logger.Logger

	// engine is set before Start: readers resolve it per event, mirroring the
	// production dispatcher's registry lookup at delivery time.
	engine atomic.Pointer[v2.Engine]

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

var _ v2.Acknowledger = (*testDispatcher)(nil)

// newDispatchedEngine wires a testDispatcher into cfg (as its acknowledger and
// subscriptions hook) and constructs the engine. The caller still calls
// engine.Start — registration happens when the engine fires
// OnSubscriptionsReady during init, exactly as in production.
func newDispatchedEngine(t *testing.T, cfg *v2.EngineConfig) *v2.Engine {
	t.Helper()
	clock := cfg.Clock
	if clock == nil {
		clock = clockwork.NewRealClock()
	}
	// Wrap cfg.Lggr through the same beholder logger the engine builds
	// internally, so the dispatcher's "Registering trigger"/"All triggers
	// registered successfully" logs reach beholder observers in tests that
	// assert on those messages, exactly as the production dispatcher's logs
	// do (via its own beholder-wired lggr).
	dispatcherLggr := logger.Logger(nil)
	if cfg.Lggr != nil {
		dispatcherLggr = logger.Sugared(custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter))
	}
	d := &testDispatcher{
		t:        t,
		cfg:      cfg,
		capReg:   cfg.CapRegistry,
		clock:    clock,
		lggr:     dispatcherLggr,
		triggers: make(map[string]capabilities.TriggerCapability),
		handles:  make(map[string]*testHandle),
		methods:  make(map[string]string),
	}

	// The dispatcher is the engine's acknowledger, unless a test has already
	// wired its own in place of defaultTestConfig's noopAcknowledger sentinel
	// (e.g. a recordingAcknowledger for a unit test that drives ExecuteTrigger
	// directly, without going through registration).
	if _, isDefault := cfg.TriggerAcknowledger.(noopAcknowledger); cfg.TriggerAcknowledger == nil || isDefault {
		cfg.TriggerAcknowledger = d
	}

	// Set the subscriptions hook: register synchronously so registration
	// errors propagate to OnInitialized, matching the engine's previous
	// behavior (and the tests' error-propagation assertions). After
	// registration, fire OnSubscribedToTriggers with the registered trigger
	// capability IDs — in production the syncer does this via
	// Engine.FireOnSubscribedToTriggers after the dispatcher registers.
	//
	// OnSubscriptionsReady is only ever set by this function, never directly
	// by a test, so it is overwritten rather than composed: several tests
	// share one cfg across t.Run subtests, each calling newDispatchedEngine
	// again, and composing onto whatever the previous subtest's call left
	// behind would re-invoke that stale dispatcher's registration against
	// already-exhausted mock expectations. OnSubscribedToTriggers, in
	// contrast, is set directly by tests (constant across subtests sharing a
	// cfg) and is preserved.
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

	engine, err := v2.NewEngine(cfg)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}
	d.engine.Store(engine)

	return engine
}

// unregisterOnCleanup registers a t.Cleanup that unregisters every handle
// register has stored, so the tests' UnregisterTrigger expectations are met.
// It is called from register, not from newDispatchedEngine: t.Cleanup runs in
// LIFO order, and register only runs once the engine starts — by which point
// every test has already constructed and configured the trigger mocks it
// depends on (Start needs them wired first). Registering here, rather than at
// newDispatchedEngine's call time, therefore guarantees this cleanup runs
// after those mocks' own construction-time AssertExpectations cleanup was
// registered, so LIFO runs ours first, regardless of whether a given test
// happens to create its mocks before or after calling newDispatchedEngine.
func (d *testDispatcher) unregisterOnCleanup() {
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
func (d *testDispatcher) register(subs []*sdkpb.TriggerSubscription) error {
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

// startReader reads one trigger's event channel and feeds the engine through
// its transitional Put method, mirroring the production dispatcher's readers.
// The reader outlives the registering test context (production readers use
// context.WithoutCancel); it exits when the channel closes or the engine is
// closed and drained.
func (d *testDispatcher) startReader(idx int, sub *sdkpb.TriggerSubscription, ch <-chan capabilities.TriggerResponse) {
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
				// Mirror the production dispatcher's reader: events carrying a
				// trigger-response error are dropped before delivery (the old
				// engine reader did the same before this responsibility moved).
				if event.Err != nil {
					continue
				}
				engine := d.engine.Load()
				if engine == nil {
					return
				}
				// The engine's scoped limiters (queue, execution
				// concurrency) key on the workflow tenant, so the delivery
				// context must carry it — in production the syncer's context
				// already does.
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
				if err := engine.Put(putCtx, routed); err != nil {
					// Draining is expected during workflow deletion.
					d.t.Logf("Put failed for event %s: %v", event.Event.ID, err)
				}
			}
		}
	}()
}

// registeredTriggerIDs returns the trigger capability IDs for the given
// subscriptions, in subscription order — the IDs RegisterTriggers returns in
// production.
func (d *testDispatcher) registeredTriggerIDs(subs []*sdkpb.TriggerSubscription) []string {
	ids := make([]string, len(subs))
	for i, sub := range subs {
		ids[i] = sub.Id
	}
	return ids
}

// Ack forwards the engine's acknowledgement to the trigger capability.
func (d *testDispatcher) Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error {
	d.mu.Lock()
	cap, ok := d.triggers[triggerCapID]
	method := d.methods[triggerCapID]
	d.mu.Unlock()
	if !ok {
		return fmt.Errorf("failed to find trigger %s", triggerCapID)
	}
	return cap.AckEvent(ctx, triggerRegistrationID, eventID, method)
}
