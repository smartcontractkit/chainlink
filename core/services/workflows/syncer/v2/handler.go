package v2

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/resourcemanager"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	generichost "github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

type ORM interface {
	artifacts.WorkflowSpecsDS
}

// engineFactoryFn creates a workflow engine. The initDone channel is used to signal when the engine
// has completed initialization (including trigger subscriptions). For v2 engines, this is wired to
// the OnInitialized lifecycle hook. For v1 legacy DAG engines, nil is sent immediately after engine
// creation since they don't support async initialization hooks.
type engineFactoryFn func(ctx context.Context, wfid, owner string, name types.WorkflowName, tag string, config, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error)

type DrainableService interface {
	Drain() bool
	ActiveExecutions() int32
	DrainStartedAt() (time.Time, bool)
}

var ErrDrainInProgress = errors.New("drain in progress")

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding method that handles the event.
type eventHandler struct {
	services.Service
	eng *services.Engine

	lggr logger.Logger

	workflowStore          store.Store
	capRegistry            core.CapabilitiesRegistry
	executionHandlers      *confidentialrelay.ExecutionHandlers
	donTimeStore           *dontime.Store
	useLocalTimeProvider   bool
	engineRegistry         *EngineRegistry
	emitter                custmsg.MessageEmitter
	emitterMu              sync.RWMutex
	engineFactory          engineFactoryFn
	engineLimiters         *v2.EngineLimiters
	featureFlags           *v2.EngineFeatureFlags
	ratelimiter            *ratelimiter.RateLimiter
	workflowLimits         limits.ResourceLimiter[int]
	workflowArtifactsStore WorkflowArtifactsStore
	workflowEncryptionKey  workflowkey.Key
	workflowDonSubscriber  capabilities.DonSubscriber
	billingClient          metering.BillingClient

	// specMeter owns all metering for durable workflow-spec storage (the
	// ResourceManager lifecycle, identity, snapshots). Nil when metering is
	// disabled: its handler-facing methods are nil-receiver safe no-ops.
	specMeter *SpecMeter

	orgResolver    orgresolver.OrgResolver
	secretsFetcher v2.SecretsFetcher
	// localSecretOverrides is keyed by owner address; values are secret id -> secret value
	localSecretOverrides map[string]map[string]string

	moduleLRU           *ModuleLRU
	moduleStore         artifacts.SerialisedModuleStore
	cacheMetrics        *CacheMetrics
	moduleEngineVersion string

	// WorkflowRegistryAddress is the address of the workflow registry contract
	workflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	workflowRegistryChainSelector string

	// debugMode enables additional OTel tracing for workflow engines and syncer.
	// When enabled, traces are created for workflow execution and syncer events.
	debugMode bool

	// tracer is the OTel tracer for this handler. It's a noop tracer when debug mode is disabled.
	tracer trace.Tracer

	shardOrchestratorClient shardorchestrator.ClientInterface
	shardingEnabled         bool
	myShardID               uint32
	shardRoutingSteady      *shardownership.SteadySignal
	shardResolver           shardownership.ShardResolver

	metrics *metrics
}

// EventHandlerOption is a functional option for configuring an eventHandler.
type EventHandlerOption = func(*eventHandler)

func WithEngineRegistry(er *EngineRegistry) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineRegistry = er
	}
}

// WithEngineFactoryFn allows for overriding the engine factory function.
// if in doubt, close initDone channel immediately in tests to prevent deadlocks.
func WithEngineFactoryFn(efn engineFactoryFn) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineFactory = efn
	}
}

func WithStaticEngine(engine services.Service) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineFactory = func(_ context.Context, _, _ string, _ types.WorkflowName, _ string, _, _ []byte, _ string, initDone chan<- error) (services.Service, error) {
			// For static engines (used in tests), signal immediate initialization success
			if initDone != nil {
				initDone <- nil
			}
			return engine, nil
		}
	}
}

func WithBillingClient(client metering.BillingClient) func(*eventHandler) {
	return func(e *eventHandler) {
		e.billingClient = client
	}
}

// WithSpecMeter supplies the SpecMeter that emits metering.v1.MeterRecord
// events for the workflow_specs_v2 storage lifecycle. The handler runs it as a
// sub-service and reports storage transitions through EmitSpecDelta; all other
// metering concerns (ResourceManager lifecycle, identity, snapshots) live on
// the meter. A nil meter (metering disabled) is a valid no-op.
func WithSpecMeter(sm *SpecMeter) func(*eventHandler) {
	return func(e *eventHandler) {
		e.specMeter = sm
	}
}

func WithShardExecutionGuard(client shardorchestrator.ClientInterface, shardingEnabled bool, shardID uint32) func(*eventHandler) {
	return func(e *eventHandler) {
		e.shardOrchestratorClient = client
		e.shardingEnabled = shardingEnabled
		e.myShardID = shardID
	}
}

func WithShardRoutingSteady(signal *shardownership.SteadySignal) func(*eventHandler) {
	return func(e *eventHandler) {
		e.shardRoutingSteady = signal
	}
}

func WithShardResolver(resolver shardownership.ShardResolver) func(*eventHandler) {
	return func(e *eventHandler) {
		e.shardResolver = resolver
	}
}

func WithWorkflowRegistry(address, chainSelector string) func(*eventHandler) {
	return func(e *eventHandler) {
		e.workflowRegistryAddress = address
		e.workflowRegistryChainSelector = chainSelector
	}
}

func WithOrgResolver(orgResolver orgresolver.OrgResolver) func(*eventHandler) {
	return func(e *eventHandler) {
		e.orgResolver = orgResolver
	}
}

// WithDebugMode enables OTel tracing when debugMode is true.
// When disabled (default), a noop tracer is used for zero overhead.
// The debugMode is also propagated to workflow engines created by this handler.
func WithDebugMode(debugMode bool) func(*eventHandler) {
	return func(e *eventHandler) {
		e.lggr.Infow("Setting debug mode for workflow syncer", "debugMode", debugMode)
		e.debugMode = debugMode
		if debugMode {
			e.lggr.Errorw("WARNING: Debug mode is enabled for workflow syncer, this is not suitable for production")
			e.tracer = otel.Tracer("workflow_syncer")
		} else {
			// set to no-op just in case a real tracer was initialised elsewhere
			e.tracer = noop.NewTracerProvider().Tracer("")
		}
	}
}

func WithSecretsFetcher(sf v2.SecretsFetcher) func(*eventHandler) {
	return func(e *eventHandler) {
		e.secretsFetcher = sf
	}
}

