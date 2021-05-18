package pipeline

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	Vars map[string]interface{}

	PipelineParamUnmarshaler interface {
		UnmarshalPipelineParam(val interface{}, vars Vars) error
	}
)

var variableRegexp = regexp.MustCompile(`\$\(([a-zA-Z0-9_\.]+)\)`)

func NewVars() Vars {
	return make(Vars)
}

func (v Vars) Get(keypath string) (interface{}, error) {
	keypathParts := strings.Split(keypath, ".")
	return v.traverse(keypathParts, false)
}

func (v Vars) Set(keypath string, val interface{}) error {
	keypathParts := strings.Split(keypath, ".")

	last, err := v.traverse(keypathParts[:len(keypathParts)-1], true)
	if err != nil {
		return err
	}
	lastKey := keypathParts[len(keypathParts)-1]

	switch typed := last.(type) {
	case map[string]interface{}:
		typed[lastKey] = val

	case []interface{}:
		idx, err := strconv.Atoi(lastKey)
		if err != nil {
			return err
		} else if len(typed) <= idx {
			return errors.New("index out of range")
		}
		typed[idx] = val

	default:
		return errors.New("encountered non-map/non-slice")
	}
	return nil
}

func (v Vars) traverse(keypathParts []string, create bool) (interface{}, error) {
	type M = map[string]interface{}
	var cur interface{} = M(v)

	for _, key := range keypathParts {
		switch typed := cur.(type) {
		case map[string]interface{}:
			var exists bool
			cur, exists = typed[key]
			if !exists && !create {
				return nil, errors.Errorf("not found: key %v keypathParts %v", key, keypathParts)
			} else if !exists {
				typed[key] = make(map[string]interface{})
				cur = typed[key]
			}

		case []interface{}:
			idx, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				return nil, err
			} else if idx > int64(len(typed)-1) {
				return nil, errors.New("index out of range")
			}
			cur = typed[idx]

		default:
			return nil, errors.New("encountered non-map/non-slice")
		}
	}
	return cur, nil
}

func (vars Vars) ResolveValue(out PipelineParamUnmarshaler, getters GetterFuncs, validators ...ValidatorFunc) error {
	var val interface{}
	var err error
	for _, get := range getters {
		val, err = get()
		if errors.Cause(err) == ErrParameterEmpty {
			continue
		} else if err != nil {
			return err
		}
		break
	}

	err = out.UnmarshalPipelineParam(val, vars)
	if err != nil {
		return err
	}

	for _, validate := range validators {
		err := validate(out)
		if err != nil {
			return err
		}
	}

	return nil
}

type GetterFuncs []GetterFunc

func From(getters ...interface{}) GetterFuncs {
	var gfs GetterFuncs
	for _, g := range getters {
		switch v := g.(type) {
		case GetterFunc:
			gfs = append(gfs, v)
		default:
			gfs = append(gfs, func() (interface{}, error) {
				return v, nil
			})
		}
	}
	return gfs
}

func (gf GetterFuncs) Or(getter GetterFunc) GetterFuncs {
	return append(gf, getter)
}

type GetterFunc func() (interface{}, error)

func VariableExpr(s string) GetterFunc {
	return func() (interface{}, error) {
		is, _ := isPureVariableExprString(s)
		if !is {
			return nil, ErrParameterEmpty
		}
		return s, nil
	}
}

func NonemptyString(s string) GetterFunc {
	return func() (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		return trimmed, nil
	}
}

func Input(inputs []Result, index int) GetterFunc {
	return func() (interface{}, error) {
		if len(inputs)-1 < index {
			return nil, ErrParameterEmpty
		}
		return inputs[index].Value, inputs[index].Error
	}
}

type ValidatorFunc func(val interface{}) error

func Validate(validators ...ValidatorFunc) []ValidatorFunc {
	return validators
}

func MapKeys(keys ...string) ValidatorFunc {
	return func(val interface{}) error {
		asMap, isMap := val.(map[string]interface{})
		if !isMap {
			return ErrBadInput
		}
		for _, k := range keys {
			_, exists := asMap[k]
			if !exists {
				return ErrBadInput
			}
		}
		return nil
	}
}

// func resolve(from, to interface{}, vars Vars) (interface{}, error) {
// 	switch val := from.(type) {
// 	case string:
// 		return resolveString(val, vars)
// 	case map[string]interface{}:
// 		return resolveMap(val, vars)
// 	// case []interface{}:
// 	// 	return resolveSlice(val, vars)
// 	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
// 		return val, nil
// 	default:
// 		return nil, ErrBadInput
// 	}
// }

// func resolveString(s string, vars Vars) (string, error) {
// 	var err error
// 	resolved := variableRegexp.ReplaceAllStringFunc(s, func(keypath string) string {
// 		val, err2 := vars.Get(keypath)
// 		if err2 != nil {
// 			err = multierr.Append(err, err2)
// 			return ""
// 		}
// 		return fmt.Sprintf("%s", val)
// 	})
// 	return resolved, err
// }

func isPureVariableExprString(s string) (is bool, keypath string) {
	trimmed := strings.TrimSpace(s)
	if strings.Count(trimmed, "$") == 1 && trimmed[:2] == "$(" && trimmed[len(trimmed)-1] == ')' {
		return true, trimmed[2 : len(trimmed)-1]
	}
	return false, ""
}

// func resolveMap(m map[string]interface{}, vars Vars) (map[string]interface{}, error) {
// 	m2 := make(map[string]interface{})
// 	for k, v := range m {
// 		resolvedKey, err := resolveString(k, vars)
// 		if err != nil {
// 			return nil, err
// 		}

// 		switch val := v.(type) {
// 		case string:

// 		}

// 		resolvedVal, err := resolve(v, vars)
// 		if err != nil {
// 			return nil, err
// 		}
// 		m2[resolvedKey] = resolvedVal
// 	}
// 	return m2, nil
// }

// func resolveValue(s string, vars Vars) interface{} {

// }

func trimStrings(strs []string) []string {
	trimmed := make([]string, len(strs))
	for i := range strs {
		trimmed[i] = strings.TrimSpace(strs[i])
	}
	return trimmed
}
