package pipeline

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func mustDecimalStaleness(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	require.NoError(t, err)
	return d
}

func runMedianWithTimestamps(t *testing.T, task MedianTask, vars Vars, inputs []Result) Result {
	t.Helper()
	result, runInfo := task.Run(t.Context(), logger.TestLogger(t), vars, inputs)
	require.False(t, runInfo.IsPending)
	require.False(t, runInfo.IsRetryable)
	return result
}

func TestMedianTask_StalenessGate(t *testing.T) {
	t.Parallel()

	const nowMs = uint64(1_700_000_000_000)
	values := map[string]any{
		"ds1": mustDecimalStaleness(t, "1"),
		"ds2": mustDecimalStaleness(t, "2"),
		"ds3": mustDecimalStaleness(t, "3"),
	}

	tests := []struct {
		name                 string
		timestamps           map[string]uint64
		allowedFaults        string
		stalenessGateSeconds string
		want                 string // expected median as decimal string
	}{
		{
			name: "gate excludes input older than max staleness",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "30",
			want:                 "1.5",
		},
		{
			name: "staleness fallback when exclusions would exceed allowed faults",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "0",
			stalenessGateSeconds: "30",
			want:                 "2",
		},
		{
			name: "unset allowedFaults tolerates stale faults by default",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "",
			stalenessGateSeconds: "30",
			want:                 "1.5",
		},
		{
			name: "zero timestamp is never excluded",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": 0,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "30",
			want:                 "1.5",
		},
		{
			name: "all timestamps zero leaves the gate inert",
			timestamps: map[string]uint64{
				"ds1": 0,
				"ds2": 0,
				"ds3": 0,
			},
			allowedFaults:        "0",
			stalenessGateSeconds: "30",
			want:                 "2",
		},
		{
			name: "mostly stale triggers fallback when exclusions exceed threshold (freshest is never stale)",
			timestamps: map[string]uint64{
				"ds1": nowMs - 120_000,
				"ds2": nowMs - 240_000,
				"ds3": nowMs - 360_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "30",
			want:                 "2",
		},
		{
			name: "explicit zero means zero tolerance: only the freshest survives",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "2",
			stalenessGateSeconds: "0",
			want:                 "1",
		},
		{
			name: "zero tolerance still never excludes unknown timestamps",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": 0,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "0",
			want:                 "1.5",
		},
		{
			name: "gate off when the parameter is absent",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "",
			want:                 "2",
		},
		{
			name: "gate seconds can be a var expression",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "$(gateSeconds)",
			want:                 "1.5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := NewVarsFrom(map[string]any{
				"ds1":         values["ds1"],
				"ds2":         values["ds2"],
				"ds3":         values["ds3"],
				"gateSeconds": uint64(30),
			})
			for dotID, ts := range test.timestamps {
				vars.SetTimestamp(dotID, ts)
			}
			task := MedianTask{
				BaseTask:             NewBaseTask(0, "task", nil, nil, 0),
				Values:               "[$(ds1),$(ds2),$(ds3)]",
				AllowedFaults:        test.allowedFaults,
				StalenessGateSeconds: test.stalenessGateSeconds,
			}
			result := runMedianWithTimestamps(t, task, vars, nil)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Value.(decimal.Decimal).String())
		})
	}
}

func TestMedianTask_StalenessGate_InputsFallback(t *testing.T) {
	t.Parallel()

	nowMs := uint64(1_700_000_000_000)
	inputs := []Result{
		{Value: mustDecimalStaleness(t, "1"), Timestamp: nowMs},
		{Value: mustDecimalStaleness(t, "2"), Timestamp: nowMs - 10_000},
		{Value: mustDecimalStaleness(t, "3"), Timestamp: nowMs - 120_000},
	}
	task := MedianTask{
		BaseTask:             NewBaseTask(0, "task", nil, nil, 0),
		AllowedFaults:        "1",
		StalenessGateSeconds: "30",
	}
	result := runMedianWithTimestamps(t, task, NewVarsFrom(nil), inputs)
	require.NoError(t, result.Error)
	// ds3 is 120s stale and excluded; ds1 and ds2 survive.
	require.Equal(t, "1.5", result.Value.(decimal.Decimal).String())
}

