package gen

import "github.com/leanovate/gopter"

// RetryUntil creates a generator that retries a given generator until a condition in met.
// condition: has to be a function with one parameter (matching the generated value of gen) returning a bool.
// Note: The new generator will only create an empty result once maxRetries is reached.
// Depending on the hit-ratio of the condition is may result in long running tests, use with care.
func RetryUntil(gen gopter.Gen, condition interface{}, maxRetries int) gopter.Gen {
	genWithSieve := gen.SuchThat(condition)
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		for i := 0; i < maxRetries; i++ {
			result := genWithSieve(genParams)
			if _, ok := result.Retrieve(); ok {
				return result
			}
		}
		resultType := gen(genParams).ResultType
		return gopter.NewEmptyResult(resultType)
	}
}
