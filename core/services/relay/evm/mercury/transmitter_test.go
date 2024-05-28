package mercury

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mercurytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

type mockCfg struct{}

func (m mockCfg) TransmitQueueMaxSize() uint32 {
	return 10_000
}

func (m mockCfg) TransmitTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(1 * time.Hour)
}

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	var jobID int32
	pgtest.MustExec(t, db, `SET CONSTRAINTS mercury_transmit_requests_job_id_fkey DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS feed_latest_reports_job_id_fkey DEFERRED`)
	codec := new(mockCodec)
	orm := NewORM(db)
	clients := map[string]wsrpc.Client{}

	t.Run("with one mercury server", func(t *testing.T) {
		t.Run("v1 report transmission successfully enqueued", func(t *testing.T) {
			report := sampleV1Report
			c := &mocks.MockWSRPCClient{}
			clients[sURL] = c
			mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
			// init the queue since we skipped starting transmitter
			mt.servers[sURL].q.Init([]*Transmission{})
			err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)
			require.NoError(t, err)

			// ensure it was added to the queue
			require.Equal(t, mt.servers[sURL].q.(*transmitQueue).pq.Len(), 1)
			assert.Subset(t, mt.servers[sURL].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
		})
		t.Run("v2 report transmission successfully enqueued", func(t *testing.T) {
			report := sampleV2Report
			c := &mocks.MockWSRPCClient{}
			clients[sURL] = c
			mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
			// init the queue since we skipped starting transmitter
			mt.servers[sURL].q.Init([]*Transmission{})
			err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)
			require.NoError(t, err)

			// ensure it was added to the queue
			require.Equal(t, mt.servers[sURL].q.(*transmitQueue).pq.Len(), 1)
			assert.Subset(t, mt.servers[sURL].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
		})
		t.Run("v3 report transmission successfully enqueued", func(t *testing.T) {
			report := sampleV3Report
			c := &mocks.MockWSRPCClient{}
			clients[sURL] = c
			mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
			// init the queue since we skipped starting transmitter
			mt.servers[sURL].q.Init([]*Transmission{})
			err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)
			require.NoError(t, err)

			// ensure it was added to the queue
			require.Equal(t, mt.servers[sURL].q.(*transmitQueue).pq.Len(), 1)
			assert.Subset(t, mt.servers[sURL].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
		})
		t.Run("v3 report transmission sent only to trigger service", func(t *testing.T) {
			report := sampleV3Report
			c := &mocks.MockWSRPCClient{}
			clients[sURL] = c
			triggerService := triggers.NewMercuryTriggerService(0, lggr)
			mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, triggerService)
			// init the queue since we skipped starting transmitter
			mt.servers[sURL].q.Init([]*Transmission{})
			err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)
			require.NoError(t, err)
			// queue is empty
			require.Equal(t, mt.servers[sURL].q.(*transmitQueue).pq.Len(), 0)
		})
	})

	t.Run("with multiple mercury servers", func(t *testing.T) {
		report := sampleV3Report
		c := &mocks.MockWSRPCClient{}
		clients[sURL] = c
		clients[sURL2] = c
		clients[sURL3] = c

		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		// init the queue since we skipped starting transmitter
		mt.servers[sURL].q.Init([]*Transmission{})
		mt.servers[sURL2].q.Init([]*Transmission{})
		mt.servers[sURL3].q.Init([]*Transmission{})

		err := mt.Transmit(testutils.Context(t), sampleReportContext, report, sampleSigs)
		require.NoError(t, err)

		// ensure it was added to the queue
		require.Equal(t, mt.servers[sURL].q.(*transmitQueue).pq.Len(), 1)
		assert.Subset(t, mt.servers[sURL].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
		require.Equal(t, mt.servers[sURL2].q.(*transmitQueue).pq.Len(), 1)
		assert.Subset(t, mt.servers[sURL2].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
		require.Equal(t, mt.servers[sURL3].q.(*transmitQueue).pq.Len(), 1)
		assert.Subset(t, mt.servers[sURL3].q.(*transmitQueue).pq.Pop().(*Transmission).Req.Payload, report)
	})
}

func Test_MercuryTransmitter_LatestTimestamp(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	var jobID int32
	codec := new(mockCodec)

	orm := NewORM(db)
	clients := map[string]wsrpc.Client{}

	t.Run("successful query", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.ObservationsTimestamp = 42
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, int64(42), ts)
	})

	t.Run("successful query returning nil report (new feed) gives latest timestamp = -1", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, int64(-1), ts)
	})

	t.Run("failing query", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		_, err := mt.LatestTimestamp(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})

	t.Run("with multiple servers, uses latest", func(t *testing.T) {
		clients[sURL] = &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		clients[sURL2] = &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.ObservationsTimestamp = 42
				return out, nil
			},
		}
		clients[sURL3] = &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.ObservationsTimestamp = 41
				return out, nil
			},
		}
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		ts, err := mt.LatestTimestamp(testutils.Context(t))
		require.NoError(t, err)

		assert.Equal(t, int64(42), ts)
	})
}

