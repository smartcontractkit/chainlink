package functions

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	type_and_version "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/type_and_version_interface_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type logPollerWrapper struct {
	services.StateMachine

	routerContract            *functions_router.FunctionsRouter
	pluginConfig              config.PluginConfig
	client                    client.Client
	logPoller                 logpoller.LogPoller
	subscribers               map[string]evmRelayTypes.RouteUpdateSubscriber
	activeCoordinator         Coordinator
	proposedCoordinator       Coordinator
	requestBlockOffset        int64
	responseBlockOffset       int64
	pastBlocksToPoll          int64
	logPollerCacheDurationSec int64
	detectedRequests          detectedEvents
	detectedResponses         detectedEvents
	abiTypes                  *abiTypes
	mu                        sync.Mutex
	closeWait                 sync.WaitGroup
	stopCh                    services.StopChan
	lggr                      logger.Logger
}

type detectedEvent struct {
	requestId    [32]byte
	timeDetected time.Time
}

type detectedEvents struct {
	isPreviouslyDetected  map[[32]byte]struct{}
	detectedEventsOrdered []detectedEvent
}

type abiTypes struct {
	uint32Type  abi.Type
	uint40Type  abi.Type
	uint64Type  abi.Type
	uint72Type  abi.Type
	uint96Type  abi.Type
	addressType abi.Type
	bytes32Type abi.Type
}

type Coordinator interface {
	Address() common.Address
	RegisterFilters() error
	OracleRequestLogTopic() (common.Hash, error)
	OracleResponseLogTopic() (common.Hash, error)
	LogsToRequests(requestLogs []logpoller.Log) ([]evmRelayTypes.OracleRequest, error)
	LogsToResponses(responseLogs []logpoller.Log) ([]evmRelayTypes.OracleResponse, error)
}

const FUNCTIONS_COORDINATOR_VERSION_1_SUBSTRING = "Functions Coordinator v1"
const FUNCTIONS_COORDINATOR_VERSION_2_SUBSTRING = "Functions Coordinator v2"

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

	abiTypes, err := initAbiTypes()
	if err != nil {
		lggr.Errorf("failed to initialize abi types: %w", err)
		return nil, err
	}

	return &logPollerWrapper{
		routerContract:            routerContract,
		pluginConfig:              pluginConfig,
		requestBlockOffset:        requestBlockOffset,
		responseBlockOffset:       responseBlockOffset,
		pastBlocksToPoll:          pastBlocksToPoll,
		logPollerCacheDurationSec: logPollerCacheDurationSec,
		detectedRequests:          detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})},
		detectedResponses:         detectedEvents{isPreviouslyDetected: make(map[[32]byte]struct{})},
		abiTypes:                  abiTypes,
		logPoller:                 logPoller,
		client:                    client,
		subscribers:               make(map[string]evmRelayTypes.RouteUpdateSubscriber),
		stopCh:                    make(services.StopChan),
		lggr:                      lggr.Named("LogPollerWrapper"),
	}, nil
}

func (l *logPollerWrapper) Start(context.Context) error {
	return l.StartOnce("LogPollerWrapper", func() error {
		l.lggr.Infow("starting LogPollerWrapper", "routerContract", l.routerContract.Address().Hex(), "contractVersion", l.pluginConfig.ContractVersion)
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.pluginConfig.ContractVersion != 1 {
			return errors.New("only contract version 1 is supported")
		}
		l.closeWait.Add(1)
		go l.checkForRouteUpdates()
		return nil
	})
}

func (l *logPollerWrapper) Close() error {
	return l.StopOnce("LogPollerWrapper", func() (err error) {
		l.lggr.Info("closing LogPollerWrapper")
		close(l.stopCh)
		l.closeWait.Wait()
		return nil
	})
}

func (l *logPollerWrapper) HealthReport() map[string]error {
	return map[string]error{l.Name(): l.Ready()}
}

func (l *logPollerWrapper) Name() string { return l.lggr.Name() }

