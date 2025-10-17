package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

type ORM interface {
	artifacts.WorkflowSpecsDS
}

type engineFactoryFn func(ctx context.Context, wfid string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte) (services.Service, error)

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding method that handles the event.
type eventHandler struct {
	lggr logger.Logger

	workflowStore          store.Store
	capRegistry            core.CapabilitiesRegistry
	donTimeStore           *dontime.Store
	useLocalTimeProvider   bool
	engineRegistry         *EngineRegistry
	emitter                custmsg.MessageEmitter
	engineFactory          engineFactoryFn
	engineLimiters         *v2.EngineLimiters
	ratelimiter            limits.RateLimiter
	workflowLimits         limits.ResourceLimiter[int]
	workflowArtifactsStore WorkflowArtifactsStore
	workflowEncryptionKey  workflowkey.Key
	billingClient          metering.BillingClient
	orgResolver            orgresolver.OrgResolver

	// WorkflowRegistryAddress is the address of the workflow registry contract
	workflowRegistryAddress string
	// WorkflowRegistryChainSelector is the chain selector for the workflow registry
	workflowRegistryChainSelector string
}

func WithEngineRegistry(er *EngineRegistry) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineRegistry = er
	}
}

func WithEngineFactoryFn(efn engineFactoryFn) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineFactory = efn
	}
}

func WithStaticEngine(engine services.Service) func(*eventHandler) {
	return func(e *eventHandler) {
		e.engineFactory = func(_ context.Context, _ string, _ string, _ types.WorkflowName, _ string, _ []byte, _ []byte) (services.Service, error) {
			return engine, nil
		}
	}
}

func WithBillingClient(client metering.BillingClient) func(*eventHandler) {
	return func(e *eventHandler) {
		e.billingClient = client
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

type WorkflowArtifactsStore interface {
	FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryIdentifier, configIdentifier string) ([]byte, []byte, error)
	GetWorkflowSpec(ctx context.Context, workflowID string) (*job.WorkflowSpec, error)
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)
	DeleteWorkflowArtifacts(ctx context.Context, workflowID string) error
}

// NewEventHandler returns a new eventHandler instance.
func NewEventHandler(
	lggr logger.Logger,
	workflowStore store.Store,
	donTimeStore *dontime.Store,
	useLocalTimeProvider bool,
	capRegistry core.CapabilitiesRegistry,
	engineRegistry *EngineRegistry,
	emitter custmsg.MessageEmitter,
	engineLimiters *v2.EngineLimiters,
	ratelimiter limits.RateLimiter,
	workflowLimits limits.ResourceLimiter[int],
	workflowArtifacts WorkflowArtifactsStore,
	workflowEncryptionKey workflowkey.Key,
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
		donTimeStore:           donTimeStore,
		useLocalTimeProvider:   useLocalTimeProvider,
		engineRegistry:         engineRegistry,
		emitter:                emitter,
		engineLimiters:         engineLimiters,
		ratelimiter:            ratelimiter,
		workflowLimits:         workflowLimits,
		workflowArtifactsStore: workflowArtifacts,
		workflowEncryptionKey:  workflowEncryptionKey,
	}
	eh.engineFactory = eh.engineFactoryFn
	for _, o := range opts {
		o(eh)
	}

	return eh, nil
}

func (h *eventHandler) Close() error {
	es := h.engineRegistry.PopAll()
	return services.MultiCloser(es).Close()
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

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, wfOwner,
			platform.KeyWorkflowTag, payload.WorkflowTag,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
		)

		var err error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), payload.BinaryURL, payload.ConfigURL, err); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()
		err = h.workflowActivatedEvent(ctx, payload)

		if err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow activated event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner),
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

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
			platform.KeyWorkflowTag, payload.Tag,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
		)

		var err error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), payload.BinaryURL, payload.ConfigURL, err); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()

		if err := h.workflowPausedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow paused event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner),
			"workflowTag", payload.Tag, "type", event.Name)
		return nil
	case WorkflowDeleted:
		payload, ok := event.Data.(WorkflowDeletedEvent)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		wfID := payload.WorkflowID.Hex()

		// Get workflow spec from database to get owner and name info for organization lookup
		// Alternative: wire through workflowOwner into the Event, but that requires a lot more surgery
		spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, wfID)
		var wfOwner, wfName, orgID string
		if err != nil {
			// Workflow spec not found, proceed with deletion but without event metadata
			h.lggr.Warnw("Workflow spec not found during deletion, proceeding without org info", "workflowID", wfID, "error", err)
		} else {
			wfOwner = spec.WorkflowOwner
			wfName = spec.WorkflowName
			if wfOwner != "" {
				orgID, err = h.fetchOrganizationID(ctx, wfOwner)
				if err != nil {
					h.lggr.Warnw("Failed to get organization from linking service", "workflowOwner", wfOwner, "error", err)
				}
			}
		}

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, wfName,
			platform.KeyWorkflowOwner, wfOwner,
			platform.KeyOrganizationID, orgID,
			platform.WorkflowRegistryAddress, h.workflowRegistryAddress,
			platform.WorkflowRegistryChainSelector, h.workflowRegistryChainSelector,
		)

		var herr error
		defer func() {
			if err2 := events.EmitWorkflowStatusChangedEventV2(ctx, cma.Labels(), toCommonHead(event.Head), string(event.Name), "", "", herr); err2 != nil {
				h.lggr.Errorf("failed to emit status changed event: %+v", err2)
			}
		}()

		if herr := h.workflowDeletedEvent(ctx, payload); herr != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow deleted event: %v", herr), h.lggr)
			return herr
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "workflowName", wfName, "workflowOwner", wfOwner, "organizationID", orgID, "type", event.Name)
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

