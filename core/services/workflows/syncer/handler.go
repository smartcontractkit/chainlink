package syncer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// WorkflowRegistryrEventType is the type of event that is emitted by the WorkflowRegistry
type WorkflowRegistryEventType string

var (
	// ForceUpdateSecretsEvent is emitted when a request to force update a workflows secrets is made
	ForceUpdateSecretsEvent WorkflowRegistryEventType = "WorkflowForceUpdateSecretsRequestedV1"

	// WorkflowRegisteredEvent is emitted when a workflow is registered
	WorkflowRegisteredEvent WorkflowRegistryEventType = "WorkflowRegisteredV1"

	// WorkflowUpdatedEvent is emitted when a workflow is updated
	WorkflowUpdatedEvent WorkflowRegistryEventType = "WorkflowUpdatedV1"

	// WorkflowPausedEvent is emitted when a workflow is paused
	WorkflowPausedEvent WorkflowRegistryEventType = "WorkflowPausedV1"

	// WorkflowActivatedEvent is emitted when a workflow is activated
	WorkflowActivatedEvent WorkflowRegistryEventType = "WorkflowActivatedV1"

	// WorkflowDeletedEvent is emitted when a workflow is deleted
	WorkflowDeletedEvent WorkflowRegistryEventType = "WorkflowDeletedV1"
)

type ORM interface {
	artifacts.WorkflowSecretsDS
	artifacts.WorkflowSpecsDS
}

// ForceUpdateSecretsRequestedV1 is a chain agnostic definition of the WorkflowRegistry
// ForceUpdateSecretsRequested event.
type ForceUpdateSecretsRequestedV1 struct {
	SecretsURLHash []byte
	Owner          []byte
	WorkflowName   string
}

type WorkflowRegisteredV1 struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	DonID         uint32
	Status        uint8
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
}

type WorkflowUpdatedV1 struct {
	OldWorkflowID types.WorkflowID
	WorkflowOwner []byte
	DonID         uint32
	NewWorkflowID types.WorkflowID
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
	Status        uint8
}

type WorkflowPausedV1 struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowActivatedV1 struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowDeletedV1 struct {
	WorkflowID    types.WorkflowID
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type engineFactoryFn func(ctx context.Context, wfid string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error)

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding
// method that handles the event.
type eventHandler struct {
	lggr logger.Logger

	workflowStore          store.Store
	capRegistry            core.CapabilitiesRegistry
	engineRegistry         *EngineRegistry
	emitter                custmsg.MessageEmitter
	engineFactory          engineFactoryFn
	ratelimiter            *ratelimiter.RateLimiter
	workflowLimits         *syncerlimiter.Limits
	workflowArtifactsStore WorkflowArtifactsStore
	billingClient          workflows.BillingClient
}

type Event struct {
	EventType WorkflowRegistryEventType
	Data      any
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
		e.engineFactory = func(_ context.Context, _ string, _ string, _ types.WorkflowName, _ []byte, _ []byte) (services.Service, error) {
			return engine, nil
		}
	}
}

func WithBillingClient(client workflows.BillingClient) func(*eventHandler) {
	return func(e *eventHandler) {
		e.billingClient = client
	}
}

type WorkflowArtifactsStore interface {
	FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error)
	GetWorkflowSpec(ctx context.Context, workflowOwner string, workflowName string) (*job.WorkflowSpec, error)
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)
	UpsertWorkflowSpecWithSecrets(ctx context.Context, entry *job.WorkflowSpec, secretsURL, urlHash, secrets string) (int64, error)
	DeleteWorkflowArtifacts(ctx context.Context, workflowOwner string, workflowName string, workflowID string) error

	// Secrets methods
	GetSecrets(ctx context.Context, secretsURL string, WorkflowID [32]byte, WorkflowOwner []byte) ([]byte, error)
	ValidateSecrets(ctx context.Context, workflowID, workflowOwner string) error
	SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error)
	GetSecretsURLHash(workflowOwner []byte, secretsURL []byte) ([]byte, error)
	GetSecretsURLByID(ctx context.Context, id int64) (string, error)

	ForceUpdateSecrets(ctx context.Context, secretsURLHash []byte, owner []byte) (string, error)
}

