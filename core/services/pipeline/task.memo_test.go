package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMemoTask_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		output string
	}{
		{"nil", nil, "<nil>"},

		{"bool", true, "true"},
		{"bool false", false, "false"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			var value pipeline.ObjectParam
			err := value.UnmarshalPipelineParam(test.input)
			require.NoError(t, err)

			task := pipeline.MemoTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
			result := task.Run(context.Background(), vars, []pipeline.Result{{Value: test.input}})
			require.NoError(t, result.Error)
			require.Equal(t, test.output, result.Value.(pipeline.ObjectParam).String())
		})
	}

	// for _, test := range tests {
	// 	test := test
	// 	t.Run(test.name+" (with pipeline.Vars)", func(t *testing.T) {
	// 		vars := pipeline.NewVarsFrom(map[string]interface{}{
	// 			"foo":   map[string]interface{}{"bar": test.input},
	// 			"chain": map[string]interface{}{"link": test.times},
	// 		})
	// 		task := pipeline.MultiplyTask{
	// 			BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
	// 			Input:    "$(foo.bar)",
	// 			Times:    "$(chain.link)",
	// 		}
	// 		result := task.Run(context.Background(), vars, []pipeline.Result{})
	// 		require.NoError(t, result.Error)
	// 		require.Equal(t, test.want.String(), result.Value.(decimal.Decimal).String())
	// 	})
	// }
}

// func TestMultiplyTask_Unhappy(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name              string
// 		times             string
// 		input             string
// 		inputs            []pipeline.Result
// 		vars              pipeline.Vars
// 		wantErrorCause    error
// 		wantErrorContains string
// 	}{
// 		{"map as input from inputs", "100", "", []pipeline.Result{{Value: map[string]interface{}{"chain": "link"}}}, pipeline.NewVarsFrom(nil), pipeline.ErrBadInput, "input"},
// 		{"map as input from var", "100", "$(foo)", nil, pipeline.NewVarsFrom(map[string]interface{}{"foo": map[string]interface{}{"chain": "link"}}), pipeline.ErrBadInput, "input"},
// 		{"slice as input from inputs", "100", "", []pipeline.Result{{Value: []interface{}{"chain", "link"}}}, pipeline.NewVarsFrom(nil), pipeline.ErrBadInput, "input"},
// 		{"slice as input from var", "100", "$(foo)", nil, pipeline.NewVarsFrom(map[string]interface{}{"foo": []interface{}{"chain", "link"}}), pipeline.ErrBadInput, "input"},
// 		{"input as missing var", "100", "$(foo)", nil, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "input"},
// 		{"times as missing var", "$(foo)", "", []pipeline.Result{{Value: "123"}}, pipeline.NewVarsFrom(nil), pipeline.ErrKeypathNotFound, "times"},
// 	}

// 	for _, tt := range tests {
// 		test := tt
// 		t.Run(test.name, func(t *testing.T) {
// 			t.Parallel()

// 			task := pipeline.MultiplyTask{
// 				BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0),
// 				Input:    test.input,
// 				Times:    test.times,
// 			}
// 			result := task.Run(context.Background(), test.vars, test.inputs)
// 			require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
// 			if test.wantErrorContains != "" {
// 				require.Contains(t, result.Error.Error(), test.wantErrorContains)
// 			}
// 		})
// 	}
// }
