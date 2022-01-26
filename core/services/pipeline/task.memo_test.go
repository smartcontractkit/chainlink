package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMemoTask_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  interface{}
		output string
	}{
		{"identity", pipeline.ObjectParam{Type: pipeline.BoolType, BoolValue: true}, "true"},

		{"nil", nil, "null"},

		{"bool", true, "true"},
		{"bool false", false, "false"},

		{"integer", 17, `"17"`},
		{"negative integer", -19, `"-19"`},
		{"uint", uint(17), `"17"`},
		{"float", 17.3, `"17.3"`},
		{"negative float", -17.3, `"-17.3"`},

		{"string", "hello world", `"hello world"`},

		{"array", []int{17, 19}, "[17,19]"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			vars := pipeline.NewVarsFrom(nil)
			var value pipeline.ObjectParam
			err := value.UnmarshalPipelineParam(test.input)
			require.NoError(t, err)

			task := pipeline.MemoTask{BaseTask: pipeline.NewBaseTask(0, "task", nil, nil, 0)}
			result, _ := task.Run(context.Background(), logger.TestLogger(t), vars, []pipeline.Result{{Value: test.input}})
			require.NoError(t, result.Error)
			marshalledValue, err := result.Value.(pipeline.ObjectParam).Marshal()
			require.NoError(t, err)
			assert.Equal(t, test.output, marshalledValue)
		})
	}
}
