package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/sqlx"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
	"github.com/smartcontractkit/chainlink-relay/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	maxTransmitQueueSize = 10_000
	transmitTimeout      = 5 * time.Second
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
		[]string{"feedID"},
	)
	transmitDuplicateCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_duplicate_count",
		Help: "Number of transmissions where the server told us it was a duplicate",
	},
		[]string{"feedID"},
	)
	transmitConnectionErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_connection_error_count",
		Help: "Number of errored transmissions that failed due to problem with the connection",
	},
		[]string{"feedID"},
	)
	transmitServerErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_transmit_server_error_count",
		Help: "Number of errored transmissions that failed due to an error returned by the mercury server",
	},
		[]string{"feedID", "code"},
	)
)

type Transmitter interface {
	relaymercury.Transmitter
	services.Service
}

type ConfigTracker interface {
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

type TransmitterReportDecoder interface {
	BenchmarkPriceFromReport(report ocrtypes.Report) (*big.Int, error)
}

var _ Transmitter = (*mercuryTransmitter)(nil)

type mercuryTransmitter struct {
	utils.StartStopOnce
	lggr               logger.Logger
	rpcClient          wsrpc.Client
	cfgTracker         ConfigTracker
	persistenceManager *PersistenceManager
	codec              TransmitterReportDecoder

	feedID      mercuryutils.FeedID
	jobID       int32
	fromAccount string

	stopCh utils.StopChan
	queue  *TransmitQueue
	wg     sync.WaitGroup

	transmitSuccessCount         prometheus.Counter
	transmitDuplicateCount       prometheus.Counter
	transmitConnectionErrorCount prometheus.Counter
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

func NewTransmitter(lggr logger.Logger, cfgTracker ConfigTracker, rpcClient wsrpc.Client, fromAccount ed25519.PublicKey, jobID int32, feedID [32]byte, db *sqlx.DB, cfg pg.QConfig, codec TransmitterReportDecoder) *mercuryTransmitter {
	feedIDHex := fmt.Sprintf("0x%x", feedID[:])
	persistenceManager := NewPersistenceManager(lggr, NewORM(db, lggr, cfg), jobID, maxTransmitQueueSize, flushDeletesFrequency, pruneFrequency)
	return &mercuryTransmitter{
		utils.StartStopOnce{},
		lggr.Named("MercuryTransmitter").With("feedID", feedIDHex),
		rpcClient,
		cfgTracker,
		persistenceManager,
		codec,
		feedID,
		jobID,
		fmt.Sprintf("%x", fromAccount),
		make(chan (struct{})),
		NewTransmitQueue(lggr, feedIDHex, maxTransmitQueueSize, nil, persistenceManager),
		sync.WaitGroup{},
		transmitSuccessCount.WithLabelValues(feedIDHex),
		transmitDuplicateCount.WithLabelValues(feedIDHex),
		transmitConnectionErrorCount.WithLabelValues(feedIDHex),
	}
}

func (mt *mercuryTransmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("MercuryTransmitter", func() error {
		mt.lggr.Debugw("Loading transmit requests from database")
		if err := mt.persistenceManager.Start(ctx); err != nil {
			return err
		}
		transmissions, err := mt.persistenceManager.Load(ctx)
		if err != nil {
			return err
		}
		mt.queue = NewTransmitQueue(mt.lggr, mt.feedID.String(), maxTransmitQueueSize, transmissions, mt.persistenceManager)

		if err := mt.rpcClient.Start(ctx); err != nil {
			return err
		}
		if err := mt.queue.Start(ctx); err != nil {
			return err
		}
		mt.wg.Add(1)
		go mt.runQueueLoop()
		return nil
	})
}

func (mt *mercuryTransmitter) Close() error {
	return mt.StopOnce("MercuryTransmitter", func() error {
		if err := mt.queue.Close(); err != nil {
			return err
		}
		if err := mt.persistenceManager.Close(); err != nil {
			return err
		}
		close(mt.stopCh)
		mt.wg.Wait()
		return mt.rpcClient.Close()
	})
}

func (mt *mercuryTransmitter) Name() string { return mt.lggr.Name() }

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	report := map[string]error{mt.Name(): mt.Healthy()}
	services.CopyHealth(report, mt.rpcClient.HealthReport())
	services.CopyHealth(report, mt.queue.HealthReport())
	return report
}

func (mt *mercuryTransmitter) runQueueLoop() {
	defer mt.wg.Done()
	// Exponential backoff with very short retry interval (since latency is a priority)
	// 5ms, 10ms, 20ms, 40ms etc
	b := backoff.Backoff{
		Min:    5 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: true,
	}
	runloopCtx, cancel := mt.stopCh.Ctx(context.Background())
	defer cancel()
	for {
		t := mt.queue.BlockingPop()
		if t == nil {
			// queue was closed
			return
		}
		ctx, cancel := context.WithTimeout(runloopCtx, utils.WithJitter(transmitTimeout))
		res, err := mt.rpcClient.Transmit(ctx, t.Req)
		cancel()
		if runloopCtx.Err() != nil {
			// runloop context is only canceled on transmitter close so we can
			// exit the runloop here
			return
		} else if err != nil {
			mt.transmitConnectionErrorCount.Inc()
			mt.lggr.Errorw("Transmit report failed", "err", err, "reportCtx", t.ReportCtx)
			if ok := mt.queue.Push(t.Req, t.ReportCtx); !ok {
				mt.lggr.Error("Failed to push report to transmit queue; queue is closed")
				return
			}
			// Wait a backoff duration before pulling the most recent transmission
			// the heap
			select {
			case <-time.After(b.Duration()):
				continue
			case <-mt.stopCh:
				return
			}
		}

		b.Reset()
		if res.Error == "" {
			mt.transmitSuccessCount.Inc()
			mt.lggr.Tracew("Transmit report success", "req", t.Req, "response", res, "reportCtx", t.ReportCtx)
		} else {
			// We don't need to retry here because the mercury server
			// has confirmed it received the report. We only need to retry
			// on networking/unknown errors
			switch res.Code {
			case DuplicateReport:
				mt.transmitSuccessCount.Inc()
				mt.transmitDuplicateCount.Inc()
				mt.lggr.Tracew("Transmit report succeeded; duplicate report", "code", res.Code)
			default:
				transmitServerErrorCount.WithLabelValues(mt.feedID.String(), fmt.Sprintf("%d", res.Code)).Inc()
				mt.lggr.Errorw("Transmit report failed; mercury server returned error", "response", res, "reportCtx", t.ReportCtx, "err", res.Error, "code", res.Code)
			}
		}

		if err := mt.persistenceManager.Delete(runloopCtx, t.Req); err != nil {
			mt.lggr.Errorw("Failed to delete transmit request record", "error", err, "reportCtx", t.ReportCtx)
			return
		}
	}
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

	payload, err := PayloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return pkgerrors.Wrap(err, "abi.Pack failed")
	}

	req := &pb.TransmitRequest{
		Payload: payload,
	}

	mt.lggr.Tracew("Transmit enqueue", "req", req, "report", report, "reportCtx", reportCtx, "signatures", signatures)

	if err := mt.persistenceManager.Insert(ctx, req, reportCtx); err != nil {
		return err
	}
	if ok := mt.queue.Push(req, reportCtx); !ok {
		return errors.New("transmit queue is closed")
	}
	return nil
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
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		mt.lggr.Warnw("latestReport failed", "err", err)
		return nil, pkgerrors.Wrap(err, "latestReport failed")
	}
	if resp == nil {
		return nil, errors.New("latestReport expected non-nil response")
	}
	if resp.Error != "" {
		err = errors.New(resp.Error)
		mt.lggr.Warnw("latestReport failed; mercury server returned error", "err", err)
		return nil, err
	}
	if resp.Report == nil {
		mt.lggr.Tracew("latestReport success: returned nil")
		return nil, nil
	} else if !bytes.Equal(resp.Report.FeedId, feedID[:]) {
		err = fmt.Errorf("latestReport failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID[:], resp.Report.FeedId[:])
		mt.lggr.Errorw("latestReport failed", "err", err)
		return nil, err
	}

	mt.lggr.Tracew("latestReport success", "currentBlockNum", resp.Report.CurrentBlockNumber)

	return resp.Report, nil
}
