package syncer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var ErrNotImplemented = errors.New("not implemented")

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

type lastFetchedAtMap struct {
	m map[string]time.Time
	sync.RWMutex
}

func (l *lastFetchedAtMap) Set(url string, at time.Time) {
	l.Lock()
	defer l.Unlock()
	l.m[url] = at
}

func (l *lastFetchedAtMap) Get(url string) (time.Time, bool) {
	l.RLock()
	defer l.RUnlock()
	got, ok := l.m[url]
	return got, ok
}

func newLastFetchedAtMap() *lastFetchedAtMap {
	return &lastFetchedAtMap{
		m: map[string]time.Time{},
	}
}

type engineFactoryFn func(ctx context.Context, wfid string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error)

type ArtifactConfig struct {
	MaxConfigSize  uint64
	MaxSecretsSize uint64
	MaxBinarySize  uint64
}

// By default, if type is unknown, the largest artifact size is 26.4KB.  Configure the artifact size
// via the ArtifactConfig to override this default.
const defaultMaxArtifactSizeBytes = 26.4 * utils.KB

func (cfg *ArtifactConfig) ApplyDefaults() {
	if cfg.MaxConfigSize == 0 {
		cfg.MaxConfigSize = defaultMaxArtifactSizeBytes
	}
	if cfg.MaxSecretsSize == 0 {
		cfg.MaxSecretsSize = defaultMaxArtifactSizeBytes
	}
	if cfg.MaxBinarySize == 0 {
		cfg.MaxBinarySize = defaultMaxArtifactSizeBytes
	}
}

// eventHandler is a handler for WorkflowRegistryEvent events.  Each event type has a corresponding
// method that handles the event.
type eventHandler struct {
	lggr logger.Logger
	orm  WorkflowRegistryDS

	// fetchFn is a function that fetches the contents of a URL with a limit on the size of the response.
	fetchFn FetcherFunc

	// limits sets max artifact sizes to fetch when handling events
	limits                   *ArtifactConfig
	workflowStore            store.Store
	capRegistry              core.CapabilitiesRegistry
	engineRegistry           *EngineRegistry
	emitter                  custmsg.MessageEmitter
	lastFetchedAtMap         *lastFetchedAtMap
	clock                    clockwork.Clock
	secretsFreshnessDuration time.Duration
	encryptionKey            workflowkey.Key
	engineFactory            engineFactoryFn
	ratelimiter              *ratelimiter.RateLimiter
	workflowLimits           *syncerlimiter.Limits
	decryptSecrets           func(data []byte, owner string) (map[string]string, error)
}

type Event interface {
	GetEventType() WorkflowRegistryEventType
	GetData() any
}

var defaultSecretsFreshnessDuration = 24 * time.Hour

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

func WithMaxArtifactSize(cfg ArtifactConfig) func(*eventHandler) {
	return func(eh *eventHandler) {
		eh.limits = &cfg
	}
}

// NewEventHandler returns a new eventHandler instance.
func NewEventHandler(
	lggr logger.Logger,
	orm ORM,
	fetchFn FetcherFunc,
	workflowStore store.Store,
	capRegistry core.CapabilitiesRegistry,
	emitter custmsg.MessageEmitter,
	clock clockwork.Clock,
	encryptionKey workflowkey.Key,
	ratelimiter *ratelimiter.RateLimiter,
	workflowLimits *syncerlimiter.Limits,
	opts ...func(*eventHandler),
) *eventHandler {
	eh := &eventHandler{
		lggr:                     lggr,
		orm:                      orm,
		workflowStore:            workflowStore,
		capRegistry:              capRegistry,
		fetchFn:                  fetchFn,
		engineRegistry:           NewEngineRegistry(),
		emitter:                  emitter,
		lastFetchedAtMap:         newLastFetchedAtMap(),
		clock:                    clock,
		limits:                   &ArtifactConfig{},
		secretsFreshnessDuration: defaultSecretsFreshnessDuration,
		encryptionKey:            encryptionKey,
		ratelimiter:              ratelimiter,
		workflowLimits:           workflowLimits,
		decryptSecrets: func(data []byte, owner string) (map[string]string, error) {
			secretsPayload := secrets.EncryptedSecretsResult{}
			err := json.Unmarshal(data, &secretsPayload)
			if err != nil {
				return nil, fmt.Errorf("could not unmarshal secrets: %w", err)
			}

			return secrets.DecryptSecretsForNode(secretsPayload, encryptionKey, owner)
		},
	}
	eh.engineFactory = eh.engineFactoryFn
	eh.limits.ApplyDefaults()

	for _, o := range opts {
		o(eh)
	}

	return eh
}

