package functions

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-integrations/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/threshold"
	evmrelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

var (
	sizeBuckets = []float64{
		1024,
		1024 * 4,
		1024 * 8,
		1024 * 16,
		1024 * 64,
		1024 * 256,
	}

	promRequestReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_received",
		Help: "Metric to track received request events",
	}, []string{"router"})

	promRequestInternalError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_internal_error",
		Help: "Metric to track internal errors",
	}, []string{"router"})

	promRequestComputationError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_computation_error",
		Help: "Metric to track computation errors",
	}, []string{"router"})

	promRequestComputationSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_computation_success",
		Help: "Metric to track number of computed requests",
	}, []string{"router"})

	promRequestTimeout = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_timeout",
		Help: "Metric to track number of timed out requests",
	}, []string{"router"})

	promRequestConfirmed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_confirmed",
		Help: "Metric to track number of confirmed requests",
	}, []string{"router"})

	promRequestDataSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "functions_request_data_size",
		Help:    "Metric to track request data size",
		Buckets: sizeBuckets,
	}, []string{"router"})

	promComputationResultSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "functions_request_computation_result_size",
		Help: "Metric to track computation result size in bytes",
	}, []string{"router"})

	promComputationErrorSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "functions_request_computation_error_size",
		Help: "Metric to track computation error size in bytes",
	}, []string{"router"})

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
	}, []string{"router"})

	promPrunedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "functions_request_pruned",
		Help: "Metric to track number of requests pruned from the DB",
	}, []string{"router"})
)

const (
	DefaultPruneMaxStoredRequests uint32 = 20_000
	DefaultPruneCheckFrequencySec uint32 = 60 * 10
	DefaultPruneBatchSize         uint32 = 500

	// Used in place of OnchainMetadata for all offchain requests.
	OffchainRequestMarker string = "OFFCHAIN_REQUEST"

	FlagCBORMaxSize    uint32 = 1
	FlagSecretsMaxSize uint32 = 2
)

type FunctionsListener interface {
	job.ServiceCtx

	HandleOffchainRequest(ctx context.Context, request *OffchainRequest) error
}

type functionsListener struct {
	services.StateMachine
	client             client.Client
	contractAddressHex string
	job                job.Job
	bridgeAccessor     BridgeAccessor
	shutdownWaitGroup  sync.WaitGroup
	chStop             services.StopChan
	pluginORM          ORM
	pluginConfig       config.PluginConfig
	s4Storage          s4.Storage
	logger             logger.Logger
	urlsMonEndpoint    commontypes.MonitoringEndpoint
	decryptor          threshold.Decryptor
	logPollerWrapper   evmrelayTypes.LogPollerWrapper
}

var _ FunctionsListener = &functionsListener{}

func (l *functionsListener) HealthReport() map[string]error {
	return map[string]error{l.Name(): l.Healthy()}
}

func (l *functionsListener) Name() string { return l.logger.Name() }

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
	lggr logger.Logger,
	urlsMonEndpoint commontypes.MonitoringEndpoint,
	decryptor threshold.Decryptor,
	logPollerWrapper evmrelayTypes.LogPollerWrapper,
) *functionsListener {
	return &functionsListener{
		client:             client,
		contractAddressHex: contractAddressHex,
		job:                job,
		bridgeAccessor:     bridgeAccessor,
		chStop:             make(chan struct{}),
		pluginORM:          pluginORM,
		pluginConfig:       pluginConfig,
		s4Storage:          s4Storage,
		logger:             lggr,
		urlsMonEndpoint:    urlsMonEndpoint,
		decryptor:          decryptor,
		logPollerWrapper:   logPollerWrapper,
	}
}

