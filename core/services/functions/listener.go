package functions

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/ocr2dr_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	evmrelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
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

	FlagCBORMaxSize    uint32 = 1
	FlagSecretsMaxSize uint32 = 2
)

type FunctionsListener struct {
	utils.StartStopOnce
	client             client.Client
	contractAddressHex string
	job                job.Job
	bridgeAccessor     BridgeAccessor
	logBroadcaster     log.Broadcaster
	shutdownWaitGroup  sync.WaitGroup
	mbOracleEvents     *utils.Mailbox[log.Broadcast]
	serviceContext     context.Context
	serviceCancel      context.CancelFunc
	chStop             chan struct{}
	pluginORM          ORM
	pluginConfig       config.PluginConfig
	s4Storage          s4.Storage
	logger             logger.Logger
	mailMon            *utils.MailboxMonitor
	urlsMonEndpoint    commontypes.MonitoringEndpoint
	decryptor          threshold.Decryptor
	logPollerWrapper   evmrelayTypes.LogPollerWrapper
}

func formatRequestId(requestId [32]byte) string {
	return fmt.Sprintf("0x%x", requestId)
}

func NewFunctionsListener(
	job job.Job,
	client client.Client,
	contractAddressHex string,
	bridgeAccessor BridgeAccessor,
	pluginORM ORM,
	pluginConfig config.PluginConfig,
	s4Storage s4.Storage,
	logBroadcaster log.Broadcaster,
	lggr logger.Logger,
	mailMon *utils.MailboxMonitor,
	urlsMonEndpoint commontypes.MonitoringEndpoint,
	decryptor threshold.Decryptor,
	logPollerWrapper evmrelayTypes.LogPollerWrapper,
) *FunctionsListener {
	return &FunctionsListener{
		client:             client,
		contractAddressHex: contractAddressHex,
		job:                job,
		bridgeAccessor:     bridgeAccessor,
		logBroadcaster:     logBroadcaster,
		mbOracleEvents:     utils.NewHighCapacityMailbox[log.Broadcast](),
		chStop:             make(chan struct{}),
		pluginORM:          pluginORM,
		pluginConfig:       pluginConfig,
		s4Storage:          s4Storage,
		logger:             lggr,
		mailMon:            mailMon,
		urlsMonEndpoint:    urlsMonEndpoint,
		decryptor:          decryptor,
		logPollerWrapper:   logPollerWrapper,
	}
}

