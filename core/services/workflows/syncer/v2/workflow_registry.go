package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/big"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/versioning"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval          = 12 * time.Second
	defaultRetryInterval         = 12 * time.Second
	defaultMaxRetryInterval      = 5 * time.Minute
	defaultMaxConcurrency        = 12
	WorkflowRegistryContractName = "WorkflowRegistry"

	GetWorkflowsByDONMethodName                   = "getWorkflowListByDON"
	GetActiveAllowlistedRequestsReverseMethodName = "getActiveAllowlistedRequestsReverse"
	TotalAllowlistedRequestsMethodName            = "totalAllowlistedRequests"

	defaultTickIntervalForAllowlistedRequests = 5 * time.Second

	// MaxResultsPerQuery defines the maximum number of results that can be queried in a single request.
	// The default value of 1,000 was chosen based on expected system performance and typical use cases.
	MaxResultsPerQuery = int64(1_000)
)

// WorkflowRegistrySyncer is the public interface of the package.
type WorkflowRegistrySyncer interface {
	services.Service

	// GetAllowlistedRequests returns the latest list of allowlisted requests. This list is fetched periodically
	// from the workflow registry contract.
	GetAllowlistedRequests(ctx context.Context) []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
}

// workflowRegistry is the implementation of the WorkflowRegistrySyncer interface.
type workflowRegistry struct {
	services.StateMachine

	// close stopCh to stop the workflowRegistry.
	stopCh services.StopChan

	// all goroutines are waited on with wg.
	wg sync.WaitGroup

	// ticker is the interval at which the workflowRegistry will
	// poll the contract for events, and poll the contract for the latest workflow metadata.
	ticker <-chan time.Time

	lggr                    logger.Logger
	workflowRegistryAddress string
	chainSelector           string

	// lastSeenAllowlistedRequestsCount tracks the last seen allowlisted requests count to avoid fetching the same allowlisted requests multiple times.
	// This value is stored in memory and not persisted to the database.
	lastSeenAllowlistedRequestsCount *big.Int
	allowListedRequests              []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	allowListedMu                    sync.RWMutex

	contractReaderFn versioning.ContractReaderFactory

	// contractReader is used exclusively for fetching allowlisted requests from the WorkflowRegistry
	// contract. This data is consumed by Vault DON nodes to authorize incoming vault requests.
	// Workflow metadata is fetched separately via workflowSources (see below).
	contractReader types.ContractReader

	// workflowSources holds workflow metadata sources (contract, file, gRPC).
	workflowSources []WorkflowMetadataSource

	config Config

	handler evtHandler

	workflowDonNotifier donNotifier

	metrics *metrics

	engineRegistry *EngineRegistry

	retryInterval    time.Duration
	maxRetryInterval time.Duration
	maxConcurrency   int
	clock            clockwork.Clock

	hooks Hooks

	shardOrchestratorClient shardorchestrator.ClientInterface
	shardRoutingSteady      shardRoutingSteadyObserver

	// myShardID is the shard index this syncer belongs to. Used to filter workflows.
	myShardID       uint32
	shardingEnabled bool
}

type shardRoutingSteadyObserver interface {
	ObserveRoutingSteady(steady bool)
	Invalidate()
}

type Hooks struct {
	OnStartFailure func(error)
}

type evtHandler interface {
	io.Closer
	Start(context.Context) error

	Handle(ctx context.Context, event Event) error
}

type donNotifier interface {
	WaitForDon(ctx context.Context) (capabilities.DON, error)
}

// WithTicker allows external callers to provide a ticker to the workflowRegistry.  This is useful
// for overriding the default tick interval.
func WithTicker(ticker <-chan time.Time) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
	}
}

func WithRetryInterval(retryInterval time.Duration) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.retryInterval = retryInterval
	}
}

func WithMaxConcurrency(maxConcurrency int) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		if maxConcurrency > 0 {
			wr.maxConcurrency = maxConcurrency
		}
	}
}

// AdditionalSourceConfig holds configuration for an additional workflow source.
type AdditionalSourceConfig struct {
	URL          string
	Name         string
	TLSEnabled   bool
	JWTGenerator nodeauthjwt.JWTGenerator
}