// workflowRegisteredEvent handles the WorkflowRegisteredEvent event type.
// This method must remain idempotent and must not error if retried multiple times.
// workflowRegisteredEvent proceeds in two phases:
// - phase 1 synchronizes the database state
// - phase 2 synchronizes the state of the engine registry.
func (h *eventHandler) workflowRegisteredEvent(
	ctx context.Context,
	payload WorkflowRegisteredEvent,
) error {
	status := toSpecStatus(payload.Status)

	// First, let's synchronize the database state.
	// We need to handle three cases:
	// - new registration, without an existing DB record
	// - existing registration that has been updated with new artifacts, and potentially also the status
	// - existing registration that has been updated with a new status
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, payload.WorkflowID.Hex())
	switch {
	case err != nil:
		newSpec, innerErr := h.createWorkflowSpec(ctx, payload)
		if innerErr != nil {
			return innerErr
		}

		spec = newSpec
	case spec.WorkflowID != payload.WorkflowID.Hex():
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
		// This is the happy-path, we're done.
		return nil
	}

	// Any other case ->
	// - engine in registry, but service isn't running
	// - state isn't active
	// Let's clean up and recreate

	cleanupErr := h.tryEngineCleanup(payload.WorkflowID)
	if cleanupErr != nil {
		return fmt.Errorf("could not clean up old engine: %w", cleanupErr)
	}

	return h.tryEngineCreate(ctx, spec)
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
	wfID := payload.WorkflowID.Hex()
	owner := hex.EncodeToString(payload.WorkflowOwner)

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

func (h *eventHandler) engineFactoryFn(ctx context.Context, workflowID string, owner string, name types.WorkflowName, tag string, config []byte, binary []byte) (services.Service, error) {
	lggr := h.lggr.Named("WorkflowEngine.Module").With("workflowID", workflowID, "workflowName", name, "workflowOwner", owner)
	moduleConfig := &host.ModuleConfig{Logger: lggr, Labeler: h.emitter}

	h.lggr.Debugf("Creating module for workflowID %s", workflowID)
	module, err := host.NewModule(moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, fmt.Errorf("could not instantiate module: %w", err)
	}
	h.lggr.Debugf("Finished creating module for workflowID %s", workflowID)

	if module.IsLegacyDAG() { // V1 aka "DAG"
		sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow sdk spec: %w", err)
		}

		// WorkflowRegistry V2 contract does not contain secrets
		emptySecretsFetcher := func(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
			return map[string]string{}, nil
		}

		cfg := workflows.Config{
			Lggr:           h.lggr,
			Workflow:       *sdkSpec,
			WorkflowID:     workflowID,
			WorkflowOwner:  owner, // this gets hex encoded in the engine.
			WorkflowName:   name,
			Registry:       h.capRegistry,
			Store:          h.workflowStore,
			Config:         config,
			Binary:         binary,
			SecretsFetcher: emptySecretsFetcher,
			RateLimiter:    h.ratelimiter,
			WorkflowLimits: h.workflowLimits,

			BillingClient: h.billingClient,
		}
		return workflows.NewEngine(ctx, cfg)
	}

	// V2 aka "NoDAG"
	cfg := &v2.EngineConfig{
		Lggr:                  h.lggr,
		Module:                module,
		WorkflowConfig:        config,
		CapRegistry:           h.capRegistry,
		UseLocalTimeProvider:  h.useLocalTimeProvider,
		DonTimeStore:          h.donTimeStore,
		ExecutionsStore:       h.workflowStore,
		WorkflowID:            workflowID,
		WorkflowOwner:         owner,
		WorkflowName:          name,
		WorkflowTag:           tag,
		WorkflowEncryptionKey: h.workflowEncryptionKey,

		LocalLimits:                       v2.EngineLimits{}, // all defaults
		LocalLimiters:                     h.engineLimiters,
		GlobalExecutionConcurrencyLimiter: h.workflowLimits,
		GlobalExecutionRateLimiter:        h.ratelimiter,

		BeholderEmitter: h.emitter,
		BillingClient:   h.billingClient,

		WorkflowRegistryAddress:       h.workflowRegistryAddress,
		WorkflowRegistryChainSelector: h.workflowRegistryChainSelector,
		OrgResolver:                   h.orgResolver,
	}
	return v2.NewEngine(cfg)
}

