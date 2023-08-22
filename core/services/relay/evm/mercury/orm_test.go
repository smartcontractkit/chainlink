package mercury

import (
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

func TestORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := NewORM(db, lggr, pgtest.NewQConfig(true))

	reportContext := ocrtypes.ReportContext{
		ReportTimestamp: ocrtypes.ReportTimestamp{
			ConfigDigest: ocrtypes.ConfigDigest{'1'},
			Epoch:        10,
			Round:        20,
		},
		ExtraHash: [32]byte{'2'},
	}

	// Test insert and get requests.
	err := orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-1")}, reportContext)
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-2")}, reportContext)
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-3")}, reportContext)
	require.NoError(t, err)

	transmissions, err := orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-2")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
	})

	// Test requests can be deleted.
	err = orm.DeleteTransmitRequests([]*pb.TransmitRequest{{Payload: []byte("report-2")}})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
	})

	// Test deleting non-existent requests does not error.
	err = orm.DeleteTransmitRequests([]*pb.TransmitRequest{{Payload: []byte("does-not-exist")}})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
	})

	// Test deleting multiple requests.
	err = orm.DeleteTransmitRequests([]*pb.TransmitRequest{
		{Payload: []byte("report-1")},
		{Payload: []byte("report-3")},
	})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Empty(t, transmissions)

	// More inserts.
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-4")}, reportContext)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-4")}, ReportCtx: reportContext},
	})

	// Duplicate requests are ignored.
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-4")}, reportContext)
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-4")}, reportContext)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-4")}, ReportCtx: reportContext},
	})
}

func TestORM_PruneTransmitRequests(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := NewORM(db, lggr, pgtest.NewQConfig(true))

	makeReportContext := func(epoch uint32, round uint8) ocrtypes.ReportContext {
		return ocrtypes.ReportContext{
			ReportTimestamp: ocrtypes.ReportTimestamp{
				ConfigDigest: ocrtypes.ConfigDigest{'1'},
				Epoch:        epoch,
				Round:        round,
			},
			ExtraHash: [32]byte{'2'},
		}
	}

	err := orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("1")}, makeReportContext(1, 1))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("2")}, makeReportContext(1, 2))
	require.NoError(t, err)

	// Max size greater than table size, expect no-op
	err = orm.PruneTransmitRequests(5)
	require.NoError(t, err)

	transmissions, err := orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("2")}, ReportCtx: makeReportContext(1, 2)},
		{Req: &pb.TransmitRequest{Payload: []byte("1")}, ReportCtx: makeReportContext(1, 1)},
	})

	// Max size equal to table size, expect no-op
	err = orm.PruneTransmitRequests(2)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("2")}, ReportCtx: makeReportContext(1, 2)},
		{Req: &pb.TransmitRequest{Payload: []byte("1")}, ReportCtx: makeReportContext(1, 1)},
	})

	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("3")}, makeReportContext(2, 1))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("4")}, makeReportContext(2, 2))
	require.NoError(t, err)

	// Max size is table size + 1, expect the oldest row to be pruned.
	err = orm.PruneTransmitRequests(3)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("4")}, ReportCtx: makeReportContext(2, 2)},
		{Req: &pb.TransmitRequest{Payload: []byte("3")}, ReportCtx: makeReportContext(2, 1)},
		{Req: &pb.TransmitRequest{Payload: []byte("2")}, ReportCtx: makeReportContext(1, 2)},
	})
}
