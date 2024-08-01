package llo

import (
	"errors"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type AggregatorFunc func(values []StreamValue, f int) (StreamValue, error)

func GetAggregatorFunc(a llotypes.Aggregator) AggregatorFunc {
	switch a {
	case llotypes.AggregatorMedian:
		return MedianAggregator
	case llotypes.AggregatorMode:
		return ModeAggregator
	case llotypes.AggregatorQuote:
		return QuoteAggregator
	default:
		return nil
	}
}

func MedianAggregator(values []StreamValue, f int) (StreamValue, error) {
	observations := make([]decimal.Decimal, 0, len(values))
	for _, value := range values {
		switch v := value.(type) {
		case *Decimal:
			observations = append(observations, v.Decimal())
		default:
			// Unexpected type, skip
			continue
		}
	}
	if len(observations) <= f {
		// In the worst case, we have 2f+1 observations, of which up to f
		// are allowed to be invalid/missing. If we have less than f+1
		// usable observations, we cannot securely generate a median at
		// all.
		return nil, fmt.Errorf("not enough observations to calculate median, expected at least f+1, got %d", len(observations))
	}
	sort.Slice(observations, func(i, j int) bool { return observations[i].Cmp(observations[j]) < 0 })
	// We use a "rank-k" median here, instead one could average in case of
	// an even number of observations.
	// In the case of an even number, the higher value is chosen.
	// e.g. [1, 2, 3, 4] -> 3
	return ToDecimal(observations[len(observations)/2]), nil
}

func ModeAggregator(values []StreamValue, f int) (StreamValue, error) {
	return nil, errors.New("not implemented")
}

func QuoteAggregator(values []StreamValue, f int) (StreamValue, error) {
	var observations []*Quote
	for _, value := range values {
		if v, ok := value.(*Quote); !ok {
			// Unexpected type, skip
			continue
		} else if v.IsValid() {
			observations = append(observations, v)
		}
		// Exclude Quotes that violate bid<=mid<=ask
	}
	if len(observations) <= f {
		// In the worst case, we have 2f+1 observations, of which up to f
		// are allowed to be invalid/missing. If we have less than f+1
		// usable observations, we cannot securely generate a median at
		// all.
		return nil, fmt.Errorf("not enough valid observations to aggregate quote, expected at least f+1, got %d", len(observations))
	}
	// Calculate "rank-k" median for benchmark, bid and ask separately.
	// This is guaranteed not to return values that violate bid<=mid<=ask due
	// to the filter of observations above.
	q := Quote{}
	sort.Slice(observations, func(i, j int) bool { return observations[i].Benchmark.Cmp(observations[j].Benchmark) < 0 })
	q.Benchmark = observations[len(observations)/2].Benchmark
	sort.Slice(observations, func(i, j int) bool { return observations[i].Bid.Cmp(observations[j].Bid) < 0 })
	q.Bid = observations[len(observations)/2].Bid
	sort.Slice(observations, func(i, j int) bool { return observations[i].Ask.Cmp(observations[j].Ask) < 0 })
	q.Ask = observations[len(observations)/2].Ask
	return &q, nil
}