type mockCodec struct {
	val *big.Int
	err error
}

var _ mercurytypes.ReportCodec = &mockCodec{}

func (m *mockCodec) BenchmarkPriceFromReport(_ ocrtypes.Report) (*big.Int, error) {
	return m.val, m.err
}

func (m *mockCodec) ObservationTimestampFromReport(report ocrtypes.Report) (uint32, error) {
	return 0, nil
}

func Test_MercuryTransmitter_LatestPrice(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	var jobID int32

	codec := new(mockCodec)
	orm := NewORM(db)
	clients := map[string]wsrpc.Client{}

	t.Run("successful query", func(t *testing.T) {
		originalPrice := big.NewInt(123456789)
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.Payload = buildSamplePayload([]byte("doesn't matter"))
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)

		t.Run("BenchmarkPriceFromReport succeeds", func(t *testing.T) {
			codec.val = originalPrice
			codec.err = nil

			price, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
			require.NoError(t, err)

			assert.Equal(t, originalPrice, price)
		})
		t.Run("BenchmarkPriceFromReport fails", func(t *testing.T) {
			codec.val = nil
			codec.err = errors.New("something exploded")

			_, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
			require.Error(t, err)

			assert.EqualError(t, err, "something exploded")
		})
	})

	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		price, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.NoError(t, err)

		assert.Nil(t, price)
	})

	t.Run("failing query", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		_, err := mt.LatestPrice(testutils.Context(t), sampleFeedID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
}

func Test_MercuryTransmitter_FetchInitialMaxFinalizedBlockNumber(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	var jobID int32
	codec := new(mockCodec)
	orm := NewORM(db)
	clients := map[string]wsrpc.Client{}

	t.Run("successful query", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.FeedId = sampleFeedID[:]
				out.Report.CurrentBlockNumber = 42
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		bn, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.NoError(t, err)

		require.NotNil(t, bn)
		assert.Equal(t, 42, int(*bn))
	})
	t.Run("successful query returning nil report (new feed)", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				out = new(pb.LatestReportResponse)
				out.Report = nil
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		bn, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.NoError(t, err)

		assert.Nil(t, bn)
	})
	t.Run("failing query", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				return nil, errors.New("something exploded")
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})
	t.Run("return feed ID is wrong", func(t *testing.T) {
		c := &mocks.MockWSRPCClient{
			LatestReportF: func(ctx context.Context, in *pb.LatestReportRequest) (out *pb.LatestReportResponse, err error) {
				require.NotNil(t, in)
				assert.Equal(t, hexutil.Encode(sampleFeedID[:]), hexutil.Encode(in.FeedId))
				out = new(pb.LatestReportResponse)
				out.Report = new(pb.Report)
				out.Report.CurrentBlockNumber = 42
				out.Report.FeedId = []byte{1, 2}
				return out, nil
			},
		}
		clients[sURL] = c
		mt := NewTransmitter(lggr, mockCfg{}, clients, sampleClientPubKey, jobID, sampleFeedID, orm, codec, nil)
		_, err := mt.FetchInitialMaxFinalizedBlockNumber(testutils.Context(t))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "latestReport failed; mismatched feed IDs, expected: 0x1c916b4aa7e57ca7b68ae1bf45653f56b656fd3aa335ef7fae696b663f1b8472, got: 0x")
	})
}

func Test_sortReportsLatestFirst(t *testing.T) {
	reports := []*pb.Report{
		nil,
		{ObservationsTimestamp: 1},
		{ObservationsTimestamp: 1},
		{ObservationsTimestamp: 2},
		{CurrentBlockNumber: 1},
		nil,
		{CurrentBlockNumber: 2},
		{},
	}

	sortReportsLatestFirst(reports)

	assert.Equal(t, int64(2), reports[0].ObservationsTimestamp)
	assert.Equal(t, int64(1), reports[1].ObservationsTimestamp)
	assert.Equal(t, int64(1), reports[2].ObservationsTimestamp)
	assert.Equal(t, int64(0), reports[3].ObservationsTimestamp)
	assert.Equal(t, int64(2), reports[3].CurrentBlockNumber)
	assert.Equal(t, int64(0), reports[4].ObservationsTimestamp)
	assert.Equal(t, int64(1), reports[4].CurrentBlockNumber)
	assert.Equal(t, int64(0), reports[5].ObservationsTimestamp)
	assert.Equal(t, int64(0), reports[5].CurrentBlockNumber)
	assert.Nil(t, reports[6])
	assert.Nil(t, reports[7])
}

type mockQ struct {
	ch chan *Transmission
}

func newMockQ() *mockQ {
	return &mockQ{make(chan *Transmission, 100)}
}

