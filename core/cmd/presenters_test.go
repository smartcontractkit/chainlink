package cmd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/cmd"
)

func TestJAID(t *testing.T) {
	t.Parallel()

	jaid := cmd.JAID{ID: "1"}

	t.Run("GetID", func(t *testing.T) { assert.Equal(t, "1", jaid.GetID()) })
	t.Run("SetID", func(t *testing.T) {
		jaid.SetID("2")
		assert.Equal(t, "2", jaid.GetID())
	})
}

func TestJobType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "directrequest", cmd.DirectRequestJob.String())
}

func TestJob_GetName(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{}

	assert.Equal(t, "specDBs", job.GetName())
}

func TestJob_GetTasks(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{}

	t.Run("gets the tasks from the DAG in reverse order", func(t *testing.T) {
		job.PipelineSpec = cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		}

		tasks, err := job.GetTasks()

		assert.NoError(t, err)
		assert.Equal(t, []string{
			"ds1 http",
			"ds1_parse jsonparse",
			"ds1_multiply multiply",
		}, tasks)
	})

	t.Run("parse error", func(t *testing.T) {
		job.PipelineSpec = cmd.PipelineSpec{
			DotDAGSource: "invalid dot",
		}

		tasks, err := job.GetTasks()

		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestJob_FriendlyTasks(t *testing.T) {
	t.Parallel()

	job := &cmd.Job{}

	t.Run("gets the tasks in a printable format", func(t *testing.T) {
		job.PipelineSpec = cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		}

		assert.Equal(t, []string{
			"ds1 http",
			"ds1_parse jsonparse",
			"ds1_multiply multiply",
		}, job.FriendlyTasks())
	})

	t.Run("parse error", func(t *testing.T) {
		job.PipelineSpec = cmd.PipelineSpec{
			DotDAGSource: "invalid dot",
		}

		assert.Equal(t, []string{"error parsing DAG"}, job.FriendlyTasks())
	})
}

func TestJob_FriendlyCreatedAt(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name   string
		job    *cmd.Job
		result string
	}{
		{
			"gets the direct request spec created at timestamp",
			&cmd.Job{
				Type: cmd.DirectRequestJob,
				DirectRequestSpec: &cmd.DirectRequestSpec{
					CreatedAt: now,
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the flux monitor spec created at timestamp",
			&cmd.Job{
				Type: cmd.FluxMonitorJob,
				FluxMonitorSpec: &cmd.FluxMonitorSpec{
					CreatedAt: now,
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the off chain reporting spec created at timestamp",
			&cmd.Job{
				Type: cmd.OffChainReportingJob,
				OffChainReportingSpec: &cmd.OffChainReportingSpec{
					CreatedAt: now,
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"invalid type",
			&cmd.Job{
				Type: "invalid type",
			},
			"unknown",
		},
		{
			"no spec exists",
			&cmd.Job{
				Type: cmd.DirectRequestJob,
			},
			"N/A",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.result, tc.job.FriendlyCreatedAt())
		})
	}
}

func TestJob_ToRow(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.Job{
		JAID: cmd.JAID{ID: "1"},
		Name: "Test Job",
		Type: cmd.DirectRequestJob,
		DirectRequestSpec: &cmd.DirectRequestSpec{
			CreatedAt: now,
		},
		PipelineSpec: cmd.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		},
	}

	assert.Equal(t, [][]string{
		{"1", "Test Job", "directrequest", "ds1 http", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_parse jsonparse", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_multiply multiply", now.Format(time.RFC3339)},
	}, job.ToRow())
}
