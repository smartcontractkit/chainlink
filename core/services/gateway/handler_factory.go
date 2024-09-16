package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	FunctionsHandlerType HandlerType = "functions"
	DummyHandlerType     HandlerType = "dummy"
	WorkflowHandlerType  HandlerType = "workflow"
)

type handlerFactory struct {
	legacyChains legacyevm.LegacyChainContainer
	ds           sqlutil.DataSource
	lggr         logger.Logger
	httpClient   network.HttpClient
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, httpClient network.HttpClient, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{
		legacyChains,
		ds,
		lggr,
		httpClient,
	}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.ds, hf.lggr)
	case WorkflowHandlerType:
		return workflow.NewWorkflowHandler(donConfig, don, hf.httpClient, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
