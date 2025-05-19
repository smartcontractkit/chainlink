package aggregation

import (
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
