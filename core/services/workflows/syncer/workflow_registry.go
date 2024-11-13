package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
)

// WorkflowRegistryrEventType is the type of event that is emitted by the WorkflowRegistry
type WorkflowRegistryEventType string

var (
	// ForceUpdateSecretsEvent is emitted when a request to force update a workflows secrets is made
	ForceUpdateSecretsEvent WorkflowRegistryEventType = "WorkflowForceUpdateSecretsRequestedV1"
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

// ContractReader is the subset of methods needed to query the events of a contract.
type ContractReader interface {
	Bind(context.Context, []types.BoundContract) error
	QueryKey(
		context.Context,
		types.BoundContract,
		query.KeyFilter,
		query.LimitAndSort,
		any,
	) ([]types.Sequence, error)
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

// WorkflowRegistrySyncer is the public interface of the package.
type WorkflowRegistrySyncer interface {
	services.Service
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

// workflowRegistry is the implementation of the WorkflowRegistrySyncer interface.
type workflowRegistry struct {
	services.StateMachine
	stopCh     services.StopChan
	lggr       logger.Logger
	orm        WorkflowRegistryDS
	reader     ContractReader
	initReader func(context.Context, logger.Logger, ContractReaderFactory, types.BoundContract) (ContractReader, error)
	relayer    ContractReaderFactory
	gateway    FetcherFunc
	wg         sync.WaitGroup
	ticker     <-chan struct{}
	eventsCh   chan WorkflowRegistryEventResponse
	handler    handler
	cfg        ContractEventPollerConfig
}

// WithTicker allows external callers to provide a ticker to the workflowRegistry.  This is useful
// for overriding the default tick interval.
func WithTicker(ticker <-chan struct{}) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
	}
}

func WithReader(reader ContractReader) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.reader = reader
	}
}

// NewWorkflowRegistry returns a new workflowRegistry.
// Only queries for WorkflowRegistryForceUpdateSecretsRequestedV1 events.
func NewWorkflowRegistry(
	lggr logger.Logger,
	orm WorkflowRegistryDS,
	reader ContractReader,
	gateway FetcherFunc,
	addr string,
	opts ...func(*workflowRegistry),
) *workflowRegistry {
	wr := &workflowRegistry{
		lggr:    lggr.Named(name),
		orm:     orm,
		reader:  reader,
		gateway: gateway,
		cfg: ContractEventPollerConfig{
			ContractName:    ContractName,
			ContractAddress: addr,
			QueryCount:      20,
			StartBlockNum:   0,
		},
		initReader: newReader,
		stopCh:     make(services.StopChan),
		eventsCh:   make(chan WorkflowRegistryEventResponse),
	}
	wr.handler = newEventHandler(wr.lggr, wr.orm, wr.gateway)
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

			w.queryLoop(ctx, w.cfg)
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
	return w.orm.SecretsFor(ctx, workflowOwner, workflowName)
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
			if err := w.handler.Handle(ctx, *event); err != nil {
				w.lggr.Errorf("failed to handle event: %+v", event)
				continue
			}
		}
	}
}

