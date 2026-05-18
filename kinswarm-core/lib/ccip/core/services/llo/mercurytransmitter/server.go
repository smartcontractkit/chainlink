package mercurytransmitter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	transmitQueueDeleteErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_queue_delete_error_count",
		Help: "Running count of DB errors when trying to delete an item from the queue DB",
	},
		[]string{"donID", "serverURL"},
	)
	transmitQueueInsertErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_queue_insert_error_count",
		Help: "Running count of DB errors when trying to insert an item into the queue DB",
	},
		[]string{"donID", "serverURL"},
	)
	transmitQueuePushErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_queue_push_error_count",
		Help: "Running count of DB errors when trying to push an item onto the queue",
	},
		[]string{"donID", "serverURL"},
	)
	transmitServerErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_mercury_transmit_server_error_count",
		Help: "Number of errored transmissions that failed due to an error returned by the mercury server",
	},
		[]string{"donID", "serverURL", "code"},
	)
)

// A server handles the queue for a given mercury server

type server struct {
	lggr logger.SugaredLogger

	transmitTimeout time.Duration

	c  wsrpc.Client
	pm *persistenceManager
	q  TransmitQueue

	deleteQueue chan [32]byte

	url string

	transmitSuccessCount          prometheus.Counter
	transmitDuplicateCount        prometheus.Counter
	transmitConnectionErrorCount  prometheus.Counter
	transmitQueueDeleteErrorCount prometheus.Counter
	transmitQueueInsertErrorCount prometheus.Counter
	transmitQueuePushErrorCount   prometheus.Counter
}

type QueueConfig interface {
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

func newServer(lggr logger.Logger, cfg QueueConfig, client wsrpc.Client, orm ORM, serverURL string) *server {
	pm := NewPersistenceManager(lggr, orm, serverURL, int(cfg.TransmitQueueMaxSize()), flushDeletesFrequency, pruneFrequency)
	donIDStr := fmt.Sprintf("%d", pm.DonID())
	return &server{
		logger.Sugared(lggr),
		cfg.TransmitTimeout().Duration(),
		client,
		pm,
		NewTransmitQueue(lggr, serverURL, int(cfg.TransmitQueueMaxSize()), pm),
		make(chan [32]byte, int(cfg.TransmitQueueMaxSize())),
		serverURL,
		transmitSuccessCount.WithLabelValues(donIDStr, serverURL),
		transmitDuplicateCount.WithLabelValues(donIDStr, serverURL),
		transmitConnectionErrorCount.WithLabelValues(donIDStr, serverURL),
		transmitQueueDeleteErrorCount.WithLabelValues(donIDStr, serverURL),
		transmitQueueInsertErrorCount.WithLabelValues(donIDStr, serverURL),
		transmitQueuePushErrorCount.WithLabelValues(donIDStr, serverURL),
	}
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
		case hash := <-s.deleteQueue:
			for {
				if err := s.pm.orm.Delete(runloopCtx, [][32]byte{hash}); err != nil {
					s.lggr.Errorw("Failed to delete transmission record", "err", err, "transmissionHash", hash)
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

func (s *server) runQueueLoop(stopCh services.StopChan, wg *sync.WaitGroup, donIDStr string) {
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
		res, err := s.transmit(ctx, t)
		cancel()
		if runloopCtx.Err() != nil {
			// runloop context is only canceled on transmitter close so we can
			// exit the runloop here
			return
		} else if err != nil {
			s.transmitConnectionErrorCount.Inc()
			s.lggr.Errorw("Transmit report failed", "err", err, "transmission", t)
			if ok := s.q.Push(t); !ok {
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
			s.lggr.Debugw("Transmit report success", "transmission", t, "response", res)
		} else {
			// We don't need to retry here because the mercury server
			// has confirmed it received the report. We only need to retry
			// on networking/unknown errors
			switch res.Code {
			case DuplicateReport:
				s.transmitSuccessCount.Inc()
				s.transmitDuplicateCount.Inc()
				s.lggr.Debugw("Transmit report success; duplicate report", "transmission", t, "response", res)
			default:
				transmitServerErrorCount.WithLabelValues(donIDStr, s.url, fmt.Sprintf("%d", res.Code)).Inc()
				s.lggr.Errorw("Transmit report failed; mercury server returned error", "response", res, "transmission", t, "err", res.Error, "code", res.Code)
			}
		}

		select {
		case s.deleteQueue <- t.Hash():
		default:
			s.lggr.Criticalw("Delete queue is full", "transmission", t)
		}
	}
}

func (s *server) transmit(ctx context.Context, t *Transmission) (*pb.TransmitResponse, error) {
	var payload []byte
	var err error

	switch t.Report.Info.ReportFormat {
	case llotypes.ReportFormatJSON:
		// TODO: exactly how to handle JSON here?
		// https://smartcontract-it.atlassian.net/browse/MERC-3659
		fallthrough
	case llotypes.ReportFormatEVMPremiumLegacy:
		payload, err = evm.ReportCodecPremiumLegacy{}.Pack(t.ConfigDigest, t.SeqNr, t.Report.Report, t.Sigs)
	default:
		return nil, fmt.Errorf("Transmit failed; unsupported report format: %q", t.Report.Info.ReportFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("Transmit: encode failed; %w", err)
	}

	req := &pb.TransmitRequest{
		Payload:      payload,
		ReportFormat: uint32(t.Report.Info.ReportFormat),
	}

	return s.c.Transmit(ctx, req)
}
