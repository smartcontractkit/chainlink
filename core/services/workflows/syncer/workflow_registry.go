package syncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval                    = 12 * time.Second
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
	WorkflowID   [32]byte
	Owner        []byte
	DonID        uint32
	Status       uint8
	WorkflowName string
	BinaryURL    string
	ConfigURL    string
	SecretsURL   string
}

type GetWorkflowMetadataListByDONReturnVal struct {
	WorkflowMetadataList []GetWorkflowMetadata
}

// WorkflowRegistryEvent is an event emitted by the WorkflowRegistry.  Each event is typed
// so that the consumer can determine how to handle the event.
type WorkflowRegistryEvent struct {
	Cursor    string
	Data      any
	EventType WorkflowRegistryEventType
	Head      Head
}

func (we WorkflowRegistryEvent) GetEventType() WorkflowRegistryEventType {
	return we.EventType
}

func (we WorkflowRegistryEvent) GetData() any {
	return we.Data
}

// WorkflowRegistryEventResponse is a response to either parsing a queried event or handling the event.
type WorkflowRegistryEventResponse struct {
	Err   error
	Event *WorkflowRegistryEvent
}

// WorkflowEventPollerConfig is the configuration needed to poll for events on a contract.  Currently
// requires the ContractEventName.
type WorkflowEventPollerConfig struct {
	QueryCount uint64
}

type WorkflowLoadConfig struct {
	FetchBatchSize int
}

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

// ContractReader is a subset of types.ContractReader defined locally to enable mocking.
type ContractReader interface {
	Start(ctx context.Context) error
	Close() error
	Bind(context.Context, []types.BoundContract) error
	QueryKey(context.Context, types.BoundContract, query.KeyFilter, query.LimitAndSort, any) ([]types.Sequence, error)
	GetLatestValueWithHeadData(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *types.Head, err error)
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

	// all goroutines are waited on with wg.
	wg sync.WaitGroup

	// ticker is the interval at which the workflowRegistry will poll the contract for events.
	ticker <-chan time.Time

	lggr                    logger.Logger
	workflowRegistryAddress string

	newContractReaderFn newContractReaderFn

	eventPollerCfg WorkflowEventPollerConfig
	eventTypes     []WorkflowRegistryEventType

	// eventsCh is read by the handler and each event is handled once received.
	eventsCh                    chan WorkflowRegistryEventResponse
	handler                     evtHandler
	initialWorkflowsStateLoader initialWorkflowsStateLoader

	// batchCh is a channel that receives batches of events from the contract query goroutines.
	batchCh chan []WorkflowRegistryEventResponse

	// heap is a min heap that merges batches of events from the contract query goroutines.  The
	// default min heap is sorted by block height.
	heap Heap

	workflowDonNotifier donNotifier

	reader ContractReader
}

// WithTicker allows external callers to provide a ticker to the workflowRegistry.  This is useful
// for overriding the default tick interval.
func WithTicker(ticker <-chan time.Time) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
	}
}

type evtHandler interface {
	Handle(ctx context.Context, event Event) error
}

type initialWorkflowsStateLoader interface {
	// LoadWorkflows loads all the workflows for the given donID from the contract.  Returns the head of the chain as of the
	// point in time at which the load occurred.
	LoadWorkflows(ctx context.Context, don capabilities.DON) (*types.Head, error)
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
	eventPollerConfig WorkflowEventPollerConfig,
	handler evtHandler,
	initialWorkflowsStateLoader initialWorkflowsStateLoader,
	workflowDonNotifier donNotifier,
	opts ...func(*workflowRegistry),
) *workflowRegistry {
	ets := []WorkflowRegistryEventType{
		ForceUpdateSecretsEvent,
		WorkflowActivatedEvent,
		WorkflowDeletedEvent,
		WorkflowPausedEvent,
		WorkflowRegisteredEvent,
		WorkflowUpdatedEvent,
	}
	wr := &workflowRegistry{
		lggr:                        lggr,
		newContractReaderFn:         newContractReaderFn,
		workflowRegistryAddress:     addr,
		eventPollerCfg:              eventPollerConfig,
		heap:                        newBlockHeightHeap(),
		stopCh:                      make(services.StopChan),
		eventTypes:                  ets,
		eventsCh:                    make(chan WorkflowRegistryEventResponse),
		batchCh:                     make(chan []WorkflowRegistryEventResponse, len(ets)),
		handler:                     handler,
		initialWorkflowsStateLoader: initialWorkflowsStateLoader,
		workflowDonNotifier:         workflowDonNotifier,
	}

	for _, opt := range opts {
		opt(wr)
	}
	return wr
}

