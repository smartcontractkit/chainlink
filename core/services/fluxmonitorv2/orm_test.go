package fluxmonitorv2_test

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	bptxmmocks "github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager/mocks"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/stretchr/testify/require"
)

func TestORM_MostRecentFluxMonitorRoundID(t *testing.T) {
	t.Parallel()

	corestore, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := fluxmonitorv2.NewORM(corestore.DB, nil, nil)

	address := cltest.NewAddress()

	// Setup the rounds
	for round := uint32(0); round < 10; round++ {
		_, err := orm.FindOrCreateFluxMonitorRoundStats(address, round)
		require.NoError(t, err)
	}

	count, err := orm.CountFluxMonitorRoundStats()
	require.NoError(t, err)
	require.Equal(t, 10, count)

	// Ensure round stats are not created again for the same address/roundID
	stats, err := orm.FindOrCreateFluxMonitorRoundStats(address, uint32(0))
	require.NoError(t, err)
	require.Equal(t, uint32(0), stats.RoundID)
	require.Equal(t, address, stats.Aggregator)

	count, err = orm.CountFluxMonitorRoundStats()
	require.NoError(t, err)
	require.Equal(t, 10, count)

	roundID, err := orm.MostRecentFluxMonitorRoundID(cltest.NewAddress())
	require.Error(t, err)
	require.Equal(t, uint32(0), roundID)

	roundID, err = orm.MostRecentFluxMonitorRoundID(address)
	require.NoError(t, err)
	require.Equal(t, uint32(9), roundID)

	// Deleting rounds against a new address should incur no changes
	err = orm.DeleteFluxMonitorRoundsBackThrough(cltest.NewAddress(), 5)
	require.NoError(t, err)

	count, err = orm.CountFluxMonitorRoundStats()
	require.NoError(t, err)
	require.Equal(t, 10, count)

	// Deleting rounds against the address
	err = orm.DeleteFluxMonitorRoundsBackThrough(address, 5)
	require.NoError(t, err)

	count, err = orm.CountFluxMonitorRoundStats()
	require.NoError(t, err)
	require.Equal(t, 5, count)
}

func TestORM_UpdateFluxMonitorRoundStats(t *testing.T) {
	t.Parallel()

	corestore, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	// Instantiate a real pipeline ORM because we need to create a pipeline run
	// for the foreign key constraint of the stats record
	eventBroadcaster := postgres.NewEventBroadcaster(
		corestore.Config.DatabaseURL(),
		corestore.Config.DatabaseListenerMinReconnectInterval(),
		corestore.Config.DatabaseListenerMaxReconnectDuration(),
	)
	pipelineORM := pipeline.NewORM(corestore.DB)
	// Instantiate a real job ORM because we need to create a job to satisfy
	// a check in pipeline.CreateRun
	jobORM := job.NewORM(corestore.ORM.DB, corestore.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	orm := fluxmonitorv2.NewORM(corestore.DB, nil, nil)

	address := cltest.NewAddress()
	var roundID uint32 = 1

	j := makeJob(t)
	jb, err := jobORM.CreateJob(context.Background(), j, pipeline.Pipeline{})
	require.NoError(t, err)

	for expectedCount := uint64(1); expectedCount < 4; expectedCount++ {
		f := time.Now()
		runID, err := pipelineORM.InsertFinishedRun(
			corestore.DB,
			pipeline.Run{
				State:          pipeline.RunStatusCompleted,
				PipelineSpecID: jb.PipelineSpec.ID,
				PipelineSpec:   *jb.PipelineSpec,
				CreatedAt:      time.Now(),
				FinishedAt:     null.TimeFrom(f),
				Errors:         pipeline.RunErrors{null.String{}},
				Outputs:        pipeline.JSONSerializable{Val: []interface{}{10}},
			}, pipeline.TaskRunResults{
				{
					ID:         uuid.NewV4(),
					Task:       &pipeline.HTTPTask{},
					Result:     pipeline.Result{Value: 10},
					CreatedAt:  f,
					FinishedAt: null.TimeFrom(f),
				},
			}, true)
		require.NoError(t, err)

		err = orm.UpdateFluxMonitorRoundStats(corestore.DB, address, roundID, runID)
		require.NoError(t, err)

		stats, err := orm.FindOrCreateFluxMonitorRoundStats(address, roundID)
		require.NoError(t, err)
		require.Equal(t, expectedCount, stats.NumSubmissions)
		require.True(t, stats.PipelineRunID.Valid)
		require.Equal(t, runID, stats.PipelineRunID.Int64)
	}
}

func makeJob(t *testing.T) *job.Job {
	t.Helper()

	return &job.Job{
		ID:            1,
		Type:          "fluxmonitor",
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
		FluxMonitorSpec: &job.FluxMonitorSpec{
			ID:                2,
			ContractAddress:   cltest.NewEIP55Address(),
			Threshold:         0.5,
			PollTimerPeriod:   1 * time.Second,
			PollTimerDisabled: false,
			IdleTimerPeriod:   1 * time.Minute,
			IdleTimerDisabled: false,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
}

func TestORM_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	corestore, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	strategy := new(bptxmmocks.TxStrategy)

	var (
		txm = new(bptxmmocks.TxManager)
		orm = fluxmonitorv2.NewORM(corestore.DB, txm, strategy)

		key      = cltest.MustInsertRandomKey(t, corestore.DB, 0)
		from     = key.Address.Address()
		to       = cltest.NewAddress()
		payload  = []byte{1, 0, 0}
		gasLimit = uint64(21000)
	)

	txm.On("CreateEthTransaction", corestore.DB, from, to, payload, gasLimit, nil, strategy).Return(bulletprooftxmanager.EthTx{}, nil).Once()

	orm.CreateEthTransaction(corestore.DB, from, to, payload, gasLimit)

	txm.AssertExpectations(t)
}
