package pipeline

import (
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestLogTaskRunToPrometheus(t *testing.T) {
	t.Parallel()

	p, err := Parse(`
memo [type=memo value="1"];
b [type=bridge name=metricsbridge];
memo -> b;
`)
	require.NoError(t, err)
	// A job ID no other test uses, since the metric vectors are global.
	spec := Spec{JobID: 987654321, JobName: "task-metrics"}
	memo, bridge := p.Tasks[0], p.Tasks[1]
	require.Equal(t, TaskTypeMemo, memo.Type())
	require.Equal(t, TaskTypeBridge, bridge.Type())

	start := time.Now()
	run := func(task Task, elapsed time.Duration, err error) {
		logTaskRunToPrometheus(TaskRunResult{
			Task:       task,
			Result:     Result{Error: err},
			CreatedAt:  start,
			FinishedAt: null.TimeFrom(start.Add(elapsed)),
		}, spec)
	}
	run(memo, time.Millisecond, nil)
	run(memo, 2*time.Millisecond, nil)
	run(memo, 3*time.Millisecond, errors.New("boom"))
	run(bridge, 4*time.Millisecond, nil)

	finished := func(dotID string, taskType TaskType, bridgeName, status string) float64 {
		return testutil.ToFloat64(PromPipelineTasksTotalFinished.WithLabelValues("987654321", "task-metrics", dotID, string(taskType), bridgeName, status))
	}
	assert.InDelta(t, 2, finished("memo", TaskTypeMemo, "", "completed"), 0)
	assert.InDelta(t, 1, finished("memo", TaskTypeMemo, "", "error"), 0)
	assert.InDelta(t, 1, finished("b", TaskTypeBridge, "metricsbridge", "completed"), 0)
	assert.InDelta(t, float64(3*time.Millisecond), testutil.ToFloat64(PromPipelineTaskExecutionTime.WithLabelValues("987654321", "task-metrics", "memo", string(TaskTypeMemo))), 0)

	// Handles are looked up once per task and status, then reused.
	first, second := taskMetricsFor(spec, memo, "completed"), taskMetricsFor(spec, memo, "completed")
	assert.Same(t, first, second)
	assert.NotSame(t, taskMetricsFor(spec, memo, "completed"), taskMetricsFor(spec, memo, "error"))
}
