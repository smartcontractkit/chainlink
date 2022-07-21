package pipeline

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type event struct {
	expected string
	result   Result
}

func TestScheduler(t *testing.T) {
	// NOTE: task type does not matter in the test cases, it's just there so it's parsed successfully
	tests := []struct {
		name      string
		spec      string
		events    []event
		assertion func(t *testing.T, p Pipeline, results map[int]TaskRunResult)
	}{
		{
			name: "fail early immediately cancels subsequent tasks",
			spec: `
			a [type=median failEarly=true]
			b [type=median index=0]
			a -> b`,
			events: []event{
				{
					expected: "a",
					result:   Result{Error: ErrTaskRunFailed},
				},
				// no further events for `b`
			},
			assertion: func(t *testing.T, p Pipeline, results map[int]TaskRunResult) {
				result := results[p.ByDotID("b").ID()]
				// b is marked as cancelled
				require.Equal(t, uint(0), result.Attempts)
				require.Equal(t, ErrCancelled, result.Result.Error)
			},
		},
		{
			name: "retry: try task N times, then fail it",
			spec: `
			a [type=median retries=3 minBackoff="1us" maxBackoff="1us"]
			b [type=median index=0]
			a -> b`,
			events: []event{
				{
					expected: "a",
					result:   Result{Error: ErrTaskRunFailed},
				},
				{
					expected: "a",
					result:   Result{Error: ErrTaskRunFailed},
				},
				{
					expected: "a",
					result:   Result{Error: ErrTimeout},
				},
				{
					expected: "b",
					result:   Result{Value: 1},
				},
			},
			assertion: func(t *testing.T, p Pipeline, results map[int]TaskRunResult) {
				result := results[p.ByDotID("a").ID()]
				// a is marked as errored with the last error in sequence
				require.Equal(t, uint(3), result.Attempts)
				require.Equal(t, ErrTimeout, result.Result.Error)
			},
		},
		{
			name: "retry task: proceed when it succeeds",
			spec: `
			a [type=median retries=3 minBackoff="1us" maxBackoff="1us"]
			b [type=median index=0]
			a -> b`,
			events: []event{
				{
					expected: "a",
					result:   Result{Error: ErrTaskRunFailed},
				},
				{
					expected: "a",
					result:   Result{Value: 1},
				},
				{
					expected: "b",
					result:   Result{Value: 1},
				},
			},
			assertion: func(t *testing.T, p Pipeline, results map[int]TaskRunResult) {
				result := results[p.ByDotID("a").ID()]
				// a has no errors
				require.Equal(t, nil, result.Result.Error)
				require.Equal(t, 1, result.Result.Value)
				require.Equal(t, uint(2), result.Attempts)
			},
		},
		{
			name: "retry task + failEarly: cancel pending retries",
			spec: `
			a [type=median retries=3 minBackoff="10ms" maxBackoff="10ms" index=0]
			b [type=median failEarly=true index=1]
			`,
			events: []event{
				{
					expected: "a",
					result:   Result{Error: ErrTaskRunFailed},
				},
				{
					expected: "b",
					result:   Result{Error: ErrTaskRunFailed},
				},
				// now `b` failing early should stop any retries on `a`
			},
			assertion: func(t *testing.T, p Pipeline, results map[int]TaskRunResult) {
				result := results[p.ByDotID("a").ID()]
				// a only has a single attempt and it got cancelled
				require.Equal(t, uint(1), result.Attempts)
				require.Equal(t, ErrCancelled, result.Result.Error)
			},
		},
	}

	for _, test := range tests {
		p, err := Parse(test.spec)
		require.NoError(t, err)
		vars := NewVarsFrom(nil)
		run := NewRun(Spec{}, vars)
		s := newScheduler(p, &run, vars, logger.TestLogger(t))

		go s.Run()

		for _, event := range test.events {
			select {
			case taskRun := <-s.taskCh:
				require.Equal(t, event.expected, taskRun.task.DotID())
				now := time.Now()
				s.report(testutils.Context(t), TaskRunResult{
					ID:         uuid.NewV4(),
					Task:       taskRun.task,
					Result:     event.result,
					FinishedAt: null.TimeFrom(now),
					CreatedAt:  now,
				})
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for task run")
			}
		}

		select {
		case _, ok := <-s.taskCh:
			// channel is now closed, if it's not that means there's more tasks
			require.Falsef(t, ok, "scheduler has more tasks to schedule")
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for scheduler to halt")
		}

		test.assertion(t, *p, s.results)

	}
}
