package gen

import (
	"fmt"
	"sort"

	"github.com/leanovate/gopter"
)

// WeightedGen adds a weight number to a generator.
// To be used as parameter to gen.Weighted
type WeightedGen struct {
	Weight int
	Gen    gopter.Gen
}

// Weighted combines multiple generators, where each generator has a weight.
// The weight of a generator is proportional to the probability that the
// generator gets selected.
func Weighted(weightedGens []WeightedGen) gopter.Gen {
	if len(weightedGens) == 0 {
		panic("weightedGens must be non-empty")
	}
	weights := make(sort.IntSlice, 0, len(weightedGens))

	totalWeight := 0
	for _, weightedGen := range weightedGens {
		w := weightedGen.Weight
		if w <= 0 {
			panic(fmt.Sprintf(
				"weightedGens must have positive weights; got %d",
				w))
		}
		totalWeight += weightedGen.Weight
		weights = append(weights, totalWeight)
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		idx := weights.Search(1 + genParams.Rng.Intn(totalWeight))
		gen := weightedGens[idx].Gen
		result := gen(genParams)
		result.Sieve = nil
		return result
	}
}
