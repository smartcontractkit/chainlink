package httptrigger

import (
	"context"
	"fmt"
	"maps"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

const flag = cre.HTTPTriggerCapability

type HTTPTrigger struct{}

func (o *HTTPTrigger) Flag() cre.CapabilityFlag {
	return flag
}

func (o *HTTPTrigger) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	// add 'http-capabilities' handler to gateway config
	// add gateway connector to to node TOML config, so that node can route http trigger requests to the gateway
	hErr := topology.AddGatewayHandlers(*don, []string{pkg.GatewayHandlerTypeHTTPCapabilities})
	if hErr != nil {
		return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config for don %s ", don.Name)
	}

	// Gateway connector injection scoped to this workflow DON's don_family (see topology_don_family.go).
	cErr := don.ConfigureForGatewayAccess(registryChainID, topology.GatewayConnectorsForDonFamily(don.DonFamily()))
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", don.Name)
	}

	capabilityConfig, ok := don.CapabilityConfigs[flag]
	if !ok {
		return nil, fmt.Errorf("config for %q capability not found for DON %s", flag, don.Name)
	}
	registryConfig, err := BuildCapabilityConfig(capabilityConfig.Values, don.HasOnlyLocalCapabilities())
	if err != nil {
		return nil, err
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "http-trigger",
			Version:        "1.0.0-alpha",
			CapabilityType: 0, // TRIGGER
		},
		Config: registryConfig,
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

// BuildCapabilityConfig supplies service configuration through registry spec_config.
func BuildCapabilityConfig(input map[string]any, localOnly bool) (*capabilitiespb.CapabilityConfig, error) {
	config := make(map[string]any)
	for source, target := range map[string]string{
		"GatewayConnection":            "gatewayConnection",
		"MetadataBatchSize":            "metadataBatchSize",
		"SendChannelBufferSize":        "sendChannelBufferSize",
		"MaxAuthorizedKeysPerWorkflow": "maxAuthorizedKeysPerWorkflow",
		"RequestCacheTTL":              "requestCacheTTL",
	} {
		if v, ok := input[target]; ok {
			config[target] = v
		} else if v, ok := input[source]; ok {
			config[target] = v
		}
	}
	// Metadata timeouts belong inside gatewayConnection in the capability API.
	gateway, _ := config["gatewayConnection"].(map[string]any)
	gateway = maps.Clone(gateway)
	for source, target := range map[string]string{
		"MaxPushMetadataDurationMs": "maxPushMetadataDurationMs",
		"MaxPullMetadataDurationMs": "maxPullMetadataDurationMs",
	} {
		if v, ok := input[source]; ok {
			if gateway == nil {
				gateway = make(map[string]any)
			}
			if _, exists := gateway[target]; !exists {
				gateway[target] = v
			}
		}
	}
	if len(gateway) > 0 {
		config["gatewayConnection"] = gateway
	}
	for _, direction := range []struct{ source, target string }{
		{"Incoming", "incomingRateLimiter"},
		{"Outgoing", "outgoingRateLimiter"},
	} {
		if v, ok := input[direction.target]; ok {
			config[direction.target] = v
			continue
		}
		limiter := make(map[string]any)
		for source, target := range map[string]string{
			"GlobalBurst": "globalBurst", "GlobalRPS": "globalRPS",
			"PerSenderBurst": "perSenderBurst", "PerSenderRPS": "perSenderRPS",
		} {
			if v, ok := input[direction.source+source]; ok {
				limiter[target] = v
			}
		}
		if len(limiter) > 0 {
			config[direction.target] = limiter
		}
	}
	specConfig, err := values.NewMap(config)
	if err != nil {
		return nil, fmt.Errorf("build HTTP trigger spec_config: %w", err)
	}
	return &capabilitiespb.CapabilityConfig{
		LocalOnly:  localOnly,
		SpecConfig: values.ProtoMap(specConfig),
	}, nil
}

func (o *HTTPTrigger) PostEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	_ *cre.Don,
	_ *cre.Dons,
	_ *cre.Environment,
) error {
	// The capability manager launches the service from registry spec_config.
	return nil
}
