package gateway

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	workflowsyncerv2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/v2"
)

const (
	FunctionsHandlerType   HandlerType = "functions"
	DummyHandlerType       HandlerType = "dummy"
	WebAPICapabilitiesType HandlerType = "web-api-capabilities" //  Handler for v0.1 HTTP capabilities for DAG workflows
	HTTPCapabilityType     HandlerType = "http-capabilities"    // Handler for v1.0 HTTP capabilities for NoDAG workflows
	VaultHandlerType       HandlerType = "vault"
)

type handlerFactory struct {
	legacyChains           legacyevm.LegacyChainContainer
	ds                     sqlutil.DataSource
	lggr                   logger.Logger
	httpClient             network.HTTPClient
	capabilitiesRegistry   core.CapabilitiesRegistry
	workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer
	lf                     limits.Factory
}

var _ HandlerFactory = (*handlerFactory)(nil)

func NewHandlerFactory(legacyChains legacyevm.LegacyChainContainer, ds sqlutil.DataSource, httpClient network.HTTPClient, capabilitiesRegistry core.CapabilitiesRegistry, workflowRegistrySyncer workflowsyncerv2.WorkflowRegistrySyncer, lggr logger.Logger, lf limits.Factory) HandlerFactory {
	return &handlerFactory{
		legacyChains,
		ds,
		lggr,
		httpClient,
		capabilitiesRegistry,
		workflowRegistrySyncer,
		lf,
	}
}

func (hf *handlerFactory) NewHandler(
	handlerType HandlerType,
	handlerConfig json.RawMessage,
	shardedDONs []config.ShardedDONConfig,
	shardsConnMgrs [][]handlers.DON,
) (handlers.Handler, error) {
	if len(shardedDONs) == 0 || len(shardsConnMgrs) == 0 {
		return nil, errors.New("at least one DON and connection manager required")
	}
	if len(shardsConnMgrs[0]) == 0 {
		return nil, errors.New("at least one shard connection manager required")
	}

	// For backward compatibility, convert sharded config to legacy DONConfig
	// using the first DON's first shard. Handlers that support sharding can
	// use the full shardedDONs and shardsConnMgrs directly.
	donConfig := shardedDONsToLegacy(shardedDONs[0])
	don := shardsConnMgrs[0][0]

	switch handlerType {
	case FunctionsHandlerType:
		return functions.NewFunctionsHandlerFromConfig(handlerConfig, donConfig, don, hf.legacyChains, hf.ds, hf.lggr)
	case DummyHandlerType:
		return handlers.NewDummyHandler(donConfig, don, hf.lggr)
	case WebAPICapabilitiesType:
		return capabilities.NewHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr)
	case HTTPCapabilityType:
		return v2.NewGatewayHandler(handlerConfig, donConfig, don, hf.httpClient, hf.lggr, hf.lf)
	case VaultHandlerType:
		requestAuthorizer := vaultcap.NewRequestAuthorizer(hf.lggr, hf.workflowRegistrySyncer)
		return vault.NewHandler(handlerConfig, donConfig, don, hf.capabilitiesRegistry, requestAuthorizer, hf.lggr, clockwork.NewRealClock(), hf.lf)
	default:
		return nil, fmt.Errorf("unsupported handler type %s", handlerType)
	}
}

func shardedDONsToLegacy(shardedDON config.ShardedDONConfig) *config.DONConfig {
	var members []config.NodeConfig
	if len(shardedDON.Shards) > 0 {
		members = shardedDON.Shards[0].Nodes
	}
	return &config.DONConfig{
		DonId:   shardedDON.DonName,
		F:       shardedDON.F,
		Members: members,
	}
}
