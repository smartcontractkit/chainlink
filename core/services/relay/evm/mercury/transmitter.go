package mercury

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jpillora/backoff"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	MaxTransmitQueueSize = 10_000
	TransmitTimeout      = 5 * time.Second
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
	services.ServiceCtx
}

type ConfigTracker interface {
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error)
}

var _ Transmitter = &mercuryTransmitter{}

type mercuryTransmitter struct {
	utils.StartStopOnce
	lggr               logger.Logger
	rpcClient          wsrpc.Client
	cfgTracker         ConfigTracker
	initialBlockNumber int64

	feedID      [32]byte
	feedIDHex   string
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

func NewTransmitter(lggr logger.Logger, cfgTracker ConfigTracker, rpcClient wsrpc.Client, fromAccount ed25519.PublicKey, feedID [32]byte, initialBlockNumber int64) *mercuryTransmitter {
	feedIDHex := fmt.Sprintf("0x%x", feedID[:])
	return &mercuryTransmitter{
		utils.StartStopOnce{},
		lggr.Named("MercuryTransmitter").With("feedID", feedIDHex),
		rpcClient,
		cfgTracker,
		initialBlockNumber,
		feedID,
		feedIDHex,
		fmt.Sprintf("%x", fromAccount),
		make(chan (struct{})),
		NewTransmitQueue(lggr, feedIDHex, MaxTransmitQueueSize),
		sync.WaitGroup{},
		transmitSuccessCount.WithLabelValues(feedIDHex),
		transmitDuplicateCount.WithLabelValues(feedIDHex),
		transmitConnectionErrorCount.WithLabelValues(feedIDHex),
	}
}

func (mt *mercuryTransmitter) Start(ctx context.Context) (err error) {
	return mt.StartOnce("MercuryTransmitter", func() error {
		if err := mt.rpcClient.Start(ctx); err != nil {
			return err
		}
		if err := mt.queue.Start(ctx); err != nil {
			return err
		}
		mt.wg.Add(1)
		go mt.runloop()
		return nil
	})
}

func (mt *mercuryTransmitter) Close() error {
	return mt.StopOnce("MercuryTransmitter", func() error {
		mt.queue.Close()
		close(mt.stopCh)
		mt.wg.Wait()
		return mt.rpcClient.Close()
	})
}
func (mt *mercuryTransmitter) Ready() error { return mt.StartStopOnce.Ready() }
func (mt *mercuryTransmitter) Name() string { return mt.lggr.Name() }

func (mt *mercuryTransmitter) HealthReport() map[string]error {
	report := map[string]error{mt.Name(): mt.StartStopOnce.Healthy()}
	maps.Copy(report, mt.rpcClient.HealthReport())
	maps.Copy(report, mt.queue.HealthReport())
	return report
}

func (mt *mercuryTransmitter) runloop() {
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
		ctx, cancel := context.WithTimeout(runloopCtx, utils.WithJitter(TransmitTimeout))
		res, err := mt.rpcClient.Transmit(ctx, t.Req)
		cancel()
		if runloopCtx.Err() != nil {
			// runloop context is only canceled on transmitter close so we can
			// exit the runloop here
			return
		} else if err != nil {
			mt.transmitConnectionErrorCount.Inc()
			mt.lggr.Errorw("Transmit report failed", "req", t.Req, "error", err, "reportCtx", t.ReportCtx)
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
				elems := map[string]interface{}{}
				var validFrom int64
				var currentBlock int64
				var unpackErr error
				if err = PayloadTypes.UnpackIntoMap(elems, t.Req.Payload); err != nil {
					unpackErr = err
				} else {
					report := elems["report"].([]byte)
					validFrom, err = (&reportcodec.EVMReportCodec{}).ValidFromBlockNumFromReport(report)
					if err != nil {
						unpackErr = err
					}
					currentBlock, err = (&reportcodec.EVMReportCodec{}).CurrentBlockNumFromReport(report)
					if err != nil {
						unpackErr = errors.Join(unpackErr, err)
					}
				}
				transmitServerErrorCount.WithLabelValues(mt.feedIDHex, fmt.Sprintf("%d", res.Code)).Inc()
				mt.lggr.Errorw("Transmit report failed; mercury server returned error", "unpackErr", unpackErr, "validFromBlock", validFrom, "currentBlock", currentBlock, "req", t.Req, "response", res, "reportCtx", t.ReportCtx, "err", res.Error, "code", res.Code)
			}
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

	mt.lggr.Debugw("Transmit enqueue", "req", req, "report", report, "reportCtx", reportCtx, "signatures", signatures)

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

func (mt *mercuryTransmitter) FetchInitialMaxFinalizedBlockNumber(ctx context.Context) (int64, error) {
	mt.lggr.Debug("FetchInitialMaxFinalizedBlockNumber")
	req := &pb.LatestReportRequest{
		FeedId: mt.feedID[:],
	}
	resp, err := mt.rpcClient.LatestReport(ctx, req)
	if err != nil {
		mt.lggr.Errorw("FetchInitialMaxFinalizedBlockNumber failed", "err", err)
		return 0, pkgerrors.Wrap(err, "FetchInitialMaxFinalizedBlockNumber failed to fetch LatestReport")
	}
	if resp == nil {
		return 0, errors.New("FetchInitialMaxFinalizedBlockNumber expected LatestReport to return non-nil response")
	}
	if resp.Error != "" {
		err = errors.New(resp.Error)
		mt.lggr.Errorw("FetchInitialMaxFinalizedBlockNumber failed; mercury server returned error", "err", err)
		return 0, err
	}
	if resp.Report == nil {
		maxFinalizedBlockNumber := mt.initialBlockNumber - 1
		mt.lggr.Infof("FetchInitialMaxFinalizedBlockNumber returned empty LatestReport; this is a new feed so maxFinalizedBlockNumber=%d (initialBlockNumber=%d)", maxFinalizedBlockNumber, mt.initialBlockNumber)
		// NOTE: It's important to return -1 if the server is missing any past
		// report (brand new feed) since we will add 1 to the
		// maxFinalizedBlockNumber to get the first validFromBlockNum, which
		// ought to be zero.
		//
		// If "initialBlockNumber" is unset, this will give a starting block of zero.
		return maxFinalizedBlockNumber, nil
	} else if !bytes.Equal(resp.Report.FeedId, mt.feedID[:]) {
		return 0, fmt.Errorf("FetchInitialMaxFinalizedBlockNumber failed; mismatched feed IDs, expected: 0x%x, got: 0x%x", mt.feedID, resp.Report.FeedId)
	}

	mt.lggr.Debugw("FetchInitialMaxFinalizedBlockNumber success", "currentBlockNum", resp.Report.CurrentBlockNumber)

	return resp.Report.CurrentBlockNumber, nil
}
