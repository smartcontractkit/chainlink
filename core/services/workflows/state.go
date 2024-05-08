package workflows

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

// copyState returns a deep copy of the input executionState
func copyState(es store.WorkflowExecution) store.WorkflowExecution {
	steps := map[string]*store.WorkflowExecutionStep{}
	for ref, step := range es.Steps {
		var mval *values.Map
		if step.Inputs != nil {
			mp := values.Proto(step.Inputs).GetMapValue()
			mval = values.FromMapValueProto(mp)
		}

		op := values.Proto(step.Outputs.Value)
		copiedov := values.FromProto(op)

		newState := &store.WorkflowExecutionStep{
			ExecutionID: step.ExecutionID,
			Ref:         step.Ref,
			Status:      step.Status,

			Outputs: &store.StepOutput{
				Err:   step.Outputs.Err,
				Value: copiedov,
			},

			Inputs: mval,
		}

		steps[ref] = newState
	}
	return store.WorkflowExecution{
		ExecutionID: es.ExecutionID,
		WorkflowID:  es.WorkflowID,
		Status:      es.Status,
		Steps:       steps,
	}
}

// interpolateKey takes a multi-part, dot-separated key and attempts to replace
// it with its corresponding value in `state`.
//
// A key is valid if it contains at least two parts, with:
//   - the first part being the workflow step's `ref` variable
//   - the second part being one of `inputs` or `outputs`
//
// If a key has more than two parts, then we traverse the parts
// to find the value we want to replace.
// We support traversing both nested maps and lists and any combination of the two.
func interpolateKey(key string, state store.WorkflowExecution) (any, error) {
	parts := strings.Split(key, ".")

	if len(parts) < 2 {
		return "", fmt.Errorf("cannot interpolate %s: must have at least two parts", key)
	}

	// lookup the step we want to get either input or output state from
	sc, ok := state.Steps[parts[0]]
	if !ok {
		return "", fmt.Errorf("could not find ref `%s`", parts[0])
	}

	var value values.Value
	switch parts[1] {
	case "inputs":
		value = sc.Inputs
	case "outputs":
		if sc.Outputs.Err != nil {
			return "", fmt.Errorf("cannot interpolate ref part `%s` in `%+v`: step has errored", parts[1], sc)
		}

		value = sc.Outputs.Value
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
				return "", fmt.Errorf("could not find ref part `%s` (ref: `%s`) in `%+v`", r, key, v)
			}

			val = inner
		case []any:
			i, err := strconv.Atoi(r)
			if err != nil {
				return "", fmt.Errorf("could not interpolate ref part `%s` (ref: `%s`) in `%+v`: `%s` is not convertible to an int", r, key, v, r)
			}

			if (i > len(v)-1) || (i < 0) {
				return "", fmt.Errorf("could not interpolate ref part `%s` (ref: `%s`) in `%+v`: index out of bounds %d", r, key, v, i)
			}

			val = v[i]
		default:
			return "", fmt.Errorf("could not interpolate ref part `%s` (ref: `%s`) in `%+v`", r, key, val)
		}
	}

	return val, nil
}

var (
	interpolationTokenRe = regexp.MustCompile(`^\$\((\S+)\)$`)
)

// findAndInterpolateAllKeys takes an `input` any value, and recursively
// identifies any values that should be replaced from `state`.
//
// A value `v` should be replaced if it is wrapped as follows: `$(v)`.
func findAndInterpolateAllKeys(input any, state store.WorkflowExecution) (any, error) {
	return deepMap(
		input,
		func(el string) (any, error) {
			matches := interpolationTokenRe.FindStringSubmatch(el)
			if len(matches) < 2 {
				return el, nil
			}

			interpolatedVar := matches[1]
			return interpolateKey(interpolatedVar, state)
		},
	)
}

// findRefs takes an `inputs` map and returns a list of all the step references
// contained within it.
func findRefs(inputs map[string]any) ([]string, error) {
	refs := []string{}
	_, err := deepMap(
		inputs,
		// This function is called for each string in the map
		// for each string, we iterate over each match of the interpolation token
		// - if there are no matches, return no reference
		// - if there is one match, return the reference
		// - if there are multiple matches (in the case of a multi-part state reference), return just the step ref
		func(el string) (any, error) {
			matches := interpolationTokenRe.FindStringSubmatch(el)
			if len(matches) < 2 {
				return el, nil
			}

			m := matches[1]
			parts := strings.Split(m, ".")
			if len(parts) < 1 {
				return nil, fmt.Errorf("invalid ref %s", m)
			}

			refs = append(refs, parts[0])
			return el, nil
		},
	)
	return refs, err
}

// deepMap recursively applies a transformation function
// over each string within:
//
//   - a map[string]any
//   - a []any
//   - a string
func deepMap(input any, transform func(el string) (any, error)) (any, error) {
	// in the case of a string, simply apply the transformation
	// in the case of a map, recurse and apply the transformation to each value
	// in the case of a list, recurse and apply the transformation to each element
	switch tv := input.(type) {
	case string:
		nv, err := transform(tv)
		if err != nil {
			return nil, err
		}

		return nv, nil
	case mapping:
		// coerce mapping to map[string]any
		mp := map[string]any(tv)

		nm := map[string]any{}
		for k, v := range mp {
			nv, err := deepMap(v, transform)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case map[string]any:
		nm := map[string]any{}
		for k, v := range tv {
			nv, err := deepMap(v, transform)
			if err != nil {
				return nil, err
			}

			nm[k] = nv
		}
		return nm, nil
	case []any:
		a := []any{}
		for _, el := range tv {
			ne, err := deepMap(el, transform)
			if err != nil {
				return nil, err
			}

			a = append(a, ne)
		}
		return a, nil
	}

	return nil, fmt.Errorf("cannot traverse item %+v of type %T", input, input)
}