// Start starts the workflowRegistry.  It starts two goroutines, one for querying the contract
// and one for handling the events.
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

			w.lggr.Debugw("Loading initial workflows for DON", "DON", don.ID)
			loadWorkflowsHead, err := w.initialWorkflowsStateLoader.LoadWorkflows(ctx, don)
			if err != nil {
				w.lggr.Errorw("failed to load workflows", "err", err)
				return
			}

			w.syncEventsLoop(ctx, loadWorkflowsHead.Height)
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			defer cancel()

			w.handlerLoop(ctx)
		}()

		return nil
	})
}

func (w *workflowRegistry) Close() error {
	return w.StopOnce(w.Name(), func() error {
		close(w.stopCh)
		w.wg.Wait()
		return nil
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

// handlerLoop handles the events that are emitted by the contract.
func (w *workflowRegistry) handlerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, open := <-w.eventsCh:
			if !open {
				return
			}

			if resp.Err != nil || resp.Event == nil {
				w.lggr.Errorw("failed to handle event", "err", resp.Err)
				continue
			}

			event := resp.Event
			w.lggr.Debugf("handling event: %+v", event)
			if err := w.handler.Handle(ctx, *event); err != nil {
				w.lggr.Errorw("failed to handle event", "event", event, "err", err)
				continue
			}
		}
	}
}

// syncEventsLoop polls the contract for events and passes them to a channel for handling.
func (w *workflowRegistry) syncEventsLoop(ctx context.Context, lastReadBlockNumber string) {
	var (
		// sendLog is a helper that sends a WorkflowRegistryEventResponse to the eventsCh in a
		// blocking way that will send the response or be canceled.
		sendLog = func(resp WorkflowRegistryEventResponse) {
			select {
			case w.eventsCh <- resp:
			case <-ctx.Done():
			}
		}

		ticker = w.getTicker()

		signals = make(map[WorkflowRegistryEventType]chan struct{}, 0)
	)

	// critical failure if there is no reader, the loop will exit and the parent context will be
	// canceled.
	reader, err := w.getContractReader(ctx)
	if err != nil {
		w.lggr.Criticalf("contract reader unavailable : %s", err)
		return
	}

	// fan out and query for each event type
	for i := 0; i < len(w.eventTypes); i++ {
		signal := make(chan struct{}, 1)
		signals[w.eventTypes[i]] = signal
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()

			queryEvent(
				ctx,
				signal,
				w.lggr,
				reader,
				lastReadBlockNumber,
				queryEventConfig{
					ContractName:              WorkflowRegistryContractName,
					ContractAddress:           w.workflowRegistryAddress,
					WorkflowEventPollerConfig: w.eventPollerCfg,
				},
				w.eventTypes[i],
				w.batchCh,
			)
		}()
	}

	// Periodically send a signal to all the queryEvent goroutines to query the contract
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			w.lggr.Debugw("Syncing with WorkflowRegistry")
			// for each event type, send a signal for it to execute a query and produce a new
			// batch of event logs
			for i := 0; i < len(w.eventTypes); i++ {
				signal := signals[w.eventTypes[i]]
				select {
				case signal <- struct{}{}:
				case <-ctx.Done():
					return
				}
			}

			// block on fan-in until all fetched event logs are sent to the handlers
			w.orderAndSend(
				ctx,
				len(w.eventTypes),
				w.batchCh,
				sendLog,
			)
		}
	}
}

// orderAndSend reads n batches from the batch channel, heapifies all the batches then dequeues
// the min heap via the sendLog function.
func (w *workflowRegistry) orderAndSend(
	ctx context.Context,
	batchCount int,
	batchCh <-chan []WorkflowRegistryEventResponse,
	sendLog func(WorkflowRegistryEventResponse),
) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-batchCh:
			for _, response := range batch {
				w.heap.Push(response)
			}
			batchCount--

			// If we have received responses for all the events, then we can drain the heap.
			if batchCount == 0 {
				for w.heap.Len() > 0 {
					sendLog(w.heap.Pop())
				}
				return
			}
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

