package syncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"math"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	wftypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval                    = 12 * time.Second
	defaultRetryInterval                   = 12 * time.Second
	defaultMaxRetryInterval                = 5 * time.Minute
	WorkflowRegistryContractName           = "WorkflowRegistry"
	GetWorkflowMetadataListByDONMethodName = "getWorkflowMetadataListByDON"
)

type Head struct {
	Hash      string
	Height    string
	Timestamp uint64
}

type GetWorkflowMetadataListByDONParams struct {
	DonID uint32
	Start uint64
	Limit uint64
}

type GetWorkflowMetadata struct {
	WorkflowID   wftypes.WorkflowID
	Owner        []byte
	DonID        uint32
	Status       uint8
	WorkflowName string
	BinaryURL    string
	ConfigURL    string
	SecretsURL   string
}

type SyncStrategy string

const (
	SyncStrategyEvent          = "event"
	SyncStrategyReconciliation = "reconciliation"
	defaultSyncStrategy        = SyncStrategyEvent
)

const (
	WorkflowStatusActive uint8 = iota
	WorkflowStatusPaused
)

type GetWorkflowMetadataListByDONReturnVal struct {
	WorkflowMetadataList []GetWorkflowMetadata
}

// workflowRegistryEvent is the event emitted by the WorkflowRegistry in events mode.
// Each event is typed so that the consumer can determine how to handle the event.
type workflowRegistryEvent struct {
	Event
	Cursor string
	Head   Head
	DonID  *uint32
}

type Config struct {
	QueryCount   uint64
	SyncStrategy SyncStrategy
}

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error)

// ContractReader is a subset of types.ContractReader defined locally to enable mocking.
type ContractReader interface {
	Start(ctx context.Context) error
	Close() error
	Bind(context.Context, []types.BoundContract) error
	QueryKeys(ctx context.Context, keyQueries []types.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, types.Sequence], error)
	GetLatestValueWithHeadData(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *types.Head, err error)
}

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

// WorkflowRegistrySyncer is the public interface of the package.
type WorkflowRegistrySyncer interface {
	services.Service
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

// workflowRegistry is the implementation of the WorkflowRegistrySyncer interface.
type workflowRegistry struct {
	services.StateMachine

	// close stopCh to stop the workflowRegistry.
	stopCh services.StopChan

	// events sent to the event channel to be handled.
	eventCh chan Event

	// all goroutines are waited on with wg.
	wg sync.WaitGroup

	// ticker is the interval at which the workflowRegistry will
	// poll the contract for events, and poll the contract for the latest workflow metadata.
	ticker <-chan time.Time

	lggr                    logger.Logger
	workflowRegistryAddress string

	newContractReaderFn newContractReaderFn

	config Config

	handler evtHandler

	workflowDonNotifier donNotifier

	metrics *metrics

	engineRegistry *EngineRegistry

	retryInterval    time.Duration
	maxRetryInterval time.Duration
	clock            clockwork.Clock
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

type evtHandler interface {
	io.Closer
	Handle(ctx context.Context, event Event) error
}

type donNotifier interface {
	WaitForDon(ctx context.Context) (capabilities.DON, error)
}

type newContractReaderFn func(context.Context, []byte) (ContractReader, error)

// NewWorkflowRegistry returns a new workflowRegistry.
// Only queries for WorkflowRegistryForceUpdateSecretsRequestedV1 events.
func NewWorkflowRegistry(
	lggr logger.Logger,
	newContractReaderFn newContractReaderFn,
	addr string,
	config Config,
	handler evtHandler,
	workflowDonNotifier donNotifier,
	engineRegistry *EngineRegistry,
	opts ...func(*workflowRegistry),
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
		newContractReaderFn:     newContractReaderFn,
		workflowRegistryAddress: addr,
		config:                  config,
		eventCh:                 make(chan Event),
		stopCh:                  make(services.StopChan),
		handler:                 handler,
		workflowDonNotifier:     workflowDonNotifier,
		metrics:                 m,
		engineRegistry:          engineRegistry,
		retryInterval:           defaultRetryInterval,
		maxRetryInterval:        defaultMaxRetryInterval,
		clock:                   clockwork.NewRealClock(),
	}

	for _, opt := range opts {
		opt(wr)
	}

	switch wr.config.SyncStrategy {
	case SyncStrategyEvent:
	case SyncStrategyReconciliation:
		break
	default:
		return nil, fmt.Errorf("SyncStrategy must be one of: %s, %s", SyncStrategyEvent, SyncStrategyReconciliation)
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

			w.lggr.Debugw("Waiting for DON...")
			don, err := w.workflowDonNotifier.WaitForDon(ctx)
			if err != nil {
				w.lggr.Errorw("failed to wait for don", "err", err)
				return
			}

			reader, err := w.newWorkflowRegistryContractReader(ctx)
			if err != nil {
				w.lggr.Criticalf("contract reader unavailable : %s", err)
				return
			}

			// Start goroutines to gather changes from Workflow Registry contract
			switch w.config.SyncStrategy {
			case SyncStrategyEvent:
				w.syncUsingEventStrategy(ctx, don, reader)
			case SyncStrategyReconciliation:
				w.syncUsingReconciliationStrategy(ctx, don, reader)
			}

		}()

		return nil
	})
}

