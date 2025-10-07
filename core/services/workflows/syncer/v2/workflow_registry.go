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
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/versioning"
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
	contractReader   types.ContractReader

	config Config

	handler evtHandler

	workflowDonNotifier donNotifier

	metrics *metrics

	engineRegistry *EngineRegistry

	retryInterval    time.Duration
	maxRetryInterval time.Duration
	clock            clockwork.Clock

	hooks Hooks
}

type Hooks struct {
	OnStartFailure func(error)
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
		stopCh:                  make(services.StopChan),
		handler:                 handler,
		workflowDonNotifier:     workflowDonNotifier,
		metrics:                 m,
		engineRegistry:          engineRegistry,
		retryInterval:           defaultRetryInterval,
		maxRetryInterval:        defaultMaxRetryInterval,
		clock:                   clockwork.NewRealClock(),
		hooks: Hooks{
			OnStartFailure: func(_ error) {},
		},
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
		initDoneCh := make(chan struct{})

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()

			w.lggr.Debugw("Waiting for DON...")
			if _, err := w.workflowDonNotifier.WaitForDon(ctx); err != nil {
				w.lggr.Errorw("failed to wait for don", "err", err)
				w.hooks.OnStartFailure(err)
				return
			}

			// Async initialization of contract reader because there is an on-chain
			// call dependency.  Blocking on initialization results in a
			// deadlock.  Instead wait until the node has identified it's DON
			// as a proxy for a DON and on-chain ready state .
			reader, err := w.newWorkflowRegistryContractReader(ctx)
			if err != nil {
				w.lggr.Criticalf("contract reader unavailable : %s", err)
				w.hooks.OnStartFailure(err)
				cancel()
				return
			}

			w.contractReader = reader
			close(initDoneCh)
			w.lggr.Debugw("Received DON and set ContractReader")
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
			don, err := w.workflowDonNotifier.WaitForDon(ctx)
			if err != nil {
				w.hooks.OnStartFailure(fmt.Errorf("failed to start workflow sync strategy: %w", err))
				return
			}
			w.syncUsingReconciliationStrategy(ctx, don)
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

// generateReconciliationEvents compares the workflow registry workflow metadata state against the engine registry's state.
// Differences are handled by the event handler by creating events that are sent to the events channel for handling.
func (w *workflowRegistry) generateReconciliationEvents(_ context.Context, pendingEvents map[string]*reconciliationEvent, workflowMetadata []WorkflowMetadataView, head *types.Head) ([]*reconciliationEvent, error) {
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
				}
				events = append(events, &reconciliationEvent{
					Event: Event{
						Data: toActivatedEvent,
						Name: WorkflowActivated,
						Head: localHead,
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
				}
				events = append(
					[]*reconciliationEvent{
						{
							Event: Event{
								Data: toPausedEvent,
								Name: WorkflowPaused,
								Head: localHead,
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
							Head: localHead,
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

func (w *workflowRegistry) syncAllowlistedRequests(ctx context.Context) {
	ticker := w.getTicker(defaultTickIntervalForAllowlistedRequests)
	w.lggr.Debug("starting syncAllowlistedRequests")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down syncAllowlistedRequests, %s", ctx.Err())
			return
		case <-ticker:
			allowListedRequests, head, err := w.getAllowlistedRequests(ctx, w.contractReader)
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
func (w *workflowRegistry) syncUsingReconciliationStrategy(ctx context.Context, don capabilities.DON) {
	ticker := w.getTicker(defaultTickInterval)
	pendingEvents := map[string]*reconciliationEvent{}
	w.lggr.Debug("running readRegistryStateLoop")
	for {
		select {
		case <-ctx.Done():
			w.lggr.Debug("shutting down readRegistryStateLoop")
			return
		case <-ticker:
			w.lggr.Debugw("fetching workflow registry metadata", "don", don.Families)
			workflowMetadata, head, err := w.getWorkflowMetadata(ctx, don, w.contractReader)
			if err != nil {
				w.lggr.Errorw("failed to get registry state", "err", err)
				continue
			}
			w.lggr.Debugw("preparing events to reconcile", "numWorkflowMetadata", len(workflowMetadata), "blockHeight", head.Height, "numPendingEvents", len(pendingEvents))
			events, err := w.generateReconciliationEvents(ctx, pendingEvents, workflowMetadata, head)
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

// validateWorkflowMetadata logs warnings for incomplete workflow metadata from contract
func validateWorkflowMetadata(wfMeta workflow_registry_wrapper_v2.WorkflowRegistryWorkflowMetadataView, lggr logger.Logger) {
	if isEmptyWorkflowID(wfMeta.WorkflowId) {
		lggr.Warnw("Workflow has empty WorkflowID from contract",
			"workflowName", wfMeta.WorkflowName,
			"owner", hex.EncodeToString(wfMeta.Owner.Bytes()),
			"binaryURL", wfMeta.BinaryUrl,
			"configURL", wfMeta.ConfigUrl)
	}

	if len(wfMeta.Owner.Bytes()) == 0 {
		lggr.Warnw("Workflow has empty Owner from contract",
			"workflowID", hex.EncodeToString(wfMeta.WorkflowId[:]),
			"workflowName", wfMeta.WorkflowName,
			"binaryURL", wfMeta.BinaryUrl,
			"configURL", wfMeta.ConfigUrl)
	}

	if wfMeta.BinaryUrl == "" || wfMeta.ConfigUrl == "" {
		lggr.Warnw("Workflow has empty BinaryURL or ConfigURL from contract",
			"workflowID", hex.EncodeToString(wfMeta.WorkflowId[:]),
			"workflowName", wfMeta.WorkflowName,
			"owner", hex.EncodeToString(wfMeta.Owner.Bytes()),
			"binaryURL", wfMeta.BinaryUrl,
			"configURL", wfMeta.ConfigUrl)
	}
}

func (w *workflowRegistry) newWorkflowRegistryContractReader(
	ctx context.Context,
) (types.ContractReader, error) {
	contractReaderCfg := config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper_v2.WorkflowRegistryABI,
				Configs: map[string]*config.ChainReaderDefinition{
					GetWorkflowsByDONMethodName: {
						ChainSpecificName: GetWorkflowsByDONMethodName,
						ReadType:          config.Method,
					},
					GetAllowlistedRequestsMethodName: {
						ChainSpecificName: GetAllowlistedRequestsMethodName,
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

// getWorkflowMetadata uses contract reader to query the contract for all workflow metadata using the method getWorkflowListByDON
func (w *workflowRegistry) getWorkflowMetadata(ctx context.Context, don capabilities.DON, contractReader types.ContractReader) ([]WorkflowMetadataView, *types.Head, error) {
	if contractReader == nil {
		return nil, nil, errors.New("cannot fetch workflow metadata: nil contract reader")
	}
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
				// Log warnings for incomplete metadata but don't skip processing
				validateWorkflowMetadata(wfMeta, w.lggr)

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
	if contractReader == nil {
		return nil, nil, errors.New("cannot fetch allow listed requests: nil contract reader")
	}
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