// methods of LogPollerWrapper
func (l *logPollerWrapper) LatestEvents() ([]evmRelayTypes.OracleRequest, []evmRelayTypes.OracleResponse, error) {
	l.mu.Lock()
	coordinators := []Coordinator{}
	if l.activeCoordinator != nil && l.activeCoordinator.Address() != (common.Address{}) {
		coordinators = append(coordinators, l.activeCoordinator)
	}
	if l.proposedCoordinator != nil && l.proposedCoordinator.Address() != (common.Address{}) && l.activeCoordinator != l.proposedCoordinator {
		coordinators = append(coordinators, l.proposedCoordinator)
	}
	latest, err := l.logPoller.LatestBlock()
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
		l.lggr.Debug("LatestEvents: no non-zero coordinators to check")
		return resultsReq, resultsResp, errors.New("no non-zero coordinators to check")
	}

	for _, coordinator := range coordinators {
		requestEndBlock := latestBlockNum - l.requestBlockOffset
		requestLogTopic, err := coordinator.OracleRequestLogTopic()
		if err != nil {
			l.lggr.Errorw("LatestEvents: ", err)
			return nil, nil, err
		}
		requestLogs, err := l.logPoller.Logs(startBlockNum, requestEndBlock, requestLogTopic, coordinator.Address())
		if err != nil {
			l.lggr.Errorw("LatestEvents: fetching request logs from LogPoller failed", "startBlock", startBlockNum, "endBlock", requestEndBlock)
			return nil, nil, err
		}
		l.lggr.Debugw("LatestEvents: fetched request logs", "nRequestLogs", len(requestLogs), "latestBlock", latest, "startBlock", startBlockNum, "endBlock", requestEndBlock)
		requestLogs = l.filterPreviouslyDetectedEvents(requestLogs, &l.detectedRequests, "requests")
		responseEndBlock := latestBlockNum - l.responseBlockOffset
		responseLogTopic, err := coordinator.OracleResponseLogTopic()
		if err != nil {
			l.lggr.Errorw("LatestEvents: ", err)
			return nil, nil, err
		}
		responseLogs, err := l.logPoller.Logs(startBlockNum, responseEndBlock, responseLogTopic, coordinator.Address())
		if err != nil {
			l.lggr.Errorw("LatestEvents: fetching response logs from LogPoller failed", "startBlock", startBlockNum, "endBlock", responseEndBlock)
			return nil, nil, err
		}
		l.lggr.Debugw("LatestEvents: fetched request logs", "nResponseLogs", len(responseLogs), "latestBlock", latest, "startBlock", startBlockNum, "endBlock", responseEndBlock)
		responseLogs = l.filterPreviouslyDetectedEvents(responseLogs, &l.detectedResponses, "responses")

		l.lggr.Debugw("LatestEvents: parsing logs", "nRequestLogs", len(requestLogs), "nResponseLogs", len(responseLogs), "coordinatorAddress", coordinator.Address().Hex())
		requests, err := coordinator.LogsToRequests(requestLogs)
		if err != nil {
			l.lggr.Errorf("LatestEvents: fetched oracle request: %w", err)
			return nil, nil, err
		}
		resultsReq = append(resultsReq, requests...)
		responses, err := coordinator.LogsToResponses(responseLogs)
		if err != nil {
			l.lggr.Errorf("LatestEvents: fetched oracle response: %w", err)
			return nil, nil, err
		}
		resultsResp = append(resultsResp, responses...)
	}

	l.lggr.Debugw("LatestEvents: done", "nRequestLogs", len(resultsReq), "nResponseLogs", len(resultsResp), "startBlock", startBlockNum, "endBlock", latestBlockNum)
	return resultsReq, resultsResp, nil
}

func (l *logPollerWrapper) filterPreviouslyDetectedEvents(logs []logpoller.Log, detectedEvents *detectedEvents, filterType string) []logpoller.Log {
	if len(logs) > maxLogsToProcess {
		l.lggr.Errorw("filterPreviouslyDetectedEvents: too many logs to process, only processing latest maxLogsToProcess logs", "filterType", filterType, "nLogs", len(logs), "maxLogsToProcess", maxLogsToProcess)
		logs = logs[len(logs)-maxLogsToProcess:]
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	filteredLogs := []logpoller.Log{}
	for _, log := range logs {
		var requestId [32]byte
		if len(log.Topics) < 2 || len(log.Topics[1]) != 32 {
			l.lggr.Errorw("filterPreviouslyDetectedEvents: invalid log, skipping", "filterType", filterType, "log", log)
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
	l.lggr.Debugw("filterPreviouslyDetectedEvents: done", "filterType", filterType, "nLogs", len(logs), "nFilteredLogs", len(filteredLogs), "nExpiredRequests", expiredRequests, "previouslyDetectedCacheSize", len(detectedEvents.detectedEventsOrdered))
	return filteredLogs
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.pluginConfig.ContractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		if err := subscriber.UpdateRoutes(l.routerContract.Address(), l.routerContract.Address()); err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "err", err)
		}
	} else if l.pluginConfig.ContractVersion == 1 {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.subscribers[subscriberName] = subscriber
	}
}

