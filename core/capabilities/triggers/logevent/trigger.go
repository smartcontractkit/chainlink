package logevent

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// Log Event Trigger Capability Request Config Details
type RequestConfig struct {
	ContractName         string                     `json:"contractName"`
	ContractAddress      common.Address             `json:"contractAddress"`
	ContractEventName    string                     `json:"contractEventName"`
	ContractReaderConfig evmtypes.ChainReaderConfig `json:"contractReaderConfig"`
}

// LogEventTrigger struct to listen for Contract events using ContractReader gRPC client
// in a loop with a periodic delay of pollPeriod milliseconds, which is specified in
// the job spec
type logEventTrigger struct {
	ch   chan<- capabilities.TriggerResponse
	lggr logger.Logger
	ctx  context.Context

	// Contract address and Event Signature to monitor for
	reqConfig      *RequestConfig
	contractReader types.ContractReader
	relayer        core.Relayer
	startBlockNum  uint64

	// Log Event Trigger config with pollPeriod and lookbackBlocks
	logEventConfig LogEventConfig
	ticker         *time.Ticker
	done           chan bool
}

// Construct for logEventTrigger struct
func newLogEventTrigger(ctx context.Context,
	reqConfig *RequestConfig,
	logEventConfig LogEventConfig,
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
			fmt.Errorf("error fetching contractReader for chainID %d from relayerSet: %v", logEventConfig.ChainId, err)
	}

	// Bind Contract in ContractReader
	boundContracts := []types.BoundContract{{Name: reqConfig.ContractName, Address: reqConfig.ContractAddress.Hex()}}
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

	// Setup callback channel, logger and ticker to poll ContractReader
	callbackCh := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	ticker := time.NewTicker(time.Duration(logEventConfig.PollPeriod) * time.Millisecond)
	done := make(chan bool)
	lggr, err := logger.New()
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialise logger for LogEventTrigger")
	}

	// Initialise a Log Event Trigger
	l := &logEventTrigger{
		ch:   callbackCh,
		lggr: logger.Named(lggr, "LogEventTrigger: "),
		ctx:  ctx,

		reqConfig:      reqConfig,
		contractReader: contractReader,
		relayer:        relayer,
		startBlockNum:  height,

		logEventConfig: logEventConfig,
		ticker:         ticker,
		done:           done,
	}
	go l.Listen()

	return l, callbackCh, nil
}

// Listen to contract events and trigger workflow runs
func (l *logEventTrigger) Listen() {
	// Listen for events from lookbackPeriod
	var logs []types.Sequence
	var err error
	logData := make(map[string]any)
	limitAndSort := query.LimitAndSort{
		SortBy: []query.SortBy{query.NewSortByTimestamp(query.Asc)},
	}
	for {
		select {
		case <-l.done:
			return
		case t := <-l.ticker.C:
			l.lggr.Infof("Polling event logs from ContractReader using QueryKey at", t)
			logs, err = l.contractReader.QueryKey(
				l.ctx,
				types.BoundContract{Name: l.reqConfig.ContractName, Address: l.reqConfig.ContractAddress.Hex()},
				query.KeyFilter{
					Key: l.reqConfig.ContractEventName,
					Expressions: []query.Expression{
						query.Confidence(primitives.Finalized),
						query.Block(fmt.Sprintf("%d", l.startBlockNum-l.logEventConfig.LookbackBlocks), primitives.Gte),
					},
				},
				limitAndSort,
				logData,
			)
			if err != nil {
				l.lggr.Fatalw("QueryKey failure", "err", err)
				continue
			}
			for _, log := range logs {
				triggerResp := createTriggerResponse(log)
				go func(resp capabilities.TriggerResponse) {
					l.ch <- resp
				}(triggerResp)
				limitAndSort.Limit = query.Limit{Cursor: log.Cursor}
			}
		}
	}
}

// Create log event trigger capability response
func createTriggerResponse(log types.Sequence) capabilities.TriggerResponse {
	wrappedPayload, err := values.WrapMap(log.Data)
	if err != nil {
		return capabilities.TriggerResponse{
			Err: fmt.Errorf("error wrapping trigger event: %s", err),
		}
	}
	return capabilities.TriggerResponse{
		Event: capabilities.TriggerEvent{
			TriggerType: ID,
			ID:          log.Cursor,
			Outputs:     wrappedPayload,
		},
	}
}

// Stop contract event listener for the current contract
func (l *logEventTrigger) Stop() {
	close(l.ch)
	l.done <- true
}
