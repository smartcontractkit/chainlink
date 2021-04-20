package cmd_test

import (
	"flag"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestJobRenderer_GetTasks(t *testing.T) {
	t.Parallel()

	r := &cmd.JobPresenter{}

	t.Run("gets the tasks from the DAG in reverse order", func(t *testing.T) {
		r.PipelineSpec = presenters.PipelineSpec{
			DotDAGSource: "ds1 [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		}

		tasks, err := r.GetTasks()

		assert.NoError(t, err)
		assert.Equal(t, []string{
			"ds1 http",
			"ds1_parse jsonparse",
			"ds1_multiply multiply",
		}, tasks)
	})

	t.Run("parse error", func(t *testing.T) {
		r.PipelineSpec = presenters.PipelineSpec{
			DotDAGSource: "invalid dot",
		}

		tasks, err := r.GetTasks()

		assert.Error(t, err)
		assert.Nil(t, tasks)
	})
}

func TestJob_FriendlyTasks(t *testing.T) {
	t.Parallel()

	r := &cmd.JobPresenter{}

	t.Run("gets the tasks in a printable format", func(t *testing.T) {
		r.PipelineSpec = presenters.PipelineSpec{
			DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
		}

		assert.Equal(t, []string{
			"ds1 http",
			"ds1_parse jsonparse",
			"ds1_multiply multiply",
		}, r.FriendlyTasks())
	})

	t.Run("parse error", func(t *testing.T) {
		r.PipelineSpec = presenters.PipelineSpec{
			DotDAGSource: "invalid dot",
		}

		assert.Equal(t, []string{"error parsing DAG"}, r.FriendlyTasks())
	})
}

func TestJob_FriendlyCreatedAt(t *testing.T) {
	t.Parallel()

	now := time.Now()

	testCases := []struct {
		name   string
		job    *cmd.JobPresenter
		result string
	}{
		{
			"gets the direct request spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.DirectRequestJobSpec,
					DirectRequestSpec: &presenters.DirectRequestSpec{
						CreatedAt: now,
					},
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the flux monitor spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.FluxMonitorJobSpec,
					FluxMonitorSpec: &presenters.FluxMonitorSpec{
						CreatedAt: now,
					},
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the off chain reporting spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.OffChainReportingJobSpec,
					OffChainReportingSpec: &presenters.OffChainReportingSpec{
						CreatedAt: now,
					},
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"invalid type",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: "invalid type",
				},
			},
			"unknown",
		},
		{
			"no spec exists",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.DirectRequestJobSpec,
				},
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

func TestJob_ToRows(t *testing.T) {
	t.Parallel()

	now := time.Now()

	job := &cmd.JobPresenter{
		JAID: cmd.NewJAID("1"),
		JobResource: presenters.JobResource{
			Name: "Test Job",
			Type: presenters.DirectRequestJobSpec,
			DirectRequestSpec: &presenters.DirectRequestSpec{
				CreatedAt: now,
			},
			PipelineSpec: presenters.PipelineSpec{
				DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
			},
		},
	}

	assert.Equal(t, [][]string{
		{"1", "Test Job", "directrequest", "ds1 http", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_parse jsonparse", now.Format(time.RFC3339)},
		{"1", "Test Job", "directrequest", "ds1_multiply multiply", now.Format(time.RFC3339)},
	}, job.ToRows())

	// Produce a single row even if there is not DAG
	job.PipelineSpec.DotDAGSource = ""
	assert.Equal(t, [][]string{
		{"1", "Test Job", "directrequest", "", now.Format(time.RFC3339)},
	}, job.ToRows())
}

func TestClient_ListJobsV2(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	// Create the job
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Parse([]string{"../testdata/tomlspecs/direct-request-spec.toml"})
	err := client.CreateJobV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
	createOutput := *r.Renders[0].(*cmd.JobPresenter)

	require.Nil(t, client.ListJobsV2(cltest.EmptyCLIContext()))
	jobs := *r.Renders[1].(*[]cmd.JobPresenter)
	require.Equal(t, 1, len(jobs))
	assert.Equal(t, createOutput.ID, jobs[0].ID)
}

func TestClient_CreateJobV2(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	requireJobsCount(t, app.JobORM, 0)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Parse([]string{"../testdata/tomlspecs/ocr-bootstrap-spec.toml"})
	err := client.CreateJobV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	requireJobsCount(t, app.JobORM, 1)

	output := *r.Renders[0].(*cmd.JobPresenter)
	assert.Equal(t, presenters.JobSpecType("offchainreporting"), output.Type)
	assert.Equal(t, uint32(1), output.SchemaVersion)
	assert.Equal(t, "0x27548a32b9aD5D64c5945EaE9Da5337bc3169D15", output.OffChainReportingSpec.ContractAddress.String())
}

func TestClient_DeleteJobV2(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	// Create the job
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Parse([]string{"../testdata/tomlspecs/direct-request-spec.toml"})
	err := client.CreateJobV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	output := *r.Renders[0].(*cmd.JobPresenter)

	requireJobsCount(t, app.JobORM, 1)

	// Must supply job id
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must pass the job id to be archived", client.DeleteJobV2(c).Error())

	set = flag.NewFlagSet("test", 0)
	set.Parse([]string{output.ID})
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.DeleteJobV2(c))

	requireJobsCount(t, app.JobORM, 0)
}

func requireJobsCount(t *testing.T, orm job.ORM, expected int) {
	jobs, err := orm.JobsV2()
	require.NoError(t, err)
	require.Len(t, jobs, expected)
}