// WithAdditionalSources adds additional workflow sources to the registry.
// Sources are detected by URL scheme:
//   - file:// prefix -> FileWorkflowSource (reads from local JSON file)
//   - Otherwise -> GRPCWorkflowSource (connects to GRPC server)
//
// These sources supplement or replace the primary contract source.
func WithAdditionalSources(sources []AdditionalSourceConfig) Option {
	return func(wr *workflowRegistry) {
		successCount := 0
		failedSources := []string{}

		for _, src := range sources {
			// Detect source type by URL scheme
			if strings.HasPrefix(src.URL, "file://") {
				// File source - extract path from file:// URL
				filePath := strings.TrimPrefix(src.URL, "file://")
				fileSource, err := NewFileWorkflowSourceWithPath(wr.lggr, src.Name, filePath)
				if err != nil {
					wr.lggr.Errorw("Failed to create file workflow source",
						"name", src.Name,
						"path", filePath,
						"error", err)
					failedSources = append(failedSources, src.Name)
					continue
				}
				wr.workflowSources = append(wr.workflowSources, fileSource)
				successCount++
				wr.lggr.Infow("Added file workflow source",
					"name", src.Name,
					"path", filePath)
			} else {
				// GRPC source (default)
				grpcSource, err := NewGRPCWorkflowSource(wr.lggr, GRPCWorkflowSourceConfig{
					URL:          src.URL,
					TLSEnabled:   src.TLSEnabled,
					Name:         src.Name,
					JWTGenerator: src.JWTGenerator,
				})
				if err != nil {
					wr.lggr.Errorw("Failed to create GRPC workflow source",
						"name", src.Name,
						"url", src.URL,
						"error", err)
					failedSources = append(failedSources, src.Name)
					continue
				}
				wr.workflowSources = append(wr.workflowSources, grpcSource)
				successCount++
				wr.lggr.Infow("Added GRPC workflow source",
					"name", src.Name,
					"url", src.URL,
					"tls", src.TLSEnabled)
			}
		}

		// Log summary if any sources failed to initialize
		if len(failedSources) > 0 {
			wr.lggr.Warnw("Some additional sources failed to initialize",
				"expected", len(sources),
				"active", successCount,
				"failed", failedSources)
		}
	}
}

// Option is a functional option for configuring a workflowRegistry.
type Option func(*workflowRegistry)

func WithShardOrchestratorClient(client shardorchestrator.ClientInterface) Option {
	return func(wr *workflowRegistry) {
		wr.shardOrchestratorClient = client
	}
}

func WithShardEnabled(shardingEnabled bool) Option {
	return func(wr *workflowRegistry) {
		wr.shardingEnabled = shardingEnabled
	}
}

// WithShardID enables shard filtering and sets the shard ID for this syncer.
func WithShardID(shardID uint32) Option {
	return func(wr *workflowRegistry) {
		wr.myShardID = shardID
	}
}

func WithRegistryShardRoutingObserver(signal shardRoutingSteadyObserver) Option {
	return func(wr *workflowRegistry) {
		wr.shardRoutingSteady = signal
	}
}

// NewWorkflowRegistry returns a new v2 workflowRegistry.
// The addr parameter is optional - if empty, no contract source will be created,
// enabling pure GRPC-only or file-only workflow deployments.
// The chainSelector parameter identifies the chain where the workflow registry contract is deployed.
func NewWorkflowRegistry(
	lggr logger.Logger,
	contractReaderFn versioning.ContractReaderFactory,
	addr string,
	chainSelector string,
	config Config,
	handler evtHandler,
	workflowDonNotifier donNotifier,
	engineRegistry *EngineRegistry,
	opts ...Option,
) (*workflowRegistry, error) {
	if engineRegistry == nil {
		return nil, errors.New("engine registry must be provided")
	}

	m, err := newMetrics()
	if err != nil {
		return nil, err
	}

	var workflowSources []WorkflowMetadataSource

	// Only add contract source if address is configured
	if addr != "" {
		contractSource := NewContractWorkflowSource(lggr, contractReaderFn, addr, chainSelector)
		workflowSources = append(workflowSources, contractSource)
		lggr.Infow("Added contract workflow source",
			"contractAddress", addr,
			"chainSelector", chainSelector)
	} else {
		lggr.Infow("No contract address configured, skipping contract workflow source")
	}

	wr := &workflowRegistry{
		lggr:                             lggr,
		contractReaderFn:                 contractReaderFn,
		workflowRegistryAddress:          addr,
		chainSelector:                    chainSelector,
		lastSeenAllowlistedRequestsCount: big.NewInt(0),
		config:                           config,
		stopCh:                           make(services.StopChan),
		handler:                          handler,
		workflowDonNotifier:              workflowDonNotifier,
		metrics:                          m,
		engineRegistry:                   engineRegistry,
		retryInterval:                    defaultRetryInterval,
		maxRetryInterval:                 defaultMaxRetryInterval,
		maxConcurrency:                   defaultMaxConcurrency,
		clock:                            clockwork.NewRealClock(),
		hooks: Hooks{
			OnStartFailure: func(_ error) {},
		},
		workflowSources: workflowSources,
	}

	for _, opt := range opts {
		opt(wr)
	}

	lggr.Infow("Initialized workflow registry with multi-source support",
		"sourceCount", len(wr.workflowSources),
		"hasContractSource", addr != "")

	switch wr.config.SyncStrategy {
	case SyncStrategyReconciliation:
		break
	default:
		return nil, fmt.Errorf("WorkflowRegistry v2 contracts must use a SyncStrategy of: %s", SyncStrategyReconciliation)
	}

	return wr, nil
}

