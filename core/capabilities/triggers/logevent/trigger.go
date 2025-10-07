package logevent

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/triggers/logevent/logeventcap"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

// LogEventTrigger struct to listen for Contract events using ContractReader gRPC client
// in a loop with a periodic delay of pollPeriod milliseconds, which is specified in
// the job spec
type logEventTrigger struct {
	ch   chan<- capabilities.TriggerResponse
	lggr logger.Logger

	// Workflow metadata
	metadata capabilities.RequestMetadata

	// Contract address and Event Signature to monitor for
	reqConfig      *logeventcap.Config
	contractReader types.ContractReader
	relayer        core.Relayer
	startBlockNum  uint64

	// Log Event Trigger config with pollPeriod and lookbackBlocks
	logEventConfig Config
	ticker         *time.Ticker
	stopChan       services.StopChan
	done           chan bool
}

// Construct for logEventTrigger struct
func newLogEventTrigger(ctx context.Context,
	lggr logger.Logger,
	metadata capabilities.RequestMetadata,
	reqConfig *logeventcap.Config,
	logEventConfig Config,
	relayer core.Relayer) (*logEventTrigger, chan capabilities.TriggerResponse, error) {
	jsonBytes, err := json.Marshal(reqConfig.ContractReaderConfig)
	if err != nil {
		return nil, nil, err
	}

	// Create a New Contract Reader client, which brings a corresponding ContractReader gRPC service
	// in Chainlink Core service
	contractReader, err := relayer.NewContractReader(ctx, jsonBytes)
	if err != nil {
		return nil, nil,
			fmt.Errorf("error fetching contractReader for chainID %s from relayerSet: %w", logEventConfig.ChainID, err)
	}

	// Bind Contract in ContractReader
	boundContracts := []types.BoundContract{{Name: reqConfig.ContractName, Address: reqConfig.ContractAddress}}
	err = contractReader.Bind(ctx, boundContracts)
	if err != nil {
		return nil, nil, err
	}

	err = contractReader.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get current block HEAD/tip of the blockchain to start polling from
	latestHead, err := relayer.LatestHead(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting latestHead from relayer client: %w", err)
	}
	height, err := strconv.ParseUint(latestHead.Height, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid height in latestHead from relayer client: %w", err)
	}
	startBlockNum := uint64(0)
	if height > logEventConfig.LookbackBlocks {
		startBlockNum = height - logEventConfig.LookbackBlocks
	}

	// Setup callback channel, logger and ticker to poll ContractReader
	callbackCh := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	ticker := time.NewTicker(time.Duration(logEventConfig.PollPeriod) * time.Millisecond)

	if logEventConfig.QueryCount == 0 {
		logEventConfig.QueryCount = 20
	}

	// Initialise a Log Event Trigger
	l := &logEventTrigger{
		ch:   callbackCh,
		lggr: logger.Named(lggr, "LogEventTrigger."+metadata.WorkflowID),

		metadata:       metadata,
		reqConfig:      reqConfig,
		contractReader: contractReader,
		relayer:        relayer,
		startBlockNum:  startBlockNum,

		logEventConfig: logEventConfig,
		ticker:         ticker,
		stopChan:       make(services.StopChan),
		done:           make(chan bool),
	}
	return l, callbackCh, nil
}

func (l *logEventTrigger) Start(ctx context.Context) error {
	go l.listen()
	return nil
}

