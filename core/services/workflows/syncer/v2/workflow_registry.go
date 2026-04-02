package v2

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodeauthjwt "github.com/smartcontractkit/chainlink-common/pkg/nodeauth/jwt"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

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

	GetWorkflowsByDONMethodName = "getWorkflowListByDON"

	// MaxResultsPerQuery defines the maximum number of results that can be queried in a single request.
	// The default value of 1,000 was chosen based on expected system performance and typical use cases.
	MaxResultsPerQuery = int64(1_000)
)

// WorkflowRegistrySyncer is the public interface of the package.
type WorkflowRegistrySyncer interface {
	services.Service
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

	contractReaderFn versioning.ContractReaderFactory

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

	// myShardID is the shard index this syncer belongs to. Used to filter workflows.
	myShardID       uint32
	shardingEnabled bool
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

// NewWorkflowRegistry returns a new v2 workflowRegistry.
// The addr parameter is optional - if empty, no contract source will be created,
// enabling pure GRPC-only or file-only workflow deployments.
// The chainSelector parameter identifies the chain where the workflow registry contract is deployed.
func NewWorkflowRegistry(
	lggr logger.Logger,
	contractReaderFn versioning.ContractReaderFactory,
	workflowSources []WorkflowMetadataSource,
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

	wr := &workflowRegistry{
		lggr:                    lggr,
		contractReaderFn:        contractReaderFn,
		workflowRegistryAddress: addr,
		config:                  config,
		stopCh:                  make(services.StopChan),
		handler:                 handler,
		workflowDonNotifier:     workflowDonNotifier,
		metrics:                 m,
		engineRegistry:          engineRegistry,
		retryInterval:           defaultRetryInterval,
		maxRetryInterval:        defaultMaxRetryInterval,
		maxConcurrency:          defaultMaxConcurrency,
		clock:                   clockwork.NewRealClock(),
		hooks: Hooks{
			OnStartFailure: func(_ error) {},
		},
		workflowSources: workflowSources,
	}

	for _, opt := range opts {
		opt(wr)
	}

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

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			_, err := w.workflowDonNotifier.WaitForDon(ctx)
			if err != nil {
				w.hooks.OnStartFailure(fmt.Errorf("failed to start workflow sync strategy: %w", err))
				return
			}
			w.syncUsingReconciliationStrategy(ctx)
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
	// workflowMetadataMap is only used for lookups; disregard when reading the state machine.
	workflowMetadataMap := make(map[string]WorkflowMetadataView)
	for _, wfMeta := range workflowMetadata {
		workflowMetadataMap[wfMeta.WorkflowID.Hex()] = wfMeta
	}

	// Keep track of which of the engines in the engineRegistry have been touched
	workflowsSeen := map[string]bool{}
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
			if _, ok := workflowMetadataMap[event.Data.(WorkflowActivatedEvent).WorkflowID.Hex()]; !ok {
				delete(pendingEvents, id)
			}
		}
	}

	if len(pendingEvents) != 0 {
		return nil, fmt.Errorf("invariant violation: some pending events were not handled in the reconcile loop: keys=%+v, len=%d", maps.Keys(pendingEvents), len(pendingEvents))
	}

	return events, nil
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
		return nil, fmt.Errorf("shard mapping unavailable: %w", err)
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
