package syncer

import (
	"bytes"
	"context"

	"encoding/hex"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
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

// WorkflowRegistryForceUpdateSecretsRequestedV1 is a chain agnostic definition of the WorkflowRegistry
// ForceUpdateSecretsRequested event.
type WorkflowRegistryForceUpdateSecretsRequestedV1 struct {
	SecretsURLHash []byte
	Owner          []byte
	WorkflowName   string
}

type WorkflowRegistryWorkflowRegisteredV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	Status        uint8
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
}

type WorkflowRegistryWorkflowUpdatedV1 struct {
	OldWorkflowID [32]byte
	WorkflowOwner []byte
	DonID         uint32
	NewWorkflowID [32]byte
	WorkflowName  string
	BinaryURL     string
	ConfigURL     string
	SecretsURL    string
}

type WorkflowRegistryWorkflowPausedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowRegistryWorkflowActivatedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type WorkflowRegistryWorkflowDeletedV1 struct {
	WorkflowID    [32]byte
	WorkflowOwner []byte
	DonID         uint32
	WorkflowName  string
}

type engineFactoryFn func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error)

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
}

type Event interface {
	GetEventType() WorkflowRegistryEventType
	GetData() any
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

type WorkflowArtifactsStore interface {
	FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error)
	GetWorkflowSpec(ctx context.Context, workflowOwner string, workflowName string) (*job.WorkflowSpec, error)
	UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error)
	UpsertWorkflowSpecWithSecrets(ctx context.Context, entry *job.WorkflowSpec, secretsURL, urlHash, secrets string) (int64, error)
	DeleteWorkflowArtifacts(ctx context.Context, workflowOwner string, workflowName string, workflowID string) error

	// Secrets methods
	GetSecrets(ctx context.Context, secretsURL string, WorkflowID [32]byte, WorkflowOwner []byte) ([]byte, error)
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
	emitter custmsg.MessageEmitter,
	ratelimiter *ratelimiter.RateLimiter,
	workflowLimits *syncerlimiter.Limits,
	workflowArtifacts WorkflowArtifactsStore,
	opts ...func(*eventHandler),
) *eventHandler {
	eh := &eventHandler{
		lggr:                   lggr,
		workflowStore:          workflowStore,
		capRegistry:            capRegistry,
		engineRegistry:         NewEngineRegistry(),
		emitter:                emitter,
		ratelimiter:            ratelimiter,
		workflowLimits:         workflowLimits,
		workflowArtifactsStore: workflowArtifacts,
	}
	eh.engineFactory = eh.engineFactoryFn
	for _, o := range opts {
		o(eh)
	}

	return eh
}

func (h *eventHandler) Close() error {
	return services.MultiCloser(h.engineRegistry.PopAll()).Close()
}

func (h *eventHandler) Handle(ctx context.Context, event Event) error {
	switch event.GetEventType() {
	case ForceUpdateSecretsEvent:
		payload, ok := event.GetData().(WorkflowRegistryForceUpdateSecretsRequestedV1)
		if !ok {
			return newHandlerTypeError(event.GetData())
		}

		cma := h.emitter.With(
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.Owner),
		)

		if _, err := h.workflowArtifactsStore.ForceUpdateSecrets(ctx, payload.SecretsURLHash, payload.Owner); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle force update secrets event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled force update secrets events for URL hash", "urlHash", payload.SecretsURLHash)
		return nil
	case WorkflowRegisteredEvent:
		payload, ok := event.GetData().(WorkflowRegistryWorkflowRegisteredV1)
		if !ok {
			return newHandlerTypeError(event.GetData())
		}
		wfID := hex.EncodeToString(payload.WorkflowID[:])

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowRegisteredEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow registered event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled workflow registration event", "workflowID", wfID)
		return nil
	case WorkflowUpdatedEvent:
		payload, ok := event.GetData().(WorkflowRegistryWorkflowUpdatedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.GetData())
		}

		newWorkflowID := hex.EncodeToString(payload.NewWorkflowID[:])
		cma := h.emitter.With(
			platform.KeyWorkflowID, newWorkflowID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowUpdatedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow updated event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled workflow updated event", "workflowID", newWorkflowID)
		return nil
	case WorkflowPausedEvent:
		payload, ok := event.GetData().(WorkflowRegistryWorkflowPausedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.GetData())
		}

		wfID := hex.EncodeToString(payload.WorkflowID[:])

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowPausedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow paused event: %v", err), h.lggr)
			return err
		}
		h.lggr.Debugw("handled workflow paused event", "workflowID", wfID)
		return nil
	case WorkflowActivatedEvent:
		payload, ok := event.GetData().(WorkflowRegistryWorkflowActivatedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.GetData())
		}

		wfID := hex.EncodeToString(payload.WorkflowID[:])

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)
		if err := h.workflowActivatedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow activated event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled workflow activated event", "workflowID", wfID)
		return nil
	case WorkflowDeletedEvent:
		payload, ok := event.GetData().(WorkflowRegistryWorkflowDeletedV1)
		if !ok {
			return fmt.Errorf("invalid data type %T for event", event.GetData())
		}

		wfID := hex.EncodeToString(payload.WorkflowID[:])

		cma := h.emitter.With(
			platform.KeyWorkflowID, wfID,
			platform.KeyWorkflowName, payload.WorkflowName,
			platform.KeyWorkflowOwner, hex.EncodeToString(payload.WorkflowOwner),
		)

		if err := h.workflowDeletedEvent(ctx, payload); err != nil {
			logCustMsg(ctx, cma, fmt.Sprintf("failed to handle workflow deleted event: %v", err), h.lggr)
			return err
		}

		h.lggr.Debugw("handled workflow deleted event", "workflowID", wfID)
		return nil
	default:
		return fmt.Errorf("event type unsupported: %v", event.GetEventType())
	}
}

