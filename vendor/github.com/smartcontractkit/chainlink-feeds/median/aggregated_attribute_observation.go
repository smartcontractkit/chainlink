package median

import (
	"cmp"
	"math/big"
	"slices"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

type aggregatedAttributedObservation struct {
	Timestamp       uint32
	Observers       [32]commontypes.OracleID
	Observations    []*big.Int
	JuelsPerFeeCoin *big.Int
	GasPriceSubunit *big.Int
}

func aggregate(observations []median.ParsedAttributedObservation) *aggregatedAttributedObservation {
	// defensive copy
	n := len(observations)
	observations = slices.Clone(observations)

	aggregated := &aggregatedAttributedObservation{Observations: make([]*big.Int, len(observations))}

	slices.SortFunc(observations, func(a, b median.ParsedAttributedObservation) int {
		return cmp.Compare(a.Timestamp, b.Timestamp)
	})
	aggregated.Timestamp = observations[n/2].Timestamp

	slices.SortFunc(observations, func(a, b median.ParsedAttributedObservation) int {
		return a.JuelsPerFeeCoin.Cmp(b.JuelsPerFeeCoin)
	})
	aggregated.JuelsPerFeeCoin = observations[n/2].JuelsPerFeeCoin

	slices.SortFunc(observations, func(a, b median.ParsedAttributedObservation) int {
		return a.GasPriceSubunits.Cmp(b.GasPriceSubunits)
	})
	aggregated.GasPriceSubunit = observations[n/2].GasPriceSubunits

	slices.SortFunc(observations, func(a, b median.ParsedAttributedObservation) int {
		return a.Value.Cmp(b.Value)
	})

	for i, o := range observations {
		aggregated.Observers[i] = o.Observer
		aggregated.Observations[i] = o.Value
	}
	return aggregated
}
