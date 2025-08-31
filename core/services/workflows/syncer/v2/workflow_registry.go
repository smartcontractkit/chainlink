package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math/big"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/versioning"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval          = 12 * time.Second
	defaultRetryInterval         = 12 * time.Second
	defaultMaxRetryInterval      = 5 * time.Minute
	WorkflowRegistryContractName = "WorkflowRegistry"

	GetWorkflowsByDONMethodName      = "getWorkflowListByDON"
	GetAllowlistedRequestsMethodName = "getAllowlistedRequests"

	defaultTickIntervalForAllowlistedRequests = 5 * time.Second

	// MaxResultsPerQuery defines the maximum number of results that can be queried in a single request.
	// The default value of 1,000 was chosen based on expected system performance and typical use cases.
	MaxResultsPerQuery = uint64(1_000)
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

	// events sent to the event channel to be handled.
	eventCh chan Event

	// all goroutines are waited on with wg.
	wg sync.WaitGroup

	// ticker is the interval at which the workflowRegistry will
	// poll the contract for events, and poll the contract for the latest workflow metadata.
	ticker <-chan time.Time

	lggr                    logger.Logger
	workflowRegistryAddress string

	allowListedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	allowListedMu       sync.RWMutex

	contractReaderFn versioning.ContractReaderFactory

	config Config

	handler evtHandler

	workflowDonNotifier donNotifier

	metrics *metrics

	engineRegistry *EngineRegistry

	retryInterval    time.Duration
	maxRetryInterval time.Duration
	clock            clockwork.Clock
}

