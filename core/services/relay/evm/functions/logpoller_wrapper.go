package functions

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type logPollerWrapper struct {
	services.Service
	eng *services.Engine

	routerContract            *functions_router.FunctionsRouter
	pluginConfig              config.PluginConfig
	client                    client.Client
	logPoller                 logpoller.LogPoller
	subscribers               map[string]evmRelayTypes.RouteUpdateSubscriber
	activeCoordinator         common.Address
	proposedCoordinator       common.Address
	requestBlockOffset        int64
	responseBlockOffset       int64
	pastBlocksToPoll          int64
	logPollerCacheDurationSec int64
	detectedRequests          detectedEvents
	detectedResponses         detectedEvents
	mu                        sync.Mutex
}

type detectedEvent struct {
	requestId    [32]byte
	timeDetected time.Time
}

type detectedEvents struct {
	isPreviouslyDetected  map[[32]byte]struct{}
	detectedEventsOrdered []detectedEvent
}

const logPollerCacheDurationSecDefault = 300
const pastBlocksToPollDefault = 50
const maxLogsToProcess = 1000

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(routerContractAddress common.Address, pluginConfig config.PluginConfig, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) (evmRelayTypes.LogPollerWrapper, error) {
	routerContract, err := functions_router.NewFunctionsRouter(routerContractAddress, client)
	if err != nil {
		return nil, err
	}
	blockOffset := int64(pluginConfig.MinIncomingConfirmations) - 1
	if blockOffset < 0 {
		lggr.Warnw("invalid minIncomingConfirmations, using 1 instead", "minIncomingConfirmations", pluginConfig.MinIncomingConfirmations)
		blockOffset = 0
	}
	requestBlockOffset := int64(pluginConfig.MinRequestConfirmations) - 1
	if requestBlockOffset < 0 {
		lggr.Warnw("invalid minRequestConfirmations, using minIncomingConfirmations instead", "minRequestConfirmations", pluginConfig.MinRequestConfirmations)
		requestBlockOffset = blockOffset
	}
	responseBlockOffset := int64(pluginConfig.MinResponseConfirmations) - 1
	if responseBlockOffset < 0 {
		lggr.Warnw("invalid minResponseConfirmations, using minIncomingConfirmations instead", "minResponseConfirmations", pluginConfig.MinResponseConfirmations)
		responseBlockOffset = blockOffset
	}
	logPollerCacheDurationSec := int64(pluginConfig.LogPollerCacheDurationSec)
	if logPollerCacheDurationSec <= 0 {
		lggr.Warnw("invalid logPollerCacheDuration, using 300 instead", "logPollerCacheDurationSec", logPollerCacheDurationSec)
		logPollerCacheDurationSec = logPollerCacheDurationSecDefault
	}
	pastBlocksToPoll := int64(pluginConfig.PastBlocksToPoll)
	if pastBlocksToPoll <= 0 {
		lggr.Warnw("invalid pastBlocksToPoll, using 50 instead", "pastBlocksToPoll", pastBlocksToPoll)
		pastBlocksToPoll = pastBlocksToPollDefault
	}
	if blockOffset >= pastBlocksToPoll || requestBlockOffset >= pastBlocksToPoll || responseBlockOffset >= pastBlocksToPoll {
		lggr.Errorw("invalid config: number of required confirmation blocks >= pastBlocksToPoll", "pastBlocksToPoll", pastBlocksToPoll, "minIncomingConfirmations", pluginConfig.MinIncomingConfirmations, "minRequestConfirmations", pluginConfig.MinRequestConfirmations, "minResponseConfirmations", pluginConfig.MinResponseConfirmations)
		return nil, errors.Errorf("invalid config: number of required confirmation blocks >= pastBlocksToPoll")
	}

	w := &logPollerWrapper{
		routerContract:            routerContract,
		pluginConfig:              pluginConfig,
		requestBlockOffset:        requestBlockOffset,
		responseBlockOffset:       responseBlockOffset,
		pastBlocksToPoll:          pastBlocksToPoll,
		logPollerCacheDurationSec: logPollerCacheDurationSec,
		detectedRequests:          detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})},
		detectedResponses:         detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})},
		logPoller:                 logPoller,
		client:                    client,
		subscribers:               make(map[string]evmRelayTypes.RouteUpdateSubscriber),
	}
	w.Service, w.eng = services.Config{
		Name:  "LoggPollerWrapper",
		Start: w.start,
	}.NewServiceEngine(lggr)
	return w, nil
}

func (l *logPollerWrapper) start(context.Context) error {
	l.eng.Infow("starting LogPollerWrapper", "routerContract", l.routerContract.Address().Hex(), "contractVersion", l.pluginConfig.ContractVersion)
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.pluginConfig.ContractVersion != 1 {
		return errors.New("only contract version 1 is supported")
	}
	l.eng.Go(l.checkForRouteUpdates)
	return nil
}