// WithLocalSecretOverrides wires [CRE.LocalSecretOverrides]: per-workflow-owner name->secret map
func WithLocalSecretOverrides(lggr logger.Logger, perOwner map[string]map[string]string) func(*eventHandler) {
	return func(e *eventHandler) {
		if len(perOwner) == 0 {
			return
		}
		e.localSecretOverrides = make(map[string]map[string]string, len(perOwner))
		for k, m := range perOwner {
			e.localSecretOverrides[k] = maps.Clone(m)
		}
		owners := make([]string, 0, len(e.localSecretOverrides))
		for owner := range e.localSecretOverrides {
			owners = append(owners, owner)
		}
		lggr.Warnw("Per-owner local secret overrides are active; vault is used for secret IDs not listed under each owner",
			"numOwners", len(e.localSecretOverrides),
			"owners", owners)
	}
}

func WithModuleLRU(lru *ModuleLRU) func(*eventHandler) {
	return func(e *eventHandler) {
		e.moduleLRU = lru
	}
}

func WithModuleStore(store artifacts.SerialisedModuleStore) func(*eventHandler) {
	return func(e *eventHandler) {
		e.moduleStore = store
	}
}

func WithModuleCacheMetrics(cm *CacheMetrics) func(*eventHandler) {
	return func(e *eventHandler) {
		e.cacheMetrics = cm
	}
}

// WithModuleEngineVersion sets the engine version tag persisted alongside cached module
// binaries and compared on reload. Reloads from an older version are rejected with a log
// and a metric, ensuring process upgrades never reuse ABI-incompatible binaries.
func WithModuleEngineVersion(v string) func(*eventHandler) {
	return func(e *eventHandler) {
		e.moduleEngineVersion = v
	}
}

type WorkflowArtifactsStore interface {
	FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryIdentifier, configIdentifier string) ([]byte, []byte, error)
	GetWorkflowSpec(ctx context.Context, workflowID string) (*job.WorkflowSpec, error)
	ListWorkflowSpecs(ctx context.Context) ([]*job.WorkflowSpec, error)
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)
	DeleteWorkflowArtifacts(ctx context.Context, workflowID string) (*job.WorkflowSpec, error)
	PauseWorkflowArtifacts(ctx context.Context, workflowID string) error
	DeleteWorkflowArtifactsBatch(ctx context.Context, workflowIDs []string) error
}

// NewEventHandler returns a new eventHandler instance.
func NewEventHandler(
	lggr logger.Logger,
	workflowStore store.Store,
	donTimeStore *dontime.Store,
	useLocalTimeProvider bool,
	capRegistry core.CapabilitiesRegistry,
	executionHandlers *confidentialrelay.ExecutionHandlers,
	engineRegistry *EngineRegistry,
	emitter custmsg.MessageEmitter,
	engineLimiters *v2.EngineLimiters,
	featureFlags *v2.EngineFeatureFlags,
	ratelimiter *ratelimiter.RateLimiter,
	workflowLimits limits.ResourceLimiter[int],
	workflowArtifacts WorkflowArtifactsStore,
	workflowEncryptionKey workflowkey.Key,
	workflowDonSubscriber capabilities.DonSubscriber,
	opts ...func(*eventHandler),
) (*eventHandler, error) {
	if workflowStore == nil {
		return nil, errors.New("workflow store must be provided")
	}
	if capRegistry == nil {
		return nil, errors.New("capabilities registry must be provided")
	}
	if engineRegistry == nil {
		return nil, errors.New("engine registry must be provided")
	}
	if donTimeStore == nil && !useLocalTimeProvider {
		return nil, errors.New("donTimeStore must be provided")
	}

	eh := &eventHandler{
		lggr:                   lggr,
		workflowStore:          workflowStore,
		capRegistry:            capRegistry,
		executionHandlers:      executionHandlers,
		donTimeStore:           donTimeStore,
		useLocalTimeProvider:   useLocalTimeProvider,
		engineRegistry:         engineRegistry,
		emitter:                emitter,
		engineLimiters:         engineLimiters,
		featureFlags:           featureFlags,
		ratelimiter:            ratelimiter,
		workflowLimits:         workflowLimits,
		workflowArtifactsStore: workflowArtifacts,
		workflowEncryptionKey:  workflowEncryptionKey,
		workflowDonSubscriber:  workflowDonSubscriber,
		// default, enable via WithDebugMode
		tracer: noop.NewTracerProvider().Tracer(""),
	}
	metricsInst, metricsErr := newMetrics()
	if metricsErr != nil {
		return nil, fmt.Errorf("new metrics: %w", metricsErr)
	}
	eh.metrics = metricsInst
	eh.engineFactory = eh.engineFactoryFn
	for _, o := range opts {
		o(eh)
	}

	eh.Service, eh.eng = services.Config{
		Name: "EventHandler",
		// The workflow store and the spec meter are started and stopped
		// alongside the handler.
		NewSubServices: func(logger.Logger) []services.Service {
			var subs []services.Service
			if eh.workflowStore != nil {
				subs = append(subs, eh.workflowStore)
			}
			if eh.specMeter != nil {
				subs = append(subs, eh.specMeter)
			}
			return subs
		},
		Start: eh.start,
		Close: eh.close,
	}.NewServiceEngine(lggr)

	return eh, nil
}

func (h *eventHandler) start(context.Context) error {
	if h.moduleLRU != nil {
		h.moduleLRU.Start()
	}
	return nil
}

// SetWorkflowDon supplies the launcher-resolved workflow DON identity for
// metering. Called by the registry after WaitForDon, before any event is
// dispatched; the value is static for the life of the node.
func (h *eventHandler) SetWorkflowDon(don commoncap.DON) {
	h.specMeter.SetWorkflowDon(don)
}

func (h *eventHandler) close() error {
	if h.moduleLRU != nil {
		h.moduleLRU.Close()
	}
	es := h.engineRegistry.PopAll()
	// No metering is emitted on close: meter records anchor on workflow-spec
	// storage transitions, not engine lifecycle, so stopping an engine at
	// shutdown leaves the persisted spec (and therefore its metered level)
	// untouched. A spec that is genuinely released stops being snapshotted
	// (the spec meter sub-service unregisters itself after this hook runs).
	cs := make([]io.Closer, 0, len(es)+1)
	cs = append(cs, h.engineLimiters)
	for _, e := range es {
		cs = append(cs, e)
	}
	return services.CloseAll(cs...)
}

