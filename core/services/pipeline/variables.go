package pipeline

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//go:generate mockery --name PipelineParamUnmarshaler --output ./mocks/ --case=underscore

type (
	Vars map[string]interface{}

	PipelineParamUnmarshaler interface {
		UnmarshalPipelineParam(val interface{}, vars Vars) error
	}
)

var (
	variableRegexp = regexp.MustCompile(`\$\(([a-zA-Z0-9_\.]+)\)`)

	ErrKeypathNotFound = errors.New("keypath not found")
)

func NewVars() Vars {
	return make(Vars)
}

func (vars Vars) Get(keypath string) (interface{}, error) {
	parts := keypathParts(keypath)
	if len(parts) == 0 {
		return (map[string]interface{})(vars), nil
	}
	return vars.traverse(parts, false)
}

func (vars Vars) Set(keypath string, val interface{}) error {
	parts := keypathParts(keypath)
	if len(parts) == 0 {
		return errors.New("can't set the root of a Vars")
	}

	last, err := vars.traverse(parts[:len(parts)-1], true)
	if err != nil {
		return err
	}
	lastKey := parts[len(parts)-1]

	switch typed := last.(type) {
	case map[string]interface{}:
		typed[string(lastKey)] = val

	case []interface{}:
		idx, err := strconv.Atoi(string(lastKey))
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

var keypathSeparator = []byte(".")

func keypathParts(keypath string) [][]byte {
	if len(keypath) == 0 {
		return nil
	}
	// The bytes package uses platform-dependent hardware optimizations and
	// avoids the extra allocations that are required to work with strings.
	// Keypaths have to be parsed quite a bit, so let's do it well.
	kp := []byte(keypath)
	n := 1 + bytes.Count(kp, keypathSeparator)
	parts := make([][]byte, n)
	for i := 0; i < n-1; i++ {
		nextSep := bytes.IndexByte(kp, keypathSeparator[0])
		parts[i] = kp[:nextSep]
		kp = kp[nextSep+1:]
	}
	parts[len(parts)-1] = kp
	return parts
}

func (vars Vars) traverse(keypathParts [][]byte, create bool) (interface{}, error) {
	var cur interface{} = (map[string]interface{})(vars)

	for _, key := range keypathParts {
		switch typed := cur.(type) {
		case map[string]interface{}:
			var exists bool
			cur, exists = typed[string(key)] // Converting []byte to string to access a map is a special-case optimization in Go
			if !exists && !create {
				return nil, errors.Wrapf(ErrKeypathNotFound, "key %v / keypath %v", string(key), bytesToStrings(keypathParts))
			} else if !exists {
				typed[string(key)] = make(map[string]interface{})
				cur = typed[string(key)]
			}

		case []interface{}:
			idx, err := strconv.ParseInt(string(key), 10, 64)
			if err != nil {
				return nil, err
			} else if idx > int64(len(typed)-1) {
				return nil, errors.Wrapf(ErrKeypathNotFound, "index %v out of range (length %v / keypath %v)", idx, len(typed), bytesToStrings(keypathParts))
			}
			cur = typed[idx]

		default:
			return nil, errors.Wrapf(ErrKeypathNotFound, "encountered non-map/non-slice (keypath %v)", bytesToStrings(keypathParts))
		}
	}
	return cur, nil
}

func bytesToStrings(bs [][]byte) []string {
	var s []string
	for _, b := range bs {
		s = append(s, string(b))
	}
	return s
}

func (vars Vars) ResolveValue(out PipelineParamUnmarshaler, getters GetterFuncs) error {
	var val interface{}
	var err error
	var found bool
	for _, get := range getters {
		val, err = get(vars)
		if errors.Cause(err) == ErrParameterEmpty {
			continue
		} else if err != nil {
			return err
		}
		found = true
		break
	}
	if !found {
		return ErrParameterEmpty
	}

	err = out.UnmarshalPipelineParam(val, vars)
	if err != nil {
		return err
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
			// If a bare value is passed in, create a simple getter
			gfs = append(gfs, func(_ Vars) (interface{}, error) {
				return v, nil
			})
		}
	}
	return gfs
}

type GetterFunc func(vars Vars) (interface{}, error)

func VariableExpr(s string) GetterFunc {
	return func(vars Vars) (interface{}, error) {
		keypath, ok := variableExprKeypath(s)
		if !ok {
			return nil, ErrParameterEmpty
		}
		return vars.Get(keypath)
	}
}

func NonemptyString(s string) GetterFunc {
	return func(_ Vars) (interface{}, error) {
		trimmed := strings.TrimSpace(s)
		if len(trimmed) == 0 {
			return nil, ErrParameterEmpty
		}
		return trimmed, nil
	}
}

func Input(inputs []Result, index int) GetterFunc {
	return func(_ Vars) (interface{}, error) {
		if len(inputs)-1 < index {
			return nil, ErrParameterEmpty
		}
		return inputs[index].Value, inputs[index].Error
	}
}

func Inputs(inputs []Result) GetterFunc {
	return func(_ Vars) (interface{}, error) {
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

func variableExprKeypath(s string) (keypath string, ok bool) {
	trimmed := strings.TrimSpace(s)
	if strings.Count(trimmed, "$") == 1 && trimmed[:2] == "$(" && trimmed[len(trimmed)-1] == ')' {
		return strings.TrimSpace(trimmed[2 : len(trimmed)-1]), true
	}
	return "", false
}

func CheckInputs(inputs []Result, minLen, maxLen, maxErrors int) ([]interface{}, error) {
	if minLen >= 0 && len(inputs) < minLen {
		return nil, ErrWrongInputCardinality
	} else if maxLen >= 0 && len(inputs) > maxLen {
		return nil, ErrWrongInputCardinality
	}
	var vals []interface{}
	var errs int
	for _, input := range inputs {
		if input.Error != nil {
			errs++
			continue
		}
		vals = append(vals, input.Value)
	}
	if maxErrors >= 0 && errs > maxErrors {
		return nil, ErrTooManyErrors
	}
	return vals, nil
}