func (w *workflowRegistry) Close() error {
	return w.StopOnce(w.Name(), func() error {
		close(w.stopCh)
		w.wg.Wait()
		return w.handler.Close()
	})
}

func (w *workflowRegistry) Ready() error {
	return nil
}

func (w *workflowRegistry) HealthReport() map[string]error {
	return nil
}

func (w *workflowRegistry) Name() string {
	return name
}

// readRegistryEventsLoop polls the contract for events and sends them to the events channel for handling.
func (w *workflowRegistry) readRegistryEventsLoop(ctx context.Context, eventTypes []WorkflowRegistryEventType, don capabilities.DON, reader ContractReader, lastReadBlockNumber string) {
	ticker := w.getTicker()

	var keyQueries = make([]types.ContractKeyFilter, 0, len(eventTypes))
	for _, et := range eventTypes {
		var logData values.Value
		keyQueries = append(keyQueries, types.ContractKeyFilter{
			KeyFilter: query.KeyFilter{
				Key: string(et),
				Expressions: []query.Expression{
					query.Confidence(primitives.Finalized),
					query.Block(lastReadBlockNumber, primitives.Gt),
				},
			},
			Contract: types.BoundContract{
				Name:    WorkflowRegistryContractName,
				Address: w.workflowRegistryAddress,
			},
			SequenceDataType: &logData,
		})
	}

	cursor := ""
	w.lggr.Debug("running readRegistryEventsLoop")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down readRegistryEventsLoop")
			return
		case <-ticker:
			limitAndSort := query.LimitAndSort{
				SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
				Limit:  query.Limit{Count: w.config.QueryCount},
			}
			if cursor != "" {
				limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, w.config.QueryCount)
			}

			logsIter, err := reader.QueryKeys(ctx, keyQueries, limitAndSort)
			if err != nil {
				w.lggr.Errorw("failed to query keys", "err", err)
				continue
			}

			var logs []sequenceWithEventType
			for eventType, log := range logsIter {
				logs = append(logs, sequenceWithEventType{
					Sequence:  log,
					EventType: WorkflowRegistryEventType(eventType),
				})
			}
			w.lggr.Debugw("QueryKeys called", "logs", len(logs), "eventTypes", eventTypes, "lastReadBlockNumber", lastReadBlockNumber, "logCursor", cursor)

			// ChainReader QueryKey API provides logs including the cursor value and not
			// after the cursor value. If the response only consists of the log corresponding
			// to the cursor and no log after it, then we understand that there are no new
			// logs
			if len(logs) == 1 && logs[0].Sequence.Cursor == cursor {
				w.lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}

			var events []workflowRegistryEvent
			for _, log := range logs {
				if log.Sequence.Cursor == cursor {
					continue
				}

				event, err := toWorkflowRegistryEventResponse(log.Sequence, log.EventType, w.lggr)
				if err != nil {
					w.lggr.Errorw("failed to convert log to workflow registry event, skipping...", "err", err)
					continue
				}

				switch {
				case event.DonID == nil:
					// event is missing a DonID, so don't filter it out;
					// it applies to all Dons
					events = append(events, event)
				case *event.DonID == don.ID:
					// event has a DonID and matches, so it applies to this DON.
					events = append(events, event)
				default:
					// event doesn't match, let's skip it
					donID := "MISSING_DON_ID"
					if event.DonID != nil {
						donID = strconv.FormatUint(uint64(*event.DonID), 10)
					}
					w.lggr.Debugw("event belongs to a different don, skipping...", "don", don.ID, "gotDON", donID)
				}

				cursor = log.Sequence.Cursor
			}

			for _, event := range events {
				select {
				case <-ctx.Done():
					w.lggr.Debug("readRegistryEventsLoop stopped during processing")
					return
				default:
					err := w.handleWithMetrics(ctx, event.Event)
					if err != nil {
						w.lggr.Errorw("failed to handle event", "err", err, "type", event.EventType)
					}
				}
			}
		}
	}
}

