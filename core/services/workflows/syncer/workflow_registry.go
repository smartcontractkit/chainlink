package syncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
)

type Head struct {
	Hash      string
	Height    string
	Timestamp uint64
}

// WorkflowRegistryEvent is an event emitted by the WorkflowRegistry.  Each event is typed
// so that the consumer can determine how to handle the event.
type WorkflowRegistryEvent struct {
	Cursor    string
	Data      any
	EventType WorkflowRegistryEventType
	Head      Head
}

// WorkflowRegistryEventResponse is a response to either parsing a queried event or handling the event.
type WorkflowRegistryEventResponse struct {
	Err   error
	Event *WorkflowRegistryEvent
}

// ContractEventPollerConfig is the configuration needed to poll for events on a contract.  Currently
// requires the ContractEventName.
//
// TODO(mstreet3): Use LookbackBlocks instead of StartBlockNum
type ContractEventPollerConfig struct {
	ContractName    string
	ContractAddress string
	StartBlockNum   uint64
	QueryCount      uint64
}

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type ContractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

// ContractReader is a subset of types.ContractReader defined locally to enable mocking.
type ContractReader interface {
	Bind(context.Context, []types.BoundContract) error
	QueryKey(context.Context, types.BoundContract, query.KeyFilter, query.LimitAndSort, any) ([]types.Sequence, error)
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

	lggr    logger.Logger
	orm     WorkflowRegistryDS
	reader  ContractReader
	gateway FetcherFunc

	// initReader allows the workflowRegistry to initialize a contract reader if one is not provided
	// and separates the contract reader initialization from the workflowRegistry start up.
	initReader func(context.Context, logger.Logger, ContractReaderFactory, types.BoundContract) (types.ContractReader, error)
	relayer    ContractReaderFactory

	cfg        ContractEventPollerConfig
	eventTypes []WorkflowRegistryEventType

	// eventsCh is read by the handler and each event is handled once received.
	eventsCh chan WorkflowRegistryEventResponse
	handler  *eventHandler

	// batchCh is a channel that receives batches of events from the contract query goroutines.
	batchCh chan []WorkflowRegistryEventResponse

	// heap is a min heap that merges batches of events from the contract query goroutines.  The
	// default min heap is sorted by block height.
	heap Heap

	workflowStore  store.Store
	capRegistry    core.CapabilitiesRegistry
	engineRegistry *engineRegistry
}

// WithTicker allows external callers to provide a ticker to the workflowRegistry.  This is useful
// for overriding the default tick interval.
func WithTicker(ticker <-chan time.Time) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
	}
}

func WithReader(reader types.ContractReader) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.reader = reader
	}
}

// NewWorkflowRegistry returns a new workflowRegistry.
// Only queries for WorkflowRegistryForceUpdateSecretsRequestedV1 events.
func NewWorkflowRegistry[T ContractReader](
	lggr logger.Logger,
	orm WorkflowRegistryDS,
	reader T,
	gateway FetcherFunc,
	addr string,
	workflowStore store.Store,
	capRegistry core.CapabilitiesRegistry,
	opts ...func(*workflowRegistry),
) *workflowRegistry {
	ets := []WorkflowRegistryEventType{ForceUpdateSecretsEvent}
	wr := &workflowRegistry{
		lggr:           lggr.Named(name),
		orm:            orm,
		reader:         reader,
		gateway:        gateway,
		workflowStore:  workflowStore,
		capRegistry:    capRegistry,
		engineRegistry: newEngineRegistry(),
		cfg: ContractEventPollerConfig{
			ContractName:    ContractName,
			ContractAddress: addr,
			QueryCount:      20,
			StartBlockNum:   0,
		},
		initReader: newReader,
		heap:       newBlockHeightHeap(),
		stopCh:     make(services.StopChan),
		eventTypes: ets,
		eventsCh:   make(chan WorkflowRegistryEventResponse),
		batchCh:    make(chan []WorkflowRegistryEventResponse, len(ets)),
	}
	wr.handler = newEventHandler(wr.lggr, wr.orm, wr.gateway, wr.workflowStore, wr.capRegistry,
		wr.engineRegistry,
	)
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

			w.syncEventsLoop(ctx)
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

func (w *workflowRegistry) SecretsFor(ctx context.Context, workflowOwner, workflowName string) (map[string]string, error) {
	return nil, errors.New("not implemented")
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
				w.lggr.Errorf("failed to handle event: %+v", resp.Err)
				continue
			}

			event := resp.Event
			w.lggr.Debugf("handling event: %+v", event)
			if err := w.handler.Handle(ctx, *event); err != nil {
				w.lggr.Errorf("failed to handle event: %+v", event)
				continue
			}
		}
	}
}