// Start begins the workflowRegistry service.
func (w *workflowRegistry) Start(_ context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		ctx, cancel := w.stopCh.NewCtx()
		initDoneCh := make(chan struct{})

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer w.lggr.Debugw("Successfully set ContractReader")
			defer close(initDoneCh)

			ticker := w.getTicker(defaultTickInterval)
			for w.contractReader == nil {
				select {
				case <-ctx.Done():
					w.lggr.Debug("shutting down workflowregistry, %s", ctx.Err())
					return
				case <-ticker:
					// Async initialization of contract reader for allowlisted requests.
					// There is an on-chain call dependency that would cause a deadlock if we block.
					// Instead, we poll until the contract reader is ready.
					reader, err := w.newAllowlistedRequestsContractReader(ctx)
					if err != nil {
						w.lggr.Infow("contract reader unavailable", "error", err.Error())
						break
					}
					w.contractReader = reader
				}
			}
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			// Start goroutines to gather changes from Workflow Registry contract
			select {
			case <-initDoneCh:
			case <-ctx.Done():
				return
			}
			w.lggr.Debugw("read from don received channel while waiting to start reconciliation sync")
			_, err := w.workflowDonNotifier.WaitForDon(ctx)
			if err != nil {
				w.hooks.OnStartFailure(fmt.Errorf("failed to start workflow sync strategy: %w", err))
				return
			}
			w.syncUsingReconciliationStrategy(ctx)
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			// Start goroutines to gather allowlisted requests from Workflow Registry contract
			select {
			case <-initDoneCh:
			case <-ctx.Done():
				return
			}
			w.syncAllowlistedRequests(ctx)
		}()

		return w.handler.Start(ctx)
	})
}

func (w *workflowRegistry) Close() error {
	return w.StopOnce(w.Name(), func() error {
		close(w.stopCh)
		w.wg.Wait()
		svcs := []io.Closer{w.handler}
		if w.shardOrchestratorClient != nil {
			svcs = append(svcs, w.shardOrchestratorClient)
		}
		return services.CloseAll(svcs...)
	})
}

func (w *workflowRegistry) Ready() error {
	return nil
}

func (w *workflowRegistry) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

func (w *workflowRegistry) Name() string {
	return name
}

func (w *workflowRegistry) handleWithMetrics(ctx context.Context, event Event) error {
	start := time.Now()
	err := w.handler.Handle(ctx, event)
	totalDuration := time.Since(start)
	w.metrics.recordHandleDuration(ctx, totalDuration, string(event.Name), err == nil)
	return err
}

// toLocalHead converts a chainlink-common Head to our local Head struct
func toLocalHead(head *types.Head) Head {
	return Head{
		Hash:      string(head.Hash),
		Height:    head.Height,
		Timestamp: head.Timestamp,
	}
}