// Start complies with job.Service
func (l *functionsListener) Start(context.Context) error {
	return l.StartOnce("FunctionsListener", func() error {
		switch l.pluginConfig.ContractVersion {
		case 1:
			l.shutdownWaitGroup.Add(1)
			go l.processOracleEventsV1()
		default:
			return fmt.Errorf("unsupported contract version: %d", l.pluginConfig.ContractVersion)
		}

		if l.pluginConfig.ListenerEventHandlerTimeoutSec == 0 {
			l.logger.Warn("listenerEventHandlerTimeoutSec set to zero! ORM calls will never time out.")
		}
		l.shutdownWaitGroup.Add(3)
		go l.timeoutRequests()
		go l.pruneRequests()
		go func() {
			<-l.chStop
			l.shutdownWaitGroup.Done()
		}()
		return nil
	})
}

// Close complies with job.Service
func (l *functionsListener) Close() error {
	return l.StopOnce("FunctionsListener", func() error {
		close(l.chStop)
		l.shutdownWaitGroup.Wait()
		return nil
	})
}

func (l *functionsListener) processOracleEventsV1() {
	defer l.shutdownWaitGroup.Done()
	ctx, cancel := l.chStop.NewCtx()
	defer cancel()
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
			requests, responses, err := l.logPollerWrapper.LatestEvents(ctx)
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

func (l *functionsListener) getNewHandlerContext() (context.Context, context.CancelFunc) {
	ctx, cancel := l.chStop.NewCtx()
	timeoutSec := l.pluginConfig.ListenerEventHandlerTimeoutSec
	if timeoutSec == 0 {
		return ctx, cancel
	}
	var cancel2 func()
	ctx, cancel2 = context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	return ctx, func() {
		cancel2()
		cancel()
	}
}

func (l *functionsListener) setError(ctx context.Context, requestId RequestID, errType ErrType, errBytes []byte) {
	if errType == INTERNAL_ERROR {
		promRequestInternalError.WithLabelValues(l.contractAddressHex).Inc()
	} else {
		promRequestComputationError.WithLabelValues(l.contractAddressHex).Inc()
	}
	readyForProcessing := errType != INTERNAL_ERROR
	if err := l.pluginORM.SetError(ctx, requestId, errType, errBytes, time.Now(), readyForProcessing); err != nil {
		l.logger.Errorw("call to SetError failed", "requestID", formatRequestId(requestId), "err", err)
	}
}

func (l *functionsListener) getMaxCBORsize(flags RequestFlags) uint32 {
	idx := flags[FlagCBORMaxSize]
	if int(idx) >= len(l.pluginConfig.MaxRequestSizesList) {
		return l.pluginConfig.MaxRequestSizeBytes // deprecated
	}
	return l.pluginConfig.MaxRequestSizesList[idx]
}

func (l *functionsListener) getMaxSecretsSize(flags RequestFlags) uint32 {
	idx := flags[FlagSecretsMaxSize]
	if int(idx) >= len(l.pluginConfig.MaxSecretsSizesList) {
		return math.MaxUint32 // not enforced if not configured
	}
	return l.pluginConfig.MaxSecretsSizesList[idx]
}

func (l *functionsListener) HandleOffchainRequest(ctx context.Context, request *OffchainRequest) error {
	if request == nil {
		return errors.New("HandleOffchainRequest: received nil request")
	}
	if len(request.RequestId) != RequestIDLength {
		return fmt.Errorf("HandleOffchainRequest: invalid request ID length %d", len(request.RequestId))
	}
	if len(request.SubscriptionOwner) != common.AddressLength || len(request.RequestInitiator) != common.AddressLength {
		return fmt.Errorf("HandleOffchainRequest: SubscriptionOwner and RequestInitiator must be set to valid addresses")
	}
	if request.Timestamp < uint64(time.Now().Unix()-int64(l.pluginConfig.RequestTimeoutSec)) {
		return fmt.Errorf("HandleOffchainRequest: request timestamp is too old")
	}

	var requestId RequestID
	copy(requestId[:], request.RequestId[:32])
	subscriptionOwner := common.BytesToAddress(request.SubscriptionOwner)
	senderAddr := common.BytesToAddress(request.RequestInitiator)
	emptyTxHash := common.Hash{}
	zeroCallbackGasLimit := uint32(0)
	newReq := &Request{
		RequestID:        requestId,
		RequestTxHash:    &emptyTxHash,
		ReceivedAt:       time.Now(),
		Flags:            []byte{},
		CallbackGasLimit: &zeroCallbackGasLimit,
		// use sender address in place of coordinator contract to keep batches uniform
		CoordinatorContractAddress: &senderAddr,
		OnchainMetadata:            []byte(OffchainRequestMarker),
	}
	if err := l.pluginORM.CreateRequest(ctx, newReq); err != nil {
		if errors.Is(err, ErrDuplicateRequestID) {
			l.logger.Warnw("HandleOffchainRequest: received duplicate request ID", "requestID", formatRequestId(requestId), "err", err)
		} else {
			l.logger.Errorw("HandleOffchainRequest: failed to create a DB entry for new request", "requestID", formatRequestId(requestId), "err", err)
		}
		return err
	}
	return l.handleRequest(ctx, requestId, request.SubscriptionId, subscriptionOwner, RequestFlags{}, &request.Data)
}

func (l *functionsListener) handleOracleRequestV1(request *evmrelayTypes.OracleRequest) {
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
	if err := l.pluginORM.CreateRequest(ctx, newReq); err != nil {
		if errors.Is(err, ErrDuplicateRequestID) {
			l.logger.Warnw("handleOracleRequestV1: received a log with duplicate request ID", "requestID", formatRequestId(request.RequestId), "err", err)
		} else {
			l.logger.Errorw("handleOracleRequestV1: failed to create a DB entry for new request", "requestID", formatRequestId(request.RequestId), "err", err)
		}
		return
	}

	promRequestReceived.WithLabelValues(l.contractAddressHex).Inc()
	promRequestDataSize.WithLabelValues(l.contractAddressHex).Observe(float64(len(request.Data)))
	requestData, err := l.parseCBOR(request.RequestId, request.Data, l.getMaxCBORsize(request.Flags))
	if err != nil {
		l.setError(ctx, request.RequestId, USER_ERROR, []byte(err.Error()))
		return
	}
	err = l.handleRequest(ctx, request.RequestId, request.SubscriptionId, request.SubscriptionOwner, request.Flags, requestData)
	if err != nil {
		l.logger.Errorw("handleOracleRequestV1: error in handleRequest()", "requestID", formatRequestId(request.RequestId), "err", err)
	}
}

func (l *functionsListener) parseCBOR(requestId RequestID, cborData []byte, maxSizeBytes uint32) (*RequestData, error) {
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

// Handle secret fetching/decryption and functions computation. Return error only for internal errors.
func (l *functionsListener) handleRequest(ctx context.Context, requestID RequestID, subscriptionId uint64, subscriptionOwner common.Address, flags RequestFlags, requestData *RequestData) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		promComputationDuration.WithLabelValues(l.contractAddressHex).Observe(float64(duration.Milliseconds()))
	}()
	requestIDStr := formatRequestId(requestID)
	l.logger.Infow("processing request", "requestID", requestIDStr)

	eaClient, err := l.bridgeAccessor.NewExternalAdapterClient(ctx)
	if err != nil {
		l.logger.Errorw("failed to create ExternalAdapterClient", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(err.Error()))
		return err
	}

	nodeProvidedSecrets, userErr, internalErr := l.getSecrets(ctx, eaClient, requestID, subscriptionOwner, requestData)
	if internalErr != nil {
		l.logger.Errorw("internal error during getSecrets", "requestID", requestIDStr, "err", internalErr)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(internalErr.Error()))
		return internalErr
	}
	if userErr != nil {
		l.logger.Debugw("user error during getSecrets", "requestID", requestIDStr, "err", userErr)
		l.setError(ctx, requestID, USER_ERROR, []byte(userErr.Error()))
		return nil // user error
	}

	maxSecretsSize := l.getMaxSecretsSize(flags)
	if uint32(len(nodeProvidedSecrets)) > maxSecretsSize {
		l.logger.Errorw("secrets size too big", "requestID", requestIDStr, "secretsSize", len(nodeProvidedSecrets), "maxSecretsSize", maxSecretsSize)
		l.setError(ctx, requestID, USER_ERROR, []byte("secrets size too big"))
		return nil // user error
	}

	computationResult, computationError, domains, err := eaClient.RunComputation(ctx, requestIDStr, l.job.Name.ValueOrZero(), subscriptionOwner.Hex(), subscriptionId, flags, nodeProvidedSecrets, requestData)

	if err != nil {
		l.logger.Errorw("internal adapter error", "requestID", requestIDStr, "err", err)
		l.setError(ctx, requestID, INTERNAL_ERROR, []byte(err.Error()))
		return err
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
		if err2 := l.pluginORM.SetResult(ctx, requestID, computationResult, time.Now()); err2 != nil {
			l.logger.Errorw("call to SetResult failed", "requestID", requestIDStr, "err", err2)
			return err2
		}
	}
	return nil
}

