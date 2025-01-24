package deployment_common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
)

type ChangesetLinkInput struct {
	MintAmount *big.Int
	Amount     *big.Int
	To         common.Address
	chainId    uint64
}

// This changeset deploys and transfers an specific amount of LINK to an address
var LinkExampleChangeset = func(e deployment.Environment, config ChangesetLinkInput) (deployment.ChangesetOutput, error) {

	// TODO: Previous reports should come from the user. If present, would allow retryability
	reporter := deployment.NewMemoryReporter([]deployment.Report[any, any, any]{})

	// Prepare operation context
	auth := e.Chains[config.chainId].DeployerKey
	client := e.Chains[config.chainId].Client
	ethCtx := deployment.Context[deployment_ethereum.EthereumDeps]{
		Log: e.Logger,
		Deps: deployment_ethereum.EthereumDeps{
			Auth:    auth,
			Client:  client,
			Confirm: e.Chains[config.chainId].Confirm,
		},
	}

	linkDeployReport, err := deployment.Execute(reporter, DeployLinkOp, ethCtx, deployment.EmptyInput{})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	grantMintInput := GrantLinkInput{
		contractAddress: linkDeployReport.Output.Address,
		To:              auth.From,
	}

	_, err = deployment.Execute(reporter, GrantMintLinkOp, ethCtx, grantMintInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	mintInput := MintLinkInput{
		contractAddress: linkDeployReport.Output.Address,
		To:              auth.From,
		Amount:          config.MintAmount,
	}

	_, err = deployment.Execute(reporter, MintLinkOp, ethCtx, mintInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	transferInput := TransferLinkInput{
		contractAddress: linkDeployReport.Output.Address,
		To:              config.To,
		Amount:          config.Amount,
	}

	_, err = deployment.Execute(reporter, TransferLinkOp, ethCtx, transferInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// TODO: Changeset should return its own Report with a unique ID, storing low level operation reports
	return deployment.ChangesetOutput{
		// Should include Address Book and other relevant information
		Reports: reporter.GetReports(),
	}, err
}