// toCommonHead converts our local Head struct back to chainlink-common Head
func toCommonHead(localHead Head) *commontypes.Head {
	return &commontypes.Head{
		Hash:      []byte(localHead.Hash),
		Height:    localHead.Height,
		Timestamp: localHead.Timestamp,
	}
}

func (h *eventHandler) Handle(ctx context.Context, event Event) error {
	ctx, span := h.tracer.Start(ctx, "handle_event",
		trace.WithAttributes(
			attribute.String("component", "workflow_syncer"),
			attribute.String("event_type", string(event.Name)),
		))
	defer span.End()

	switch event.Name {
	case WorkflowActivated:
		payload, ok := event.Data.(WorkflowActivatedEvent)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		wfID := payload.WorkflowID.Hex()
		wfOwner := hex.EncodeToString(payload.WorkflowOwner)
		orgID, ferr := h.fetchOrganizationID(ctx, wfOwner)
		if ferr != nil {
			h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", wfOwner, "error", ferr)
		}
		ctx = contexts.WithCRE(ctx, contexts.CRE{Org: orgID, Owner: wfOwner, Workflow: wfID})

		h.emitterMu.RLock()
		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, wfOwner,
			platform.KeyWorkflowTag, payload.WorkflowTag,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
			platform.KeyWorkflowSource, payload.Source,
		)
		h.emitterMu.RUnlock()

		var err error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), payload.BinaryURL, payload.ConfigURL, customerFacingError(err)); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()
		err = h.workflowActivatedEvent(ctx, payload)
		if err != nil {
			h.lggr.Errorw("failed to handle workflow activated event", "error", err, "workflowID", wfID)
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow activated event: %v", customerFacingError(err)), h.lggr)
			return err
		}

		h.lggr.Debugw("handled event (WorkflowActivated)", "workflowID", wfID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner),
			"workflowTag", payload.WorkflowTag, "type", event.Name)
		return nil
	case WorkflowPaused:
		payload, ok := event.Data.(WorkflowPausedEvent)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		wfID := payload.WorkflowID.Hex()
		wfOwner := hex.EncodeToString(payload.WorkflowOwner)
		orgID, ferr := h.fetchOrganizationID(ctx, wfOwner)
		if ferr != nil {
			h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", wfOwner, "error", ferr)
		}
		ctx = contexts.WithCRE(ctx, contexts.CRE{Org: orgID, Owner: wfOwner, Workflow: wfID})

		h.emitterMu.RLock()
		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
			platform.KeyWorkflowTag, payload.WorkflowTag,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
			platform.KeyWorkflowSource, payload.Source,
		)
		h.emitterMu.RUnlock()

		var err error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), payload.BinaryURL, payload.ConfigURL, customerFacingError(err)); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()

		if err = h.workflowPausedEvent(ctx, payload); err != nil {
			if errors.Is(err, ErrDrainInProgress) {
				logCustMsg(ctx, cma, fmt.Sprintf("workflow pause deferred: %v", err), h.lggr)
			} else {
				logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow paused event: %v", err), h.lggr)
			}
			return err
		}

		h.lggr.Debugw("handled event (WorkflowPaused)", "workflowID", wfID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner),
			"workflowTag", payload.WorkflowTag, "type", event.Name)
		return nil
	case WorkflowDeleted:
		payload, ok := event.Data.(WorkflowDeletedEvent)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		wfID := payload.WorkflowID.Hex()

		// Get workflow spec from database to get owner and name info for organization lookup,
		// and to check if the spec exists to determine if a MeterRecord is warranted.
		var wfOwner, wfName, orgID string
		if spec, gerr := h.workflowArtifactsStore.GetWorkflowSpec(ctx, wfID); gerr == nil && spec != nil {
			wfOwner = spec.WorkflowOwner
			wfName = spec.WorkflowName
		} else if gerr != nil && !errors.Is(gerr, sql.ErrNoRows) {
			h.lggr.Errorw("failed to read workflow spec during deletion, proceeding without metadata", "workflowID", wfID, "error", gerr)
		}
		if wfOwner != "" {
			if resolvedOrgID, orgErr := h.fetchOrganizationID(ctx, wfOwner); orgErr != nil {
				h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", wfOwner, "error", orgErr)
			} else {
				orgID = resolvedOrgID
			}
		}
		ctx = contexts.WithCRE(ctx, contexts.CRE{Org: orgID, Owner: wfOwner, Workflow: wfID})

		h.emitterMu.RLock()
		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, wfName,
			platform.KeyWorkflowOwner, wfOwner,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
			platform.KeyWorkflowSource, payload.Source,
		)
		h.emitterMu.RUnlock()

		var herr error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), "", "", herr); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()

		if herr = h.workflowDeletedEvent(ctx, payload, wfOwner); herr != nil {
			if errors.Is(herr, ErrDrainInProgress) {
				logCustMsg(ctx, cma, fmt.Sprintf("workflow deletion deferred: %v", herr), h.lggr)
			} else {
				logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow deleted event: %v", herr), h.lggr)
			}
			return herr
		}

		h.lggr.Debugw("handled event (WorkflowDeleted)", "workflowID", wfID, "workflowName", wfName, "workflowOwner", wfOwner, "organizationID", orgID, "type", event.Name)
		return nil
	default:
		return fmt.Errorf("event type unsupported: %v", event.Name)
	}
}

// workflowActivatedEvent handles the WorkflowActivatedEvent event type.
// This method redirects to workflowRegisteredEvent since they have identical processing logic.
func (h *eventHandler) workflowActivatedEvent(
	ctx context.Context,
	payload WorkflowActivatedEvent,
) error {
	// Convert WorkflowActivatedEvent to WorkflowRegisteredEvent since they have identical fields
	registeredPayload := WorkflowRegisteredEvent(payload)
	return h.workflowRegisteredEvent(ctx, registeredPayload)
}

// workflowTagBackfillActive reports whether the cresettings
// PerWorkflow.FeatureWorkflowTagBackfillActivePeriod window covers time.Now().
// Fail-closed: any error or missing limiter is treated as inactive so a
// misconfigured deploy cannot silently start rewriting workflow_tag.
func (h *eventHandler) workflowTagBackfillActive(ctx context.Context) bool {
	if h.featureFlags == nil || h.featureFlags.WorkflowTagBackfill == nil {
		return false
	}
	return h.featureFlags.WorkflowTagBackfill.Check(ctx, config.Timestamp(time.Now().Unix())) == nil
}

