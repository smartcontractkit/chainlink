package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

const (
	FunctionsHandlerType HandlerType = "functions"
	DummyHandlerType     HandlerType = "dummy"
)

type handlerFactory struct {
	chains evm.ChainSet
	lggr   logger.Logger
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(chains evm.ChainSet, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{chains, lggr}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandler(handlerConfig, donConfig, don, hf.chains, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