type evtHandler interface {
	io.Closer
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

// NewWorkflowRegistry returns a new v2 workflowRegistry.
func NewWorkflowRegistry(
	lggr logger.Logger,
	contractReaderFn versioning.ContractReaderFactory,
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
		contractReaderFn:        contractReaderFn,
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

		contractReader, err := w.newWorkflowRegistryContractReader(context.Background())
		if err != nil {
			w.lggr.Criticalf("contract reader unavailable : %s", err)
			return errors.New("failed to create contract reader: " + err.Error())
		}

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

			// Start goroutines to gather changes from Workflow Registry contract
			w.syncUsingReconciliationStrategy(ctx, don, contractReader)
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()
			// Start goroutines to gather allowlisted requests from Workflow Registry contract
			w.syncAllowlistedRequests(ctx, contractReader)
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

func (w *workflowRegistry) handleWithMetrics(ctx context.Context, event Event) error {
	start := time.Now()
	err := w.handler.Handle(ctx, event)
	totalDuration := time.Since(start)
	w.metrics.recordHandleDuration(ctx, totalDuration, string(event.Name), err == nil)
	return err
}

// generateReconciliationEvents compares the workflow registry workflow metadata state against the engine registry's state.
// Differences are handled by the event handler by creating events that are sent to the events channel for handling.
func (w *workflowRegistry) generateReconciliationEvents(_ context.Context, pendingEvents map[string]*reconciliationEvent, workflowMetadata []WorkflowMetadataView) ([]*reconciliationEvent, error) {
	var events []*reconciliationEvent

	// Keep track of which of the engines in the engineRegistry have been touched
	workflowsSeen := map[string]bool{}
	for _, wfMeta := range workflowMetadata {
		id := wfMeta.WorkflowID.Hex()
		engineFound := w.engineRegistry.Contains(wfMeta.WorkflowID)

		switch wfMeta.Status {
		case WorkflowStatusActive:
			switch engineFound {
			// if the workflow is active, but unable to get engine from the engine registry
			// then handle as registered event
			case false:
				signature := fmt.Sprintf("%s-%s-%s", WorkflowRegistered, id, toSpecStatus(wfMeta.Status))

				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature == signature {
					events = append(events, pendingEvents[id])
					delete(pendingEvents, id)
					continue
				}

				delete(pendingEvents, id)

				toRegisteredEvent := WorkflowRegisteredEvent{
					WorkflowID:    wfMeta.WorkflowID,
					WorkflowOwner: wfMeta.Owner,
					CreatedAt:     wfMeta.CreatedAt,
					Status:        wfMeta.Status,
					WorkflowName:  wfMeta.WorkflowName,
					BinaryURL:     wfMeta.BinaryURL,
					ConfigURL:     wfMeta.ConfigURL,
					Tag:           wfMeta.Tag,
					Attributes:    wfMeta.Attributes,
				}
				events = append(events, &reconciliationEvent{
					Event: Event{
						Data: toRegisteredEvent,
						Name: WorkflowRegistered,
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
			switch engineFound {
			case false:
				// Account for a state change from active to paused, by checking
				// whether an existing pendingEvent exists.
				// We do this regardless of whether we have an event to handle or not, since this ensures
				// we correctly handle the state of pending events in the following situation:
				// - we registered an active workflow, but it failed to process successfully
				// - we then paused the workflow; this should clear the pending event
				signature := fmt.Sprintf("%s-%s-%s", WorkflowPaused, id, toSpecStatus(wfMeta.Status))
				if _, ok := pendingEvents[id]; ok && pendingEvents[id].signature != signature {
					delete(pendingEvents, id)
				}
			case true:
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
			}
			events = append(
				[]*reconciliationEvent{
					{
						Event: Event{
							Data: toDeletedEvent,
							Name: WorkflowDeleted,
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
	// the workflow no longer exists in the workflow registry contract
	for id, event := range pendingEvents {
		if event.Name == WorkflowRegistered {
			existsInMetadata := false
			for _, wfMeta := range workflowMetadata {
				if wfMeta.WorkflowID.Hex() == event.Data.(WorkflowRegisteredEvent).WorkflowID.Hex() {
					existsInMetadata = true
					break
				}
			}
			if !existsInMetadata {
				delete(pendingEvents, id)
			}
		}
	}

	if len(pendingEvents) != 0 {
		return nil, fmt.Errorf("invariant violation: some pending events were not handled in the reconcile loop: keys=%+v, len=%d", maps.Keys(pendingEvents), len(pendingEvents))
	}

	return events, nil
}

func (w *workflowRegistry) syncAllowlistedRequests(ctx context.Context, contractReader types.ContractReader) {
	ticker := time.NewTicker(defaultTickIntervalForAllowlistedRequests).C
	w.lggr.Debug("starting syncAllowlistedRequests")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down syncAllowlistedRequests, %s", ctx.Err())
			return
		case <-ticker:
			allowListedRequests, head, err := w.getAllowlistedRequests(ctx, contractReader)
			if err != nil {
				w.lggr.Errorw("failed to call getAllowlistedRequests", "err", err)
				continue
			}
			w.allowListedMu.Lock()
			w.allowListedRequests = allowListedRequests
			w.allowListedMu.Unlock()
			w.lggr.Debugw("fetched allowlisted requests", "numRequests", len(allowListedRequests), "blockHeight", head.Height)
		}
	}
}

// syncUsingReconciliationStrategy syncs workflow registry contract state by polling the workflow metadata state and comparing to local state.
// NOTE: In this mode paused states will be treated as a deleted workflow. Workflows will not be registered as paused.
func (w *workflowRegistry) syncUsingReconciliationStrategy(ctx context.Context, don capabilities.DON, contractReader types.ContractReader) {
	ticker := w.getTicker()
	pendingEvents := map[string]*reconciliationEvent{}
	w.lggr.Debug("running readRegistryStateLoop")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down readRegistryStateLoop")
			return
		case <-ticker:
			workflowMetadata, head, err := w.getWorkflowMetadata(ctx, don, contractReader)
			if err != nil {
				w.lggr.Errorw("failed to get registry state", "err", err)
				continue
			}
			w.lggr.Debugw("preparing events to reconcile", "numWorkflowMetadata", len(workflowMetadata), "blockHeight", head.Height, "numPendingEvents", len(pendingEvents))
			events, err := w.generateReconciliationEvents(ctx, pendingEvents, workflowMetadata)
			if err != nil {
				w.lggr.Errorw("failed to generate reconciliation events", "err", err)
				continue
			}
			w.lggr.Debugw("generated events to reconcile", "num", len(events), "events", events)

			pendingEvents = map[string]*reconciliationEvent{}

			// Send events generated from differences to the handler
			reconcileReport := newReconcileReport()
			for _, event := range events {
				select {
				case <-ctx.Done():
					w.lggr.Debug("readRegistryStateLoop stopped during processing")
					return
				default:
					reconcileReport.NumEventsByType[string(event.Name)]++

					if event.retryCount == 0 || w.clock.Now().After(event.nextRetryAt) {
						err := w.handleWithMetrics(ctx, event.Event)
						if err != nil {
							event.updateNextRetryFor(w.clock, w.retryInterval, w.maxRetryInterval)

							pendingEvents[event.id] = event

							reconcileReport.Backoffs[event.id] = event.nextRetryAt
							w.lggr.Errorw("failed to handle event, backing off...", "err", err, "type", event.Name, "nextRetryAt", event.nextRetryAt, "retryCount", event.retryCount)
						}
					} else {
						// It's not ready to execute yet, let's put it back on the pending queue.
						pendingEvents[event.id] = event

						reconcileReport.Backoffs[event.id] = event.nextRetryAt
						w.lggr.Debugw("skipping event, still in backoff", "nextRetryAt", event.nextRetryAt, "event", event.Name, "id", event.id, "signature", event.signature)
					}
				}
			}

			w.lggr.Debugw("reconciled events", "report", reconcileReport)
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

func (w *workflowRegistry) newWorkflowRegistryContractReader(
	ctx context.Context,
) (types.ContractReader, error) {
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper_v2.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					GetWorkflowsByDONMethodName: {
						ChainSpecificName: GetWorkflowsByDONMethodName,
						ReadType:          evmtypes.Method,
					},
					GetAllowlistedRequestsMethodName: {
						ChainSpecificName: GetAllowlistedRequestsMethodName,
						ReadType:          evmtypes.Method,
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

// getWorkflowMetadata uses contract reader to query the contract for all workflow metadata using the method getWorkflowListByDON
func (w *workflowRegistry) getWorkflowMetadata(ctx context.Context, don capabilities.DON, contractReader types.ContractReader) ([]WorkflowMetadataView, *types.Head, error) {
	contractBinding := types.BoundContract{
		Address: w.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	readIdentifier := contractBinding.ReadIdentifier(GetWorkflowsByDONMethodName)
	var headAtLastRead *types.Head
	var allWorkflows []WorkflowMetadataView

	for _, family := range don.Families {
		params := GetWorkflowListByDONParams{
			DonFamily: family,
			Start:     big.NewInt(0),
			Limit:     big.NewInt(int64(MaxResultsPerQuery)), //nolint:gosec // safe conversion
		}

		for {
			var err error
			var workflows struct {
				List []workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView
			}

			headAtLastRead, err = contractReader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Finalized, params, &workflows)
			if err != nil {
				return []WorkflowMetadataView{}, &types.Head{Height: "0"}, fmt.Errorf("failed to get lastest value with head data %w", err)
			}

			for _, wfMeta := range workflows.List {
				// TODO: https://smartcontract-it.atlassian.net/browse/CAPPL-1021 load balance across workflow nodes in DON Family
				allWorkflows = append(allWorkflows, WorkflowMetadataView{
					WorkflowID:   wfMeta.WorkflowId,
					Owner:        wfMeta.Owner.Bytes(),
					CreatedAt:    wfMeta.CreatedAt,
					Status:       wfMeta.Status,
					WorkflowName: wfMeta.WorkflowName,
					BinaryURL:    wfMeta.BinaryUrl,
					ConfigURL:    wfMeta.ConfigUrl,
					Tag:          wfMeta.Tag,
					Attributes:   wfMeta.Attributes,
					DonFamily:    wfMeta.DonFamily,
				})
			}

			// if less workflows than limit, then we have reached the end of the list
			if uint64(len(workflows.List)) < MaxResultsPerQuery {
				break
			}

			// otherwise, increment the start parameter and continue to fetch more workflows
			params.Start.Add(params.Start, big.NewInt(int64(len(workflows.List))))
		}
	}

	if headAtLastRead == nil {
		return allWorkflows, &types.Head{Height: "0"}, nil
	}

	return allWorkflows, headAtLastRead, nil
}

func (w *workflowRegistry) GetAllowlistedRequests(_ context.Context) []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest {
	w.allowListedMu.RLock()
	defer w.allowListedMu.RUnlock()
	return w.allowListedRequests
}

// GetAllowlistedRequests uses contract reader to query the contract for all allowlisted requests
func (w *workflowRegistry) getAllowlistedRequests(ctx context.Context, contractReader types.ContractReader) ([]workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest, *types.Head, error) {
	contractBinding := types.BoundContract{
		Address: w.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	readIdentifier := contractBinding.ReadIdentifier(GetAllowlistedRequestsMethodName)
	var headAtLastRead *types.Head
	var allAllowlistedRequests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
	params := GetAllowlistedRequestsParams{
		Start: big.NewInt(0),
		Limit: big.NewInt(int64(MaxResultsPerQuery)), //nolint:gosec // safe conversion
	}

	for {
		var err error
		var results struct {
			Requests []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest
			err      error
		}

		// Read at confidenceLevel Unconfirmed, since we want to see new allowlisted requests asap, even if awaiting finalization.
		// The delay in detecting allowlisted requests will directly affect user experience.
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Unconfirmed, params, &results)
		if err != nil {
			return []workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{}, &types.Head{Height: "0"}, errors.New("failed to get lastest value with head data. error: " + err.Error())
		}

		for _, request := range results.Requests {
			// TODO: https://smartcontract-it.atlassian.net/browse/CAPPL-1021 load balance across workflow nodes in DON Family
			allAllowlistedRequests = append(allAllowlistedRequests, workflow_registry_wrapper_v2.WorkflowRegistryOwnerAllowlistedRequest{
				RequestDigest:   request.RequestDigest,
				Owner:           request.Owner,
				ExpiryTimestamp: request.ExpiryTimestamp,
			})
		}

		// if less results than limit, then we have reached the end of the list
		if uint64(len(results.Requests)) < MaxResultsPerQuery {
			break
		}

		// otherwise, increment the start parameter and continue to fetch more workflows
		params.Start = params.Start.Add(params.Start, big.NewInt(int64(len(results.Requests))))
	}

	if headAtLastRead == nil {
		return allAllowlistedRequests, &types.Head{Height: "0"}, nil
	}

	return allAllowlistedRequests, headAtLastRead, nil
}
