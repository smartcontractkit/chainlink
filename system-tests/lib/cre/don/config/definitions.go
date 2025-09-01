package config

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
)

const (
	// Template for EVM workflow configuration
	evmWorkflowConfigTemplate = `
		[EVM.Workflow]
		FromAddress = '{{.FromAddress}}'
		ForwarderAddress = '{{.ForwarderAddress}}'
		GasLimitDefault = {{.GasLimitDefault}}
		TxAcceptanceState = {{.TxAcceptanceState}}
		PollPeriod = '{{.PollPeriod}}'
		AcceptanceTimeout = '{{.AcceptanceTimeout}}'

		[EVM.Transactions]
		ForwardersEnabled = true
`

	solWorkflowConfigTemplate = `
		Enabled = true
		TxRetentionTimeout = '{{.TxRetentionTimeout}}'

		[Solana.Workflow]
		Enabled = true
		ForwarderAddress = '{{.ForwarderAddress}}'
		FromAddress      = '{{.FromAddress}}'
		ForwarderState   = '{{.ForwarderState}}'
		PollPeriod = '{{.PollPeriod}}'
		AcceptanceTimeout = '{{.AcceptanceTimeout}}'
		TxAcceptanceState = {{.TxAcceptanceState}}
		Local = {{.Local}}
	`
)

func BootstrapEVM(donBootstrapNodePeerID string, homeChainID uint64, capabilitiesRegistryAddress common.Address, chains []*WorkerEVMInput) string {
	evmChainsConfig := ""
	for _, chain := range chains {
		evmChainsConfig += fmt.Sprintf(`
	[[EVM]]
	ChainID = '%d'
	AutoCreateKey = false

	[[EVM.Nodes]]
	Name = '%s'
	WSURL = '%s'
	HTTPURL = '%s'
`,
			chain.ChainID,
			chain.Name,
			chain.WSRPC,
			chain.HTTPRPC,
		)
	}
	return fmt.Sprintf(`
	[Feature]
	LogPoller = true

	[OCR2]
	Enabled = true
	DatabaseTimeout = '1s'
	ContractPollInterval = '1s'

	[P2P.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:5001']
	# bootstrap node in the DON always points to itself as the OCR peering bootstrapper
	DefaultBootstrappers = ['%s@localhost:5001']

%s
	# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
	# Required even, when all capabilities are local to DON in a single DON scenario
	[Capabilities.ExternalRegistry]
	Address = '%s'
	NetworkID = 'evm'
	ChainID = '%d'
`,
		donBootstrapNodePeerID,
		evmChainsConfig,
		capabilitiesRegistryAddress,
		homeChainID,
	)
}

func BoostrapDon2DonPeering(peeringData cre.CapabilitiesPeeringData) string {
	return fmt.Sprintf(`
	[Capabilities.Peering.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:%d']
	DefaultBootstrappers = ['%s@%s:%d']
`,
		peeringData.Port,
		peeringData.GlobalBootstraperPeerID,
		"localhost", // bootstrap node should always point to itself as the don2don peering bootstrapper
		peeringData.Port,
	)
}

type WorkerEVMInput struct {
	Name             string
	ChainID          uint64
	ChainSelector    uint64
	HTTPRPC          string
	WSRPC            string
	FromAddress      common.Address
	ForwarderAddress string
	WritesToEVM      bool
	WorkflowConfig   map[string]any // Configuration for EVM.Workflow section
}

func WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost string, ocrPeeringData cre.OCRPeeringData, capabilitiesPeeringData cre.CapabilitiesPeeringData, capabilitiesRegistryAddress common.Address, homeChainID uint64, chains []*WorkerEVMInput) (string, error) {
	evmChainsConfig := ""
	for _, chain := range chains {
		evmChainsConfig += fmt.Sprintf(`
	[[EVM]]
	ChainID = '%d'
	AutoCreateKey = false
	# reduce workflow registry sync time to minimum to speed up tests & local environment
	FinalityDepth = 1
	LogPollInterval = '5s'

	[[EVM.Nodes]]
	Name = '%s'
	WSURL = '%s'
	HTTPURL = '%s'
`,
			chain.ChainID,
			chain.Name,
			chain.WSRPC,
			chain.HTTPRPC,
		)

		// won't move this to a separate factory function, because this bit needs to be added in the very specific part of the node config
		// it can't be just concatenated to the config in any random place
		if chain.WritesToEVM {
			// Execute template with chain's workflow configuration
			tmpl, err := template.New("evmWorkflowConfig").Parse(evmWorkflowConfigTemplate)
			if err != nil {
				return "", errors.Wrap(err, "failed to parse evm workflow config template")
			}
			var configBuffer bytes.Buffer
			if executeErr := tmpl.Execute(&configBuffer, chain.WorkflowConfig); executeErr != nil {
				return "", errors.Wrap(executeErr, "failed to execute evm workflow config template")
			}

			flag := cre.WriteEVMCapability
			configStr := configBuffer.String()

			if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
				return "", errors.Wrapf(err, "%s template validation failed", flag)
			}

			evmChainsConfig += configStr
		}
	}

	return fmt.Sprintf(`
	[Feature]
	LogPoller = true

	[OCR2]
	Enabled = true
	DatabaseTimeout = '1s'
	ContractPollInterval = '1s'

	[P2P.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:%d']
	DefaultBootstrappers = ['%s@%s:%d']

	[Capabilities.Peering.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:%d']
	DefaultBootstrappers = ['%s@%s:%d']

%s
	# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
	# Required even, when all capabilities are local to DON in a single DON scenario
	[Capabilities.ExternalRegistry]
	Address = '%s'
	NetworkID = 'evm'
	ChainID = '%d'
`,
		ocrPeeringData.Port,
		donBootstrapNodePeerID,
		donBootstrapNodeHost,
		ocrPeeringData.Port,
		capabilitiesPeeringData.Port,
		capabilitiesPeeringData.GlobalBootstraperPeerID,
		capabilitiesPeeringData.GlobalBootstraperHost,
		capabilitiesPeeringData.Port,
		evmChainsConfig,
		capabilitiesRegistryAddress,
		homeChainID,
	), nil
}

