package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capMercury "github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/triggers"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	// Mercury server error codes
	DuplicateReport = 2
)

var (
	transmitSuccessCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_success_count",
		Help: "Number of successful transmissions (duplicates are counted as success)",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitDuplicateCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_duplicate_count",
		Help: "Number of transmissions where the server told us it was a duplicate",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitConnectionErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_connection_error_count",
		Help: "Number of errored transmissions that failed due to problem with the connection",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitQueueDeleteErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_queue_delete_error_count",
		Help: "Running count of DB errors when trying to delete an item from the queue DB",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitQueueInsertErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_queue_insert_error_count",
		Help: "Running count of DB errors when trying to insert an item into the queue DB",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitQueuePushErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_queue_push_error_count",
		Help: "Running count of DB errors when trying to push an item onto the queue",
	},
		[]string{"feedID", "serverURL"},
	)
	transmitServerErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_server_error_count",
		Help: "Number of errored transmissions that failed due to an error returned by the mercury server",
	},
		[]string{"feedID", "serverURL", "code"},
	)
)

type Transmitter interface {
	mercury.Transmitter
	services.Service
}

type ConfigTracker interface {
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

type TransmitterReportDecoder interface {
	BenchmarkPriceFromReport(report ocrtypes.Report) (*big.Int, error)
	ObservationTimestampFromReport(report ocrtypes.Report) (uint32, error)
}

var _ Transmitter = (*mercuryTransmitter)(nil)

type TransmitterConfig interface {
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

type mercuryTransmitter struct {
	services.StateMachine
	lggr logger.Logger
	cfg  TransmitterConfig

	orm     ORM
	servers map[string]*server

	codec             TransmitterReportDecoder
	triggerCapability *triggers.MercuryTriggerService

	feedID      mercuryutils.FeedID
	jobID       int32
	fromAccount string

	stopCh services.StopChan
	wg     *sync.WaitGroup
}

var PayloadTypes = getPayloadTypes()

func getPayloadTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "reportContext", Type: mustNewType("bytes32[3]")},
		{Name: "report", Type: mustNewType("bytes")},
		{Name: "rawRs", Type: mustNewType("bytes32[]")},
		{Name: "rawSs", Type: mustNewType("bytes32[]")},
		{Name: "rawVs", Type: mustNewType("bytes32")},
	})
}

type server struct {
	lggr logger.Logger

	transmitTimeout time.Duration

	c  wsrpc.Client
	pm *PersistenceManager
	q  *TransmitQueue

	deleteQueue chan *pb.TransmitRequest

	transmitSuccessCount          prometheus.Counter
	transmitDuplicateCount        prometheus.Counter
	transmitConnectionErrorCount  prometheus.Counter
	transmitQueueDeleteErrorCount prometheus.Counter
	transmitQueueInsertErrorCount prometheus.Counter
	transmitQueuePushErrorCount   prometheus.Counter
}

func (s *server) HealthReport() map[string]error {
	report := map[string]error{}
	services.CopyHealth(report, s.c.HealthReport())
	services.CopyHealth(report, s.q.HealthReport())
	return report
}