func (m *mockQ) Start(context.Context) error { return nil }
func (m *mockQ) Close() error {
	m.ch <- nil
	return nil
}
func (m *mockQ) Ready() error                   { return nil }
func (m *mockQ) HealthReport() map[string]error { return nil }
func (m *mockQ) Name() string                   { return "" }
func (m *mockQ) BlockingPop() (t *Transmission) {
	val := <-m.ch
	return val
}
func (m *mockQ) Push(req *pb.TransmitRequest, reportCtx ocrtypes.ReportContext) (ok bool) {
	m.ch <- &Transmission{Req: req, ReportCtx: reportCtx}
	return true
}
func (m *mockQ) Init(transmissions []*Transmission) {}
func (m *mockQ) IsEmpty() bool                      { return false }

func Test_MercuryTransmitter_runQueueLoop(t *testing.T) {
	feedIDHex := utils.NewHash().Hex()
	lggr := logger.TestLogger(t)
	c := &mocks.MockWSRPCClient{}
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db)
	pm := NewPersistenceManager(lggr, sURL, orm, 0, 0, 0, 0)
	cfg := mockCfg{}

	s := newServer(lggr, cfg, c, pm, sURL, feedIDHex)

	req := &pb.TransmitRequest{
		Payload:      []byte{1, 2, 3},
		ReportFormat: 32,
	}

	t.Run("pulls from queue and transmits successfully", func(t *testing.T) {
		transmit := make(chan *pb.TransmitRequest, 1)
		c.TransmitF = func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
			transmit <- in
			return &pb.TransmitResponse{Code: 0, Error: ""}, nil
		}
		q := newMockQ()
		s.q = q
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go s.runQueueLoop(nil, wg, feedIDHex)

		q.Push(req, sampleReportContext)

		select {
		case tr := <-transmit:
			assert.Equal(t, []byte{1, 2, 3}, tr.Payload)
			assert.Equal(t, 32, int(tr.ReportFormat))
			// case <-time.After(testutils.WaitTimeout(t)):
		case <-time.After(1 * time.Second):
			t.Fatal("expected a transmit request to be sent")
		}

		q.Close()
		wg.Wait()
	})

	t.Run("on duplicate, success", func(t *testing.T) {
		transmit := make(chan *pb.TransmitRequest, 1)
		c.TransmitF = func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
			transmit <- in
			return &pb.TransmitResponse{Code: DuplicateReport, Error: ""}, nil
		}
		q := newMockQ()
		s.q = q
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go s.runQueueLoop(nil, wg, feedIDHex)

		q.Push(req, sampleReportContext)

		select {
		case tr := <-transmit:
			assert.Equal(t, []byte{1, 2, 3}, tr.Payload)
			assert.Equal(t, 32, int(tr.ReportFormat))
			// case <-time.After(testutils.WaitTimeout(t)):
		case <-time.After(1 * time.Second):
			t.Fatal("expected a transmit request to be sent")
		}

		q.Close()
		wg.Wait()
	})
	t.Run("on server-side error, does not retry", func(t *testing.T) {
		transmit := make(chan *pb.TransmitRequest, 1)
		c.TransmitF = func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
			transmit <- in
			return &pb.TransmitResponse{Code: DuplicateReport, Error: ""}, nil
		}
		q := newMockQ()
		s.q = q
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go s.runQueueLoop(nil, wg, feedIDHex)

		q.Push(req, sampleReportContext)

		select {
		case tr := <-transmit:
			assert.Equal(t, []byte{1, 2, 3}, tr.Payload)
			assert.Equal(t, 32, int(tr.ReportFormat))
			// case <-time.After(testutils.WaitTimeout(t)):
		case <-time.After(1 * time.Second):
			t.Fatal("expected a transmit request to be sent")
		}

		q.Close()
		wg.Wait()
	})
	t.Run("on transmit error, retries", func(t *testing.T) {
		transmit := make(chan *pb.TransmitRequest, 1)
		c.TransmitF = func(ctx context.Context, in *pb.TransmitRequest) (*pb.TransmitResponse, error) {
			transmit <- in
			return &pb.TransmitResponse{}, errors.New("transmission error")
		}
		q := newMockQ()
		s.q = q
		wg := &sync.WaitGroup{}
		wg.Add(1)
		stopCh := make(chan struct{}, 1)

		go s.runQueueLoop(stopCh, wg, feedIDHex)

		q.Push(req, sampleReportContext)

		cnt := 0
	Loop:
		for {
			select {
			case tr := <-transmit:
				assert.Equal(t, []byte{1, 2, 3}, tr.Payload)
				assert.Equal(t, 32, int(tr.ReportFormat))
				if cnt > 2 {
					break Loop
				}
				cnt++
				// case <-time.After(testutils.WaitTimeout(t)):
			case <-time.After(1 * time.Second):
				t.Fatal("expected 3 transmit requests to be sent")
			}
		}

		close(stopCh)
		wg.Wait()
	})
}