// getContractReader initializes a contract reader if needed, otherwise returns the existing
// reader.
func (w *workflowRegistry) getContractReader(ctx context.Context) (ContractReader, error) {
	c := types.BoundContract{
		Name:    WorkflowRegistryContractName,
		Address: w.workflowRegistryAddress,
	}

	if w.reader == nil {
		reader, err := getWorkflowRegistryEventReader(ctx, w.newContractReaderFn, c)
		if err != nil {
			return nil, err
		}

		w.reader = reader
	}

	return w.reader, nil
}

type queryEventConfig struct {
	ContractName    string
	ContractAddress string
	WorkflowEventPollerConfig
}

// queryEvent queries the contract for events of the given type on each tick from the ticker.
// Sends a batch of event logs to the batch channel.  The batch represents all the
// event logs read since the last query.  Loops until the context is canceled.
func queryEvent(
	ctx context.Context,
	ticker <-chan struct{},
	lggr logger.Logger,
	reader ContractReader,
	lastReadBlockNumber string,
	cfg queryEventConfig,
	et WorkflowRegistryEventType,
	batchCh chan<- []WorkflowRegistryEventResponse,
) {
	// create query
	var (
		logData      values.Value
		cursor       = ""
		limitAndSort = query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: cfg.QueryCount},
		}
		bc = types.BoundContract{
			Name:    cfg.ContractName,
			Address: cfg.ContractAddress,
		}
	)

	// Loop until canceled
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			responseBatch := []WorkflowRegistryEventResponse{}

			if cursor != "" {
				limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, cfg.QueryCount)
			}

			logs, err := reader.QueryKey(
				ctx,
				bc,
				query.KeyFilter{
					Key: string(et),
					Expressions: []query.Expression{
						query.Confidence(primitives.Finalized),
						query.Block(lastReadBlockNumber, primitives.Gte),
					},
				},
				limitAndSort,
				&logData,
			)
			lcursor := cursor
			if lcursor == "" {
				lcursor = "empty"
			}
			lggr.Debugw("QueryKeys called", "logs", len(logs), "eventType", et, "lastReadBlockNumber", lastReadBlockNumber, "logCursor", lcursor)

			if err != nil {
				lggr.Errorw("QueryKey failure", "err", err)
				continue
			}

			// ChainReader QueryKey API provides logs including the cursor value and not
			// after the cursor value. If the response only consists of the log corresponding
			// to the cursor and no log after it, then we understand that there are no new
			// logs
			if len(logs) == 1 && logs[0].Cursor == cursor {
				lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}

			for _, log := range logs {
				if log.Cursor == cursor {
					continue
				}

				responseBatch = append(responseBatch, toWorkflowRegistryEventResponse(log, et, lggr))
				cursor = log.Cursor
			}
			batchCh <- responseBatch
		}
	}
}