func (s *server) runDeleteQueueLoop(stopCh services.StopChan, wg *sync.WaitGroup) {
	defer wg.Done()
	runloopCtx, cancel := stopCh.Ctx(context.Background())
	defer cancel()

	// Exponential backoff for very rarely occurring errors (DB disconnect etc)
	b := backoff.Backoff{
		Min:    1 * time.Second,
		Max:    120 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	for {
		select {
		case req := <-s.deleteQueue:
			for {
				if err := s.pm.Delete(runloopCtx, req); err != nil {
					s.lggr.Errorw("Failed to delete transmit request record", "err", err, "req.Payload", req.Payload)
					s.transmitQueueDeleteErrorCount.Inc()
					select {
					case <-time.After(b.Duration()):
						// Wait a backoff duration before trying to delete again
						continue
					case <-stopCh:
						// abort and return immediately on stop even if items remain in queue
						return
					}
				}
				break
			}
			// success
			b.Reset()
		case <-stopCh:
			// abort and return immediately on stop even if items remain in queue
			return
		}
	}
}

func (s *server) runQueueLoop(stopCh services.StopChan, wg *sync.WaitGroup, feedIDHex string) {
	defer wg.Done()
	// Exponential backoff with very short retry interval (since latency is a priority)
	// 5ms, 10ms, 20ms, 40ms etc
	b := backoff.Backoff{
		Min:    5 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: true,
	}
	runloopCtx, cancel := stopCh.Ctx(context.Background())
	defer cancel()
	for {
		t := s.q.BlockingPop()
		if t == nil {
			// queue was closed
			return
		}
		ctx, cancel := context.WithTimeout(runloopCtx, utils.WithJitter(s.transmitTimeout))
		res, err := s.c.Transmit(ctx, t.Req)
		cancel()
		if runloopCtx.Err() != nil {
			// runloop context is only canceled on transmitter close so we can
			// exit the runloop here
			return
		} else if err != nil {
			s.transmitConnectionErrorCount.Inc()
			s.lggr.Errorw("Transmit report failed", "err", err, "reportCtx", t.ReportCtx)
			if ok := s.q.Push(t.Req, t.ReportCtx); !ok {
				s.lggr.Error("Failed to push report to transmit queue; queue is closed")
				return
			}
			// Wait a backoff duration before pulling the most recent transmission
			// the heap
			select {
			case <-time.After(b.Duration()):
				continue
			case <-stopCh:
				return
			}
		}

		b.Reset()
		if res.Error == "" {
			s.transmitSuccessCount.Inc()
			s.lggr.Debugw("Transmit report success", "payload", hexutil.Encode(t.Req.Payload), "response", res, "repts", t.ReportCtx.ReportTimestamp)
		} else {
			// We don't need to retry here because the mercury server
			// has confirmed it received the report. We only need to retry
			// on networking/unknown errors
			switch res.Code {
			case DuplicateReport:
				s.transmitSuccessCount.Inc()
				s.transmitDuplicateCount.Inc()
				s.lggr.Debugw("Transmit report success; duplicate report", "payload", hexutil.Encode(t.Req.Payload), "response", res, "repts", t.ReportCtx.ReportTimestamp)
			default:
				transmitServerErrorCount.WithLabelValues(feedIDHex, fmt.Sprintf("%d", res.Code)).Inc()
				s.lggr.Errorw("Transmit report failed; mercury server returned error", "response", res, "reportCtx", t.ReportCtx, "err", res.Error, "code", res.Code)
			}
		}

		select {
		case s.deleteQueue <- t.Req:
		default:
			s.lggr.Criticalw("Delete queue is full", "reportCtx", t.ReportCtx)
		}
	}
}

func NewTransmitter(lggr logger.Logger, cfg TransmitterConfig, clients map[string]wsrpc.Client, fromAccount ed25519.PublicKey, jobID int32, feedID [32]byte, orm ORM, codec TransmitterReportDecoder, triggerCapability *triggers.MercuryTriggerService) *mercuryTransmitter {
	feedIDHex := fmt.Sprintf("0x%x", feedID[:])
	servers := make(map[string]*server, len(clients))
	for serverURL, client := range clients {
		cLggr := lggr.Named(serverURL).With("serverURL", serverURL)
		pm := NewPersistenceManager(cLggr, serverURL, orm, jobID, int(cfg.TransmitQueueMaxSize()), flushDeletesFrequency, pruneFrequency)
		servers[serverURL] = &server{
			cLggr,
			cfg.TransmitTimeout().Duration(),
			client,
			pm,
			NewTransmitQueue(cLggr, serverURL, feedIDHex, int(cfg.TransmitQueueMaxSize()), pm),
			make(chan *pb.TransmitRequest, int(cfg.TransmitQueueMaxSize())),
			transmitSuccessCount.WithLabelValues(feedIDHex, serverURL),
			transmitDuplicateCount.WithLabelValues(feedIDHex, serverURL),
			transmitConnectionErrorCount.WithLabelValues(feedIDHex, serverURL),
			transmitQueueDeleteErrorCount.WithLabelValues(feedIDHex, serverURL),
			transmitQueueInsertErrorCount.WithLabelValues(feedIDHex, serverURL),
			transmitQueuePushErrorCount.WithLabelValues(feedIDHex, serverURL),
		}
	}
	return &mercuryTransmitter{
		services.StateMachine{},
		lggr.Named("MercuryTransmitter").With("feedID", feedIDHex),
		cfg,
		orm,
		servers,
		codec,
		triggerCapability,
		feedID,
		jobID,
		fmt.Sprintf("%x", fromAccount),
		make(services.StopChan),
		&sync.WaitGroup{},
	}
}

func (mt *mercuryTransmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("MercuryTransmitter", func() error {
		mt.lggr.Debugw("Loading transmit requests from database")

		{
			var startClosers []services.StartClose
			for _, s := range mt.servers {
				transmissions, err := s.pm.Load(ctx)
				if err != nil {
					return err
				}
				s.q.Init(transmissions)
				// starting pm after loading from it is fine because it simply spawns some garbage collection/prune goroutines
				startClosers = append(startClosers, s.c, s.q, s.pm)

				mt.wg.Add(2)
				go s.runDeleteQueueLoop(mt.stopCh, mt.wg)
				go s.runQueueLoop(mt.stopCh, mt.wg, mt.feedID.Hex())
			}
			if err := (&services.MultiStart{}).Start(ctx, startClosers...); err != nil {
				return err
			}
		}

		return nil
	})
}

func (mt *mercuryTransmitter) Close() error {
	return mt.StopOnce("MercuryTransmitter", func() error {
		// Drain all the queues first
		var qs []io.Closer
		for _, s := range mt.servers {
			qs = append(qs, s.q)
		}
		if err := services.CloseAll(qs...); err != nil {
			return err
		}

		close(mt.stopCh)
		mt.wg.Wait()

		// Close all the persistence managers
		// Close all the clients
		var closers []io.Closer
		for _, s := range mt.servers {
			closers = append(closers, s.pm)
			closers = append(closers, s.c)
		}
		return services.CloseAll(closers...)
	})
}

func (mt *mercuryTransmitter) Name() string { return mt.lggr.Name() }

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	report := map[string]error{mt.Name(): mt.Healthy()}
	for _, s := range mt.servers {
		services.CopyHealth(report, s.HealthReport())
	}
	return report
}