// generateReconciliationEvents compares workflow metadata from a specific source against the engine registry's state.
// It only considers engines from the specified source when determining deletions. This ensures that when a source
// fails to fetch, we don't incorrectly delete engines from other sources.
func (w *workflowRegistry) generateReconciliationEvents(
	_ context.Context,
	pendingEvents map[string]*reconciliationEvent,
	workflowMetadata []WorkflowMetadataView,
	head *types.Head,
	sourceName string,
) ([]*reconciliationEvent, error) {
	var events []*reconciliationEvent
	localHead := toLocalHead(head)
	// workflowMetadataIDs is a set of workflow IDs present in this tick's metadata
	workflowMetadataIDs := make(map[string]struct{}, len(workflowMetadata))
	for _, wfMeta := range workflowMetadata {
		workflowMetadataIDs[wfMeta.WorkflowID.Hex()] = struct{}{}
	}

	// Keep track of which of the engines in the engineRegistry have been touched.
	workflowsSeen := make(map[string]bool, len(workflowMetadata))
	for _, wfMeta := range workflowMetadata {
		id := wfMeta.WorkflowID.Hex()
		engineFound := w.engineRegistry.Contains(wfMeta.WorkflowID)

		switch wfMeta.Status {
		case WorkflowStatusActive:
			switch engineFound {
			// we can't tell the difference between an activation and registration without holding
			// state in the db; so we handle as an activation event.
			case false:
				signature := fmt.Sprintf("%s-%s-%s", WorkflowActivated, id, toSpecStatus(wfMeta.Status))

				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
					events = append(events, pendingEvents[id])
					delete(pendingEvents, id)
					continue
				}

				delete(pendingEvents, id)

				toActivatedEvent := WorkflowActivatedEvent{
					WorkflowID:    wfMeta.WorkflowID,
					WorkflowOwner: wfMeta.Owner,
					CreatedAt:     wfMeta.CreatedAt,
					Status:        wfMeta.Status,
					WorkflowName:  wfMeta.WorkflowName,
					BinaryURL:     wfMeta.BinaryURL,
					ConfigURL:     wfMeta.ConfigURL,
					Tag:           wfMeta.Tag,
					Attributes:    wfMeta.Attributes,
					Source:        wfMeta.Source,
				}
				events = append(events, &reconciliationEvent{
					Event: Event{
						Data: toActivatedEvent,
						Name: WorkflowActivated,
						Head: localHead,
						Info: fmt.Sprintf("[ID: %s, Name: %s, Owner: %s, Source: %s]", wfMeta.WorkflowID.Hex(), wfMeta.WorkflowName, hex.EncodeToString(wfMeta.Owner), sourceName),
					},
					signature: signature,
					id:        id,
				})
				workflowsSeen[id] = true
			// if the workflow is active, the workflow engine is in the engine registry, and the metadata has not changed
			// then we don't need to action the event further. Mark as seen and continue.
			case true:
				workflowsSeen[id] = true
			}
		case WorkflowStatusPaused:
			signature := fmt.Sprintf("%s-%s-%s", WorkflowPaused, id, toSpecStatus(wfMeta.Status))
			switch engineFound {
			case false:
				// Account for a state change from active to paused, by checking
				// whether an existing pendingEvent exists.
				// We do this regardless of whether we have an event to handle or not, since this ensures
				// we correctly handle the state of pending events in the following situation:
				// - we registered an active workflow, but it failed to process successfully
				// - we then paused the workflow; this should clear the pending event
				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature != signature {
					delete(pendingEvents, id)
				}
			case true:
				// Will be handled in the event handler as a deleted event and will clear the DB workflow spec.
				workflowsSeen[id] = true

				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
					events = append(events, pendingEvents[id])
					delete(pendingEvents, id)
					continue
				}

				delete(pendingEvents, id)

				toPausedEvent := WorkflowPausedEvent{
					WorkflowID:    wfMeta.WorkflowID,
					WorkflowOwner: wfMeta.Owner,
					CreatedAt:     wfMeta.CreatedAt,
					Status:        wfMeta.Status,
					WorkflowName:  wfMeta.WorkflowName,
					Source:        wfMeta.Source,
				}
				events = append(
					[]*reconciliationEvent{
						{
							Event: Event{
								Data: toPausedEvent,
								Name: WorkflowPaused,
								Head: localHead,
								Info: fmt.Sprintf("[ID: %s, Name: %s, Owner: %s, Source: %s]", wfMeta.WorkflowID.Hex(), wfMeta.WorkflowName, hex.EncodeToString(wfMeta.Owner), sourceName),
							},
							signature: signature,
							id:        id,
						},
					},
					events...,
				)
			}
		default:
			return nil, fmt.Errorf("invariant violation: unable to determine difference from workflow metadata (status=%d)", wfMeta.Status)
		}
	}

	// Shut down engines that are no longer in the contract's latest workflow metadata state
	sourceEngines := w.engineRegistry.GetBySource(sourceName)
	for _, engine := range sourceEngines {
		id := engine.WorkflowID.Hex()
		if !workflowsSeen[id] {
			signature := fmt.Sprintf("%s-%s", WorkflowDeleted, id)

			if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
				events = append(events, pendingEvents[id])
				delete(pendingEvents, id)
				continue
			}

			delete(pendingEvents, id)

			toDeletedEvent := WorkflowDeletedEvent{
				WorkflowID: engine.WorkflowID,
				Source:     sourceName,
			}
			events = append(
				[]*reconciliationEvent{
					{
						Event: Event{
							Data: toDeletedEvent,
							Name: WorkflowDeleted,
							Head: localHead,
							Info: fmt.Sprintf("[ID: %s, Source: %s]", id, sourceName),
						},
						signature: signature,
						id:        id,
					},
				},
				events...,
			)
		}
	}

	// Clean up create events which no longer need to be attempted because
	// the workflow no longer exists in this source's metadata
	for id, event := range pendingEvents {
		if event.Name == WorkflowActivated {
			if _, ok := workflowMetadataIDs[event.Data.(WorkflowActivatedEvent).WorkflowID.Hex()]; !ok {
				delete(pendingEvents, id)
			}
		}
	}

	if len(pendingEvents) != 0 {
		return nil, fmt.Errorf("invariant violation: some pending events were not handled in the reconcile loop: keys=%+v, len=%d", maps.Keys(pendingEvents), len(pendingEvents))
	}

	return events, nil
}

