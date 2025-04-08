package v1_5_1

import (
	"github.com/smartcontractkit/chainlink/deployment"
)

var DeployTokenPoolFactory = deployment.CreateChangeSet(deployTokenPoolFactoryLogic, deployTokenPoolFactoryPrecondition)

type DeployTokenPoolFactoryConfig struct {
	// Chains is the list of chains on which to deploy the token pool factory.
	Chains []uint64
}

func deployTokenPoolFactoryPrecondition(e deployment.Environment, config DeployTokenPoolFactoryConfig) error {
	return nil
}

func deployTokenPoolFactoryLogic(e deployment.Environment, config DeployTokenPoolFactoryConfig) (deployment.ChangesetOutput, error) {
	return deployment.ChangesetOutput{}, nil
}