func getWorkflowRegistryEventReader(
	ctx context.Context,
	newReaderFn newContractReaderFn,
	bc types.BoundContract,
) (ContractReader, error) {
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

	reader, err := newReaderFn(ctx, marshalledCfg)
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

type workflowAsEvent struct {
	Data      WorkflowRegistryWorkflowRegisteredV1
	EventType WorkflowRegistryEventType
}

func (r workflowAsEvent) GetEventType() WorkflowRegistryEventType {
	return r.EventType
}

func (r workflowAsEvent) GetData() any {
	return r.Data
}

type workflowRegistryContractLoader struct {
	lggr                    logger.Logger
	workflowRegistryAddress string
	newContractReaderFn     newContractReaderFn
	handler                 evtHandler
}

func NewWorkflowRegistryContractLoader(
	lggr logger.Logger,
	workflowRegistryAddress string,
	newContractReaderFn newContractReaderFn,
	handler evtHandler,
) *workflowRegistryContractLoader {
	return &workflowRegistryContractLoader{
		lggr:                    lggr.Named("WorkflowRegistryContractLoader"),
		workflowRegistryAddress: workflowRegistryAddress,
		newContractReaderFn:     newContractReaderFn,
		handler:                 handler,
	}
}

func (l *workflowRegistryContractLoader) LoadWorkflows(ctx context.Context, don capabilities.DON) (*types.Head, error) {
	// Build the ContractReader config
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			WorkflowRegistryContractName: {
				ContractABI: workflow_registry_wrapper.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					GetWorkflowMetadataListByDONMethodName: {
						ChainSpecificName: GetWorkflowMetadataListByDONMethodName,
					},
				},
			},
		},
	}

	contractReaderCfgBytes, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config: %w", err)
	}

	contractReader, err := l.newContractReaderFn(ctx, contractReaderCfgBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract reader: %w", err)
	}

	err = contractReader.Bind(ctx, []types.BoundContract{{Name: WorkflowRegistryContractName, Address: l.workflowRegistryAddress}})
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract reader: %w", err)
	}

	contractBinding := types.BoundContract{
		Address: l.workflowRegistryAddress,
		Name:    WorkflowRegistryContractName,
	}

	readIdentifier := contractBinding.ReadIdentifier(GetWorkflowMetadataListByDONMethodName)
	params := GetWorkflowMetadataListByDONParams{
		DonID: don.ID,
		Start: 0,
		Limit: 0, // 0 tells the contract to return max pagination limit workflows on each call
	}

	var headAtLastRead *types.Head
	for {
		var err error
		var workflows GetWorkflowMetadataListByDONReturnVal
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Finalized, params, &workflows)
		if err != nil {
			return nil, fmt.Errorf("failed to get workflow metadata for don %w", err)
		}

		l.lggr.Debugw("Rehydrating existing workflows", "len", len(workflows.WorkflowMetadataList))
		for _, workflow := range workflows.WorkflowMetadataList {
			toRegisteredEvent := WorkflowRegistryWorkflowRegisteredV1{
				WorkflowID:    workflow.WorkflowID,
				WorkflowOwner: workflow.Owner,
				DonID:         workflow.DonID,
				Status:        workflow.Status,
				WorkflowName:  workflow.WorkflowName,
				BinaryURL:     workflow.BinaryURL,
				ConfigURL:     workflow.ConfigURL,
				SecretsURL:    workflow.SecretsURL,
			}
			if err = l.handler.Handle(ctx, workflowAsEvent{
				Data:      toRegisteredEvent,
				EventType: WorkflowRegisteredEvent,
			}); err != nil {
				l.lggr.Errorf("failed to handle workflow registration: %s", err)
			}
		}

		if len(workflows.WorkflowMetadataList) == 0 {
			break
		}

		params.Start += uint64(len(workflows.WorkflowMetadataList))
	}

	return headAtLastRead, nil
}

// toWorkflowRegistryEventResponse converts a types.Sequence to a WorkflowRegistryEventResponse.
func toWorkflowRegistryEventResponse(
	log types.Sequence,
	evt WorkflowRegistryEventType,
	lggr logger.Logger,
) WorkflowRegistryEventResponse {
	resp := WorkflowRegistryEventResponse{
		Event: &WorkflowRegistryEvent{
			Cursor:    log.Cursor,
			EventType: evt,
			Head: Head{
				Hash:      hex.EncodeToString(log.Hash),
				Height:    log.Height,
				Timestamp: log.Timestamp,
			},
		},
	}

	dataAsValuesMap, err := values.WrapMap(log.Data)
	if err != nil {
		return WorkflowRegistryEventResponse{
			Err: err,
		}
	}

	switch evt {
	case ForceUpdateSecretsEvent:
		var data WorkflowRegistryForceUpdateSecretsRequestedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	case WorkflowRegisteredEvent:
		var data WorkflowRegistryWorkflowRegisteredV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	case WorkflowUpdatedEvent:
		var data WorkflowRegistryWorkflowUpdatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	case WorkflowPausedEvent:
		var data WorkflowRegistryWorkflowPausedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	case WorkflowActivatedEvent:
		var data WorkflowRegistryWorkflowActivatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	case WorkflowDeletedEvent:
		var data WorkflowRegistryWorkflowDeletedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
	default:
		lggr.Errorf("unknown event type: %s", evt)
		resp.Event = nil
		resp.Err = fmt.Errorf("unknown event type: %s", evt)
	}

	return resp
}
