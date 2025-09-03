package config

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"slices"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

type bootstrapConfigTemplateData struct {
	OCRPeeringData              cre.OCRPeeringData
	CapabilitiesPeeringData     cre.CapabilitiesPeeringData
	RegistryChainID             uint64
	CapabilitiesRegistryAddress string
	EVMChains                   []*evmChain
}

const bootstrapConfigTemplate = `
[Feature]
LogPoller = true

[OCR2]
Enabled = true
DatabaseTimeout = '1s'
ContractPollInterval = '1s'

[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:{{.OCRPeeringData.Port}}']
# bootstrap node in the DON always points to itself as the OCR peering bootstrapper
DefaultBootstrappers = ['{{.OCRPeeringData.OCRBootstraperPeerID}}@localhost:{{.OCRPeeringData.Port}}']

# support multiple chains, even though bootstrap node only needs to be connected to the registry chain
{{range .EVMChains}}
[[EVM]]
ChainID = '{{.ChainID}}'
AutoCreateKey = false

[[EVM.Nodes]]
Name = '{{.Name}}'
WSURL = '{{.WSRPC}}'
HTTPURL = '{{.HTTPRPC}}'
{{end}}

# we assume that this bootstrap node is also the capabilities peering bootstrapper
[Capabilities.Peering.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:{{.CapabilitiesPeeringData.Port}}']
DefaultBootstrappers = ['{{.CapabilitiesPeeringData.GlobalBootstraperPeerID}}@localhost:{{.CapabilitiesPeeringData.Port}}']

# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
# Required even, when all capabilities are local to DON in a single DON scenario
[Capabilities.ExternalRegistry]
Address = '{{.CapabilitiesRegistryAddress}}'
NetworkID = 'evm'
ChainID = '{{.RegistryChainID}}'
`

// GatewayData represents a single gateway configuration
type gatewayData struct {
	ID  string
	URL string
}

type workerConfigTemplateData struct {
	OCRPeeringData          cre.OCRPeeringData
	CapabilitiesPeeringData cre.CapabilitiesPeeringData
	RegistryChainID         uint64

	NodeETHAddress              string
	CapabilitiesRegistryAddress string
	WorkflowRegistryAddress     string

	EVMChains   []*evmChain
	SolanaChain *solanaChain

	DonID    string
	Gateways []gatewayData

	// Template control flags
	IncludeSolanaChain      bool
	IncludeWorkflowRegistry bool // required by workers in the workflow DON only
	IncludeGatewayConnector bool // required by workers in the workflow DON and if DON has certain capabilities
}

const workerConfigTemplate = `
[Feature]
LogPoller = true

[OCR2]
Enabled = true
DatabaseTimeout = '1s'
ContractPollInterval = '1s'

[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:{{.OCRPeeringData.Port}}']
DefaultBootstrappers = ['{{.OCRPeeringData.OCRBootstraperPeerID}}@{{.OCRPeeringData.OCRBootstraperHost}}:{{.OCRPeeringData.Port}}']

[Capabilities.Peering.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:{{.CapabilitiesPeeringData.Port}}']
DefaultBootstrappers = ['{{.CapabilitiesPeeringData.GlobalBootstraperPeerID}}@{{.CapabilitiesPeeringData.GlobalBootstraperHost}}:{{.CapabilitiesPeeringData.Port}}']

{{range .EVMChains}}
[[EVM]]
ChainID = '{{.ChainID}}'
AutoCreateKey = false
# reduce workflow registry sync time to minimum to speed up tests & local environment
FinalityDepth = 1
LogPollInterval = '5s'

[[EVM.Nodes]]
Name = '{{.Name}}'
WSURL = '{{.WSRPC}}'
HTTPURL = '{{.HTTPRPC}}'
{{end}}

{{if .IncludeSolanaChain}}
[[Solana]]
ChainID = '{{.SolanaChain.ChainID}}'
Enabled = true

[[Solana.Nodes]]
Name = '{{.SolanaChain.Name}}'
URL = '{{.SolanaChain.NodeURL}}'
{{end}}

# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
# Required even, when all capabilities are local to DON in a single DON scenario
[Capabilities.ExternalRegistry]
Address = '{{.CapabilitiesRegistryAddress}}'
NetworkID = 'evm'
ChainID = '{{.RegistryChainID}}'

{{if .IncludeWorkflowRegistry}}
[Capabilities.WorkflowRegistry]
Address = "{{.WorkflowRegistryAddress}}"
NetworkID = "evm"
ChainID = "{{.RegistryChainID}}"
# there are two strategies for syncing workflow registry:
# - reconciliation: poll the contract for events
# - event: watch events on the contract
# SyncStrategy = "reconciliation"
{{end}}

{{ if.IncludeGatewayConnector}}
[Capabilities.GatewayConnector]
DonID = "{{.DonID}}"
ChainIDForNodeKey = "{{.RegistryChainID}}"
NodeAddress = "{{.NodeETHAddress}}"
{{range .Gateways}}
[[Capabilities.GatewayConnector.Gateways]]
Id = "{{.ID}}"
URL = "{{.URL}}"
{{end}}
{{end}}
`

