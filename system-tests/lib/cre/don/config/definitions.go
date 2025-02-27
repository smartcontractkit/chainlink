package config

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func BootstrapEVM(donBootstrapNodePeerID string, chainID uint64, capabilitiesRegistryAddress common.Address, httpRPC, wsRPC string) string {
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

	[[EVM]]
	ChainID = '%d'

	[[EVM.Nodes]]
	Name = 'anvil'
	WSURL = '%s'
	HTTPURL = '%s'

	# Capabilities registry address, required for do2don p2p mesh to work and for capabilities discovery
	# Required even, when all capabilities are local to DON in a single DON scenario
	[Capabilities.ExternalRegistry]
	Address = '%s'
	NetworkID = 'evm'
	ChainID = '%d'
`,
		donBootstrapNodePeerID,
		chainID,
		wsRPC,
		httpRPC,
		capabilitiesRegistryAddress,
		chainID,
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

// could add multichain with something like this:
//
//	type EVMChain struct {
//		ChainID uint64
//		HTTPRPC string
//		WSRPC   string
//	}
//
// so that we are future-proof (for bootstrap too!)
// we'd need to have capabilitiesRegistryChainID too
func WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost string, peeringData types.CapabilitiesPeeringData, chainID uint64, capabilitiesRegistryAddress common.Address, httpRPC, wsRPC string) string {
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

	[[EVM]]
	ChainID = '%d'
	AutoCreateKey = false

	[[EVM.Nodes]]
	Name = 'anvil'
	WSURL = '%s'
	HTTPURL = '%s'

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
		chainID,
		wsRPC,
		httpRPC,
		capabilitiesRegistryAddress,
		chainID,
	)
}

func WorkerWriteEMV(nodeAddress, forwarderAddress common.Address) string {
	return fmt.Sprintf(`
	# Required for the target capability to be initialized
	[EVM.Workflow]
	FromAddress = '%s'
	ForwarderAddress = '%s'
	GasLimitDefault = 400_000
`,
		nodeAddress.Hex(),
		forwarderAddress.Hex(),
	)
}

func WorkerWorkflowRegistry(workflowRegistryAddr common.Address, chainID uint64) string {
	return fmt.Sprintf(`
	[Capabilities.WorkflowRegistry]
	Address = "%s"
	NetworkID = "evm"
	ChainID = "%d"
`,
		workflowRegistryAddr.Hex(),
		chainID,
	)
}

func WorkerGateway(nodeAddress common.Address, chainID uint64, donID uint32, gatewayConnectorData types.GatewayConnectorOutput) string {
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
		chainID,
		nodeAddress,
		gatewayURL,
	)
}