// NewEventHandler returns a new eventHandler instance.
func NewEventHandler(
	lggr logger.Logger,
	workflowStore store.Store,
	capRegistry core.CapabilitiesRegistry,
	engineRegistry *EngineRegistry,
	emitter custmsg.MessageEmitter,
	ratelimiter *ratelimiter.RateLimiter,
	workflowLimits *syncerlimiter.Limits,
	workflowArtifacts WorkflowArtifactsStore,
	opts ...func(*eventHandler),
) (*eventHandler, error) {
	if engineRegistry == nil {
		return nil, errors.New("engine registry must be provided")
	}

	eh := &eventHandler{
		lggr:                   lggr,
		workflowStore:          workflowStore,
		capRegistry:            capRegistry,
		engineRegistry:         engineRegistry,
		emitter:                emitter,
		ratelimiter:            ratelimiter,
		workflowLimits:         workflowLimits,
		workflowArtifactsStore: workflowArtifacts,
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

func (h *eventHandler) Handle(ctx context.Context, event Event) error {
	switch event.EventType {
	case ForceUpdateSecretsEvent:
		payload, ok := event.Data.(ForceUpdateSecretsRequestedV1)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		cma := h.emitter.With(
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.Owner),
		)

		if _, err := h.workflowArtifactsStore.ForceUpdateSecrets(ctx, payload.SecretsURLHash, payload.Owner); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle force update secrets event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled event", "urlHash", payload.SecretsURLHash, "workflowOwner", hex.EncodeToString(payload.Owner), "type", ForceUpdateSecretsEvent)
		return nil
	case WorkflowRegisteredEvent:
		payload, ok := event.Data.(WorkflowRegisteredV1)
		if !ok {
			return newHandlerTypeError(event.Data)
		}

		wfID := payload.WorkflowID.Hex()

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowRegisteredEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow registered event: %v", err), h.lggr)
			return err
		}

		if err := events.EmitWorkflowStatusChangedEvent(ctx, h.emitter, string(event.EventType)); err != nil {
			h.lggr.Errorf("failed to emit status changed event: %+v", err)
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner), "type", WorkflowRegisteredEvent)
		return nil
	case WorkflowUpdatedEvent:
		payload, ok := event.Data.(WorkflowUpdatedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.Data)
		}

		newWorkflowID := payload.NewWorkflowID.Hex()
		oldWorkflowID := payload.OldWorkflowID.Hex()
		cma := h.emitter.With(
			platform.KeyWorkflowID, newWorkflowID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowUpdatedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow updated event: %v", err), h.lggr)
			return err
		}

		if err := events.EmitWorkflowStatusChangedEvent(ctx, h.emitter, string(event.EventType)); err != nil {
			h.lggr.Errorf("failed to emit status changed event: %+v", err)
		}

		h.lggr.Debugw("handled event", "newWorkflowID", newWorkflowID, "oldWorkflowID", oldWorkflowID, "workflowName", payload.WorkflowName, "workflowOwner", hex.EncodeToString(payload.WorkflowOwner), "type", WorkflowUpdatedEvent)
		return nil
	case WorkflowPausedEvent:
		payload, ok := event.Data.(WorkflowPausedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.Data)
		}

		wfID := payload.WorkflowID.Hex()

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowPausedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow paused event: %v", err), h.lggr)
			return err
		}

		if err := events.EmitWorkflowStatusChangedEvent(ctx, h.emitter, string(event.EventType)); err != nil {
			h.lggr.Errorf("failed to emit status changed event: %+v", err)
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "type", WorkflowPausedEvent)
		return nil
	case WorkflowActivatedEvent:
		payload, ok := event.Data.(WorkflowActivatedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.Data)
		}

		wfID := payload.WorkflowID.Hex()
		wfOwner := hex.EncodeToString(payload.WorkflowOwner)
		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, wfOwner,
		)
		if err := h.workflowActivatedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow activated event: %v", err), h.lggr)
			return err
		}

		if err := events.EmitWorkflowStatusChangedEvent(ctx, h.emitter, string(event.EventType)); err != nil {
			h.lggr.Errorf("failed to emit status changed event: %+v", err)
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "type", WorkflowActivatedEvent, "workflowName", payload.WorkflowName, "workflowOwner", wfOwner)
		return nil
	case WorkflowDeletedEvent:
		payload, ok := event.Data.(WorkflowDeletedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.Data)
		}

		wfID := payload.WorkflowID.Hex()
		wfOwner := hex.EncodeToString(payload.WorkflowOwner)

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, wfOwner,
		)

		if err := h.workflowDeletedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow deleted event: %v", err), h.lggr)
			return err
		}

		if err := events.EmitWorkflowStatusChangedEvent(ctx, h.emitter, string(event.EventType)); err != nil {
			h.lggr.Errorf("failed to emit status changed event: %+v", err)
		}

		h.lggr.Debugw("handled event", "workflowID", wfID, "type", WorkflowDeletedEvent, "workflowName", payload.WorkflowName, "workflowOwner", wfOwner)
		return nil
	default:
		return fmt.Errorf("event type unsupported: %v", event.EventType)
	}
}