// workflowRegisteredEvent handles the WorkflowRegisteredEvent event type.
// This method must remain idempotent and must not error if retried multiple times.
// workflowRegisteredEvent proceeds in two phases:
// - phase 1 synchronizes the database state
// - phase 2 synchronizes the state of the engine registry.
func (h *eventHandler) workflowRegisteredEvent(
	ctx context.Context,
	payload WorkflowRegisteredEvent,
) error {
	ctx, span := h.tracer.Start(ctx, "workflow_registered",
		trace.WithAttributes(
			attribute.String("component", "workflow_syncer"),
			attribute.String("workflow_name", payload.WorkflowName),
		))
	defer span.End()

	status := toSpecStatus(payload.Status)

	// First, let's synchronize the database state.
	// We need to handle three cases:
	// - new registration, without an existing DB record
	// - existing registration that has been updated with new artifacts, and potentially also the status
	// - existing registration that has been updated with a new status
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, payload.WorkflowID.Hex())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		newSpec, innerErr := h.createWorkflowSpec(ctx, payload)
		if innerErr != nil {
			return innerErr
		}

		h.specMeter.EmitSpecDelta(ctx, 1, payload.WorkflowID.Hex(), hex.EncodeToString(payload.WorkflowOwner),
			resourcemanager.EventID("workflow-spec-register", payload.WorkflowID.Hex(), strconv.FormatUint(payload.CreatedAt, 10)))

		spec = newSpec
	case err != nil:
		return fmt.Errorf("failed to get workflow spec: %w", err)
	case spec.WorkflowID != payload.WorkflowID.Hex() ||
		spec.WorkflowOwner != hex.EncodeToString(payload.WorkflowOwner) ||
		spec.WorkflowName != payload.WorkflowName:
		newSpec, innerErr := h.createWorkflowSpec(ctx, payload)
		if innerErr != nil {
			return innerErr
		}
		// A different spec's artifacts were persisted under this key: the newly
		// stored spec is a fresh durable resource, so emit a +1 delta.
		h.specMeter.EmitSpecDelta(ctx, 1, payload.WorkflowID.Hex(), hex.EncodeToString(payload.WorkflowOwner),
			resourcemanager.EventID("workflow-spec-register", payload.WorkflowID.Hex(), strconv.FormatUint(payload.CreatedAt, 10)))

		spec = newSpec
	case status == job.WorkflowSpecStatusActive && spec.Workflow == "":
		// Activating a paused tombstone: the artifact payload was cleared at
		// pause time, so refetch and re-persist it. Level-neutral for metering
		// (the registration generation was never released), so no delta.
		newSpec, innerErr := h.createWorkflowSpec(ctx, payload)
		if innerErr != nil {
			return innerErr
		}
		spec = newSpec
	case spec.Status != status:
		spec.Status = status
		if _, innerErr := h.workflowArtifactsStore.UpsertWorkflowSpec(ctx, spec); innerErr != nil {
			return fmt.Errorf("failed to update workflow spec: %w", innerErr)
		}
		// Status-only flip: no artifact-persistence transition, no delta.
	}

	// backfill registered_at, source, workflow_tag when necessary
	backfill := false
	if spec.RegisteredAt == 0 && payload.CreatedAt > 0 {
		spec.RegisteredAt = int64(payload.CreatedAt) //nolint:gosec // G115: CreatedAt is a timestamp that cannot overflow int64
		backfill = true
	}
	if spec.Source == "" && payload.Source != "" {
		spec.Source = payload.Source
		backfill = true
	}
	// WorkflowTag is used by remote capability requests; an empty local value
	// diverges the request hash from nodes that have the on-chain tag. Gated by
	// cresettings PerWorkflow.FeatureWorkflowTagBackfillActivePeriod: the
	// backfill only fires when time.Now() falls inside the configured window,
	// so ops can coordinate a healing pass across the DON. Default window is
	// far-future, so a fresh deploy is a no-op.
	var (
		tagBackfilled bool
		tagBefore     string
	)
	if h.workflowTagBackfillActive(ctx) &&
		spec.WorkflowTag != payload.WorkflowTag && payload.WorkflowTag != "" {
		tagBefore = spec.WorkflowTag
		spec.WorkflowTag = payload.WorkflowTag
		backfill = true
		tagBackfilled = true
	}
	if backfill {
		if _, err := h.workflowArtifactsStore.UpsertWorkflowSpec(ctx, spec); err != nil {
			h.lggr.Warnw("failed to backfill registered_at/source/workflow_tag", "workflowID", spec.WorkflowID, "err", err)
		} else if tagBackfilled {
			reason := "refreshed"
			if tagBefore == "" {
				reason = "filled_empty"
			}
			h.lggr.Infow("backfilled workflow_tag",
				"workflowID", spec.WorkflowID,
				"workflowOwner", spec.WorkflowOwner,
				"workflowName", spec.WorkflowName,
				"tagBefore", tagBefore,
				"tagAfter", spec.WorkflowTag,
				"reason", reason,
			)
			h.metrics.incrementWorkflowTagBackfill(ctx, reason)
		}
	}

	// Next, let's synchronize the engine.
	// If the state isn't active, we shouldn't have an engine running.

	// Let's try to clean one up if it exists
	if spec.Status != job.WorkflowSpecStatusActive {
		return h.tryEngineCleanup(payload.WorkflowID)
	}

	// We know we need an engine, let's make sure that there isn't already one running for this workflow ID.
	prevEngine, ok := h.engineRegistry.Get(payload.WorkflowID)
	if ok && prevEngine.Ready() == nil && spec.Status == job.WorkflowSpecStatusActive {
		drainable, isDrainable := prevEngine.Service.(DrainableService)
		isDraining := false
		if isDrainable {
			_, isDraining = drainable.DrainStartedAt()
		}
		if isDrainable && isDraining {
			h.lggr.Infow("engine is draining, replacing with a new engine", "workflowID", payload.WorkflowID.Hex())
		}

		// This is the happy-path, we're done.
		if !isDrainable || !isDraining {
			return nil
		}
	}

	// Any other case ->
	// - engine in registry, but service isn't running
	// - engine in registry and service is running, but it's draining and must be replaced
	// - state isn't active
	// Let's clean up and recreate

	cleanupErr := h.tryEngineCleanup(payload.WorkflowID)
	if cleanupErr != nil {
		return fmt.Errorf("could not clean up old engine: %w", cleanupErr)
	}

	return h.tryEngineCreate(ctx, spec, payload.Source)
}