// workflowPausedEvent handles the WorkflowPausedEvent event type. This method must remain idempotent.
func (h *eventHandler) workflowPausedEvent(
	ctx context.Context,
	payload WorkflowPausedEvent,
) error {
	return h.workflowDeletedEvent(ctx, WorkflowDeletedEvent{WorkflowID: payload.WorkflowID})
}

// workflowDeletedEvent handles the WorkflowDeletedEvent event type. This method must remain idempotent.
func (h *eventHandler) workflowDeletedEvent(
	ctx context.Context,
	payload WorkflowDeletedEvent,
) error {
	// The order in the handler is slightly different to the order in `tryEngineCleanup`.
	// This is because the engine requires its corresponding DB record to be present to be successfully
	// closed.
	// At the same time, popping the engine should occur last to allow deletes to be retried if any of the
	// prior steps fail.
	e, ok := h.engineRegistry.Get(payload.WorkflowID)
	if ok {
		if innerErr := e.Close(); innerErr != nil {
			return fmt.Errorf("failed to close workflow engine: %w", innerErr)
		}
	}

	if err := h.workflowArtifactsStore.DeleteWorkflowArtifacts(ctx, payload.WorkflowID.Hex()); err != nil {
		return fmt.Errorf("failed to delete workflow artifacts: %w", err)
	}

	_, err := h.engineRegistry.Pop(payload.WorkflowID)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
}

// tryEngineCleanup attempts to stop the workflow engine for the given workflow ID.  Does nothing if the
// workflow engine is not running.
func (h *eventHandler) tryEngineCleanup(workflowID types.WorkflowID) error {
	e, ok := h.engineRegistry.Get(workflowID)
	if ok {
		// Stop the engine
		if err := e.Close(); err != nil {
			return fmt.Errorf("failed to close workflow engine: %w", err)
		}

		// Remove the engine from the registry
		_, err := h.engineRegistry.Pop(workflowID)
		if err != nil {
			return fmt.Errorf("failed to remove workflow engine: %w", err)
		}
	}
	return nil
}

// tryEngineCreate attempts to create a new workflow engine, start it, and register it with the engine registry
func (h *eventHandler) tryEngineCreate(ctx context.Context, spec *job.WorkflowSpec) error {
	// Ensure the capabilities registry is ready before creating any Engine instances.
	// This should be guaranteed by the Workflow Registry Syncer.
	if err := h.ensureCapRegistryReady(ctx); err != nil {
		return fmt.Errorf("failed to ensure capabilities registry is ready: %w", err)
	}

	decodedBinary, err := hex.DecodeString(spec.Workflow)
	if err != nil {
		return fmt.Errorf("failed to decode workflow spec binary: %w", err)
	}

	// Workflow Registry version >2 no longer handles secrets
	secretsURL := ""

	// Before running the engine, handle validations
	// Workflow ID should match what is generated from the stored artifacts
	ownerBytes, err := hex.DecodeString(spec.WorkflowOwner)
	if err != nil {
		return fmt.Errorf("failed to decode owner: %w", err)
	}
	hash, err := pkgworkflows.GenerateWorkflowID(ownerBytes, spec.WorkflowName, decodedBinary, []byte(spec.Config), secretsURL)
	if err != nil {
		return fmt.Errorf("failed to generate workflow id: %w", err)
	}
	wid, err := types.WorkflowIDFromHex(spec.WorkflowID)
	if err != nil {
		return fmt.Errorf("invalid workflow id: %w", err)
	}
	if !types.WorkflowID(hash).Equal(wid) {
		return fmt.Errorf("workflowID mismatch: %x != %x", hash, wid)
	}

	// Start a new WorkflowEngine instance, and add it to local engine registry
	workflowName, err := types.NewWorkflowName(spec.WorkflowName)
	if err != nil {
		return fmt.Errorf("invalid workflow name: %w", err)
	}
	engine, err := h.engineFactory(
		ctx,
		spec.WorkflowID,
		spec.WorkflowOwner,
		workflowName,
		spec.WorkflowTag,
		[]byte(spec.Config),
		decodedBinary,
	)
	if err != nil {
		return fmt.Errorf("failed to create workflow engine: %w", err)
	}

	if err = engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start workflow engine: %w", err)
	}

	if err := h.engineRegistry.Add(wid, engine); err != nil {
		if closeErr := engine.Close(); closeErr != nil {
			return fmt.Errorf("failed to close workflow engine: %w during invariant violation: %w", closeErr, err)
		}
		// This shouldn't happen because we call the handler serially and
		// check for running engines above, see the call to engineRegistry.Contains.
		return fmt.Errorf("invariant violation: %w", err)
	}
	return nil
}

// logCustMsg emits a custom message to the external sink and logs an error if that fails.
func logCustMsg(ctx context.Context, cma custmsg.MessageEmitter, msg string, log logger.Logger) {
	err := cma.Emit(ctx, msg)
	if err != nil {
		log.Helper(1).Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
	}
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
		})
}

func newHandlerTypeError(data any) error {
	return fmt.Errorf("invalid data type %T for event", data)
}