// queryLoop polls the contract for events.
func (w *workflowRegistry) queryLoop(
	ctx context.Context,
	cfg ContractEventPollerConfig,
) {
	// setup helpers
	var (
		sendLog = func(resp WorkflowRegistryEventResponse) {
			select {
			case w.eventsCh <- resp:
			case <-ctx.Done():
			}
		}

		sendErr = func(err error) {
			select {
			case w.eventsCh <- WorkflowRegistryEventResponse{
				Err: err,
			}:
			case <-ctx.Done():
			default:
			}
		}

		ticker = w.getTicker(ctx)
	)

	// setup contract reader
	boundContract := types.BoundContract{
		Name:    cfg.ContractName,
		Address: cfg.ContractAddress,
	}
	reader, err := w.getContractReader(ctx, boundContract)
	if err != nil {
		sendErr(err)
		return
	}

	// bind contract to contract reader
	if err := reader.Bind(ctx, []types.BoundContract{boundContract}); err != nil {
		sendErr(err)
		return
	}

	/* // setup log handling
	events := []WorkflowRegistryEventType{ForceUpdateSecretsEvent}
	logCh := make(chan []types.Sequence, len(events)) */

	// create query
	var (
		logs         []types.Sequence
		logData      WorkflowRegistryForceUpdateSecretsRequestedV1
		cursor       = ""
		limitAndSort = query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: cfg.QueryCount},
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

			logs, err = reader.QueryKey(
				ctx,
				boundContract,
				query.KeyFilter{
					Key: string(ForceUpdateSecretsEvent),
					Expressions: []query.Expression{
						query.Confidence(primitives.Finalized),
						query.Block(strconv.FormatUint(cfg.StartBlockNum, 10), primitives.Gte),
					},
				},
				limitAndSort,
				&logData,
			)

			if err != nil {
				w.lggr.Errorw("QueryKey failure", "err", err)
				sendErr(err)
				continue
			}

			// ChainReader QueryKey API provides logs including the cursor value and not
			// after the cursor value. If the response only consists of the log corresponding
			// to the cursor and no log after it, then we understand that there are no new
			// logs
			if len(logs) == 1 && logs[0].Cursor == cursor {
				w.lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}

			for _, log := range logs {
				if log.Cursor == cursor {
					continue
				}

				sendLog(toWorkflowRegistryEventResponse(log, ForceUpdateSecretsEvent, w.lggr))

				cursor = log.Cursor
			}
		}
	}
}

func queryEvent[T any](
	ctx context.Context,
	lggr logger.Logger,
	reader ContractReader,
	ticker <-chan struct{},
	et WorkflowRegistryEventType,
	bc types.BoundContract,
	cfg ContractEventPollerConfig,
	logsCh chan<- []types.Sequence,
) {
	// create query
	var (
		logsToSend   []types.Sequence
		logData      WorkflowRegistryForceUpdateSecretsRequestedV1
		cursor       = ""
		limitAndSort = query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: cfg.QueryCount},
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
					Key: string(ForceUpdateSecretsEvent),
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
				logsToSend = append(logsToSend, log)
				cursor = log.Cursor
			}
			logsCh <- logsToSend
			logsToSend = make([]types.Sequence, 0)
		}
	}
}

// getTicker returns the ticker that the workflowRegistry will use to poll for events.  If the ticker
// is nil, then a default ticker is returned.
func (w *workflowRegistry) getTicker(ctx context.Context) <-chan struct{} {
	if w.ticker == nil {
		return signalers.MakeTicker(ctx.Done(), defaultTickInterval)
	}

	return w.ticker
}

func (w *workflowRegistry) getContractReader(ctx context.Context, c types.BoundContract) (ContractReader, error) {
	if w.reader == nil {
		reader, err := w.initReader(ctx, w.lggr, w.relayer, c)
		if err != nil {
			return nil, err
		}
		w.reader = reader
	}
	return w.reader, nil
}

func newReader(
	ctx context.Context,
	lggr logger.Logger,
	factory ContractReaderFactory,
	contract types.BoundContract,
) (ContractReader, error) {
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
						OutputModifications: commoncodec.ModifiersConfig{
							&commoncodec.AddressBytesToStringModifierConfig{
								Fields: WorkflowRegistryForceUpdateSecretsRequestedV1ModifyFields,
							},
						},
					},
				},
			},
		},
	}

	marshalledCfg, err := json.Marshal(contractReaderCfg)
	if err != nil {
		return nil, err
	}

	return factory.NewContractReader(ctx, marshalledCfg)
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
				Hash:      fmt.Sprintf("%x", log.Hash),
				Height:    log.Height,
				Timestamp: log.Timestamp,
			},
		},
	}

	switch evt {
	case ForceUpdateSecretsEvent:
		var data WorkflowRegistryForceUpdateSecretsRequestedV1
		if err := mapstructure.Decode(log.Data, &data); err != nil {
			lggr.Errorf("failed to decode data: %+v", log.Data)
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
