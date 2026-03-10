package aggregation

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

// Default MODE Aggregator needs a configurable number of identical responses for aggregation to succeed
type defaultModeAggregator struct {
	minIdenticalResponses uint32
}

var _ remotetypes.Aggregator = &defaultModeAggregator{}

func NewDefaultModeAggregator(minIdenticalResponses uint32) *defaultModeAggregator {
	return &defaultModeAggregator{
		minIdenticalResponses: minIdenticalResponses,
	}
}

func (a *defaultModeAggregator) Aggregate(_ string, responses [][]byte) (commoncap.TriggerResponse, error) {
	found, err := AggregateModeRaw(responses, a.minIdenticalResponses)
	if err != nil {
		return commoncap.TriggerResponse{}, fmt.Errorf("failed to aggregate responses, err: %w", err)
	}

	unmarshaled, err := pb.UnmarshalTriggerResponse(found)
	if err != nil {
		return commoncap.TriggerResponse{}, fmt.Errorf("failed to unmarshal aggregated responses, err: %w", err)
	}
	return unmarshaled, nil
}

func AggregateModeRaw(elemList [][]byte, minIdenticalResponses uint32) ([]byte, error) {
	// Fast path: if elemList[0] holds a strict majority (>len/2) hashing won't be needed.
	n := len(elemList)
	majorityThreshold := uint32(n/2 + 1)
	if minIdenticalResponses > majorityThreshold {
		majorityThreshold = minIdenticalResponses
	}
	var matchCount uint32 = 1
	for i := 1; i < n; i++ {
		if bytes.Equal(elemList[0], elemList[i]) {
			matchCount++
			if matchCount >= majorityThreshold {
				return elemList[0], nil
			}
		}
	}

	// Fallback: elemList[0] is not the majority, use hash-based mode finding.
	hashToCount := make(map[[32]byte]uint32, n)
	var found []byte
	hasher := sha256.New()
	for _, elem := range elemList {
		hasher.Reset()
		hasher.Write(elem)
		var key [32]byte
		hasher.Sum(key[:0])
		hashToCount[key]++
		if hashToCount[key] >= minIdenticalResponses {
			found = elem
			minIdenticalResponses = hashToCount[key]
		}
	}
	if found == nil {
		return nil, errors.New("not enough identical responses found")
	}
	return found, nil
}
