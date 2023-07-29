package functions

import (
	"context"
	"encoding/hex"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type logPollerWrapper struct {
	utils.StartStopOnce
	contractAddress     common.Address
	client              client.Client
	pluginConfig        config.PluginConfig
	logPoller           logpoller.LogPoller
	subscribers         map[string]evmRelayTypes.RouteUpdateSubscriber
	mu                  sync.RWMutex
	lggr                logger.Logger
	shutdownWaitGroup   sync.WaitGroup
	serviceContext      context.Context
	serviceCancel       context.CancelFunc
	chStop              chan struct{}
	previousCoordinator common.Address
	activeCoordinator   common.Address
	proposedCoordinator common.Address
}

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(contractAddress common.Address, client client.Client, pluginConfig config.PluginConfig, logPoller logpoller.LogPoller, lggr logger.Logger) evmRelayTypes.LogPollerWrapper {
	return &logPollerWrapper{
		contractAddress: contractAddress,
		client:          client,
		pluginConfig:    pluginConfig,
		logPoller:       logPoller,
		subscribers:     make(map[string]evmRelayTypes.RouteUpdateSubscriber),
		lggr:            lggr,
		chStop:          make(chan struct{}),
	}
}

// TODO(FUN-381): Implement LogPollerWrapper with an API suiting all users.
func (l *logPollerWrapper) Start(context.Context) error {
	return l.StartOnce("LogPollerWrapper", func() error {
		l.serviceContext, l.serviceCancel = context.WithCancel(context.Background())
		l.shutdownWaitGroup.Add(4)

		if l.pluginConfig.ContractVersion > 0 {
			// Set up initial Log Poller filters
			l.registerRouter(l.contractAddress)

			activeCoordinator, proposedCoordinator := l.getRouteContracts()
			l.registerCoordinator(activeCoordinator)
			if activeCoordinator != proposedCoordinator {
				// Active may equal Proposed, only watch if they are different
				l.registerCoordinator(proposedCoordinator)
			}

			// Periodically check for updates
			l.checkForRouteUpdate()
		}

		go func() {
			<-l.chStop
			l.unregisterRouter(l.contractAddress)
			zeroAddress := common.Address{}
			if l.previousCoordinator != zeroAddress {
				l.unregisterCoordinator(l.previousCoordinator)
			}
			if l.activeCoordinator != zeroAddress {
				l.unregisterCoordinator(l.activeCoordinator)
			}
			if l.proposedCoordinator != zeroAddress {
				l.unregisterCoordinator(l.proposedCoordinator)
			}
			l.shutdownWaitGroup.Done()
		}()

		return nil
	})
}

func (l *logPollerWrapper) Close() error {
	return nil
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

// "public" methods of LogPollerWrapper
func (l *logPollerWrapper) LatestRequests() ([]evmRelayTypes.OracleRequest, error) {
	// TODO poll both active and proposed coordinators for requests and parse them

	// logs, err := lp.Logs(
	// 	end-lookbackBlocks,
	// 	end,
	// 	iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
	// 	registryAddress,
	// 	pg.WithParentCtx(ctx),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	// }

	// TODO:
	// format logs

	return nil, nil
}

func (l *logPollerWrapper) LatestResponses() ([]evmRelayTypes.OracleResponse, error) {
	// TODO poll both active and proposed coordinators for requests and parse them

	// logs, err := lp.Logs(
	// 	end-lookbackBlocks,
	// 	end,
	// 	iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
	// 	registryAddress,
	// 	pg.WithParentCtx(ctx),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	// }

	// TODO:
	// format logs

	return nil, nil
}

func (l *logPollerWrapper) LatestUpdates() ([]evmRelayTypes.RouteUpdate, error) {
	// TODO poll both active and proposed coordinators for requests and parse them

	// logs, err := lp.Logs(
	// 	end-lookbackBlocks,
	// 	end,
	// 	iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
	// 	registryAddress,
	// 	pg.WithParentCtx(ctx),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	// }

	// TODO:
	// format logs

	return nil, nil
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.pluginConfig.ContractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		err := subscriber.UpdateRoutes(l.contractAddress, l.contractAddress)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "error", err)
		}
	} else if l.pluginConfig.ContractVersion == 1 {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.subscribers[subscriberName] = subscriber
		err := subscriber.UpdateRoutes(l.contractAddress, l.contractAddress)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "error", err)
		}
	}
}

func (l *logPollerWrapper) checkForRouteUpdate() {
	defer l.shutdownWaitGroup.Done()
	freqSec := l.pluginConfig.ContractUpdateCheckFrequencySec
	if freqSec == 0 {
		l.lggr.Errorw("ContractUpdateCheckFrequencySec must set to more than 0 in PluginConfig")
		return
	}
	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.chStop:
			return
		case <-ticker.C:
			updates, err := l.LatestUpdates()
			if err != nil {
				l.lggr.Errorw("LogPoller Wrapper: unable to get latest update logs", "error", err)
				break
			}
			for _, update := range updates {
				activeAddress := common.Address{}
				activeAddress.SetBytes(update.ActiveAddress)
				proposedAddress := common.Address{}
				proposedAddress.SetBytes(update.ProposedAddress)
				l.handleRouteUpdate(activeAddress, proposedAddress)
			}
		}
	}
}

