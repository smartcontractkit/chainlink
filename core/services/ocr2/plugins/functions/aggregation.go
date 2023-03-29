package functions

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
)

func CanAggregate(N int, F int, observations []*ProcessedRequest) bool {
	return N > 0 && F >= 0 && len(observations) > 0 && len(observations) <= N && len(observations) >= 2*F+1
}

func Aggregate(aggMethod config.AggregationMethod, observations []*ProcessedRequest) (*ProcessedRequest, error) {
	if len(observations) == 0 {
		return nil, fmt.Errorf("empty observation list passed for aggregation")
	}
	var errored []*ProcessedRequest
	var successful []*ProcessedRequest
	reqId := observations[0].RequestID
	finalResult := ProcessedRequest{
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
	var rawData [][]byte
	if len(errored) > len(successful) {
		for _, item := range errored {
			rawData = append(rawData, item.Error)
		}
		// Errors are always aggregated using MODE method
		finalResult.Error = aggregateMode(rawData)
		return &finalResult, nil
	}
	for _, item := range successful {
		rawData = append(rawData, item.Result)
	}
	switch aggMethod {
	case config.AggregationMethod_AGGREGATION_MODE:
		finalResult.Result = aggregateMode(rawData)
		return &finalResult, nil
	case config.AggregationMethod_AGGREGATION_MEDIAN:
		finalResult.Result = aggregateMedian(rawData)
		return &finalResult, nil
	default:
		return nil, fmt.Errorf("unsupported aggregation method: %s", aggMethod)
	}
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
	return []byte(mostFrequent)
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