func generateBootstrapNodeConfig(
	ocrPeeringData cre.OCRPeeringData,
	capabilitiesPeeringData cre.CapabilitiesPeeringData,
	commonInputs *commonInputs,
) (string, error) {
	tmpl, err := template.New("bootstrapConfigTemplate").Parse(bootstrapConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse bootstrap node config template: %v", err)
	}

	data := bootstrapConfigTemplateData{
		OCRPeeringData:              ocrPeeringData,
		CapabilitiesPeeringData:     capabilitiesPeeringData,
		RegistryChainID:             commonInputs.registryChainID,
		CapabilitiesRegistryAddress: commonInputs.capabilitiesRegistryAddress.Hex(),
		EVMChains:                   commonInputs.evmChains,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(fmt.Sprintf("failed to execute bootstrap EVM template: %v", err))
	}

	rendered := buf.String()
	if err := don.ValidateTemplateSubstitution(rendered, "bootstrapConfigTemplate"); err != nil {
		return "", fmt.Errorf("failed to validate template substitution")
	}

	return rendered, nil
}

func generateWorkerNodeConfig(
	ocrPeeringData cre.OCRPeeringData,
	capabilitiesPeeringData cre.CapabilitiesPeeringData,
	commonInputs *commonInputs,
	gatewayConfigurations []*cre.GatewayConfiguration,
	donName string,
	donFlags []string,
	nodeLabels []*cre.Label,
) (string, error) {
	tmpl, err := template.New("workerNodeConfig").Parse(workerConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse worker config template: %w", err)
	}

	//find node's ETH address on the registry chain
	var nodeEthAddr string
	expectedAddressKey := node.AddressKeyFromSelector(commonInputs.registryChainSelector)
	for _, label := range nodeLabels {
		if label.Key == expectedAddressKey {
			nodeEthAddr = label.Value
			break
		}
	}

	if nodeEthAddr == "" {
		return "", errors.Errorf("no ETH address found for node for chain %d", commonInputs.registryChainID)
	}

	// convert Gateway configurations to structure expected by the template
	gateways := make([]gatewayData, len(gatewayConfigurations))
	for i, gatewayConnectorData := range gatewayConfigurations {
		gatewayURL := fmt.Sprintf("ws://%s:%d%s",
			gatewayConnectorData.Outgoing.Host,
			gatewayConnectorData.Outgoing.Port,
			gatewayConnectorData.Outgoing.Path)

		gateways[i] = gatewayData{
			ID:  gatewayConnectorData.AuthGatewayID,
			URL: gatewayURL,
		}
	}

	data := workerConfigTemplateData{
		OCRPeeringData:              ocrPeeringData,
		CapabilitiesPeeringData:     capabilitiesPeeringData,
		CapabilitiesRegistryAddress: commonInputs.capabilitiesRegistryAddress.Hex(),
		WorkflowRegistryAddress:     commonInputs.workflowRegistryAddress.Hex(),
		NodeETHAddress:              nodeEthAddr,
		DonID:                       donName,
		RegistryChainID:             commonInputs.registryChainID,
		EVMChains:                   commonInputs.evmChains,
		SolanaChain:                 commonInputs.solanaChain,
		Gateways:                    gateways,
		IncludeSolanaChain:          commonInputs.solanaChain != nil,
		IncludeWorkflowRegistry:     flags.HasFlag(donFlags, cre.WorkflowDON), // include Workflow Registry configuration only for worker nodes in the Workflow DON
		IncludeGatewayConnector:     flags.HasFlag(donFlags, cre.WorkflowDON) || don.NodeNeedsAnyGateway(donFlags),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute worker config template: %w", err)
	}

	rendered := buf.String()
	if err := don.ValidateTemplateSubstitution(rendered, "workerConfigTemplate"); err != nil {
		return "", fmt.Errorf("failed to validate template substitution")
	}

	return rendered, nil
}