// Start to contract events and trigger workflow runs
func (l *logEventTrigger) listen() {
	ctx, cancel := l.stopChan.NewCtx()
	defer cancel()
	defer close(l.done)

	// Listen for events from lookbackPeriod
	var logs []types.Sequence
	var err error
	var logData values.Value
	cursor := ""
	limitAndSort := query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
		Limit:  query.Limit{Count: l.logEventConfig.QueryCount},
	}
	for {
		select {
		case <-ctx.Done():
			l.lggr.Infow("Closing trigger server for (waiting for waitGroup)", "ChainID", l.logEventConfig.ChainID,
				"ContractName", l.reqConfig.ContractName,
				"ContractAddress", l.reqConfig.ContractAddress,
				"ContractEventName", l.reqConfig.ContractEventName)
			return
		case t := <-l.ticker.C:
			l.lggr.Infow("Polling event logs from ContractReader using QueryKey at", "time", t,
				"startBlockNum", l.startBlockNum,
				"cursor", cursor)
			if cursor != "" {
				limitAndSort.Limit = query.CursorLimit(cursor, query.CursorFollowing, l.logEventConfig.QueryCount)
			}
			logs, err = l.contractReader.QueryKey(
				ctx,
				types.BoundContract{Name: l.reqConfig.ContractName, Address: l.reqConfig.ContractAddress},
				query.KeyFilter{
					Key: l.reqConfig.ContractEventName,
					Expressions: []query.Expression{
						query.Confidence(primitives.Finalized),
						query.Block(strconv.FormatUint(l.startBlockNum, 10), primitives.Gte),
					},
				},
				limitAndSort,
				&logData,
			)
			if err != nil {
				l.lggr.Errorw("QueryKey failure", "err", err)
				continue
			}
			// ChainReader QueryKey API provides logs including the cursor value and not
			// after the cursor value. If the response only consists of the log corresponding
			// to the cursor and no log after it, then we understand that there are no new
			// logs
			if len(logs) == 1 && logs[0].Cursor == cursor {
				l.lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}
			for _, log := range logs {
				if log.Cursor == cursor {
					continue
				}
				triggerResp := createTriggerResponse(log, l.logEventConfig.Version(ID))

				// Emit trigger execution started event
				workflowExecutionID, err := events.GenerateExecutionID(l.metadata.WorkflowID, triggerResp.Event.ID)
				if err != nil {
					l.lggr.Errorw("failed to generate execution ID", "err", err)
					workflowExecutionID = ""
				}

				// Create labels map with workflow metadata
				labels := map[string]string{
					platform.KeyWorkflowID:    l.metadata.WorkflowID,
					platform.KeyWorkflowOwner: l.metadata.WorkflowOwner,
					platform.KeyWorkflowName:  l.metadata.WorkflowName,
				}

				// Add optional metadata fields if available
				if l.metadata.WorkflowTag != "" {
					labels[platform.KeyWorkflowTag] = l.metadata.WorkflowTag
				}
				if l.metadata.WorkflowDonID != 0 {
					labels[platform.KeyDonID] = strconv.FormatUint(uint64(l.metadata.WorkflowDonID), 10)
				}
				if l.metadata.WorkflowDonConfigVersion != 0 {
					labels[platform.DonVersion] = strconv.FormatUint(uint64(l.metadata.WorkflowDonConfigVersion), 10)
				}

				err = events.EmitTriggerExecutionStarted(ctx, labels, triggerResp.Event.ID, workflowExecutionID)
				if err != nil {
					l.lggr.Errorw("failed to emit trigger execution started event", "err", err)
				}

				l.ch <- triggerResp
				cursor = log.Cursor
			}
		}
	}
}

// Create log event trigger capability response
func createTriggerResponse(log types.Sequence, version string) capabilities.TriggerResponse {
	dataAsValuesMap, err := values.WrapMap(log.Data)
	if err != nil {
		return capabilities.TriggerResponse{
			Err: fmt.Errorf("error decoding log data as values.Map: %w", err),
		}
	}
	dataAsMap := map[string]any{}
	err = dataAsValuesMap.UnwrapTo(&dataAsMap)
	if err != nil {
		return capabilities.TriggerResponse{
			Err: fmt.Errorf("error decoding log data as map[string]any: %w", err),
		}
	}

	wrappedPayload, err := values.WrapMap(&logeventcap.Output{
		Cursor: log.Cursor,
		Data:   dataAsMap,
		Head: logeventcap.Head{
			Hash:      "0x" + hex.EncodeToString(log.Hash),
			Height:    log.Height,
			Timestamp: log.Timestamp,
		},
	})
	if err != nil {
		return capabilities.TriggerResponse{
			Err: fmt.Errorf("error wrapping trigger event: %w", err),
		}
	}

	return capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: version,
			ID:          log.Cursor,
			Outputs:     wrappedPayload,
		},
	}
}

// Close contract event listener for the current contract
// This function is called when UnregisterTrigger is called individually
// for a specific ContractAddress and EventName
// When the whole capability service is stopped, stopChan of the service
// is closed, which would stop all triggers
func (l *logEventTrigger) Close() error {
	close(l.stopChan)
	<-l.done
	return nil
}