func (w *workflowRegistry) syncAllowlistedRequests(ctx context.Context) {
	ticker := w.getTicker(defaultTickIntervalForAllowlistedRequests)
	w.lggr.Debug("starting syncAllowlistedRequests")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down syncAllowlistedRequests, %s", ctx.Err())
			return
		case <-ticker:
			newAllowListedRequests, totalAllowlistedRequests, head, err := w.getAllowlistedRequests(ctx, w.contractReader)
			if err != nil {
				w.lggr.Errorw("failed to call getAllowlistedRequests", "err", err)
				continue
			}
			w.allowListedMu.Lock()
			// Prune expired requests
			activeAllowlistedRequests := []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}
			expiredRequestsCount := 0
			for _, request := range w.allowListedRequests {
				if int64(request.ExpiryTimestamp) > time.Now().Unix() {
					activeAllowlistedRequests = append(activeAllowlistedRequests, request)
				} else {
					expiredRequestsCount++
				}
			}

			// Add new requests
			activeAllowlistedRequests = append(activeAllowlistedRequests, newAllowListedRequests...)
			w.allowListedRequests = activeAllowlistedRequests
			w.lastSeenAllowlistedRequestsCount = totalAllowlistedRequests
			w.lggr.Debugw("synced allowlisted requests",
				"newRequestsNum", len(newAllowListedRequests),
				"expiredRequestsNum", expiredRequestsCount,
				"activeRequestsNum", len(w.allowListedRequests),
				"lastSeenOnchainRequestsNum", w.lastSeenAllowlistedRequestsCount,
				"blockHeight", head.Height,
			)
			w.allowListedMu.Unlock()
		}
	}
}

func (w *workflowRegistry) filterWorkflowsByShard(ctx context.Context, workflows []WorkflowMetadataView) ([]WorkflowMetadataView, error) {
	if w.shardOrchestratorClient == nil {
		return workflows, nil
	}
	if len(workflows) == 0 {
		return workflows, nil
	}
	workflowIDs := make([]string, 0, len(workflows))
	for _, wf := range workflows {
		workflowIDs = append(workflowIDs, wf.WorkflowID.Hex())
	}
	resp, err := w.shardOrchestratorClient.GetWorkflowShardMapping(ctx, workflowIDs)
	if err != nil {
		if w.shardRoutingSteady != nil {
			w.shardRoutingSteady.Invalidate()
		}
		return nil, fmt.Errorf("shard mapping unavailable: %w", err)
	}
	if w.shardRoutingSteady != nil {
		w.shardRoutingSteady.ObserveRoutingSteady(resp.GetRoutingSteady())
	}
	filtered := make([]WorkflowMetadataView, 0, len(workflows))
	for _, wf := range workflows {
		id := wf.WorkflowID.Hex()
		if shardID, ok := resp.Mappings[id]; ok && shardID == w.myShardID {
			filtered = append(filtered, wf)
		}
	}
	return filtered, nil
}

