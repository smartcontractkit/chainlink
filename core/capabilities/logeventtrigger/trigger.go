package logeventtrigger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type logEventTrigger struct {
	ch chan<- capabilities.TriggerResponse

	// Contract address and Event Signature to monitor for
	contractName         string
	contractAddress      common.Address
	contractReaderConfig evmtypes.ChainReaderConfig
	contractReader       types.ContractReader

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

	callbackCh := make(chan capabilities.TriggerResponse, defaultSendChannelBufferSize)
	ticker := time.NewTicker(time.Duration(logEventConfig.PollPeriod) * time.Millisecond)
	done := make(chan bool)

	// Initialise a Log Event Trigger
	l := &logEventTrigger{
		ch:                   callbackCh,
		contractName:         reqConfig.ContractName,
		contractAddress:      reqConfig.ContractAddress,
		contractReaderConfig: reqConfig.ContractReaderConfig,
		contractReader:       contractReader,
		logEventConfig:       logEventConfig,
		ticker:               ticker,
		done:                 done,
	}
	go l.Listen()

	return l, callbackCh, nil
}

// Listen to contract events and trigger workflow runs
func (l *logEventTrigger) Listen() {
	// Listen for events from lookbackPeriod
	for {
		select {
		case <-l.done:
			return
		case t := <-l.ticker.C:
			fmt.Println("Tick at", t)
		}
	}
}

// Stop contract event listener for the current contract
func (l *logEventTrigger) Stop() {
	l.done <- true
}
