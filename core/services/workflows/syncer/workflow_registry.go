package syncer

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

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

type WorkflowRegistryEventType string

var (
	ForceUpdateSecretsEvent WorkflowRegistryEventType = "WorkflowForceUpdateSecretsRequestedV1"
)

type WorkflowRegistryEvent struct {
	logeventcap.Output
	EventType WorkflowRegistryEventType
}
type WorkflowRegistryEventResponse struct {
	Err   error
	Event WorkflowRegistryEvent
}

type Syncer interface {
	Sync(ctx context.Context, isInitialSync bool) error
}

type URLGetter interface {
	GetURLHash() string
}

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

type ContractEventPollerConfig struct {
	ContractName      string
	ContractAddress   string
	ContractEventName string
	StartBlockNum     uint64
	QueryCount        uint64
}

type FetcherFunc func(ctx context.Context, url string) ([]byte, error)

type WorkflowRegistrySyncer interface {
	services.Service
}

var _ WorkflowRegistrySyncer = (*workflowRegistry)(nil)

type workflowRegistry struct {
	services.StateMachine
	stopCh  services.StopChan
	lggr    logger.Logger
	orm     WorkflowRegistryDS
	reader  ContractReader
	gateway FetcherFunc
	wg      sync.WaitGroup
	ticker  <-chan struct{}
	cfg     ContractEventPollerConfig
}

func WithTicker(ticker <-chan struct{}) func(*workflowRegistry) {
	return func(wr *workflowRegistry) {
		wr.ticker = ticker
	}
}

func NewWorkflowRegistry(
	lggr logger.Logger,
	orm WorkflowRegistryDS,
	reader ContractReader,
	gateway FetcherFunc,
	cfg ContractEventPollerConfig,
	opts ...func(*workflowRegistry),
) *workflowRegistry {
	wr := &workflowRegistry{
		lggr:    lggr.Named(name),
		orm:     orm,
		reader:  reader,
		gateway: gateway,
		cfg:     cfg,
		stopCh:  make(services.StopChan),
	}
	for _, opt := range opts {
		opt(wr)
	}
	return wr
}

func (w *workflowRegistry) Start(_ context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		var (
			ctx, _       = w.stopCh.NewCtx()
			eventsCh     = make(chan WorkflowRegistryEventResponse)
			ticker       = w.getTicker(ctx.Done())
			eventHandler = newEventHandler(w.lggr, w.orm, w.gateway)
		)

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			queryLoop(ctx, w.cfg, w.lggr, w.reader, ticker, eventsCh)
		}()

		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			handlerLoop(ctx, w.lggr, eventHandler, eventsCh)
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

func (w *workflowRegistry) getTicker(stop <-chan struct{}) <-chan struct{} {
	if w.ticker == nil {
		return signalers.MakeTicker(stop, defaultTickInterval)
	}

	return w.ticker
}

func handlerLoop(
	ctx context.Context,
	lggr logger.Logger,
	h Handler,
	events <-chan WorkflowRegistryEventResponse,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, open := <-events:
			if !open {
				return
			}

			if resp.Err != nil {
				continue
			}

			event := resp.Event
			if err := h.Handle(ctx, event); err != nil {
				lggr.Errorf("failed to handle event: %+v", event)
				continue
			}
		}
	}
}

func queryLoop(
	ctx context.Context,
	cfg ContractEventPollerConfig,
	lggr logger.Logger,
	cr ContractReader,
	ticker <-chan struct{},
	eventsCh chan<- WorkflowRegistryEventResponse,
) {
	// setup helpers
	var (
		sendLog = func(ld WorkflowRegistryEventResponse) {
			select {
			case eventsCh <- ld:
			case <-ctx.Done():
			}
		}

		sendErr = func(err error) {
			select {
			case eventsCh <- WorkflowRegistryEventResponse{
				Err: err,
			}:
			case <-ctx.Done():
			default:
			}
		}
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
	err = cr.Bind(ctx, boundContracts)
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

			logs, err = cr.QueryKey(
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
				lggr.Errorw("QueryKey failure", "err", err)
				sendErr(err)
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

				sendLog(toWorkflowRegistryEventResponse(log, lggr))

				cursor = log.Cursor
			}
		}
	}
}

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

	lggr.Debugf("got data before decode %+v", dataAsMap)

	var event any
	if err := mapstructure.Decode(dataAsMap, &event); err != nil {
		lggr.Debugf("got error from decoding %s", err)
		return WorkflowRegistryEventResponse{
			Err: err,
		}
	}

	return WorkflowRegistryEventResponse{
		Event: WorkflowRegistryEvent{
			// TODO(mstreet3): This type should be dynamically known from the event if possible
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
