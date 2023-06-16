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
	err = orm.DeleteTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-2")})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
	})

	// Test deleting non-existent requests does not error.
	err = orm.DeleteTransmitRequest(&pb.TransmitRequest{Payload: []byte("does-not-exist")})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
	})

	// More inserts.
	err = orm.InsertTransmitRequest(&pb.TransmitRequest{Payload: []byte("report-4")}, reportContext)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests()
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
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
		{Req: &pb.TransmitRequest{Payload: []byte("report-1")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-3")}, ReportCtx: reportContext},
		{Req: &pb.TransmitRequest{Payload: []byte("report-4")}, ReportCtx: reportContext},
	})
}