func (l *logPollerWrapper) getRouteContracts() (activeCoordinator common.Address, proposedCoordinator common.Address) {
	routerContract, err := functions_router.NewFunctionsRouter(l.contractAddress, l.client)
	if err != nil {
		l.lggr.Errorw("LogPoller Wrapper: unable to initialize Router contract", err)
	}

	var donId [32]byte
	donIdBytes, err := hex.DecodeString(l.pluginConfig.DONId)
	copy(donId[:], donIdBytes)

	activeCoordinator, err = routerContract.GetContractById(&bind.CallOpts{}, donId, false)
	proposedCoordinator, err = routerContract.GetContractById(&bind.CallOpts{}, donId, true)
	return activeCoordinator, proposedCoordinator
}

func (l *logPollerWrapper) handleRouteUpdate(activeCoordinator common.Address, proposedCoordinator common.Address) {
	l.previousCoordinator = l.activeCoordinator
	l.activeCoordinator = activeCoordinator
	l.proposedCoordinator = proposedCoordinator

	// Register filters for new proposedCoordinator
	err := l.registerCoordinator(proposedCoordinator)
	if err != nil {
		l.lggr.Errorw("LogPoller Wrapper: unable to register Log Poller filters for new coordinator", "coordinatorAddress", proposedCoordinator, "error", err)
	}

	// TODO: wait request timeout seconds and then de-register old filters?

	// Notify subscribers of change
	for name, subscriber := range l.subscribers {
		err := subscriber.UpdateRoutes(activeCoordinator, proposedCoordinator)
		if err != nil {
			l.lggr.Errorw("LogPoller Wrapper: unable to notify subscriber", "subscriberName", name, "error", err)
		}
	}
}

func (l *logPollerWrapper) registerFilters(filters []logpoller.Filter) error {
	for _, lpFilter := range filters {
		if err := l.logPoller.RegisterFilter(lpFilter); err != nil {
			return err
		}
	}
	return nil
}

func (l *logPollerWrapper) unregisterFilters(filters []logpoller.Filter, q pg.Queryer) error {
	for _, lpFilter := range filters {
		if err := l.logPoller.UnregisterFilter(lpFilter.Name, q); err != nil {
			return err
		}
	}
	return nil
}

func (l *logPollerWrapper) registerCoordinator(coordinatorAddress common.Address) error {
	oracleRequestFilters := getFiltersOracleRequest(coordinatorAddress)
	err := l.registerFilters(oracleRequestFilters)
	if err != nil {
		return err
	}

	oracleResponseFilters := getFiltersOracleResponse(coordinatorAddress)
	err = l.registerFilters(oracleResponseFilters)
	if err != nil {
		return err
	}

	return nil
}

func (l *logPollerWrapper) unregisterCoordinator(coordinatorAddress common.Address) error {
	oracleRequestFilters := getFiltersOracleRequest(coordinatorAddress)
	err := l.unregisterFilters(oracleRequestFilters, nil)
	if err != nil {
		return err
	}

	oracleResponseFilters := getFiltersOracleResponse(coordinatorAddress)
	err = l.unregisterFilters(oracleResponseFilters, nil)
	if err != nil {
		return err
	}

	return nil
}

func (l *logPollerWrapper) registerRouter(routerAddress common.Address) error {
	routeUpdateFilters := getFiltersRouteUpdate(routerAddress)
	err := l.registerFilters(routeUpdateFilters)
	if err != nil {
		return err
	}
	return nil
}

func (l *logPollerWrapper) unregisterRouter(routerAddress common.Address) error {
	routeUpdateFilters := getFiltersRouteUpdate(routerAddress)
	err := l.unregisterFilters(routeUpdateFilters, nil)
	if err != nil {
		return err
	}
	return nil
}

func FilterName(addr common.Address, filterType string) string {
	return logpoller.FilterName("Chainlink Functions", filterType, addr.String())
}

func getFiltersOracleRequest(coordinatorAddress common.Address) []logpoller.Filter {
	name := "OracleRequest"
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(FilterName(coordinatorAddress, name)),
			EventSigs: []common.Hash{functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic()},
			Addresses: []common.Address{coordinatorAddress},
		},
	}
}

func getFiltersOracleResponse(coordinatorAddress common.Address) []logpoller.Filter {
	name := "OracleResponse"
	return []logpoller.Filter{
		{
			Name: logpoller.FilterName(FilterName(coordinatorAddress, name)),
			EventSigs: []common.Hash{
				functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(),
				functions_coordinator.FunctionsCoordinatorInvalidRequestID{}.Topic(),
				functions_coordinator.FunctionsCoordinatorInsufficientGasProvided{}.Topic(),
				functions_coordinator.FunctionsCoordinatorCostExceedsCommitment{}.Topic(),
				functions_coordinator.FunctionsCoordinatorInsufficientSubscriptionBalance{}.Topic(),
			},
			Addresses: []common.Address{coordinatorAddress},
		},
	}
}

func getFiltersRouteUpdate(routerAddress common.Address) []logpoller.Filter {
	name := "RouteUpdate"
	return []logpoller.Filter{
		{
			Name: logpoller.FilterName(FilterName(routerAddress, name)),
			EventSigs: []common.Hash{
				functions_router.FunctionsRouterContractProposed{}.Topic(),
				functions_router.FunctionsRouterContractUpdated{}.Topic(),
			},
			Addresses: []common.Address{routerAddress},
		},
	}
}