// Start complies with job.Service
func (l *FunctionsListener) Start(context.Context) error {
	return l.StartOnce("FunctionsListener", func() error {
		l.serviceContext, l.serviceCancel = context.WithCancel(context.Background())
		contractAddress := common.HexToAddress(l.contractAddressHex)
		var unsubscribeLogs func()

		switch l.pluginConfig.ContractVersion {
		case 0:
			oracleContract, err := ocr2dr_oracle.NewOCR2DROracle(contractAddress, l.client)
			if err != nil {
				return err
			}
			unsubscribeLogs = l.logBroadcaster.Register(l, log.ListenerOpts{
				Contract: oracleContract.Address(),
				ParseLog: oracleContract.ParseLog,
				LogsWithTopics: map[common.Hash][][]log.Topic{
					ocr2dr_oracle.OCR2DROracleOracleRequest{}.Topic():        {},
					ocr2dr_oracle.OCR2DROracleOracleResponse{}.Topic():       {},
					ocr2dr_oracle.OCR2DROracleUserCallbackError{}.Topic():    {},
					ocr2dr_oracle.OCR2DROracleUserCallbackRawError{}.Topic(): {},
					ocr2dr_oracle.OCR2DROracleResponseTransmitted{}.Topic():  {},
				},
				MinIncomingConfirmations: l.pluginConfig.MinIncomingConfirmations,
			})
			l.shutdownWaitGroup.Add(1)
			go l.processOracleEventsV0()
		case 1:
			l.shutdownWaitGroup.Add(1)
			go l.processOracleEventsV1()
		default:
			return errors.New("Functions: unsupported PluginConfig.ContractVersion")
		}

		if l.pluginConfig.ListenerEventHandlerTimeoutSec == 0 {
			l.logger.Warn("listenerEventHandlerTimeoutSec set to zero! ORM calls will never time out.")
		}
		l.shutdownWaitGroup.Add(3)
		go l.timeoutRequests()
		go l.pruneRequests()
		go func() {
			<-l.chStop
			if unsubscribeLogs != nil {
				unsubscribeLogs() // v0 only
			}
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
	case *ocr2dr_oracle.OCR2DROracleOracleRequest, *ocr2dr_oracle.OCR2DROracleOracleResponse, *ocr2dr_oracle.OCR2DROracleUserCallbackError, *ocr2dr_oracle.OCR2DROracleUserCallbackRawError, *ocr2dr_oracle.OCR2DROracleResponseTransmitted, *functions_coordinator.FunctionsCoordinatorOracleRequest, *functions_coordinator.FunctionsCoordinatorOracleResponse:
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

func (l *FunctionsListener) processOracleEventsV0() {
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
					l.logger.Errorw("Could not determine if log was already consumed", "err", err)
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
				// Version 0
				case *ocr2dr_oracle.OCR2DROracleOracleRequest:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "OracleRequest").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleRequestV0(log, lb)
				case *ocr2dr_oracle.OCR2DROracleOracleResponse:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "OracleResponse").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponseV0("OracleResponse", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleUserCallbackError:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "UserCallbackError").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponseV0("UserCallbackError", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleUserCallbackRawError:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "UserCallbackRawError").Inc()
					l.shutdownWaitGroup.Add(1)
					go l.handleOracleResponseV0("UserCallbackRawError", log.RequestId, lb)
				case *ocr2dr_oracle.OCR2DROracleResponseTransmitted:
					promOracleEvent.WithLabelValues(log.Raw.Address.Hex(), "ResponseTransmitted").Inc()
				default:
					l.logger.Warnf("Unexpected log type %T", log)
				}
			}
		}
	}
}

