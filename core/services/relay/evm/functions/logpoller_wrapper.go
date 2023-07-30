package functions

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type logPollerWrapper struct {
	routerContract  common.Address
	contractVersion uint32
	logPoller       logpoller.LogPoller
	subscribers     map[string]evmRelayTypes.RouteUpdateSubscriber
	mu              sync.RWMutex
	lggr            logger.Logger
}

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(routerContract common.Address, contractVersion uint32, logPoller logpoller.LogPoller, lggr logger.Logger) evmRelayTypes.LogPollerWrapper {
	return &logPollerWrapper{
		routerContract:  routerContract,
		contractVersion: contractVersion,
		logPoller:       logPoller,
		subscribers:     make(map[string]evmRelayTypes.RouteUpdateSubscriber),
		lggr:            lggr,
	}
}

// TODO(FUN-381): Implement LogPollerWrapper with an API suiting all users.
func (l *logPollerWrapper) Start(context.Context) error {
	// TODO periodically update active and proposed coordinator contract addresses
	// update own Logpoller filters accordingly
	// push to all subscribers
	return nil
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
	return nil, nil
}

func (l *logPollerWrapper) LatestResponses() ([]evmRelayTypes.OracleResponse, error) {
	// TODO poll both active and proposed coordinators for responses and parse them
	return nil, nil
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.contractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		err := subscriber.UpdateRoutes(l.routerContract, l.routerContract)
		l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "error", err)
	} else if l.contractVersion == 1 {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.subscribers[subscriberName] = subscriber

		// TODO remove when periodic updates are ready
		_ = subscriber.UpdateRoutes(l.routerContract, l.routerContract)
	}
}
