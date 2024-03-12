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

	var value values.Value
	switch parts[1] {
	case "inputs":
		value = sc.inputs
	case "outputs":
		if sc.outputs.err != nil {
			return "", fmt.Errorf("cannot interpolate ref part `%s` in `%+v`: step has errored", parts[1], sc)
		}

		value = sc.outputs.value
	default:
		return "", fmt.Errorf("cannot interpolate ref part `%s` in `%+v`: second part must be `inputs` or `outputs`", parts[1], sc)
	}

	val, err := values.Unwrap(value)
	if err != nil {
		return "", err
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

			if d < 0 {
				return "", fmt.Errorf("could not interpolate ref part `%s` in `%+v`: index %d must be a positive number", r, v, d)
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

// findAndInterpolateAllKeys takes an `input` any value, and recursively
// identifies any values that should be replaced from `state`.
// A value `v` should be replaced if it is wrapped as follows `$(v)`.
func findAndInterpolateAllKeys(input any, state *executionState) (any, error) {
	switch tv := input.(type) {
	case string:
		matches := interpolationTokenRe.FindStringSubmatch(tv)
		if len(matches) < 2 {
			return tv, nil
		}

		interpolatedVar := matches[1]
		nv, err := interpolateKey(interpolatedVar, state)
		if err != nil {
			return nil, err
		}

		return nv, nil
	case map[string]any:
		nm := map[string]any{}
		for k, v := range tv {
			nv, err := findAndInterpolateAllKeys(v, state)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case []any:
		a := []any{}
		for _, el := range tv {
			ne, err := findAndInterpolateAllKeys(el, state)
			if err != nil {
				return nil, err
			}

			a = append(a, ne)
		}
		return a, nil
	}

	return nil, fmt.Errorf("cannot interpolate item %+v of type %T", input, input)
}