func toSpecStatus(s uint8) job.WorkflowSpecStatus {
	switch s {
	case WorkflowStatusActive:
		return job.WorkflowSpecStatusActive
	case WorkflowStatusPaused:
		return job.WorkflowSpecStatusPaused
	default:
		return job.WorkflowSpecStatusDefault
	}
}

func (h *eventHandler) createWorkflowSpec(ctx context.Context, payload WorkflowRegisteredEvent) (*job.WorkflowSpec, error) {
	ctx, span := h.tracer.Start(ctx, "fetch_artifacts",
		trace.WithAttributes(
			attribute.String("component", "workflow_syncer"),
			attribute.String("workflow_name", payload.WorkflowName),
		))
	defer span.End()

	wfID := payload.WorkflowID.Hex()
	owner := hex.EncodeToString(payload.WorkflowOwner)
	orgID, err := h.fetchOrganizationID(ctx, owner)
	if err != nil {
		h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", owner, "error", err)
	}
	ctx = contexts.WithCRE(ctx, contexts.CRE{Org: orgID, Owner: owner, Workflow: wfID})

	// With Workflow Registry contract v2 the BinaryURL and ConfigURL are expected to be identifiers that put through the Storage Service.
	decodedBinary, config, err := h.workflowArtifactsStore.FetchWorkflowArtifacts(ctx, wfID, payload.BinaryURL, payload.ConfigURL)
	if err != nil {
		return nil, err
	}

	status := toSpecStatus(payload.Status)

	// Create a new entry in the workflow_specs_v2 table corresponding for the new workflow, with the contents of the binaryIdentifier + configIdentifier in the table
	entry := &job.WorkflowSpec{
		Workflow:      hex.EncodeToString(decodedBinary),
		Config:        string(config),
		WorkflowID:    wfID,
		Status:        status,
		WorkflowOwner: owner,
		WorkflowName:  payload.WorkflowName,
		WorkflowTag:   payload.WorkflowTag,
		SpecType:      job.WASMFile,
		BinaryURL:     payload.BinaryURL,
		ConfigURL:     payload.ConfigURL,
		Attributes:    payload.Attributes,
		RegisteredAt:  int64(payload.CreatedAt), //nolint:gosec // G115: CreatedAt is a timestamp that cannot overflow int64
		Source:        payload.Source,
	}

	if _, err = h.workflowArtifactsStore.UpsertWorkflowSpec(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to upsert workflow spec: %w", err)
	}

	return entry, nil
}

// fetchOrganizationID fetches the organization ID for the given workflow owner using the OrgResolver
func (h *eventHandler) fetchOrganizationID(ctx context.Context, workflowOwner string) (string, error) {
	if h.orgResolver == nil {
		return "", errors.New("org resolver is not available")
	}

	organizationID, err := h.orgResolver.Get(ctx, workflowOwner)
	if err != nil {
		h.lggr.Warnw("Failed to get organization ID from org resolver", "workflowOwner", workflowOwner, "error", err)
		return "", err
	}

	if organizationID == "" {
		h.lggr.Warnw("No organization ID returned from org resolver", "workflowOwner", workflowOwner)
		return "", errors.New("no organization ID returned from org resolver")
	}

	h.lggr.Debugw("Successfully retrieved organization ID from org resolver", "workflowOwner", workflowOwner, "organizationId", organizationID)
	return organizationID, nil
}

func (h *eventHandler) engineFactoryFn(ctx context.Context, workflowID, owner string, name types.WorkflowName, tag string, config, binary []byte, binaryURL string, initDone chan<- error) (services.Service, error) {
	lggr := logger.Named(h.lggr, "WorkflowEngine.Module")
	lggr = logger.With(lggr, "workflowID", workflowID, "workflowName", name, "workflowOwner", owner)
	var sdkName string
	h.emitterMu.RLock()
	labeler := h.emitter
	h.emitterMu.RUnlock()
	moduleConfig := &host.ModuleConfig{
		Logger:                               lggr,
		Labeler:                              labeler,
		MemoryLimiter:                        h.engineLimiters.WASMMemorySize,
		MaxCompressedBinaryLimiter:           h.engineLimiters.WASMCompressedBinarySize,
		MaxDecompressedBinaryLimiter:         h.engineLimiters.WASMBinarySize,
		MaxResponseSizeLimiter:               h.engineLimiters.ExecutionResponse,
		PendingCallsLimiter:                  h.engineLimiters.CapabilityConcurrency,
		EnableUserMetricsLimiter:             h.engineLimiters.UserMetricEnabled,
		MaxUserMetricPayloadLimiter:          h.engineLimiters.UserMetricPayload,
		MaxUserMetricNameLengthLimiter:       h.engineLimiters.UserMetricNameLength,
		MaxUserMetricLabelsPerMetricLimiter:  h.engineLimiters.UserMetricLabelsPerMetric,
		MaxUserMetricLabelValueLengthLimiter: h.engineLimiters.UserMetricLabelValueLength,
		SdkLabeler: func(name string) {
			sdkName = name
			h.emitterMu.Lock()
			h.emitter = h.emitter.With(platform.KeySDK, name)
			h.emitterMu.Unlock()
		},
	}

	h.lggr.Debugw("Creating module for workflowID", "workflowID", workflowID)

	module, err := host.NewModule(ctx, moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, err
	}

	h.lggr.Debugw("Finished creating module for workflowID", "workflowID", workflowID)

	if module.IsLegacyDAG() { // V1 aka "DAG"
		return nil, errors.New("legacy DAG workflows are not supported")
	}

	// V2 aka "NoDAG"
	// Wrap the local WASM module in a RequirementSelectingModule that routes
	// triggers with a TEE requirement to the ConfidentialModule (which delegates
	// to the confidential-workflows capability and runs the WASM inside the
	// enclave). Triggers without a TEE requirement continue to run locally.
	binaryHash := v2.ComputeBinaryHash(binary)
	confLggr := logger.Named(h.lggr, "WorkflowEngine.ConfidentialModule")
	confLggr = logger.With(confLggr, "workflowID", workflowID, "workflowName", name, "workflowOwner", owner)
	confidential, err := v2.NewConfidentialModule(h.capRegistry, h.executionHandlers, binaryURL, binaryHash, workflowID, owner, name.String(), tag, h.fetchOrganizationID, h.engineLimiters.ConfidentialWorkflowsEnabled, h.engineLimiters.Settings, confLggr)
	if err != nil {
		return nil, fmt.Errorf("failed to create confidential module: %w", err)
	}
	engineModule := h.createEngineModule(ctx, workflowID, binary, moduleConfig, module)
	selectingModule := generichost.NewRequirementSelectingModule(
		generichost.ModuleAndHandler{Module: engineModule},
		[]generichost.ModuleAndHandler{{
			Module:              confidential,
			RequirementsHandler: generichost.RequirementsHandler{Tee: confidential.Tee},
		}},
	)
	cfg := h.newV2EngineConfig(selectingModule, workflowID, owner, tag, sdkName, name, config)

	h.wireInitDoneHook(cfg, initDone)

	return v2.NewEngine(cfg)
}