// syncUsingReconciliationStrategy syncs workflow registry contract state by polling the workflow metadata state and comparing to local state.
// NOTE: In this mode paused states will be treated as a deleted workflow. Workflows will not be registered as paused.
// This function processes each source independently to ensure that failure in one source doesn't affect workflows from other sources.
func (w *workflowRegistry) syncUsingReconciliationStrategy(ctx context.Context) {
	ticker := w.getTicker(defaultTickInterval)
	pendingEventsBySource := make(map[string]map[string]*reconciliationEvent)
	w.lggr.Debug("running readRegistryStateLoop")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down readRegistryStateLoop")
			return
		case <-ticker:
			don, err := w.workflowDonNotifier.WaitForDon(ctx)
			if err != nil {
				w.lggr.Errorw("failed to get get don from notifier", "err", err)
				continue
			}
			w.lggr.Debugw("fetching workflow metadata from all sources", "don", don.Families)

			// Process each source independently to isolate failures
			totalWorkflowsFetched := 0
			reconcileReport := newReconcileReport()

			for _, source := range w.workflowSources {
				sourceName := source.Name()
				sourceIdentifier := source.SourceIdentifier()

				// Initialize pending events for this source if needed
				if pendingEventsBySource[sourceIdentifier] == nil {
					pendingEventsBySource[sourceIdentifier] = make(map[string]*reconciliationEvent)
				}
				pendingEvents := pendingEventsBySource[sourceIdentifier]

				// Fetch workflows from this source (each source handles lazy initialization internally)
				start := time.Now()
				workflows, head, fetchErr := source.ListWorkflowMetadata(ctx, don)
				duration := time.Since(start)

				// Record metrics for this source fetch
				w.metrics.recordSourceFetch(ctx, sourceName, len(workflows), duration, fetchErr)

				if fetchErr != nil {
					w.lggr.Errorw("Failed to fetch from source, skipping reconciliation for this source",
						"source", sourceName, "error", fetchErr, "durationMs", duration.Milliseconds())
					// KEY: Skip this source entirely - no events generated, no deletions
					continue
				}

				totalWorkflowsFetched += len(workflows)
				w.lggr.Debugw("Fetched workflows from source",
					"source", sourceName,
					"count", len(workflows),
					"durationMs", duration.Milliseconds())

				filteredWorkflowsMetadata := workflows
				if w.shardingEnabled {
					filteredWorkflowsMetadata, err = w.filterWorkflowsByShard(ctx, workflows)
					if err != nil {
						w.lggr.Errorw("failed to filter workflows by shard",
							"err", err,
							"source", sourceName)
						continue
					}
					w.lggr.Debugw("filtered workflows by shard",
						"total", len(workflows),
						"filtered", len(filteredWorkflowsMetadata),
						"shardID", w.myShardID,
						"source", sourceName,
					)
				}

				// Generate events only for this source's engines (using sourceIdentifier for engine registry lookups)
				events, genErr := w.generateReconciliationEvents(ctx, pendingEvents, filteredWorkflowsMetadata, head, sourceIdentifier)
				if genErr != nil {
					w.lggr.Errorw("Failed to generate reconciliation events for source",
						"source", sourceName, "error", genErr)
					continue
				}

				w.lggr.Debugw("Generated events for source", "source", sourceName, "num", len(events))

				// Clear pending events after successful reconciliation
				pendingEventsBySource[sourceIdentifier] = make(map[string]*reconciliationEvent)

				// Handle events concurrently — each event targets a distinct workflow ID
				var wg sync.WaitGroup
				var mu sync.Mutex
				sem := make(chan struct{}, w.maxConcurrency)
				batchStart := time.Now()
				var dispatched, backoffCount int
				for _, event := range events {
					select {
					case <-ctx.Done():
						w.lggr.Debug("readRegistryStateLoop stopped during processing")
						return
					default:
					}

					w.lggr.Debugw("processing event", "source", sourceName, "event", event.Name, "id", event.id, "signature", event.signature, "workflowInfo", event.Info)

					mu.Lock()
					reconcileReport.NumEventsByType[string(event.Name)]++
					mu.Unlock()

					if event.retryCount > 0 && !w.clock.Now().After(event.nextRetryAt) {
						backoffCount++
						mu.Lock()
						pendingEventsBySource[sourceIdentifier][event.id] = event
						reconcileReport.Backoffs[event.id] = event.nextRetryAt
						mu.Unlock()
						w.lggr.Debugw("skipping event, still in backoff", "nextRetryAt", event.nextRetryAt, "event", event.Name, "id", event.id, "signature", event.signature, "workflowInfo", event.Info)
						continue
					}

					select {
					case sem <- struct{}{}:
					case <-ctx.Done():
						w.lggr.Debug("readRegistryStateLoop stopped waiting for semaphore")
						return
					}

					dispatched++
					wg.Add(1)
					go func(evt *reconciliationEvent) {
						defer func() {
							<-sem
							wg.Done()
						}()
						handleErr := w.handleWithMetrics(ctx, evt.Event)
						if handleErr != nil {
							evt.updateNextRetryFor(w.clock, w.retryInterval, w.maxRetryInterval)
							mu.Lock()
							pendingEventsBySource[sourceIdentifier][evt.id] = evt
							reconcileReport.Backoffs[evt.id] = evt.nextRetryAt
							mu.Unlock()
							w.lggr.Errorw("failed to handle event, backing off...", "err", handleErr, "type", evt.Name, "nextRetryAt", evt.nextRetryAt, "retryCount", evt.retryCount, "workflowInfo", evt.Info)
						}
					}(event)
				}
				wg.Wait()

				// prompt the GC to reclaim transient allocations from event handling
				// that would otherwise be delayed because the dominant CGo/wasmtime memory is invisible to the Go GC
				if dispatched > 0 {
					runtime.GC()
				}

				batchDuration := time.Since(batchStart)
				w.metrics.recordReconcileBatch(ctx, sourceName, dispatched, batchDuration)
				if backoffCount > 0 {
					w.metrics.recordReconcileBackoff(ctx, sourceName, backoffCount)
				}

				w.lggr.Infow("reconciliation tick completed",
					"source", sourceName,
					"dispatched", dispatched,
					"backoffs", backoffCount,
					"failed", len(pendingEventsBySource[sourceIdentifier]),
					"durationMs", batchDuration.Milliseconds(),
					"eventsByType", reconcileReport.NumEventsByType,
				)
			}

			w.metrics.recordFetchedWorkflows(ctx, totalWorkflowsFetched)
			w.lggr.Debugw("reconciled events", "report", reconcileReport)

			runningWorkflows := w.engineRegistry.GetAll()
			w.metrics.recordRunningWorkflows(ctx, len(runningWorkflows))
			w.metrics.incrementCompletedSyncs(ctx)
		}
	}
}