func (l *functionsListener) handleOracleResponseV1(response *evmrelayTypes.OracleResponse) {
	defer l.shutdownWaitGroup.Done()
	l.logger.Infow("oracle response v1 received", "requestID", formatRequestId(response.RequestId))

	ctx, cancel := l.getNewHandlerContext()
	defer cancel()
	if err := l.pluginORM.SetConfirmed(ctx, response.RequestId); err != nil {
		l.logger.Errorw("setting CONFIRMED state failed", "requestID", formatRequestId(response.RequestId), "err", err)
	}
	promRequestConfirmed.WithLabelValues(l.contractAddressHex).Inc()
}

func (l *functionsListener) timeoutRequests() {
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
			ids, err := l.pluginORM.TimeoutExpiredResults(ctx, cutoff, batchSize)
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

func (l *functionsListener) pruneRequests() {
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
			nTotal, nPruned, err := l.pluginORM.PruneOldestRequests(ctx, maxStoredRequests, batchSize)
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

func (l *functionsListener) reportSourceCodeDomains(requestId RequestID, domains []string) {
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

func (l *functionsListener) getSecrets(ctx context.Context, eaClient ExternalAdapterClient, requestID RequestID, subscriptionOwner common.Address, requestData *RequestData) (decryptedSecrets string, userError, internalError error) {
	if l.decryptor == nil {
		l.logger.Warn("Decryptor not configured")
		return "", nil, nil
	}

	var secrets []byte
	requestIDStr := formatRequestId(requestID)

	switch requestData.SecretsLocation {
	case LocationInline:
		if len(requestData.Secrets) > 0 {
			l.logger.Warnw("request used Inline secrets location, processing with no secrets", "requestID", requestIDStr)
		} else {
			l.logger.Debugw("request does not use any secrets", "requestID", requestIDStr)
		}
		return "", nil, nil
	case LocationRemote:
		thresholdEncSecrets, userError, err := eaClient.FetchEncryptedSecrets(ctx, requestData.Secrets, requestIDStr, l.job.Name.ValueOrZero())
		if err != nil {
			return "", nil, errors.Wrap(err, "failed to fetch encrypted secrets")
		}
		if len(userError) != 0 {
			return "", errors.New(string(userError)), nil
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
			return "", errors.Wrap(err, "failed to fetch DONHosted secrets"), nil
		}
		secrets = record.Payload
	}

	if len(secrets) == 0 {
		return "", nil, nil
	}

	decryptCtx, cancel := context.WithTimeout(ctx, time.Duration(l.pluginConfig.DecryptionQueueConfig.DecryptRequestTimeoutSec)*time.Second)
	defer cancel()

	decryptedSecretsBytes, err := l.decryptor.Decrypt(decryptCtx, requestID[:], secrets)
	if err != nil {
		l.logger.Debugw("threshold decryption of secrets failed", "requestID", requestIDStr, "err", err)
		return "", errors.New("threshold decryption of secrets failed"), nil
	}
	return string(decryptedSecretsBytes), nil, nil
}
