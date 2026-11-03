package pipeline

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
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
			name: "staleness fallback with no fault budget left keeps every input",
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
			name: "staleness fallback excludes only the stalest inputs the threshold allows (freshest is never stale)",
			timestamps: map[string]uint64{
				"ds1": nowMs - 120_000,
				"ds2": nowMs - 240_000,
				"ds3": nowMs - 360_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "30",
			want:                 "1.5",
		},
		{
			name: "staleness fallback excludes equally stale inputs earliest in values first",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 120_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "1",
			stalenessGateSeconds: "30",
			want:                 "2",
		},
		{
			// Accepted risk (docs/adr/0001-median-staleness-gate.md): the
			// reference is the freshest input, so a single source reporting a
			// future timestamp makes every other input stale and, with the
			// default allowedFaults of n-1, alone becomes the median.
			name: "a single future-timestamped source becomes the reference and the median",
			timestamps: map[string]uint64{
				"ds1": nowMs + 3_600_000,
				"ds2": nowMs,
				"ds3": nowMs - 10_000,
			},
			allowedFaults:        "",
			stalenessGateSeconds: "30",
			want:                 "1",
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
			name: "absurdly large gate seconds never overflow into a small gate",
			timestamps: map[string]uint64{
				"ds1": nowMs,
				"ds2": nowMs - 10_000,
				"ds3": nowMs - 120_000,
			},
			allowedFaults:        "0",
			stalenessGateSeconds: "18446744073709551615",
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

func TestMedianTask_StalenessGate_ResultTimestamp(t *testing.T) {
	t.Parallel()

	const nowMs = uint64(1_700_000_000_000)

	tests := []struct {
		name                 string
		stalenessGateSeconds string
		want                 uint64
	}{
		{"excluded inputs do not make the result look stale", "30", nowMs - 10_000},
		{"without the gate every input counts", "", nowMs - 120_000},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vars := NewVarsFrom(map[string]any{
				"ds1": mustDecimalStaleness(t, "1"),
				"ds2": mustDecimalStaleness(t, "2"),
				"ds3": mustDecimalStaleness(t, "3"),
				"ds4": mustDecimalStaleness(t, "4"),
			})
			vars.SetTimestamp("ds1", nowMs)
			vars.SetTimestamp("ds2", nowMs-10_000)
			vars.SetTimestamp("ds3", nowMs-120_000)
			// ds4 is unknown and never lowers the result's timestamp
			task := MedianTask{
				BaseTask:             NewBaseTask(0, "task", nil, nil, 0),
				Values:               "[$(ds1),$(ds2),$(ds3),$(ds4)]",
				StalenessGateSeconds: test.stalenessGateSeconds,
			}
			result := runMedianWithTimestamps(t, task, vars, nil)
			require.NoError(t, result.Error)
			require.Equal(t, test.want, result.Timestamp)
		})
	}
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

	t.Run("literal values never take timestamps from inputs", func(t *testing.T) {
		inputs := []Result{{Timestamp: nowMs}, {Timestamp: nowMs - 120_000}}
		values := SliceParam{mustDecimalStaleness(t, "1"), mustDecimalStaleness(t, "2")}
		ts := resolveTimestamps("[1,2]", NewVarsFrom(nil), values, inputs)
		require.Nil(t, ts)
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

func TestScheduler_TimestampPropagation_ImplicitEdge(t *testing.T) {
	t.Parallel()

	// b consumes a through a var reference, which creates an implicit edge
	// (PropagateResult=false) — timestamps must still flow.
	spec := `
	a [type=median]
	b [type=median values="$(a)" index=0]`
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

	select {
	case taskRun := <-s.taskCh:
		require.Equal(t, "a", taskRun.task.DotID())
		report(taskRun, Result{Value: 1, Timestamp: 1000})
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for task run")
	}

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

	require.Equal(t, uint64(1000), s.results[p.ByDotID("b").ID()].Result.Timestamp)
	require.Equal(t, uint64(1000), s.vars.GetTimestamp("b"))
}

func TestScheduler_TimestampPropagation_OldestDependency(t *testing.T) {
	t.Parallel()

	spec := `
	a [type=median]
	b [type=median]
	c [type=median index=0]
	a -> c
	b -> c`

	tests := []struct {
		name    string
		results map[string]Result
		want    uint64
	}{
		{
			name: "derived value is only as fresh as its stalest input",
			results: map[string]Result{
				"a": {Value: 1, Timestamp: 2000},
				"b": {Value: 2, Timestamp: 1000},
				"c": {Value: 3},
			},
			want: 1000,
		},
		{
			name: "unknown dependency timestamps are ignored",
			results: map[string]Result{
				"a": {Value: 1, Timestamp: 2000},
				"b": {Value: 2},
				"c": {Value: 3},
			},
			want: 2000,
		},
		{
			name: "errored dependencies contribute no timestamp",
			results: map[string]Result{
				"a": {Value: 1, Timestamp: 2000},
				"b": {Error: errors.New("boom"), Timestamp: 1000},
				"c": {Value: 3},
			},
			want: 2000,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := Parse(spec)
			require.NoError(t, err)
			vars := NewVarsFrom(nil)
			s := newScheduler(p, NewRun(Spec{}, vars), vars, logger.TestLogger(t))

			go s.Run()

			for {
				select {
				case taskRun, ok := <-s.taskCh:
					if !ok {
						require.Equal(t, test.want, s.results[p.ByDotID("c").ID()].Result.Timestamp)
						return
					}
					now := time.Now()
					s.report(t.Context(), TaskRunResult{
						ID:         uuid.New(),
						Task:       taskRun.task,
						Result:     test.results[taskRun.task.DotID()],
						FinishedAt: null.TimeFrom(now),
						CreatedAt:  now,
					})
				case <-time.After(time.Second):
					t.Fatal("timed out waiting for scheduler")
				}
			}
		})
	}
}

func TestScheduler_ReconstructTimestamps(t *testing.T) {
	t.Parallel()

	// A resumed run: the bridge completed pre-suspension (its output is the
	// raw response body), the median downstream of it completed too. Both
	// timestamps must be recovered on reconstruction.
	spec := `
	ds1 [type=bridge name=foo]
	m  [type=median values="$(ds1.result)" index=0]`
	p, err := Parse(spec)
	require.NoError(t, err)

	run := NewRun(Spec{}, NewVarsFrom(nil))
	run.PipelineTaskRuns = []TaskRun{
		{
			DotID:  "ds1",
			Type:   TaskTypeBridge,
			Output: jsonserializable.JSONSerializable{Val: `{"data":{"mid":1},"timestamps":{"providerIndicatedTimeUnixMs":1000}}`, Valid: true},
		},
		{
			DotID:  "m",
			Type:   TaskTypeMedian,
			Output: jsonserializable.JSONSerializable{Val: "1", Valid: true},
		},
	}

	vars := NewVarsFrom(nil)
	s := newScheduler(p, run, vars, logger.TestLogger(t))

	require.Equal(t, uint64(1000), s.results[p.ByDotID("ds1").ID()].Result.Timestamp)
	require.Equal(t, uint64(1000), s.results[p.ByDotID("m").ID()].Result.Timestamp)
	require.Equal(t, uint64(1000), s.vars.GetTimestamp("ds1"))
	require.Equal(t, uint64(1000), s.vars.GetTimestamp("m"))
}