func (h *eventHandler) refreshSecrets(ctx context.Context, workflowOwner, workflowName, workflowID, secretsURLHash string) (string, error) {
	owner, err := hex.DecodeString(workflowOwner)
	if err != nil {
		return "", err
	}

	decodedHash, err := hex.DecodeString(secretsURLHash)
	if err != nil {
		return "", err
	}

	updatedSecrets, err := h.forceUpdateSecretsEvent(
		ctx,
		WorkflowRegistryForceUpdateSecretsRequestedV1{
			SecretsURLHash: decodedHash,
			Owner:          owner,
			WorkflowName:   name,
		},
	)
	if err != nil {
		return "", err
	}

	return updatedSecrets, nil
}

func (h *eventHandler) SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	secretsURLHash, secretsPayload, err := h.orm.GetContentsByWorkflowID(ctx, workflowID)
	if err != nil {
		// The workflow record was found, but secrets_id was empty.
		// Let's just stub out the response.
		if errors.Is(err, ErrEmptySecrets) {
			return map[string]string{}, nil
		}

		return nil, fmt.Errorf("failed to fetch secrets by workflow ID: %w", err)
	}

	lastFetchedAt, ok := h.lastFetchedAtMap.Get(secretsURLHash)
	if !ok || h.clock.Now().Sub(lastFetchedAt) > h.secretsFreshnessDuration {
		updatedSecrets, innerErr := h.refreshSecrets(ctx, workflowOwner, hexWorkflowName, workflowID, secretsURLHash)
		if innerErr != nil {
			msg := fmt.Sprintf("could not refresh secrets: proceeding with stale secrets for workflowID %s: %s", workflowID, innerErr)
			h.lggr.Error(msg)

			logCustMsg(
				ctx,
				h.emitter.With(
					platform.KeyWorkflowID, workflowID,
					platform.KeyWorkflowName, decodedWorkflowName,
					platform.KeyWorkflowOwner, workflowOwner,
				),
				msg,
				h.lggr,
			)
		} else {
			secretsPayload = updatedSecrets
		}
	}

	return h.decryptSecrets([]byte(secretsPayload), workflowOwner)
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

		if _, err := h.forceUpdateSecretsEvent(ctx, payload); err != nil {
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

func messageID(url string, parts ...string) string {
	h := sha256.New()
	h.Write([]byte(url))
	for _, p := range parts {
		h.Write([]byte(p))
	}
	hash := hex.EncodeToString(h.Sum(nil))
	p := []string{ghcapabilities.MethodWorkflowSyncer, hash}
	return strings.Join(p, "/")
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
		wid := hex.EncodeToString(payload.WorkflowID[:])
		req := ghcapabilities.Request{
			URL:              payload.SecretsURL,
			Method:           http.MethodGet,
			MaxResponseBytes: safeUint32(h.limits.MaxSecretsSize),
			WorkflowID:       wid,
		}
		fetchedSecrets, fetchErr := h.fetchFn(ctx, messageID(payload.SecretsURL, wid), req)
		if fetchErr != nil {
			return fmt.Errorf("failed to fetch secrets from %s : %w", payload.SecretsURL, fetchErr)
		}

		// sanity check by decoding the secrets
		_, decryptErr := h.decryptSecrets(fetchedSecrets, hex.EncodeToString(payload.WorkflowOwner))
		if decryptErr != nil {
			return fmt.Errorf("failed to decrypt secrets %s: %w", payload.SecretsURL, decryptErr)
		}
		secrets = fetchedSecrets
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
	if h.engineRegistry.IsRunning(hex.EncodeToString(payload.WorkflowID[:])) {
		return fmt.Errorf("workflow is already running, so not starting it : %s", hex.EncodeToString(payload.WorkflowID[:]))
	}

	// Save the workflow secrets
	urlHash, err := h.orm.GetSecretsURLHash(payload.WorkflowOwner, []byte(payload.SecretsURL))
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

	if _, err = h.orm.UpsertWorkflowSpecWithSecrets(ctx, entry, payload.SecretsURL, hex.EncodeToString(urlHash), string(secrets)); err != nil {
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

	// This shouldn't happen because we call the handler serially and
	// check for running engines above, see the call to engineRegistry.IsRunning.
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
	// Check if the workflow spec is already stored in the database
	if spec, err := h.orm.GetWorkflowSpecByID(ctx, hex.EncodeToString(payload.WorkflowID[:])); err == nil {
		// there is no update in the BinaryURL or ConfigURL, lets decode the stored artifacts
		decodedBinary, err := hex.DecodeString(spec.Workflow)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode stored workflow spec: %w", err)
		}
		return decodedBinary, []byte(spec.Config), nil
	}

	// Fetch the binary and config files from the specified URLs.
	var (
		binary, decodedBinary, config []byte
		err                           error
	)

	wid := hex.EncodeToString(payload.WorkflowID[:])
	req := ghcapabilities.Request{
		URL:              payload.BinaryURL,
		Method:           http.MethodGet,
		MaxResponseBytes: safeUint32(h.limits.MaxBinarySize),
		WorkflowID:       wid,
	}
	binary, err = h.fetchFn(ctx, messageID(payload.BinaryURL, wid), req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch binary from %s : %w", payload.BinaryURL, err)
	}

	if decodedBinary, err = base64.StdEncoding.DecodeString(string(binary)); err != nil {
		return nil, nil, fmt.Errorf("failed to decode binary: %w", err)
	}

	if payload.ConfigURL != "" {
		req := ghcapabilities.Request{
			URL:              payload.ConfigURL,
			Method:           http.MethodGet,
			MaxResponseBytes: safeUint32(h.limits.MaxConfigSize),
			WorkflowID:       wid,
		}
		config, err = h.fetchFn(ctx, messageID(payload.ConfigURL, wid), req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch config from %s : %w", payload.ConfigURL, err)
		}
	}
	return decodedBinary, config, nil
}

func (h *eventHandler) engineFactoryFn(ctx context.Context, id string, owner string, name workflows.WorkflowNamer, config []byte, binary []byte) (services.Service, error) {
	moduleConfig := &host.ModuleConfig{Logger: h.lggr, Labeler: h.emitter}
	sdkSpec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow sdk spec: %w", err)
	}

	cfg := workflows.Config{
		Lggr:           h.lggr,
		Workflow:       *sdkSpec,
		WorkflowID:     id,
		WorkflowOwner:  owner, // this gets hex encoded in the engine.
		WorkflowName:   name,
		Registry:       h.capRegistry,
		Store:          h.workflowStore,
		Config:         config,
		Binary:         binary,
		SecretsFetcher: h,
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

	// get existing workflow spec from DB
	spec, err := h.orm.GetWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow spec: %w", err)
	}

	// update the status of the workflow spec
	spec.Status = job.WorkflowSpecStatusPaused
	if _, err := h.orm.UpsertWorkflowSpec(ctx, spec); err != nil {
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
	spec, err := h.orm.GetWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow spec: %w", err)
	}

	// Do nothing if the workflow is already active
	if spec.Status == job.WorkflowSpecStatusActive && h.engineRegistry.IsRunning(hex.EncodeToString(payload.WorkflowID[:])) {
		return nil
	}

	// get the secrets url by the secrets id
	secretsURL, err := h.orm.GetSecretsURLByID(ctx, spec.SecretsID.Int64)
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
	if err := h.tryEngineCleanup(hex.EncodeToString(payload.WorkflowID[:])); err != nil {
		return err
	}

	err := h.orm.DeleteWorkflowSpec(ctx, hex.EncodeToString(payload.WorkflowOwner), payload.WorkflowName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.lggr.Warnw("workflow spec not found", "workflowID", hex.EncodeToString(payload.WorkflowID[:]))
			return nil
		}
		return fmt.Errorf("failed to delete workflow spec: %w", err)
	}

	return nil
}