// workflowRegisteredEvent handles the WorkflowRegisteredEvent event type.
// This method must remain idempotent and must not error if retried multiple times.
// workflowRegisteredEvent proceeds in two phases:
// - phase 1 synchronizes the database state
// - phase 2 synchronizes the state of the engine registry.
func (h *eventHandler) workflowRegisteredEvent(
	ctx context.Context,
	payload WorkflowRegisteredV1,
) error {
	wfID := payload.WorkflowID.Hex()
	owner := hex.EncodeToString(payload.WorkflowOwner)
	status := toSpecStatus(payload.Status)

	// First, let's synchronize the database state.
	// We need to handle three cases:
	// - new registration, without an existing DB record
	// - existing registration that has been updated with new artifacts, and potentially also the status
	// - existing registration that has been updated with a new status
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, owner, payload.WorkflowName)
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

	// Next, let's synchronize the engine registry.
	// If the state isn't active, we shouldn't have an engine running.
	// Let's try to clean one up if it exists
	if spec.Status != job.WorkflowSpecStatusActive {
		return h.tryEngineCleanup(payload.WorkflowOwner, payload.WorkflowName)
	}

	// We know we need an engine, let's make sure it's the right one.
	// We do this by fetching and comparing whether it's running and that the workflow ID matches
	prevEngine, ok := h.engineRegistry.Get(EngineRegistryKey{Owner: payload.WorkflowOwner, Name: payload.WorkflowName})
	if ok && prevEngine.Ready() == nil && prevEngine.WorkflowID.Hex() == wfID {
		// This is the happy-path, we're done.
		return nil
	}

	// Any other case ->
	// - engine not in workflow registry
	// - engine in registry, but workflow ID does not match
	// - engine in registry, but service isn't running
	// Let's clean up and recreate

	cleanupErr := h.tryEngineCleanup(payload.WorkflowOwner, payload.WorkflowName)
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

func (h *eventHandler) createWorkflowSpec(ctx context.Context, payload WorkflowRegisteredV1) (*job.WorkflowSpec, error) {
	wfID := payload.WorkflowID.Hex()
	owner := hex.EncodeToString(payload.WorkflowOwner)

	decodedBinary, config, err := h.workflowArtifactsStore.FetchWorkflowArtifacts(ctx, wfID, payload.BinaryURL, payload.ConfigURL)
	if err != nil {
		return nil, err
	}

	// Always fetch secrets from the SecretsURL
	var secrets []byte
	if payload.SecretsURL != "" {
		secrets, err = h.workflowArtifactsStore.GetSecrets(ctx, payload.SecretsURL, payload.WorkflowID, payload.WorkflowOwner)
		if err != nil {
			return nil, fmt.Errorf("failed to get secrets: %w", err)
		}
	}

	status := toSpecStatus(payload.Status)

	// Create a new entry in the workflow_spec table corresponding for the new workflow, with the contents of the binaryURL + configURL in the table
	entry := &job.WorkflowSpec{
		Workflow:      hex.EncodeToString(decodedBinary),
		Config:        string(config),
		WorkflowID:    wfID,
		Status:        status,
		WorkflowOwner: owner,
		WorkflowName:  payload.WorkflowName,
		SpecType:      job.WASMFile,
		BinaryURL:     payload.BinaryURL,
		ConfigURL:     payload.ConfigURL,
	}

	secretsURLHash, err := h.workflowArtifactsStore.GetSecretsURLHash(payload.WorkflowOwner, []byte(payload.SecretsURL))
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets URL hash: %w", err)
	}

	if _, err = h.workflowArtifactsStore.UpsertWorkflowSpecWithSecrets(ctx, entry, payload.SecretsURL, hex.EncodeToString(secretsURLHash), string(secrets)); err != nil {
		return nil, fmt.Errorf("failed to upsert workflow spec with secrets: %w", err)
	}

	return entry, nil
}

