package deployment_common

import (
	"github.com/smartcontractkit/chainlink/deployment"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
)

// This changeset uses the DeployLinkSequence to deploy a LINK token contract, grant minting permissions, mint some LINK to the deployer, and transfer some LINK to another address
var LinkExampleChangeset = func(e deployment.Environment, config SqDeployLinkInput) (deployment.ChangesetOutput, error) {

	// Prepare sequence deps
	auth := e.Chains[config.chainID].DeployerKey
	client := e.Chains[config.chainID].Client
	deps := deployment_ethereum.EthereumDeps{
		Auth:    auth,
		Client:  client,
		Confirm: e.Chains[config.chainID].ConfirmByHash,
	}

	// Execute Sequence
	// TODO: Build a way to re execute the sequence from an existing report. Ideally this is an extra optional parameter in the changeset
	linkDeployReport := deployment.ExecuteSeq(e.OpEnv, *DeployLinkSequence, deps, config)
	if linkDeployReport.Err != nil {
		return deployment.ChangesetOutput{}, linkDeployReport.Err
	}

	// TODO: Changeset should return its own Report with a unique ID, storing low level operation reports
	return deployment.ChangesetOutput{
		// TODO: Make a way to return only the executed reports
		Reports: e.OpEnv.Reporter.GetReports(),
	}, nil
}