func (mt *mercuryTransmitter) sendToTrigger(report ocrtypes.Report, rs [][32]byte, ss [][32]byte, vs [32]byte) error {
	var rsUnsized [][]byte
	var ssUnsized [][]byte
	for idx := range rs {
		rsUnsized = append(rsUnsized, rs[idx][:])
		ssUnsized = append(ssUnsized, ss[idx][:])
	}
	converted := capMercury.FeedReport{
		FeedID:     mt.feedID.Hex(),
		FullReport: report,
		Rs:         rsUnsized,
		Ss:         ssUnsized,
		Vs:         vs[:],
		// NOTE: Skipping fields derived from FullReport, they will be filled out at a later stage
		// after decoding and validating signatures.
	}
	return mt.triggerCapability.ProcessReport([]capMercury.FeedReport{converted})
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (mt *mercuryTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range signatures {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	if mt.triggerCapability != nil {
		// Acting as a Capability - send report to trigger service and exit.
		return mt.sendToTrigger(report, rs, ss, vs)
	}

	payload, err := PayloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return pkgerrors.Wrap(err, "abi.Pack failed")
	}

	req := &pb.TransmitRequest{
		Payload: payload,
	}

	ts, err := mt.codec.ObservationTimestampFromReport(report)
	if err != nil {
		mt.lggr.Warnw("Failed to get observation timestamp from report", "err", err)
	}
	mt.lggr.Debugw("Transmit enqueue", "req.Payload", hexutil.Encode(req.Payload), "report", report, "repts", reportCtx.ReportTimestamp, "signatures", signatures, "observationsTimestamp", ts)

	if err := mt.orm.InsertTransmitRequest(ctx, maps.Keys(mt.servers), req, mt.jobID, reportCtx); err != nil {
		return err
	}

	g := new(errgroup.Group)
	for _, s := range mt.servers {
		s := s // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			if ok := s.q.Push(req, reportCtx); !ok {
				s.transmitQueuePushErrorCount.Inc()
				return errors.New("transmit queue is closed")
			}
			return nil
		})
	}

	return g.Wait()
}

