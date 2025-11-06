package mercury

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func bootstrapPersistenceManager(t *testing.T, jobID int32, db *sqlx.DB) *PersistenceManager {
	t.Helper()
	orm := NewORM(db)
	return NewPersistenceManager(logger.Test(t), "mercuryserver.example", orm, jobID, 2, 5*time.Millisecond, 5*time.Millisecond)
}

func TestPersistenceManager(t *testing.T) {
	jobID1 := rand.Int32()
	jobID2 := jobID1 + 1

	ctx := testutils.Context(t)
	db := testutils.NewSqlxDB(t)
	testutils.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	testutils.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	pm := bootstrapPersistenceManager(t, jobID1, db)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{})
	require.NoError(t, err)

	transmissions, err := pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[0]}},
		{Req: &pb.TransmitRequest{Payload: reports[1]}},
	}, transmissions)

	err = pm.Delete(ctx, &pb.TransmitRequest{Payload: reports[0]})
	require.NoError(t, err)

	transmissions, err = pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}},
	}, transmissions)

	t.Run("scopes load to only transmissions with matching job ID", func(t *testing.T) {
		pm2 := bootstrapPersistenceManager(t, jobID2, db)
		transmissions, err = pm2.Load(ctx)
		require.NoError(t, err)

		assert.Len(t, transmissions, 0)
	})
}

func TestPersistenceManagerAsyncDelete(t *testing.T) {
	ctx := testutils.Context(t)
	jobID := rand.Int32()
	db := testutils.NewSqlxDB(t)
	testutils.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	testutils.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	pm := bootstrapPersistenceManager(t, jobID, db)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{})
	require.NoError(t, err)

	err = pm.Start(ctx)
	require.NoError(t, err)

	pm.AsyncDelete(&pb.TransmitRequest{Payload: reports[0]})

	// Wait for next poll.
	testutils.RequireEventually(t, func() bool {
		pm.deleteMu.Lock()
		defer pm.deleteMu.Unlock()
		return len(pm.deleteQueue) == 0
	})

	transmissions, err := pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}},
	}, transmissions)

	// Test AsyncDelete is a no-op after Close.
	err = pm.Close()
	require.NoError(t, err)

	pm.AsyncDelete(&pb.TransmitRequest{Payload: reports[1]})

	time.Sleep(15 * time.Millisecond)

	transmissions, err = pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}},
	}, transmissions)
}

func TestPersistenceManagerPrune(t *testing.T) {
	jobID1 := rand.Int32()
	jobID2 := jobID1 + 1
	db := testutils.NewSqlxDB(t)
	testutils.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	testutils.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)

	ctx := testutils.Context(t)

	reports := make([][]byte, 25)
	for i := range 25 {
		reports[i] = buildSampleV2Report(int64(i))
	}

	pm2 := bootstrapPersistenceManager(t, jobID2, db)
	for i := range 20 {
		err := pm2.Insert(ctx, &pb.TransmitRequest{Payload: reports[i]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: uint32(i)}}) //nolint:gosec // G115
		require.NoError(t, err)
	}

	pm := bootstrapPersistenceManager(t, jobID1, db)

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[21]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 21}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[22]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 22}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[23]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 23}})
	require.NoError(t, err)

	err = pm.Start(ctx)
	require.NoError(t, err)

	testutils.RequireEventually(t, func() bool {
		requests, err2 := pm.Load(testutils.Context(t))
		require.NoError(t, err2)
		return len(requests) == pm.maxTransmitQueueSize
	})

	transmissions, err := pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[23]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 23}}},
		{Req: &pb.TransmitRequest{Payload: reports[22]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 22}}},
	}, transmissions)

	// Test pruning stops after Close.
	err = pm.Close()
	require.NoError(t, err)

	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[24]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 24}})
	require.NoError(t, err)

	transmissions, err = pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[24]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 24}}},
		{Req: &pb.TransmitRequest{Payload: reports[23]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 23}}},
		{Req: &pb.TransmitRequest{Payload: reports[22]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 22}}},
	}, transmissions)

	t.Run("prune was scoped to job ID", func(t *testing.T) {
		transmissions, err = pm2.Load(ctx)
		require.NoError(t, err)
		assert.Len(t, transmissions, 20)
	})
}
