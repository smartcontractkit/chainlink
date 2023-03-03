package gen

import "github.com/leanovate/gopter"

// Const creates a generator for a constant value
// Not the most exciting generator, but can be helpful from time to time
func Const(value interface{}) gopter.Gen {
	return func(*gopter.GenParameters) *gopter.GenResult {
		return gopter.NewGenResult(value, gopter.NoShrinker)
	}
}
