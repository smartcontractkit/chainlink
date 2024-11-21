package syncer

import (
	"context"
	_ "embed"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

const (
	// Compute Fetch Workflow
	workflowID    = "924eef66516e5387b6e8ab8cc544685dfe50dfc837886f22beecebced5063968"
	workflowOwner = "00000000000000000000000000000000000000ab"
	workflowName  = "trueusdpor"

	// Chain Read Workflow
	workflow2ID    = "00000066516e5387b6e8ab8cc544685dfe50dfc837886f22beecebced5063968"
	workflow2Owner = "00000000000000000000000000000000000000ab"
	workflow2Name  = "ethseppor"
)

var (
	// Compute Fetch Workflow
	//go:embed config.yaml
	config []byte

	//go:embed workflow.wasm.br
	workflow []byte

	// Chain Read Workflow
	//go:embed config.yaml
	config2 []byte

	//go:embed por-read-chain.wasm.br
	workflow2 []byte
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
	wg          sync.WaitGroup
	Logger      logger.Logger
	Registry    core.CapabilitiesRegistry
	Store       store.Store
	DS          sqlutil.DataSource
	subServices []job.ServiceCtx
}

func (w *WorkflowRegistry) Start(ctx context.Context) error {
	go func() {
		timeout := time.After(5 * time.Minute)
		ticker := time.NewTicker(10 * time.Second)

		for {
			select {
			case <-timeout:
				w.Logger.Info("timed out setting up hardcoded workflow")
				return
			case <-ticker.C:
				success1 := w.trySetup(workflowID, workflowName, workflowOwner, workflow, config)
				success2 := w.trySetup(workflow2ID, workflow2Name, workflow2Owner, workflow2, config2)
				if success1 && success2 {
					return
				}
			}
		}
	}()
	return nil
}

func (w *WorkflowRegistry) trySetup(id, name, owner string, binary, config []byte) bool {
	ctx := context.Background()
	w.Logger.Info("starting hardcoded workflow...")

	// HACK: don't load the workflow if we aren't a workflow node.
	_, err := w.Registry.Get(ctx, "offchain_reporting@1.0.0")
	if err != nil {
		w.Logger.Info("not a workflow node, skipping hardcoded workflow")
		return false
	}

	jb := job.WorkflowSpec{
		Workflow:      "a string",
		Config:        "a config",
		WorkflowID:    id,
		WorkflowName:  name,
		WorkflowOwner: owner,
	}
	sql := `INSERT INTO workflow_specs (workflow, workflow_id, workflow_owner, workflow_name, created_at, updated_at, spec_type, config)
	VALUES (:workflow, :workflow_id, :workflow_owner, :workflow_name, NOW(), NOW(), :spec_type, :config)
	RETURNING id;`
	_, err = w.DS.NamedExecContext(ctx, sql, jb)
	if err != nil {
		w.Logger.Info("failed to create entry: %w", err)
	}

	moduleConfig := &host.ModuleConfig{Logger: logger.NullLogger}
	spec, err := host.GetWorkflowSpec(ctx, moduleConfig, binary, config)
	if err != nil {
		w.Logger.Errorf("failed to get workflow spec", err)
		return false
	}

	cfg := workflows.Config{
		Lggr:           w.Logger,
		Workflow:       *spec,
		WorkflowID:     id,
		WorkflowOwner:  owner,
		WorkflowName:   name,
		Registry:       w.Registry,
		Store:          w.Store,
		Config:         config,
		Binary:         binary,
		SecretsFetcher: w,
	}
	engine, err := workflows.NewEngine(ctx, cfg)
	if err != nil {
		w.Logger.Errorf("failed to create engine: %w", err)
		return false
	}
	err = engine.Start(ctx)
	if err != nil {
		w.Logger.Errorf("failed to start hardcoded workflow: %w", err)
		return false
	}
	w.subServices = []job.ServiceCtx{engine}
	return true
}

func (w *WorkflowRegistry) Close() error {
	for _, s := range w.subServices {
		err := s.Close()
		if err != nil {
			w.Logger.Errorf("could not close hardcoded engine: %w", err)
		}
	}

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
