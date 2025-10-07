package httpaction

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.HTTPActionCapability

type HTTPAction struct{}

func (o *HTTPAction) Flag() cre.CapabilityFlag {
	return flag
}

func (o *HTTPAction) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs, // move to Topology
	contractVersions map[string]string,
	gatewayConfigs map[uint64]config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	dons := topology.DonsWithFlag(o.Flag())
	if len(dons) == 0 {
		return nil, nil
	}

	for _, don := range dons {
		workers, wErr := don.Workers()
		if wErr != nil {
			return nil, wErr
		}
		evmKey, ok := workers[0].Keys.EVM[registryChainSelector]
		if !ok {
			return nil, fmt.Errorf("worker node at index %d does not have EVM key for chain selector %d", workers[0].Index, registryChainSelector)
		}

		// donID is the ID of DON that runs the gateway, not the capability
		for _, gc := range gatewayConfigs {
			// for each DON, we need to add a handler config specific for this capability

			donFound := false
			for _, maybeDON := range gc.Dons {
				// we need to find DON configuration that matches current don
				for _, member := range maybeDON.Members {
					// if any of the member's address matches the EVM key of the worker node, we found the right DON
					if member.Address == evmKey.PublicAddress.Hex() {
						donFound = true
						break
					}
				}

				if donFound {
					maybeDON.Handlers = append(maybeDON.Handlers, config.Handler{
						ServiceName: "workflows",
						// TODO - figure out correct json syntax
						Config: []byte(`
[gatewayConfig.Dons.Handlers.Config]
maxTriggerRequestDurationMs = 5_000
metadataPullIntervalMs = 1_000
metadataAggregationIntervalMs = 1_000
[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`),
					})

					break
				}
			}

			// TODO this is not entirely correct, if don doesn't need web api gateway it will be missing from this configuration and thus should be added
			// TODO make web api gateway also a feature! <--------
			if !donFound {
				return nil, fmt.Errorf("could not find DON configuration for DON named %s in gateway config, even though it requires a gateway", don.Name)
			}
		}
	}

	return nil, nil
}

func (o *HTTPAction) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
) error {
	return nil
}
