package mercury

import (
	"context"
	"testing"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func bootstrapPersistenceManager(t *testing.T) (*PersistenceManager, *observer.ObservedLogs) {
	t.Helper()
	db := pgtest.NewSqlxDB(t)
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	orm := NewORM(db, lggr, pgtest.NewQConfig(true))
	return NewPersistenceManager(lggr, orm, 0, 2, 5*time.Millisecond, 5*time.Millisecond), observedLogs
}

func TestPersistenceManager(t *testing.T) {
	ctx := context.Background()
	pm, _ := bootstrapPersistenceManager(t)

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
}

func TestPersistenceManagerAsyncDelete(t *testing.T) {
	ctx := context.Background()
	pm, observedLogs := bootstrapPersistenceManager(t)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{})
	require.NoError(t, err)

	err = pm.Start(ctx)
	require.NoError(t, err)

	pm.AsyncDelete(&pb.TransmitRequest{Payload: reports[0]})

	// Wait for next poll.
	observedLogs.TakeAll()
	testutils.WaitForLogMessage(t, observedLogs, "Deleted queued transmit requests")

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
	ctx := context.Background()
	pm, observedLogs := bootstrapPersistenceManager(t)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 1}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 2}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[2]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 3}})
	require.NoError(t, err)

	err = pm.Start(ctx)
	require.NoError(t, err)

	// Wait for next poll.
	observedLogs.TakeAll()
	testutils.WaitForLogMessage(t, observedLogs, "Pruned transmit requests table")

	transmissions, err := pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 3}}},
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 2}}},
	}, transmissions)

	// Test pruning stops after Close.
	err = pm.Close()
	require.NoError(t, err)

	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[3]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 4}})
	require.NoError(t, err)

	time.Sleep(15 * time.Millisecond)

	transmissions, err = pm.Load(ctx)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[3]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 4}}},
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 3}}},
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 2}}},
	}, transmissions)
}
