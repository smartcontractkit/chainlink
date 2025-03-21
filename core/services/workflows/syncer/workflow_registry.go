package syncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
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
	DonID     *uint32
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

	// all goroutines are waited on with wg.
	wg sync.WaitGroup

	// ticker is the interval at which the workflowRegistry will poll the contract for events.
	ticker <-chan time.Time

	lggr                    logger.Logger
	workflowRegistryAddress string

	newContractReaderFn newContractReaderFn

	eventPollerCfg WorkflowEventPollerConfig
	eventTypes     []WorkflowRegistryEventType
	handler        evtHandler

	workflowDonNotifier donNotifier
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
		lggr:                    lggr,
		newContractReaderFn:     newContractReaderFn,
		workflowRegistryAddress: addr,
		eventPollerCfg:          eventPollerConfig,
		stopCh:                  make(services.StopChan),
		eventTypes:              ets,
		handler:                 handler,
		workflowDonNotifier:     workflowDonNotifier,
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

			reader, err := w.newWorkflowRegistryContractReader(ctx)
			if err != nil {
				w.lggr.Criticalf("contract reader unavailable : %s", err)
				return
			}

			w.lggr.Debugw("Loading initial workflows for DON", "DON", don.ID)
			loadWorkflowsHead, err := w.loadWorkflows(ctx, don, reader)
			if err != nil {
				// TODO - this is a temporary fix to handle the case where the chainreader errors because the contract
				// contains no workflows.  To track: https://smartcontract-it.atlassian.net/browse/CAPPL-393
				if !strings.Contains(err.Error(), "attempting to unmarshal an empty string while arguments are expected") {
					w.lggr.Errorw("failed to load workflows", "err", err)
					return
				}

				loadWorkflowsHead = &types.Head{
					Height: "0",
				}
			}

			w.readRegistryEvents(ctx, don, reader, loadWorkflowsHead.Height)
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

// readRegistryEvents polls the contract for events and send them to the events channel.
func (w *workflowRegistry) readRegistryEvents(ctx context.Context, don capabilities.DON, reader ContractReader, lastReadBlockNumber string) {
	ticker := w.getTicker()

	var keyQueries = make([]types.ContractKeyFilter, 0, len(w.eventTypes))
	for _, et := range w.eventTypes {
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
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			limitAndSort := query.LimitAndSort{
				SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
				Limit:  query.Limit{Count: w.eventPollerCfg.QueryCount},
			}
			if cursor != "" {
				limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, w.eventPollerCfg.QueryCount)
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
			w.lggr.Debugw("QueryKeys called", "logs", len(logs), "eventTypes", w.eventTypes, "lastReadBlockNumber", lastReadBlockNumber, "logCursor", cursor)

			// ChainReader QueryKey API provides logs including the cursor value and not
			// after the cursor value. If the response only consists of the log corresponding
			// to the cursor and no log after it, then we understand that there are no new
			// logs
			if len(logs) == 1 && logs[0].Sequence.Cursor == cursor {
				w.lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}

			var events []WorkflowRegistryEventResponse
			for _, log := range logs {
				if log.Sequence.Cursor == cursor {
					continue
				}

				event := toWorkflowRegistryEventResponse(log.Sequence, log.EventType, w.lggr)

				switch {
				case event.Event.DonID == nil:
					// event is missing a DonID, so don't filter it out;
					// it applies to all Dons
					events = append(events, event)
				case *event.Event.DonID == don.ID:
					// event has a DonID and matches, so it applies to this DON.
					events = append(events, event)
				default:
					// event doesn't match, let's skip it
					donID := "MISSING_DON_ID"
					if event.Event.DonID != nil {
						donID = strconv.FormatUint(uint64(*event.Event.DonID), 10)
					}
					w.lggr.Debugw("event belongs to a different don, skipping...", "don", don.ID, "gotDON", donID)
				}

				cursor = log.Sequence.Cursor
			}

			for _, event := range events {
				err := w.handler.Handle(ctx, event.Event)
				if err != nil {
					w.lggr.Errorw("failed to handle event", "err", err, "type", event.Event.EventType)
				}
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

func (w *workflowRegistry) loadWorkflows(ctx context.Context, don capabilities.DON, contractReader ContractReader) (*types.Head, error) {
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
	for {
		var err error
		var workflows GetWorkflowMetadataListByDONReturnVal
		headAtLastRead, err = contractReader.GetLatestValueWithHeadData(ctx, readIdentifier, primitives.Finalized, params, &workflows)
		if err != nil {
			return nil, fmt.Errorf("failed to get lastest value with head data %w", err)
		}

		w.lggr.Debugw("Rehydrating existing workflows", "len", len(workflows.WorkflowMetadataList))
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
			if err = w.handler.Handle(ctx, workflowAsEvent{
				Data:      toRegisteredEvent,
				EventType: WorkflowRegisteredEvent,
			}); err != nil {
				w.lggr.Errorf("failed to handle workflow registration: %s", err)
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
		resp.Event.DonID = &data.DonID
	case WorkflowUpdatedEvent:
		var data WorkflowRegistryWorkflowUpdatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
		resp.Event.DonID = &data.DonID
	case WorkflowPausedEvent:
		var data WorkflowRegistryWorkflowPausedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
		resp.Event.DonID = &data.DonID
	case WorkflowActivatedEvent:
		var data WorkflowRegistryWorkflowActivatedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
		resp.Event.DonID = &data.DonID
	case WorkflowDeletedEvent:
		var data WorkflowRegistryWorkflowDeletedV1
		if err := dataAsValuesMap.UnwrapTo(&data); err != nil {
			lggr.Errorf("failed to unwrap data: %+v", log.Data)
			resp.Event = nil
			resp.Err = err
			return resp
		}
		resp.Event.Data = data
		resp.Event.DonID = &data.DonID
	default:
		lggr.Errorf("unknown event type: %s", evt)
		resp.Event = nil
		resp.Err = fmt.Errorf("unknown event type: %s", evt)
	}

	return resp
}