func (l *logPollerWrapper) checkForRouteUpdates() {
	defer l.closeWait.Done()
	freqSec := l.pluginConfig.ContractUpdateCheckFrequencySec
	if freqSec == 0 {
		l.lggr.Errorw("LogPollerWrapper: ContractUpdateCheckFrequencySec is zero - route update checks disabled")
		return
	}

	updateOnce := func() {
		// NOTE: timeout == frequency here, could be changed to a separate config value
		timeoutCtx, cancel := utils.ContextFromChanWithTimeout(l.stopCh, time.Duration(l.pluginConfig.ContractUpdateCheckFrequencySec)*time.Second)
		defer cancel()
		active, proposed, err := l.getCurrentCoordinators(timeoutCtx)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: error calling getCurrentCoordinators", "err", err)
			return
		}
		l.handleRouteUpdate(active, proposed)
	}

	updateOnce() // update once right away
	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.stopCh:
			return
		case <-ticker.C:
			updateOnce()
		}
	}
}

func (l *logPollerWrapper) getCurrentCoordinators(ctx context.Context) (common.Address, common.Address, error) {
	if l.pluginConfig.ContractVersion == 0 {
		return l.routerContract.Address(), l.routerContract.Address(), nil
	}
	var donId [32]byte
	copy(donId[:], []byte(l.pluginConfig.DONID))

	activeCoordinatorAddress, err := l.routerContract.GetContractById(&bind.CallOpts{
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
		return activeCoordinatorAddress, l.proposedCoordinator.Address(), nil
	}

	return activeCoordinatorAddress, proposedCoordinator, nil
}

func (l *logPollerWrapper) getTypeAndVersion(contractAddress common.Address) (string, error) {
	if contractAddress == (common.Address{}) {
		l.lggr.Debug("LogPollerWrapper: cannot get typeAndVersion from an unset address")
		return "", nil
	}

	contract, err := type_and_version.NewTypeAndVersionInterface(contractAddress, l.client)
	if err != nil {
		l.lggr.Error("LogPollerWrapper: could not initialize contract with typeAndVersion interface", contractAddress)
		return "", err
	}

	typeAndVersion, err := contract.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		l.lggr.Error("LogPollerWrapper: could not get typeAndVersion from contract ", contractAddress)
		return "", err
	}

	return typeAndVersion, nil
}

