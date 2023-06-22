package functions

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/libocr/commontypes"
)

var (
	_ log.Listener   = &FunctionsListener{}
	_ job.ServiceCtx = &FunctionsListener{}

	sizeBuckets = []float64{
		1024,
		1024 * 4,
		1024 * 8,
		1024 * 16,
		1024 * 64,
		1024 * 256,
	}

	promOracleEvent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_oracle_event",
		Help: "Metric to track received oracle events",
	}, []string{"oracle", "event"})

	promRequestInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_internal_error",
		Help: "Metric to track internal errors",
	}, []string{"oracle"})

	promRequestComputationError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_computation_error",
		Help: "Metric to track computation errors",
	}, []string{"oracle"})

	promRequestComputationSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_computation_success",
		Help: "Metric to track number of computed requests",
	}, []string{"oracle"})

	promRequestTimeout = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_timeout",
		Help: "Metric to track number of timed out requests",
	}, []string{"oracle"})

	promRequestConfirmed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_confirmed",
		Help: "Metric to track number of confirmed requests",
	}, []string{"oracle", "responseType"})

	promRequestDataSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "functions_request_data_size",
		Help:    "Metric to track request data size",
		Buckets: sizeBuckets,
	}, []string{"oracle"})

	promComputationResultSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "functions_request_computation_result_size",
		Help: "Metric to track computation result size in bytes",
	}, []string{"oracle"})

	promComputationErrorSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "functions_request_computation_error_size",
		Help: "Metric to track computation error size in bytes",
	}, []string{"oracle"})

	promComputationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "functions_request_computation_duration",
		Help: "Metric to track computation duration in ms",
		Buckets: []float64{
			float64(10 * time.Millisecond),
			float64(100 * time.Millisecond),
			float64(500 * time.Millisecond),
			float64(time.Second),
			float64(10 * time.Second),
			float64(30 * time.Second),
			float64(60 * time.Second),
		},
	}, []string{"oracle"})

	promPrunedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_pruned",
		Help: "Metric to track number of requests pruned from the DB",
	}, []string{"oracle"})
)

const (
	DefaultPruneMaxStoredRequests uint32 = 20_000
	DefaultPruneCheckFrequencySec uint32 = 60 * 10
	DefaultPruneBatchSize         uint32 = 500
)

type FunctionsListener struct {
	utils.StartStopOnce
	oracle            *ocr2dr_oracle.OCR2DROracle
	oracleHexAddr     string
	job               job.Job
	bridgeAccessor    BridgeAccessor
	logBroadcaster    log.Broadcaster
	shutdownWaitGroup sync.WaitGroup
	mbOracleEvents    *utils.Mailbox[log.Broadcast]
	serviceContext    context.Context
	serviceCancel     context.CancelFunc
	chStop            chan struct{}
	pluginORM         ORM
	pluginConfig      config.PluginConfig
	logger            logger.Logger
	mailMon           *utils.MailboxMonitor
	urlsMonEndpoint   commontypes.MonitoringEndpoint
}

func formatRequestId(requestId [32]byte) string {
	return fmt.Sprintf("0x%x", requestId)
}

func NewFunctionsListener(oracle *ocr2dr_oracle.OCR2DROracle, job job.Job, bridgeAccessor BridgeAccessor, pluginORM ORM, pluginConfig config.PluginConfig, logBroadcaster log.Broadcaster, lggr logger.Logger, mailMon *utils.MailboxMonitor, urlsMonEndpoint commontypes.MonitoringEndpoint) *FunctionsListener {
	return &FunctionsListener{
		oracle:          oracle,
		oracleHexAddr:   oracle.Address().Hex(),
		job:             job,
		bridgeAccessor:  bridgeAccessor,
		logBroadcaster:  logBroadcaster,
		mbOracleEvents:  utils.NewHighCapacityMailbox[log.Broadcast](),
		chStop:          make(chan struct{}),
		pluginORM:       pluginORM,
		pluginConfig:    pluginConfig,
		logger:          lggr,
		mailMon:         mailMon,
		urlsMonEndpoint: urlsMonEndpoint,
	}
}

