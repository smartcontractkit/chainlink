package workflows

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

func TestInterpolateKey(t *testing.T) {
	t.Parallel()
	val, err := values.NewMap(
		map[string]any{
			"reports": map[string]any{
				"inner": "key",
			},
			"reportsList": []any{
				"listElement",
			},
		},
	)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		key      string
		state    store.WorkflowExecution
		expected any
		errMsg   string
	}{
		{
			name: "digging into a string",
			key:  "evm_median.outputs.reports",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: values.NewString("<a report>"),
						},
					},
				},
			},
			errMsg: "could not interpolate ref part `reports` (ref: `evm_median.outputs.reports`) in `<a report>`",
		},
		{
			name: "ref doesn't exist",
			key:  "evm_median.outputs.reports",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{},
			},
			errMsg: "could not find ref `evm_median`",
		},
		{
			name: "less than 2 parts",
			key:  "evm_median",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{},
			},
			errMsg: "must have at least two parts",
		},
		{
			name: "second part isn't `inputs` or `outputs`",
			key:  "evm_median.foo",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: values.NewString("<a report>"),
						},
					},
				},
			},
			errMsg: "second part must be `inputs` or `outputs`",
		},
		{
			name: "outputs has errored",
			key:  "evm_median.outputs",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Err: errors.New("catastrophic error"),
						},
					},
				},
			},
			errMsg: "step has errored",
		},
		{
			name: "digging into a recursive map",
			key:  "evm_median.outputs.reports.inner",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			expected: "key",
		},
		{
			name: "missing key in map",
			key:  "evm_median.outputs.reports.missing",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			errMsg: "could not find ref part `missing` (ref: `evm_median.outputs.reports.missing`) in",
		},
		{
			name: "digging into an array",
			key:  "evm_median.outputs.reportsList.0",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			expected: "listElement",
		},
		{
			name: "digging into an array that's too small",
			key:  "evm_median.outputs.reportsList.2",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			errMsg: "index out of bounds 2",
		},
		{
			name: "digging into an array with a string key",
			key:  "evm_median.outputs.reportsList.notAString",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			errMsg: "could not interpolate ref part `notAString` (ref: `evm_median.outputs.reportsList.notAString`) in `[listElement]`: `notAString` is not convertible to an int",
		},
		{
			name: "digging into an array with a negative index",
			key:  "evm_median.outputs.reportsList.-1",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			errMsg: "could not interpolate ref part `-1` (ref: `evm_median.outputs.reportsList.-1`) in `[listElement]`: index out of bounds -1",
		},
		{
			name: "empty element",
			key:  "evm_median.outputs..notAString",
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: val,
						},
					},
				},
			},
			errMsg: "could not find ref part `` (ref: `evm_median.outputs..notAString`) in",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			got, err := interpolateKey(tc.key, tc.state)
			if tc.errMsg != "" {
				require.ErrorContains(st, err, tc.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestInterpolateInputsFromState(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		inputs   map[string]any
		state    store.WorkflowExecution
		expected any
		errMsg   string
	}{
		{
			name: "substituting with a variable that exists",
			inputs: map[string]any{
				"shouldnotinterpolate": map[string]any{
					"shouldinterpolate": "$(evm_median.outputs)",
				},
			},
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: values.NewString("<a report>"),
						},
					},
				},
			},
			expected: map[string]any{
				"shouldnotinterpolate": map[string]any{
					"shouldinterpolate": "<a report>",
				},
			},
		},
		{
			name: "no substitution required",
			inputs: map[string]any{
				"foo": "bar",
			},
			state: store.WorkflowExecution{
				Steps: map[string]*store.WorkflowExecutionStep{
					"evm_median": {
						Outputs: &store.StepOutput{
							Value: values.NewString("<a report>"),
						},
					},
				},
			},
			expected: map[string]any{
				"foo": "bar",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			got, err := findAndInterpolateAllKeys(tc.inputs, tc.state)
			if tc.errMsg != "" {
				require.ErrorContains(st, err, tc.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