func (l *logPollerWrapper) handleRouteUpdate(activeCoordinatorAddress common.Address, proposedCoordinatorAddress common.Address) {
	l.mu.Lock()
	defer l.mu.Unlock()

	commitmentABIV1 := abi.Arguments{
		{Type: l.abiTypes.bytes32Type}, // RequestId
		{Type: l.abiTypes.addressType}, // Coordinator
		{Type: l.abiTypes.uint96Type},  // EstimatedTotalCostJuels
		{Type: l.abiTypes.addressType}, // Client
		{Type: l.abiTypes.uint64Type},  // SubscriptionId
		{Type: l.abiTypes.uint32Type},  // CallbackGasLimit
		{Type: l.abiTypes.uint72Type},  // AdminFee
		{Type: l.abiTypes.uint72Type},  // DonFee
		{Type: l.abiTypes.uint40Type},  // GasOverheadBeforeCallback
		{Type: l.abiTypes.uint40Type},  // GasOverheadAfterCallback
		{Type: l.abiTypes.uint32Type},  // TimeoutTimestamp
	}

	commitmentABIV2 := append(commitmentABIV1,
		abi.Argument{Type: l.abiTypes.uint72Type}) // OperationFee

	if activeCoordinatorAddress == (common.Address{}) {
		l.lggr.Error("LogPollerWrapper: cannot update activeCoordinator to zero address")
		return
	}

	if (l.activeCoordinator != nil && l.activeCoordinator.Address() == activeCoordinatorAddress) &&
		(l.proposedCoordinator != nil && l.proposedCoordinator.Address() == proposedCoordinatorAddress) {
		l.lggr.Debug("LogPollerWrapper: no changes to routes")
		return
	}

	activeCoordinatorTypeAndVersion, err := l.getTypeAndVersion(activeCoordinatorAddress)
	if err != nil {
		l.lggr.Errorf("LogPollerWrapper: failed to get active coordinatorTypeAndVersion: %w", err)
		return
	}
	var activeCoordinator Coordinator
	switch {
	case strings.Contains(activeCoordinatorTypeAndVersion, FUNCTIONS_COORDINATOR_VERSION_1_SUBSTRING):
		activeCoordinator = NewCoordinatorV1(activeCoordinatorAddress, commitmentABIV1, l.client, l.logPoller, l.lggr)
	case strings.Contains(activeCoordinatorTypeAndVersion, FUNCTIONS_COORDINATOR_VERSION_2_SUBSTRING):
		activeCoordinator = NewCoordinatorV2(activeCoordinatorAddress, commitmentABIV2, l.client, l.logPoller, l.lggr)
	default:
		l.lggr.Errorf("LogPollerWrapper: Invalid active coordinator type and version: %q", activeCoordinatorTypeAndVersion)
		return
	}

	if activeCoordinator != nil {
		err = activeCoordinator.RegisterFilters()
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to register active coordinator filters", err)
			return
		}
		l.activeCoordinator = activeCoordinator
		l.lggr.Debugw("LogPollerWrapper: new routes", "activeCoordinator", activeCoordinator.Address().Hex())
	}

	proposedCoordinatorTypeAndVersion, err := l.getTypeAndVersion(proposedCoordinatorAddress)
	if err != nil {
		l.lggr.Errorf("LogPollerWrapper: failed to get proposed coordinatorTypeAndVersion: %w", err)
		return
	}

	var proposedCoordinator Coordinator
	switch {
	// proposedCoordinatorTypeAndVersion can be empty due to an empty proposedCoordinatorAddress
	case proposedCoordinatorTypeAndVersion == "":
		proposedCoordinator = NewCoordinatorV1(proposedCoordinatorAddress, commitmentABIV1, l.client, l.logPoller, l.lggr)
	case strings.Contains(proposedCoordinatorTypeAndVersion, FUNCTIONS_COORDINATOR_VERSION_1_SUBSTRING):
		proposedCoordinator = NewCoordinatorV1(proposedCoordinatorAddress, commitmentABIV1, l.client, l.logPoller, l.lggr)
	case strings.Contains(proposedCoordinatorTypeAndVersion, FUNCTIONS_COORDINATOR_VERSION_2_SUBSTRING):
		proposedCoordinator = NewCoordinatorV2(proposedCoordinatorAddress, commitmentABIV2, l.client, l.logPoller, l.lggr)

	}

	if proposedCoordinator != nil {
		err = proposedCoordinator.RegisterFilters()
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to register proposed coordinator filters", err)
			return
		}
		l.proposedCoordinator = proposedCoordinator
		l.lggr.Debugw("LogPollerWrapper: new routes", "proposedCoordinator", proposedCoordinator.Address().Hex())
	}

	for _, subscriber := range l.subscribers {
		err := subscriber.UpdateRoutes(activeCoordinator.Address(), proposedCoordinator.Address())
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "err", err)
		}
	}
}

func initAbiTypes() (*abiTypes, error) {
	var err error
	uint32Type, err := abi.NewType("uint32", "uint32", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize uint32Type type: %w", err)
	}
	uint40Type, err := abi.NewType("uint40", "uint40", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize uint40Type type: %w", err)
	}
	uint64Type, err := abi.NewType("uint64", "uint64", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize uint64Type type: %w", err)
	}
	uint72Type, err := abi.NewType("uint72", "uint72", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize uint72Type type: %w", err)
	}
	uint96Type, err := abi.NewType("uint96", "uint96", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize uint96Type type: %w", err)
	}
	addressType, err := abi.NewType("address", "address", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize addressType type: %w", err)
	}
	bytes32Type, err := abi.NewType("bytes32", "bytes32", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bytes32Type type: %w", err)
	}

	at := &abiTypes{
		uint32Type:  uint32Type,
		uint40Type:  uint40Type,
		uint64Type:  uint64Type,
		uint72Type:  uint72Type,
		uint96Type:  uint96Type,
		addressType: addressType,
		bytes32Type: bytes32Type,
	}

	return at, nil
}