// syncEventsLoop polls the contract for events and passes them to a channel for handling.
func (w *workflowRegistry) syncEventsLoop(ctx context.Context) {
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
				w.cfg,
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
		Name:    w.cfg.ContractName,
		Address: w.cfg.ContractAddress,
	}

	if w.reader == nil {
		reader, err := w.initReader(ctx, w.lggr, w.relayer, c)
		if err != nil {
			return nil, err
		}

		w.reader = reader
	}

	return w.reader, nil
}

// queryEvent queries the contract for events of the given type on each tick from the ticker.
// Sends a batch of event logs to the batch channel.  The batch represents all the
// event logs read since the last query.  Loops until the context is canceled.
func queryEvent(
	ctx context.Context,
	ticker <-chan struct{},
	lggr logger.Logger,
	reader ContractReader,
	cfg ContractEventPollerConfig,
	et WorkflowRegistryEventType,
	batchCh chan<- []WorkflowRegistryEventResponse,
) {
	// create query
	var (
		responseBatch []WorkflowRegistryEventResponse
		logData       values.Value
		cursor        = ""
		limitAndSort  = query.LimitAndSort{
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
						query.Block(strconv.FormatUint(cfg.StartBlockNum, 10), primitives.Gte),
					},
				},
				limitAndSort,
				&logData,
			)

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

func newReader(
	ctx context.Context,
	lggr logger.Logger,
	factory ContractReaderFactory,
	bc types.BoundContract,
) (types.ContractReader, error) {
	contractReaderCfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			ContractName: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{string(ForceUpdateSecretsEvent)},
				},
				ContractABI: workflow_registry_wrapper.WorkflowRegistryABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					string(ForceUpdateSecretsEvent): {
						ChainSpecificName: string(ForceUpdateSecretsEvent),
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

	reader, err := factory.NewContractReader(ctx, marshalledCfg)
	if err != nil {
		return nil, err
	}

	// bind contract to contract reader
	if err := reader.Bind(ctx, []types.BoundContract{bc}); err != nil {
		return nil, err
	}

	return reader, nil
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
	default:
		lggr.Errorf("unknown event type: %s", evt)
		resp.Event = nil
		resp.Err = fmt.Errorf("unknown event type: %s", evt)
	}

	return resp
}

type nullWorkflowRegistrySyncer struct {
	services.Service
}

func NewNullWorkflowRegistrySyncer() *nullWorkflowRegistrySyncer {
	return &nullWorkflowRegistrySyncer{}
}

// Start
func (u *nullWorkflowRegistrySyncer) Start(context.Context) error {
	return nil
}

// Close
func (u *nullWorkflowRegistrySyncer) Close() error {
	return nil
}

// SecretsFor
func (u *nullWorkflowRegistrySyncer) SecretsFor(context.Context, string, string) (map[string]string, error) {
	return nil, nil
}

func (u *nullWorkflowRegistrySyncer) Ready() error {
	return nil
}

func (u *nullWorkflowRegistrySyncer) HealthReport() map[string]error {
	return nil
}

func (u *nullWorkflowRegistrySyncer) Name() string {
	return "Null" + name
}