// FromAccount returns the stringified (hex) CSA public key
func (mt *mercuryTransmitter) FromAccount() (ocrtypes.Account, error) {
	return ocrtypes.Account(mt.fromAccount), nil
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
func (mt *mercuryTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (cd ocrtypes.ConfigDigest, epoch uint32, err error) {
	panic("not needed for OCR3")
}

func (mt *mercuryTransmitter) FetchInitialMaxFinalizedBlockNumber(ctx context.Context) (*int64, error) {
	mt.lggr.Trace("FetchInitialMaxFinalizedBlockNumber")

	report, err := mt.latestReport(ctx, mt.feedID)
	if err != nil {
		return nil, err
	}

	if report == nil {
		mt.lggr.Debugw("FetchInitialMaxFinalizedBlockNumber success; got nil report")
		return nil, nil
	}

	mt.lggr.Debugw("FetchInitialMaxFinalizedBlockNumber success", "currentBlockNum", report.CurrentBlockNumber)

	return &report.CurrentBlockNumber, nil
}

func (mt *mercuryTransmitter) LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error) {
	mt.lggr.Trace("LatestPrice")

	fullReport, err := mt.latestReport(ctx, feedID)
	if err != nil {
		return nil, err
	}
	if fullReport == nil {
		return nil, nil
	}
	payload := fullReport.Payload
	m := make(map[string]interface{})
	if err := PayloadTypes.UnpackIntoMap(m, payload); err != nil {
		return nil, err
	}
	report, is := m["report"].([]byte)
	if !is {
		return nil, fmt.Errorf("expected report to be []byte, but it was %T", m["report"])
	}
	return mt.codec.BenchmarkPriceFromReport(report)
}

// LatestTimestamp will return -1, nil if the feed is missing
func (mt *mercuryTransmitter) LatestTimestamp(ctx context.Context) (int64, error) {
	mt.lggr.Trace("LatestTimestamp")

	report, err := mt.latestReport(ctx, mt.feedID)
	if err != nil {
		return 0, err
	}

	if report == nil {
		mt.lggr.Debugw("LatestTimestamp success; got nil report")
		return -1, nil
	}

	mt.lggr.Debugw("LatestTimestamp success", "timestamp", report.ObservationsTimestamp)

	return report.ObservationsTimestamp, nil
}

func (mt *mercuryTransmitter) latestReport(ctx context.Context, feedID [32]byte) (*pb.Report, error) {
	mt.lggr.Trace("latestReport")

	req := &pb.LatestReportRequest{
		FeedId: feedID[:],
	}

	var reports []*pb.Report
	mu := sync.Mutex{}
	var g errgroup.Group
	for _, s := range mt.servers {
		s := s
		g.Go(func() error {
			resp, err := s.c.LatestReport(ctx, req)
			if err != nil {
				s.lggr.Warnw("latestReport failed", "err", err)
				return err
			}
			if resp == nil {
				err = errors.New("latestReport expected non-nil response from server")
				s.lggr.Warn(err.Error())
				return err
			}
			if resp.Error != "" {
				err = errors.New(resp.Error)
				s.lggr.Warnw("latestReport failed; mercury server returned error", "err", err)
				return fmt.Errorf("latestReport failed; mercury server returned error: %s", resp.Error)
			}
			if resp.Report == nil {
				s.lggr.Tracew("latestReport success: returned nil")
			} else if !bytes.Equal(resp.Report.FeedId, feedID[:]) {
				err = fmt.Errorf("latestReport failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID[:], resp.Report.FeedId[:])
				s.lggr.Errorw("latestReport failed", "err", err)
				return err
			} else {
				s.lggr.Tracew("latestReport success", "observationsTimestamp", resp.Report.ObservationsTimestamp, "currentBlockNum", resp.Report.CurrentBlockNumber)
			}
			mu.Lock()
			defer mu.Unlock()
			reports = append(reports, resp.Report)
			return nil
		})
	}
	err := g.Wait()

	if len(reports) == 0 {
		return nil, fmt.Errorf("latestReport failed; all servers returned an error: %w", err)
	}

	sortReportsLatestFirst(reports)

	return reports[0], nil
}

func sortReportsLatestFirst(reports []*pb.Report) {
	sort.Slice(reports, func(i, j int) bool {
		// nils are "earliest" so they go to the end
		if reports[i] == nil {
			return false
		} else if reports[j] == nil {
			return true
		}
		// Handle block number case
		if reports[i].ObservationsTimestamp == reports[j].ObservationsTimestamp {
			return reports[i].CurrentBlockNumber > reports[j].CurrentBlockNumber
		}
		// Timestamp case
		return reports[i].ObservationsTimestamp > reports[j].ObservationsTimestamp
	})
}