func (h *eventHandler) createEngineModule(
	ctx context.Context,
	workflowID string,
	binary []byte,
	moduleConfig *host.ModuleConfig,
	module generichost.Module,
) generichost.Module {
	engineModule := module
	if h.moduleLRU != nil && h.moduleStore != nil {
		storeStart := time.Now()
		storeErr := h.moduleStore.StoreModule(workflowID, binary, h.moduleEngineVersion)
		if h.metrics != nil {
			h.metrics.recordModuleStore(ctx, time.Since(storeStart), storeErr == nil)
		}
		if storeErr != nil {
			h.lggr.Warnw("Failed to cache module binary to disk, LRU eviction disabled for this workflow", "workflowID", workflowID, "err", storeErr)
		} else {
			evictable := NewEvictableModule(module, moduleConfig, h.moduleStore, workflowID, h.moduleEngineVersion, nil, h.cacheMetrics, int64(len(binary)), host.WithDeterminism())
			h.moduleLRU.Register(workflowID, evictable)
			engineModule = evictable
		}
	}

	return engineModule
}

// stopEngine drains (returning ErrDrainInProgress while executions remain) and
// closes the engine for workflowID if one is registered.
// Returns the drainable handle (nil-able) for drain-completed metrics.
func (h *eventHandler) stopEngine(ctx context.Context, workflowID types.WorkflowID) (DrainableService, error) {
	e, ok := h.engineRegistry.Get(workflowID)
	var drainable DrainableService
	if ok {
		var isDrainable bool
		if drainable, isDrainable = e.Service.(DrainableService); isDrainable {
			if started := drainable.Drain(); started {
				h.lggr.Infow("initiated drain for workflow in workflow engine", "workflowID", workflowID.String())
				if h.metrics != nil {
					h.metrics.incrementDrainStarted(ctx)
				}
			}

			if active := drainable.ActiveExecutions(); active > 0 {
				if h.metrics != nil {
					h.metrics.incrementDeleteDeferred(ctx, "drain_in_progress")
				}
				h.lggr.Infow("workflow deletion deferred: active executions still running",
					"workflowID", workflowID.String(),
					"activeExecutions", active)
				return nil, fmt.Errorf("%w: %d active executions still running", ErrDrainInProgress, active)
			}
		}

		if innerErr := e.Close(); innerErr != nil && !errors.Is(innerErr, services.ErrAlreadyStopped) {
			return nil, fmt.Errorf("failed to close workflow engine: %w", innerErr)
		}
	}
	return drainable, nil
}

// releaseSpecStorage deletes the persisted spec row and, iff a row was
// actually removed, emits the -1 generation delta. RowsAffected is the
// exactly-once gate: redeliveries observe no row and emit nothing; a transient
// DELETE error returns before emission so the retry emits when the delete
// lands. The event_id is derived AFTER the delete from the row's persisted
// registered_at, so all nodes emit the identical id. Module-cache cleanup
// always runs. owner may be empty; org resolution is fail-open.
func (h *eventHandler) releaseSpecStorage(ctx context.Context, workflowID, owner string) error {
	deletedSpec, err := h.workflowArtifactsStore.DeleteWorkflowArtifacts(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to delete workflow artifacts: %w", err)
	}
	if deletedSpec != nil {
		parts := []string{workflowID}
		if deletedSpec.RegisteredAt > 0 {
			parts = append(parts, strconv.FormatInt(deletedSpec.RegisteredAt, 10))
		}
		h.specMeter.EmitSpecDelta(ctx, -1, workflowID, owner, resourcemanager.EventID("workflow-spec-delete", parts...))
	}
	h.cleanupModuleCache(workflowID)
	return nil
}

// workflowPausedEvent handles WorkflowPaused. Idempotent.
//
// Pause is level-neutral for metering: the spec row survives as a tombstone
// (status=paused, artifact payload cleared to free storage) so the
// registration generation stays metered until deletion — mirroring the
// on-chain registry, which retains the paused registration. Emitting per-cycle
// pause/activate deltas is impossible without drift: no DON-consistent
// per-occurrence discriminator exists, so a -1 here would force the
// re-activation +1 to reuse the register event_id and be dropped by consumer
// dedup. Deltas therefore fire only at generation boundaries (register/delete).
func (h *eventHandler) workflowPausedEvent(
	ctx context.Context,
	payload WorkflowPausedEvent,
) error {
	workflowID := payload.WorkflowID.Hex()
	if _, err := h.stopEngine(ctx, payload.WorkflowID); err != nil {
		return err
	}
	if err := h.workflowArtifactsStore.PauseWorkflowArtifacts(ctx, workflowID); err != nil {
		return fmt.Errorf("failed to pause workflow artifacts: %w", err)
	}
	h.cleanupModuleCache(workflowID)
	if _, err := h.engineRegistry.Pop(payload.WorkflowID); err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}
	return nil
}

// workflowDeletedEvent handles the WorkflowDeletedEvent event type. This method must remain idempotent.
func (h *eventHandler) workflowDeletedEvent(
	ctx context.Context,
	payload WorkflowDeletedEvent,
	owner string,
) error {
	workflowID := payload.WorkflowID.Hex()
	drainable, err := h.stopEngine(ctx, payload.WorkflowID)
	if err != nil {
		return err
	}
	if err = h.releaseSpecStorage(ctx, workflowID, owner); err != nil {
		return err
	}
	_, err = h.engineRegistry.Pop(payload.WorkflowID)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if drainable != nil {
		startedAt, exists := drainable.DrainStartedAt()
		if exists && h.metrics != nil {
			h.metrics.recordDrainCompleted(ctx, time.Since(startedAt))
		}
	}
	return nil
}

