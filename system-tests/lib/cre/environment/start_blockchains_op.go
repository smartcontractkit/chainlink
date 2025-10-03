package environment

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
)

type StartBlockchainsOpDeps struct {
	Deployers map[blockchains.ChainFamily]blockchains.Deployer
}

type StartBlockchainsOpInput struct {
	Inputs []*blockchain.Input
}

type StartBlockchainsOpOutput struct {
	DeployedBlockchains *blockchains.DeployedBlockchains `json:"-"`
	OpBundle            operations.Bundle                `json:"-"`
}

var StartBlockchainsOp = operations.NewOperation(
	"start-blockchains-op",
	semver.MustParse("1.0.0"),
	"Starts blockchains using provided deployers and returns their outputs",
	func(b operations.Bundle, deps StartBlockchainsOpDeps, input StartBlockchainsOpInput) (StartBlockchainsOpOutput, error) {
		output, err := blockchains.Start(
			input.Inputs,
			deps.Deployers,
		)

		if err != nil {
			return StartBlockchainsOpOutput{}, err
		}

		return StartBlockchainsOpOutput{
			DeployedBlockchains: output,
			OpBundle:            b,
		}, nil
	},
)
