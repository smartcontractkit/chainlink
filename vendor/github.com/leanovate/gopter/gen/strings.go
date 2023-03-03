package gen

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/leanovate/gopter"
)

// RuneRange generates runes within a given range
func RuneRange(min, max rune) gopter.Gen {
	return genRune(Int64Range(int64(min), int64(max)))
}

// Rune generates an arbitrary character rune
func Rune() gopter.Gen {
	return genRune(Frequency(map[int]gopter.Gen{
		0xD800:                Int64Range(0, 0xD800),
		utf8.MaxRune - 0xDFFF: Int64Range(0xDFFF, int64(utf8.MaxRune)),
	}))
}

// RuneNoControl generates an arbitrary character rune that is not a control character
func RuneNoControl() gopter.Gen {
	return genRune(Frequency(map[int]gopter.Gen{
		0xD800:                Int64Range(32, 0xD800),
		utf8.MaxRune - 0xDFFF: Int64Range(0xDFFF, int64(utf8.MaxRune)),
	}))
}

func genRune(int64Gen gopter.Gen) gopter.Gen {
	return int64Gen.Map(func(value int64) rune {
		return rune(value)
	}).SuchThat(func(v rune) bool {
		return utf8.ValidRune(v)
	})
}

// NumChar generates arbitrary numberic character runes
func NumChar() gopter.Gen {
	return RuneRange('0', '9')
}

// AlphaUpperChar generates arbitrary uppercase alpha character runes
func AlphaUpperChar() gopter.Gen {
	return RuneRange('A', 'Z')
}

// AlphaLowerChar generates arbitrary lowercase alpha character runes
func AlphaLowerChar() gopter.Gen {
	return RuneRange('a', 'z')
}

// AlphaChar generates arbitrary character runes (upper- and lowercase)
func AlphaChar() gopter.Gen {
	return Frequency(map[int]gopter.Gen{
		0: AlphaUpperChar(),
		9: AlphaLowerChar(),
	})
}

// AlphaNumChar generates arbitrary alpha-numeric character runes
func AlphaNumChar() gopter.Gen {
	return Frequency(map[int]gopter.Gen{
		0: NumChar(),
		9: AlphaChar(),
	})
}

// UnicodeChar generates arbitrary character runes with a given unicode table
func UnicodeChar(table *unicode.RangeTable) gopter.Gen {
	if table == nil || len(table.R16)+len(table.R32) == 0 {
		return Fail(reflect.TypeOf(rune('a')))
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		tableIdx := genParams.Rng.Intn(len(table.R16) + len(table.R32))

		var selectedRune rune
		if tableIdx < len(table.R16) {
			r := table.R16[tableIdx]
			runeOffset := uint16(genParams.Rng.Int63n(int64((r.Hi-r.Lo+1)/r.Stride))) * r.Stride
			selectedRune = rune(runeOffset + r.Lo)
		} else {
			r := table.R32[tableIdx-len(table.R16)]
			runeOffset := uint32(genParams.Rng.Int63n(int64((r.Hi-r.Lo+1)/r.Stride))) * r.Stride
			selectedRune = rune(runeOffset + r.Lo)
		}
		genResult := gopter.NewGenResult(selectedRune, gopter.NoShrinker)
		genResult.Sieve = func(v interface{}) bool {
			return unicode.Is(table, v.(rune))
		}
		return genResult
	}
}

// AnyString generates an arbitrary string
func AnyString() gopter.Gen {
	return genString(Rune(), utf8.ValidRune)
}

// AlphaString generates an arbitrary string with letters
func AlphaString() gopter.Gen {
	return genString(AlphaChar(), unicode.IsLetter)
}

// NumString generates an arbitrary string with digits
func NumString() gopter.Gen {
	return genString(NumChar(), unicode.IsDigit)
}

// Identifier generates an arbitrary identifier string
// Identitiers are supporsed to start with a lowercase letter and contain only
// letters and digits
func Identifier() gopter.Gen {
	return gopter.CombineGens(
		AlphaLowerChar(),
		SliceOf(AlphaNumChar()),
	).Map(func(values []interface{}) string {
		first := values[0].(rune)
		tail := values[1].([]rune)
		result := make([]rune, 0, len(tail)+1)
		return string(append(append(result, first), tail...))
	}).SuchThat(func(str string) bool {
		if len(str) < 1 || !unicode.IsLower(([]rune(str))[0]) {
			return false
		}
		for _, ch := range str {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
				return false
			}
		}
		return true
	}).WithShrinker(StringShrinker)
}

// UnicodeString generates an arbitrary string from a given
// unicode table.
func UnicodeString(table *unicode.RangeTable) gopter.Gen {
	return genString(UnicodeChar(table), func(ch rune) bool {
		return unicode.Is(table, ch)
	})
}

func genString(runeGen gopter.Gen, runeSieve func(ch rune) bool) gopter.Gen {
	return SliceOf(runeGen).Map(runesToString).SuchThat(func(v string) bool {
		for _, ch := range v {
			if !runeSieve(ch) {
				return false
			}
		}
		return true
	}).WithShrinker(StringShrinker)
}

func runesToString(v []rune) string {
	return string(v)
}