// Start complies with job.Service
func (l *FunctionsListener) Start(context.Context) error {
	return l.StartOnce("FunctionsListener", func() error {
		l.serviceContext, l.serviceCancel = context.WithCancel(context.Background())
		unsubscribeLogs := l.logBroadcaster.Register(l, log.ListenerOpts{
			Contract: l.oracle.Address(),
			ParseLog: l.oracle.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				ocr2dr_oracle.OCR2DROracleOracleRequest{}.Topic():        {},
				ocr2dr_oracle.OCR2DROracleOracleResponse{}.Topic():       {},
				ocr2dr_oracle.OCR2DROracleUserCallbackError{}.Topic():    {},
				ocr2dr_oracle.OCR2DROracleUserCallbackRawError{}.Topic(): {},
				ocr2dr_oracle.OCR2DROracleResponseTransmitted{}.Topic():  {},
			},
			MinIncomingConfirmations: l.pluginConfig.MinIncomingConfirmations,
		})
		if l.pluginConfig.ListenerEventHandlerTimeoutSec == 0 {
			l.logger.Warn("listenerEventHandlerTimeoutSec set to zero! ORM calls will never time out.")
		}
		l.shutdownWaitGroup.Add(4)
		go l.processOracleEvents()
		go l.timeoutRequests()
		go l.pruneRequests()
		go func() {
			<-l.chStop
			unsubscribeLogs()
			l.shutdownWaitGroup.Done()
		}()

		l.mailMon.Monitor(l.mbOracleEvents, "FunctionsListener", "OracleEvents", fmt.Sprint(l.job.ID))

		return nil
	})
}

// Close complies with job.Service
func (l *FunctionsListener) Close() error {
	return l.StopOnce("FunctionsListener", func() error {
		l.serviceCancel()
		close(l.chStop)
		l.shutdownWaitGroup.Wait()

		return l.mbOracleEvents.Close()
	})
}

// HandleLog implements log.Listener
func (l *FunctionsListener) HandleLog(lb log.Broadcast) {
	log := lb.DecodedLog()
	if log == nil || reflect.ValueOf(log).IsNil() {
		l.logger.Error("HandleLog: ignoring nil value")
		return
	}

	switch log := log.(type) {
	case *ocr2dr_oracle.OCR2DROracleOracleRequest, *ocr2dr_oracle.OCR2DROracleOracleResponse, *ocr2dr_oracle.OCR2DROracleUserCallbackError, *ocr2dr_oracle.OCR2DROracleUserCallbackRawError, *ocr2dr_oracle.OCR2DROracleResponseTransmitted:
		wasOverCapacity := l.mbOracleEvents.Deliver(lb)
		if wasOverCapacity {
			l.logger.Error("OracleRequest log mailbox is over capacity - dropped the oldest log")
		}
	default:
		l.logger.Errorf("Unexpected log type %T", log)
	}
}

// JobID() complies with log.Listener
func (l *FunctionsListener) JobID() int32 {
	return l.job.ID
}

func (l *FunctionsListener) processOracleEvents() {
	defer l.shutdownWaitGroup.Done()
	for {
		select {
		case <-l.chStop:
			return
		case <-l.mbOracleEvents.Notify():
			for {
				select {
				case <-l.chStop:
					return
				default:
				}
				lb, exists := l.mbOracleEvents.Retrieve()
				if !exists {
					break
				}
				was, err := l.logBroadcaster.WasAlreadyConsumed(lb)
				if err != nil {
					l.logger.Errorw("Could not determine if log was already consumed", "error", err)
					continue
				} else if was {
					continue
				}

				log := lb.DecodedLog()
				if log == nil || reflect.ValueOf(log).IsNil() {
					l.logger.Error("processOracleEvents: ignoring nil value")
					continue
				}

				switch log := log.(type) {
				case *ocr2dr_oracle.OCR2DROracleOracleRequest:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "OracleRequest").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleRequest(log, lb)
				case *ocr2dr_oracle.OCR2DROracleOracleResponse:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "OracleResponse").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponse("OracleResponse", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleUserCallbackError:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "UserCallbackError").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponse("UserCallbackError", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleUserCallbackRawError:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "UserCallbackRawError").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponse("UserCallbackRawError", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleResponseTransmitted:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "ResponseTransmitted").Inc()
				default:
					l.logger.Warnf("Unexpected log type %T", log)
				}
			}
		}
	}
}

func (l *FunctionsListener) getNewHandlerContext() (context.Context, context.CancelFunc) {
	timeoutSec := l.pluginConfig.ListenerEventHandlerTimeoutSec
	if timeoutSec == 0 {
		return context.WithCancel(l.serviceContext)
	}
	return context.WithTimeout(l.serviceContext, time.Duration(timeoutSec)*time.Second)
}

func (l *FunctionsListener) setError(ctx context.Context, requestId RequestID, runId int64, errType ErrType, errBytes []byte) {
	if errType == INTERNAL_ERROR {
		promRequestInternalError.WithLabelValues(l.oracleHexAddr).Inc()
	} else {
		promRequestComputationError.WithLabelValues(l.oracleHexAddr).Inc()
	}
	readyForProcessing := errType != INTERNAL_ERROR
	if err := l.pluginORM.SetError(requestId, runId, errType, errBytes, time.Now(), readyForProcessing, pg.WithParentCtx(ctx)); err != nil {
		l.logger.Errorw("call to SetError failed", "requestID", formatRequestId(requestId), "err", err)
	}
}

