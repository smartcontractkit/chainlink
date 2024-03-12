package workflows

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type stepOutput struct {
	err   error
	value values.Value
}

type stepState struct {
	inputs  *values.Map
	outputs *stepOutput
}

type executionState struct {
	steps       map[string]*stepState
	executionID string
	workflowID  string
}

// interpolateKey takes a multi-part, dot-separated key and attempts to replace
// it with its corresponding value in `state`.
// A key is valid if:
// - it contains at least two parts, with the first part being the workflow step's `ref` variable, and the second being one of `inputs` or `outputs`
// - any subsequent parts will be processed as a list index (if the current element is a list) or a map key (if it's a map)
func interpolateKey(key string, state *executionState) (any, error) {
	parts := strings.Split(key, ".")

	if len(parts) < 2 {
		return "", fmt.Errorf("cannot interpolate %s: must have at least two parts", key)
	}

	sc, ok := state.steps[parts[0]]
	if !ok {
		return "", fmt.Errorf("could not find ref `%s`", parts[0])
	}

	var val any
	switch parts[1] {
	case "inputs":
		v, err := values.Unwrap(sc.inputs)
		if err != nil {
			return nil, err
		}
		val = v
	case "outputs":
		if sc.outputs.err != nil {
			return "", fmt.Errorf("cannot interpolate ref part `%s` in `%+v`: step has errored", parts[1], val)
		}

		v, err := values.Unwrap(sc.outputs.value)
		if err != nil {
			return "", err
		}

		val = v
	default:
		return "", fmt.Errorf("cannot interpolate ref part `%s` in `%+v`: second part must be `inputs` or `outputs`", parts[1], val)
	}

	remainingParts := parts[2:]
	for _, r := range remainingParts {
		switch v := val.(type) {
		case map[string]any:
			inner, ok := v[r]
			if !ok {
				return "", fmt.Errorf("could not find ref part `%s` in `%+v`", r, v)
			}

			val = inner
		case []any:
			d, err := strconv.Atoi(r)
			if err != nil {
				return "", fmt.Errorf("could not interpolate ref part `%s` in `%+v`: `%s` is not convertible to an int", r, v, r)
			}

			if d > len(v)-1 {
				return "", fmt.Errorf("could not interpolate ref part `%s` in `%+v`: cannot fetch index %d", r, v, d)
			}

			val = v[d]
		default:
			return "", fmt.Errorf("could not interpolate ref part `%s` in `%+v`", r, val)
		}
	}

	return val, nil
}

var (
	interpolationTokenRe = regexp.MustCompile(`^\$\((\S+)\)$`)
)

// interpolateInputsFromState takes an `m` any value, and recursively
// identifies any values that should be replaced from `es`.
// A value `v` should be replaced if it is wrapped as follows `$(v)`.
func interpolateInputsFromState(m any, es *executionState) (any, error) {
	switch tv := m.(type) {
	case string:
		matches := interpolationTokenRe.FindStringSubmatch(tv)
		if len(matches) < 2 {
			return tv, nil
		}

		interpolatedVar := matches[1]
		nv, err := interpolateKey(interpolatedVar, es)
		if err != nil {
			return nil, err
		}

		return nv, nil
	case map[string]any:
		nm := map[string]any{}
		for k, v := range tv {
			nv, err := interpolateInputsFromState(v, es)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case []any:
		a := []any{}
		for _, el := range tv {
			ne, err := interpolateInputsFromState(el, es)
			if err != nil {
				return nil, err
			}

			a = append(a, ne)
		}
		return a, nil
	}

	return nil, fmt.Errorf("cannot interpolate item %+v of type %T", m, m)
}
