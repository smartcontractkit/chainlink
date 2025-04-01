package mercurytransmitter

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jpillora/backoff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"
	"github.com/smartcontractkit/chainlink-data-streams/rpc"

	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/grpc"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	promTransmitQueueInsertErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_queue_insert_error_count",
		Help:      "Running count of DB errors when trying to insert an item into the queue DB",
	},
		[]string{"donID", "serverURL"},
	)
	promTransmitQueuePushErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_queue_push_error_count",
		Help:      "Running count of DB errors when trying to push an item onto the queue",
	},
		[]string{"donID", "serverURL"},
	)
	promTransmitServerErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "transmit_server_error_count",
		Help:      "Number of errored transmissions that failed due to an error returned by the mercury server",
	},
		[]string{"donID", "serverURL", "code"},
	)
	promTransmitConcurrentTransmitGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "llo",
		Subsystem: "mercurytransmitter",
		Name:      "concurrent_transmit_gauge",
		Help:      "Gauge that measures the number of transmit threads currently waiting on a remote transmit call. You may wish to alert if this exceeds some number for a given period of time, or if it ever reaches its max.",
	},
		[]string{"donID", "serverURL"},
	)
)

type ReportPacker interface {
	Pack(digest types.ConfigDigest, seqNr uint64, report ocr2types.Report, sigs []ocr2types.AttributedOnchainSignature) ([]byte, error)
}

// A server handles the queue for a given mercury server

type server struct {
	lggr           logger.SugaredLogger
	verboseLogging bool

	transmitTimeout time.Duration

	c  grpc.Client
	pm *persistenceManager
	q  TransmitQueue

	url string

	evmPremiumLegacyPacker ReportPacker
	evmStreamlinedPacker   ReportPacker
	jsonPacker             ReportPacker

	transmitSuccessCount            prometheus.Counter
	transmitDuplicateCount          prometheus.Counter
	transmitConnectionErrorCount    prometheus.Counter
	transmitQueueInsertErrorCount   prometheus.Counter
	transmitQueuePushErrorCount     prometheus.Counter
	transmitConcurrentTransmitGauge prometheus.Gauge

	transmitThreadBusyCount         atomic.Int32
	consecutiveTransmitErrorCount   int
	consecutiveTransmitUniqueErrors map[string]struct{}
	consecutiveTransmitErrorMu      sync.Mutex
}

type QueueConfig interface {
	ReaperMaxAge() commonconfig.Duration
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

func newServer(lggr logger.Logger, verboseLogging bool, cfg QueueConfig, client grpc.Client, orm ORM, serverURL string) *server {
	pm := NewPersistenceManager(lggr, orm, serverURL, int(cfg.TransmitQueueMaxSize()), FlushDeletesFrequency, PruneFrequency, cfg.ReaperMaxAge().Duration())
	donIDStr := strconv.FormatUint(uint64(pm.DonID()), 10)
	var codecLggr logger.Logger
	if verboseLogging {
		codecLggr = lggr
	} else {
		codecLggr = corelogger.NullLogger
	}

	s := &server{
		logger.Sugared(lggr),
		verboseLogging,
		cfg.TransmitTimeout().Duration(),
		client,
		pm,
		NewTransmitQueue(lggr, serverURL, int(cfg.TransmitQueueMaxSize()), pm),
		serverURL,
		evm.NewReportCodecPremiumLegacy(codecLggr, pm.DonID()),
		evm.NewReportCodecStreamlined(),
		llo.JSONReportCodec{},
		promTransmitSuccessCount.WithLabelValues(donIDStr, serverURL),
		promTransmitDuplicateCount.WithLabelValues(donIDStr, serverURL),
		promTransmitConnectionErrorCount.WithLabelValues(donIDStr, serverURL),
		promTransmitQueueInsertErrorCount.WithLabelValues(donIDStr, serverURL),
		promTransmitQueuePushErrorCount.WithLabelValues(donIDStr, serverURL),
		promTransmitConcurrentTransmitGauge.WithLabelValues(donIDStr, serverURL),
		atomic.Int32{},
		0,
		make(map[string]struct{}),
		sync.Mutex{},
	}

	return s
}

func (s *server) HealthReport() map[string]error {
	report := map[string]error{}
	services.CopyHealth(report, s.c.HealthReport())
	services.CopyHealth(report, s.q.HealthReport())
	return report
}

func (s *server) transmitThreadBusyCountInc() {
	val := s.transmitThreadBusyCount.Add(1)
	s.transmitConcurrentTransmitGauge.Set(float64(val))
}
func (s *server) transmitThreadBusyCountDec() {
	val := s.transmitThreadBusyCount.Add(-1)
	s.transmitConcurrentTransmitGauge.Set(float64(val))
}

func (s *server) spawnTransmitLoops(stopCh services.StopChan, wg *sync.WaitGroup, donID uint32, n int) {
	donIDStr := strconv.FormatUint(uint64(donID), 10)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go s.spawnTransmitLoop(stopCh, wg, donIDStr)
	}
}