// getTicker returns the ticker that the workflowRegistry will use to poll for events.  If the ticker
// is nil, then a default ticker is returned.
func (w *workflowRegistry) getTicker(d time.Duration) <-chan time.Time {
	if w.ticker == nil {
		return time.NewTicker(d).C
	}

	return w.ticker
}

// isEmptyWorkflowID checks if a WorkflowID is empty (all zeros)
func isEmptyWorkflowID(wfID [32]byte) bool {
	emptyID := [32]byte{}
	return wfID == emptyID
}

// isZeroOwner checks if a workflow owner address is the zero address (all zeros).
// This can indicate stale metadata from deleted workflows in the contract - there's a known
// bug where deleted workflows aren't always fully removed from the contract state.
func isZeroOwner(owner []byte) bool {
	// does not contain non-zero bytes
	return !slices.ContainsFunc(owner, func(b byte) bool { return b != 0 })
}

// newAllowlistedRequestsContractReader creates a contract reader specifically for fetching
// allowlisted requests from the WorkflowRegistry contract. This is used by Vault DON nodes
// to verify that incoming vault requests have been pre-authorized on-chain by workflow owners.
//
// Note: Workflow metadata is fetched separately via ContractWorkflowSource, which maintains
// its own contract reader. The two concerns are separated because:
//   - Allowlisted requests: Used by Vault DON for request authorization
//   - Workflow metadata: Used by workflow engine for deployment/reconciliation
func (w *workflowRegistry) newAllowlistedRequestsContractReader(
	ctx context.Context,
) (types.ContractReader, error) {
	contractReaderCfg := config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper_v2.WorkflowRegistryABI,
				Configs: map[string]*config.ChainReaderDefinition{
					GetActiveAllowlistedRequestsReverseMethodName: {
						ChainSpecificName: GetActiveAllowlistedRequestsReverseMethodName,
						ReadType:          config.Method,
					},
					TotalAllowlistedRequestsMethodName: {
						ChainSpecificName: TotalAllowlistedRequestsMethodName,
						ReadType:          config.Method,
					},
				},
			},
		},
	}

	marshalledCfg, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, err
	}

	reader, err := w.contractReaderFn(ctx, marshalledCfg)
	if err != nil {
		return nil, err
	}

	bc := types.BoundContract{
		Name:    WorkflowRegistryContractName,
		Address: w.workflowRegistryAddress,
	}

	// bind contract to contract reader
	if err := reader.Bind(ctx, []types.BoundContract{bc}); err != nil {
		return nil, err
	}

	if err := reader.Start(ctx); err != nil {
		return nil, err
	}

	return reader, nil
}

