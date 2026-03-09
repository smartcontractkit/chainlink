package aggregation

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	hashToCount := make(map[string]uint32)
	var found []byte
	for _, elem := range elemList {
		hasher := sha256.New()
		hasher.Write(elem)
		sha := hex.EncodeToString(hasher.Sum(nil))
		hashToCount[sha]++
		if hashToCount[sha] >= minIdenticalResponses {
			found = elem
			// update in case we find another elem with an even higher count
			minIdenticalResponses = hashToCount[sha]
		}
	}
	if found == nil {
		return nil, errors.New("not enough identical responses found")
	}
	return found, nil
}

func AggregateModeRawFast(elemList [][]byte, minIdenticalResponses uint32) ([]byte, error) {
	// Fast path: if all elements are byte-equal, the mode is trivially elemList[0].
	allEqual := len(elemList) >= int(minIdenticalResponses)
	for i := 1; i < len(elemList) && allEqual; i++ {
		allEqual = bytes.Equal(elemList[0], elemList[i])
	}
	if allEqual {
		return elemList[0], nil
	}

	hashToCount := make(map[[32]byte]uint32, len(elemList))
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
