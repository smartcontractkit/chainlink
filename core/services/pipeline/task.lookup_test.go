package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
)

func Test_LookupTask(t *testing.T) {
	task := pipeline.LookupTask{}
	m := map[string]interface{}{
		"foo": 42,
		"bar": "baz",
	}
	var vars pipeline.Vars
	var inputs []pipeline.Result

	t.Run("with valid key for map", func(t *testing.T) {
		task.Key = "foo"
		inputs = []pipeline.Result{{Value: m, Error: nil}}

		res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), vars, inputs)

		assert.Equal(t, 42, res.Value)
		assert.Nil(t, res.Error)
	})
	t.Run("returns nil if key is missing", func(t *testing.T) {
		task.Key = "qux"
		inputs = []pipeline.Result{{Value: m, Error: nil}}

		res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), vars, inputs)

		assert.Nil(t, res.Error)
		assert.Nil(t, res.Value)
	})
	t.Run("errors when input is not a map", func(t *testing.T) {
		task.Key = "qux"
		inputs = []pipeline.Result{{Value: "something", Error: nil}}

		res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), vars, inputs)

		assert.EqualError(t, res.Error, "unexpected input type: string")
		assert.Nil(t, res.Value)
	})
	t.Run("errors when input is error", func(t *testing.T) {
		task.Key = "qux"
		inputs = []pipeline.Result{{Value: nil, Error: errors.New("something blew up")}}

		res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), vars, inputs)

		assert.EqualError(t, res.Error, "task inputs: too many errors")
		assert.Nil(t, res.Value)
	})
	t.Run("errors with too many inputs", func(t *testing.T) {
		task.Key = "qux"
		inputs = []pipeline.Result{{Value: m, Error: nil}, {Value: nil, Error: errors.New("something blew up")}}

		res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), vars, inputs)

		assert.EqualError(t, res.Error, "task inputs: min: 1 max: 1 (got 2): wrong number of task inputs")
		assert.Nil(t, res.Value)
	})
}