type WorkerSolanaInput struct {
	Name             string
	ChainID          string
	ChainSelector    uint64
	NodeURL          string
	FromAddress      solana.PublicKey
	ForwarderAddress string
	ForwarderState   string
	HasWrite         bool
	WorkflowConfig   map[string]any // Configuration for Solana.Workflow section
}

func BootstrapSolana(chains []*WorkerSolanaInput) string {
	var ret string
	for _, chain := range chains {
		ret += fmt.Sprintf(`
		[[Solana]]
		ChainID = '%s'
		Enabled = true

		[[Solana.Nodes]]
		Name = '%s'
		URL = '%s'
		`, chain.ChainID, chain.Name, chain.NodeURL)
	}

	return ret
}

func WorkerSolana(chains []*WorkerSolanaInput) (string, error) {
	var ret string
	for _, chain := range chains {
		ret += fmt.Sprintf(`
		[[Solana]]
		ChainID = '%s'
		`, chain.ChainID)

		// won't move this to a separate factory function, because this bit needs to be added in the very specific part of the node config
		// it can't be just concatenated to the config in any random place
		if chain.HasWrite {
			// Execute template with chain's workflow configuration
			tmpl, err := template.New("solanaWorkflowConfig").Parse(solWorkflowConfigTemplate)
			if err != nil {
				return "", errors.Wrap(err, "failed to parse solana workflow config template")
			}
			var configBuffer bytes.Buffer
			if executeErr := tmpl.Execute(&configBuffer, chain.WorkflowConfig); executeErr != nil {
				return "", errors.Wrap(executeErr, "failed to execute solana workflow config template")
			}

			flag := cre.WriteSolanaCapability
			configStr := configBuffer.String()

			if err := don.ValidateTemplateSubstitution(configStr, flag); err != nil {
				return "", errors.Wrapf(err, "%s template validation failed", flag)
			}

			ret += configStr
		}

		ret += fmt.Sprintf(`
		[[Solana.Nodes]]
		Name = '%s'
		URL = '%s'
		`, chain.Name, chain.NodeURL)
	}

	return ret, nil
}

func WorkerWorkflowRegistry(workflowRegistryAddr common.Address, homeChainID uint64) string {
	return fmt.Sprintf(`
	# there are two strategies for syncing workflow registry:
	# - reconciliation: poll the contract for events
	# - event: watch events on the contract
	[Capabilities.WorkflowRegistry]
	Address = "%s"
	NetworkID = "evm"
	ChainID = "%d"
	# SyncStrategy = "reconciliation"
`,
		workflowRegistryAddr.Hex(),
		homeChainID,
	)
}

func WorkerGateway(nodeAddress common.Address, homeChainID uint64, donID string, gatewayConfiguration []*cre.GatewayConfiguration) string {
	config := fmt.Sprintf(`
	[Capabilities.GatewayConnector]
	DonID = "%s"
	ChainIDForNodeKey = "%d"
	NodeAddress = '%s'
`,
		donID,
		homeChainID,
		nodeAddress,
	)

	for _, gatewayConnectorData := range gatewayConfiguration {
		gatewayURL := fmt.Sprintf("ws://%s:%d%s", gatewayConnectorData.Outgoing.Host, gatewayConnectorData.Outgoing.Port, gatewayConnectorData.Outgoing.Path)
		config += fmt.Sprintf(`
	[[Capabilities.GatewayConnector.Gateways]]
	Id = "%s"
	URL = "%s"
`,
			gatewayConnectorData.AuthGatewayID,
			gatewayURL,
		)
	}

	return config
}