func (w *workflowRegistry) handleWithMetrics(ctx context.Context, event Event) error {
	start := time.Now()
	err := w.handler.Handle(ctx, event)
	totalDuration := time.Since(start)
	w.metrics.recordHandleDuration(ctx, totalDuration, string(event.EventType), err == nil)
	return err
}

// syncUsingEventStrategy syncs workflow registry contract state by watching for events on the contract.
// It first loads the workflow metadata from the contract state as WorkflowRegistered events,
// and then starts a goroutine with one loop for handling the events.
func (w *workflowRegistry) syncUsingEventStrategy(ctx context.Context, don capabilities.DON, reader ContractReader) {
	w.lggr.Debugw("Loading initial workflows for DON", "DON", don.ID)

	workflowMetadata, loadWorkflowsHead, err := w.getWorkflowMetadata(ctx, don, reader)
	if err != nil {
		w.lggr.Errorw("failed to load initial workflows", "err", err)
	}

	w.lggr.Debugw("Rehydrating existing workflows", "len", len(workflowMetadata))
	for _, workflow := range workflowMetadata {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shut down during initial workflow registration")
			return

		default:
			err := w.handleWithMetrics(ctx, Event{
				Data: WorkflowRegisteredV1{
					WorkflowID:    workflow.WorkflowID,
					WorkflowOwner: workflow.Owner,
					DonID:         workflow.DonID,
					Status:        workflow.Status,
					WorkflowName:  workflow.WorkflowName,
					BinaryURL:     workflow.BinaryURL,
					ConfigURL:     workflow.ConfigURL,
					SecretsURL:    workflow.SecretsURL,
				},
				EventType: WorkflowRegisteredEvent,
			})
			if err != nil {
				w.lggr.Errorw("failed to handle event", "err", err)
			}

		}
	}

	// Poll for all workflow related events
	ets := []WorkflowRegistryEventType{
		ForceUpdateSecretsEvent,
		WorkflowActivatedEvent,
		WorkflowDeletedEvent,
		WorkflowPausedEvent,
		WorkflowRegisteredEvent,
		WorkflowUpdatedEvent,
	}

	w.readRegistryEventsLoop(ctx, ets, don, reader, loadWorkflowsHead.Height)
}

type reconciliationEvent struct {
	Event
	id          string
	signature   string
	nextRetryAt time.Time
	retryCount  int
}

