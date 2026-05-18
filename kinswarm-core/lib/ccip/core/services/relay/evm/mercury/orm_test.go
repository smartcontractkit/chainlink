package mercury

import (
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var (
	sURL  = "wss://example.com/mercury"
	sURL2 = "wss://mercuryserver.test"
	sURL3 = "wss://mercuryserver.example/foo"
)

func TestORM(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)

	jobID := rand.Int32() // foreign key constraints disabled so value doesn't matter
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	orm := NewORM(db)
	feedID := sampleFeedID

	reports := sampleReports
	reportContexts := make([]ocrtypes.ReportContext, 4)
	for i := range reportContexts {
		reportContexts[i] = ocrtypes.ReportContext{
			ReportTimestamp: ocrtypes.ReportTimestamp{
				ConfigDigest: ocrtypes.ConfigDigest{'1'},
				Epoch:        10,
				Round:        uint8(i),
			},
			ExtraHash: [32]byte{'2'},
		}
	}

	l, err := orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Nil(t, l)

	// Test insert and get requests.
	// s1
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[0]}, jobID, reportContexts[0])
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[1]}, jobID, reportContexts[1])
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[2]}, jobID, reportContexts[2])
	require.NoError(t, err)

	// s2
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[3]}, jobID, reportContexts[0])
	require.NoError(t, err)

	transmissions, err := orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: reportContexts[2]},
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: reportContexts[1]},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]},
	})
	transmissions, err = orm.GetTransmitRequests(ctx, sURL2, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[3]}, ReportCtx: reportContexts[0]},
	})

	l, err = orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.NotEqual(t, reports[0], l)
	assert.Equal(t, reports[2], l)

	// Test requests can be deleted.
	err = orm.DeleteTransmitRequests(ctx, sURL, []*pb.TransmitRequest{{Payload: reports[1]}})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: reportContexts[2]},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]},
	})

	l, err = orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Equal(t, reports[2], l)

	// Test deleting non-existent requests does not error.
	err = orm.DeleteTransmitRequests(ctx, sURL, []*pb.TransmitRequest{{Payload: []byte("does-not-exist")}})
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: reportContexts[2]},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]},
	})

	// Test deleting multiple requests.
	err = orm.DeleteTransmitRequests(ctx, sURL, []*pb.TransmitRequest{
		{Payload: reports[0]},
		{Payload: reports[2]},
	})
	require.NoError(t, err)

	l, err = orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Equal(t, reports[2], l)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Empty(t, transmissions)

	// More inserts.
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[3]}, jobID, reportContexts[3])
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[3]}, ReportCtx: reportContexts[3]},
	})

	// Duplicate requests are ignored.
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[3]}, jobID, reportContexts[3])
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[3]}, jobID, reportContexts[3])
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[3]}, ReportCtx: reportContexts[3]},
	})

	l, err = orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Equal(t, reports[3], l)

	// s2 not affected by deletion
	transmissions, err = orm.GetTransmitRequests(ctx, sURL2, jobID)
	require.NoError(t, err)
	require.Len(t, transmissions, 1)
}

func TestORM_InsertTransmitRequest_MultipleServerURLs(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)

	jobID := rand.Int32() // foreign key constraints disabled so value doesn't matter
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	orm := NewORM(db)
	feedID := sampleFeedID

	reports := sampleReports
	reportContexts := make([]ocrtypes.ReportContext, 4)
	for i := range reportContexts {
		reportContexts[i] = ocrtypes.ReportContext{
			ReportTimestamp: ocrtypes.ReportTimestamp{
				ConfigDigest: ocrtypes.ConfigDigest{'1'},
				Epoch:        10,
				Round:        uint8(i),
			},
			ExtraHash: [32]byte{'2'},
		}
	}
	err := orm.InsertTransmitRequest(ctx, []string{sURL, sURL2, sURL3}, &pb.TransmitRequest{Payload: reports[0]}, jobID, reportContexts[0])
	require.NoError(t, err)

	transmissions, err := orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Len(t, transmissions, 1)
	assert.Equal(t, transmissions[0], &Transmission{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]})

	transmissions, err = orm.GetTransmitRequests(ctx, sURL2, jobID)
	require.NoError(t, err)
	require.Len(t, transmissions, 1)
	assert.Equal(t, transmissions[0], &Transmission{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]})

	transmissions, err = orm.GetTransmitRequests(ctx, sURL3, jobID)
	require.NoError(t, err)
	require.Len(t, transmissions, 1)
	assert.Equal(t, transmissions[0], &Transmission{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: reportContexts[0]})

	l, err := orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Equal(t, reports[0], l)
}

