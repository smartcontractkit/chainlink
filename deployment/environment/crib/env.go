package crib

const (
	AddressBookFileName   = "ccip-v2-scripts-address-book.json"
	NodesDetailsFileName  = "ccip-v2-scripts-nodes-details.json"
	ChainsConfigsFileName = "ccip-v2-scripts-chains-details.json"
)

type CRIBEnv struct {
	envStateDir string
}

func NewDevspaceEnvFromStateDir(envStateDir string) CRIBEnv {
	return CRIBEnv{
		envStateDir: envStateDir,
	}
}

func (c CRIBEnv) GetConfig() DeployOutput {
	reader := NewOutputReader(c.envStateDir)
	nodesDetails := reader.ReadNodesDetails()
	chainConfigs := reader.ReadChainConfigs()
	return DeployOutput{
		AddressBook: reader.ReadAddressBook(),
		NodeIDs:     nodesDetails.NodeIDs,
		Chains:      chainConfigs,
	}
}

type RPC struct {
	External *string
	Internal *string
}

type ChainConfig struct {
	ChainID   uint64 // chain id as per EIP-155, mainly applicable for EVM chains
	ChainName string // name of the chain populated from chainselector repo
	ChainType string // should denote the chain family. Acceptable values are EVM, COSMOS, SOLANA, STARKNET, APTOS etc
	WSRPCs    []RPC  // websocket rpcs to connect to the chain
	HTTPRPCs  []RPC  // http rpcs to connect to the chain
}

type NodesDetails struct {
	NodeIDs []string
}
