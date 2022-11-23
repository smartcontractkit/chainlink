package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestMergeTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		left              string
		right             string
		vars              pipeline.Vars
		inputs            []pipeline.Result
		wantData          interface{}
		wantError         bool
		wantErrorContains string
	}{
		{
			"implicit left explicit right",
			"",
			`{"foo": "42", "bar": null, "blobber": false}`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"foo": "baz", "qux": 99, "flibber": null, "baz": true}`}},
			map[string]interface{}{
				"foo":     "42",
				"qux":     float64(99),
				"bar":     nil,
				"flibber": nil,
				"baz":     true,
				"blobber": false,
			},
			false,
			"",
		},
		{
			"explicit left explicit right",
			`{"foo": "baz", "qux": 99, "flibber": null}`,
			`{"foo": 42, "qux": null}`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"ignored": true}`}},
			map[string]interface{}{
				"foo":     int64(42),
				"qux":     nil,
				"flibber": nil,
			},
			false,
			"",
		},
		{
			"directions reversed",
			`{"foo": 42, "bar": null}`,
			`{"foo": "baz", "qux": 99, "flibber": null}`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"ignored": true}`}},
			map[string]interface{}{
				"foo":     "baz",
				"qux":     int64(99),
				"bar":     nil,
				"flibber": nil,
			},
			false,
			"",
		},
		{
			"invalid implicit left explicit right",
			``,
			`{"foo": 42, "bar": null}`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `"not a map"`}},
			nil,
			true,
			"left-side: json: cannot unmarshal string",
		},
		{
			"implicit left invalid explicit right",
			"",
			`not a map`,
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"foo": "baz", "qux": 99, "flibber": null, "baz": true}`}},
			nil,
			true,
			`right-side`,
		},
		{
			"explicit left variable data on right",
			`{"foo": 42, "bar": null}`,
			"$(someInput)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"someInput": map[string]interface{}{
					"foo":     "baz",
					"qux":     99,
					"flibber": nil,
				},
			}),
			[]pipeline.Result{},
			map[string]interface{}{
				"foo":     "baz",
				"qux":     99,
				"bar":     nil,
				"flibber": nil,
			},
			false,
			"",
		},
		{
			"explicit left invalid variable data on right",
			`{"foo": 42, "bar": null}`,
			"$(someInput)",
			pipeline.NewVarsFrom(map[string]interface{}{
				"someInput": "this is a string",
			}),
			[]pipeline.Result{},
			nil,
			true,
			`right-side`,
		},
		{
			"variable in left",
			`{"foo": 42, "bar": null}`,
			`{"flibber": $(someInput)}`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"someInput": "this is a string",
			}),
			[]pipeline.Result{},
			map[string]interface{}{
				"foo":     int64(42),
				"bar":     nil,
				"flibber": "this is a string",
			},
			false,
			"",
		},
		{
			"variable in right",
			`{"flibber": $(someInput)}`,
			`{"foo": 42, "bar": null}`,
			pipeline.NewVarsFrom(map[string]interface{}{
				"someInput": "this is a string",
			}),
			[]pipeline.Result{},
			map[string]interface{}{
				"foo":     int64(42),
				"bar":     nil,
				"flibber": "this is a string",
			},
			false,
			"",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.MergeTask{
				BaseTask: pipeline.NewBaseTask(0, "merge", nil, nil, 0),
				Left:     test.left,
				Right:    test.right,
			}
			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.wantError {
				if test.wantErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.wantErrorContains)
				}

				require.Nil(t, result.Value)
				val, err := test.vars.Get("merge")
				require.Equal(t, pipeline.ErrKeypathNotFound, errors.Cause(err))
				require.Nil(t, val)
			} else {
				assert.NoError(t, result.Error)
				assert.Equal(t, test.wantData, result.Value)
			}
		})
	}
}