func (l *FunctionsListener) processOracleEventsV1() {
	defer l.shutdownWaitGroup.Done()
	freqMillis := l.pluginConfig.ListenerEventsCheckFrequencyMillis
	if freqMillis == 0 {
		l.logger.Errorw("ListenerEventsCheckFrequencyMillis must set to more than 0 in PluginConfig")
		return
	}
	ticker := time.NewTicker(time.Duration(freqMillis) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-l.chStop:
			return
		case <-ticker.C:
			requests, responses, err := l.logPollerWrapper.LatestEvents()
			if err != nil {
				l.logger.Errorw("error when calling LatestEvents()", "err", err)
				break
			}
			l.logger.Debugw("processOracleEventsV1: processing v1 events", "nRequests", len(requests), "nResponses", len(responses))
			for _, request := range requests {
				request := request
				l.shutdownWaitGroup.Add(1)
				go l.handleOracleRequestV1(&request)
			}
			for _, response := range responses {
				response := response
				l.shutdownWaitGroup.Add(1)
				go l.handleOracleResponseV1(&response)
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

func (l *FunctionsListener) setError(ctx context.Context, requestId RequestID, errType ErrType, errBytes []byte) {
	if errType == INTERNAL_ERROR {
		promRequestInternalError.WithLabelValues(l.contractAddressHex).Inc()
	} else {
		promRequestComputationError.WithLabelValues(l.contractAddressHex).Inc()
	}
	readyForProcessing := errType != INTERNAL_ERROR
	if err := l.pluginORM.SetError(requestId, errType, errBytes, time.Now(), readyForProcessing, pg.WithParentCtx(ctx)); err != nil {
		l.logger.Errorw("call to SetError failed", "requestID", formatRequestId(requestId), "err", err)
	}
}

func (l *FunctionsListener) getMaxCBORsize(flags RequestFlags) uint32 {
	idx := flags[FlagCBORMaxSize]
	if int(idx) >= len(l.pluginConfig.MaxRequestSizesList) {
		return l.pluginConfig.MaxRequestSizeBytes // deprecated
	}
	return l.pluginConfig.MaxRequestSizesList[idx]
}

func (l *FunctionsListener) getMaxSecretsSize(flags RequestFlags) uint32 {
	idx := flags[FlagSecretsMaxSize]
	if int(idx) >= len(l.pluginConfig.MaxSecretsSizesList) {
		return math.MaxUint32 // not enforced if not configured
	}
	return l.pluginConfig.MaxSecretsSizesList[idx]
}

func (l *FunctionsListener) handleOracleRequestV1(request *evmrelayTypes.OracleRequest) {
	defer l.shutdownWaitGroup.Done()
	l.logger.Infow("handleOracleRequestV1: oracle request v1 received", "requestID", formatRequestId(request.RequestId))
	ctx, cancel := l.getNewHandlerContext()
	defer cancel()

	callbackGasLimit := uint32(request.CallbackGasLimit)
	newReq := &Request{
		RequestID:                  request.RequestId,
		RequestTxHash:              &request.TxHash,
		ReceivedAt:                 time.Now(),
		Flags:                      request.Flags[:],
		CallbackGasLimit:           &callbackGasLimit,
		CoordinatorContractAddress: &request.CoordinatorContract,
		OnchainMetadata:            request.OnchainMetadata,
	}
	if err := l.pluginORM.CreateRequest(newReq, pg.WithParentCtx(ctx)); err != nil {
		if errors.Is(err, ErrDuplicateRequestID) {
			l.logger.Warnw("handleOracleRequestV1: received a log with duplicate request ID", "requestID", formatRequestId(request.RequestId), "err", err)
		} else {
			l.logger.Errorw("handleOracleRequestV1: failed to create a DB entry for new request", "requestID", formatRequestId(request.RequestId), "err", err)
		}
		return
	}

	promRequestDataSize.WithLabelValues(l.contractAddressHex).Observe(float64(len(request.Data)))
	requestData, err := l.parseCBOR(request.RequestId, request.Data, l.getMaxCBORsize(request.Flags))
	if err != nil {
		l.setError(ctx, request.RequestId, USER_ERROR, []byte(err.Error()))
		return
	}
	l.handleRequest(ctx, request.RequestId, request.SubscriptionId, request.SubscriptionOwner, request.Flags, requestData)
}

// deprecated
func (l *FunctionsListener) handleOracleRequestV0(request *ocr2dr_oracle.OCR2DROracleOracleRequest, lb log.Broadcast) {
	defer l.shutdownWaitGroup.Done()
	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	l.logger.Infow("oracle request received", "requestID", formatRequestId(request.RequestId))

	newReq := &Request{RequestID: request.RequestId, RequestTxHash: &request.Raw.TxHash, ReceivedAt: time.Now()}
	if err := l.pluginORM.CreateRequest(newReq, pg.WithParentCtx(ctx)); err != nil {
		if errors.Is(err, ErrDuplicateRequestID) {
			l.logger.Warnw("received a log with duplicate request ID", "requestID", formatRequestId(request.RequestId), "err", err)
			l.markLogConsumed(lb, pg.WithParentCtx(ctx))
		} else {
			l.logger.Errorw("failed to create a DB entry for new request", "requestID", formatRequestId(request.RequestId), "err", err)
		}
		return
	}
	l.markLogConsumed(lb, pg.WithParentCtx(ctx))

	promRequestDataSize.WithLabelValues(l.contractAddressHex).Observe(float64(len(request.Data)))
	requestData, err := l.parseCBOR(request.RequestId, request.Data, l.pluginConfig.MaxRequestSizeBytes)
	if err != nil {
		l.setError(ctx, request.RequestId, USER_ERROR, []byte(err.Error()))
		return
	}
	l.handleRequest(ctx, request.RequestId, request.SubscriptionId, request.SubscriptionOwner, [32]byte{}, requestData)
}

func (l *FunctionsListener) parseCBOR(requestId RequestID, cborData []byte, maxSizeBytes uint32) (*RequestData, error) {
	if maxSizeBytes > 0 && uint32(len(cborData)) > maxSizeBytes {
		l.logger.Errorw("request too big", "requestID", formatRequestId(requestId), "requestSize", len(cborData), "maxRequestSize", maxSizeBytes)
		return nil, fmt.Errorf("request too big (max %d bytes)", maxSizeBytes)
	}

	var requestData RequestData
	if err := cbor.ParseDietCBORToStruct(cborData, &requestData); err != nil {
		l.logger.Errorw("failed to parse CBOR", "requestID", formatRequestId(requestId), "err", err)
		return nil, errors.New("CBOR parsing error")
	}

	return &requestData, nil
}

func (l *FunctionsListener) handleRequest(ctx context.Context, requestID RequestID, subscriptionId uint64, subscriptionOwner common.Address, flags RequestFlags, requestData *RequestData) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		promComputationDuration.WithLabelValues(l.contractAddressHex).Observe(float64(duration.Milliseconds()))
	}()
	requestIDStr := formatRequestId(requestID)
	l.logger.Infow("processing request", "requestID", requestIDStr)

	if l.pluginConfig.ContractVersion == 1 && l.pluginConfig.EnableRequestSignatureCheck {
		err := VerifyRequestSignature(subscriptionOwner, requestData)
		if err != nil {
			l.logger.Errorw("invalid request signature", "requestID", requestIDStr, "err", err)
			l.setError(ctx, requestID, USER_ERROR, []byte(err.Error()))
			return
		}
	}

	eaClient, err := l.bridgeAccessor.NewExternalAdapterClient()
	if err != nil {
		l.logger.Errorw("failed to create ExternalAdapterClient", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(err.Error()))
		return
	}

	nodeProvidedSecrets, userErr, internalErr := l.getSecrets(ctx, eaClient, requestIDStr, subscriptionOwner, requestData)
	if internalErr != nil {
		l.logger.Errorw("internal error during getSecrets", "requestID", requestIDStr, "err", internalErr)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(internalErr.Error()))
		return
	}
	if userErr != nil {
		l.logger.Debugw("user error during getSecrets", "requestID", requestIDStr, "err", userErr)
		l.setError(ctx, requestID, USER_ERROR, []byte(userErr.Error()))
		return
	}

	maxSecretsSize := l.getMaxSecretsSize(flags)
	if uint32(len(nodeProvidedSecrets)) > maxSecretsSize {
		l.logger.Errorw("secrets size too big", "requestID", requestIDStr, "secretsSize", len(nodeProvidedSecrets), "maxSecretsSize", maxSecretsSize)
		l.setError(ctx, requestID, USER_ERROR, []byte("secrets size too big"))
		return
	}

	computationResult, computationError, domains, err := eaClient.RunComputation(ctx, requestIDStr, l.job.Name.ValueOrZero(), subscriptionOwner.Hex(), subscriptionId, flags, nodeProvidedSecrets, requestData)

	if err != nil {
		l.logger.Errorw("internal adapter error", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(err.Error()))
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
		l.setError(ctx, requestID, USER_ERROR, computationError)
		promComputationErrorSize.WithLabelValues(l.contractAddressHex).Set(float64(len(computationError)))
	} else {
		promRequestComputationSuccess.WithLabelValues(l.contractAddressHex).Inc()
		promComputationResultSize.WithLabelValues(l.contractAddressHex).Set(float64(len(computationResult)))
		l.logger.Debugw("saving computation result", "requestID", requestIDStr)
		if err2 := l.pluginORM.SetResult(requestID, computationResult, time.Now(), pg.WithParentCtx(ctx)); err2 != nil {
			l.logger.Errorw("call to SetResult failed", "requestID", requestIDStr, "err", err2)
		}
	}
}

