package functions

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type logPollerWrapper struct {
	logPoller logpoller.LogPoller
	lggr      logger.Logger
}

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(logPoller logpoller.LogPoller, lggr logger.Logger) evmRelayTypes.LogPollerWrapper {
	return &logPollerWrapper{
		logPoller: logPoller,
		lggr:      lggr,
	}
}

// TODO(FUN-381): Implement LogPollerWrapper with an API suiting all users.
func (l *logPollerWrapper) Start(context.Context) error {
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

func (l *logPollerWrapper) LatestRoutes() (activeCoordinator common.Address, proposedCoordinator common.Address, err error) {
	return common.Address{}, common.Address{}, nil
}