func TestResolveTimestamps(t *testing.T) {
	t.Parallel()

	const nowMs = uint64(1_700_000_000_000)

	t.Run("single ref to a slice assigns its timestamp to every element", func(t *testing.T) {
		vars := NewVarsFrom(map[string]any{
			"ds1": []any{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2"), mustDecimalStaleness(t, "3")},
		})
		vars.SetTimestamp("ds1", nowMs)
		values := SliceParam{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2"), mustDecimalStaleness(t, "3")}
		ts := resolveTimestamps("$(ds1)", vars, values, nil)
		require.Equal(t, []uint64{nowMs, nowMs, nowMs}, ts)
	})

	t.Run("nested keypath resolves to the root dotID", func(t *testing.T) {
		vars := NewVarsFrom(map[string]any{
			"ds1": map[string]any{"result": mustDecimalStaleness(t, "1")},
			"ds2": map[string]any{"result": mustDecimalStaleness(t, "2")},
		})
		vars.SetTimestamp("ds1", nowMs)
		values := SliceParam{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2")}
		ts := resolveTimestamps("[$(ds1.result),$(ds2.result)]", vars, values, nil)
		require.Equal(t, []uint64{nowMs, 0}, ts)
	})

	t.Run("unmappable values disable the gate", func(t *testing.T) {
		vars := NewVarsFrom(map[string]any{"ds1": mustDecimalStaleness(t, "1")})
		values := SliceParam{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2")}
		ts := resolveTimestamps("[$(ds1)]", vars, values, nil)
		require.Nil(t, ts)
	})

	t.Run("no refs falls back to inputs 1:1", func(t *testing.T) {
		inputs := []Result{{Timestamp: nowMs}, {Timestamp: 0}}
		values := SliceParam{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2")}
		ts := resolveTimestamps("", NewVarsFrom(nil), values, inputs)
		require.Equal(t, []uint64{nowMs, 0}, ts)
	})
}

func TestBridgeResponseTimestamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want uint64
	}{
		{
			name: "provider-indicated preferred over provider-data-received",
			body: `{"data":{"mid":1},"timestamps":{"providerIndicatedTimeUnixMs":100,"providerDataReceivedUnixMs":200}}`,
			want: 100,
		},
		{
			name: "provider-data-received used when indicated is missing",
			body: `{"data":{"mid":1},"timestamps":{"providerDataReceivedUnixMs":200}}`,
			want: 200,
		},
		{
			name: "string timestamps parse",
			body: `{"timestamps":{"providerIndicatedTimeUnixMs":"300"}}`,
			want: 300,
		},
		{
			name: "missing timestamps",
			body: `{"data":{"mid":1}}`,
			want: 0,
		},
		{
			name: "garbage timestamps",
			body: `{"timestamps":{"providerIndicatedTimeUnixMs":"not-a-number"}}`,
			want: 0,
		},
		{
			name: "zero timestamps are unknown",
			body: `{"timestamps":{"providerIndicatedTimeUnixMs":0}}`,
			want: 0,
		},
		{
			name: "invalid body",
			body: `not json`,
			want: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, bridgeResponseTimestamp([]byte(test.body)))
		})
	}
}

func TestScheduler_TimestampPropagation(t *testing.T) {
	t.Parallel()

	spec := `
	a [type=median]
	b [type=median index=0]
	a -> b`
	p, err := Parse(spec)
	require.NoError(t, err)
	vars := NewVarsFrom(nil)
	run := NewRun(Spec{}, vars)
	s := newScheduler(p, run, vars, logger.TestLogger(t))

	go s.Run()

	report := func(taskRun *memoryTaskRun, result Result) {
		now := time.Now()
		s.report(t.Context(), TaskRunResult{
			ID:         uuid.New(),
			Task:       taskRun.task,
			Result:     result,
			FinishedAt: null.TimeFrom(now),
			CreatedAt:  now,
		})
	}

	// a reports with its own timestamp
	select {
	case taskRun := <-s.taskCh:
		require.Equal(t, "a", taskRun.task.DotID())
		report(taskRun, Result{Value: 1, Timestamp: 1000})
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for task run")
	}

	// b reports without a timestamp and must inherit a's
	select {
	case taskRun := <-s.taskCh:
		require.Equal(t, "b", taskRun.task.DotID())
		report(taskRun, Result{Value: 2})
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for task run")
	}

	select {
	case _, ok := <-s.taskCh:
		require.Falsef(t, ok, "scheduler has more tasks to schedule")
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for scheduler to halt")
	}

	require.Equal(t, uint64(1000), s.results[p.ByDotID("a").ID()].Result.Timestamp)
	require.Equal(t, uint64(1000), s.results[p.ByDotID("b").ID()].Result.Timestamp)
	// tasks receive copies of the scheduler's vars, which carry the registry
	require.Equal(t, uint64(1000), s.vars.GetTimestamp("b"))
}
