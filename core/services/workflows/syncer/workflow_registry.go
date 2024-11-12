package syncer

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	types "github.com/smartcontractkit/chainlink-common/pkg/types"
	query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils/signalers"
)

const name = "WorkflowRegistrySyncer"

var (
	defaultTickInterval = 12 * time.Second
	ContractName        = "WorkflowRegistry"
	ContractEventName   = "WorkflowForceUpdateSecretsRequestedV1"
)

// WorkflowRegistryrEventType is the type of event that is emitted by the WorkflowRegistry
type WorkflowRegistryEventType string

var (
	// ForceUpdateSecretsEvent is emitted when a request to force update a workflows secrets is made
	ForceUpdateSecretsEvent WorkflowRegistryEventType = "WorkflowForceUpdateSecretsRequestedV1"
)

// WorkflowRegistryEvent is an event emitted by the WorkflowRegistry.  Each event is typed
// so that the consumer can determine how to handle the event.
type WorkflowRegistryEvent struct {
	logeventcap.Output
	EventType WorkflowRegistryEventType
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
// TODO(mstreet3): Can we query all contract events at once and tag each with an event type?
//
// TODO(mstreet3): Use LookbackBlocks instead of StartBlockNum
type ContractEventPollerConfig struct {
	ContractName      string
	ContractAddress   string
	ContractEventName string
	StartBlockNum     uint64
	QueryCount        uint64
}

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

// WorkflowRegistrySyncer is the public interface of the package.
type WorkflowRegistrySyncer interface {
	services.Service
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

// workflowRegistry is the implementation of the WorkflowRegistrySyncer interface.
type workflowRegistry struct {
	services.StateMachine
	stopCh   services.StopChan
	lggr     logger.Logger
	orm      WorkflowRegistryDS
	reader   ContractReader
	gateway  FetcherFunc
	wg       sync.WaitGroup
	ticker   <-chan struct{}
	eventsCh chan WorkflowRegistryEventResponse
	handler  handler
	cfg      ContractEventPollerConfig
}

// WithTicker allows external callers to provide a ticker to the workflowRegistry.  This is useful
// for overriding the default tick interval.
func WithTicker(ticker <-chan struct{}) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
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
			ContractName:      ContractName,
			ContractAddress:   addr,
			ContractEventName: ContractEventName,
			QueryCount:        20,
			StartBlockNum:     0,
		},
		stopCh:   make(services.StopChan),
		eventsCh: make(chan WorkflowRegistryEventResponse),
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

		ticker = w.getTicker(ctx.Done())
	)

	// create query
	var (
		logs    []types.Sequence
		err     error
		logData values.Value

		boundContracts = []types.BoundContract{
			{
				Name:    cfg.ContractName,
				Address: cfg.ContractAddress,
			},
		}
		cursor       = ""
		limitAndSort = query.LimitAndSort{
			SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
			Limit:  query.Limit{Count: cfg.QueryCount},
		}
	)

	// bind contracts to contract reader
	err = w.reader.Bind(ctx, boundContracts)
	if err != nil {
		sendErr(err)
		return
	}

	// Loop until canceled
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			if cursor != "" {
				limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, cfg.QueryCount)
			}

			logs, err = w.reader.QueryKey(
				ctx,
				types.BoundContract{
					Name:    cfg.ContractName,
					Address: cfg.ContractAddress,
				},
				query.KeyFilter{
					Key: cfg.ContractEventName,
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

				sendLog(toWorkflowRegistryEventResponse(log, w.lggr))

				cursor = log.Cursor
			}
		}
	}
}

// getTicker returns the ticker that the workflowRegistry will use to poll for events.  If the ticker
// is nil, then a default ticker is returned.
func (w *workflowRegistry) getTicker(stop <-chan struct{}) <-chan struct{} {
	if w.ticker == nil {
		return signalers.MakeTicker(stop, defaultTickInterval)
	}

	return w.ticker
}

// toWorkflowRegistryEventResponse converts a types.Sequence to a WorkflowRegistryEventResponse.
//
// TODO(mstreet3): The event type is hardcoded to ForceUpdateSecretsEvent.  This should be a parameter
// or function of the incoming log.
func toWorkflowRegistryEventResponse(
	log types.Sequence,
	lggr logger.Logger,
) WorkflowRegistryEventResponse {
	dataAsValuesMap, err := values.WrapMap(log.Data)
	if err != nil {
		lggr.Debugf("failed to wrap data as map : %+v", err)
		return WorkflowRegistryEventResponse{
			Err: err,
		}
	}

	var dataAsMap map[string]any
	err = dataAsValuesMap.UnwrapTo(&dataAsMap)
	if err != nil {
		lggr.Debugf("failed to unwrap to map[string]any : %+v", err)
		return WorkflowRegistryEventResponse{
			Err: err,
		}
	}

	return WorkflowRegistryEventResponse{
		Event: &WorkflowRegistryEvent{
			EventType: ForceUpdateSecretsEvent,
			Output: logeventcap.Output{
				Cursor: log.Cursor,
				Data:   dataAsMap,
				Head: logeventcap.Head{
					Hash:      fmt.Sprintf("0x%x", log.Hash),
					Height:    log.Height,
					Timestamp: log.Timestamp,
				},
			},
		},
	}
}
