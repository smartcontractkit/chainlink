package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// GetterFunc is a function that either returns a value or an error.
type GetterFunc func() (interface{}, error)

// From creates []GetterFunc from a mix of getters or bare values.
func From(getters ...interface{}) []GetterFunc {
	var gfs []GetterFunc
	for _, g := range getters {
		switch v := g.(type) {
		case GetterFunc:
			gfs = append(gfs, v)

		default:
			// If a bare value is passed in, create a simple getter
			gfs = append(gfs, func() (interface{}, error) {
				return v, nil
			})
		}
	}
	return gfs
}

// NonemptyString creates a getter to ensure the string is non-empty.
func NonemptyString(s string) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		return trimmed, nil
	}
}

// ValidDurationInSeconds creates a getter to ensure the string is a valid duration and return duration in seconds.
func ValidDurationInSeconds(s string) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		dr, err := time.ParseDuration(s)
		if err != nil {
			return nil, err
		}
		return int(dr.Seconds()), nil
	}
}

// Input creates a getter returning inputs[index] value, or error if index is out of range.
func Input(inputs []Result, index int) GetterFunc {
	return func() (interface{}, error) {
		if index < 0 || index >= len(inputs) {
			return nil, ErrIndexOutOfRange
		}
		return inputs[index].Value, inputs[index].Error
	}
}

// Inputs creates a getter returning array of Result.Value (or Result.Error where not nil).
func Inputs(inputs []Result) GetterFunc {
	return func() (interface{}, error) {
		var vals []interface{}
		for _, input := range inputs {
			if input.Error != nil {
				vals = append(vals, input.Error)
			} else {
				vals = append(vals, input.Value)
			}
		}
		return vals, nil
	}
}

// VarExpr creates a getter interpolating expr value using the given Vars.
// The expression allows whitespace on both ends that will be trimmed.
// Expr examples: $(foo.bar), $(arr.1), $(bar)
func VarExpr(expr string, vars Vars) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(expr)
		if len(trimmed) < 3 {
			return nil, ErrParameterEmpty
		}
		isVariableExpr := strings.Count(trimmed, "$") == 1 && trimmed[:2] == "$(" && trimmed[len(trimmed)-1] == ')'
		if !isVariableExpr {
			return nil, ErrParameterEmpty
		}
		keypath := strings.TrimSpace(trimmed[2 : len(trimmed)-1])
		if len(keypath) == 0 {
			return nil, ErrParameterEmpty
		}
		val, err := vars.Get(keypath)
		if err != nil {
			return nil, err
		} else if as, is := val.(error); is {
			return nil, errors.Wrapf(ErrTooManyErrors, "VarExpr: %v", as)
		}
		return val, nil
	}
}

// JSONWithVarExprs creates a getter that unmarshals jsExpr string as JSON, and
// interpolates all variables expressions found in jsExpr from Vars.
// The getter returns the unmarshalled object having expressions interpolated from Vars.
// allowErrors flag indicates if interpolating values stored in Vars can be errors.
// jsExpr example: {"requestId": $(decode_log.requestId), "payment": $(decode_log.payment)}
func JSONWithVarExprs(jsExpr string, vars Vars, allowErrors bool) GetterFunc {
	return func() (interface{}, error) {
		if strings.TrimSpace(jsExpr) == "" {
			return nil, ErrParameterEmpty
		}
		const chainlinkKeyPath = "__chainlink_key_path__"
		replaced := variableRegexp.ReplaceAllFunc([]byte(jsExpr), func(expr []byte) []byte {
			keypathStr := strings.TrimSpace(string(expr[2 : len(expr)-1]))
			return []byte(fmt.Sprintf(`{ "%s": "%s" }`, chainlinkKeyPath, keypathStr))
		})

		var val interface{}
		jd := json.NewDecoder(bytes.NewReader(replaced))
		jd.UseNumber()
		if err := jd.Decode(&val); err != nil {
			return nil, errors.Wrapf(ErrBadInput, "while unmarshalling JSON: %v; js: %s", err, string(replaced))
		}
		reinterpreted, err := reinterpetJsonNumbers(val)
		if err != nil {
			return nil, errors.Wrapf(ErrBadInput, "while processing json.Number: %v; js: %s", err, string(replaced))
		}
		val = reinterpreted

		return mapGoValue(val, func(val interface{}) (interface{}, error) {
			if m, is := val.(map[string]interface{}); is {
				maybeKeypath, exists := m[chainlinkKeyPath]
				if !exists {
					return val, nil
				}
				keypath, is := maybeKeypath.(string)
				if !is {
					return nil, errors.Wrapf(ErrBadInput, fmt.Sprintf("you cannot use %s in your JSON", chainlinkKeyPath))
				}
				newVal, err := vars.Get(keypath)
				if err != nil {
					return nil, err
				} else if err, is := newVal.(error); is && !allowErrors {
					return nil, errors.Wrapf(ErrBadInput, "error is not allowed: %v", err)
				}
				return newVal, nil
			}
			return val, nil
		})
	}
}

// mapGoValue iterates on v object recursively and calls fn for each value.
// Used by JSONWithVarExprs to interpolate all variables expressions.
func mapGoValue(v interface{}, fn func(val interface{}) (interface{}, error)) (x interface{}, err error) {
	type item struct {
		val         interface{}
		parentMap   map[string]interface{}
		parentKey   string
		parentSlice []interface{}
		parentIdx   int
	}

	stack := []item{{val: v}}
	var current item

	for len(stack) > 0 {
		current = stack[0]
		stack = stack[1:]

		val, err := fn(current.val)
		if err != nil {
			return nil, err
		}

		if current.parentMap != nil {
			current.parentMap[current.parentKey] = val
		} else if current.parentSlice != nil {
			current.parentSlice[current.parentIdx] = val
		}

		if asMap, isMap := val.(map[string]interface{}); isMap {
			for key := range asMap {
				stack = append(stack, item{val: asMap[key], parentMap: asMap, parentKey: key})
			}
		} else if asSlice, isSlice := val.([]interface{}); isSlice {
			for i := range asSlice {
				stack = append(stack, item{val: asSlice[i], parentSlice: asSlice, parentIdx: i})
			}
		}
	}
	return v, nil
}