func Generate(input cre.GenerateConfigsInput, nodeConfigFns []cre.NodeConfigTransformerFn) (cre.NodeIndexToConfigOverride, error) {
	configOverrides := make(cre.NodeIndexToConfigOverride)

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	commonInputs, inputsErr := gatherCommonInputs(input)
	if inputsErr != nil {
		return nil, errors.Wrap(inputsErr, "failed to gather common inputs")
	}

	for nodeIdx, nodeMetadata := range input.DonMetadata.NodesMetadata {
		nodeType, typeErr := node.FindLabelValue(nodeMetadata, node.NodeTypeKey)
		if typeErr != nil {
			return nil, errors.Wrap(typeErr, "failed to find node type")
		}

		var nodeConfig string
		var configErr error
		switch nodeType {
		case cre.BootstrapNode:
			nodeConfig, configErr = generateBootstrapNodeConfig(input.OCRPeeringData, input.CapabilitiesPeeringData, commonInputs)
		case cre.WorkerNode:
			nodeConfig, configErr = generateWorkerNodeConfig(input.OCRPeeringData, input.CapabilitiesPeeringData, commonInputs, input.GatewayConnectorOutput.Configurations, input.DonMetadata.Name, input.DonMetadata.Flags, nodeMetadata.Labels)
		default:
			return nil, fmt.Errorf("unsupported node type found for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
		}

		if configErr != nil {
			return nil, errors.Wrapf(configErr, "failed to generate node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
		}

		configOverrides[nodeIdx] = nodeConfig
	}

	// execute capability-provided functions that transform the node config (currently: write-evm, write-solana)
	// these functions must return whole node configs after transforming them, instead of just returning configuration parts
	// that need to be merged into the existing config
	for _, configFn := range nodeConfigFns {
		if configFn == nil {
			continue
		}

		modifiedConfigs, err := configFn(input, configOverrides)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}

		maps.Copy(configOverrides, modifiedConfigs)
	}

	return configOverrides, nil
}

type commonInputs struct {
	registryChainID             uint64
	registryChainSelector       uint64
	capabilitiesRegistryAddress common.Address
	workflowRegistryAddress     common.Address

	evmChains   []*evmChain
	solanaChain *solanaChain
}

func gatherCommonInputs(input cre.GenerateConfigsInput) (*commonInputs, error) {
	registryChainID, homeErr := chain_selectors.ChainIdFromSelector(input.HomeChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	// prepare chains, we need chainIDs and URLs
	evmChains := findEVMChains(input)
	solanaChain, solErr := findOneSolanaChain(input)
	if solErr != nil {
		return nil, errors.Wrap(solErr, "failed to find Solana chain in the environment configuration")
	}

	// find contract addresses
	capabilitiesRegistryAddress, capErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capErr != nil {
		return nil, errors.Wrap(capErr, "failed to find CapabilitiesRegistry address")
	}

	workflowRegistryAddress, capErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.WorkflowRegistry.String())
	if capErr != nil {
		return nil, errors.Wrap(capErr, "failed to find WorkflowRegistry address")
	}

	return &commonInputs{
		registryChainID:             registryChainID,
		registryChainSelector:       input.HomeChainSelector,
		capabilitiesRegistryAddress: capabilitiesRegistryAddress,
		workflowRegistryAddress:     workflowRegistryAddress,
		evmChains:                   evmChains,
		solanaChain:                 solanaChain,
	}, nil
}

type evmChain struct {
	Name    string
	ChainID uint64
	HTTPRPC string
	WSRPC   string
}

func findEVMChains(input cre.GenerateConfigsInput) []*evmChain {
	evmChains := make([]*evmChain, 0)
	for chainSelector, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain != nil {
			continue
		}

		// if the DON doesn't support the chain, we skip it; if slice is empty, it means that the DON supports all chains
		// TODO: review if we really need this SupportedChains functionality
		if len(input.DonMetadata.SupportedChains) > 0 && !slices.Contains(input.DonMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		evmChains = append(evmChains, &evmChain{
			Name:    fmt.Sprintf("node-%d", chainSelector),
			ChainID: bcOut.ChainID,
			HTTPRPC: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
			WSRPC:   bcOut.BlockchainOutput.Nodes[0].InternalWSUrl,
		})
	}
	return evmChains
}

type solanaChain struct {
	Name    string
	ChainID string
	NodeURL string
}

func findOneSolanaChain(input cre.GenerateConfigsInput) (*solanaChain, error) {
	var solChain *solanaChain
	chainsFound := 0

	for _, bcOut := range input.BlockchainOutput {
		if bcOut.SolChain == nil {
			continue
		}

		chainsFound++
		if chainsFound > 1 {
			return nil, errors.New("multiple Solana chains found, expected only one")
		}

		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		chainID, err := bcOut.SolClient.GetGenesisHash(ctx)
		if err != nil {
			cancelFn()
			return nil, errors.Wrap(err, "failed to get chainID for Solana")
		}
		cancelFn()

		solChain = &solanaChain{
			Name:    fmt.Sprintf("node-%d", bcOut.SolChain.ChainSelector),
			ChainID: chainID.String(),
			NodeURL: bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
		}
	}

	return solChain, nil
}
