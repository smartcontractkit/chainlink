package types

import "github.com/smartcontractkit/chainlink/deployment"

// data streams contract types
const (
	ChannelConfigStore deployment.ContractType = "ChannelConfigStore"
	Configurator       deployment.ContractType = "Configurator"
	FeeManager         deployment.ContractType = "FeeManager"
	RewardManager      deployment.ContractType = "RewardManager"
	Verifier           deployment.ContractType = "Verifier"
	VerifierProxy      deployment.ContractType = "VerifierProxy"
)