func (r *reconciliationEvent) updateNextRetryFor(clock clockwork.Clock, retryInterval time.Duration, maxRetryInterval time.Duration) {
	r.retryCount++
	nextRetry := math.Pow(2, float64(r.retryCount)) * float64(retryInterval)
	nextRetry = math.Min(float64(maxRetryInterval), nextRetry)
	r.nextRetryAt = clock.Now().Add(time.Duration(nextRetry))
}

func idFor(owner []byte, name string) string {
	return hex.EncodeToString(owner) + "-" + name
}

// generateReconciliationEvents compares the workflow registry workflow metadata state against the engine registry's state.
// Differences are handled by the event handler by creating events that are sent to the events channel for handling.
func (w *workflowRegistry) generateReconciliationEvents(ctx context.Context, pendingEvents map[string]*reconciliationEvent, workflowMetadata []GetWorkflowMetadata, donID uint32) ([]*reconciliationEvent, error) {
	var events []*reconciliationEvent

	// Keep track of which of the engines in the engineRegistry have been touched
	workflowsSeen := map[string]bool{}
	for _, wfMeta := range workflowMetadata {
		engine, engineFound := w.engineRegistry.Get(EngineRegistryKey{Owner: wfMeta.Owner, Name: wfMeta.WorkflowName})

		currWfID := wfMeta.WorkflowID.Hex()
		prevWfID := engine.WorkflowID.Hex()

		id := idFor(wfMeta.Owner, wfMeta.WorkflowName)
		switch {
		case wfMeta.Status == WorkflowStatusActive:
			switch {
			// if the workflow is active, but unable to get engine from the engine registry
			// then handle as registered event
			case !engineFound:
				signature := fmt.Sprintf("%s-%s-%s", WorkflowRegisteredEvent, currWfID, toSpecStatus(wfMeta.Status))

				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
					events = append(events, pendingEvents[id])
					delete(pendingEvents, id)
					continue
				}

				delete(pendingEvents, id)

				toRegisteredEvent := WorkflowRegisteredV1{
					WorkflowID:    wfMeta.WorkflowID,
					WorkflowOwner: wfMeta.Owner,
					DonID:         wfMeta.DonID,
					Status:        wfMeta.Status,
					WorkflowName:  wfMeta.WorkflowName,
					BinaryURL:     wfMeta.BinaryURL,
					ConfigURL:     wfMeta.ConfigURL,
					SecretsURL:    wfMeta.SecretsURL,
				}
				events = append(events, &reconciliationEvent{
					Event: Event{
						Data:      toRegisteredEvent,
						EventType: WorkflowRegisteredEvent,
					},
					signature: signature,
					id:        id,
				})
				workflowsSeen[id] = true
			// if the workflow is active, the workflow engine is in the engine registry, but the metadata has changed
			// then handle as updated event
			case currWfID != prevWfID:
				signature := fmt.Sprintf("%s-%s-%s-%s", WorkflowUpdatedEvent, engine.WorkflowID.Hex(), wfMeta.WorkflowID.Hex(), toSpecStatus(wfMeta.Status))

				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
					events = append(events, pendingEvents[id])
					delete(pendingEvents, id)
					continue
				}

				delete(pendingEvents, id)

				toUpdatedEvent := WorkflowUpdatedV1{
					OldWorkflowID: engine.WorkflowID,
					NewWorkflowID: wfMeta.WorkflowID,
					WorkflowOwner: wfMeta.Owner,
					DonID:         wfMeta.DonID,
					WorkflowName:  wfMeta.WorkflowName,
					BinaryURL:     wfMeta.BinaryURL,
					ConfigURL:     wfMeta.ConfigURL,
					SecretsURL:    wfMeta.SecretsURL,
					Status:        wfMeta.Status,
				}
				events = append(events, &reconciliationEvent{
					Event: Event{
						Data:      toUpdatedEvent,
						EventType: WorkflowUpdatedEvent,
					},
					signature: signature,
					id:        id,
				})
				workflowsSeen[id] = true
			// if the workflow is active, the workflow engine is in the engine registry, and the metadata has not changed
			// then we don't need to action the event further. Mark as seen and continue.
			case currWfID == prevWfID:
				workflowsSeen[id] = true
			default:
				return nil, fmt.Errorf("invariant violation: could not handle workflow (currWfID=%s; prevWfID=%s, engineFound=%t) in active status", currWfID, prevWfID, engineFound)
			}
		case wfMeta.Status == WorkflowStatusPaused:
			switch {
			case !engineFound:
				// Account for a state change from active to paused, by checking
				// whether an existing pendingEvent exists.
				// We do this regardless of whether we have an event to handle or not, since this ensures
				// we correctly handle the state of pending events in the following situation:
				// - we registered an active workflow but it failed to process successfully
				// - we then paused the workflow; this should clear the pending event
				signature := fmt.Sprintf("%s-%s-%s", WorkflowPausedEvent, currWfID, toSpecStatus(wfMeta.Status))
				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature != signature {
					delete(pendingEvents, id)
				}
			default:
				// Paused means we skip for processing as a deleted event
				// To be handled below as a deleted event, which clears the DB workflow spec.
			}
		default:
			return nil, fmt.Errorf("invariant violation: unable to determine difference from workflow metadata (status=%d)", wfMeta.Status)
		}
	}

	// Shut down engines that are no longer in the contract's latest workflow metadata state
	allEngines := w.engineRegistry.GetAll()
	for _, engine := range allEngines {
		id := idFor(engine.WorkflowOwner, engine.WorkflowName)
		if !workflowsSeen[id] {
			signature := fmt.Sprintf("%s-%s", WorkflowDeletedEvent, engine.WorkflowID.Hex())

			if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
				events = append(events, pendingEvents[id])
				delete(pendingEvents, id)
				continue
			}

			delete(pendingEvents, id)

			toDeletedEvent := WorkflowDeletedV1{
				WorkflowID:    engine.WorkflowID,
				WorkflowOwner: engine.WorkflowOwner,
				DonID:         donID,
				WorkflowName:  engine.WorkflowName,
			}
			events = append(events, &reconciliationEvent{
				Event: Event{
					Data:      toDeletedEvent,
					EventType: WorkflowDeletedEvent,
				},
				signature: signature,
				id:        id,
			})
		}
	}

	if len(pendingEvents) != 0 {
		return nil, fmt.Errorf("invariant violation: some pending events were not handled in the reconcile loop: keys=%+v, len=%d", maps.Keys(pendingEvents), len(pendingEvents))
	}

	return events, nil
}

