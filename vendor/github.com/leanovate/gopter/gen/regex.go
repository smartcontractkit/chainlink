package gen

import (
	"reflect"
	"regexp"
	"regexp/syntax"
	"strings"

	"github.com/leanovate/gopter"
)

// RegexMatch generates matches for a given regular expression
// regexStr is supposed to conform to the perl regular expression syntax
func RegexMatch(regexStr string) gopter.Gen {
	regexSyntax, err1 := syntax.Parse(regexStr, syntax.Perl)
	regex, err2 := regexp.Compile(regexStr)
	if err1 != nil || err2 != nil {
		return Fail(reflect.TypeOf(""))
	}
	return regexMatchGen(regexSyntax.Simplify()).SuchThat(func(v string) bool {
		return regex.MatchString(v)
	}).WithShrinker(StringShrinker)
}

func regexMatchGen(regex *syntax.Regexp) gopter.Gen {
	switch regex.Op {
	case syntax.OpLiteral:
		return Const(string(regex.Rune))
	case syntax.OpCharClass:
		gens := make([]gopter.Gen, 0, len(regex.Rune)/2)
		for i := 0; i+1 < len(regex.Rune); i += 2 {
			gens = append(gens, RuneRange(regex.Rune[i], regex.Rune[i+1]).Map(runeToString))
		}
		return OneGenOf(gens...)
	case syntax.OpAnyChar:
		return Rune().Map(runeToString)
	case syntax.OpAnyCharNotNL:
		return RuneNoControl().Map(runeToString)
	case syntax.OpCapture:
		return regexMatchGen(regex.Sub[0])
	case syntax.OpStar:
		elementGen := regexMatchGen(regex.Sub[0])
		return SliceOf(elementGen).Map(func(v []string) string {
			return strings.Join(v, "")
		})
	case syntax.OpPlus:
		elementGen := regexMatchGen(regex.Sub[0])
		return gopter.CombineGens(elementGen, SliceOf(elementGen)).Map(func(vs []interface{}) string {
			return vs[0].(string) + strings.Join(vs[1].([]string), "")
		})
	case syntax.OpQuest:
		elementGen := regexMatchGen(regex.Sub[0])
		return OneGenOf(Const(""), elementGen)
	case syntax.OpConcat:
		gens := make([]gopter.Gen, len(regex.Sub))
		for i, sub := range regex.Sub {
			gens[i] = regexMatchGen(sub)
		}
		return gopter.CombineGens(gens...).Map(func(v []interface{}) string {
			result := ""
			for _, str := range v {
				result += str.(string)
			}
			return result
		})
	case syntax.OpAlternate:
		gens := make([]gopter.Gen, len(regex.Sub))
		for i, sub := range regex.Sub {
			gens[i] = regexMatchGen(sub)
		}
		return OneGenOf(gens...)
	}
	return Const("")
}

func runeToString(v rune) string {
	return string(v)
}