func TestORM_PruneTransmitRequests(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	jobID := rand.Int32() // foreign key constraints disabled so value doesn't matter
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)

	orm := NewORM(db)

	reports := sampleReports

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

	// s1
	err := orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[0]}, jobID, makeReportContext(1, 1))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[1]}, jobID, makeReportContext(1, 2))
	require.NoError(t, err)
	// s2 - should not be touched
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[0]}, jobID, makeReportContext(1, 0))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[0]}, jobID, makeReportContext(1, 1))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[1]}, jobID, makeReportContext(1, 2))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[2]}, jobID, makeReportContext(1, 3))
	require.NoError(t, err)

	// Max size greater than number of records, expect no-op
	err = orm.PruneTransmitRequests(ctx, sURL, jobID, 5)
	require.NoError(t, err)

	transmissions, err := orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: makeReportContext(1, 2)},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: makeReportContext(1, 1)},
	})

	// Max size equal to number of records, expect no-op
	err = orm.PruneTransmitRequests(ctx, sURL, jobID, 2)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, transmissions, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: makeReportContext(1, 2)},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: makeReportContext(1, 1)},
	})

	// Max size is number of records + 1, but jobID differs, expect no-op
	err = orm.PruneTransmitRequests(ctx, sURL, -1, 2)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: makeReportContext(1, 2)},
		{Req: &pb.TransmitRequest{Payload: reports[0]}, ReportCtx: makeReportContext(1, 1)},
	}, transmissions)

	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[2]}, jobID, makeReportContext(2, 1))
	require.NoError(t, err)
	err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[3]}, jobID, makeReportContext(2, 2))
	require.NoError(t, err)

	// Max size is table size - 1, expect the oldest row to be pruned.
	err = orm.PruneTransmitRequests(ctx, sURL, jobID, 3)
	require.NoError(t, err)

	transmissions, err = orm.GetTransmitRequests(ctx, sURL, jobID)
	require.NoError(t, err)
	require.Equal(t, []*Transmission{
		{Req: &pb.TransmitRequest{Payload: reports[3]}, ReportCtx: makeReportContext(2, 2)},
		{Req: &pb.TransmitRequest{Payload: reports[2]}, ReportCtx: makeReportContext(2, 1)},
		{Req: &pb.TransmitRequest{Payload: reports[1]}, ReportCtx: makeReportContext(1, 2)},
	}, transmissions)

	// s2 not touched
	transmissions, err = orm.GetTransmitRequests(ctx, sURL2, jobID)
	require.NoError(t, err)
	assert.Len(t, transmissions, 3)
}

func TestORM_InsertTransmitRequest_LatestReport(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	jobID := rand.Int32() // foreign key constraints disabled so value doesn't matter
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)

	orm := NewORM(db)
	feedID := sampleFeedID

	reports := sampleReports

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

	err := orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[0]}, jobID, makeReportContext(
		0, 0,
	))
	require.NoError(t, err)

	// this should be ignored, because report context is the same
	err = orm.InsertTransmitRequest(ctx, []string{sURL2}, &pb.TransmitRequest{Payload: reports[1]}, jobID, makeReportContext(
		0, 0,
	))
	require.NoError(t, err)

	l, err := orm.LatestReport(testutils.Context(t), feedID)
	require.NoError(t, err)
	assert.Equal(t, reports[0], l)

	t.Run("replaces if epoch and round are larger", func(t *testing.T) {
		err = orm.InsertTransmitRequest(ctx, []string{"foo"}, &pb.TransmitRequest{Payload: reports[1]}, jobID, makeReportContext(1, 1))
		require.NoError(t, err)

		l, err = orm.LatestReport(testutils.Context(t), feedID)
		require.NoError(t, err)
		assert.Equal(t, reports[1], l)
	})
	t.Run("replaces if epoch is the same but round is greater", func(t *testing.T) {
		err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[2]}, jobID, makeReportContext(1, 2))
		require.NoError(t, err)

		l, err = orm.LatestReport(testutils.Context(t), feedID)
		require.NoError(t, err)
		assert.Equal(t, reports[2], l)
	})
	t.Run("replaces if epoch is larger but round is smaller", func(t *testing.T) {
		err = orm.InsertTransmitRequest(ctx, []string{"bar"}, &pb.TransmitRequest{Payload: reports[3]}, jobID, makeReportContext(2, 1))
		require.NoError(t, err)

		l, err = orm.LatestReport(testutils.Context(t), feedID)
		require.NoError(t, err)
		assert.Equal(t, reports[3], l)
	})
	t.Run("does not overwrite if epoch/round is the same", func(t *testing.T) {
		err = orm.InsertTransmitRequest(ctx, []string{sURL}, &pb.TransmitRequest{Payload: reports[0]}, jobID, makeReportContext(2, 1))
		require.NoError(t, err)

		l, err = orm.LatestReport(testutils.Context(t), feedID)
		require.NoError(t, err)
		assert.Equal(t, reports[3], l)
	})
}

func Test_ReportCodec_FeedIDFromReport(t *testing.T) {
	t.Run("FeedIDFromReport extracts the current block number from a valid report", func(t *testing.T) {
		report := buildSampleV1Report(42)

		f, err := FeedIDFromReport(report)
		require.NoError(t, err)

		assert.Equal(t, sampleFeedID[:], f[:])
	})
	t.Run("FeedIDFromReport returns error if report is invalid", func(t *testing.T) {
		report := []byte{1}

		_, err := FeedIDFromReport(report)
		assert.EqualError(t, err, "invalid length for report: 1")
	})
}
