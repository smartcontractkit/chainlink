package cmd_test

import (
	"bytes"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestJobPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id              = "1"
		name            = "Job 1"
		jobSpecType     = "fluxmonitor"
		schemaVersion   = uint32(1)
		maxTaskDuration = models.Interval(1 * time.Second)

		createdAt = time.Now()
		updatedAt = time.Now().Add(time.Second)
		buffer    = bytes.NewBufferString("")
		r         = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.JobPresenter{
		JobResource: presenters.JobResource{
			JAID:              presenters.NewJAID(id),
			Name:              name,
			Type:              presenters.JobSpecType(jobSpecType),
			SchemaVersion:     schemaVersion,
			MaxTaskDuration:   maxTaskDuration,
			DirectRequestSpec: nil,
			FluxMonitorSpec: &presenters.FluxMonitorSpec{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			OffChainReportingSpec: nil,
			KeeperSpec:            nil,
			PipelineSpec: presenters.PipelineSpec{
				ID:           1,
				DotDAGSource: "ds1 [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
			},
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, name)
	assert.Contains(t, output, jobSpecType)
	assert.Contains(t, output, "ds1 http")
	assert.Contains(t, output, "ds1_parse jsonparse")
	assert.Contains(t, output, "ds1_multiply multiply")
	assert.Contains(t, output, createdAt.Format(time.RFC3339))

	// Render many resources
	buffer.Reset()
	ps := cmd.JobPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, name)
	assert.Contains(t, output, jobSpecType)
	assert.Contains(t, output, "ds1 http")
	assert.Contains(t, output, "ds1_parse jsonparse")
	assert.Contains(t, output, "ds1_multiply multiply")
	assert.Contains(t, output, createdAt.Format(time.RFC3339))
}

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
			"gets the cron spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.CronJobSpec,
					CronSpec: &presenters.CronSpec{
						CreatedAt: now,
					},
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the vrf spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.VRFJobSpec,
					VRFSpec: &presenters.VRFSpec{
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
			"gets the blockhash store spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.BlockhashStoreJobSpec,
					BlockhashStoreSpec: &presenters.BlockhashStoreSpec{
						CreatedAt: now,
					},
				},
			},
			now.Format(time.RFC3339),
		},
		{
			"gets the blockheaderfeeder spec created at timestamp",
			&cmd.JobPresenter{
				JobResource: presenters.JobResource{
					Type: presenters.BlockHeaderFeederJobSpec,
					BlockHeaderFeederSpec: &presenters.BlockHeaderFeederSpec{
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

func TestClient_ListFindJobs(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
	})
	client, r := app.NewClientAndRenderer()

	// Create the job
	fs := flag.NewFlagSet("", flag.ExitOnError)
	cltest.FlagSetApplyFromAction(client.CreateJob, fs, "")

	require.NoError(t, fs.Parse([]string{"../testdata/tomlspecs/direct-request-spec.toml"}))

	err := client.CreateJob(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
	require.Len(t, r.Renders, 1)
	createOutput, ok := r.Renders[0].(*cmd.JobPresenter)
	require.True(t, ok, "Expected Renders[0] to be *cmd.JobPresenter, got %T", r.Renders[0])

	require.Nil(t, client.ListJobs(cltest.EmptyCLIContext()))
	jobs := *r.Renders[1].(*cmd.JobPresenters)
	require.Equal(t, 1, len(jobs))
	assert.Equal(t, createOutput.ID, jobs[0].ID)
}

func TestClient_ShowJob(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
	})
	client, r := app.NewClientAndRenderer()

	// Create the job
	fs := flag.NewFlagSet("", flag.ExitOnError)
	cltest.FlagSetApplyFromAction(client.CreateJob, fs, "")

	require.NoError(t, fs.Parse([]string{"../testdata/tomlspecs/direct-request-spec.toml"}))

	err := client.CreateJob(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
	require.Len(t, r.Renders, 1)
	createOutput, ok := r.Renders[0].(*cmd.JobPresenter)
	require.True(t, ok, "Expected Renders[0] to be *cmd.JobPresenter, got %T", r.Renders[0])

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{createOutput.ID})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.ShowJob(c))
	job := *r.Renders[0].(*cmd.JobPresenter)
	assert.Equal(t, createOutput.ID, job.ID)
}

func TestClient_CreateJobV2(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(100 * time.Millisecond)
		c.OCR.Enabled = ptr(true)
		c.P2P.V1.Enabled = ptr(true)
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	}, func(opts *startOptions) {
		opts.FlagsAndDeps = append(opts.FlagsAndDeps, cltest.DefaultP2PKey)
	})
	client, r := app.NewClientAndRenderer()

	requireJobsCount(t, app.JobORM(), 0)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	cltest.FlagSetApplyFromAction(client.CreateJob, fs, "")

	require.NoError(t, fs.Parse([]string{"../testdata/tomlspecs/ocr-bootstrap-spec.toml"}))

	err := client.CreateJob(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	requireJobsCount(t, app.JobORM(), 1)

	output := *r.Renders[0].(*cmd.JobPresenter)
	assert.Equal(t, presenters.JobSpecType("offchainreporting"), output.Type)
	assert.Equal(t, uint32(1), output.SchemaVersion)
	assert.Equal(t, "0x27548a32b9aD5D64c5945EaE9Da5337bc3169D15", output.OffChainReportingSpec.ContractAddress.String())
}

func TestClient_DeleteJob(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})
	client, r := app.NewClientAndRenderer()

	// Create the job
	fs := flag.NewFlagSet("", flag.ExitOnError)
	cltest.FlagSetApplyFromAction(client.CreateJob, fs, "")

	require.NoError(t, fs.Parse([]string{"../testdata/tomlspecs/direct-request-spec.toml"}))

	err := client.CreateJob(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
	require.NotEmpty(t, r.Renders)

	output := *r.Renders[0].(*cmd.JobPresenter)

	requireJobsCount(t, app.JobORM(), 1)

	jobs, _, err := app.JobORM().FindJobs(0, 1000)
	require.NoError(t, err)
	jobID := jobs[0].ID
	cltest.AwaitJobActive(t, app.JobSpawner(), jobID, 3*time.Second)

	// Must supply job id
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteJob, set, "")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, "must pass the job id to be archived", client.DeleteJob(c).Error())

	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteJob, set, "")

	require.NoError(t, set.Parse([]string{output.ID}))

	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.DeleteJob(c))

	requireJobsCount(t, app.JobORM(), 0)
}

func requireJobsCount(t *testing.T, orm job.ORM, expected int) {
	jobs, _, err := orm.FindJobs(0, 1000)
	require.NoError(t, err)
	require.Len(t, jobs, expected)
}