// methods of LogPollerWrapper
func (l *logPollerWrapper) LatestEvents(ctx context.Context) ([]evmRelayTypes.OracleRequest, []evmRelayTypes.OracleResponse, error) {
	l.mu.Lock()
	coordinators := []common.Address{}
	if l.activeCoordinator != (common.Address{}) {
		coordinators = append(coordinators, l.activeCoordinator)
	}
	if l.proposedCoordinator != (common.Address{}) && l.activeCoordinator != l.proposedCoordinator {
		coordinators = append(coordinators, l.proposedCoordinator)
	}
	latest, err := l.logPoller.LatestBlock(ctx)
	if err != nil {
		l.mu.Unlock()
		return nil, nil, err
	}
	latestBlockNum := latest.BlockNumber
	startBlockNum := latestBlockNum - l.pastBlocksToPoll
	if startBlockNum < 0 {
		startBlockNum = 0
	}
	l.mu.Unlock()

	// outside of the lock
	resultsReq := []evmRelayTypes.OracleRequest{}
	resultsResp := []evmRelayTypes.OracleResponse{}
	if len(coordinators) == 0 {
		l.eng.Debug("LatestEvents: no non-zero coordinators to check")
		return resultsReq, resultsResp, errors.New("no non-zero coordinators to check")
	}

	for _, coordinator := range coordinators {
		requestEndBlock := latestBlockNum - l.requestBlockOffset
		requestLogs, err := l.logPoller.Logs(ctx, startBlockNum, requestEndBlock, functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(), coordinator)
		if err != nil {
			l.eng.Errorw("LatestEvents: fetching request logs from LogPoller failed", "startBlock", startBlockNum, "endBlock", requestEndBlock)
			return nil, nil, err
		}
		l.eng.Debugw("LatestEvents: fetched request logs", "nRequestLogs", len(requestLogs), "latestBlock", latest, "startBlock", startBlockNum, "endBlock", requestEndBlock)
		requestLogs = l.filterPreviouslyDetectedEvents(requestLogs, &l.detectedRequests, "requests")
		responseEndBlock := latestBlockNum - l.responseBlockOffset
		responseLogs, err := l.logPoller.Logs(ctx, startBlockNum, responseEndBlock, functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(), coordinator)
		if err != nil {
			l.eng.Errorw("LatestEvents: fetching response logs from LogPoller failed", "startBlock", startBlockNum, "endBlock", responseEndBlock)
			return nil, nil, err
		}
		l.eng.Debugw("LatestEvents: fetched request logs", "nResponseLogs", len(responseLogs), "latestBlock", latest, "startBlock", startBlockNum, "endBlock", responseEndBlock)
		responseLogs = l.filterPreviouslyDetectedEvents(responseLogs, &l.detectedResponses, "responses")

		parsingContract, err := functions_coordinator.NewFunctionsCoordinator(coordinator, l.client)
		if err != nil {
			l.eng.Error("LatestEvents: creating a contract instance for parsing failed")
			return nil, nil, err
		}

		l.eng.Debugw("LatestEvents: parsing logs", "nRequestLogs", len(requestLogs), "nResponseLogs", len(responseLogs), "coordinatorAddress", coordinator.Hex())
		for _, log := range requestLogs {
			gethLog := log.ToGethLog()
			oracleRequest, err := parsingContract.ParseOracleRequest(gethLog)
			if err != nil {
				l.eng.Errorw("LatestEvents: failed to parse a request log, skipping", "err", err)
				continue
			}

			uint32Type, errType1 := abi.NewType("uint32", "uint32", nil)
			uint40Type, errType2 := abi.NewType("uint40", "uint40", nil)
			uint64Type, errType3 := abi.NewType("uint64", "uint64", nil)
			uint72Type, errType4 := abi.NewType("uint72", "uint72", nil)
			uint96Type, errType5 := abi.NewType("uint96", "uint96", nil)
			addressType, errType6 := abi.NewType("address", "address", nil)
			bytes32Type, errType7 := abi.NewType("bytes32", "bytes32", nil)

			if errType1 != nil || errType2 != nil || errType3 != nil || errType4 != nil || errType5 != nil || errType6 != nil || errType7 != nil {
				l.eng.Errorw("LatestEvents: failed to initialize types", "errType1", errType1,
					"errType2", errType2, "errType3", errType3, "errType4", errType4, "errType5", errType5, "errType6", errType6, "errType7", errType7,
				)
				continue
			}
			commitmentABI := abi.Arguments{
				{Type: bytes32Type}, // RequestId
				{Type: addressType}, // Coordinator
				{Type: uint96Type},  // EstimatedTotalCostJuels
				{Type: addressType}, // Client
				{Type: uint64Type},  // SubscriptionId
				{Type: uint32Type},  // CallbackGasLimit
				{Type: uint72Type},  // AdminFee
				{Type: uint72Type},  // DonFee
				{Type: uint40Type},  // GasOverheadBeforeCallback
				{Type: uint40Type},  // GasOverheadAfterCallback
				{Type: uint32Type},  // TimeoutTimestamp
			}
			commitmentBytes, err := commitmentABI.Pack(
				oracleRequest.Commitment.RequestId,
				oracleRequest.Commitment.Coordinator,
				oracleRequest.Commitment.EstimatedTotalCostJuels,
				oracleRequest.Commitment.Client,
				oracleRequest.Commitment.SubscriptionId,
				oracleRequest.Commitment.CallbackGasLimit,
				oracleRequest.Commitment.AdminFee,
				oracleRequest.Commitment.DonFee,
				oracleRequest.Commitment.GasOverheadBeforeCallback,
				oracleRequest.Commitment.GasOverheadAfterCallback,
				oracleRequest.Commitment.TimeoutTimestamp,
			)
			if err != nil {
				l.eng.Errorw("LatestEvents: failed to pack commitment bytes, skipping", "err", err)
			}

			resultsReq = append(resultsReq, evmRelayTypes.OracleRequest{
				RequestId:           oracleRequest.RequestId,
				RequestingContract:  oracleRequest.RequestingContract,
				RequestInitiator:    oracleRequest.RequestInitiator,
				SubscriptionId:      oracleRequest.SubscriptionId,
				SubscriptionOwner:   oracleRequest.SubscriptionOwner,
				Data:                oracleRequest.Data,
				DataVersion:         oracleRequest.DataVersion,
				Flags:               oracleRequest.Flags,
				CallbackGasLimit:    oracleRequest.CallbackGasLimit,
				TxHash:              oracleRequest.Raw.TxHash,
				OnchainMetadata:     commitmentBytes,
				CoordinatorContract: coordinator,
			})
		}
		for _, log := range responseLogs {
			gethLog := log.ToGethLog()
			oracleResponse, err := parsingContract.ParseOracleResponse(gethLog)
			if err != nil {
				l.eng.Errorw("LatestEvents: failed to parse a response log, skipping")
				continue
			}
			resultsResp = append(resultsResp, evmRelayTypes.OracleResponse{
				RequestId: oracleResponse.RequestId,
			})
		}
	}

	l.eng.Debugw("LatestEvents: done", "nRequestLogs", len(resultsReq), "nResponseLogs", len(resultsResp), "startBlock", startBlockNum, "endBlock", latestBlockNum)
	return resultsReq, resultsResp, nil
}