func (h *eventHandler) engineFactoryFn(ctx context.Context, workflowID string, owner string, name types.WorkflowName, config []byte, binary []byte) (services.Service, error) {
	moduleConfig := &host.ModuleConfig{Logger: h.lggr, Labeler: h.emitter}
	module, err := host.NewModule(moduleConfig, binary, host.WithDeterminism())
	if err != nil {
		return nil, fmt.Errorf("could not instantiate module: %w", err)
	}

	if module.IsLegacyDAG() { // V1 aka "DAG"
		sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow sdk spec: %w", err)
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
			SecretsFetcher: h.workflowArtifactsStore.SecretsFor,
			RateLimiter:    h.ratelimiter,
			WorkflowLimits: h.workflowLimits,
			BillingClient:  h.billingClient,
		}
		return workflows.NewEngine(ctx, cfg)
	}

	// V2 aka "NoDAG"
	cfg := &v2.EngineConfig{
		Lggr:            h.lggr,
		Module:          module,
		CapRegistry:     h.capRegistry,
		ExecutionsStore: h.workflowStore,

		WorkflowID:    workflowID,
		WorkflowOwner: owner,
		WorkflowName:  name,

		LocalLimits:          v2.EngineLimits{}, // all defaults
		GlobalLimits:         h.workflowLimits,
		ExecutionRateLimiter: h.ratelimiter,
	}
	return v2.NewEngine(cfg)
}

// workflowUpdatedEvent handles the WorkflowUpdatedEvent event type by first finding the
// current workflow engine, stopping it, and then starting a new workflow engine with the
// updated workflow spec.
func (h *eventHandler) workflowUpdatedEvent(
	ctx context.Context,
	payload WorkflowUpdatedV1,
) error {
	// Remove the old workflow engine from the local registry if it exists
	if err := h.tryEngineCleanup(payload.WorkflowOwner, payload.WorkflowName); err != nil {
		return err
	}

	registeredEvent := WorkflowRegisteredV1{
		WorkflowID:    payload.NewWorkflowID,
		WorkflowOwner: payload.WorkflowOwner,
		DonID:         payload.DonID,
		Status:        payload.Status,
		WorkflowName:  payload.WorkflowName,
		BinaryURL:     payload.BinaryURL,
		ConfigURL:     payload.ConfigURL,
		SecretsURL:    payload.SecretsURL,
	}

	return h.workflowRegisteredEvent(ctx, registeredEvent)
}

// workflowPausedEvent handles the WorkflowPausedEvent event type.
func (h *eventHandler) workflowPausedEvent(
	ctx context.Context,
	payload WorkflowPausedV1,
) error {
	// Remove the workflow engine from the local registry if it exists
	if err := h.tryEngineCleanup(payload.WorkflowOwner, payload.WorkflowName); err != nil {
		return err
	}

	// get existing workflow spec
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow spec: %w", err)
	}

	// update the status of the workflow spec
	spec.Status = job.WorkflowSpecStatusPaused
	if _, err := h.workflowArtifactsStore.UpsertWorkflowSpec(ctx, spec); err != nil {
		return fmt.Errorf("failed to update workflow spec: %w", err)
	}

	return nil
}