func (l *FunctionsListener) handleOracleRequest(request *ocr2dr_oracle.OCR2DROracleOracleRequest, lb log.Broadcast) {
	defer l.shutdownWaitGroup.Done()
	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	l.logger.Infow("oracle request received", "requestID", formatRequestId(request.RequestId))

	err := l.pluginORM.CreateRequest(request.RequestId, time.Now(), &request.Raw.TxHash, pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, ErrDuplicateRequestID) {
			l.logger.Warnw("received a log with duplicate request ID", "requestID", formatRequestId(request.RequestId), "err", err)
			l.markLogConsumed(lb, pg.WithParentCtx(ctx))
		} else {
			l.logger.Errorw("failed to create a DB entry for new request", "requestID", formatRequestId(request.RequestId), "err", err)
		}
		return
	}
	l.markLogConsumed(lb, pg.WithParentCtx(ctx))

	promRequestDataSize.WithLabelValues(l.oracleHexAddr).Observe(float64(len(request.Data)))

	if l.pluginConfig.MaxRequestSizeBytes > 0 && uint32(len(request.Data)) > l.pluginConfig.MaxRequestSizeBytes {
		l.logger.Errorw("request too big", "requestID", formatRequestId(request.RequestId), "requestSize", len(request.Data), "maxRequestSize", l.pluginConfig.MaxRequestSizeBytes)
		l.setError(ctx, request.RequestId, 0, USER_ERROR, []byte(fmt.Sprintf("request too big (max %d bytes)", l.pluginConfig.MaxRequestSizeBytes)))
		return
	}

	var requestData RequestData
	cborParseErr := cbor.ParseDietCBORToStruct(request.Data, &requestData)
	if cborParseErr != nil {
		l.logger.Errorw("failed to parse CBOR", "requestID", formatRequestId(request.RequestId), "err", cborParseErr)
		l.setError(ctx, request.RequestId, 0, USER_ERROR, []byte("CBOR parsing error"))
		return
	}

	l.handleRequest(ctx, request.RequestId, request.SubscriptionId, request.SubscriptionOwner, &requestData)
}

func (l *FunctionsListener) handleRequest(ctx context.Context, requestID [32]byte, subscriptionId uint64, subscriptionOwner common.Address, requestData *RequestData) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		promComputationDuration.WithLabelValues(l.oracleHexAddr).Observe(float64(duration.Milliseconds()))
	}()
	requestIDStr := formatRequestId(requestID)
	l.logger.Infow("processing request", "requestID", requestIDStr)

	eaClient, err := l.bridgeAccessor.NewExternalAdapterClient()
	if err != nil {
		l.logger.Errorw("failed to create ExternalAdapterClient", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, 0, INTERNAL_ERROR, []byte(err.Error()))
		return
	}

	computationResult, computationError, domains, err := eaClient.RunComputation(ctx, requestIDStr, l.job.Name.ValueOrZero(), subscriptionOwner.Hex(), subscriptionId, "", requestData)

	if err != nil {
		l.logger.Errorw("internal adapter error", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, 0, INTERNAL_ERROR, []byte(err.Error()))
		return
	}

	if len(computationError) == 0 && len(computationResult) == 0 {
		l.logger.Errorw("both result and error are empty - saving result", "requestID", requestIDStr)
		computationResult = []byte{}
		computationError = []byte{}
	}

	if len(domains) > 0 {
		l.reportSourceCodeDomains(requestID, domains)
	}

	if len(computationError) != 0 {
		if len(computationResult) != 0 {
			l.logger.Warnw("both result and error are non-empty - using error", "requestID", requestIDStr)
		}
		l.logger.Debugw("saving computation error", "requestID", requestIDStr)
		l.setError(ctx, requestID, 0, USER_ERROR, computationError)
		promComputationErrorSize.WithLabelValues(l.oracleHexAddr).Set(float64(len(computationError)))
	} else {
		promRequestComputationSuccess.WithLabelValues(l.oracleHexAddr).Inc()
		promComputationResultSize.WithLabelValues(l.oracleHexAddr).Set(float64(len(computationResult)))
		l.logger.Debugw("saving computation result", "requestID", requestIDStr)
		if err2 := l.pluginORM.SetResult(requestID, 0, computationResult, time.Now(), pg.WithParentCtx(ctx)); err2 != nil {
			l.logger.Errorw("call to SetResult failed", "requestID", requestIDStr, "err", err2)
		}
	}
}

