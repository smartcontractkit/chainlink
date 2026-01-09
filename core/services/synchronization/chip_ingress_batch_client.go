package synchronization

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
)

type chipIngressBatchClient struct {
	services.Service
	eng *services.Engine

	chipClient chipingress.Client

	logging bool

	telemBufferSize   uint
	telemMaxBatchSize uint
	telemSendInterval time.Duration
	telemSendTimeout  time.Duration

	workers      map[string]*chipIngressBatchWorker
	workersMutex sync.RWMutex
}

// NewChipIngressBatchClient returns a client backed by chipingress.Client that
// can send telemetry to the chip ingress server
func NewChipIngressBatchClient(chipClient chipingress.Client, logging bool, lggr logger.Logger, telemBufferSize uint, telemMaxBatchSize uint, telemSendInterval time.Duration, telemSendTimeout time.Duration) ChipIngressService {
	c := &chipIngressBatchClient{
		telemBufferSize:   telemBufferSize,
		telemMaxBatchSize: telemMaxBatchSize,
		telemSendInterval: telemSendInterval,
		telemSendTimeout:  telemSendTimeout,
		chipClient:        chipClient,
		logging:           logging,
		workers:           make(map[string]*chipIngressBatchWorker),
	}
	c.Service, c.eng = services.Config{
		Name:  "ChipIngressBatchClient",
		Start: c.start,
	}.NewServiceEngine(lggr)

	return c
}

// start initializes the chip ingress batch client and starts health monitoring
func (cc *chipIngressBatchClient) start(ctx context.Context) error {
	cc.startHealthMonitoring(ctx, cc.chipClient)
	return nil
}

// Send directs incoming telmetry messages to the worker responsible for pushing it to
// the ingress server. If the worker telemetry buffer is full, messages are dropped
// and a warning is logged.
func (cc *chipIngressBatchClient) Send(ctx context.Context, payload TelemPayload) {
	worker := cc.findOrCreateWorker(payload)

	select {
	case worker.chTelemetry <- payload:
		worker.dropMessageCount.Store(0)
	case <-ctx.Done():
		return
	default:
		worker.logBufferFullWithExpBackoff(payload)
	}
}

// findOrCreateWorker finds a worker by ContractID or creates a new one if none exists.
// The number of workers is naturally bounded by the number of unique (ContractID, TelemetryType)
// pairs, which is determined by the jobs configured on the node. Each job has a fixed contract
// and telemetry type, so the worker count is proportional to the number of active jobs.
func (cc *chipIngressBatchClient) findOrCreateWorker(payload TelemPayload) *chipIngressBatchWorker {
	cc.workersMutex.Lock()
	defer cc.workersMutex.Unlock()

	workerKey := fmt.Sprintf("%s_%s", payload.ContractID, payload.TelemType)
	worker, found := cc.workers[workerKey]

	if !found {
		worker = NewChipIngressBatchWorker(
			cc.telemMaxBatchSize,
			cc.telemSendTimeout,
			cc.chipClient,
			make(chan TelemPayload, cc.telemBufferSize),
			payload.ContractID,
			payload.TelemType,
			cc.eng,
			cc.logging,
		)
		cc.eng.GoTick(timeutil.NewTicker(func() time.Duration {
			return cc.telemSendInterval
		}), worker.Send)
		cc.workers[workerKey] = worker

		TelemetryClientWorkers.WithLabelValues(chipIngress, string(payload.TelemType)).Inc()
	}

	return worker
}

// startHealthMonitoring starts a goroutine to monitor the connection state and update other relevant metrics every 5 seconds
func (cc *chipIngressBatchClient) startHealthMonitoring(_ context.Context, chipClient chipingress.Client) {
	cc.eng.GoTick(timeutil.NewTicker(func() time.Duration {
		return 5 * time.Second
	}), func(ctx context.Context) {
		connected := float64(0)
		pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
		_, err := chipClient.Ping(pingCtx, &pb.EmptyRequest{})
		pingCancel()
		if err == nil {
			connected = float64(1)
		} else {
			cc.eng.EmitHealthErr(err)
		}
		TelemetryClientConnectionStatus.WithLabelValues(chipIngress).Set(connected)
	})
}