// workflowActivatedEvent handles the WorkflowActivatedEvent event type.
func (h *eventHandler) workflowActivatedEvent(
	ctx context.Context,
	payload WorkflowActivatedV1,
) error {
	// fetch the workflow spec from the DB
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow spec: %w", err)
	}

	// Do nothing if the workflow is already active
	key := EngineRegistryKey{
		Owner: payload.WorkflowOwner, Name: payload.WorkflowName,
	}
	if spec.Status == job.WorkflowSpecStatusActive && h.engineRegistry.Contains(key) {
		return nil
	}

	// get the secrets url by the secrets id
	secretsURL, err := h.workflowArtifactsStore.GetSecretsURLByID(ctx, spec.SecretsID.Int64)
	if err != nil {
		return fmt.Errorf("failed to get secrets URL by ID: %w", err)
	}

	// start a new workflow engine
	registeredEvent := WorkflowRegisteredV1{
		WorkflowID:    payload.WorkflowID,
		WorkflowOwner: payload.WorkflowOwner,
		DonID:         payload.DonID,
		Status:        0,
		WorkflowName:  payload.WorkflowName,
		BinaryURL:     spec.BinaryURL,
		ConfigURL:     spec.ConfigURL,
		SecretsURL:    secretsURL,
	}

	return h.workflowRegisteredEvent(ctx, registeredEvent)
}

// workflowDeletedEvent handles the WorkflowDeletedEvent event type. This method must remain idempotent.
func (h *eventHandler) workflowDeletedEvent(
	ctx context.Context,
	payload WorkflowDeletedV1,
) error {
	// The order in the handler is slightly different to the order in `tryEngineCleanup`.
	// This is because the engine requires its corresponding DB record to be present to be successfully
	// closed.
	// At the same time, popping the engine should occur last to allow deletes to be retried if any of the
	// prior steps fail.
	key := EngineRegistryKey{Owner: payload.WorkflowOwner, Name: payload.WorkflowName}

	e, ok := h.engineRegistry.Get(key)
	if ok {
		if innerErr := e.Close(); innerErr != nil {
			return fmt.Errorf("failed to close workflow engine: %w", innerErr)
		}
	}

	if err := h.workflowArtifactsStore.DeleteWorkflowArtifacts(ctx, hex.EncodeToString(payload.WorkflowOwner),
		payload.WorkflowName, payload.WorkflowID.Hex()); err != nil {
		return fmt.Errorf("failed to delete workflow artifacts: %w", err)
	}

	_, err := h.engineRegistry.Pop(key)
	if errors.Is(err, errNotFound) {
		return nil
	}
	return err
}

// tryEngineCleanup attempts to stop the workflow engine for the given workflow ID.  Does nothing if the
// workflow engine is not running.
func (h *eventHandler) tryEngineCleanup(workflowOwner []byte, workflowName string) error {
	key := EngineRegistryKey{
		Owner: workflowOwner, Name: workflowName,
	}
	if h.engineRegistry.Contains(key) {
		// This shouldn't error since we just checked that the key existed above.
		e, ok := h.engineRegistry.Get(key)
		if !ok {
			return errors.New("invariant violation: failed to get workflow engine")
		}

		// Stop the engine
		if err := e.Close(); err != nil {
			return fmt.Errorf("failed to close workflow engine: %w", err)
		}

		// Remove the engine from the registry
		_, err := h.engineRegistry.Pop(key)
		if err != nil {
			return fmt.Errorf("failed to remove workflow engine: %w", err)
		}
	}
	return nil
}

// tryEngineCreate attempts to create a new workflow engine, start it, and register it with the engine registry
func (h *eventHandler) tryEngineCreate(ctx context.Context, spec *job.WorkflowSpec) error {
	decodedBinary, err := hex.DecodeString(spec.Workflow)
	if err != nil {
		return fmt.Errorf("failed to decode workflow spec binary: %w", err)
	}

	secretsURL := ""
	if spec.SecretsID.Valid {
		secretsURL, err = h.workflowArtifactsStore.GetSecretsURLByID(ctx, spec.SecretsID.Int64)
		if err != nil {
			return err
		}
	}

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

	// Secrets should be valid
	err = h.workflowArtifactsStore.ValidateSecrets(ctx, spec.WorkflowID, spec.WorkflowOwner)
	if err != nil {
		return err
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
		[]byte(spec.Config),
		decodedBinary,
	)
	if err != nil {
		return fmt.Errorf("failed to create workflow engine: %w", err)
	}

	if err = engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start workflow engine: %w", err)
	}

	if err := h.engineRegistry.Add(EngineRegistryKey{Owner: ownerBytes, Name: spec.WorkflowName}, engine, wid); err != nil {
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

func newHandlerTypeError(data any) error {
	return fmt.Errorf("invalid data type %T for event", data)
}
