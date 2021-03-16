package fluxmonitorv2_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func TestORM_MostRecentFluxMonitorRoundID(t *testing.T) {
	t.Parallel()

	corestore, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	orm := fluxmonitorv2.NewORM(corestore.DB)

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
	pipelineORM := pipeline.NewORM(corestore.ORM.DB, corestore.Config, eventBroadcaster)
	// Instantiate a real job ORM because we need to create a job to satisfy
	// a check in pipeline.CreateRun
	jobORM := job.NewORM(corestore.ORM.DB, corestore.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	orm := fluxmonitorv2.NewORM(corestore.DB)

	address := cltest.NewAddress()
	var roundID uint32 = 1

	j := makeJob(t)
	err := jobORM.CreateJob(context.Background(), j, *pipeline.NewTaskDAG())
	require.NoError(t, err)

	for expectedCount := uint64(1); expectedCount < 4; expectedCount++ {
		runID, err := pipelineORM.CreateRun(context.Background(), j.ID, map[string]interface{}{})
		require.NoError(t, err)

		err = orm.UpdateFluxMonitorRoundStats(address, roundID, runID)
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
		IDEmbed:       job.IDEmbed{ID: 1},
		Type:          "fluxmonitor",
		SchemaVersion: 1,
		Pipeline:      *pipeline.NewTaskDAG(),
		FluxMonitorSpec: &job.FluxMonitorSpec{
			IDEmbed:           job.IDEmbed{ID: 2},
			ContractAddress:   cltest.NewEIP55Address(),
			Precision:         2,
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

	var (
		orm = fluxmonitorv2.NewORM(corestore.DB)

		key      = cltest.MustInsertRandomKey(t, corestore.DB, 0)
		from     = key.Address.Address()
		to       = cltest.NewAddress()
		payload  = []byte{1, 0, 0}
		gasLimit = uint64(21000)
	)

	orm.CreateEthTransaction(from, to, payload, gasLimit, 0)

	etx := models.EthTx{}
	require.NoError(t, corestore.ORM.DB.First(&etx).Error)

	require.Equal(t, gasLimit, etx.GasLimit)
	require.Equal(t, from, etx.FromAddress)
	require.Equal(t, to, etx.ToAddress)
	require.Equal(t, payload, etx.EncodedPayload)
	require.Equal(t, assets.NewEthValue(0), etx.Value)
}

func TestORM_CreateEthTransaction_OutOfEth(t *testing.T) {
	t.Parallel()

	corestore, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)

	var (
		orm = fluxmonitorv2.NewORM(corestore.DB)

		key      = cltest.MustInsertRandomKey(t, corestore.DB, 1)
		otherKey = cltest.MustInsertRandomKey(t, corestore.DB, 1)
		from     = key.Address.Address()
		to       = cltest.NewAddress()
		payload  = []byte{1, 0, 0}
		gasLimit = uint64(21000)
	)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, corestore, 0, otherKey.Address.Address())

		err := orm.CreateEthTransaction(from, to, payload, gasLimit, 0)
		require.NoError(t, err)

		etx := models.EthTx{}
		require.NoError(t, corestore.ORM.DB.First(&etx, "nonce IS NULL AND from_address = ?", from).Error)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, corestore.DB.Exec(`DELETE FROM eth_txes WHERE from_address = ?`, from).Error)

	t.Run("if this key has any transactions with insufficient eth errors, skips transmission entirely", func(t *testing.T) {
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, corestore, 0, from)

		err := orm.CreateEthTransaction(from, to, payload, gasLimit, 0)
		require.EqualError(t, err, fmt.Sprintf("Skipped Flux Monitor submission because wallet is out of eth: %s", from))
	})

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		require.NoError(t, corestore.DB.Exec(`UPDATE eth_tx_attempts SET state = 'broadcast'`).Error)
		require.NoError(t, corestore.DB.Exec(`UPDATE eth_txes SET nonce = 0, state = 'confirmed', broadcast_at = NOW()`).Error)

		err := orm.CreateEthTransaction(from, to, payload, gasLimit, 0)
		require.NoError(t, err)

		etx := models.EthTx{}
		require.NoError(t, corestore.ORM.DB.First(&etx, "nonce IS NULL AND from_address = ?", from).Error)
		require.Equal(t, payload, etx.EncodedPayload)
	})
}
