package workflows

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
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

			Outputs: store.StepOutput{
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

// findAndInterpolateAllKeys takes an `input` any value, and recursively
// identifies any values that should be replaced from `state`.
//
// A value `v` should be replaced if it is wrapped as follows: `$(v)`.
func findAndInterpolateAllKeys(input any, state store.WorkflowExecution) (any, error) {
	return workflows.DeepMap(
		input,
		func(el string) (any, error) {
			matches := workflows.InterpolationTokenRe.FindStringSubmatch(el)
			if len(matches) < 2 {
				return el, nil
			}

			interpolatedVar := matches[1]
			return interpolateKey(interpolatedVar, state)
		},
	)
}