// ListWorkflowSpecs backs the orphan sweep.
func (h *eventHandler) ListWorkflowSpecs(ctx context.Context) ([]*job.WorkflowSpec, error) {
	return h.workflowArtifactsStore.ListWorkflowSpecs(ctx)
}

// tryEngineCleanup attempts to stop the workflow engine for the given workflow ID.  Does nothing if the
// workflow engine is not running.
func (h *eventHandler) tryEngineCleanup(workflowID types.WorkflowID) error {
	e, ok := h.engineRegistry.Get(workflowID)
	if !ok {
		return nil
	}

	// Close the engine, then remove it from the registry only once cleanup has succeeded. The
	// registry entry is what tells us this workflow still needs cleanup, so if a step fails we
	// leave it in place and return the error, allowing a future attempt to retry in case the
	// failure was spurious. Close is idempotent, so ErrAlreadyStopped means the engine was
	// already closed on a prior attempt and is treated as success rather than a failure.
	if err := e.Close(); err != nil && !errors.Is(err, services.ErrAlreadyStopped) {
		return fmt.Errorf("failed to close workflow engine: %w", err)
	}

	h.cleanupModuleCache(workflowID.Hex())

	if _, err := h.engineRegistry.Pop(workflowID); err != nil {
		return fmt.Errorf("failed to remove workflow engine: %w", err)
	}
	return nil
}

func (h *eventHandler) cleanupModuleCache(workflowID string) {
	if h.moduleLRU != nil {
		h.moduleLRU.Deregister(workflowID)
	}
	if h.moduleStore != nil {
		if err := h.moduleStore.DeleteModule(workflowID); err != nil {
			h.lggr.Warnw("Failed to delete cached module binary", "workflowID", workflowID, "err", err)
		}
	}
}

// tryEngineCreate attempts to create a new workflow engine, start it, and register it with the engine registry.
// This function waits for the engine to complete initialization (including trigger subscriptions) before returning,
// ensuring that the workflowActivated event accurately reflects the deployment status including trigger registration.
func (h *eventHandler) tryEngineCreate(ctx context.Context, spec *job.WorkflowSpec, source string) error {
	ctx, span := h.tracer.Start(ctx, "engine_create",
		trace.WithAttributes(
			attribute.String("component", "workflow_syncer"),
			attribute.String("workflow_name", spec.WorkflowName),
			attribute.String("source", source),
		))
	defer span.End()

	// Ensure the capabilities registry is ready before creating any Engine instances.
	// This should be guaranteed by the Workflow Registry Syncer.
	if err := h.ensureCapRegistryReady(ctx); err != nil {
		return fmt.Errorf("failed to ensure capabilities registry is ready: %w", err)
	}

	decodedBinary, err := hex.DecodeString(spec.Workflow)
	if err != nil {
		return nonRetryable(fmt.Errorf("failed to decode workflow spec binary: %w", err))
	}
	// Free the hex-encoded binary string as it is not needed beyond this decode
	spec.Workflow = ""

	// Workflow Registry version >2 no longer handles secrets
	secretsURL := ""

	// Before running the engine, handle validations
	// Workflow ID should match what is generated from the stored artifacts
	ownerBytes, err := hex.DecodeString(spec.WorkflowOwner)
	if err != nil {
		return nonRetryable(fmt.Errorf("failed to decode owner: %w", err))
	}
	configBytes := []byte(spec.Config)
	hash, err := pkgworkflows.GenerateWorkflowID(ownerBytes, spec.WorkflowName, decodedBinary, configBytes, secretsURL)
	if err != nil {
		return fmt.Errorf("failed to generate workflow id: %w", err)
	}
	wid, err := types.WorkflowIDFromHex(spec.WorkflowID)
	if err != nil {
		return nonRetryable(fmt.Errorf("invalid workflow id: %w", err))
	}
	if !types.WorkflowID(hash).Equal(wid) {
		return nonRetryable(fmt.Errorf("workflowID mismatch: %x != %x", hash, wid))
	}

	// Start a new WorkflowEngine instance, and add it to local engine registry
	workflowName, err := types.NewWorkflowName(spec.WorkflowName)
	if err != nil {
		return nonRetryable(fmt.Errorf("invalid workflow name: %w", err))
	}

	// Create a channel to receive the initialization result.
	// This allows us to wait for the engine to complete initialization (including trigger subscriptions)
	// before emitting the workflowActivated event, ensuring the event accurately reflects deployment status.
	initDone := make(chan error, 1)
	var engine services.Service

	engine, err = h.engineFactory(ctx, spec.WorkflowID, spec.WorkflowOwner, workflowName, spec.WorkflowTag, configBytes, decodedBinary, spec.BinaryURL, initDone)
	if err != nil {
		return fmt.Errorf("failed to create workflow engine: %w", err)
	}

	if err = engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start workflow engine: %w", err)
	}

	// Wait for the engine to complete initialization (including trigger subscriptions).
	// This ensures we don't emit workflowActivated events before the engine initializes successfully.
	select {
	case <-ctx.Done():
		// Context cancelled while waiting for initialization
		if closeErr := engine.Close(); closeErr != nil {
			h.lggr.Errorw("failed to close engine after context cancellation", "error", closeErr, "workflowID", spec.WorkflowID)
		}
		return fmt.Errorf("context cancelled while waiting for engine initialization: %w", ctx.Err())
	case initErr := <-initDone:
		if initErr != nil {
			// Engine initialization failed (e.g., trigger subscription failed)
			if closeErr := engine.Close(); closeErr != nil {
				h.lggr.Errorw("failed to close engine after initialization failure", "error", closeErr, "workflowID", spec.WorkflowID)
			}
			return fmt.Errorf("engine initialization failed: %w", initErr)
		}
	}

	// Engine is fully initialized, add to registry with source tracking and identity fingerprint
	reconcileKey, err := ReconcileKey(ownerBytes, spec.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to compute reconcile key: %w", err)
	}
	if err := h.engineRegistry.AddWithReconcileKey(wid, source, reconcileKey, engine); err != nil {
		if closeErr := engine.Close(); closeErr != nil {
			return fmt.Errorf("failed to close workflow engine: %w during invariant violation: %w", closeErr, err)
		}

		// Check for WorkflowID collision across sources
		if errors.Is(err, ErrAlreadyExists) {
			existingEntry, found := h.engineRegistry.Get(wid)
			if found {
				h.lggr.Warnw("WorkflowID collision detected: workflow already exists from different source",
					"workflowID", wid.Hex(),
					"attemptedSource", source,
					"existingSource", existingEntry.Source,
					"hint", "Each workflow ID should only be registered from a single source. Check your workflow configurations for duplicates.")
			}
		}

		// This shouldn't happen because we call the handler serially and
		// check for running engines above, see the call to engineRegistry.Contains.
		return fmt.Errorf("invariant violation: %w", err)
	}
	return nil
}

