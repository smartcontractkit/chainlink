package config

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func BootstrapEVM(donBootstrapNodePeerID string, homeChainID uint64, capabilitiesRegistryAddress common.Address, chains []*WorkerEVMInput) string {
	evmChainsConfig := ""
	for _, chain := range chains {
		evmChainsConfig += fmt.Sprintf(`
	[[EVM]]
	ChainID = '%s'
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

func BoostrapDon2DonPeering(peeringData types.CapabilitiesPeeringData) string {
	return fmt.Sprintf(`
	[Capabilities.Peering.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:6690']
	DefaultBootstrappers = ['%s@%s:6690']
`,
		peeringData.GlobalBootstraperPeerID,
		"localhost", // bootstrap node should always point to itself as the don2don peering bootstrapper
	)
}

type WorkerEVMInput struct {
	Name             string
	ChainID          string
	ChainSelector    uint64
	HTTPRPC          string
	WSRPC            string
	FromAddress      common.Address
	ForwarderAddress string
}

func WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost string, peeringData types.CapabilitiesPeeringData, capabilitiesRegistryAddress common.Address, homeChainID uint64, chains []*WorkerEVMInput) string {
	evmChainsConfig := ""
	for _, chain := range chains {
		evmChainsConfig += fmt.Sprintf(`
	[[EVM]]
	ChainID = '%s'
	AutoCreateKey = false
	# reduce workflow registry sync time to minimum to speed up tests & local environment
	FinalityDepth = 1
	LogPollInterval = '5s'

	[[EVM.Nodes]]
	Name = '%s'
	WSURL = '%s'
	HTTPURL = '%s'

	[EVM.Workflow]
	FromAddress = '%s'
	ForwarderAddress = '%s'
	GasLimitDefault = 400_000

	[EVM.Transactions]
	ForwardersEnabled = true

`,
			chain.ChainID,
			chain.Name,
			chain.WSRPC,
			chain.HTTPRPC,
			chain.FromAddress,
			chain.ForwarderAddress,
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
	DefaultBootstrappers = ['%s@%s:5001']

	[Capabilities.Peering.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:6690']
	DefaultBootstrappers = ['%s@%s:6690']

%s
	# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
	# Required even, when all capabilities are local to DON in a single DON scenario
	[Capabilities.ExternalRegistry]
	Address = '%s'
	NetworkID = 'evm'
	ChainID = '%d'
`,
		donBootstrapNodePeerID,
		donBootstrapNodeHost,
		peeringData.GlobalBootstraperPeerID,
		peeringData.GlobalBootstraperHost,
		evmChainsConfig,
		capabilitiesRegistryAddress,
		homeChainID,
	)
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

func WorkerGateway(nodeAddress common.Address, homeChainID uint64, donID uint32, gatewayConnectorData types.GatewayConnectorOutput) string {
	gatewayURL := fmt.Sprintf("ws://%s:%d/%s", gatewayConnectorData.Host, 5003, "node")

	return fmt.Sprintf(`
	[Capabilities.GatewayConnector]
	DonID = "%s"
	ChainIDForNodeKey = "%d"
	NodeAddress = '%s'

	[[Capabilities.GatewayConnector.Gateways]]
	Id = "por_gateway"
	URL = "%s"
`,
		strconv.FormatUint(uint64(donID), 10),
		homeChainID,
		nodeAddress,
		gatewayURL,
	)
}
