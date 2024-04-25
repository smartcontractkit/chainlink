package job_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func clearJobsDb(t *testing.T, db *sqlx.DB) {
	cltest.ClearDBTables(t, db, "flux_monitor_round_stats_v2", "jobs", "pipeline_runs", "pipeline_specs", "pipeline_task_runs")
}

func TestPipelineORM_Integration(t *testing.T) {
	ctx := testutils.Context(t)
	const DotStr = `
        // data source 1
        ds1          [type=bridge name=voter_turnout];
        ds1_parse    [type=jsonparse path="one,two"];
        ds1_multiply [type=multiply times=1.23];

        // data source 2
        ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];
        ds2_parse    [type=jsonparse path="three.four" separator="."];
        ds2_multiply [type=multiply times=4.56];

        ds1 -> ds1_parse -> ds1_multiply -> answer1;
        ds2 -> ds2_parse -> ds2_multiply -> answer1;

        answer1 [type=median                      index=0];
        answer2 [type=bridge name=election_winner index=1];
    `

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = commonconfig.MustNewDuration(30 * time.Millisecond)
	})
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	ethKeyStore := keyStore.Eth()

	_, transmitterAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	require.NoError(t, keyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	var specID int32

	answer1 := &pipeline.MedianTask{
		AllowedFaults: "",
	}
	answer2 := &pipeline.BridgeTask{
		Name: "election_winner",
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times: "1.23",
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path: "one,two",
	}
	ds1 := &pipeline.BridgeTask{
		Name: "voter_turnout",
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times: "4.56",
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path: "three,four",
	}
	ds2 := &pipeline.HTTPTask{
		URL:         "https://chain.link/voter_turnout/USA-2020",
		Method:      "GET",
		RequestData: `{"hi": "hello"}`,
	}

	answer1.BaseTask = pipeline.NewBaseTask(
		6,
		"answer1",
		[]pipeline.TaskDependency{
			{PropagateResult: true, InputTask: pipeline.Task(ds1_multiply)},
			{PropagateResult: true, InputTask: pipeline.Task(ds2_multiply)}},
		nil,
		0)
	answer2.BaseTask = pipeline.NewBaseTask(7, "answer2", nil, nil, 1)
	ds1_multiply.BaseTask = pipeline.NewBaseTask(
		2,
		"ds1_multiply",
		[]pipeline.TaskDependency{{PropagateResult: true, InputTask: pipeline.Task(ds1_parse)}},
		[]pipeline.Task{answer1},
		0)
	ds2_multiply.BaseTask = pipeline.NewBaseTask(
		5,
		"ds2_multiply",
		[]pipeline.TaskDependency{{PropagateResult: true, InputTask: pipeline.Task(ds2_parse)}},
		[]pipeline.Task{answer1},
		0)
	ds1_parse.BaseTask = pipeline.NewBaseTask(
		1,
		"ds1_parse",
		[]pipeline.TaskDependency{{PropagateResult: true, InputTask: pipeline.Task(ds1)}},
		[]pipeline.Task{ds1_multiply},
		0)
	ds2_parse.BaseTask = pipeline.NewBaseTask(
		4,
		"ds2_parse",
		[]pipeline.TaskDependency{{PropagateResult: true, InputTask: pipeline.Task(ds2)}},
		[]pipeline.Task{ds2_multiply},
		0)
	ds1.BaseTask = pipeline.NewBaseTask(0, "ds1", nil, []pipeline.Task{ds1_parse}, 0)
	ds2.BaseTask = pipeline.NewBaseTask(3, "ds2", nil, []pipeline.Task{ds2_parse}, 0)
	expectedTasks := []pipeline.Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{})

	t.Run("creates task DAGs", func(t *testing.T) {
		ctx := testutils.Context(t)
		clearJobsDb(t, db)

		orm := pipeline.NewORM(db, logger.TestLogger(t), config.JobPipeline().MaxSuccessfulRuns())

		p, err := pipeline.Parse(DotStr)
		require.NoError(t, err)

		specID, err = orm.CreateSpec(ctx, *p, models.Interval(0))
		require.NoError(t, err)

		var pipelineSpecs []pipeline.Spec
		sql := `SELECT * FROM pipeline_specs;`
		err = db.Select(&pipelineSpecs, sql)
		require.NoError(t, err)
		require.Len(t, pipelineSpecs, 1)
		require.Equal(t, specID, pipelineSpecs[0].ID)
		require.Equal(t, DotStr, pipelineSpecs[0].DotDagSource)

		_, err = db.Exec(`DELETE FROM pipeline_specs`)
		require.NoError(t, err)
	})

	t.Run("creates runs", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		cfg := configtest.NewTestGeneralConfig(t)
		clearJobsDb(t, db)
		orm := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
		btORM := bridges.NewORM(db)
		relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{Client: evmtest.NewEthClientMockWithDefaultChain(t), DB: db, GeneralConfig: config, KeyStore: ethKeyStore})
		legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
		runner := pipeline.NewRunner(orm, btORM, config.JobPipeline(), cfg.WebServer(), legacyChains, nil, nil, lggr, nil, nil)

		jobORM := NewTestORM(t, db, orm, btORM, keyStore)

		dbSpec := makeVoterTurnoutOCRJobSpec(t, transmitterAddress, bridge.Name.String(), bridge2.Name.String())

		// Need a job in order to create a run
		require.NoError(t, jobORM.CreateJob(testutils.Context(t), dbSpec))

		var pipelineSpecs []pipeline.Spec
		sql := `SELECT pipeline_specs.*, job_pipeline_specs.job_id FROM pipeline_specs JOIN job_pipeline_specs ON (pipeline_specs.id = job_pipeline_specs.pipeline_spec_id);`
		require.NoError(t, db.Select(&pipelineSpecs, sql))
		require.Len(t, pipelineSpecs, 1)
		require.Equal(t, dbSpec.PipelineSpecID, pipelineSpecs[0].ID)
		pipelineSpecID := pipelineSpecs[0].ID

		// Create the run
		runID, _, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), pipelineSpecs[0], pipeline.NewVarsFrom(nil), lggr, true)
		require.NoError(t, err)

		// Check the DB for the pipeline.Run
		var pipelineRuns []pipeline.Run
		sql = `SELECT * FROM pipeline_runs WHERE id = $1;`
		err = db.Select(&pipelineRuns, sql, runID)
		require.NoError(t, err)
		require.Len(t, pipelineRuns, 1)
		require.Equal(t, pipelineSpecID, pipelineRuns[0].PipelineSpecID)
		require.Equal(t, runID, pipelineRuns[0].ID)

		// Check the DB for the pipeline.TaskRuns
		var taskRuns []pipeline.TaskRun
		sql = `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1;`
		err = db.Select(&taskRuns, sql, runID)
		require.NoError(t, err)
		require.Len(t, taskRuns, len(expectedTasks))

		for _, taskRun := range taskRuns {
			assert.Equal(t, runID, taskRun.PipelineRunID)
			assert.False(t, taskRun.Output.Valid)
			assert.False(t, taskRun.Error.IsZero())
		}
	})
}