func (l *FunctionsListener) handleOracleResponseV1(response *evmrelayTypes.OracleResponse) {
	defer l.shutdownWaitGroup.Done()
	l.logger.Infow("oracle response v1 received", "requestID", formatRequestId(response.RequestId))

	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	if err := l.pluginORM.SetConfirmed(response.RequestId, pg.WithParentCtx(ctx)); err != nil {
		l.logger.Errorw("setting CONFIRMED state failed", "requestID", formatRequestId(response.RequestId), "err", err)
	}
	promRequestConfirmed.WithLabelValues(l.contractAddressHex, "OracleResponse").Inc()
}

func (l *FunctionsListener) handleOracleResponseV0(responseType string, requestID [32]byte, lb log.Broadcast) {
	defer l.shutdownWaitGroup.Done()
	l.logger.Infow("oracle response received", "type", responseType, "requestID", formatRequestId(requestID))

	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	if err := l.pluginORM.SetConfirmed(requestID, pg.WithParentCtx(ctx)); err != nil {
		l.logger.Errorw("setting CONFIRMED state failed", "requestID", formatRequestId(requestID), "err", err)
	}
	promRequestConfirmed.WithLabelValues(l.contractAddressHex, responseType).Inc()
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
				promRequestTimeout.WithLabelValues(l.contractAddressHex).Add(float64(len(ids)))
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
				promPrunedRequests.WithLabelValues(l.contractAddressHex).Add(float64(nPruned))
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

func (l *FunctionsListener) getSecrets(ctx context.Context, eaClient ExternalAdapterClient, requestID string, subscriptionOwner common.Address, requestData *RequestData) (decryptedSecrets string, userError, internalError error) {
	if l.decryptor == nil {
		l.logger.Warn("Decryptor not configured")
		return "", nil, nil
	}

	var secrets []byte

	switch requestData.SecretsLocation {
	case LocationInline:
		l.logger.Warnw("request used Inline secrets location, processing with no secrets", "requestID", requestID)
		return "", nil, nil
	case LocationRemote:
		thresholdEncSecrets, userError, err := eaClient.FetchEncryptedSecrets(ctx, requestData.Secrets, requestID, l.job.Name.ValueOrZero())
		if err != nil {
			return "", nil, errors.Wrap(err, "failed to fetch encrypted secrets")
		}
		if len(userError) != 0 {
			l.logger.Debugw("no valid threshold encrypted secrets detected, falling back to legacy secrets", "requestID", requestID, "err", string(userError))
		}
		secrets = thresholdEncSecrets
	case LocationDONHosted:
		if l.s4Storage == nil {
			return "", nil, errors.New("S4 storage not configured")
		}
		var donSecrets DONHostedSecrets
		if err := cbor.ParseDietCBORToStruct(requestData.Secrets, &donSecrets); err != nil {
			return "", errors.Wrap(err, "failed to parse DONHosted secrets CBOR"), nil
		}
		record, _, err := l.s4Storage.Get(ctx, &s4.Key{
			Address: subscriptionOwner,
			SlotId:  donSecrets.SlotID,
			Version: donSecrets.Version,
		})
		if err != nil {
			return "", errors.Wrap(err, "failed to fetch S4 record for a secret"), nil
		}
		secrets = record.Payload
	}

	if len(secrets) == 0 {
		return "", nil, nil
	}

	decryptCtx, cancel := context.WithTimeout(ctx, time.Duration(l.pluginConfig.DecryptionQueueConfig.DecryptRequestTimeoutSec)*time.Second)
	defer cancel()

	decryptedSecretsBytes, err := l.decryptor.Decrypt(decryptCtx, []byte(requestID), secrets)
	if err != nil {
		return "", errors.New("threshold decryption of secrets failed"), nil
	}
	return string(decryptedSecretsBytes), nil, nil
}
