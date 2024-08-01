package gen

import "github.com/leanovate/gopter"

var runeSliceShrinker = SliceShrinker(gopter.NoShrinker)

// StringShrinker is a shrinker for strings.
// It is very similar to a slice shrinker just that the elements themselves will not be shrunk.
func StringShrinker(v interface{}) gopter.Shrink {
	return runeSliceShrinker([]rune(v.(string))).Map(runesToString)
}
