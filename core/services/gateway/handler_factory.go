package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	FunctionsHandlerType   HandlerType = "functions"
	DummyHandlerType       HandlerType = "dummy"
	WebAPICapabilitiesType HandlerType = "web-api-capabilities" //  Handler for v0.1 HTTP capabilities for DAG workflows
	HTTPCapabilityType     HandlerType = "http-capabilities"    // Handler for v1.0 HTTP capabilities for NoDAG workflows
	VaultHandlerType       HandlerType = "vault"
)

type handlerFactory struct {
	legacyChains legacyevm.LegacyChainContainer
	ds           sqlutil.DataSource
	lggr         logger.Logger
	httpClient   network.HTTPClient
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, httpClient network.HTTPClient, lggr logger.Logger) HandlerFactory {
	return &handlerFactory{
		legacyChains,
		ds,
		lggr,
		httpClient,
	}
}

func HandlerTypeForMethod(method string) (HandlerType, error) {
	switch method {
	case functions.MethodSecretsSet, functions.MethodSecretsList, functions.MethodHeartbeat:
		return FunctionsHandlerType, nil
	case capabilities.MethodWebAPITarget, capabilities.MethodWebAPITrigger, capabilities.MethodComputeAction, capabilities.MethodWorkflowSyncer:
		return WebAPICapabilitiesType, nil
	case gateway_common.MethodHTTPAction, gateway_common.MethodPushWorkflowMetadata, gateway_common.MethodPullWorkflowMetadata:
		return HTTPCapabilityType, nil
	case vault.MethodSecretsCreate:
		return VaultHandlerType, nil
	case "dummy":
		return DummyHandlerType, nil
	default:
		return HandlerType(""), fmt.Errorf("unknown handler type for method %s", method)
	}
}

func (hf *handlerFactory) NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error) {
	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.ds, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	case WebAPICapabilitiesType:
		return capabilities.NewHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr)
	case HTTPCapabilityType:
		return v2.NewGatewayHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr)
	case VaultHandlerType:
		return vault.NewHandler(handlerConfig, donConfig, don, hf.lggr)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}