type reconcileReport struct {
	// events is a map of event type to the number of events of that type
	NumEventsByType map[string]int
	// id -> nextRetry time
	Backoffs map[string]time.Time
}

func newReconcileReport() *reconcileReport {
	return &reconcileReport{
		NumEventsByType: map[string]int{},
		Backoffs:        map[string]time.Time{},
	}
}

// syncUsingReconciliationStrategy syncs workflow registry contract state by polling the workflow metadata state and comparing to local state.
// It still watches for ForceUpdateSecretsEvents, which can't be reconciled through workflow metadata state.
// NOTE: In this mode paused states will be treated as a deleted workflow. Workflows will not be registered as paused.
func (w *workflowRegistry) syncUsingReconciliationStrategy(ctx context.Context, don capabilities.DON, reader ContractReader) {
	_, loadWorkflowsHead, err := w.getWorkflowMetadata(ctx, don, reader)
	if err != nil {
		w.lggr.Errorw("failed to load workflows head", "err", err)
	}

	// Poll for events for only ForceUpdateSecretsEvent
	ets := []WorkflowRegistryEventType{
		ForceUpdateSecretsEvent,
	}
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.readRegistryEventsLoop(ctx, ets, don, reader, loadWorkflowsHead.Height)
	}()

	ticker := w.getTicker()
	pendingEvents := map[string]*reconciliationEvent{}
	w.lggr.Debug("running readRegistryStateLoop")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down readRegistryStateLoop")
			return
		case <-ticker:
			workflowMetadata, _, err := w.getWorkflowMetadata(ctx, don, reader)
			if err != nil {
				w.lggr.Errorw("failed to get registry state", "err", err)
				continue
			}
			events, err := w.generateReconciliationEvents(ctx, pendingEvents, workflowMetadata, don.ID)
			if err != nil {
				w.lggr.Errorw("failed to generate reconciliation events", "err", err)
				continue
			}

			pendingEvents = map[string]*reconciliationEvent{}

			// Send events generated from differences to the event channel to be handled
			reconcileReport := newReconcileReport()
			for _, event := range events {
				select {
				case <-ctx.Done():
					w.lggr.Debug("readRegistryStateLoop stopped during processing")
					return
				default:
					reconcileReport.NumEventsByType[string(event.EventType)]++

					if event.retryCount == 0 || w.clock.Now().After(event.nextRetryAt) {
						err := w.handleWithMetrics(ctx, event.Event)
						if err != nil {
							event.updateNextRetryFor(w.clock, w.retryInterval, w.maxRetryInterval)

							pendingEvents[event.id] = event

							reconcileReport.Backoffs[event.id] = event.nextRetryAt
							w.lggr.Errorw("failed to handle event, backing off...", "err", err, "type", event.EventType, "nextRetryAt", event.nextRetryAt, "retryCount", event.retryCount)
						}
					} else {
						// It's not ready to execute yet, let's put it back on the pending queue.
						pendingEvents[event.id] = event

						reconcileReport.Backoffs[event.id] = event.nextRetryAt
						w.lggr.Debugw("skipping event, still in backoff", "nextRetryAt", event.nextRetryAt, "event", event.EventType, "id", event.id, "signature", event.signature)
					}
				}
			}

			w.lggr.Debugw("generated events to reconcile", "num", len(events), "report", reconcileReport)
		}
	}
}