type workflowName struct {
	name string
}

func (w workflowName) String() string {
	return w.name
}

func (w workflowName) Hex() string {
	// Internal workflow names must not exceed 10 bytes for workflow engine and on-chain use.
	// A name is used internally that is first hashed to avoid collisions,
	// hex encoded to ensure UTF8 encoding, then truncated to 10 bytes.
	truncatedName := pkgworkflows.HashTruncateName(w.name)
	hexName := hex.EncodeToString([]byte(truncatedName))
	return hexName
}

// workflowRegisteredEvent handles the WorkflowRegisteredEvent event type.
func (h *eventHandler) workflowRegisteredEvent(
	ctx context.Context,
	payload WorkflowRegistryWorkflowRegisteredV1,
) error {
	// Fetch the workflow artifacts from the database or download them from the specified URLs
	decodedBinary, config, err := h.getWorkflowArtifacts(ctx, payload)
	if err != nil {
		return err
	}

	// Always fetch secrets from the SecretsURL
	var secrets []byte
	if payload.SecretsURL != "" {
		secrets, err = h.workflowArtifactsStore.GetSecrets(ctx, payload.SecretsURL, payload.WorkflowID, payload.WorkflowOwner)
		if err != nil {
			return fmt.Errorf("failed to get secrets: %w", err)
		}
	}

	// Calculate the hash of the binary and config files
	hash, err := pkgworkflows.GenerateWorkflowID(payload.WorkflowOwner, payload.WorkflowName, decodedBinary, config, payload.SecretsURL)
	if err != nil {
		return fmt.Errorf("failed to generate workflow id: %w", err)
	}

	// Pre-check: verify that the workflowID matches; if it doesn’t abort and log an error via Beholder.
	if !bytes.Equal(hash[:], payload.WorkflowID[:]) {
		return fmt.Errorf("workflowID mismatch: %x != %x", hash, payload.WorkflowID)
	}

	// Ensure that there is no running workflow engine for the given workflow ID.
	if h.engineRegistry.Contains(hex.EncodeToString(payload.WorkflowID[:])) {
		return fmt.Errorf("workflow is already running, so not starting it : %s", hex.EncodeToString(payload.WorkflowID[:]))
	}

	// Save the workflow secrets
	urlHash, err := h.workflowArtifactsStore.GetSecretsURLHash(payload.WorkflowOwner, []byte(payload.SecretsURL))
	if err != nil {
		return fmt.Errorf("failed to get secrets URL hash: %w", err)
	}

	// Create a new entry in the workflow_spec table corresponding for the new workflow, with the contents of the binaryURL + configURL in the table
	status := job.WorkflowSpecStatusActive
	if payload.Status == 1 {
		status = job.WorkflowSpecStatusPaused
	}

	wfID := hex.EncodeToString(payload.WorkflowID[:])
	owner := hex.EncodeToString(payload.WorkflowOwner)
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

	if _, err = h.workflowArtifactsStore.UpsertWorkflowSpecWithSecrets(ctx, entry, payload.SecretsURL, hex.EncodeToString(urlHash), string(secrets)); err != nil {
		return fmt.Errorf("failed to upsert workflow spec with secrets: %w", err)
	}

	if status != job.WorkflowSpecStatusActive {
		h.lggr.Debugw("workflow is marked as paused, so not starting it", "workflow", wfID)
		return nil
	}

	// If status == active, start a new WorkflowEngine instance, and add it to local engine registry
	engine, err := h.engineFactory(
		ctx,
		wfID,
		owner,
		workflowName{
			name: payload.WorkflowName,
		},
		config,
		decodedBinary,
	)
	if err != nil {
		return fmt.Errorf("failed to create workflow engine: %w", err)
	}

	if err := engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start workflow engine: %w", err)
	}

	// This shouldn't fail because we call the handler serially and
	// check for running engines above, see the call to engineRegistry.Contains.
	if err := h.engineRegistry.Add(wfID, engine); err != nil {
		return fmt.Errorf("invariant violation: %w", err)
	}

	return nil
}

