package median

import (
	"math/big"
	"sort"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

type AggregatedAttributedObservation struct {
	Timestamp       uint32
	Observers       [32]commontypes.OracleID
	Observations    []*big.Int
	JuelsPerFeeCoin *big.Int
}

func aggregate(observations []median.ParsedAttributedObservation) *AggregatedAttributedObservation {
	// defensive copy
	n := len(observations)
	observations = append([]median.ParsedAttributedObservation{}, observations...)

	aggregated := &AggregatedAttributedObservation{Observations: make([]*big.Int, len(observations))}

	sort.Slice(observations, func(i, j int) bool {
		return observations[i].Timestamp < observations[j].Timestamp
	})
	aggregated.Timestamp = observations[n/2].Timestamp

	// get median juelsPerFeeCoin
	sort.Slice(observations, func(i, j int) bool {
		return observations[i].JuelsPerFeeCoin.Cmp(observations[j].JuelsPerFeeCoin) < 0
	})
	aggregated.JuelsPerFeeCoin = observations[n/2].JuelsPerFeeCoin

	// sort by values
	sort.Slice(observations, func(i, j int) bool {
		return observations[i].Value.Cmp(observations[j].Value) < 0
	})

	for i, o := range observations {
		aggregated.Observers[i] = o.Observer
		aggregated.Observations[i] = o.Value
	}
	return aggregated
}