func (l *logPollerWrapper) filterPreviouslyDetectedEvents(logs []logpoller.Log, detectedEvents *detectedEvents, filterType string) []logpoller.Log {
	if len(logs) > maxLogsToProcess {
		l.eng.Errorw("filterPreviouslyDetectedEvents: too many logs to process, only processing latest maxLogsToProcess logs", "filterType", filterType, "nLogs", len(logs), "maxLogsToProcess", maxLogsToProcess)
		logs = logs[len(logs)-maxLogsToProcess:]
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	filteredLogs := []logpoller.Log{}
	for _, log := range logs {
		var requestId [32]byte
		if len(log.Topics) < 2 || len(log.Topics[1]) != 32 {
			l.eng.Errorw("filterPreviouslyDetectedEvents: invalid log, skipping", "filterType", filterType, "log", log)
			continue
		}
		copy(requestId[:], log.Topics[1]) // requestId is the second topic (1st topic is the event signature)
		if _, ok := detectedEvents.isPreviouslyDetected[requestId]; !ok {
			filteredLogs = append(filteredLogs, log)
			detectedEvents.isPreviouslyDetected[requestId] = struct{}{}
			detectedEvents.detectedEventsOrdered = append(detectedEvents.detectedEventsOrdered, detectedEvent{requestId: requestId, timeDetected: time.Now()})
		}
	}
	expiredRequests := 0
	for _, detectedEvent := range detectedEvents.detectedEventsOrdered {
		expirationTime := time.Now().Add(-time.Second * time.Duration(l.logPollerCacheDurationSec))
		if !detectedEvent.timeDetected.Before(expirationTime) {
			break
		}
		delete(detectedEvents.isPreviouslyDetected, detectedEvent.requestId)
		expiredRequests++
	}
	detectedEvents.detectedEventsOrdered = detectedEvents.detectedEventsOrdered[expiredRequests:]
	l.eng.Debugw("filterPreviouslyDetectedEvents: done", "filterType", filterType, "nLogs", len(logs), "nFilteredLogs", len(filteredLogs), "nExpiredRequests", expiredRequests, "previouslyDetectedCacheSize", len(detectedEvents.detectedEventsOrdered))
	return filteredLogs
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(ctx context.Context, subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.pluginConfig.ContractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		if err := subscriber.UpdateRoutes(ctx, l.routerContract.Address(), l.routerContract.Address()); err != nil {
			l.eng.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "err", err)
		}
	} else if l.pluginConfig.ContractVersion == 1 {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.subscribers[subscriberName] = subscriber
	}
}

