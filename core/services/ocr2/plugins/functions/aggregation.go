package functions

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
)

func CanAggregate(N int, F int, observations []*encoding.ProcessedRequest) bool {
	return N > 0 && F >= 0 && len(observations) > 0 && len(observations) <= N && len(observations) >= 2*F+1
}

func Aggregate(aggMethod config.AggregationMethod, observations []*encoding.ProcessedRequest) (*encoding.ProcessedRequest, error) {
	if len(observations) == 0 {
		return nil, errors.New("empty observation list passed for aggregation")
	}
	var errored []*encoding.ProcessedRequest
	var successful []*encoding.ProcessedRequest
	reqId := observations[0].RequestID
	finalResult := encoding.ProcessedRequest{
		RequestID: reqId,
		Result:    []byte{},
		Error:     []byte{},
	}
	for _, obs := range observations {
		if !bytes.Equal(obs.RequestID, reqId) {
			return nil, fmt.Errorf("inconsistent request IDs in aggregated observations %v vs %v", formatRequestId(obs.RequestID), formatRequestId(reqId))
		}
		if obs.GetError() != nil && len(obs.GetError()) > 0 {
			errored = append(errored, obs)
		} else {
			successful = append(successful, obs)
		}
	}
	resultIsError := len(errored) > len(successful)
	var toAggregate []*encoding.ProcessedRequest
	var rawData [][]byte
	if resultIsError {
		toAggregate = errored
		for _, item := range errored {
			rawData = append(rawData, item.Error)
		}
	} else {
		toAggregate = successful
		for _, item := range successful {
			rawData = append(rawData, item.Result)
		}
	}
	// Metadata (CallbackGasLimit, CoordinatorContract and OnchainMetadata) is aggregated using MODE method
	finalResult.CallbackGasLimit, finalResult.CoordinatorContract, finalResult.OnchainMetadata = aggregateMetadata(toAggregate)
	if resultIsError {
		// Errors are always aggregated using MODE method
		finalResult.Error = aggregateMode(rawData)
	} else {
		switch aggMethod {
		case config.AggregationMethod_AGGREGATION_MODE:
			finalResult.Result = aggregateMode(rawData)
		case config.AggregationMethod_AGGREGATION_MEDIAN:
			finalResult.Result = aggregateMedian(rawData)
		default:
			return nil, fmt.Errorf("unsupported aggregation method: %s", aggMethod)
		}
	}
	return &finalResult, nil
}

func aggregateMetadata(items []*encoding.ProcessedRequest) (uint32, []byte, []byte) {
	gasLimitBytes := make([][]byte, len(items))
	coordinatorContracts := make([][]byte, len(items))
	onchainMetadata := make([][]byte, len(items))
	for i, item := range items {
		gasLimitBytes[i] = make([]byte, 4)
		binary.BigEndian.PutUint32(gasLimitBytes[i], item.CallbackGasLimit)
		coordinatorContracts[i] = item.CoordinatorContract
		if coordinatorContracts[i] == nil {
			coordinatorContracts[i] = []byte{}
		}
		onchainMetadata[i] = item.OnchainMetadata
		if onchainMetadata[i] == nil {
			onchainMetadata[i] = []byte{}
		}
	}
	aggGasLimitBytes := aggregateMode(gasLimitBytes)
	aggGasLimitUint32 := binary.BigEndian.Uint32(aggGasLimitBytes)
	aggCoordinatorContract := aggregateMode(coordinatorContracts)
	aggOnchainMetadata := aggregateMode(onchainMetadata)
	return aggGasLimitUint32, aggCoordinatorContract, aggOnchainMetadata
}

func aggregateMode(items [][]byte) []byte {
	counts := make(map[string]int)
	var mostFrequent []byte
	highestFreq := 0
	for _, item := range items {
		str := string(item)
		currCount := counts[str] + 1
		counts[str] = currCount
		if currCount > highestFreq {
			highestFreq = currCount
			mostFrequent = item
		}
	}
	return mostFrequent
}

func aggregateMedian(items [][]byte) []byte {
	sort.Slice(items, func(i, j int) bool {
		if len(items[i]) != len(items[j]) {
			// NOTE: this doesn't account for extra leading zeros
			return len(items[i]) < len(items[j])
		}
		return bytes.Compare(items[i], items[j]) < 0
	})
	return items[(len(items)-1)/2]
}