func (w *workflowRegistry) GetAllowlistedRequests(_ context.Context) []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	w.allowListedMu.RLock()
	defer w.allowListedMu.RUnlock()
	allowListedRequests := make([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, len(w.allowListedRequests))
	copy(allowListedRequests, w.allowListedRequests)
	return allowListedRequests
}

func (w *workflowRegistry) GetLastSeenOnchainAllowlistedRequestsCount(_ context.Context) *big.Int {
	w.allowListedMu.RLock()
	defer w.allowListedMu.RUnlock()
	if w.lastSeenAllowlistedRequestsCount == nil {
		return nil
	}
	return new(big.Int).Set(w.lastSeenAllowlistedRequestsCount)
}

// GetAllowlistedRequests uses contract reader to query the contract for all allowlisted requests
func (w *workflowRegistry) getAllowlistedRequests(ctx context.Context, contractReader types.ContractReader) ([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, *big.Int, *types.Head, error) {
	if contractReader == nil {
		return nil, nil, nil, errors.New("cannot fetch allow listed requests: nil contract reader")
	}
	contractBinding := types.BoundContract{
		Address: w.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	// Read current total allowlisted requests
	var headAtLastRead *types.Head
	var totalAllowlistedRequestsResult *big.Int
	readIdentifier := contractBinding.ReadIdentifier(TotalAllowlistedRequestsMethodName)
	headAtLastRead, err := contractReader.GetLatestValueWithHeadData(
		ctx, readIdentifier, primitives.Unconfirmed, nil, &totalAllowlistedRequestsResult,
	)
	if err != nil {
		return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, w.lastSeenAllowlistedRequestsCount, &types.Head{Height: "0"}, errors.New("failed to get latest value with head data. error: " + err.Error())
	}

	if w.lastSeenAllowlistedRequestsCount.Cmp(totalAllowlistedRequestsResult) == 0 {
		return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, totalAllowlistedRequestsResult, headAtLastRead, nil
	}

	var newAllowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	readIdentifier = contractBinding.ReadIdentifier(GetActiveAllowlistedRequestsReverseMethodName)
	var endIndex = new(big.Int).Sub(totalAllowlistedRequestsResult, big.NewInt(1))
	var startIndex *big.Int

	for {
		var err error
		var response struct {
			AllowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
			SearchComplete      bool
			err                 error
		}

		// Start index should be no more than MaxResultsPerQuery away from end index
		startIndex = new(big.Int).Sub(endIndex, big.NewInt(MaxResultsPerQuery-1))
		// If start index is less than last seen allowlisted requests count, set it to last seen allowlisted requests
		// count to avoid duplicate requests
		if startIndex.Cmp(w.lastSeenAllowlistedRequestsCount) < 0 {
			startIndex = w.lastSeenAllowlistedRequestsCount
		}

		params := GetActiveAllowlistedRequestsReverseParams{
			EndIndex:   endIndex,
			StartIndex: startIndex,
		}
		w.lggr.Debugw("getting active allowlisted requests",
			"endIndex", endIndex,
			"startIndex", startIndex,
		)
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(
			ctx, readIdentifier, primitives.Unconfirmed, params, &response,
		)
		if err != nil {
			return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, w.lastSeenAllowlistedRequestsCount, &types.Head{Height: "0"}, errors.New("failed to get lastest value with head data. error: " + err.Error())
		}

		w.lggr.Debugw("contract call response",
			"fetchedAllowlistedRequestsNum", len(response.AllowlistedRequests),
			"searchComplete", response.SearchComplete,
			"error", response.err,
			"blockHeight", headAtLastRead.Height)

		for _, request := range response.AllowlistedRequests {
			newAllowlistedRequests = append(newAllowlistedRequests, workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
				RequestDigest:   request.RequestDigest,
				Owner:           request.Owner,
				ExpiryTimestamp: request.ExpiryTimestamp,
			})
		}

		// We can break early if the search is complete even if we haven't
		// looked at all the allowlisted requests. This is because the contract
		// method determines if there are more allowlisted requests to fetch.
		if response.SearchComplete {
			break
		}

		// If search is not complete, set the end index to the start index minus MaxResultsPerQuery
		// to continue fetching the next batch of allowlisted requests
		endIndex = endIndex.Sub(endIndex, big.NewInt(MaxResultsPerQuery))
		// Ensure endIndex doesn't go below zero
		if endIndex.Cmp(big.NewInt(0)) < 0 {
			endIndex = big.NewInt(0)
		}
	}

	return newAllowlistedRequests, totalAllowlistedRequestsResult, headAtLastRead, nil
}
