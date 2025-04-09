package fluxmonitorv2_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	commontxmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestORM_MostRecentFluxMonitorRoundID(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	db := pgtest.NewSqlxDB(t)
	orm := newORM(t, db, nil)

	address := testutils.NewAddress()

	// Setup the rounds
	for round := uint32(0); round < 10; round++ {
		_, err := orm.FindOrCreateFluxMonitorRoundStats(ctx, address, round, 1)
		require.NoError(t, err)
	}

	count, err := orm.CountFluxMonitorRoundStats(ctx)
	require.NoError(t, err)
	require.Equal(t, 10, count)

	// Ensure round stats are not created again for the same address/roundID
	stats, err := orm.FindOrCreateFluxMonitorRoundStats(ctx, address, uint32(0), 1)
	require.NoError(t, err)
	require.Equal(t, uint32(0), stats.RoundID)
	require.Equal(t, address, stats.Aggregator)
	require.Equal(t, uint64(1), stats.NumNewRoundLogs)

	count, err = orm.CountFluxMonitorRoundStats(ctx)
	require.NoError(t, err)
	require.Equal(t, 10, count)

	roundID, err := orm.MostRecentFluxMonitorRoundID(ctx, testutils.NewAddress())
	require.Error(t, err)
	require.Equal(t, uint32(0), roundID)

	roundID, err = orm.MostRecentFluxMonitorRoundID(ctx, address)
	require.NoError(t, err)
	require.Equal(t, uint32(9), roundID)

	// Deleting rounds against a new address should incur no changes
	err = orm.DeleteFluxMonitorRoundsBackThrough(ctx, testutils.NewAddress(), 5)
	require.NoError(t, err)

	count, err = orm.CountFluxMonitorRoundStats(ctx)
	require.NoError(t, err)
	require.Equal(t, 10, count)

	// Deleting rounds against the address
	err = orm.DeleteFluxMonitorRoundsBackThrough(ctx, address, 5)
	require.NoError(t, err)

	count, err = orm.CountFluxMonitorRoundStats(ctx)
	require.NoError(t, err)
	require.Equal(t, 5, count)
}

func TestORM_UpdateFluxMonitorRoundStats(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	cfg := configtest.NewGeneralConfig(t, nil)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db)
	lggr := logger.TestLogger(t)

	// Instantiate a real pipeline ORM because we need to create a pipeline run
	// for the foreign key constraint of the stats record
	pipelineORM := pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns())
	bridgeORM := bridges.NewORM(db)

	// Instantiate a real job ORM because we need to create a job to satisfy
	// a check in pipeline.CreateRun
	jobORM := job.NewORM(db, pipelineORM, bridgeORM, keyStore, lggr)
	orm := newORM(t, db, nil)

	address := testutils.NewAddress()
	var roundID uint32 = 1

	jb := makeJob(t)
	require.NoError(t, jobORM.CreateJob(ctx, jb))

	for expectedCount := uint64(1); expectedCount < 4; expectedCount++ {
		f := time.Now()
		run :=
			&pipeline.Run{
				State:          pipeline.RunStatusCompleted,
				PipelineSpecID: jb.PipelineSpec.ID,
				PruningKey:     jb.ID,
				PipelineSpec:   *jb.PipelineSpec,
				CreatedAt:      time.Now(),
				FinishedAt:     null.TimeFrom(f),
				AllErrors:      pipeline.RunErrors{null.String{}},
				FatalErrors:    pipeline.RunErrors{null.String{}},
				Outputs:        jsonserializable.JSONSerializable{Val: []interface{}{10}, Valid: true},
				PipelineTaskRuns: []pipeline.TaskRun{
					{
						ID:         uuid.New(),
						Type:       pipeline.TaskTypeHTTP,
						Output:     jsonserializable.JSONSerializable{Val: 10, Valid: true},
						CreatedAt:  f,
						FinishedAt: null.TimeFrom(f),
					},
				},
			}
		err := pipelineORM.InsertFinishedRun(ctx, run, true)
		require.NoError(t, err)

		err = orm.UpdateFluxMonitorRoundStats(ctx, address, roundID, run.ID, 0)
		require.NoError(t, err)

		stats, err := orm.FindOrCreateFluxMonitorRoundStats(ctx, address, roundID, 0)
		require.NoError(t, err)
		require.Equal(t, expectedCount, stats.NumSubmissions)
		require.True(t, stats.PipelineRunID.Valid)
		require.Equal(t, run.ID, stats.PipelineRunID.Int64)
	}
}

func makeJob(t *testing.T) *job.Job {
	t.Helper()

	return &job.Job{
		ID:            1,
		Type:          "fluxmonitor",
		SchemaVersion: 1,
		ExternalJobID: uuid.New(),
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
			EVMChainID:        (*big.Big)(testutils.FixtureChainID),
		},
	}
}

func TestORM_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	strategy := commontxmmocks.NewTxStrategy(t)

	var (
		txm = txmmocks.NewMockEvmTxManager(t)
		orm = fluxmonitorv2.NewORM(db, logger.TestLogger(t), txm, strategy, txmgr.TransmitCheckerSpec{})

		_, from  = cltest.MustInsertRandomKey(t, ethKeyStore)
		to       = testutils.NewAddress()
		payload  = []byte{1, 0, 0}
		gasLimit = uint64(21000)
	)
	idempotencyKey := uuid.New().String()
	txm.On("CreateTransaction", mock.Anything, txmgr.TxRequest{
		IdempotencyKey: &idempotencyKey,
		FromAddress:    from,
		ToAddress:      to,
		EncodedPayload: payload,
		FeeLimit:       gasLimit,
		Meta:           nil,
		Strategy:       strategy,
	}).Return(txmgr.Tx{}, nil).Once()

	require.NoError(t, orm.CreateEthTransaction(testutils.Context(t), from, to, payload, gasLimit, &idempotencyKey))
}