// getWorkflowArtifacts retrieves the workflow artifacts from the database if they exist,
// or downloads them from the specified URLs if they are not found in the database.
func (h *eventHandler) getWorkflowArtifacts(
	ctx context.Context,
	payload WorkflowRegistryWorkflowRegisteredV1,
) ([]byte, []byte, error) {
	workflowID := hex.EncodeToString(payload.WorkflowID[:])

	return h.workflowArtifactsStore.FetchWorkflowArtifacts(ctx, workflowID, payload.BinaryURL, payload.ConfigURL)
}

func (h *eventHandler) engineFactoryFn(ctx context.Context, workflowID string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
	moduleConfig := &host.ModuleConfig{Logger: h.lggr, Labeler: h.emitter}
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
	}
	return workflows.NewEngine(ctx, cfg)
}

// workflowUpdatedEvent handles the WorkflowUpdatedEvent event type by first finding the
// current workflow engine, stopping it, and then starting a new workflow engine with the
// updated workflow spec.
func (h *eventHandler) workflowUpdatedEvent(
	ctx context.Context,
	payload WorkflowRegistryWorkflowUpdatedV1,
) error {
	// Remove the old workflow engine from the local registry if it exists
	if err := h.tryEngineCleanup(hex.EncodeToString(payload.OldWorkflowID[:])); err != nil {
		return err
	}

	registeredEvent := WorkflowRegistryWorkflowRegisteredV1{
		WorkflowID:    payload.NewWorkflowID,
		WorkflowOwner: payload.WorkflowOwner,
		DonID:         payload.DonID,
		Status:        0,
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
	payload WorkflowRegistryWorkflowPausedV1,
) error {
	// Remove the workflow engine from the local registry if it exists
	if err := h.tryEngineCleanup(hex.EncodeToString(payload.WorkflowID[:])); err != nil {
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
	payload WorkflowRegistryWorkflowActivatedV1,
) error {
	// fetch the workflow spec from the DB
	spec, err := h.workflowArtifactsStore.GetWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow spec: %w", err)
	}

	// Do nothing if the workflow is already active
	if spec.Status == job.WorkflowSpecStatusActive && h.engineRegistry.Contains(hex.EncodeToString(payload.WorkflowID[:])) {
		return nil
	}

	// get the secrets url by the secrets id
	secretsURL, err := h.workflowArtifactsStore.GetSecretsURLByID(ctx, spec.SecretsID.Int64)
	if err != nil {
		return fmt.Errorf("failed to get secrets URL by ID: %w", err)
	}

	// start a new workflow engine
	registeredEvent := WorkflowRegistryWorkflowRegisteredV1{
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

// workflowDeletedEvent handles the WorkflowDeletedEvent event type.
func (h *eventHandler) workflowDeletedEvent(
	ctx context.Context,
	payload WorkflowRegistryWorkflowDeletedV1,
) error {
	workflowID := hex.EncodeToString(payload.WorkflowID[:])

	if err := h.tryEngineCleanup(workflowID); err != nil {
		return err
	}

	if err := h.workflowArtifactsStore.DeleteWorkflowArtifacts(ctx, hex.EncodeToString(payload.WorkflowOwner),
		payload.WorkflowName, workflowID); err != nil {
		return fmt.Errorf("failed to delete workflow artifacts: %w", err)
	}

	return nil
}

// tryEngineCleanup attempts to stop the workflow engine for the given workflow ID.  Does nothing if the
// workflow engine is not running.
func (h *eventHandler) tryEngineCleanup(wfID string) error {
	if h.engineRegistry.Contains(wfID) {
		// Remove the engine from the registry
		e, err := h.engineRegistry.Pop(wfID)
		if err != nil {
			return fmt.Errorf("failed to get workflow engine: %w", err)
		}

		// Stop the engine
		if err := e.Close(); err != nil {
			return fmt.Errorf("failed to close workflow engine: %w", err)
		}
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
