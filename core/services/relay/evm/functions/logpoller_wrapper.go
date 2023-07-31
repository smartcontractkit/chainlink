package functions

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type logPollerWrapper struct {
	utils.StartStopOnce

	routerContract      *functions_router.FunctionsRouter
	pluginConfig        config.PluginConfig
	client              client.Client
	logPoller           logpoller.LogPoller
	subscribers         map[string]evmRelayTypes.RouteUpdateSubscriber
	activeCoordinator   common.Address
	proposedCoordinator common.Address
	nextBlock           int64
	mu                  sync.Mutex
	closeWait           sync.WaitGroup
	stopCh              utils.StopChan
	lggr                logger.Logger
}

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(routerContractAddress common.Address, pluginConfig config.PluginConfig, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) (evmRelayTypes.LogPollerWrapper, error) {
	routerContract, err := functions_router.NewFunctionsRouter(routerContractAddress, client)
	if err != nil {
		return nil, err
	}

	return &logPollerWrapper{
		routerContract: routerContract,
		pluginConfig:   pluginConfig,
		logPoller:      logPoller,
		client:         client,
		subscribers:    make(map[string]evmRelayTypes.RouteUpdateSubscriber),
		stopCh:         make(utils.StopChan),
		lggr:           lggr,
	}, nil
}

func (l *logPollerWrapper) Start(context.Context) error {
	return l.StartOnce("LogPollerWrapper", func() error {
		l.lggr.Info("starting LogPollerWrapper")
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.pluginConfig.ContractVersion == 0 {
			l.activeCoordinator = l.routerContract.Address()
			l.proposedCoordinator = l.routerContract.Address()
		} else if l.pluginConfig.ContractVersion == 1 {
			nextBlock, err := l.logPoller.LatestBlock()
			l.nextBlock = nextBlock
			if err != nil {
				l.lggr.Error("LogPollerWrapper: LatestBlock() failed, starting from 0")
			} else {
				l.lggr.Debugw("LogPollerWrapper: LatestBlock() got starting block", "block", nextBlock)
			}
			l.closeWait.Add(1)
			go l.checkForRouteUpdates()
		}
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
	return make(map[string]error)
}

func (l *logPollerWrapper) Name() string {
	return "LogPollerWrapper"
}

func (l *logPollerWrapper) Ready() error {
	return nil
}

// methods of LogPollerWrapper
func (l *logPollerWrapper) LatestEvents() ([]evmRelayTypes.OracleRequest, []evmRelayTypes.OracleResponse, error) {
	l.mu.Lock()
	coordinators := []common.Address{l.activeCoordinator}
	if l.activeCoordinator != l.proposedCoordinator {
		coordinators = append(coordinators, l.proposedCoordinator)
	}
	nextBlock := l.nextBlock
	latest, err := l.logPoller.LatestBlock()
	if err != nil {
		l.mu.Unlock()
		return nil, nil, err
	}
	if latest >= nextBlock {
		l.nextBlock = latest + 1
	}
	l.mu.Unlock()

	// outside of the lock
	resultsReq := []evmRelayTypes.OracleRequest{}
	resultsResp := []evmRelayTypes.OracleResponse{}
	if latest < nextBlock {
		l.lggr.Debugw("LatestEvents: no new blocks to check", "latest", latest, "nextBlock", nextBlock)
		return resultsReq, resultsResp, nil
	}

	for _, coordinator := range coordinators {
		requestLogs, err := l.logPoller.Logs(nextBlock, latest, functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(), coordinator)
		if err != nil {
			l.lggr.Errorw("LatestEvents: fetching request logs from LogPoller failed", "latest", latest, "nextBlock", nextBlock)
			return nil, nil, err
		}
		responseLogs, err := l.logPoller.Logs(nextBlock, latest, functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(), coordinator)
		if err != nil {
			l.lggr.Errorw("LatestEvents: fetching response logs from LogPoller failed", "latest", latest, "nextBlock", nextBlock)
			return nil, nil, err
		}

		parsingContract, err := functions_coordinator.NewFunctionsCoordinator(coordinator, l.client)
		if err != nil {
			l.lggr.Error("LatestEvents: creating a contract instance for parsing failed")
			return nil, nil, err
		}

		l.lggr.Debugw("LatestEvents: parsing logs", "nRequestLogs", len(requestLogs), "nResponseLogs", len(responseLogs), "coordinatorAddress", coordinator.Hex())
		for _, log := range requestLogs {
			gethLog := log.ToGethLog()
			oracleRequest, err := parsingContract.ParseOracleRequest(gethLog)
			if err != nil {
				l.lggr.Errorw("LatestEvents: failed to parse a request log, skipping")
				continue
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
				CoordinatorContract: coordinator,
			})
		}
		for _, log := range responseLogs {
			gethLog := log.ToGethLog()
			oracleResponse, err := parsingContract.ParseOracleResponse(gethLog)
			if err != nil {
				l.lggr.Errorw("LatestEvents: failed to parse a response log, skipping")
				continue
			}
			resultsResp = append(resultsResp, evmRelayTypes.OracleResponse{
				RequestId: oracleResponse.RequestId,
			})
		}
	}

	l.lggr.Debugw("LatestEvents: done", "nRequestLogs", len(resultsReq), "nResponseLogs", len(resultsResp), "nextBlock", nextBlock, "latest", latest)
	return resultsReq, resultsResp, nil
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.pluginConfig.ContractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		err := subscriber.UpdateRoutes(l.routerContract.Address(), l.routerContract.Address())
		l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "error", err)
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
		l.lggr.Errorw("ContractUpdateCheckFrequencySec is zero - route update checks disabled")
		return
	}

	updateOnce := func() {
		// NOTE: timeout == frequency here, could be changed to a separate config value
		timeoutCtx, cancel := utils.ContextFromChanWithTimeout(l.stopCh, time.Duration(l.pluginConfig.ContractUpdateCheckFrequencySec)*time.Second)
		defer cancel()
		active, proposed, err := l.getCurrentCoordinators(timeoutCtx)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: error calling getCurrentCoordinators", "err", err)
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
	copy(donId[:], []byte(l.pluginConfig.DONId))

	activeCoordinator, err := l.routerContract.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId, false)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	proposedCoordinator, err := l.routerContract.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId, true)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	return activeCoordinator, proposedCoordinator, nil
}

