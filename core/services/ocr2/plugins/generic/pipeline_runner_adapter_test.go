package generic_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	_ "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const spec = `
answer [type=sum values=<[ $(val), 2 ]>]
answer;
`

func TestAdapter_Integration(t *testing.T) {
	testutils.SkipShortDB(t)
	ctx := testutils.Context(t)
	logger := logger.TestLogger(t)
	cfg := configtest.NewTestGeneralConfig(t)
	url := cfg.Database().URL()
	db, err := pg.NewConnection(ctx, url.String(), cfg.Database().DriverName(), cfg.Database())
	require.NoError(t, err)

	keystore := keystore.NewInMemory(db, utils.FastScryptParams, logger)
	pipelineORM := pipeline.NewORM(db, logger, cfg.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)
	jobORM := job.NewORM(db, pipelineORM, bridgesORM, keystore, logger)
	pr := pipeline.NewRunner(
		pipelineORM,
		bridgesORM,
		cfg.JobPipeline(),
		cfg.WebServer(),
		nil,
		keystore.Eth(),
		keystore.VRF(),
		logger,
		http.DefaultClient,
		http.DefaultClient,
	)
	err = keystore.Unlock(ctx, cfg.Password().Keystore())
	require.NoError(t, err)
	jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), cfg.OCR2(), cfg.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)

	const juelsPerFeeCoinSource = `
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;`

	_, address := cltest.MustInsertRandomKey(t, keystore.Eth())
	jb.Name = null.StringFrom("Job 1")
	jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource
	err = jobORM.CreateJob(ctx, &jb)
	require.NoError(t, err)
	pra := generic.NewPipelineRunnerAdapter(logger, jb, pr)
	results, err := pra.ExecuteRun(testutils.Context(t), spec, core.Vars{Vars: map[string]interface{}{"val": 1}}, core.Options{})
	require.NoError(t, err)

	finalResult := results[0].Value.Val.(decimal.Decimal)

	assert.True(t, decimal.NewFromInt(3).Equal(finalResult))
}

func newMockPipelineRunner() *mockPipelineRunner {
	return &mockPipelineRunner{}
}

type mockPipelineRunner struct {
	results pipeline.TaskRunResults
	err     error
	spec    pipeline.Spec
	vars    pipeline.Vars
}

func (m *mockPipelineRunner) ExecuteAndInsertFinishedRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, saveSuccessfulTaskRuns bool) (runID int64, results pipeline.TaskRunResults, err error) {
	m.spec = spec
	m.vars = vars
	// We never attach a run to the mock, so we can't return a runID
	return 0, m.results, m.err
}

func TestAdapter_AddsDefaultVars(t *testing.T) {
	logger := logger.TestLogger(t)
	mpr := newMockPipelineRunner()
	jobID, externalJobID, name := int32(100), uuid.New(), null.StringFrom("job-name")
	pra := generic.NewPipelineRunnerAdapter(logger, job.Job{ID: jobID, ExternalJobID: externalJobID, Name: name}, mpr)

	_, err := pra.ExecuteRun(testutils.Context(t), spec, core.Vars{}, core.Options{})
	require.NoError(t, err)

	gotName, err := mpr.vars.Get("jb.name")
	require.NoError(t, err)
	assert.Equal(t, name.String, gotName)

	gotID, err := mpr.vars.Get("jb.databaseID")
	require.NoError(t, err)
	assert.Equal(t, jobID, gotID)

	gotExternalID, err := mpr.vars.Get("jb.externalJobID")
	require.NoError(t, err)
	assert.Equal(t, externalJobID, gotExternalID)
}

func TestPipelineRunnerAdapter_SetsVarsOnSpec(t *testing.T) {
	logger := logger.TestLogger(t)
	mpr := newMockPipelineRunner()
	jobID, externalJobID, name, jobType := int32(100), uuid.New(), null.StringFrom("job-name"), job.Type("generic")
	pra := generic.NewPipelineRunnerAdapter(logger, job.Job{ID: jobID, ExternalJobID: externalJobID, Name: name, Type: jobType}, mpr)

	maxDuration := 100 * time.Second
	_, err := pra.ExecuteRun(testutils.Context(t), spec, core.Vars{}, core.Options{MaxTaskDuration: maxDuration})
	require.NoError(t, err)

	assert.Equal(t, jobID, mpr.spec.JobID)
	assert.Equal(t, name.ValueOrZero(), mpr.spec.JobName)
	assert.Equal(t, string(jobType), mpr.spec.JobType)
	assert.Equal(t, maxDuration, mpr.spec.MaxTaskDuration.Duration())
}