func (l *logPollerWrapper) checkForRouteUpdates(ctx context.Context) {
	freqSec := l.pluginConfig.ContractUpdateCheckFrequencySec
	if freqSec == 0 {
		l.eng.Errorw("LogPollerWrapper: ContractUpdateCheckFrequencySec is zero - route update checks disabled")
		return
	}

	updateOnce := func(ctx context.Context) {
		// NOTE: timeout == frequency here, could be changed to a separate config value
		timeout := time.Duration(l.pluginConfig.ContractUpdateCheckFrequencySec) * time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		active, proposed, err := l.getCurrentCoordinators(ctx)
		if err != nil {
			l.eng.Errorw("LogPollerWrapper: error calling getCurrentCoordinators", "err", err)
			return
		}

		l.handleRouteUpdate(ctx, active, proposed)
	}

	updateOnce(ctx) // update once right away
	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateOnce(ctx)
		}
	}
}

func (l *logPollerWrapper) getCurrentCoordinators(ctx context.Context) (common.Address, common.Address, error) {
	if l.pluginConfig.ContractVersion == 0 {
		return l.routerContract.Address(), l.routerContract.Address(), nil
	}
	var donId [32]byte
	copy(donId[:], []byte(l.pluginConfig.DONID))

	activeCoordinator, err := l.routerContract.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	proposedCoordinator, err := l.routerContract.GetProposedContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId)
	if err != nil {
		return activeCoordinator, l.proposedCoordinator, nil
	}

	return activeCoordinator, proposedCoordinator, nil
}

func (l *logPollerWrapper) handleRouteUpdate(ctx context.Context, activeCoordinator common.Address, proposedCoordinator common.Address) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if activeCoordinator == (common.Address{}) {
		l.eng.Error("LogPollerWrapper: cannot update activeCoordinator to zero address")
		return
	}

	if activeCoordinator == l.activeCoordinator && proposedCoordinator == l.proposedCoordinator {
		l.eng.Debug("LogPollerWrapper: no changes to routes")
		return
	}
	errActive := l.registerFilters(ctx, activeCoordinator)
	errProposed := l.registerFilters(ctx, proposedCoordinator)
	if errActive != nil || errProposed != nil {
		l.eng.Errorw("LogPollerWrapper: Failed to register filters", "errorActive", errActive, "errorProposed", errProposed)
		return
	}

	l.eng.Debugw("LogPollerWrapper: new routes", "activeCoordinator", activeCoordinator.Hex(), "proposedCoordinator", proposedCoordinator.Hex())

	l.activeCoordinator = activeCoordinator
	l.proposedCoordinator = proposedCoordinator

	for _, subscriber := range l.subscribers {
		err := subscriber.UpdateRoutes(ctx, activeCoordinator, proposedCoordinator)
		if err != nil {
			l.eng.Errorw("LogPollerWrapper: Failed to update routes", "err", err)
		}
	}

	filters := l.logPoller.GetFilters()
	for _, filter := range filters {
		if filter.Name[:len(l.filterPrefix())] != l.filterPrefix() {
			continue
		}
		if filter.Name == l.filterName(l.activeCoordinator) || filter.Name == l.filterName(l.proposedCoordinator) {
			continue
		}
		if err := l.logPoller.UnregisterFilter(ctx, filter.Name); err != nil {
			l.eng.Errorw("LogPollerWrapper: Failed to unregister filter", "filterName", filter.Name, "err", err)
		}
		l.eng.Debugw("LogPollerWrapper: Successfully unregistered filter", "filterName", filter.Name)
	}
}

func (l *logPollerWrapper) filterPrefix() string {
	return "FunctionsLogPollerWrapper:" + l.pluginConfig.DONID
}

func (l *logPollerWrapper) filterName(addr common.Address) string {
	return logpoller.FilterName(l.filterPrefix(), addr.String())
}

func (l *logPollerWrapper) registerFilters(ctx context.Context, coordinatorAddress common.Address) error {
	if (coordinatorAddress == common.Address{}) {
		return nil
	}
	return l.logPoller.RegisterFilter(
		ctx,
		logpoller.Filter{
			Name: l.filterName(coordinatorAddress),
			EventSigs: []common.Hash{
				functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(),
				functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(),
			},
			Addresses: []common.Address{coordinatorAddress},
		})
}