func (h *eventHandler) overrideFetcherForOwner(owner string) v2.SecretsFetcher {
	if h.localSecretOverrides == nil {
		return nil
	}
	key, err := v2.LocalSecretOverrideOwnerKey(owner)
	if err != nil {
		h.lggr.Errorw("invalid workflow owner for local secret overrides", "owner", owner, "err", err)
		return nil
	}
	overrides := h.localSecretOverrides[key]
	if len(overrides) == 0 {
		return nil
	}
	return v2.NewLocalSecretsFetcher(owner, overrides)
}

// newV2EngineConfig builds the common EngineConfig shared by both the normal
// WASM engine and the confidential engine paths. Caller supplies the module.
func (h *eventHandler) newV2EngineConfig(
	module host.ModuleV2,
	workflowID, owner, tag, sdkName string,
	name types.WorkflowName,
	config []byte,
) *v2.EngineConfig {
	return &v2.EngineConfig{
		Lggr:                  h.lggr,
		Module:                module,
		WorkflowConfig:        config,
		CapRegistry:           h.capRegistry,
		DonSubscriber:         h.workflowDonSubscriber,
		UseLocalTimeProvider:  h.useLocalTimeProvider,
		DonTimeStore:          h.donTimeStore,
		ExecutionsStore:       h.workflowStore,
		WorkflowID:            workflowID,
		WorkflowOwner:         owner,
		WorkflowName:          name,
		WorkflowTag:           tag,
		WorkflowEncryptionKey: h.workflowEncryptionKey,

		LocalLimits:         v2.EngineLimits{}, // all defaults
		LocalLimiters:       h.engineLimiters,
		FeatureFlags:        h.featureFlags,
		GlobalWorkflowLimit: h.workflowLimits,

		BeholderEmitter: func() custmsg.MessageEmitter {
			h.emitterMu.RLock()
			defer h.emitterMu.RUnlock()
			return h.emitter
		}(),
		BillingClient: h.billingClient,

		WorkflowRegistryAddress:       h.workflowRegistryAddress,
		WorkflowRegistryChainSelector: h.workflowRegistryChainSelector,
		OrgResolver:                   h.orgResolver,
		SecretsFetcher:                h.secretsFetcher,
		OverrideFetcher:               h.overrideFetcherForOwner(owner),
		DebugMode:                     h.debugMode,
		SdkName:                       sdkName,

		ShardOrchestratorClient: h.shardOrchestratorClient,
		ShardingEnabled:         h.shardingEnabled,
		MyShardID:               h.myShardID,
		ShardRoutingSteady:      h.shardRoutingSteady,
		ShardResolver:           h.shardResolver,
	}
}

// wireInitDoneHook wires the initDone channel to the OnInitialized lifecycle hook.
// This will be called when the engine completes initialization (including trigger subscriptions).
// We compose with any existing hook to avoid overwriting test hooks or other user-provided hooks.
func (h *eventHandler) wireInitDoneHook(cfg *v2.EngineConfig, initDone chan<- error) {
	if initDone == nil {
		return
	}
	existingHook := cfg.Hooks.OnInitialized
	cfg.Hooks.OnInitialized = func(err error) {
		// Signal completion to the handler first
		initDone <- err
		// Then call any existing hook (e.g., from tests)
		if existingHook != nil {
			existingHook(err)
		}
	}
}

// logCustMsg emits a custom message to the external sink and logs an error if that fails.
func logCustMsg(ctx context.Context, cma custmsg.MessageEmitter, msg string, log logger.Logger) {
	err := cma.Emit(ctx, msg)
	if err != nil {
		logger.Helper(log, 1).Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
	}
}

func (h *eventHandler) EmitActivationAbandoned(
	ctx context.Context,
	event Event,
	reason eventsv2.ActivationAbandonReason,
	activationErr error,
	retryCount int32,
) error {
	if event.Name != WorkflowActivated || h == nil || h.emitter == nil {
		return nil
	}

	payload, ok := event.Data.(WorkflowActivatedEvent)
	if !ok {
		return newHandlerTypeError(event.Data)
	}

	wfID := payload.WorkflowID.Hex()
	wfOwner := hex.EncodeToString(payload.WorkflowOwner)
	orgID, ferr := h.fetchOrganizationID(ctx, wfOwner)
	if ferr != nil {
		h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", wfOwner, "error", ferr)
	}

	h.emitterMu.RLock()
	cma := h.emitter.With(
		platform.KeyWorkflowID, wfID,
		platform.KeyWorkflowName, payload.WorkflowName,
		platform.KeyWorkflowOwner, wfOwner,
		platform.KeyWorkflowTag, payload.WorkflowTag,
		platform.KeyOrganizationID, orgID,
		platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
		platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
		platform.KeyWorkflowSource, payload.Source,
	)
	h.emitterMu.RUnlock()

	return events.EmitWorkflowActivationAbandonedV2(
		ctx,
		cma.Labels(),
		payload.BinaryURL,
		payload.ConfigURL,
		reason,
		customerFacingError(activationErr),
		retryCount,
	)
}

func (h *eventHandler) ensureCapRegistryReady(ctx context.Context) error {
	// Check every 500ms until the capabilities registry is ready.
	retryInterval := time.Millisecond * time.Duration(500)
	return internal.RunWithRetries(
		ctx,
		h.lggr,
		retryInterval,
		0, // infinite retries, until context is done
		func() error {
			// Test that the registry is ready by attempting to get the local node
			_, err := h.capRegistry.LocalNode(ctx)
			if err != nil {
				return fmt.Errorf("capabilities registry not ready: %w", err)
			}
			return nil
		},
	)
}

// customerFacingError returns a deterministic, user-actionable error for beholder emission.
// Internal errors (e.g. ArtifactFetchError with per-node signed URLs) are replaced with a
// clean message so that workflow-service can aggregate error_message across nodes.
func customerFacingError(err error) error {
	if err == nil {
		return nil
	}
	if fetchErr, ok := errors.AsType[*types.ArtifactFetchError](err); ok {
		return errors.New(fetchErr.CustomerError())
	}
	return err
}

func newHandlerTypeError(data any) error {
	return fmt.Errorf("invalid data type %T for event", data)
}