func (l *logPollerWrapper) handleRouteUpdate(activeCoordinator common.Address, proposedCoordinator common.Address) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if activeCoordinator == l.activeCoordinator && proposedCoordinator == l.proposedCoordinator {
		l.lggr.Debug("LogPollerWrapper: no changes to routes")
		return
	}
	errActive := l.registerFilters(activeCoordinator)
	errProposed := l.registerFilters(proposedCoordinator)
	if errActive != nil || errProposed != nil {
		l.lggr.Errorw("LogPollerWrapper: Failed to register filters", "errorActive", errActive, "errorProposed", errProposed)
		return
	}

	l.lggr.Debugw("LogPollerWrapper: new routes", "activeCoordinator", activeCoordinator.Hex(), "proposedCoordinator", proposedCoordinator.Hex())
	l.activeCoordinator = activeCoordinator
	l.proposedCoordinator = proposedCoordinator

	for _, subscriber := range l.subscribers {
		err := subscriber.UpdateRoutes(activeCoordinator, proposedCoordinator)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "error", err)
		}
	}
}

func filterName(addr common.Address) string {
	return logpoller.FilterName("FunctionsLogPollerWrapper", addr.String())
}

func (l *logPollerWrapper) registerFilters(coordinatorAddress common.Address) error {
	return l.logPoller.RegisterFilter(
		logpoller.Filter{
			Name: filterName(coordinatorAddress),
			EventSigs: []common.Hash{
				functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(),
				functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(),
			},
			Addresses: []common.Address{coordinatorAddress},
		})
}