func (l *FunctionsListener) handleOracleResponse(responseType string, requestID [32]byte, lb log.Broadcast) {
	defer l.shutdownWaitGroup.Done()
	l.logger.Infow("oracle response received", "type", responseType, "requestID", formatRequestId(requestID))

	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	if err := l.pluginORM.SetConfirmed(requestID, pg.WithParentCtx(ctx)); err != nil {
		l.logger.Errorw("setting CONFIRMED state failed", "requestID", formatRequestId(requestID), "err", err)
	}
	promRequestConfirmed.WithLabelValues(l.oracleHexAddr, responseType).Inc()
	l.markLogConsumed(lb, pg.WithParentCtx(ctx))
}

func (l *FunctionsListener) markLogConsumed(lb log.Broadcast, qopts ...pg.QOpt) {
	if err := l.logBroadcaster.MarkConsumed(lb, qopts...); err != nil {
		l.logger.Errorw("unable to mark log consumed", "err", err, "log", lb.String())
	}
}

func (l *FunctionsListener) timeoutRequests() {
	defer l.shutdownWaitGroup.Done()
	timeoutSec, freqSec, batchSize := l.pluginConfig.RequestTimeoutSec, l.pluginConfig.RequestTimeoutCheckFrequencySec, l.pluginConfig.RequestTimeoutBatchLookupSize
	if timeoutSec == 0 || freqSec == 0 || batchSize == 0 {
		l.logger.Warn("request timeout checker not configured - disabling it")
		return
	}
	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.chStop:
			return
		case <-ticker.C:
			cutoff := time.Now().Add(-(time.Duration(timeoutSec) * time.Second))
			ctx, cancel := l.getNewHandlerContext()
			ids, err := l.pluginORM.TimeoutExpiredResults(cutoff, batchSize, pg.WithParentCtx(ctx))
			cancel()
			if err != nil {
				l.logger.Errorw("error when calling FindExpiredResults", "err", err)
				break
			}
			if len(ids) > 0 {
				promRequestTimeout.WithLabelValues(l.oracleHexAddr).Add(float64(len(ids)))
				var idStrs []string
				for _, id := range ids {
					idStrs = append(idStrs, formatRequestId(id))
				}
				l.logger.Debugw("timed out requests", "requestIDs", idStrs)
			} else {
				l.logger.Debug("no requests to time out")
			}
		}
	}
}

func (l *FunctionsListener) pruneRequests() {
	defer l.shutdownWaitGroup.Done()
	maxStoredRequests, freqSec, batchSize := l.pluginConfig.PruneMaxStoredRequests, l.pluginConfig.PruneCheckFrequencySec, l.pluginConfig.PruneBatchSize
	if maxStoredRequests == 0 {
		l.logger.Warnw("pruneMaxStoredRequests not configured - using default", "DefaultPruneMaxStoredRequests", DefaultPruneMaxStoredRequests)
		maxStoredRequests = DefaultPruneMaxStoredRequests
	}
	if freqSec == 0 {
		l.logger.Warnw("pruneCheckFrequencySec not configured - using default", "DefaultPruneCheckFrequencySec", DefaultPruneCheckFrequencySec)
		freqSec = DefaultPruneCheckFrequencySec
	}
	if batchSize == 0 {
		l.logger.Warnw("pruneBatchSize not configured - using default", "DefaultPruneBatchSize", DefaultPruneBatchSize)
		batchSize = DefaultPruneBatchSize
	}

	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.chStop:
			return
		case <-ticker.C:
			ctx, cancel := l.getNewHandlerContext()
			startTime := time.Now()
			nTotal, nPruned, err := l.pluginORM.PruneOldestRequests(maxStoredRequests, batchSize, pg.WithParentCtx(ctx))
			cancel()
			elapsedMillis := time.Since(startTime).Milliseconds()
			if err != nil {
				l.logger.Errorw("error when calling PruneOldestRequests", "err", err, "elapsedMillis", elapsedMillis)
				break
			}
			if nPruned > 0 {
				promPrunedRequests.WithLabelValues(l.oracleHexAddr).Add(float64(nPruned))
				l.logger.Debugw("pruned requests from the DB", "nTotal", nTotal, "nPruned", nPruned, "elapsedMillis", elapsedMillis)
			} else {
				l.logger.Debugw("no pruned requests at this time", "nTotal", nTotal, "elapsedMillis", elapsedMillis)
			}
		}
	}
}

func (l *FunctionsListener) reportSourceCodeDomains(requestId RequestID, domains []string) {
	r := &telem.FunctionsRequest{
		RequestId:   formatRequestId(requestId),
		NodeAddress: l.job.OCR2OracleSpec.TransmitterID.ValueOrZero(),
		Domains:     domains,
	}

	bytes, err := proto.Marshal(r)
	if err != nil {
		l.logger.Warnw("telem.FunctionsRequest marshal error", "err", err)
	} else {
		l.urlsMonEndpoint.SendLog(bytes)
	}
}
