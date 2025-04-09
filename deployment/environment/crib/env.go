package crib

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	AddressBookFileName       = "address-book.json"
	NodesDetailsFileName      = "nodes-details.json"
	ChainsConfigsFileName     = "chains-details.json"
	RMNNodeIdentitiesFileName = "rmn-node-identities.json"
)

type CRIBEnv struct {
	lggr                logger.Logger
	cribEnvStateDirPath string
}

func NewDevspaceEnvFromStateDir(lggr logger.Logger, envStateDir string) CRIBEnv {
	return CRIBEnv{
		lggr:                lggr,
		cribEnvStateDirPath: envStateDir,
	}
}

func (c CRIBEnv) GetConfig(key string) (DeployOutput, error) {
	reader := NewOutputReader(c.cribEnvStateDirPath)
	nodesDetails, err := reader.ReadNodesDetails()
	if err != nil {
		c.lggr.Warn("No nodes details found, not necessary for testing.. continuing...", err)
	}
	chainConfigs, err := reader.ReadChainConfigs()
	if err != nil {
		return DeployOutput{}, err
	}
	for i, chain := range chainConfigs {
		err := chain.SetDeployerKey(&key)
		if err != nil {
			return DeployOutput{}, err
		}
		chainConfigs[i] = chain
	}

	addressBook, err := reader.ReadAddressBook()
	if err != nil {
		return DeployOutput{}, err
	}

	return DeployOutput{
		AddressBook: addressBook,
		NodeIDs:     nodesDetails.NodeIDs,
		Chains:      chainConfigs,
	}, nil
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

type BootstrapNode struct {
	P2PID        string
	InternalHost string
	Port         string
}

type NodesDetails struct {
	NodeIDs       []string
	BootstrapNode BootstrapNode
}