// forceUpdateSecretsEvent handles the ForceUpdateSecretsEvent event type.
func (h *eventHandler) forceUpdateSecretsEvent(
	ctx context.Context,
	payload WorkflowRegistryForceUpdateSecretsRequestedV1,
) (string, error) {
	// Get the URL of the secrets file from the event data
	hash := hex.EncodeToString(payload.SecretsURLHash)

	url, err := h.orm.GetSecretsURLByHash(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("failed to get URL by hash %s : %w", hash, err)
	}

	ownerHex := hex.EncodeToString(payload.Owner)
	req := ghcapabilities.Request{
		URL:              url,
		Method:           http.MethodGet,
		MaxResponseBytes: safeUint32(h.limits.MaxSecretsSize),
		// TODO -- fix, but this is used for rate limiting purposes
		WorkflowID: hex.EncodeToString(payload.Owner),
	}
	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetchFn(ctx, messageID(url, ownerHex), req)
	if err != nil {
		return "", err
	}

	// Sanity check the payload and ensure we can decrypt it.
	// If we can't, let's return an error and we won't store the result in the DB.
	_, err = h.decryptSecrets(secrets, hex.EncodeToString(payload.Owner))
	if err != nil {
		return "", fmt.Errorf("failed to validate secrets: could not decrypt: %w", err)
	}

	h.lastFetchedAtMap.Set(hash, h.clock.Now())

	// Update the secrets in the ORM
	if _, err := h.orm.Update(ctx, hash, string(secrets)); err != nil {
		return "", fmt.Errorf("failed to update secrets: %w", err)
	}

	return string(secrets), nil
}

// tryEngineCleanup attempts to stop the workflow engine for the given workflow ID.  Does nothing if the
// workflow engine is not running.
func (h *eventHandler) tryEngineCleanup(wfID string) error {
	if h.engineRegistry.IsRunning(wfID) {
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

func safeUint32(n uint64) uint32 {
	if n > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(n)
}