func (s *server) spawnTransmitLoop(stopCh services.StopChan, wg *sync.WaitGroup, donIDStr string) {
	defer wg.Done()
	s.transmitConcurrentTransmitGauge.Set(0) // initial set to populate metric

	// Exponential backoff with very short retry interval (since latency is a priority)
	// 5ms, 10ms, 20ms, 40ms etc
	b := backoff.Backoff{
		Min:    5 * time.Millisecond,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: true,
	}
	ctx, cancel := stopCh.NewCtx()
	defer cancel()
	cont := true
	for cont {
		cont = func() bool {
			t := s.q.BlockingPop()
			if t == nil {
				// queue was closed
				return false
			}

			s.transmitThreadBusyCountInc()
			defer s.transmitThreadBusyCountDec()

			req, res, err := func(ctx context.Context) (*rpc.TransmitRequest, *rpc.TransmitResponse, error) {
				ctx, cancelFn := context.WithTimeout(ctx, utils.WithJitter(s.transmitTimeout))
				defer cancelFn()
				return s.transmit(ctx, t)
			}(ctx)

			lggr := s.lggr.With("transmission", t, "response", res, "transmissionHash", fmt.Sprintf("%x", t.Hash()))
			if req != nil {
				lggr = s.lggr.With("req.Payload", req.Payload, "req.ReportFormat", req.ReportFormat)
			}

			if ctx.Err() != nil {
				// only canceled on transmitter close so we can exit
				return false
			} else if err != nil {
				s.transmitConnectionErrorCount.Inc()
				s.rateLimitedLogError(lggr, "Transmit report failed", err.Error())
				if ok := s.q.Push(t); !ok {
					s.lggr.Error("Failed to push report to transmit queue; queue is closed")
					return false
				}
				// Wait a backoff duration before pulling the most recent transmission
				// the heap
				select {
				case <-time.After(b.Duration()):
					return true
				case <-stopCh:
					return false
				}
			}

			b.Reset()
			if res.Error == "" {
				s.transmitSuccessCount.Inc()
				s.resetConsecutiveTransmitFailures()
				lggr.Debug("Transmit report success")
			} else {
				// We don't need to retry here because the mercury server
				// has confirmed it received the report. We only need to retry
				// on networking/unknown errors
				switch res.Code {
				case DuplicateReport:
					s.transmitSuccessCount.Inc()
					s.transmitDuplicateCount.Inc()
					s.resetConsecutiveTransmitFailures()
					lggr.Debug("Transmit report success; duplicate report")
				default:
					promTransmitServerErrorCount.WithLabelValues(donIDStr, s.url, strconv.FormatInt(int64(res.Code), 10)).Inc()
					s.rateLimitedLogError(lggr, "Transmit report failed; mercury server returned error", fmt.Sprintf("mercury server returned error: %q, statusCode: %d", res.Error, res.Code))
				}
			}

			s.pm.AsyncDelete(t.Hash())
			return true
		}()
	}
}

func (s *server) rateLimitedLogError(lggr logger.Logger, msg string, err string) {
	cnt, uniqueErrors := s.incConsecutiveTransmitErrorCount(err)
	switch {
	case cnt < 10:
		// Log first 10 errors individually
		lggr.Errorw(msg, "nErrs", 1, "err", err)
		return
	case cnt < 10_000:
		// Log errors up to 10k in batches of 100
		if cnt%100 == 0 {
			lggr.Errorw(msg+" (100 failures)", "nErrs", 100, "uniqueErrors", uniqueErrors)
		}
		return
	default:
		// After that, log every 10k errors
		if cnt%10_000 == 0 {
			lggr.Errorw(msg+" (10,000 failures)", "nErrs", 10_000, "uniqueErrors", uniqueErrors)
		}
		return
	}
}

func (s *server) incConsecutiveTransmitErrorCount(errStr string) (int, []string) {
	s.consecutiveTransmitErrorMu.Lock()
	defer s.consecutiveTransmitErrorMu.Unlock()
	s.consecutiveTransmitErrorCount++
	s.consecutiveTransmitUniqueErrors[errStr] = struct{}{}
	return s.consecutiveTransmitErrorCount, slices.Sorted(maps.Keys(s.consecutiveTransmitUniqueErrors))
}

func (s *server) resetConsecutiveTransmitFailures() {
	s.consecutiveTransmitErrorMu.Lock()
	s.consecutiveTransmitErrorCount = 0
	s.consecutiveTransmitUniqueErrors = make(map[string]struct{})
	s.consecutiveTransmitErrorMu.Unlock()
}

func (s *server) transmit(ctx context.Context, t *Transmission) (*rpc.TransmitRequest, *rpc.TransmitResponse, error) {
	var payload []byte
	var err error

	switch t.Report.Info.ReportFormat {
	case llotypes.ReportFormatJSON:
		payload, err = s.jsonPacker.Pack(t.ConfigDigest, t.SeqNr, t.Report.Report, t.Sigs)
	case llotypes.ReportFormatEVMPremiumLegacy, llotypes.ReportFormatEVMABIEncodeUnpacked:
		payload, err = s.evmPremiumLegacyPacker.Pack(t.ConfigDigest, t.SeqNr, t.Report.Report, t.Sigs)
	case llotypes.ReportFormatEVMStreamlined:
		payload, err = s.evmStreamlinedPacker.Pack(t.ConfigDigest, t.SeqNr, t.Report.Report, t.Sigs)
	default:
		return nil, nil, fmt.Errorf("Transmit failed; don't know how to Pack unsupported report format: %q", t.Report.Info.ReportFormat)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("Transmit: encode failed; %w", err)
	}

	req := &rpc.TransmitRequest{
		Payload:      payload,
		ReportFormat: uint32(t.Report.Info.ReportFormat),
	}

	resp, err := s.c.Transmit(ctx, req)
	return req, resp, err
}
