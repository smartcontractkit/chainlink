package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func BootstrapEVM(donBootstrapNodePeerID string, homeChainID uint64, capabilitiesRegistryAddress common.Address, chains []*EVMChain) string {
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

type EVMChain struct {
	Name    string
	ChainID uint64
	HTTPRPC string
	WSRPC   string
}

func WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost string, ocrPeeringData cre.OCRPeeringData, capabilitiesPeeringData cre.CapabilitiesPeeringData, capabilitiesRegistryAddress common.Address, homeChainID uint64, chains []*EVMChain) (string, error) {
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

type SolanaChain struct {
	Name    string
	ChainID string
	NodeURL string
}

func WorkerSolana(chain *SolanaChain) string {
	if chain == nil {
		return ""
	}

	return fmt.Sprintf(`
		[[Solana]]
		ChainID = '%s'
		Enabled = true

		[[Solana.Nodes]]
		Name = '%s'
		URL = '%s'
		`, chain.ChainID, chain.Name, chain.NodeURL)
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
