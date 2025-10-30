package gateway

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"dario.cat/mergo"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

type WhitelistConfig struct {
	ExtraAllowedPorts                    []int
	ExtraAllowedIPs, ExtraAllowedIPsCIDR []string
}

func Configs(topology *cre.Topology) ([]cre.GatewayConfig, error) {
	workflowDON, donErr := topology.DonsMetadata.WorkflowDON()
	if donErr != nil {
		return nil, errors.Wrap(donErr, "failed to find workflow DON")
	}

	return []cre.GatewayConfig{{
		Name:     workflowDON.Name,
		Handlers: []string{pkg.GatewayHandlerTypeWebAPICapabilities},
	}}, nil
}

func CreateJobs(ctx context.Context, creEnv *cre.Environment, dons *cre.Dons, gatewayConfigs []cre.GatewayConfig, whitelistConfig WhitelistConfig) error {
	specs := make(map[string][]string)

	for _, config := range dons.GatewayConnectors.Configurations {
		gatewayNode, ok := dons.NodeWithUUID(config.NodeUUID)
		if !ok {
			return fmt.Errorf("could not find gateway node with UUID %s in DON topology", config.NodeUUID)
		}

		workerInput := cre_jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     gatewayNode.DON.Name,
			JobName:     "gateway-worker",
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: gatewayNode.DON.Name},
				{Key: "p2p_id", Value: gatewayNode.Keys.PeerID()},
			},
			Template: job_types.Gateway,
			Inputs: job_types.JobSpecInput{
				"dons":                    gatewayConfigs,
				"allowedPorts":            append(whitelistConfig.ExtraAllowedPorts, DefaultAllowedPorts...),
				"allowedSchemes":          []string{"http", "https"},
				"allowedIPsCIDR":          whitelistConfig.ExtraAllowedIPsCIDR,
				"gatewayKeyChainSelector": creEnv.RegistryChainSelector,
				"authGatewayID":           config.AuthGatewayID,
			},
		}

		workerVerErr := cre_jobs.ProposeJobSpec{}.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput)
		if workerVerErr != nil {
			return fmt.Errorf("precondition verification failed for Custom Compute worker job: %w", workerVerErr)
		}

		workerReport, workerErr := cre_jobs.ProposeJobSpec{}.Apply(*creEnv.CldfEnvironment, workerInput)
		if workerErr != nil {
			return fmt.Errorf("failed to propose Custom Compute worker job spec: %w", workerErr)
		}

		for _, r := range workerReport.Reports {
			out, ok := r.Output.(cre_jobs_ops.ProposeGatewayJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeGatewayJobOutput, actual type: %T", r.Output)
			}
			mErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice)
			if mErr != nil {
				return fmt.Errorf("failed to merge worker job specs: %w", mErr)
			}
		}
	}

	approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs)
	if approveErr != nil {
		return fmt.Errorf("failed to approve Custom Compute jobs: %w", approveErr)
	}

	return nil
}

// AddHandlers adds the given handler name to the gateway config of the given DON. It only adds handlers, if they are not already present.
func AddHandlers(donMetadata cre.DonMetadata, gatewayConfigs []cre.GatewayConfig, handlers []string) ([]cre.GatewayConfig, error) {
	donFound := false

	for idx, gc := range gatewayConfigs {
		if gc.Name == donMetadata.Name {
			donFound = true
		}

		if donFound {
			for _, handlerName := range handlers {
				alreadyPresent := false
				for _, existingHandler := range gc.Handlers {
					if strings.EqualFold(existingHandler, handlerName) {
						alreadyPresent = true
						break
					}
				}
				if !alreadyPresent {
					gatewayConfigs[idx].Handlers = append(gatewayConfigs[idx].Handlers, handlerName)
				}
			}
			break
		}
	}

	// if we did not find the DON in the gateway config, we need to add it
	if !donFound {
		gatewayConfigs = append(gatewayConfigs, cre.GatewayConfig{
			Name:     donMetadata.Name,
			Handlers: handlers,
		})
	}

	return gatewayConfigs, nil
}

// AddConnectors adds gateway connector configuration to the node TOML config of each node in the given DON. It only adds connectors, if they are not already present.
func AddConnectors(donMetadata *cre.DonMetadata, registryChainID uint64, connectors cre.GatewayConnectors) error {
	workers, wErr := donMetadata.Workers()
	if wErr != nil {
		return wErr
	}

	for _, workerNode := range workers {
		currentConfig := donMetadata.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
		}

		evmKey, ok := workerNode.Keys.EVM[registryChainID]
		if !ok {
			return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", registryChainID, workerNode.Index)
		}

		// if no gateways are configured, then gateway connector config is most probably also not configured
		if len(typedConfig.Capabilities.GatewayConnector.Gateways) == 0 {
			typedConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
				DonID:             ptr.Ptr(donMetadata.Name),
				ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(registryChainID, 10)),
				NodeAddress:       ptr.Ptr(evmKey.PublicAddress.Hex()),
			}
		}

		// make sure that all other gateways are also present in the config
		for _, gatewayConnector := range connectors.Configurations {
			alreadyPresent := false
			for _, existingGateway := range typedConfig.Capabilities.GatewayConnector.Gateways {
				if gatewayConnector.AuthGatewayID == *existingGateway.ID {
					alreadyPresent = true
					continue
				}
			}

			if !alreadyPresent {
				typedConfig.Capabilities.GatewayConnector.Gateways = append(typedConfig.Capabilities.GatewayConnector.Gateways, coretoml.ConnectorGateway{
					ID: ptr.Ptr(gatewayConnector.AuthGatewayID),
					URL: ptr.Ptr(fmt.Sprintf("ws://%s:%d%s",
						gatewayConnector.Outgoing.Host,
						gatewayConnector.Outgoing.Port,
						gatewayConnector.Outgoing.Path)),
				})
			}
		}

		stringifiedConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
		}

		donMetadata.NodeSets().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
	}

	return nil
}
