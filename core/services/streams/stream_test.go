package streams

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var UUID = uuid.New()

type mockRunner struct {
	p    *pipeline.Pipeline
	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

func (m *mockRunner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	return m.run, m.trrs, m.err
}
func (m *mockRunner) InitializePipeline(spec pipeline.Spec) (p *pipeline.Pipeline, err error) {
	return m.p, m.err
}
func (m *mockRunner) InsertFinishedRun(run *pipeline.Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	return m.err
}

type MockTask struct {
	result pipeline.Result
}

func (m *MockTask) Type() pipeline.TaskType { return "MockTask" }
func (m *MockTask) ID() int                 { return 0 }
func (m *MockTask) DotID() string           { return "" }
func (m *MockTask) Run(ctx context.Context, lggr logger.Logger, vars pipeline.Vars, inputs []pipeline.Result) (pipeline.Result, pipeline.RunInfo) {
	return m.result, pipeline.RunInfo{}
}
func (m *MockTask) Base() *pipeline.BaseTask           { return nil }
func (m *MockTask) Outputs() []pipeline.Task           { return nil }
func (m *MockTask) Inputs() []pipeline.TaskDependency  { return nil }
func (m *MockTask) OutputIndex() int32                 { return 0 }
func (m *MockTask) TaskTimeout() (time.Duration, bool) { return 0, false }
func (m *MockTask) TaskRetries() uint32                { return 0 }
func (m *MockTask) TaskMinBackoff() time.Duration      { return 0 }
func (m *MockTask) TaskMaxBackoff() time.Duration      { return 0 }

func Test_Stream(t *testing.T) {
	lggr := logger.TestLogger(t)
	runner := &mockRunner{}
	spec := pipeline.Spec{}
	id := StreamID(123)
	ctx := testutils.Context(t)

	t.Run("Run", func(t *testing.T) {
		strm := newStream(lggr, id, spec, runner, nil)

		t.Run("errors with empty pipeline", func(t *testing.T) {
			_, _, err := strm.Run(ctx)
			assert.EqualError(t, err, "Run failed: Run failed due to unparseable pipeline: empty pipeline")
		})

		spec.DotDagSource = `
succeed             [type=memo value=42]
succeed;
`

		strm = newStream(lggr, id, spec, runner, nil)

		t.Run("executes the pipeline (success)", func(t *testing.T) {
			runner.run = &pipeline.Run{ID: 42}
			runner.trrs = []pipeline.TaskRunResult{pipeline.TaskRunResult{ID: UUID}}
			runner.err = nil

			run, trrs, err := strm.Run(ctx)
			assert.NoError(t, err)

			assert.Equal(t, int64(42), run.ID)
			require.Len(t, trrs, 1)
			assert.Equal(t, UUID, trrs[0].ID)
		})
		t.Run("executes the pipeline (failure)", func(t *testing.T) {
			runner.err = errors.New("something exploded")

			_, _, err := strm.Run(ctx)
			require.Error(t, err)

			assert.EqualError(t, err, "Run failed: error executing run for spec ID 0: something exploded")
		})
	})
}

func Test_ExtractBigInt(t *testing.T) {
	t.Run("wrong number of inputs", func(t *testing.T) {
		trrs := []pipeline.TaskRunResult{}

		_, err := ExtractBigInt(trrs)
		assert.EqualError(t, err, "invalid number of results, expected: 1, got: 0")
	})
	t.Run("wrong type", func(t *testing.T) {
		trrs := []pipeline.TaskRunResult{
			{
				Result: pipeline.Result{Value: []byte{1, 2, 3}},
				Task:   &MockTask{},
			},
		}

		_, err := ExtractBigInt(trrs)
		assert.EqualError(t, err, "failed to parse BenchmarkPrice: type []uint8 cannot be converted to decimal.Decimal ([1 2 3])")
	})
	t.Run("correct inputs", func(t *testing.T) {
		trrs := []pipeline.TaskRunResult{
			{
				Result: pipeline.Result{Value: "122.345"},
				Task:   &MockTask{},
			},
		}

		val, err := ExtractBigInt(trrs)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(122), val)
	})
}
