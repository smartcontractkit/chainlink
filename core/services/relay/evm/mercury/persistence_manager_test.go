package mercury

import (
	"context"
	"testing"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func bootstrapPersistenceManager(t *testing.T) *PersistenceManager {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := NewORM(db, lggr, pgtest.NewQConfig(true))
	return NewPersistenceManager(lggr, orm)
}

func TestPersistenceManager(t *testing.T) {
	ctx := context.Background()
	pm := bootstrapPersistenceManager(t)

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
	pm := bootstrapPersistenceManager(t)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{})
	require.NoError(t, err)

	flushDeletesFrequency = 10 * time.Millisecond
	err = pm.Start(ctx)
	require.NoError(t, err)

	pm.AsyncDelete(&pb.TransmitRequest{Payload: reports[0]})

	time.Sleep(15 * time.Millisecond)

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
	pm := bootstrapPersistenceManager(t)

	reports := sampleReports

	err := pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[0]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 1}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[1]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 2}})
	require.NoError(t, err)
	err = pm.Insert(ctx, &pb.TransmitRequest{Payload: reports[2]}, ocrtypes.ReportContext{ReportTimestamp: ocrtypes.ReportTimestamp{Epoch: 3}})
	require.NoError(t, err)

	maxTransmitQueueSize = 2
	pruneFrequency = 10 * time.Millisecond
	err = pm.Start(ctx)
	require.NoError(t, err)

	time.Sleep(15 * time.Millisecond)

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