// getTicker returns the ticker that the workflowRegistry will use to poll for events.  If the ticker
// is nil, then a default ticker is returned.
func (w *workflowRegistry) getTicker() <-chan time.Time {
	if w.ticker == nil {
		return time.NewTicker(defaultTickInterval).C
	}

	return w.ticker
}

type sequenceWithEventType struct {
	Sequence  types.Sequence
	EventType WorkflowRegistryEventType
}

func (w *workflowRegistry) newWorkflowRegistryContractReader(
	ctx context.Context,
) (ContractReader, error) {
	bc := types.BoundContract{
		Name:    WorkflowRegistryContractName,
		Address: w.workflowRegistryAddress,
	}

	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{
						string(ForceUpdateSecretsEvent),
						string(WorkflowActivatedEvent),
						string(WorkflowDeletedEvent),
						string(WorkflowPausedEvent),
						string(WorkflowRegisteredEvent),
						string(WorkflowUpdatedEvent),
					},
				},
				ContractABI: workflow_registry_wrapper.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					GetWorkflowMetadataListByDONMethodName: {
						ChainSpecificName: GetWorkflowMetadataListByDONMethodName,
					},
					string(ForceUpdateSecretsEvent): {
						ChainSpecificName: string(ForceUpdateSecretsEvent),
						ReadType:          evmtypes.Event,
					},
					string(WorkflowActivatedEvent): {
						ChainSpecificName: string(WorkflowActivatedEvent),
						ReadType:          evmtypes.Event,
					},
					string(WorkflowDeletedEvent): {
						ChainSpecificName: string(WorkflowDeletedEvent),
						ReadType:          evmtypes.Event,
					},
					string(WorkflowPausedEvent): {
						ChainSpecificName: string(WorkflowPausedEvent),
						ReadType:          evmtypes.Event,
					},
					string(WorkflowRegisteredEvent): {
						ChainSpecificName: string(WorkflowRegisteredEvent),
						ReadType:          evmtypes.Event,
					},
					string(WorkflowUpdatedEvent): {
						ChainSpecificName: string(WorkflowUpdatedEvent),
						ReadType:          evmtypes.Event,
					},
				},
			},
		},
	}

	marshalledCfg, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, err
	}

	reader, err := w.newContractReaderFn(ctx, marshalledCfg)
	if err != nil {
		return nil, err
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

// getWorkflowMetadata uses contract reader to query the contract for all workflow metadata using the method GetWorkflowMetadataListByDONMethodName
func (w *workflowRegistry) getWorkflowMetadata(ctx context.Context, don capabilities.DON, contractReader ContractReader) ([]GetWorkflowMetadata, *types.Head, error) {
	contractBinding := types.BoundContract{
		Address: w.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	readIdentifier := contractBinding.ReadIdentifier(GetWorkflowMetadataListByDONMethodName)
	params := GetWorkflowMetadataListByDONParams{
		DonID: don.ID,
		Start: 0,
		Limit: 0, // 0 tells the contract to return max pagination limit workflows on each call
	}

	var headAtLastRead *types.Head
	var allWorkflows []GetWorkflowMetadata

	for {
		var err error
		var workflows GetWorkflowMetadataListByDONReturnVal
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Finalized, params, &workflows)
		if err != nil {
			return []GetWorkflowMetadata{}, &types.Head{Height: "0"}, fmt.Errorf("failed to get lastest value with head data %w", err)
		}

		allWorkflows = append(allWorkflows, workflows.WorkflowMetadataList...)

		// if no workflows returned the end has been reached
		if len(workflows.WorkflowMetadataList) == 0 {
			break
		}

		params.Start += uint64(len(workflows.WorkflowMetadataList))
	}

	return allWorkflows, headAtLastRead, nil
}

// toWorkflowRegistryEventResponse converts a types.Sequence to a WorkflowRegistryEventResponse.
func toWorkflowRegistryEventResponse(
	log types.Sequence,
	evt WorkflowRegistryEventType,
	lggr logger.Logger,
) (workflowRegistryEvent, error) {
	resp := workflowRegistryEvent{
		Cursor: log.Cursor,
		Event: Event{
			EventType: evt,
		},
		Head: Head{
			Hash:      hex.EncodeToString(log.Hash),
			Height:    log.Height,
			Timestamp: log.Timestamp,
		},
	}

	dataAsValuesMap, err := values.WrapMap(log.Data)
	if err != nil {
		return workflowRegistryEvent{}, err
	}

	switch evt {
	case ForceUpdateSecretsEvent:
		var data ForceUpdateSecretsRequestedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
	case WorkflowRegisteredEvent:
		var data WorkflowRegisteredV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
		resp.DonID = &data.DonID
	case WorkflowUpdatedEvent:
		var data WorkflowUpdatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
		resp.DonID = &data.DonID
	case WorkflowPausedEvent:
		var data WorkflowPausedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
		resp.DonID = &data.DonID
	case WorkflowActivatedEvent:
		var data WorkflowActivatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
		resp.DonID = &data.DonID
	case WorkflowDeletedEvent:
		var data WorkflowDeletedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			return workflowRegistryEvent{}, err
		}
		resp.Event.Data = data
		resp.DonID = &data.DonID
	default:
		return workflowRegistryEvent{}, fmt.Errorf("unknown event type: %s", evt)
	}

	return resp, nil
}

type metrics struct {
	handleDuration metric.Int64Histogram
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, event string, success bool) {
	// Beholder doesn't support non-string attributes
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", successStr),
		attribute.String("eventType", event),
	))
}

func newMetrics() (*metrics, error) {
	h, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	return &metrics{handleDuration: h}, nil
}
