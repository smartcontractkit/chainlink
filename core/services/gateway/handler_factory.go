package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/chain"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
)

const (
	FunctionsHandlerType HandlerType = "functions"
	ChainHandlerType     HandlerType = "chain"
	DummyHandlerType     HandlerType = "dummy"
)

type handlerFactory struct {
	legacyChains legacyevm.LegacyChainContainer
	ds           sqlutil.DataSource
	relayGetter  handlers.RelayGetter
	lggr         logger.Logger
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, relayGetter handlers.RelayGetter, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{legacyChains, ds, relayGetter, lggr}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.ds, hf.lggr)
	case ChainHandlerType:
		return chain.NewChainHandler(hf.relayGetter, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
