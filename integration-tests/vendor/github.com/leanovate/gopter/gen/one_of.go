package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

// OneConstOf generate one of a list of constant values
func OneConstOf(consts ...interface{}) gopter.Gen {
	if len(consts) == 0 {
		return Fail(reflect.TypeOf(nil))
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		idx := genParams.Rng.Intn(len(consts))
		return gopter.NewGenResult(consts[idx], gopter.NoShrinker)
	}
}

// OneGenOf generate one value from a a list of generators
func OneGenOf(gens ...gopter.Gen) gopter.Gen {
	if len(gens) == 0 {
		return Fail(reflect.TypeOf(nil))
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		idx := genParams.Rng.Intn(len(gens))
		return gens[idx](genParams)
	}
}
