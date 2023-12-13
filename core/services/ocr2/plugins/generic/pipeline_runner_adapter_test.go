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

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	_ "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/generic"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const spec = `
answer [type=sum values=<[ $(val), 2 ]>]
answer;
`

func TestAdapter_Integration(t *testing.T) {
	logger := logger.TestLogger(t)
	cfg := configtest.NewTestGeneralConfig(t)
	url := cfg.Database().URL()
	db, err := pg.NewConnection(url.String(), cfg.Database().Dialect(), cfg.Database())
	require.NoError(t, err)

	keystore := keystore.NewInMemory(db, utils.FastScryptParams, logger, cfg.Database())
	pipelineORM := pipeline.NewORM(db, logger, cfg.Database(), cfg.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db, logger, cfg.Database())
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
	pra := generic.NewPipelineRunnerAdapter(logger, job.Job{}, pr)
	results, err := pra.ExecuteRun(testutils.Context(t), spec, types.Vars{Vars: map[string]interface{}{"val": 1}}, types.Options{})
	require.NoError(t, err)

	finalResult := results[0].Value.(decimal.Decimal)

	assert.True(t, decimal.NewFromInt(3).Equal(finalResult))
}

func newMockPipelineRunner() *mockPipelineRunner {
	return &mockPipelineRunner{}
}

type mockPipelineRunner struct {
	results pipeline.TaskRunResults
	err     error
	run     *pipeline.Run
	spec    pipeline.Spec
	vars    pipeline.Vars
}

func (m *mockPipelineRunner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (*pipeline.Run, pipeline.TaskRunResults, error) {
	m.spec = spec
	m.vars = vars
	return m.run, m.results, m.err
}

func TestAdapter_AddsDefaultVars(t *testing.T) {
	logger := logger.TestLogger(t)
	mpr := newMockPipelineRunner()
	jobID, externalJobID, name := int32(100), uuid.New(), null.StringFrom("job-name")
	pra := generic.NewPipelineRunnerAdapter(logger, job.Job{ID: jobID, ExternalJobID: externalJobID, Name: name}, mpr)

	_, err := pra.ExecuteRun(testutils.Context(t), spec, types.Vars{}, types.Options{})
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
	_, err := pra.ExecuteRun(testutils.Context(t), spec, types.Vars{}, types.Options{MaxTaskDuration: maxDuration})
	require.NoError(t, err)

	assert.Equal(t, jobID, mpr.spec.JobID)
	assert.Equal(t, name.ValueOrZero(), mpr.spec.JobName)
	assert.Equal(t, string(jobType), mpr.spec.JobType)
	assert.Equal(t, maxDuration, mpr.spec.MaxTaskDuration.Duration())
}
