package logevent

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// Log Event Trigger Capability Request Config Details
type RequestConfig struct {
	ContractName      string `json:"contractName"`
	ContractAddress   string `json:"contractAddress"`
	ContractEventName string `json:"contractEventName"`
	// Log Event Trigger capability takes in a []byte as ContractReaderConfig
	// to not depend on evm ChainReaderConfig type and be chain agnostic
	ContractReaderConfig map[string]any `json:"contractReaderConfig"`
}

// LogEventTrigger struct to listen for Contract events using ContractReader gRPC client
// in a loop with a periodic delay of pollPeriod milliseconds, which is specified in
// the job spec
type logEventTrigger struct {
	ch   chan<- capabilities.TriggerResponse
	lggr logger.Logger

	// Contract address and Event Signature to monitor for
	reqConfig      *RequestConfig
	contractReader types.ContractReader
	relayer        core.Relayer
	startBlockNum  uint64

	// Log Event Trigger config with pollPeriod and lookbackBlocks
	logEventConfig Config
	ticker         *time.Ticker
	done           chan bool
	wg             sync.WaitGroup
}

// Construct for logEventTrigger struct
func newLogEventTrigger(ctx context.Context,
	reqConfig *RequestConfig,
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
			fmt.Errorf("error fetching contractReader for chainID %d from relayerSet: %v", logEventConfig.ChainID, err)
	}

	// Bind Contract in ContractReader
	boundContracts := []types.BoundContract{{Name: reqConfig.ContractName, Address: reqConfig.ContractAddress}}
	err = contractReader.Bind(ctx, boundContracts)
	if err != nil {
		return nil, nil, err
	}

	// Get current block HEAD/tip of the blockchain to start polling from
	latestHead, err := relayer.LatestHead(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting latestHead from relayer client: %v", err)
	}
	height, err := strconv.ParseUint(latestHead.Height, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid height in latestHead from relayer client: %v", err)
	}
	startBlockNum := uint64(0)
	if height > logEventConfig.LookbackBlocks {
		startBlockNum = height - logEventConfig.LookbackBlocks
	}

	// Setup callback channel, logger and ticker to poll ContractReader
	callbackCh := make(chan capabilities.TriggerResponse)
	ticker := time.NewTicker(time.Duration(int64(logEventConfig.PollPeriod)) * time.Millisecond)
	done := make(chan bool)
	lggr, err := logger.New()
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialise logger for LogEventTrigger")
	}

	// Initialise a Log Event Trigger
	l := &logEventTrigger{
		ch:   callbackCh,
		lggr: logger.Named(lggr, "LogEventTrigger: "),

		reqConfig:      reqConfig,
		contractReader: contractReader,
		relayer:        relayer,
		startBlockNum:  startBlockNum,

		logEventConfig: logEventConfig,
		ticker:         ticker,
		done:           done,
	}
	go l.Listen(ctx)

	return l, callbackCh, nil
}

// Listen to contract events and trigger workflow runs
func (l *logEventTrigger) Listen(ctx context.Context) {
	// Listen for events from lookbackPeriod
	var logs []types.Sequence
	var err error
	logData := make(map[string]any)
	cursor := ""
	limitAndSort := query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
	}
	for {
		select {
		case <-l.done:
			return
		case t := <-l.ticker.C:
			l.lggr.Infow("Polling event logs from ContractReader using QueryKey at", "time", t,
				"startBlockNum", l.startBlockNum,
				"cursor", cursor)
			if cursor != "" {
				limitAndSort.Limit = query.Limit{Cursor: cursor}
			}
			logs, err = l.contractReader.QueryKey(
				ctx,
				types.BoundContract{Name: l.reqConfig.ContractName, Address: l.reqConfig.ContractAddress},
				query.KeyFilter{
					Key: l.reqConfig.ContractEventName,
					Expressions: []query.Expression{
						query.Confidence(primitives.Unconfirmed),
						query.Block(fmt.Sprintf("%d", l.startBlockNum), primitives.Gte),
					},
				},
				limitAndSort,
				&logData,
			)
			if err != nil {
				l.lggr.Fatalw("QueryKey failure", "err", err)
				continue
			}
			if len(logs) == 1 && logs[0].Cursor == cursor {
				l.lggr.Infow("No new logs since", "cursor", cursor)
				continue
			}
			for _, log := range logs {
				triggerResp := createTriggerResponse(log, l.logEventConfig.Version(ID))
				l.wg.Add(1)
				go func(resp capabilities.TriggerResponse) {
					defer l.wg.Done()
					l.ch <- resp
				}(triggerResp)
				cursor = log.Cursor
			}
		}
	}
}

// Create log event trigger capability response
func createTriggerResponse(log types.Sequence, version string) capabilities.TriggerResponse {
	wrappedPayload, err := values.WrapMap(log)
	if err != nil {
		return capabilities.TriggerResponse{
			Err: fmt.Errorf("error wrapping trigger event: %s", err),
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

// Stop contract event listener for the current contract
func (l *logEventTrigger) Stop() {
	l.lggr.Infow("Closing trigger server for (waiting for waitGroup)", "ChainID", l.logEventConfig.ChainID,
		"ContractName", l.reqConfig.ContractName,
		"ContractAddress", l.reqConfig.ContractAddress,
		"ContractEventName", l.reqConfig.ContractEventName)
	l.wg.Wait()
	close(l.ch)
	l.done <- true
	l.lggr.Infow("Closed trigger server for", "ChainID", l.logEventConfig.ChainID,
		"ContractName", l.reqConfig.ContractName,
		"ContractAddress", l.reqConfig.ContractAddress,
		"ContractEventName", l.reqConfig.ContractEventName)
}
