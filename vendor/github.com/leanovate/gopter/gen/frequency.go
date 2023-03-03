package gen

import (
	"sort"

	"github.com/leanovate/gopter"
)

// Frequency combines multiple weighted generators of the the same result type
// The generators from weightedGens will be used accrding to the weight, i.e. generators
// with a hight weight will be used more often than generators with a low weight.
func Frequency(weightedGens map[int]gopter.Gen) gopter.Gen {
	if len(weightedGens) == 0 {
		return Fail(nil)
	}
	weights := make(sort.IntSlice, 0, len(weightedGens))
	max := 0
	for weight := range weightedGens {
		if weight > max {
			max = weight
		}
		weights = append(weights, weight)
	}
	weights.Sort()
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		idx := weights.Search(genParams.Rng.Intn(max + 1))
		gen := weightedGens[weights[idx]]

		result := gen(genParams)
		result.Sieve = nil
		return result
	}
}
